package chattoplan

import (
	"context"
	"omniport-api/internal/helper"
	"omniport-api/internal/modules/chat"
	"omniport-api/internal/modules/plan/vesselschedule"
)

type DirectVesselProvider struct {
	planService vesselschedule.VesselScheduleService
}

func NewDirectVesselProvider(planService vesselschedule.VesselScheduleService) *DirectVesselProvider {
	return &DirectVesselProvider{planService: planService}
}

func (p *DirectVesselProvider) SetPlanService(service vesselschedule.VesselScheduleService) {
	p.planService = service
}

func (p *DirectVesselProvider) GetScheduleTopicByID(ctx context.Context, scheduleID uint64) (*chat.ScheduleChat, error) {
	if p.planService == nil {
		return nil, chat.ErrInternal
	}
	schedule, err := p.planService.FindByID(ctx, scheduleID)
	if err != nil {
		return nil, err
	}

	return mapToChatSchedule(schedule), nil
}

func (p *DirectVesselProvider) GetScheduleByTopicID(ctx context.Context, topicID int64) (*chat.ScheduleChat, error) {
	if p.planService == nil {
		return nil, chat.ErrInternal
	}
	schedule, err := p.planService.FindByTopicID(ctx, topicID)
	if err != nil {
		return nil, err
	}

	return mapToChatSchedule(schedule), nil
}

func (p *DirectVesselProvider) ListScheduleTopics(ctx context.Context, status string) ([]chat.ScheduleChat, error) {
	if p.planService == nil {
		return nil, chat.ErrInternal
	}

	// Just pass the filter to the service, no business logic here
	res, _, err := p.planService.Search(ctx, helper.PaginationQuery{
		Filters: map[string]string{
			"status": status,
		},
		Limit: 100,
	})
	if err != nil {
		return nil, err
	}

	out := make([]chat.ScheduleChat, 0)
	for _, item := range res {
		// Only include items that have a topic ID
		if item.VesselSchedule.TelegramTopicID != nil && *item.VesselSchedule.TelegramTopicID != "" {
			out = append(out, *mapToChatSchedule(&item.VesselSchedule))
		}
	}
	return out, nil
}

func (p *DirectVesselProvider) UpdateScheduleTopicStatus(ctx context.Context, scheduleID uint64, status int) error {
	if p.planService == nil {
		return chat.ErrInternal
	}
	schedule, err := p.planService.FindByID(ctx, scheduleID)
	if err != nil {
		return err
	}
	if schedule.ScheduleCode == nil {
		return nil
	}
	return p.planService.UpdateStatus(ctx, *schedule.ScheduleCode, status, "SYSTEM")
}

// mapToChatSchedule is a helper for data translation (Dumb Mapping)
func mapToChatSchedule(s *vesselschedule.VesselSchedule) *chat.ScheduleChat {
	if s == nil {
		return nil
	}
	vnum := s.VoyageNumber
	return &chat.ScheduleChat{
		ID:                s.ID,
		ScheduleCode:      s.ScheduleCode,
		VesselName:        s.VesselName,
		VoyageNumber:      &vnum,
		ETA:               s.ETA,
		TelegramTopicID:   s.TelegramTopicID,
		TelegramTopicName: s.TelegramTopicName,
		LastUpdatedDate:   s.LastUpdatedDate,
	}
}
