package branch

import "time"

type Branch struct {
	ID                uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	BranchCode        string     `gorm:"column:branch_code;unique;not null" json:"branch_code"`
	BranchName        string     `gorm:"column:branch_name;not null" json:"branch_name"`
	CompanyCode       string     `gorm:"column:company_code" json:"company_code"`
	CompanyName       string     `gorm:"column:company_name" json:"company_name"`
	KdPort            string     `gorm:"column:kd_port" json:"kd_port"`
	RegionalArea      string     `gorm:"column:regional_area" json:"regional_area"`
	ProfitCenter      string     `gorm:"column:profit_center" json:"profit_center"`
	Status            string     `gorm:"column:status" json:"status"`
	CreatedBy         string     `gorm:"column:created_by" json:"created_by"`
	CreatedDate       *time.Time `gorm:"column:created_date" json:"created_date"`
	LastUpdatedBy     string     `gorm:"column:last_updated_by" json:"last_updated_by"`
	LastUpdatedDate   *time.Time `gorm:"column:last_updated_date" json:"last_updated_date"`
	ProgramName       string     `gorm:"column:program_name" json:"program_name"`
}

func (Branch) TableName() string { return "adm.posm_branches" }
