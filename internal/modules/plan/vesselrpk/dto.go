package vesselrpk

import "time"

// ─────────────────────────────────────────────────────────────
// REQUEST DTOs
// ─────────────────────────────────────────────────────────────

type CreateVesselRpkInput struct {
	NoPkk                  string               `json:"no_pkk" binding:"required"`
	NoPpk                  string               `json:"no_ppk"`
	LocationCodeInaportnet string               `json:"location_code_inaportnet"`
	RpkType                string               `json:"rpk_type"`
	BerthPosition          string               `json:"berth_position"`
	VesselPosition         string               `json:"vessel_position"`
	StartMeter             string               `json:"start_meter"`
	EndMeter               string               `json:"end_meter"`
	StartMooring           *time.Time           `json:"start_mooring"`
	EndMooring             *time.Time           `json:"end_mooring"`
	RampDoor               string               `json:"ramp_door"`
	Distribution           string               `json:"distribution"`
	Packaging              string               `json:"packaging"`
	NoRkbm                 string               `json:"no_rkbm"`
	Reason                 string               `json:"reason"`
	Notes                  string               `json:"notes"`
	OpsPlanCode            string               `json:"ops_plan_code"`
	ActivityCode           string               `json:"activity_code"`
	Payload                map[string]interface{} `json:"payload"`
	BranchCode             int64                `json:"branch_code"`
	TerminalCode           int64                `json:"terminal_code"`
	Ops                    []CreateOpInput        `json:"ops"`
}

type CreateOpInput struct {
	ID               uint64                `json:"id"`
	Pbm              string                `json:"pbm"`
	Emkl             string                `json:"emkl"`
	Shipper          string                `json:"shipper"`
	StartDischarging *time.Time            `json:"start_discharging"`
	EndDischarging   *time.Time            `json:"end_discharging"`
	StartActivityDate *time.Time           `json:"start_activity_date"`
	EndActivityDate   *time.Time           `json:"end_activity_date"`
	OpDetail         []CreateOpDetailInput `json:"op_detail"`
}

type CreateOpDetailInput struct {
	ID                uint64 `json:"id"`
	RkbmMuatNumber    string `json:"rkbm_muat_number"`
	RkbmBongkarNumber string `json:"rkbm_bongkar_number"`
	Loading           string `json:"loading"`
	Discharging       string `json:"discharging"`
	Commodity         string `json:"commodity"`
}

// ─────────────────────────────────────────────────────────────
// RESPONSE DTOs
// ─────────────────────────────────────────────────────────────

type VesselRpkResponse struct {
	ID                     uint64            `json:"id"`
	NoPkk                  string            `json:"no_pkk"`
	NoPpk                  string            `json:"no_ppk"`
	VesselName             string            `json:"vessel_name"`
	LocationCodeInaportnet string            `json:"location_code_inaportnet"`
	RpkType                string            `json:"rpk_type"`
	BerthPosition          string            `json:"berth_position"`
	VesselPosition         string            `json:"vessel_position"`
	StartMeter             string            `json:"start_meter"`
	EndMeter               string            `json:"end_meter"`
	StartMooring           *time.Time        `json:"start_mooring"`
	EndMooring             *time.Time        `json:"end_mooring"`
	RampDoor               string            `json:"ramp_door"`
	Distribution           string            `json:"distribution"`
	Packaging              string            `json:"packaging"`
	NoRkbm                 string            `json:"no_rkbm"`
	Reason                 string            `json:"reason"`
	Notes                  string            `json:"notes"`
	OpsPlanCode            string            `json:"ops_plan_code"`
	ActivityCode           string            `json:"activity_code"`
	Payload                JSONB             `json:"payload"`
	Ops                    []VesselRpkOpResponse `json:"ops"`
	BranchCode             int64             `json:"branch_code"`
	TerminalCode           int64             `json:"terminal_code"`
	CreationDate           time.Time         `json:"creation_date"`
	CreationBy             string            `json:"creation_by"`
	LastUpdatedDate        *time.Time        `json:"last_updated_date"`
	LastUpdatedBy          string            `json:"last_updated_by"`
}

type VesselRpkOpResponse struct {
	ID               uint64                  `json:"id"`
	Pbm              string                  `json:"pbm"`
	Emkl             string                  `json:"emkl"`
	Shipper          string                  `json:"shipper"`
	StartDischarging *time.Time              `json:"start_discharging"`
	EndDischarging   *time.Time              `json:"end_discharging"`
	StartActivityDate *time.Time             `json:"start_activity_date"`
	EndActivityDate   *time.Time             `json:"end_activity_date"`
	OpDetail         []VesselRpkOpDetailResponse `json:"op_detail"`
}

type VesselRpkOpDetailResponse struct {
	ID                uint64 `json:"id"`
	RkbmMuatNumber    string `json:"rkbm_muat_number"`
	RkbmBongkarNumber string `json:"rkbm_bongkar_number"`
	Loading           string `json:"loading"`
	Discharging       string `json:"discharging"`
	Commodity         string `json:"commodity"`
}
