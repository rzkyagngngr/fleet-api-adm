package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"omniport-api/internal/helper"
	"strconv"
	"strings"
	"time"
)

type VesselScheduleProvider interface {
	GetScheduleTopicByID(ctx context.Context, scheduleID uint64) (*ScheduleChat, error)
	GetScheduleByTopicID(ctx context.Context, topicID int64) (*ScheduleChat, error)
	ListScheduleTopics(ctx context.Context, status string) ([]ScheduleChat, error)
	UpdateScheduleTopicStatus(ctx context.Context, scheduleID uint64, status int) error
}

type Service interface {
	HandleWebhookUpdate(ctx context.Context, update *TelegramUpdate) error
	CreateVesselCall(ctx context.Context, req *CreateVesselCallRequest, createdBy string) (*VesselCallResponse, error)
	ListVesselCalls(ctx context.Context, status string) ([]VesselCallResponse, error)
	SuspendVesselCall(ctx context.Context, scheduleID uint64, updatedBy string) error
	ContinueVesselCall(ctx context.Context, scheduleID uint64, updatedBy string) error
	RenameVesselCall(ctx context.Context, scheduleID uint64, topicName string, updatedBy string) error
	ArchiveVesselCall(ctx context.Context, scheduleID uint64, reason string, updatedBy string) error

	// Participant management
	AddParticipant(ctx context.Context, req *AddParticipantRequest, createdBy string) (*CallParticipant, error)
	ListParticipants(ctx context.Context, scheduleID uint64) ([]CallParticipant, error)
	RemoveParticipant(ctx context.Context, scheduleID uint64, telegramUserID int64) error

	// Messaging
	SendMessage(ctx context.Context, scheduleID uint64, text string) error
	ListMessages(ctx context.Context, scheduleID uint64, query *ListMessagesQuery) ([]ChatMessage, int64, error)
	InviteByScheduleID(ctx context.Context, scheduleID uint64) (*InviteResponse, error)
}

type service struct {
	repo                 Repository
	vesselProvider       VesselScheduleProvider
	bot                  TelegramClient
	s3                   helper.StorageProvider
	s3Bucket             string
	telegramParentChatID int64
}

func NewService(repo Repository, bot TelegramClient, s3 helper.StorageProvider, s3Bucket string, telegramParentChatID int64, vesselProvider VesselScheduleProvider) Service {
	return &service{
		repo:                 repo,
		vesselProvider:       vesselProvider,
		bot:                  bot,
		s3:                   s3,
		s3Bucket:             s3Bucket,
		telegramParentChatID: telegramParentChatID,
	}
}

func (s *service) HandleWebhookUpdate(ctx context.Context, update *TelegramUpdate) error {
	if update.Message == nil {
		return nil
	}

	msg := update.Message
	if msg.Chat.Type != "supergroup" || msg.MessageThreadID == 0 {
		return nil
	}

	call, err := s.vesselProvider.GetScheduleByTopicID(ctx, int64(msg.MessageThreadID))
	if err != nil {
		return fmt.Errorf("failed to resolve call_id for topic %d: %w", msg.MessageThreadID, err)
	}

	senderName := strings.TrimSpace(msg.From.FirstName + " " + msg.From.LastName)
	if senderName == "" {
		senderName = msg.From.Username
	}

	raw, _ := json.Marshal(update)
	dbMsg := &ChatMessage{
		CallID:            call.ID,
		TelegramMessageID: int64(msg.MessageID),
		TelegramChatID:    fmt.Sprintf("%d", msg.Chat.ID),
		TelegramTopicID:   fmt.Sprintf("%d", msg.MessageThreadID),
		SenderID:          msg.From.ID,
		SenderName:        senderName,
		Text:              msg.Text,
		MessageTimestamp:  time.Unix(int64(msg.Date), 0).UTC(),
		RawPayload:        raw,
		IsAuthorized:      true,
	}

	// Media handling
	fileID, mimeType, fileName, fileSize := resolveFileFromMessage(msg)
	if fileID != "" {
		if fileSize <= 10*1024*1024 {
			fileURL, err := s.bot.GetFileURL(fileID)
			if err == nil {
				data, err := s.bot.DownloadFile(fileURL)
				if err == nil {
					key := fmt.Sprintf("chat/%d/%d/%s", msg.Chat.ID, msg.MessageThreadID, fileName)
					if err := s.s3.UploadObject(ctx, s.s3Bucket, key, mimeType, data); err == nil {
						s3URL, _ := s.s3.GeneratePresignedGetURL(ctx, s.s3Bucket, key, 7*24*time.Hour)
						dbMsg.Attatchment, _ = json.Marshal(map[string]interface{}{
							"file_id":     fileID,
							"file_name":   fileName,
							"mime_type":   mimeType,
							"file_size":   fileSize,
							"preview_url": s3URL,
						})
					}
				}
			}
		}
	}

	return s.repo.SaveChatMessage(ctx, dbMsg)
}

func (s *service) CreateVesselCall(ctx context.Context, req *CreateVesselCallRequest, createdBy string) (*VesselCallResponse, error) {
	topicID, err := s.bot.CreateForumTopic(req.TelegramChatID, req.TopicName)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &VesselCallResponse{
		ScheduleID:        req.ScheduleID,
		TelegramTopicID:   topicID,
		TelegramTopicName: req.TopicName,
		LastUpdatedDate:   &now,
	}, nil
}

func (s *service) ListVesselCalls(ctx context.Context, status string) ([]VesselCallResponse, error) {
	rows, err := s.vesselProvider.ListScheduleTopics(ctx, status)
	if err != nil {
		return nil, err
	}

	res := make([]VesselCallResponse, 0, len(rows))
	for _, row := range rows {
		res = append(res, ToVesselCallResponse(row))
	}
	return res, nil
}

func (s *service) SuspendVesselCall(ctx context.Context, scheduleID uint64, updatedBy string) error {
	call, err := s.vesselProvider.GetScheduleTopicByID(ctx, scheduleID)
	if err != nil {
		return err
	}
	parentChatID, _ := s.resolveParentChatID(0)
	if call.TelegramTopicID != nil {
		tID, _ := strconv.ParseInt(*call.TelegramTopicID, 10, 64)
		s.bot.CloseForumTopic(parentChatID, tID)
	}
	return s.vesselProvider.UpdateScheduleTopicStatus(ctx, scheduleID, 2)
}

func (s *service) ContinueVesselCall(ctx context.Context, scheduleID uint64, updatedBy string) error {
	call, err := s.vesselProvider.GetScheduleTopicByID(ctx, scheduleID)
	if err != nil {
		return err
	}
	parentChatID, _ := s.resolveParentChatID(0)
	if call.TelegramTopicID != nil {
		tID, _ := strconv.ParseInt(*call.TelegramTopicID, 10, 64)
		s.bot.ReopenForumTopic(parentChatID, tID)
	}
	return s.vesselProvider.UpdateScheduleTopicStatus(ctx, scheduleID, 1)
}

func (s *service) RenameVesselCall(ctx context.Context, scheduleID uint64, topicName string, updatedBy string) error {
	call, err := s.vesselProvider.GetScheduleTopicByID(ctx, scheduleID)
	if err != nil {
		return err
	}
	parentChatID, _ := s.resolveParentChatID(0)
	if call.TelegramTopicID != nil {
		tID, _ := strconv.ParseInt(*call.TelegramTopicID, 10, 64)
		s.bot.EditForumTopic(parentChatID, tID, strings.TrimSpace(topicName))
	}
	return nil
}

func (s *service) ArchiveVesselCall(ctx context.Context, scheduleID uint64, reason string, updatedBy string) error {
	call, err := s.vesselProvider.GetScheduleTopicByID(ctx, scheduleID)
	if err != nil {
		return err
	}
	parentChatID, _ := s.resolveParentChatID(0)
	if call.TelegramTopicID != nil {
		tID, _ := strconv.ParseInt(*call.TelegramTopicID, 10, 64)
		s.bot.CloseForumTopic(parentChatID, tID)
	}
	return s.vesselProvider.UpdateScheduleTopicStatus(ctx, scheduleID, 9)
}

func (s *service) AddParticipant(ctx context.Context, req *AddParticipantRequest, createdBy string) (*CallParticipant, error) {
	call, err := s.vesselProvider.GetScheduleTopicByID(ctx, req.ScheduleID)
	if err != nil {
		return nil, err
	}
	parentChatID, _ := s.resolveParentChatID(0)
	if call.TelegramTopicID == nil {
		return nil, fmt.Errorf("topic not initialized")
	}

	row := &CallParticipant{
		CallID:          call.ID,
		TelegramChatID:  fmt.Sprintf("%d", parentChatID),
		TelegramTopicID: *call.TelegramTopicID,
		InternalUserID:  req.InternalUserID,
		TelegramUserID:  req.TelegramUserID,
		Role:           strings.ToLower(req.Role),
		Status:         1,
		CreationBy:     createdBy,
	}

	if err := s.repo.CreateParticipant(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *service) ListParticipants(ctx context.Context, scheduleID uint64) ([]CallParticipant, error) {
	call, err := s.vesselProvider.GetScheduleTopicByID(ctx, scheduleID)
	if err != nil {
		return nil, err
	}
	parentChatID, _ := s.resolveParentChatID(0)
	if call.TelegramTopicID == nil {
		return nil, nil
	}
	tID, _ := strconv.ParseInt(*call.TelegramTopicID, 10, 64)
	return s.repo.ListParticipants(ctx, parentChatID, tID)
}

func (s *service) RemoveParticipant(ctx context.Context, scheduleID uint64, telegramUserID int64) error {
	call, err := s.vesselProvider.GetScheduleTopicByID(ctx, scheduleID)
	if err != nil {
		return err
	}
	parentChatID, _ := s.resolveParentChatID(0)
	if call.TelegramTopicID == nil {
		return nil
	}
	tID, _ := strconv.ParseInt(*call.TelegramTopicID, 10, 64)
	return s.repo.DeleteParticipant(ctx, parentChatID, tID, telegramUserID)
}

func (s *service) SendMessage(ctx context.Context, scheduleID uint64, text string) error {
	call, err := s.vesselProvider.GetScheduleTopicByID(ctx, scheduleID)
	if err != nil {
		return err
	}
	parentChatID, err := s.resolveParentChatID(0)
	if err != nil {
		return err
	}
	if call.TelegramTopicID == nil {
		return fmt.Errorf("telegram topic not initialized")
	}

	tID, _ := strconv.ParseInt(*call.TelegramTopicID, 10, 64)
	msgID, err := s.bot.SendMessage(parentChatID, tID, text)
	if err != nil {
		return err
	}

	outgoing := &ChatMessage{
		CallID:            call.ID,
		TelegramMessageID: msgID,
		TelegramChatID:    fmt.Sprintf("%d", parentChatID),
		TelegramTopicID:   *call.TelegramTopicID,
		SenderID:          0,
		SenderName:        "OMNIPORT SYSTEM",
		Text:              text,
		MessageTimestamp:  time.Now().UTC(),
		IsAuthorized:      true,
	}
	return s.repo.SaveChatMessage(ctx, outgoing)
}

func (s *service) ListMessages(ctx context.Context, scheduleID uint64, query *ListMessagesQuery) ([]ChatMessage, int64, error) {
	call, err := s.vesselProvider.GetScheduleTopicByID(ctx, scheduleID)
	if err != nil {
		return nil, 0, err
	}
	parentChatID, _ := s.resolveParentChatID(0)
	if call.TelegramTopicID == nil {
		return nil, 0, nil
	}

	limit, offset := 50, 0
	if query != nil {
		if query.Limit > 0 { limit = query.Limit }
		if query.Offset >= 0 { offset = query.Offset }
	}
	tID, _ := strconv.ParseInt(*call.TelegramTopicID, 10, 64)
	return s.repo.ListMessages(ctx, parentChatID, tID, limit, offset)
}

func (s *service) InviteByScheduleID(ctx context.Context, scheduleID uint64) (*InviteResponse, error) {
	call, err := s.vesselProvider.GetScheduleTopicByID(ctx, scheduleID)
	if err != nil {
		return nil, err
	}
	if call.TelegramTopicID == nil {
		return nil, fmt.Errorf("telegram topic not initialized for this schedule")
	}

	// In a real scenario, we would use the emergency contact phone number
	// and possibly use a WhatsApp/SMS gateway to send the link.
	// For now, we generate the deep link to the specific topic.
	
	parentChatID, _ := s.resolveParentChatID(0)
	
	// Format: https://t.me/c/CHAT_ID/TOPIC_ID
	// Note: For private groups, CHAT_ID in links usually drops the '-100' prefix
	cleanChatID := strings.TrimPrefix(fmt.Sprintf("%d", parentChatID), "-100")
	inviteLink := fmt.Sprintf("https://t.me/c/%s/%s", cleanChatID, *call.TelegramTopicID)

	return &InviteResponse{
		ScheduleID:  scheduleID,
		InviteLink:  inviteLink,
		Note:        "Invite link generated for topic. Share this with the emergency contact.",
	}, nil
}

func (s *service) resolveParentChatID(reqID int64) (int64, error) {
	if reqID != 0 { return reqID, nil }
	if s.telegramParentChatID == 0 { return 0, fmt.Errorf("parent chat not configured") }
	return s.telegramParentChatID, nil
}

func resolveFileFromMessage(msg *TelegramMessage) (fileID, mimeType, fileName string, fileSize int64) {
	if len(msg.Photo) > 0 {
		p := msg.Photo[len(msg.Photo)-1]
		return p.FileID, "image/jpeg", fmt.Sprintf("photo_%s.jpg", p.FileID), int64(p.FileSize)
	}
	if msg.Document != nil {
		return msg.Document.FileID, msg.Document.MimeType, msg.Document.FileName, int64(msg.Document.FileSize)
	}
	return "", "", "", 0
}
