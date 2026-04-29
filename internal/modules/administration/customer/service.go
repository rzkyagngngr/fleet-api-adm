package customer

import (
	"context"
	"errors"
	"fmt"
	"omniport-api/internal/helper"
	"time"

	"gorm.io/gorm"
)

var ErrCustomerAlreadyExists = errors.New("customer already exists")

type CustomerService interface {
	Search(ctx context.Context, query helper.PaginationQuery) ([]Customer, helper.PaginationMeta, error)
	Create(ctx context.Context, customer *Customer) error
	Update(ctx context.Context, id uint64, customer *Customer) error
	Delete(ctx context.Context, id uint64) error
	FindByID(ctx context.Context, id uint64) (*Customer, error)
	GetAuthLocation(ctx context.Context, userID uint64) (*CustomerAuthLocation, error)
}

type customerService struct {
	db *gorm.DB
}

type CustomerAuthLocation struct {
	BranchName   string
	TerminalName string
}

func NewCustomerService(db *gorm.DB) CustomerService {
	return &customerService{db: db}
}

func (s *customerService) Search(ctx context.Context, query helper.PaginationQuery) ([]Customer, helper.PaginationMeta, error) {
	config := helper.NativePaginationConfig{
		TableName: "adm.posm_customers",
		SelectColumns: []string{
			"id",
			"branch_code",
			"branch_name",
			"terminal_code",
			"terminal_name",
			"customer_code",
			"customer_name",
			"customer_type",
			"profit_center",
			"customer_country",
			"customer_address",
			"city",
			"contact_person",
			"phone_number",
			"email_address",
			"fax_number",
			"tax_id_number",
			"tax_id_16_digit",
			"tax_branch_code",
			"national_id_number",
			"business_license_date",
			"tax_id_document_upload",
			"registered_taxpayer_name",
			"registered_taxpayer_address",
			"business_type",
			"business_entity_type",
			"bank_code",
			"bank_account_idr",
			"foreign_currency_account",
			"program_name",
			"status",
			"internal_notes",
			"creation_date",
			"creation_by",
			"last_updated_date",
			"last_updated_by",
		},
		SearchColumns: []string{
			"branch_name",
			"terminal_name",
			"customer_code",
			"customer_name",
			"customer_type",
			"profit_center",
			"customer_country",
			"city",
			"contact_person",
			"phone_number",
			"email_address",
			"tax_id_number",
		},
		FilterableColumns: map[string]string{
			"branch_code":          "branch_code",
			"branch_name":          "branch_name",
			"terminal_code":        "terminal_code",
			"terminal_name":        "terminal_name",
			"customer_code":        "customer_code",
			"customer_name":        "customer_name",
			"customer_type":        "customer_type",
			"profit_center":        "profit_center",
			"customer_country":     "customer_country",
			"city":                 "city",
			"contact_person":       "contact_person",
			"phone_number":         "phone_number",
			"email_address":        "email_address",
			"tax_id_number":        "tax_id_number",
			"tax_id_16_digit":      "tax_id_16_digit",
			"tax_branch_code":      "tax_branch_code",
			"national_id_number":   "national_id_number",
			"business_type":        "business_type",
			"business_entity_type": "business_entity_type",
			"bank_code":            "bank_code",
			"status":               "status",
		},
		SortableColumns: map[string]string{
			"id":            "id",
			"branch_code":   "branch_code",
			"branch_name":   "branch_name",
			"terminal_code": "terminal_code",
			"terminal_name": "terminal_name",
			"customer_code": "customer_code",
			"customer_name": "customer_name",
			"customer_type": "customer_type",
			"creation_date": "creation_date",
			"last_updated":  "last_updated_date",
		},
		DefaultSortBy:    "customer_code",
		DefaultSortOrder: "ASC",
		MaxLimit:         100,
		MaxDownloadLimit: 1000,
	}


	var customers []Customer
	meta, err := helper.GetDynamicPaginatedNativeData(s.db.WithContext(ctx), config, query, &customers)
	return customers, meta, err
}

func (s *customerService) Create(ctx context.Context, c *Customer) error {
	const insertQuery = `
		INSERT INTO adm.posm_customers (
			branch_code,
			branch_name,
			terminal_code,
			terminal_name,
			customer_code,
			customer_name,
			customer_type,
			profit_center,
			customer_country,
			customer_address,
			city,
			contact_person,
			phone_number,
			email_address,
			fax_number,
			tax_id_number,
			tax_id_16_digit,
			tax_branch_code,
			national_id_number,
			business_license_date,
			tax_id_document_upload,
			registered_taxpayer_name,
			registered_taxpayer_address,
			business_type,
			business_entity_type,
			bank_code,
			bank_account_idr,
			foreign_currency_account,
			program_name,
			status,
			internal_notes,
			creation_date,
			creation_by,
			last_updated_date,
			last_updated_by
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id
	`

	const updateCodeQuery = `
		UPDATE adm.posm_customers
		SET customer_code = ?
		WHERE id = ?
	`

	now := time.Now()
	programName := "Master Customer"
	c.ProgramName = &programName
	c.CreationDate = &now
	c.LastUpdatedDate = &now

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		exists, err := s.customerExists(ctx, tx, c.TaxIDNumber, c.BranchCode, c.TerminalCode, nil)
		if err != nil {
			return err
		}
		if exists {
			return ErrCustomerAlreadyExists
		}

		if err := tx.Raw(
			insertQuery,
			c.BranchCode,
			c.BranchName,
			c.TerminalCode,
			c.TerminalName,
			c.CustomerCode,
			c.CustomerName,
			c.CustomerType,
			c.ProfitCenter,
			c.CustomerCountry,
			c.CustomerAddress,
			c.City,
			c.ContactPerson,
			c.PhoneNumber,
			c.EmailAddress,
			c.FaxNumber,
			c.TaxIDNumber,
			c.TaxID16Digit,
			c.TaxBranchCode,
			c.NationalIDNumber,
			c.BusinessLicenseDate,
			c.TaxIDDocumentUpload,
			c.RegisteredTaxpayerName,
			c.RegisteredTaxpayerAddr,
			c.BusinessType,
			c.BusinessEntityType,
			c.BankCode,
			c.BankAccountIDR,
			c.ForeignCurrencyAccount,
			c.ProgramName,
			c.Status,
			c.InternalNotes,
			c.CreationDate,
			c.CreationBy,
			c.LastUpdatedDate,
			c.LastUpdatedBy,
		).Scan(&c.ID).Error; err != nil {
			return err
		}

		customerCode := fmt.Sprintf("CUST%06d", c.ID)
		if err := tx.Exec(updateCodeQuery, customerCode, c.ID).Error; err != nil {
			return err
		}

		c.CustomerCode = &customerCode
		return nil
	})
}

func (s *customerService) Update(ctx context.Context, id uint64, c *Customer) error {
	const query = `
		UPDATE adm.posm_customers
		SET
			branch_code = ?,
			branch_name = ?,
			terminal_code = ?,
			terminal_name = ?,
			customer_code = ?,
			customer_name = ?,
			customer_type = ?,
			profit_center = ?,
			customer_country = ?,
			customer_address = ?,
			city = ?,
			contact_person = ?,
			phone_number = ?,
			email_address = ?,
			fax_number = ?,
			tax_id_number = ?,
			tax_id_16_digit = ?,
			tax_branch_code = ?,
			national_id_number = ?,
			business_license_date = ?,
			tax_id_document_upload = ?,
			registered_taxpayer_name = ?,
			registered_taxpayer_address = ?,
			business_type = ?,
			business_entity_type = ?,
			bank_code = ?,
			bank_account_idr = ?,
			foreign_currency_account = ?,
			program_name = ?,
			status = ?,
			internal_notes = ?,
			last_updated_date = ?,
			last_updated_by = ?
		WHERE id = ?
	`

	now := time.Now()
	programName := "Master Customer"
	c.ProgramName = &programName
	c.LastUpdatedDate = &now

	exists, err := s.customerExists(ctx, s.db.WithContext(ctx), c.TaxIDNumber, c.BranchCode, c.TerminalCode, &id)
	if err != nil {
		return err
	}
	if exists {
		return ErrCustomerAlreadyExists
	}

	result := s.db.WithContext(ctx).Exec(
		query,
		c.BranchCode,
		c.BranchName,
		c.TerminalCode,
		c.TerminalName,
		c.CustomerCode,
		c.CustomerName,
		c.CustomerType,
		c.ProfitCenter,
		c.CustomerCountry,
		c.CustomerAddress,
		c.City,
		c.ContactPerson,
		c.PhoneNumber,
		c.EmailAddress,
		c.FaxNumber,
		c.TaxIDNumber,
		c.TaxID16Digit,
		c.TaxBranchCode,
		c.NationalIDNumber,
		c.BusinessLicenseDate,
		c.TaxIDDocumentUpload,
		c.RegisteredTaxpayerName,
		c.RegisteredTaxpayerAddr,
		c.BusinessType,
		c.BusinessEntityType,
		c.BankCode,
		c.BankAccountIDR,
		c.ForeignCurrencyAccount,
		c.ProgramName,
		c.Status,
		c.InternalNotes,
		c.LastUpdatedDate,
		c.LastUpdatedBy,
		id,
	)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (s *customerService) Delete(ctx context.Context, id uint64) error {
	const query = `DELETE FROM adm.posm_customers WHERE id = ?`
	result := s.db.WithContext(ctx).Exec(query, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s *customerService) FindByID(ctx context.Context, id uint64) (*Customer, error) {
	const query = `
		SELECT
			id,
			branch_code,
			branch_name,
			terminal_code,
			terminal_name,
			customer_code,
			customer_name,
			customer_type,
			profit_center,
			customer_country,
			customer_address,
			city,
			contact_person,
			phone_number,
			email_address,
			fax_number,
			tax_id_number,
			tax_id_16_digit,
			tax_branch_code,
			national_id_number,
			business_license_date,
			tax_id_document_upload,
			registered_taxpayer_name,
			registered_taxpayer_address,
			business_type,
			business_entity_type,
			bank_code,
			bank_account_idr,
			foreign_currency_account,
			program_name,
			status,
			internal_notes,
			creation_date,
			creation_by,
			last_updated_date,
			last_updated_by
		FROM adm.posm_customers
		WHERE id = ?
		LIMIT 1
	`

	var customer Customer
	result := s.db.WithContext(ctx).Raw(query, id).Scan(&customer)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &customer, nil
}

func (s *customerService) GetAuthLocation(ctx context.Context, userID uint64) (*CustomerAuthLocation, error) {
	const userQuery = `
		SELECT branch_name, terminal_name
		FROM posm_users
		WHERE id = ?
		LIMIT 1
	`

	var result CustomerAuthLocation
	if err := s.db.WithContext(ctx).Raw(userQuery, userID).Scan(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *customerService) customerExists(ctx context.Context, db *gorm.DB, taxIDNumber *string, branchCode *int, terminalCode *int, excludeID *uint64) (bool, error) {
	query := `
		SELECT COUNT(1)
		FROM adm.posm_customers
		WHERE tax_id_number = ?
		  AND branch_code = ?
		  AND terminal_code = ?
	`

	args := []interface{}{taxIDNumber, branchCode, terminalCode}
	if excludeID != nil {
		query += " AND id <> ?"
		args = append(args, *excludeID)
	}

	var count int64
	if err := db.WithContext(ctx).Raw(query, args...).Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
