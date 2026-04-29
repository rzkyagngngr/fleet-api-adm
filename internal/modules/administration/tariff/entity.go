package tariff

import "time"

type Tariff struct {
	ID              uint64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	BranchCode      *int            `gorm:"column:branch_code" json:"branch_code"`
	BranchName      *string         `gorm:"column:branch_name;size:100" json:"branch_name"`
	TerminalCode    *int            `gorm:"column:terminal_code" json:"terminal_code"`
	TerminalName    *string         `gorm:"column:terminal_name;size:100" json:"terminal_name"`
	NameTariff      string          `gorm:"column:name_tariff;size:100;not null" json:"name_tariff"`
	Description     *string         `gorm:"column:description" json:"description"`
	Status          *int            `gorm:"column:status" json:"status"`
	AgreementNumber *string         `gorm:"column:agreement_number;size:100" json:"agreement_number"`
	StartDate       *time.Time      `gorm:"column:start_date" json:"start_date"`
	EndDate         *time.Time      `gorm:"column:end_date" json:"end_date"`
	CreationDate    *time.Time      `gorm:"column:creation_date" json:"creation_date"`
	CreationBy      *string         `gorm:"column:creation_by;size:100" json:"creation_by"`
	LastUpdatedDate *time.Time      `gorm:"column:last_updated_date" json:"last_updated_date"`
	LastUpdatedBy   *string         `gorm:"column:last_updated_by;size:100" json:"last_updated_by"`
	Details         []TariffService `gorm:"foreignKey:IDTariff;references:ID" json:"details,omitempty"`
}

func (Tariff) TableName() string { return "adm.posm_tariffs" }

type TariffService struct {
	IDTariff       uint64   `gorm:"column:id_tariff" json:"id_tariff"`
	BranchCode     int      `gorm:"column:branch_code" json:"branch_code"`
	BranchName     *string  `gorm:"column:branch_name;size:100" json:"branch_name"`
	TerminalCode   int      `gorm:"column:terminal_code" json:"terminal_code"`
	TerminalName   *string  `gorm:"column:terminal_name;size:100" json:"terminal_name"`
	SequenceNo     *int     `gorm:"column:sequence_no" json:"sequence_no"`
	ServiceType    *string  `gorm:"column:service_type;size:50" json:"service_type"`
	ServiceName    *string  `gorm:"column:service_name;size:50" json:"service_name"`
	CustomerName   *string  `gorm:"column:customer_name;size:100" json:"customer_name"`
	CustomerCode   *string  `gorm:"column:customer_code;size:25" json:"customer_code"`
	CargoCode      *string  `gorm:"column:cargo_code;size:25" json:"cargo_code"`
	CargoName      *string  `gorm:"column:cargo_name;size:50" json:"cargo_name"`
	CargoPackaging *string  `gorm:"column:cargo_packaging;size:12" json:"cargo_packaging"`
	CargoUnit      *string  `gorm:"column:cargo_unit;size:20" json:"cargo_unit"`
	EquipmentCode  *string  `gorm:"column:equipment_code;size:25" json:"equipment_code"`
	EquipmentName  *string  `gorm:"column:equipment_name;size:50" json:"equipment_name"`
	EquipmentGroup *string  `gorm:"column:equipment_group;size:255" json:"equipment_group"`
	EquipmentUnit  *string  `gorm:"column:equipment_unit;size:255" json:"equipment_unit"`
	BaseTariff     *float64 `gorm:"column:base_tariff" json:"base_tariff"`
	CurrencyCode   *string  `gorm:"column:currency_code;size:25" json:"currency_code"`
	Discount       *float64 `gorm:"column:discount" json:"discount"`
	Attrib1        *string  `gorm:"column:attrib1;size:50" json:"attrib1"`
	Attrib2        *string  `gorm:"column:attrib2;size:50" json:"attrib2"`
	Attrib3        *string  `gorm:"column:attrib3;size:50" json:"attrib3"`
}

func (TariffService) TableName() string { return "adm.posm_tariff_services" }
