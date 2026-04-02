package dermaga

import "time"

type Dermaga struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	KdCabang    uint      `gorm:"not null" json:"kd_cabang"`
	KdTerminal  uint      `gorm:"not null" json:"kd_terminal"`
	NmCabang    string    `gorm:"not null" json:"nm_cabang"`
	NmTerminal  string    `gorm:"not null" json:"nm_terminal"`
	NmDermaga   string    `gorm:"not null" json:"nm_dermaga"`
	KdDermaga   string    `gorm:"not null" json:"kd_dermaga"`
	PosisiAwal  uint      `gorm:"not null" json:"posisi_awal"`
	PosisiAkhir uint      `gorm:"not null" json:"posisi_akhir"`
	Keterangan  string    `gorm:"not null" json:"keterangan"`
	Status      string    `gorm:"not null" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   string    `gorm:"not null" json:"created_by"`
	UpdatedBy   string    `gorm:"not null" json:"updated_by"`
}

func (Dermaga) TableName() string { return "dermaga" }
