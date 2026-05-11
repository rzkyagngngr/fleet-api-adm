package postrequest

import (
	"context"
	"fmt"
	"omniport-api/internal/helper"

	"gorm.io/gorm"
)

// ─────────────────────────────────────────────────────────────
// INTERFACE
// ─────────────────────────────────────────────────────────────

type PostRequestRepository interface {
	Create(ctx context.Context, header *PostRequest, details []PostRequestDetail, files []PostRequestFile) error
	FindByID(ctx context.Context, id int64) (*PostRequest, []PostRequestDetail, error)
	FindDetailsByRequestCode(ctx context.Context, requestCode string, branchCode, terminalCode int) ([]PostRequestDetail, error)
	UpdateHeader(ctx context.Context, id int64, header *PostRequest) error
	ReplaceDetails(ctx context.Context, requestCode string, branchCode, terminalCode int, details []PostRequestDetail) error
	ReplaceFiles(ctx context.Context, headerID int64, files []PostRequestFile) error
	Delete(ctx context.Context, id int64) error
	Search(ctx context.Context, param helper.PaginationQuery) ([]PostRequest, helper.PaginationMeta, error)
	GetStats(ctx context.Context, branchCode, terminalCode int) (*PostRequestStatsResponse, error)
	UpdateStatus(ctx context.Context, id int64, status int, remarks string, updatedBy string) error

	// Vessel Schedule methods
	SearchVesselSchedule(ctx context.Context, param helper.PaginationQuery) ([]PostVesselSchedule, helper.PaginationMeta, error)
	FindVesselScheduleByID(ctx context.Context, id int64) (*PostVesselSchedule, error)
}

// ─────────────────────────────────────────────────────────────
// IMPLEMENTATION
// ─────────────────────────────────────────────────────────────

type postRequestRepository struct{ db *gorm.DB }

func NewPostRequestRepository(db *gorm.DB) PostRequestRepository {
	return &postRequestRepository{db: db}
}

// Create inserts the header and all detail rows inside a single transaction.
func (r *postRequestRepository) Create(ctx context.Context, header *PostRequest, details []PostRequestDetail, files []PostRequestFile) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(header).Error; err != nil {
			return err
		}
		for i := range details {
			details[i].RequestCode = header.RequestCode
			details[i].BranchCode = *header.BranchCode
			details[i].TerminalCode = *header.TerminalCode
			details[i].BranchName = header.BranchName
			details[i].TerminalName = header.TerminalName
		}
		if len(details) > 0 {
			if err := tx.CreateInBatches(details, 100).Error; err != nil {
				return err
			}
		}

		// Save Attachments
		for i := range files {
			files[i].HeaderID = header.ID
		}
		if len(files) > 0 {
			if err := tx.Create(&files).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// FindByID retrieves header + all associated detail rows.
func (r *postRequestRepository) FindByID(ctx context.Context, id int64) (*PostRequest, []PostRequestDetail, error) {
	var header PostRequest
	if err := r.db.WithContext(ctx).Preload("Files").First(&header, id).Error; err != nil {
		return nil, nil, err
	}
	var details []PostRequestDetail
	if err := r.db.WithContext(ctx).
		Where("request_code = ? AND branch_code = ? AND terminal_code = ?",
			header.RequestCode, *header.BranchCode, *header.TerminalCode).
		Order("sequence_number ASC").
		Find(&details).Error; err != nil {
		return nil, nil, err
	}
	return &header, details, nil
}

// FindDetailsByRequestCode retrieves all detail rows for a given request_code.
func (r *postRequestRepository) FindDetailsByRequestCode(ctx context.Context, requestCode string, branchCode, terminalCode int) ([]PostRequestDetail, error) {
	var details []PostRequestDetail
	err := r.db.WithContext(ctx).
		Where("request_code = ? AND branch_code = ? AND terminal_code = ?", requestCode, branchCode, terminalCode).
		Order("sequence_number ASC").
		Find(&details).Error
	return details, err
}

// UpdateHeader updates only the header fields using selective update.
func (r *postRequestRepository) UpdateHeader(ctx context.Context, id int64, header *PostRequest) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Updates(header).Error
}

// ReplaceDetails deletes old detail rows and inserts the new ones atomically.
func (r *postRequestRepository) ReplaceDetails(ctx context.Context, requestCode string, branchCode, terminalCode int, details []PostRequestDetail) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where(
			"request_code = ? AND branch_code = ? AND terminal_code = ?",
			requestCode, branchCode, terminalCode,
		).Delete(&PostRequestDetail{}).Error; err != nil {
			return err
		}
		if len(details) > 0 {
			if err := tx.CreateInBatches(details, 100).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
// ReplaceFiles deletes old file links and inserts the new ones atomically.
func (r *postRequestRepository) ReplaceFiles(ctx context.Context, headerID int64, files []PostRequestFile) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("header_id = ?", headerID).Delete(&PostRequestFile{}).Error; err != nil {
			return err
		}
		for i := range files {
			files[i].HeaderID = headerID
		}
		if len(files) > 0 {
			if err := tx.Create(&files).Error; err != nil {
				return err
			}
		}
		return nil
	})
}


// Delete removes the header. Detail rows are expected to be cascade-deleted
// or cleaned up by the service layer before calling this.
func (r *postRequestRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var header PostRequest
		if err := tx.First(&header, id).Error; err != nil {
			return err
		}
		// Delete details first
		if err := tx.Where(
			"request_code = ? AND branch_code = ? AND terminal_code = ?",
			header.RequestCode, header.BranchCode, header.TerminalCode,
		).Delete(&PostRequestDetail{}).Error; err != nil {
			return err
		}
		return tx.Delete(&PostRequest{}, id).Error
	})
}

// Search uses the generic native pagination helper.
func (r *postRequestRepository) Search(ctx context.Context, param helper.PaginationQuery) ([]PostRequest, helper.PaginationMeta, error) {
	config := helper.NativePaginationConfig{
		TableName: "plan.post_requests",
		SelectColumns: []string{
			"id", "branch_code", "terminal_code", "branch_name", "terminal_name",
			"ppk_number", "schedule_id", "schedule_code", "vessel_code", "vessel_name", "vessel_type", "voyage_type",
			"agent_name", "request_code", "request_date",
			"pbm_code", "pbm_name", "no_bc11", "date_bc11",
			"description", "status", "plan_status",
			"ref_number", "ref_date", "ref1", "ref2", "val1", "val2",
			"total_manifest", "billable_code", "billable_name",
			"vessel_code_dst", "vessel_name_dst",
			"activity_code", "activity_name", "to_ppk_number",
			"approval_date", "creation_date", "creation_by",
			"last_updated_date", "last_updated_by",
		},
		SearchColumns: []string{
			"request_code", "ppk_number", "vessel_name", "pbm_name",
			"billable_name", "agent_name",
		},
		FilterableColumns: map[string]string{
			"branch_code":   "branch_code",
			"terminal_code": "terminal_code",
			"status":        "status",
			"plan_status":   "plan_status",
			"vessel_code":   "vessel_code",
			"vessel_name":   "vessel_name",
			"pbm_name":      "pbm_name",
			"request_code":  "request_code",
			"ppk_number":    "ppk_number",
			"voyage_type":   "voyage_type",
			"request_date":  "request_date",
			"total_manifest": "total_manifest",
			"activity_code": "activity_code",
		},
		SortableColumns: map[string]string{
			"id":           "id",
			"request_code": "request_code",
			"request_date": "request_date",
			"vessel_name":  "vessel_name",
			"status":       "status",
		},
		DefaultSortBy:    "request_date",
		DefaultSortOrder: "DESC",
		MaxLimit:         200,
	}

	var rows []PostRequest
	meta, err := helper.GetDynamicPaginatedNativeData(r.db.WithContext(ctx), config, param, &rows)
	return rows, meta, err
}

// GetStats returns quick dashboard counters.
func (r *postRequestRepository) GetStats(ctx context.Context, branchCode, terminalCode int) (*PostRequestStatsResponse, error) {
	var stats PostRequestStatsResponse

	base := r.db.WithContext(ctx).Model(&PostRequest{})
	if branchCode > 0 {
		base = base.Where("branch_code = ?", branchCode)
	}
	if terminalCode > 0 {
		base = base.Where("terminal_code = ?", terminalCode)
	}

	queries := []struct {
		cond  interface{}
		dest  *int64
		label string
	}{
		{nil, &stats.Total, "total"},
		{"status = 0", &stats.Pending, "pending"},
		{fmt.Sprintf("status = %d", 1), &stats.Approved, "approved"},
		{fmt.Sprintf("status = %d", 2), &stats.Rejected, "rejected"},
	}

	for _, q := range queries {
		s := base.Session(&gorm.Session{})
		if q.cond != nil {
			s = s.Where(q.cond)
		}
		if err := s.Count(q.dest).Error; err != nil {
			return nil, fmt.Errorf("count %s: %w", q.label, err)
		}
	}
	return &stats, nil
}

func (r *postRequestRepository) UpdateStatus(ctx context.Context, id int64, status int, remarks string, updatedBy string) error {
	updates := map[string]interface{}{
		"status":            status,
		"description":       remarks,
		"last_updated_by":   updatedBy,
		"last_updated_date": gorm.Expr("NOW()"),
	}

	if status == 1 {
		updates["approval_date"] = gorm.Expr("NOW()")
	}

	return r.db.WithContext(ctx).Model(&PostRequest{}).Where("id = ?", id).Updates(updates).Error
}

func (r *postRequestRepository) SearchVesselSchedule(ctx context.Context, param helper.PaginationQuery) ([]PostVesselSchedule, helper.PaginationMeta, error) {
	config := helper.NativePaginationConfig{
		TableName: "plan.post_vessel_schedules",
		SelectColumns: []string{
			"id", "branch_code", "terminal_code", "branch_name", "terminal_name",
			"vessel_name", "vessel_code", "vessel_type", "voyage_number", "voyage_type", "grt", "loa",
			"schedule_code", "pkk_number",

			"agency_name", "port_agent", "emergency_contact",
			"origin_port_code", "origin_port_name", "destination_port_code", "destination_port_name",
			"discharge_port_code", "discharge_port_name", "assigned_berth_name",
			"dock_id", "dock_code", "dock_name", "berth_code", "berth_name",
			"berth_position", "position_range", "eta", "etb", "etc", "etd",
			"status", "creation_date", "creation_by", "last_updated_date", "last_updated_by",
		},
		SearchColumns: []string{
			"vessel_name", "vessel_code", "voyage_number", "pkk_number", "agency_name", "port_agent",
		},
		FilterableColumns: map[string]string{
			"branch_code":   "branch_code",
			"terminal_code": "terminal_code",
			"status":        "status",
		},
		SortableColumns: map[string]string{
			"id":          "id",
			"vessel_name": "vessel_name",
			"eta":         "eta",
		},
		DefaultSortBy:    "eta",
		DefaultSortOrder: "DESC",
		MaxLimit:         200,
	}

	var rows []PostVesselSchedule
	meta, err := helper.GetDynamicPaginatedNativeData(r.db.WithContext(ctx), config, param, &rows)
	return rows, meta, err
}

func (r *postRequestRepository) FindVesselScheduleByID(ctx context.Context, id int64) (*PostVesselSchedule, error) {
	var row PostVesselSchedule
	err := r.db.WithContext(ctx).First(&row, id).Error
	return &row, err
}
