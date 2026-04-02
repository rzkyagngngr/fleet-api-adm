package reference

import (
	"net/http"

	"omniport-api/internal/helper"
	"omniport-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

type ReferenceHandler interface {
	GetAllReferences(c *gin.Context)
	GetReferenceDetail(c *gin.Context)
	SaveReference(c *gin.Context)
	DeleteReference(c *gin.Context)
}

type referenceHandler struct{ service ReferenceService }

func NewReferenceHandler(service ReferenceService) ReferenceHandler {
	return &referenceHandler{service: service}
}

// GetAllReferences godoc
// @Summary Get all references
// @Description Retrieve all reference headers
// @Tags master-references
// @Produce json
// @Security BearerAuth
// @Success 200 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/references [get]
func (h *referenceHandler) GetAllReferences(c *gin.Context) {
	refs, err := h.service.GetAllHeaders()
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "references retrieved successfully", refs)
}

// GetReferenceDetail godoc
// @Summary Get reference detail
// @Description Retrieve reference header and details by id
// @Tags master-references
// @Produce json
// @Security BearerAuth
// @Param id path int true "Reference ID"
// @Success 200 {object} helper.Response
// @Failure 404 {object} helper.Response
// @Router /master/references/{id} [get]
func (h *referenceHandler) GetReferenceDetail(c *gin.Context) {
	id := c.Param("id")
	ref, err := h.service.GetHeaderWithDetails(id)
	if err != nil {
		helper.ErrorResponse(c, http.StatusNotFound, "reference not found")
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "reference detail retrieved successfully", ref)
}

// SaveReference godoc
// @Summary Save reference
// @Description Create or update reference header with details
// @Tags master-references
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body reference.PosmReference true "Reference payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/references [post]
func (h *referenceHandler) SaveReference(c *gin.Context) {
	var input PosmReference
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	employeeID, _ := c.Get(middleware.EmployeeIDKey)
	branchCode, _ := c.Get(middleware.BranchCodeKey)
	terminalCode, _ := c.Get(middleware.TerminalCodeKey)

	input.CreationBy = employeeID.(string)
	input.LastUpdatedBy = employeeID.(string)
	input.BranchCode = *branchCode.(*int64)
	input.TerminalCode = terminalCode.(*int64)

	if err := h.service.SaveReference(input, input.Details); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "reference saved successfully", input)
}

// DeleteReference godoc
// @Summary Delete reference
// @Description Delete reference by id
// @Tags master-references
// @Produce json
// @Security BearerAuth
// @Param id path int true "Reference ID"
// @Success 200 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/references/{id} [delete]
func (h *referenceHandler) DeleteReference(c *gin.Context) {
	if err := h.service.DeleteReference(c.Param("id")); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "reference deleted successfully", nil)
}
