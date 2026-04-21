package terminal

import "time"

type Terminal struct {
	ID                uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	BranchCode        int64      `gorm:"column:branch_code;not null" json:"branch_code"`
	BranchName        string     `gorm:"column:branch_name" json:"branch_name"`
	TerminalCode      int64      `gorm:"column:terminal_code;unique;not null" json:"terminal_code"`
	TerminalName      string     `gorm:"column:terminal_name" json:"terminal_name"`
	GoLiveDate        *time.Time `gorm:"column:go_live_date" json:"go_live_date"`
	IsGoLive          string     `gorm:"column:is_go_live;size:2" json:"is_go_live"`
	ProfitCenter      string     `gorm:"column:profit_center" json:"profit_center"`
	Latitude          string     `gorm:"column:latitude" json:"latitude"`
	Longitude         string     `gorm:"column:longitude" json:"longitude"`
	Status            string     `gorm:"column:status" json:"status"`
	VersionCode       int64      `gorm:"column:version_code" json:"version_code"`
	VersionName       string     `gorm:"column:version_name" json:"version_name"`
	DocumentCode      string     `gorm:"column:document_code" json:"document_code"`
	CompanyCode       string     `gorm:"column:company_code" json:"company_code"`
	CompanyName       string     `gorm:"column:company_name" json:"company_name"`
	VesselVersion     int64      `gorm:"column:vessel_version" json:"vessel_version"`
	LogoURL           string     `gorm:"column:logo_url" json:"logo_url"`
	LogoMiniURL       string     `gorm:"column:logo_mini_url" json:"logo_mini_url"`
	Address           string     `gorm:"column:address" json:"address"`
	CompanyType       int64      `gorm:"column:company_type" json:"company_type"`
	PortCode          string     `gorm:"column:port_code" json:"port_code"`
	
	// Audit Fields
	CreatedBy         string     `gorm:"column:created_by" json:"created_by"`
	CreatedDate       *time.Time `gorm:"column:created_date" json:"created_date"`
	LastUpdatedBy     string     `gorm:"column:last_updated_by" json:"last_updated_by"`
	LastUpdatedDate   *time.Time `gorm:"column:last_updated_date" json:"last_updated_date"`
	ProgramName       string     `gorm:"column:program_name" json:"program_name"`
}

func (Terminal) TableName() string { return "adm.posm_terminals" }
