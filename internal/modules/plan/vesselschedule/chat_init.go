package vesselschedule

import "context"

type ChatInitRequest struct {
	ScheduleID     uint64
	TelegramChatID int64
	TopicName      string
	Actor          string
}

type ScheduleChatInitializer interface {
	InitScheduleChat(ctx context.Context, req ChatInitRequest) (*int64, *string, error)
}
