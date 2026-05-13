package chat

import "time"

const (
	VesselCallStatusActive    = "active"
	VesselCallStatusSuspended = "suspended"

	ParticipantRoleAgent    = "agent"
	ParticipantRolePBM      = "pbm"
	ParticipantRoleOperator = "operator"
)

type ScheduleChat struct {
	ID                uint64     `gorm:"column:id" json:"id"`
	ScheduleCode      *string    `gorm:"column:schedule_code" json:"schedule_code"`
	VesselName        *string    `gorm:"column:vessel_name" json:"vessel_name"`
	VoyageNumber      *string    `gorm:"column:voyage_number" json:"voyage_number"`
	ETA               *time.Time `gorm:"column:eta" json:"eta"`
	TelegramTopicID   *string    `gorm:"column:telegram_topic_id" json:"telegram_topic_id"`

	TelegramTopicName *string    `gorm:"column:telegram_topic_name" json:"telegram_topic_name"`
	LastUpdatedDate   *time.Time `gorm:"column:last_updated_date" json:"last_updated_date"`
}

func (ScheduleChat) TableName() string { return "plan.post_vessel_schedules" }

type CallParticipant struct {
	ID              uint64     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CallID          uint64     `gorm:"column:call_id;not null" json:"call_id"`
	TelegramChatID  string     `gorm:"column:telegram_chat_id;not null" json:"telegram_chat_id"`
	TelegramTopicID string     `gorm:"column:telegram_topic_id;not null" json:"telegram_topic_id"`
	InternalUserID  *uint64    `gorm:"column:internal_user_id" json:"internal_user_id"`
	TelegramUserID  int64      `gorm:"column:telegram_user_id;not null" json:"telegram_user_id"`
	Role            string     `gorm:"column:role;size:20;not null" json:"role"`
	Status          int16      `gorm:"column:status;default:1;not null" json:"status"`
	CreationDate    time.Time  `gorm:"column:creation_date;autoCreateTime" json:"creation_date"`
	CreationBy      string     `gorm:"column:creation_by;size:100" json:"creation_by"`
	LastUpdatedDate *time.Time `gorm:"column:last_updated_date" json:"last_updated_date"`
	LastUpdatedBy   *string    `gorm:"column:last_updated_by;size:100" json:"last_updated_by"`
}

func (CallParticipant) TableName() string { return "chnl.chat_participants" }


type ChatMessage struct {
	ID                uint64    `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CallID            uint64    `gorm:"column:call_id;not null" json:"call_id"`
	TelegramMessageID int64     `gorm:"column:telegram_message_id;not null" json:"telegram_message_id"`
	TelegramChatID    string    `gorm:"column:telegram_chat_id;not null" json:"telegram_chat_id"`
	TelegramTopicID   string    `gorm:"column:telegram_topic_id;not null" json:"telegram_topic_id"`
	SenderID          int64     `gorm:"column:sender_id;not null" json:"sender_id"`
	SenderName        string    `gorm:"column:sender_name" json:"sender_name"`
	SenderNameLocal   string    `gorm:"column:sender_name_local" json:"sender_name_local"`
	Text              string    `gorm:"column:text;type:text" json:"text"`
	Attatchment       []byte    `gorm:"column:attatchment;type:jsonb" json:"attatchment,omitempty"`
	MessageTimestamp  time.Time `gorm:"column:message_timestamp;not null" json:"message_timestamp"`
	RawPayload        []byte    `gorm:"column:raw_payload;type:jsonb" json:"raw_payload,omitempty"`
	IsAuthorized      bool      `gorm:"column:is_authorized;not null;default:true" json:"is_authorized"`
	CreationDate      time.Time `gorm:"column:creation_date;autoCreateTime" json:"creation_date"`
}

func (ChatMessage) TableName() string { return "chnl.chat_messages" }

