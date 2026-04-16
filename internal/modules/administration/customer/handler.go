package customer

import (
	"errors"
	"net/http"
	"omniport-api/internal/helper"
	"omniport-api/internal/middleware"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	service CustomerService
}

func NewCustomerHandler(service CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

// SearchCustomers godoc
// @Summary Search customers
// @Description Retrieve customers with server-side pagination, filtering, sorting, and download range support
// @Tags master-customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body customer.SearchCustomerRequest true "Customer search payload"
// @Success 200 {object} helper.MetaResponse
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/customer/search [post]
func (h *CustomerHandler) SearchCustomers(c *gin.Context) {
	var input SearchCustomerRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	branchCodeVal, exists := c.Get(middleware.BranchCodeKey)
	if !exists || branchCodeVal == nil {
		helper.ErrorResponse(c, http.StatusUnauthorized, "branch code not found in token")
		return
	}
	terminalCodeVal, exists := c.Get(middleware.TerminalCodeKey)
	if !exists || terminalCodeVal == nil {
		helper.ErrorResponse(c, http.StatusUnauthorized, "terminal code not found in token")
		return
	}

	if input.Filters == nil {
		input.Filters = map[string]string{}
	}

	input.Filters["branch_code"] = strconv.FormatInt(*branchCodeVal.(*int64), 10)
	input.Filters["terminal_code"] = strconv.FormatInt(*terminalCodeVal.(*int64), 10)

	if status, ok := input.Filters["status"]; ok {
		switch strings.ToLower(strings.TrimSpace(status)) {
		case "true":
			input.Filters["status"] = "1"
		case "false":
			input.Filters["status"] = "0"
		}
	}

	customers, meta, err := h.service.Search(c.Request.Context(), input.ToPaginationQuery())
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to search customers")
		return
	}

	helper.MetaSuccessResponse(c, http.StatusOK, "customers retrieved successfully", customers, meta)
}

// CreateCustomer godoc
// @Summary Create customer
// @Description Create a new customer record
// @Tags master-customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body customer.CustomerReq true "Customer payload"
// @Success 201 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/customer [post]
func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var req CustomerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userName := middleware.GetUserEmail(c)
	if userName == "" {
		userName = "SYSTEM"
	}
	userID := middleware.GetUserID(c)
	if userID == 0 {
		helper.ErrorResponse(c, http.StatusUnauthorized, "user id not found in token")
		return
	}
	branchCodeVal, exists := c.Get(middleware.BranchCodeKey)
	if !exists || branchCodeVal == nil {
		helper.ErrorResponse(c, http.StatusUnauthorized, "branch code not found in token")
		return
	}
	terminalCodeVal, exists := c.Get(middleware.TerminalCodeKey)
	if !exists || terminalCodeVal == nil {
		helper.ErrorResponse(c, http.StatusUnauthorized, "terminal code not found in token")
		return
	}

	customerName := req.CustomerName
	branchCode := int(*branchCodeVal.(*int64))
	terminalCode := int(*terminalCodeVal.(*int64))
	authLocation, err := h.service.GetAuthLocation(c.Request.Context(), userID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to resolve auth location")
		return
	}
	branchName := authLocation.BranchName
	terminalName := authLocation.TerminalName

	customer := Customer{
		BranchCode:             &branchCode,
		BranchName:             &branchName,
		TerminalCode:           &terminalCode,
		TerminalName:           &terminalName,
		CustomerName:           &customerName,
		CustomerType:           req.CustomerType,
		ProfitCenter:           req.ProfitCenter,
		CustomerCountry:        req.CustomerCountry,
		CustomerAddress:        req.CustomerAddress,
		City:                   req.City,
		ContactPerson:          req.ContactPerson,
		PhoneNumber:            req.PhoneNumber,
		EmailAddress:           req.EmailAddress,
		FaxNumber:              req.FaxNumber,
		TaxIDNumber:            req.TaxIDNumber,
		TaxID16Digit:           req.TaxID16Digit,
		TaxBranchCode:          req.TaxBranchCode,
		NationalIDNumber:       req.NationalIDNumber,
		BusinessLicenseDate:    req.BusinessLicenseDate,
		TaxIDDocumentUpload:    req.TaxIDDocumentUpload,
		RegisteredTaxpayerName: req.RegisteredTaxpayerName,
		RegisteredTaxpayerAddr: req.RegisteredTaxpayerAddr,
		BusinessType:           req.BusinessType,
		BusinessEntityType:     req.BusinessEntityType,
		BankCode:               req.BankCode,
		BankAccountIDR:         req.BankAccountIDR,
		ForeignCurrencyAccount: req.ForeignCurrencyAccount,
		StartDate:              req.StartDate,
		EndDate:                req.EndDate,
		Status:                 req.Status,
		InternalNotes:          req.InternalNotes,
		CreationBy:             &userName,
		LastUpdatedBy:          &userName,
	}

	if err := h.service.Create(c.Request.Context(), &customer); err != nil {
		if errors.Is(err, ErrCustomerAlreadyExists) {
			helper.ErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to create customer")
		return
	}

	helper.SuccessResponse(c, http.StatusCreated, "customer created successfully", customer)
}

// UpdateCustomer godoc
// @Summary Update customer
// @Description Update an existing customer by id
// @Tags master-customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Param payload body customer.CustomerReq true "Customer payload"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/customer/{id} [put]
func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid customer id")
		return
	}

	var req CustomerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userName := middleware.GetUserEmail(c)
	if userName == "" {
		userName = "SYSTEM"
	}
	userID := middleware.GetUserID(c)
	if userID == 0 {
		helper.ErrorResponse(c, http.StatusUnauthorized, "user id not found in token")
		return
	}
	branchCodeVal, exists := c.Get(middleware.BranchCodeKey)
	if !exists || branchCodeVal == nil {
		helper.ErrorResponse(c, http.StatusUnauthorized, "branch code not found in token")
		return
	}
	terminalCodeVal, exists := c.Get(middleware.TerminalCodeKey)
	if !exists || terminalCodeVal == nil {
		helper.ErrorResponse(c, http.StatusUnauthorized, "terminal code not found in token")
		return
	}

	customerName := req.CustomerName
	branchCode := int(*branchCodeVal.(*int64))
	terminalCode := int(*terminalCodeVal.(*int64))
	authLocation, err := h.service.GetAuthLocation(c.Request.Context(), userID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to resolve auth location")
		return
	}
	branchName := authLocation.BranchName
	terminalName := authLocation.TerminalName

	customer := Customer{
		BranchCode:             &branchCode,
		BranchName:             &branchName,
		TerminalCode:           &terminalCode,
		TerminalName:           &terminalName,
		CustomerName:           &customerName,
		CustomerType:           req.CustomerType,
		ProfitCenter:           req.ProfitCenter,
		CustomerCountry:        req.CustomerCountry,
		CustomerAddress:        req.CustomerAddress,
		City:                   req.City,
		ContactPerson:          req.ContactPerson,
		PhoneNumber:            req.PhoneNumber,
		EmailAddress:           req.EmailAddress,
		FaxNumber:              req.FaxNumber,
		TaxIDNumber:            req.TaxIDNumber,
		TaxID16Digit:           req.TaxID16Digit,
		TaxBranchCode:          req.TaxBranchCode,
		NationalIDNumber:       req.NationalIDNumber,
		BusinessLicenseDate:    req.BusinessLicenseDate,
		TaxIDDocumentUpload:    req.TaxIDDocumentUpload,
		RegisteredTaxpayerName: req.RegisteredTaxpayerName,
		RegisteredTaxpayerAddr: req.RegisteredTaxpayerAddr,
		BusinessType:           req.BusinessType,
		BusinessEntityType:     req.BusinessEntityType,
		BankCode:               req.BankCode,
		BankAccountIDR:         req.BankAccountIDR,
		ForeignCurrencyAccount: req.ForeignCurrencyAccount,
		StartDate:              req.StartDate,
		EndDate:                req.EndDate,
		Status:                 req.Status,
		InternalNotes:          req.InternalNotes,
		LastUpdatedBy:          &userName,
	}

	existingCustomer, err := h.service.FindByID(c.Request.Context(), id)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to find customer")
		return
	}
	customer.CustomerCode = existingCustomer.CustomerCode

	if err := h.service.Update(c.Request.Context(), id, &customer); err != nil {
		if errors.Is(err, ErrCustomerAlreadyExists) {
			helper.ErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to update customer")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "customer updated successfully", customer)
}

// DeleteCustomer godoc
// @Summary Delete customer
// @Description Delete customer by id
// @Tags master-customers
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Success 200 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Failure 500 {object} helper.Response
// @Router /master/customer/{id} [delete]
func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid customer id")
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "failed to delete customer")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "customer deleted successfully", nil)
}
