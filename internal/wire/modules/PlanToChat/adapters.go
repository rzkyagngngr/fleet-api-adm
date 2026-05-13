package plantochat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"omniport-api/internal/modules/chat"
	"omniport-api/internal/modules/plan/vesselschedule"
	wirehelper "omniport-api/internal/wire/helper"
)

const planToChatInitPath = "/internal/chat/vessel-calls"

type directPlanToChat struct {
	service chat.Service
}

func NewDirectPlanToChatWire(service chat.Service) vesselschedule.ScheduleChatInitializer {
	return &directPlanToChat{service: service}
}

func (w *directPlanToChat) InitScheduleChat(ctx context.Context, req vesselschedule.ChatInitRequest) (*int64, *string, error) {
	if w == nil || w.service == nil {
		return nil, nil, nil
	}
	res, err := w.service.CreateVesselCall(ctx, &chat.CreateVesselCallRequest{
		ScheduleID:     req.ScheduleID,
		TelegramChatID: req.TelegramChatID,
		TopicName:      req.TopicName,
	}, req.Actor)
	if err != nil {
		if err == chat.ErrScheduleTopicExists {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	return &res.TelegramTopicID, &res.TelegramTopicName, nil
}

type restPlanToChat struct {
	client *wirehelper.InternalRESTClient
}

func NewRestPlanToChatWire(baseURL string, token string) vesselschedule.ScheduleChatInitializer {
	return &restPlanToChat{
		client: wirehelper.NewInternalRESTClient(baseURL, token, 20*time.Second),
	}
}

func (w *restPlanToChat) InitScheduleChat(ctx context.Context, req vesselschedule.ChatInitRequest) (*int64, *string, error) {
	if w == nil || !w.client.Enabled() {
		return nil, nil, nil
	}

	resp, err := w.client.PostJSON(ctx, planToChatInitPath, chat.CreateVesselCallRequest{
		ScheduleID:     req.ScheduleID,
		TelegramChatID: req.TelegramChatID,
		TopicName:      req.TopicName,
	}, map[string]string{
		"X-Internal-Actor": req.Actor,
	})
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		return nil, nil, nil
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, nil, fmt.Errorf("chat internal init returned status %d", resp.StatusCode)
	}

	var result struct {
		Data chat.VesselCallResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, nil, err
	}

	return &result.Data.TelegramTopicID, &result.Data.TelegramTopicName, nil
}
