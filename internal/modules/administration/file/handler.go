package file

import (
	"net/http"
	"omniport-api/internal/helper"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FileHandler struct {
	service FileService
}

func NewFileHandler(service FileService) *FileHandler {
	return &FileHandler{service: service}
}

// GetUploadSignature godoc
// @Summary      Create upload signature
// @Description  Generate an upload signature for a file
// @Tags         storage
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload body file.UploadSignatureRequest true "Upload signature payload"
// @Success      200 {object} helper.Response
// @Failure      400 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /storage/upload-signature [post]
func (h *FileHandler) GetUploadSignature(c *gin.Context) {
	var req UploadSignatureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID := c.GetString("user_id")
	res, err := h.service.GetUploadSignature(c.Request.Context(), req, userID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Upload signature generated", res)
}

// CommitUpload godoc
// @Summary      Commit upload
// @Description  Mark an uploaded file as committed
// @Tags         storage
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "File ID"
// @Success      200 {object} helper.Response
// @Failure      400 {object} helper.Response
// @Failure      500 {object} helper.Response
// @Router       /storage/commit/{id} [post]
func (h *FileHandler) CommitUpload(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid file ID")
		return
	}

	if err := h.service.CommitUpload(c.Request.Context(), id); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "File committed successfully", nil)
}

// GetFileDetail godoc
// @Summary      Get file detail
// @Description  Get uploaded file detail by id
// @Tags         storage
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "File ID"
// @Success      200 {object} helper.Response
// @Failure      400 {object} helper.Response
// @Failure      404 {object} helper.Response
// @Router       /storage/file/{id} [get]
func (h *FileHandler) GetFileDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Invalid file ID")
		return
	}

	res, err := h.service.GetFileDetail(c.Request.Context(), id)
	if err != nil {
		helper.ErrorResponse(c, http.StatusNotFound, "File not found")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "File detail retrieved", res)
}
