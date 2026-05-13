package plantochat

import (
	"context"
	"strings"

	"omniport-api/internal/modules/plan/vesselschedule"
)

type planToChatWire struct {
	mode         string
	monolith     vesselschedule.ScheduleChatInitializer
	microservice vesselschedule.ScheduleChatInitializer
}

func NewPlanToChatWire(mode string, monolith vesselschedule.ScheduleChatInitializer, microservice vesselschedule.ScheduleChatInitializer) vesselschedule.ScheduleChatInitializer {
	return &planToChatWire{
		mode:         strings.ToLower(strings.TrimSpace(mode)),
		monolith:     monolith,
		microservice: microservice,
	}
}

func (w *planToChatWire) InitScheduleChat(ctx context.Context, req vesselschedule.ChatInitRequest) (*int64, *string, error) {
	if w == nil || req.TelegramChatID == 0 {
		return nil, nil, nil
	}

	initializer := w.monolith
	if w.mode != "monolith" && w.microservice != nil {
		initializer = w.microservice
	}

	if initializer == nil {
		return nil, nil, nil
	}

	return initializer.InitScheduleChat(ctx, req)
}
