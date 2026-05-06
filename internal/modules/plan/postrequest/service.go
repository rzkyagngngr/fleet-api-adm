package postrequest

import (
	"context"
	"errors"
	"fmt"
	"omniport-api/internal/helper"
	"omniport-api/internal/modules/administration/file"
	"time"
)

const programName = "ADM_SERVICE"

// ─────────────────────────────────────────────────────────────
// INTERFACE
// ─────────────────────────────────────────────────────────────

type PostRequestService interface {
	Create(ctx context.Context, input *CreatePostRequestInput, branchCode, terminalCode int, branchName, terminalName, createdBy string) (*PostRequestResponse, error)
	GetByID(ctx context.Context, id int64) (*PostRequestResponse, error)
	Update(ctx context.Context, id int64, input *UpdatePostRequestInput, updatedBy string) (*PostRequestResponse, error)
	Delete(ctx context.Context, id int64) error
	Search(ctx context.Context, query helper.PaginationQuery) ([]PostRequestResponse, helper.PaginationMeta, error)
	GetStats(ctx context.Context, branchCode, terminalCode int) (*PostRequestStatsResponse, error)
	UpdateStatus(ctx context.Context, id int64, status int, remarks string, updatedBy string) error

	// Vessel Schedule methods
	SearchVesselSchedule(ctx context.Context, query helper.PaginationQuery) ([]PostVesselScheduleResponse, helper.PaginationMeta, error)
	GetVesselScheduleByID(ctx context.Context, id int64) (*PostVesselScheduleResponse, error)
}

// ─────────────────────────────────────────────────────────────
// IMPLEMENTATION
// ─────────────────────────────────────────────────────────────

type postRequestService struct {
	repo        PostRequestRepository
	fileService file.FileService
}

func NewPostRequestService(repo PostRequestRepository, fileService file.FileService) PostRequestService {
	return &postRequestService{
		repo:        repo,
		fileService: fileService,
	}
}

// generateRequestCode produces a unique code: PR-YYYYMMDD-<nano_suffix>
func generateRequestCode() string {
	now := time.Now()
	return fmt.Sprintf("PR-%s-%d", now.Format("20060102"), now.UnixNano()%1_000_000)
}

// Create validates, builds, and persists a new cargo service request.
func (s *postRequestService) Create(
	ctx context.Context,
	input *CreatePostRequestInput,
	branchCode, terminalCode int,
	branchName, terminalName, createdBy string,
) (*PostRequestResponse, error) {

	if input.VesselCode == "" || input.VesselName == "" {
		return nil, errors.New("vessel_code and vessel_name are required")
	}
	if len(input.Details) == 0 {
		return nil, errors.New("at least one manifest detail (details) is required")
	}

	now := time.Now()
	requestCode := generateRequestCode()
	statusPending := 0
	planStatusDraft := 0

	header := &PostRequest{
		BranchCode:      &branchCode,
		TerminalCode:    &terminalCode,
		BranchName:      branchName,
		TerminalName:    terminalName,
		PPKNumber:       input.PPKNumber,
		VesselCode:      input.VesselCode,
		VesselName:      input.VesselName,
		VesselType:      input.VesselType,
		VoyageType:      input.VoyageType,
		AgentName:       input.AgentName,
		RequestCode:     requestCode,
		RequestDate:     input.RequestDate,
		PBMCode:         input.PBMCode,
		PBMName:         input.PBMName,
		NoBC11:          input.NoBC11,
		DateBC11:        input.DateBC11,
		Description:     input.Description,
		Status:          &statusPending,
		PlanStatus:      &planStatusDraft,
		ProgramName:     programName,
		RefNumber:       input.RefNumber,
		RefDate:         input.RefDate,
		Ref1:            input.Ref1,
		Ref2:            input.Ref2,
		Val1:            input.Val1,
		Val2:            input.Val2,
		TotalManifest:   input.TotalManifest,
		BillableCode:    input.BillableCode,
		BillableName:    input.BillableName,
		VesselCodeDst:   input.VesselCodeDst,
		VesselNameDst:   input.VesselNameDst,
		ActivityCode:    input.ActivityCode,
		ActivityName:    input.ActivityName,
		ToPPKNumber:     input.ToPPKNumber,
		CreationDate:    &now,
		CreationBy:      createdBy,
		LastUpdatedDate: &now,
		LastUpdatedBy:   createdBy,
	}

	details := buildDetails(input.Details, requestCode, branchCode, terminalCode, branchName, terminalName, createdBy, now)
	files := buildFiles(input.Attachments)

	if err := s.repo.Create(ctx, header, details, files); err != nil {
		return nil, fmt.Errorf("create post_request: %w", err)
	}

	header.Files = files // For response
	res := header.ToResponse(details)
	s.fillFileURLs(ctx, &res)
	return &res, nil
}

func (s *postRequestService) fillFileURLsBulk(ctx context.Context, results []*PostRequestResponse) {
	if s.fileService == nil || len(results) == 0 {
		return
	}

	var allAttachments []*file.FileAttachment
	for _, res := range results {
		if res == nil { continue }
		for i := range res.Attachments {
			allAttachments = append(allAttachments, &res.Attachments[i].FileAttachment)
		}
	}

	_ = s.fileService.EnrichAttachments(ctx, allAttachments)
}

func (s *postRequestService) fillFileURLs(ctx context.Context, res *PostRequestResponse) {
	s.fillFileURLsBulk(ctx, []*PostRequestResponse{res})
}

// GetByID fetches a request with all its manifest lines.
func (s *postRequestService) GetByID(ctx context.Context, id int64) (*PostRequestResponse, error) {
	header, details, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("post_request not found: %w", err)
	}
	res := header.ToResponse(details)
	s.fillFileURLs(ctx, &res)
	return &res, nil
}

// Update patches header fields and, if details are provided, replaces manifest lines.
func (s *postRequestService) Update(ctx context.Context, id int64, input *UpdatePostRequestInput, updatedBy string) (*PostRequestResponse, error) {
	header, _, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("post_request not found")
	}

	now := time.Now()

	// Patch only non-zero/non-empty fields
	if input.PPKNumber != "" {
		header.PPKNumber = input.PPKNumber
	}
	if input.VesselCode != "" {
		header.VesselCode = input.VesselCode
	}
	if input.VesselName != "" {
		header.VesselName = input.VesselName
	}
	if input.VesselType != "" {
		header.VesselType = input.VesselType
	}
	if input.VoyageType != "" {
		header.VoyageType = input.VoyageType
	}
	if input.AgentName != "" {
		header.AgentName = input.AgentName
	}
	if input.RequestDate != nil {
		header.RequestDate = *input.RequestDate
	}
	if input.PBMCode != "" {
		header.PBMCode = input.PBMCode
	}
	if input.PBMName != "" {
		header.PBMName = input.PBMName
	}
	if input.NoBC11 != "" {
		header.NoBC11 = input.NoBC11
	}
	if input.DateBC11 != nil {
		header.DateBC11 = input.DateBC11
	}
	if input.Description != "" {
		header.Description = input.Description
	}
	if input.RefNumber != "" {
		header.RefNumber = input.RefNumber
	}
	if input.RefDate != nil {
		header.RefDate = input.RefDate
	}
	if input.Ref1 != "" {
		header.Ref1 = input.Ref1
	}
	if input.Ref2 != "" {
		header.Ref2 = input.Ref2
	}
	if input.Val1 != nil {
		header.Val1 = input.Val1
	}
	if input.Val2 != nil {
		header.Val2 = input.Val2
	}
	if input.TotalManifest != nil {
		header.TotalManifest = input.TotalManifest
	}
	if input.BillableCode != "" {
		header.BillableCode = input.BillableCode
	}
	if input.BillableName != "" {
		header.BillableName = input.BillableName
	}
	if input.VesselCodeDst != "" {
		header.VesselCodeDst = input.VesselCodeDst
	}
	if input.VesselNameDst != "" {
		header.VesselNameDst = input.VesselNameDst
	}
	if input.ActivityCode != "" {
		header.ActivityCode = input.ActivityCode
	}
	if input.ActivityName != "" {
		header.ActivityName = input.ActivityName
	}
	if input.ToPPKNumber != "" {
		header.ToPPKNumber = input.ToPPKNumber
	}

	header.LastUpdatedDate = &now
	header.LastUpdatedBy = updatedBy

	if err := s.repo.UpdateHeader(ctx, id, header); err != nil {
		return nil, fmt.Errorf("update post_request header: %w", err)
	}

	var finalDetails []PostRequestDetail
	if len(input.Details) > 0 {
		newDetails := buildDetails(input.Details, header.RequestCode, *header.BranchCode, *header.TerminalCode,
			header.BranchName, header.TerminalName, updatedBy, now)
		if err := s.repo.ReplaceDetails(ctx, header.RequestCode, *header.BranchCode, *header.TerminalCode, newDetails); err != nil {
			return nil, fmt.Errorf("replace post_request_d: %w", err)
		}
		finalDetails = newDetails
	} else {
		finalDetails, err = s.repo.FindDetailsByRequestCode(ctx, header.RequestCode, *header.BranchCode, *header.TerminalCode)
		if err != nil {
			return nil, err
		}
	}

	if len(input.Attachments) > 0 {
		newFiles := buildFiles(input.Attachments)
		if err := s.repo.ReplaceFiles(ctx, id, newFiles); err != nil {
			return nil, fmt.Errorf("replace post_request_f: %w", err)
		}
		header.Files = newFiles
	}

	res := header.ToResponse(finalDetails)
	s.fillFileURLs(ctx, &res)
	return &res, nil
}

// Delete removes the header and all associated detail rows.
func (s *postRequestService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// Search returns paginated list with metadata.
func (s *postRequestService) Search(ctx context.Context, query helper.PaginationQuery) ([]PostRequestResponse, helper.PaginationMeta, error) {
	rows, meta, err := s.repo.Search(ctx, query)
	if err != nil {
		return nil, meta, err
	}
	
	res := make([]PostRequestResponse, len(rows))
	ptrs := make([]*PostRequestResponse, len(rows))
	
	for i, h := range rows {
		// Details not loaded in list view for performance; use GetByID for full detail.
		res[i] = h.ToResponse(nil)
		ptrs[i] = &res[i]
	}

	// Fill technical metadata and URLs for the entire batch
	s.fillFileURLsBulk(ctx, ptrs)

	return res, meta, nil
}

// GetStats returns aggregated counts.
func (s *postRequestService) GetStats(ctx context.Context, branchCode, terminalCode int) (*PostRequestStatsResponse, error) {
	return s.repo.GetStats(ctx, branchCode, terminalCode)
}

func (s *postRequestService) UpdateStatus(ctx context.Context, id int64, status int, remarks string, updatedBy string) error {
	_, _, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("permohonan tidak ditemukan")
	}
	return s.repo.UpdateStatus(ctx, id, status, remarks, updatedBy)
}

func (s *postRequestService) SearchVesselSchedule(ctx context.Context, query helper.PaginationQuery) ([]PostVesselScheduleResponse, helper.PaginationMeta, error) {
	rows, meta, err := s.repo.SearchVesselSchedule(ctx, query)
	if err != nil {
		return nil, meta, err
	}
	var res []PostVesselScheduleResponse
	for _, h := range rows {
		res = append(res, h.ToResponse())
	}
	return res, meta, nil
}

func (s *postRequestService) GetVesselScheduleByID(ctx context.Context, id int64) (*PostVesselScheduleResponse, error) {
	row, err := s.repo.FindVesselScheduleByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("vessel_schedule not found: %w", err)
	}
	res := row.ToResponse()
	return &res, nil
}


// ─────────────────────────────────────────────────────────────
// HELPERS
// ─────────────────────────────────────────────────────────────

func buildDetails(
	inputs []PostRequestDetailInput,
	requestCode string,
	branchCode, terminalCode int,
	branchName, terminalName, createdBy string,
	now time.Time,
) []PostRequestDetail {
	details := make([]PostRequestDetail, 0, len(inputs))
	for i, d := range inputs {
		seq := i + 1
		if d.SequenceNumber != nil {
			seq = *d.SequenceNumber
		}
		details = append(details, PostRequestDetail{
			BranchCode:          branchCode,
			TerminalCode:        terminalCode,
			BranchName:          branchName,
			TerminalName:        terminalName,
			RequestCode:         requestCode,
			SequenceNumber:      &seq,
			StackingType:        d.StackingType,
			CargoCode:           d.CargoCode,
			CargoName:           d.CargoName,
			CargoUnit:           d.CargoUnit,
			Total:               d.Total,
			QuantityMT:          d.QuantityMT,
			Quantity:            d.Quantity,
			CargoNature:         d.CargoNature,
			CargoNatureDesc:     d.CargoNatureDesc,
			CargoPackaging:      d.CargoPackaging,
			Stowage:             d.Stowage,
			PlannedDate:         d.PlannedDate,
			WarehouseID:         d.WarehouseID,
			BLAWBNumber:         d.BLAWBNumber,
			BLAWBDate:           d.BLAWBDate,
			Description:         d.Description,
			PackageCount:        d.PackageCount,
			OriginPortCode:      d.OriginPortCode,
			DestinationPortCode: d.DestinationPortCode,
			OriginPortName:      d.OriginPortName,
			DestinationPortName: d.DestinationPortName,
			StorageReference:    d.StorageReference,
			StorageStackDate:    d.StorageStackDate,
			WarehouseDetailID:   d.WarehouseDetailID,
			WarehouseDetailName: d.WarehouseDetailName,
			WarehouseName:       d.WarehouseName,
			ConsigneeCode:       d.ConsigneeCode,
			ConsigneeName:       d.ConsigneeName,
			CreationDate:        &now,
			CreationBy:          createdBy,
			LastUpdatedDate:     &now,
			LastUpdatedBy:       createdBy,
			ProgramName:         programName,
		})
	}
	return details
}

func buildFiles(inputs []AttachmentInput) []PostRequestFile {
	files := make([]PostRequestFile, 0, len(inputs))
	for _, a := range inputs {
		files = append(files, PostRequestFile{
			FileID:  a.FileID,
			DocType: a.DocType,
			DocName: a.DocName,
		})
	}
	return files
}
