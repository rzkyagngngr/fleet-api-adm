package reference

import "time"

type PosmReference struct {
	ID              int64            `gorm:"primaryKey;autoIncrement" json:"id"`
	BranchCode      int64            `gorm:"not null" json:"branch_code"`
	TerminalCode    *int64           `json:"terminal_code"`
	IdTable         string           `gorm:"size:60" json:"id_table"`
	IdRefFile       string           `gorm:"not null;size:60" json:"id_ref_file"`
	Keterangan      string           `gorm:"size:100" json:"keterangan"`
	Ref1            string           `gorm:"size:100" json:"ref1"`
	Val1            *int64           `json:"val1"`
	KdAktif         string           `gorm:"size:1" json:"kd_aktif"`
	CreationDate    time.Time        `gorm:"not null;default:CURRENT_TIMESTAMP" json:"creation_date"`
	CreationBy      string           `gorm:"not null;size:30" json:"creation_by"`
	LastUpdatedDate *time.Time       `json:"last_updated_date"`
	LastUpdatedBy   string           `gorm:"size:10" json:"last_updated_by"`
	ProgramName     string           `gorm:"not null;size:50" json:"program_name"`
	LevelAkses      string           `gorm:"not null;size:30" json:"level_akses"`
	Details         []PosmReferenceD `gorm:"foreignKey:IdRefFile;references:IdRefFile" json:"details"`
}

func (PosmReference) TableName() string { return "posm_reference" }

type PosmReferenceD struct {
	ID              int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	BranchCode      int64      `gorm:"not null" json:"branch_code"`
	TerminalCode    *int64     `json:"terminal_code"`
	IdTable         string     `gorm:"size:20" json:"id_table"`
	IdRefFile       string     `gorm:"not null;size:100" json:"id_ref_file"`
	IdRefKey        string     `gorm:"not null;size:100" json:"id_ref_key"`
	KetRefData      string     `gorm:"size:330" json:"ket_ref_data"`
	Val1            *int64     `json:"val1"`
	Val2            *int64     `json:"val2"`
	Val3            *int64     `json:"val3"`
	Val4            *int64     `json:"val4"`
	Val5            *int64     `json:"val5"`
	Ref1            string     `gorm:"size:50" json:"ref1"`
	Ref2            string     `gorm:"size:50" json:"ref2"`
	Ref3            string     `gorm:"size:50" json:"ref3"`
	Ref4            string     `gorm:"size:50" json:"ref4"`
	Ref5            string     `gorm:"size:50" json:"ref5"`
	KdAktif         string     `gorm:"size:1" json:"kd_aktif"`
	CreationDate    time.Time  `gorm:"not null" json:"creation_date"`
	CreationBy      string     `gorm:"not null;size:30" json:"creation_by"`
	LastUpdatedDate *time.Time `json:"last_updated_date"`
	LastUpdatedBy   string     `gorm:"size:10" json:"last_updated_by"`
	ProgramName     string     `gorm:"not null;size:50" json:"program_name"`
	LevelAkses      string     `gorm:"not null;size:30" json:"level_akses"`
	Rc              *int64     `json:"rc"`
}

func (PosmReferenceD) TableName() string { return "posm_reference_d" }
