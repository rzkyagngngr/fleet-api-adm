package company

import (
	"net/http"
	"strconv"

	"omniport-api/internal/helper"
	"omniport-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

type CompanyHandler struct {
	service CompanyService
}

func NewCompanyHandler(service CompanyService) *CompanyHandler {
	return &CompanyHandler{service: service}
}

// Search godoc
// @Summary Search companies
// @Description Search and filter companies
// @Tags master-company
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body company.SearchCompaniesRequest true "Search payload"
// @Success 200 {object} helper.MetaResponse
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/companies/search [post]
func (h *CompanyHandler) Search(c *gin.Context) {
	var input SearchCompaniesRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	rows, meta, err := h.service.SearchCompanies(c.Request.Context(), input.ToPaginationQuery())
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.MetaSuccessResponse(c, http.StatusOK, "success", rows, meta)
}

// Create godoc
// @Summary Create company
// @Description Create a new company
// @Tags master-company
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body company.CompanyRequest true "Company payload"
// @Success 201 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/companies [post]
func (h *CompanyHandler) Create(c *gin.Context) {
	var input CompanyRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	employeeID, _ := c.Get(middleware.EmployeeIDKey)

	if err := h.service.CreateCompany(c.Request.Context(), &input, employeeID.(string)); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusCreated, "company created successfully", nil)
}

// Update godoc
// @Summary Update company
// @Description Update company by id
// @Tags master-company
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query int true "Company ID"
// @Param payload body company.CompanyRequest true "Company payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/companies [put]
func (h *CompanyHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid company id")
		return
	}

	var input CompanyRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ValidationErrorResponse(c, err)
		return
	}

	employeeID, _ := c.Get(middleware.EmployeeIDKey)

	if err := h.service.UpdateCompany(c.Request.Context(), id, &input, employeeID.(string)); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "company updated successfully", nil)
}

// GetByID godoc
// @Summary Get company by id
// @Description Get company by id
// @Tags master-company
// @Produce json
// @Security BearerAuth
// @Param id query int true "Company ID"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/companies [get]
func (h *CompanyHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid company id")
		return
	}

	res, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "success", res)
}

// Delete godoc
// @Summary Delete company
// @Description Delete company by id
// @Tags master-company
// @Produce json
// @Security BearerAuth
// @Param id query int true "Company ID"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/companies [delete]
func (h *CompanyHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid company id")
		return
	}

	if err := h.service.DeleteCompany(c.Request.Context(), id); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "company deleted successfully", nil)
}
