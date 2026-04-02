package reference

import "time"

type ReferenceService interface {
	GetAllHeaders() ([]PosmReference, error)
	GetHeaderWithDetails(idRefFile string) (PosmReference, error)
	SaveReference(header PosmReference, details []PosmReferenceD) error
	DeleteReference(idRefFile string) error
}

type referenceService struct{ repo ReferenceRepository }

func NewReferenceService(repo ReferenceRepository) ReferenceService {
	return &referenceService{repo: repo}
}
func (s *referenceService) GetAllHeaders() ([]PosmReference, error) { return s.repo.GetAllHeaders() }
func (s *referenceService) GetHeaderWithDetails(idRefFile string) (PosmReference, error) {
	h, err := s.repo.GetHeaderByIdRefFile(idRefFile)
	if err != nil {
		return h, err
	}
	d, err := s.repo.GetDetailsByIdRefFile(idRefFile)
	if err != nil {
		return h, err
	}
	h.Details = d
	return h, nil
}
func (s *referenceService) SaveReference(header PosmReference, details []PosmReferenceD) error {
	now := time.Now()
	header.LastUpdatedDate = &now
	for i := range details {
		details[i].IdRefFile = header.IdRefFile
		details[i].LastUpdatedDate = &now
		details[i].BranchCode = header.BranchCode
		details[i].TerminalCode = header.TerminalCode
		details[i].ProgramName = header.ProgramName
		details[i].LevelAkses = header.LevelAkses
	}
	return s.repo.SaveHeaderAndDetails(&header, details)
}
func (s *referenceService) DeleteReference(idRefFile string) error {
	return s.repo.DeleteHeaderAndDetails(idRefFile)
}
