package chat

import (
	"net/http"
	"strings"

	"omniport-api/internal/helper"
	"omniport-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	service Service
}

func NewChatHandler(service Service) *ChatHandler {
	return &ChatHandler{service: service}
}

func (h *ChatHandler) CreateVesselCall(c *gin.Context) {
	var req CreateVesselCallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	row, err := h.service.CreateVesselCall(c.Request.Context(), &req, actorFromContext(c))
	if err != nil {
		writeChatError(c, err, "failed to create schedule topic")
		return
	}

	helper.SuccessResponse(c, http.StatusCreated, "schedule topic created successfully", row)
}

func (h *ChatHandler) ListVesselCalls(c *gin.Context) {
	status := c.Query("status")
	rows, err := h.service.ListVesselCalls(c.Request.Context(), status)
	if err != nil {
		writeChatError(c, err, "failed to get schedule topics")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "schedule topics retrieved successfully", rows)
}


func (h *ChatHandler) SuspendVesselCall(c *gin.Context) {
	scheduleID, err := parseUint64(c.Param("id"))
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid schedule id")
		return
	}

	if err := h.service.SuspendVesselCall(c.Request.Context(), scheduleID, actorFromContext(c)); err != nil {
		writeChatError(c, err, "failed to suspend schedule topic")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "schedule topic suspended successfully", nil)
}

func (h *ChatHandler) ContinueVesselCall(c *gin.Context) {
	scheduleID, err := parseUint64(c.Param("id"))
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid schedule id")
		return
	}

	if err := h.service.ContinueVesselCall(c.Request.Context(), scheduleID, actorFromContext(c)); err != nil {
		writeChatError(c, err, "failed to continue schedule topic")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "schedule topic continued successfully", nil)
}

func (h *ChatHandler) RenameVesselCall(c *gin.Context) {
	scheduleID, err := parseUint64(c.Param("id"))
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid schedule id")
		return
	}

	var req RenameVesselCallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	if err := h.service.RenameVesselCall(c.Request.Context(), scheduleID, req.TopicName, actorFromContext(c)); err != nil {
		writeChatError(c, err, "failed to rename schedule topic")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "schedule topic renamed successfully", nil)
}

func (h *ChatHandler) ListParticipants(c *gin.Context) {
	scheduleID, err := parseUint64(c.Param("id"))
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid schedule id")
		return
	}

	rows, err := h.service.ListParticipants(c.Request.Context(), scheduleID)
	if err != nil {
		writeChatError(c, err, "failed to get participants")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "participants retrieved successfully", rows)
}

func (h *ChatHandler) RemoveParticipant(c *gin.Context) {
	scheduleID, err := parseUint64(c.Param("id"))
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid schedule id")
		return
	}

	telegramUserID, err := parseUint64(c.Param("user_id"))
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid telegram user id")
		return
	}

	if err := h.service.RemoveParticipant(c.Request.Context(), scheduleID, int64(telegramUserID)); err != nil {
		writeChatError(c, err, "failed to remove participant")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "participant removed successfully", nil)
}

func (h *ChatHandler) AddParticipant(c *gin.Context) {
	var req AddParticipantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	row, err := h.service.AddParticipant(c.Request.Context(), &req, actorFromContext(c))
	if err != nil {
		writeChatError(c, err, "failed to add participant")
		return
	}

	helper.SuccessResponse(c, http.StatusCreated, "participant added successfully", row)
}

func (h *ChatHandler) Invite(c *gin.Context) {
	scheduleID, err := parseUint64(c.Param("id"))
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid schedule id")
		return
	}

	row, err := h.service.InviteByScheduleID(c.Request.Context(), scheduleID)
	if err != nil {
		writeChatError(c, err, "failed to generate invite")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "invite generated successfully", row)
}


func (h *ChatHandler) SendMessage(c *gin.Context) {
	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	if err := h.service.SendMessage(c.Request.Context(), req.ScheduleID, strings.TrimSpace(req.Text)); err != nil {
		writeChatError(c, err, "failed to send message")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "message sent successfully", nil)
}

func (h *ChatHandler) HandleTelegramWebhook(c *gin.Context) {
	var update TelegramUpdate
	if err := c.ShouldBindJSON(&update); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	if err := h.service.HandleWebhookUpdate(c.Request.Context(), &update); err != nil {
		writeChatError(c, err, "failed to process telegram webhook")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "telegram webhook processed", nil)
}

func (h *ChatHandler) ListMessages(c *gin.Context) {
	scheduleID, err := parseUint64(c.Param("id"))
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid schedule id")
		return
	}

	var query ListMessagesQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	messages, total, err := h.service.ListMessages(c.Request.Context(), scheduleID, &query)
	if err != nil {
		writeChatError(c, err, "failed to get chat messages")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "chat messages retrieved successfully", gin.H{
		"data":  messages,
		"total": total,
	})
}

func (h *ChatHandler) ArchiveVesselCall(c *gin.Context) {
	scheduleID, err := parseUint64(c.Param("id"))
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid schedule id")
		return
	}

	var req ArchiveVesselCallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Optional body, but bind if present
	}

	if err := h.service.ArchiveVesselCall(c.Request.Context(), scheduleID, req.Reason, actorFromContext(c)); err != nil {
		writeChatError(c, err, "failed to archive schedule topic")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "schedule topic archived successfully", nil)
}

func actorFromContext(c *gin.Context) string {
	actor := strings.TrimSpace(middleware.GetUserEmail(c))
	if actor == "" {
		actor = strings.TrimSpace(c.GetHeader("X-Internal-Actor"))
	}
	if actor == "" {
		return "SYSTEM"
	}
	return actor
}

func writeChatError(c *gin.Context, err error, fallbackMessage string) {
	switch err {
	case ErrScheduleNotFound, ErrVesselCallNotFound, ErrParticipantNotFound:
		helper.ErrorResponse(c, http.StatusNotFound, err.Error())
	case ErrParticipantExists, ErrScheduleTopicExists:
		helper.ErrorResponse(c, http.StatusConflict, err.Error())
	case ErrInvalidCallStatus, ErrInvalidRole:
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
	default:
		helper.ErrorResponse(c, http.StatusInternalServerError, fallbackMessage)
	}
}
