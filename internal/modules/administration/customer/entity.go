package customer

import "time"

type Customer struct {
	ID                     uint64     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	BranchCode             *int       `gorm:"column:branch_code" json:"branch_code"`
	BranchName             *string    `gorm:"column:branch_name;size:100" json:"branch_name"`
	TerminalCode           *int       `gorm:"column:terminal_code" json:"terminal_code"`
	TerminalName           *string    `gorm:"column:terminal_name;size:100" json:"terminal_name"`
	CustomerCode           *string    `gorm:"column:customer_code;size:30" json:"customer_code"`
	CustomerName           *string    `gorm:"column:customer_name;size:150" json:"customer_name"`
	CustomerType           *string    `gorm:"column:customer_type;size:50" json:"customer_type"`
	ProfitCenter           *string    `gorm:"column:profit_center;size:50" json:"profit_center"`
	CustomerCountry        *string    `gorm:"column:customer_country;size:100" json:"customer_country"`
	CustomerAddress        *string    `gorm:"column:customer_address;size:255" json:"customer_address"`
	City                   *string    `gorm:"column:city;size:100" json:"city"`
	ContactPerson          *string    `gorm:"column:contact_person;size:100" json:"contact_person"`
	PhoneNumber            *string    `gorm:"column:phone_number;size:30" json:"phone_number"`
	EmailAddress           *string    `gorm:"column:email_address;size:150" json:"email_address"`
	FaxNumber              *string    `gorm:"column:fax_number;size:30" json:"fax_number"`
	TaxIDNumber            *string    `gorm:"column:tax_id_number;size:30" json:"tax_id_number"`
	TaxID16Digit           *string    `gorm:"column:tax_id_16_digit;size:30" json:"tax_id_16_digit"`
	TaxBranchCode          *string    `gorm:"column:tax_branch_code;size:50" json:"tax_branch_code"`
	NationalIDNumber       *string    `gorm:"column:national_id_number;size:30" json:"national_id_number"`
	BusinessLicenseDate    *time.Time `gorm:"column:business_license_date" json:"business_license_date"`
	TaxIDDocumentUpload    *string    `gorm:"column:tax_id_document_upload" json:"tax_id_document_upload"`
	RegisteredTaxpayerName *string    `gorm:"column:registered_taxpayer_name;size:200" json:"registered_taxpayer_name"`
	RegisteredTaxpayerAddr *string    `gorm:"column:registered_taxpayer_address;size:255" json:"registered_taxpayer_address"`
	BusinessType           *string    `gorm:"column:business_type;size:100" json:"business_type"`
	BusinessEntityType     *string    `gorm:"column:business_entity_type;size:100" json:"business_entity_type"`
	BankCode               *string    `gorm:"column:bank_code;size:50" json:"bank_code"`
	BankAccountIDR         *string    `gorm:"column:bank_account_idr;size:50" json:"bank_account_idr"`
	ForeignCurrencyAccount *string    `gorm:"column:foreign_currency_account;size:50" json:"foreign_currency_account"`
	StartDate              *time.Time `gorm:"column:start_date" json:"start_date"`
	EndDate                *time.Time `gorm:"column:end_date" json:"end_date"`
	ProgramName            *string    `gorm:"column:program_name;size:50" json:"program_name"`
	Status                 *int       `gorm:"column:status" json:"status"`
	InternalNotes          *string    `gorm:"column:internal_notes" json:"internal_notes"`
	CreationDate           *time.Time `gorm:"column:creation_date" json:"creation_date"`
	CreationBy             *string    `gorm:"column:creation_by;size:100" json:"creation_by"`
	LastUpdatedDate        *time.Time `gorm:"column:last_updated_date" json:"last_updated_date"`
	LastUpdatedBy          *string    `gorm:"column:last_updated_by;size:100" json:"last_updated_by"`
}

func (Customer) TableName() string { return "adm.posm_customers" }
