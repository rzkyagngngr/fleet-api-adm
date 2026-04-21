package cargo

import "time"

type Cargo struct {
	ID                    uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	BranchCode            int        `gorm:"column:branch_code;not null" json:"branch_code"`
	TerminalCode          int        `gorm:"column:terminal_code;not null" json:"terminal_code"`
	CargoCode             string     `gorm:"column:cargo_code;not null" json:"cargo_code"`
	CargoSitcCode         string     `gorm:"column:cargo_sitc_code" json:"cargo_sitc_code"`
	CargoHsHarmonizedCode string     `gorm:"column:cargo_hs_harmonized_code" json:"cargo_hs_harmonized_code"`
	CargoName             string     `gorm:"column:cargo_name" json:"cargo_name"`
	CargoGroup            string     `gorm:"column:cargo_group" json:"cargo_group"`
	CargoCommodity        string     `gorm:"column:cargo_commodity" json:"cargo_commodity"`
	CargoCharacteristic   string     `gorm:"column:cargo_characteristic" json:"cargo_characteristic"`
	CargoImdgCode         *int16     `gorm:"column:cargo_imdg_code" json:"cargo_imdg_code"`
	CargoImdgDescription  string     `gorm:"column:cargo_imdg_description" json:"cargo_imdg_description"`
	CargoPackaging1       string     `gorm:"column:cargo_packaging_1" json:"cargo_packaging_1"`
	CargoConversion1      float64    `gorm:"column:cargo_conversion_1" json:"cargo_conversion_1"`
	CargoDimension1       float64    `gorm:"column:cargo_dimension_1" json:"cargo_dimension_1"`
	CargoUnit1            string     `gorm:"column:cargo_unit_1" json:"cargo_unit_1"`
	CargoPackaging2       string     `gorm:"column:cargo_packaging_2" json:"cargo_packaging_2"`
	CargoConversion2      float64    `gorm:"column:cargo_conversion_2" json:"cargo_conversion_2"`
	CargoDimension2       float64    `gorm:"column:cargo_dimension_2" json:"cargo_dimension_2"`
	CargoUnit2            string     `gorm:"column:cargo_unit_2" json:"cargo_unit_2"`
	CargoPackaging3       string     `gorm:"column:cargo_packaging_3" json:"cargo_packaging_3"`
	CargoConversion3      float64    `gorm:"column:cargo_conversion_3" json:"cargo_conversion_3"`
	CargoDimension3       float64    `gorm:"column:cargo_dimension_3" json:"cargo_dimension_3"`
	CargoUnit3            string     `gorm:"column:cargo_unit_3" json:"cargo_unit_3"`
	CargoMooringType      string     `gorm:"column:cargo_mooring_type" json:"cargo_mooring_type"`
	CargoNotes            string     `gorm:"column:cargo_notes" json:"cargo_notes"`
	CargoCommodityGroup   string     `gorm:"column:cargo_commodity_group" json:"cargo_commodity_group"`
	CargoCommodityType    string     `gorm:"column:cargo_commodity_type" json:"cargo_commodity_type"`
	IsActive              string     `gorm:"column:is_active;type:char(1)" json:"is_active"`
	CargoDocument         string     `gorm:"column:cargo_document" json:"cargo_document"`
	HsCode                string     `gorm:"column:hs_code" json:"hs_code"`
	HsDescription         string     `gorm:"column:hs_description" json:"hs_description"`
	CargoProductName      string     `gorm:"column:cargo_product_name" json:"cargo_product_name"`
	CreatedBy             string     `gorm:"column:created_by" json:"created_by"`
	CreatedDate           *time.Time `gorm:"column:created_date" json:"created_date"`
	LastUpdatedDate       time.Time  `gorm:"column:last_updated_date;not null" json:"last_updated_date"`
	LastUpdatedBy         string     `gorm:"column:last_updated_by;not null" json:"last_updated_by"`
	ProgramName           string     `gorm:"column:program_name;not null" json:"program_name"`
}

func (Cargo) TableName() string { return "posm_cargos" }
