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
	Payload                map[string]interface{} `json:"payload"`
	Op                     *CreateOpInput        `json:"op"`
}

type CreateOpInput struct {
	ID               uint64                `json:"id"`
	Pbm              string                `json:"pbm"`
	Emkl             string                `json:"emkl"`
	Shipper          string                `json:"shipper"`
	StartDischarging *time.Time            `json:"start_discharging"`
	EndDischarging   *time.Time            `json:"end_discharging"`
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
	Payload                JSONB             `json:"payload"`
	Op                     *VesselRpkOpResponse `json:"op"`
	BranchCode             int64             `json:"branch_code"`
	TerminalCode           int64             `json:"terminal_code"`
	CreationDate           time.Time         `json:"creation_date"`
	CreationBy             string            `json:"creation_by"`
}

type VesselRpkOpResponse struct {
	ID               uint64                  `json:"id"`
	Pbm              string                  `json:"pbm"`
	Emkl             string                  `json:"emkl"`
	Shipper          string                  `json:"shipper"`
	StartDischarging *time.Time              `json:"start_discharging"`
	EndDischarging   *time.Time              `json:"end_discharging"`
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
