package vesselrpk

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &j)
}

type VesselRpk struct {
	ID                     uint64        `gorm:"primaryKey;autoIncrement" json:"id"`
	NoPkk                  string        `gorm:"column:no_pkk;not null" json:"no_pkk"`
	NoPpk                  string        `gorm:"column:no_ppk" json:"no_ppk"`
	LocationCodeInaportnet string        `gorm:"column:location_code_inaportnet" json:"location_code_inaportnet"`
	RpkType                string        `gorm:"column:rpk_type" json:"rpk_type"`
	BerthPosition          string        `gorm:"column:berth_position" json:"berth_position"`
	VesselPosition         string        `gorm:"column:vessel_position" json:"vessel_position"`
	StartMeter             string        `gorm:"column:start_meter" json:"start_meter"`
	EndMeter               string        `gorm:"column:end_meter" json:"end_meter"`
	StartMooring           *time.Time    `gorm:"column:start_mooring" json:"start_mooring"`
	EndMooring             *time.Time    `gorm:"column:end_mooring" json:"end_mooring"`
	RampDoor               string        `gorm:"column:ramp_door" json:"ramp_door"`
	Distribution           string        `gorm:"column:distribution" json:"distribution"`
	Packaging              string        `gorm:"column:packaging" json:"packaging"`
	NoRkbm                 string        `gorm:"column:no_rkbm" json:"no_rkbm"`
	Reason                 string        `gorm:"column:reason" json:"reason"`
	Notes                  string        `gorm:"column:notes" json:"notes"`
	Payload                JSONB         `gorm:"column:payload;type:jsonb" json:"payload"`
	
	// Relations
	Op *VesselRpkOp `gorm:"foreignKey:VesselRpkID" json:"op"`

	// Audit & Multi-tenancy
	BranchCode      int64      `gorm:"column:branch_code" json:"branch_code"`
	TerminalCode    int64      `gorm:"column:terminal_code" json:"terminal_code"`
	CreationDate    time.Time  `gorm:"column:creation_date;autoCreateTime" json:"creation_date"`
	CreationBy      string     `gorm:"column:creation_by" json:"creation_by"`
	LastUpdatedDate *time.Time `gorm:"column:last_updated_date;autoUpdateTime" json:"last_updated_date"`
	LastUpdatedBy   string     `gorm:"column:last_updated_by" json:"last_updated_by"`
}

func (VesselRpk) TableName() string { return "plan.post_vessel_rpk" }

type VesselRpkOp struct {
	ID               uint64               `gorm:"primaryKey;autoIncrement" json:"id"`
	VesselRpkID      uint64               `gorm:"column:vessel_rpk_id" json:"vessel_rpk_id"`
	Pbm              string               `gorm:"column:pbm" json:"pbm"`
	Emkl             string               `gorm:"column:emkl" json:"emkl"`
	Shipper          string               `gorm:"column:shipper" json:"shipper"`
	StartDischarging *time.Time           `gorm:"column:start_discharging" json:"start_discharging"`
	EndDischarging   *time.Time           `gorm:"column:end_discharging" json:"end_discharging"`
	OpDetail         []VesselRpkOpDetail `gorm:"foreignKey:VesselRpkOpID" json:"op_detail"`
}

func (VesselRpkOp) TableName() string { return "plan.post_vessel_rpk_op" }

type VesselRpkOpDetail struct {
	ID                uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	VesselRpkOpID      uint64 `gorm:"column:vessel_rpk_op_id" json:"vessel_rpk_op_id"`
	RkbmMuatNumber    string `gorm:"column:rkbm_muat_number" json:"rkbm_muat_number"`
	RkbmBongkarNumber string `gorm:"column:rkbm_bongkar_number" json:"rkbm_bongkar_number"`
	Loading           string `gorm:"column:loading" json:"loading"`
	Discharging       string `gorm:"column:discharging" json:"discharging"`
	Commodity         string `gorm:"column:commodity" json:"commodity"`
}

func (VesselRpkOpDetail) TableName() string { return "plan.post_vessel_rpk_op_detail" }
