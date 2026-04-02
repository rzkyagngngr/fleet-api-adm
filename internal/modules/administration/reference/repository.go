package reference

import "gorm.io/gorm"

type ReferenceRepository interface {
	GetAllHeaders() ([]PosmReference, error)
	GetHeaderByIdRefFile(idRefFile string) (PosmReference, error)
	GetDetailsByIdRefFile(idRefFile string) ([]PosmReferenceD, error)
	SaveHeaderAndDetails(header *PosmReference, details []PosmReferenceD) error
	DeleteHeaderAndDetails(idRefFile string) error
}

type referenceRepository struct{ db *gorm.DB }

func NewReferenceRepository(db *gorm.DB) ReferenceRepository { return &referenceRepository{db: db} }
func (r *referenceRepository) GetAllHeaders() ([]PosmReference, error) {
	var h []PosmReference
	err := r.db.Preload("Details").Find(&h).Error
	return h, err
}
func (r *referenceRepository) GetHeaderByIdRefFile(id string) (PosmReference, error) {
	var h PosmReference
	err := r.db.Where("id_ref_file = ?", id).First(&h).Error
	return h, err
}
func (r *referenceRepository) GetDetailsByIdRefFile(id string) ([]PosmReferenceD, error) {
	var d []PosmReferenceD
	err := r.db.Where("id_ref_file = ?", id).Find(&d).Error
	return d, err
}
func (r *referenceRepository) SaveHeaderAndDetails(header *PosmReference, details []PosmReferenceD) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(header).Error; err != nil {
			return err
		}
		if err := tx.Where("id_ref_file = ?", header.IdRefFile).Delete(&PosmReferenceD{}).Error; err != nil {
			return err
		}
		if len(details) > 0 {
			if err := tx.Create(&details).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
func (r *referenceRepository) DeleteHeaderAndDetails(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id_ref_file = ?", id).Delete(&PosmReferenceD{}).Error; err != nil {
			return err
		}
		if err := tx.Where("id_ref_file = ?", id).Delete(&PosmReference{}).Error; err != nil {
			return err
		}
		return nil
	})
}
