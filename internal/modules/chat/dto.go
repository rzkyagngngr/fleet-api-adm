package chat

import (
	"strconv"
	"strings"
	"time"
)

type CreateVesselCallRequest struct {
	ScheduleID     uint64 `json:"schedule_id" binding:"required"`
	TelegramChatID int64  `json:"telegram_chat_id" binding:"required"`
	TopicName      string `json:"topic_name" binding:"omitempty,max=200"`
}

type ListVesselCallsQuery struct {
	ScheduleID *uint64
	Status     string
}

type RenameVesselCallRequest struct {
	TopicName string `json:"topic_name" binding:"required,max=200"`
}

type ArchiveVesselCallRequest struct {
	Reason string `json:"reason" binding:"omitempty,max=500"`
}


type AddParticipantRequest struct {
	ScheduleID     uint64  `json:"schedule_id" binding:"required"`
	InternalUserID *uint64 `json:"internal_user_id"`
	TelegramUserID int64   `json:"telegram_user_id" binding:"required"`
	Role           string  `json:"role" binding:"required"`
}

type SendMessageRequest struct {
	ScheduleID uint64 `json:"schedule_id" binding:"required"`
	Text       string `json:"text" binding:"required"`
}

type InviteRequest struct {
	ScheduleID  uint64 `json:"schedule_id" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type ListMessagesQuery struct {
	Limit  int `form:"limit,default=100"`
	Offset int `form:"offset,default=0"`
}


type InviteResponse struct {
	ScheduleID  uint64 `json:"schedule_id"`
	PhoneNumber string `json:"phone_number"`
	InviteLink  string `json:"invite_link"`
	Note        string `json:"note"`
}

type VesselCallResponse struct {
	ScheduleID        uint64     `json:"schedule_id"`
	ScheduleCode      string     `json:"schedule_code"`
	VesselName        string     `json:"vessel_name"`
	VoyageNumber      string     `json:"voyage_number"`
	ETA               *time.Time `json:"eta"`
	TelegramTopicID   int64      `json:"telegram_topic_id"`
	TelegramTopicName string     `json:"telegram_topic_name"`
	LastUpdatedDate   *time.Time `json:"last_updated_date"`
}

func ToVesselCallResponse(input ScheduleChat) VesselCallResponse {
	return VesselCallResponse{
		ScheduleID:        input.ID,
		ScheduleCode:      stringValue(input.ScheduleCode),
		VesselName:        stringValue(input.VesselName),
		VoyageNumber:      stringValue(input.VoyageNumber),
		ETA:               input.ETA,
		TelegramTopicID:   parseTopicID(input.TelegramTopicID),
		TelegramTopicName: stringValue(input.TelegramTopicName),
		LastUpdatedDate:   input.LastUpdatedDate,
	}
}

func ParseListVesselCallsQuery(scheduleIDRaw, statusRaw string) (*ListVesselCallsQuery, error) {
	query := &ListVesselCallsQuery{}

	status := strings.TrimSpace(strings.ToLower(statusRaw))
	if status != "" {
		query.Status = status
	}

	if scheduleIDRaw == "" {
		return query, nil
	}

	parsed, err := parseUint64(scheduleIDRaw)
	if err != nil {
		return nil, err
	}
	query.ScheduleID = &parsed
	return query, nil
}

func stringValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func parseTopicID(v *string) int64 {
	if v == nil || *v == "" {
		return 0
	}
	id, _ := strconv.ParseInt(*v, 10, 64)
	return id
}

