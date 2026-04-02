package cabang

type Cabang struct {
	ID         uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	KdCabang   uint   `gorm:"not null" json:"kd_cabang"`
	NmCabang   string `gorm:"not null" json:"nm_cabang"`
	KdTerminal uint   `gorm:"not null" json:"kd_terminal"`
	NmTerminal string `gorm:"not null" json:"nm_terminal"`
}

func (Cabang) TableName() string { return "cabang" }
