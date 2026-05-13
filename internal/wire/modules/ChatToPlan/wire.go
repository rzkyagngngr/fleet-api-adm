package chattoplan

import (
	"context"
	"strings"

	"omniport-api/internal/modules/chat"
)


type chatToPlanWire struct {
	mode         string
	monolith     chat.VesselScheduleProvider
	microservice chat.VesselScheduleProvider
}

func NewChatToPlanWire(mode string, monolith chat.VesselScheduleProvider, microservice chat.VesselScheduleProvider) chat.VesselScheduleProvider {
	return &chatToPlanWire{
		mode:         strings.ToLower(strings.TrimSpace(mode)),
		monolith:     monolith,
		microservice: microservice,
	}
}

func (w *chatToPlanWire) GetScheduleTopicByID(ctx context.Context, scheduleID uint64) (*chat.ScheduleChat, error) {
	return w.getProvider().GetScheduleTopicByID(ctx, scheduleID)
}

func (w *chatToPlanWire) GetScheduleByTopicID(ctx context.Context, topicID int64) (*chat.ScheduleChat, error) {
	return w.getProvider().GetScheduleByTopicID(ctx, topicID)
}

func (w *chatToPlanWire) ListScheduleTopics(ctx context.Context, status string) ([]chat.ScheduleChat, error) {
	return w.getProvider().ListScheduleTopics(ctx, status)
}

func (w *chatToPlanWire) UpdateScheduleTopicStatus(ctx context.Context, scheduleID uint64, status int) error {
	return w.getProvider().UpdateScheduleTopicStatus(ctx, scheduleID, status)
}

func (w *chatToPlanWire) getProvider() chat.VesselScheduleProvider {
	if w.mode != "monolith" && w.microservice != nil {
		return w.microservice
	}
	return w.monolith
}
