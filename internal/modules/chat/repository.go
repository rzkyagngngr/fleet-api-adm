package chat

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	SaveChatMessage(ctx context.Context, message *ChatMessage) error
	ListMessages(ctx context.Context, chatID, topicID int64, limit int, offset int) ([]ChatMessage, int64, error)

	// Participant management
	CreateParticipant(ctx context.Context, participant *CallParticipant) error
	ListParticipants(ctx context.Context, chatID, topicID int64) ([]CallParticipant, error)
	DeleteParticipant(ctx context.Context, chatID, topicID int64, telegramUserID int64) error
	IsParticipantAuthorized(ctx context.Context, chatID, topicID int64, telegramUserID int64) (bool, error)
}

type repository struct {
	chatDB *gorm.DB
}

func NewRepository(chatDB *gorm.DB) Repository {
	if chatDB == nil {
		panic("chat database connection is nil")
	}
	return &repository{
		chatDB: chatDB,
	}
}

func (r *repository) SaveChatMessage(ctx context.Context, message *ChatMessage) error {
	query := `
		INSERT INTO chnl.chat_messages (
			call_id, telegram_message_id, telegram_chat_id, telegram_topic_id, 
			sender_id, sender_name, sender_name_local, text, attatchment, 
			message_timestamp, raw_payload, is_authorized, creation_date
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP
		) ON CONFLICT DO NOTHING`

	return r.chatDB.WithContext(ctx).Exec(query,
		message.CallID, message.TelegramMessageID, fmt.Sprintf("%d", message.TelegramChatID), fmt.Sprintf("%d", message.TelegramTopicID),
		message.SenderID, message.SenderName, message.SenderNameLocal, message.Text, message.Attatchment,
		message.MessageTimestamp, message.RawPayload, message.IsAuthorized,
	).Error
}

func (r *repository) ListMessages(ctx context.Context, chatID, topicID int64, limit int, offset int) ([]ChatMessage, int64, error) {
	sChatID := fmt.Sprintf("%d", chatID)
	sTopicID := fmt.Sprintf("%d", topicID)

	var total int64
	countQuery := `SELECT count(*) FROM chnl.chat_messages WHERE telegram_chat_id = ? AND telegram_topic_id = ?`
	if err := r.chatDB.WithContext(ctx).Raw(countQuery, sChatID, sTopicID).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []ChatMessage
	selectQuery := `
		SELECT * FROM chnl.chat_messages 
		WHERE telegram_chat_id = ? AND telegram_topic_id = ? 
		ORDER BY message_timestamp ASC 
		LIMIT ? OFFSET ?`

	err := r.chatDB.WithContext(ctx).Raw(selectQuery, sChatID, sTopicID, limit, offset).Scan(&rows).Error
	return rows, total, err
}

func (r *repository) CreateParticipant(ctx context.Context, p *CallParticipant) error {
	query := `
		INSERT INTO chnl.chat_participants (
			call_id, telegram_chat_id, telegram_topic_id, internal_user_id, 
			telegram_user_id, role, status, creation_date, creation_by
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, ?
		) ON CONFLICT (telegram_chat_id, telegram_topic_id, telegram_user_id) 
		DO UPDATE SET 
			status = EXCLUDED.status,
			role = EXCLUDED.role,
			last_updated_date = CURRENT_TIMESTAMP`

	return r.chatDB.WithContext(ctx).Exec(query,
		p.CallID, fmt.Sprintf("%d", p.TelegramChatID), fmt.Sprintf("%d", p.TelegramTopicID), p.InternalUserID,
		p.TelegramUserID, p.Role, p.Status, p.CreationBy,
	).Error
}

func (r *repository) ListParticipants(ctx context.Context, chatID, topicID int64) ([]CallParticipant, error) {
	sChatID := fmt.Sprintf("%d", chatID)
	sTopicID := fmt.Sprintf("%d", topicID)

	var rows []CallParticipant
	query := `SELECT * FROM chnl.chat_participants WHERE telegram_chat_id = ? AND telegram_topic_id = ? AND status = 1`
	err := r.chatDB.WithContext(ctx).Raw(query, sChatID, sTopicID).Scan(&rows).Error
	return rows, err
}

func (r *repository) DeleteParticipant(ctx context.Context, chatID, topicID int64, telegramUserID int64) error {
	sChatID := fmt.Sprintf("%d", chatID)
	sTopicID := fmt.Sprintf("%d", topicID)

	query := `
		UPDATE chnl.chat_participants 
		SET status = 0, last_updated_date = CURRENT_TIMESTAMP 
		WHERE telegram_chat_id = ? AND telegram_topic_id = ? AND telegram_user_id = ?`

	return r.chatDB.WithContext(ctx).Exec(query, sChatID, sTopicID, telegramUserID).Error
}

func (r *repository) IsParticipantAuthorized(ctx context.Context, chatID, topicID int64, telegramUserID int64) (bool, error) {
	sChatID := fmt.Sprintf("%d", chatID)
	sTopicID := fmt.Sprintf("%d", topicID)

	var count int64
	query := `
		SELECT count(*) FROM chnl.chat_participants 
		WHERE telegram_chat_id = ? AND telegram_topic_id = ? AND telegram_user_id = ? AND status = 1`
	
	err := r.chatDB.WithContext(ctx).Raw(query, sChatID, sTopicID, telegramUserID).Scan(&count).Error
	return count > 0, err
}
