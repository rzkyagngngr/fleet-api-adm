package pelabuhan

import "time"

type Port struct {
	ID              uint64     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	PortCode        string     `gorm:"column:port_code;size:10;not null" json:"port_code"`
	PortName        *string    `gorm:"column:port_name;size:40" json:"port_name"`
	PortCity        *string    `gorm:"column:port_city;size:40" json:"port_city"`
	CountryCode     *string    `gorm:"column:country_code;size:3" json:"country_code"`
	CreatedBy       *string    `gorm:"column:created_by;size:30" json:"created_by"`
	CreatedDate     *time.Time `gorm:"column:created_date" json:"created_date"`
	LastUpdatedDate time.Time  `gorm:"column:last_updated_date;not null" json:"last_updated_date"`
	LastUpdatedBy   string     `gorm:"column:last_updated_by;size:30;not null" json:"last_updated_by"`
	ProgramName     string     `gorm:"column:program_name;size:30;not null" json:"program_name"`
	Status          *string    `gorm:"column:status;type:bpchar(5)" json:"status"`
}

func (Port) TableName() string { return "adm.posm_port" }
