package pelabuhan

import (
	"net/http"
	"omniport-api/internal/helper"
	"omniport-api/internal/middleware"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PortHandler struct {
	service PortService
}

func NewPortHandler(service PortService) *PortHandler {
	return &PortHandler{service: service}
}

// SearchPorts godoc
// @Summary Search ports
// @Description Retrieve ports with server-side pagination, filtering, sorting, and download range support
// @Tags master-ports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body pelabuhan.SearchPortRequest true "Port search payload"
// @Success 200 {object} helper.MetaResponse
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/pelabuhan/search [post]
func (h *PortHandler) SearchPorts(c *gin.Context) {
	var input SearchPortRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	ports, meta, err := h.service.Search(c.Request.Context(), input.ToPaginationQuery())
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to search ports")
		return
	}

	helper.MetaSuccessResponse(c, http.StatusOK, "ports retrieved successfully", ports, meta)
}

// CreatePort godoc
// @Summary Create port
// @Description Create a new port record
// @Tags master-ports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body pelabuhan.PortReq true "Port payload"
// @Success 201 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/pelabuhan [post]
func (h *PortHandler) CreatePort(c *gin.Context) {
	var req PortReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userName := middleware.GetUserEmail(c) // Use Email as user identifier
	if userName == "" {
		userName = "SYSTEM"
	}

	port := Port{
		PortCode:    req.PortCode,
		PortName:    &req.PortName,
		PortCity:    &req.PortCity,
		CountryCode: &req.CountryCode,
		Status:      req.Status,
		CreatedBy:   &userName,
		LastUpdatedBy: userName,
	}

	if err := h.service.Create(c.Request.Context(), &port); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to create port")
		return
	}

	helper.SuccessResponse(c, http.StatusCreated, "port created successfully", port)
}

// UpdatePort godoc
// @Summary Update port
// @Description Update an existing port by id
// @Tags master-ports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Port ID"
// @Param payload body pelabuhan.PortReq true "Port payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/pelabuhan/{id} [put]
func (h *PortHandler) UpdatePort(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid port id")
		return
	}

	var req PortReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userName := middleware.GetUserEmail(c)
	if userName == "" {
		userName = "SYSTEM"
	}

	port := Port{
		PortCode:    req.PortCode,
		PortName:    &req.PortName,
		PortCity:    &req.PortCity,
		CountryCode: &req.CountryCode,
		Status:      req.Status,
		LastUpdatedBy: userName,
	}

	if err := h.service.Update(c.Request.Context(), id, &port); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to update port")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "port updated successfully", port)
}

// DeletePort godoc
// @Summary Delete port
// @Description Delete port by id
// @Tags master-ports
// @Produce json
// @Security BearerAuth
// @Param id path int true "Port ID"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/pelabuhan/{id} [delete]
func (h *PortHandler) DeletePort(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid port id")
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to delete port")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "port deleted successfully", nil)
}
