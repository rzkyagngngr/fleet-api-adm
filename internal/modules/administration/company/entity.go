package company

import "time"

type Company struct {
	ID              uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	CompanyCode     string     `gorm:"column:company_code;unique;not null" json:"company_code"`
	CompanyName     string     `gorm:"column:company_name;not null" json:"company_name"`
	Npwp            string     `gorm:"column:npwp" json:"npwp"`
	Address         string     `gorm:"column:address" json:"address"`
	Email           string     `gorm:"column:email" json:"email"`
	PhoneNumber     string     `gorm:"column:phone_number" json:"phone_number"`
	BusinessType    string     `gorm:"column:business_type" json:"business_type"`
	Status          string     `gorm:"column:status;default:1" json:"status"`
	CreatedBy       string     `gorm:"column:created_by" json:"created_by"`
	CreatedDate     *time.Time `gorm:"column:created_date;default:CURRENT_TIMESTAMP" json:"created_date"`
	LastUpdatedBy   string     `gorm:"column:last_updated_by" json:"last_updated_by"`
	LastUpdatedDate *time.Time `gorm:"column:last_updated_date" json:"last_updated_date"`
	ProgramName     string     `gorm:"column:program_name;default:OMNIPORT_ADM" json:"program_name"`
}

func (Company) TableName() string { return "adm.posm_companies" }
