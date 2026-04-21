package customer

import (
	"omniport-api/internal/helper"
	"time"
)

type CustomerReq struct {
	BranchCode             *int       `json:"branch_code"`
	BranchName             *string    `json:"branch_name" binding:"omitempty,max=100"`
	TerminalCode           *int       `json:"terminal_code"`
	TerminalName           *string    `json:"terminal_name" binding:"omitempty,max=100"`
	CustomerName           string     `json:"customer_name" binding:"required,max=150"`
	CustomerType           *string    `json:"customer_type" binding:"omitempty,max=50"`
	ProfitCenter           *string    `json:"profit_center" binding:"omitempty,max=50"`
	CustomerCountry        *string    `json:"customer_country" binding:"omitempty,max=100"`
	CustomerAddress        *string    `json:"customer_address" binding:"omitempty,max=255"`
	City                   *string    `json:"city" binding:"omitempty,max=100"`
	ContactPerson          *string    `json:"contact_person" binding:"omitempty,max=100"`
	PhoneNumber            *string    `json:"phone_number" binding:"omitempty,max=30"`
	EmailAddress           *string    `json:"email_address" binding:"omitempty,max=150,email"`
	FaxNumber              *string    `json:"fax_number" binding:"omitempty,max=30"`
	TaxIDNumber            *string    `json:"tax_id_number" binding:"required,max=30"`
	TaxID16Digit           *string    `json:"tax_id_16_digit" binding:"omitempty,max=30"`
	TaxBranchCode          *string    `json:"tax_branch_code" binding:"omitempty,max=50"`
	NationalIDNumber       *string    `json:"national_id_number" binding:"omitempty,max=30"`
	BusinessLicenseDate    *time.Time `json:"business_license_date"`
	TaxIDDocumentUpload    *string    `json:"tax_id_document_upload"`
	RegisteredTaxpayerName *string    `json:"registered_taxpayer_name" binding:"omitempty,max=200"`
	RegisteredTaxpayerAddr *string    `json:"registered_taxpayer_address" binding:"omitempty,max=255"`
	BusinessType           *string    `json:"business_type" binding:"omitempty,max=100"`
	BusinessEntityType     *string    `json:"business_entity_type" binding:"omitempty,max=100"`
	BankCode               *string    `json:"bank_code" binding:"omitempty,max=50"`
	BankAccountIDR         *string    `json:"bank_account_idr" binding:"omitempty,max=50"`
	ForeignCurrencyAccount *string    `json:"foreign_currency_account" binding:"omitempty,max=50"`
	Status                 *int       `json:"status"`
	InternalNotes          *string    `json:"internal_notes"`
}

type SearchCustomerRequest struct {
	Page     int                  `json:"page"`
	Limit    int                  `json:"limit"`
	Search   string               `json:"search"`
	Filters  map[string]string    `json:"filters"`
	Sort     helper.SortQuery     `json:"sort"`
	Download helper.DownloadQuery `json:"download"`
}

func (r SearchCustomerRequest) ToPaginationQuery() helper.PaginationQuery {
	return helper.PaginationQuery{
		Page:     r.Page,
		Limit:    r.Limit,
		Search:   r.Search,
		Filters:  r.Filters,
		Sort:     r.Sort,
		Download: r.Download,
	}
}
