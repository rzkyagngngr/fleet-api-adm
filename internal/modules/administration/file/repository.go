package file

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileRepository interface {
	Create(ctx context.Context, file *FileRecord) error
	FindByID(ctx context.Context, id uuid.UUID) (*FileRecord, error)
	FindByIDs(ctx context.Context, ids []uuid.UUID) ([]FileRecord, error)
	FindByHash(ctx context.Context, hash string) (*FileRecord, error)
	Update(ctx context.Context, file *FileRecord) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type fileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{db: db}
}

func (r *fileRepository) Create(ctx context.Context, file *FileRecord) error {
	return r.db.WithContext(ctx).Create(file).Error
}

func (r *fileRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]FileRecord, error) {
	var files []FileRecord
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

func (r *fileRepository) FindByID(ctx context.Context, id uuid.UUID) (*FileRecord, error) {
	var file FileRecord
	if err := r.db.WithContext(ctx).First(&file, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *fileRepository) FindByHash(ctx context.Context, hash string) (*FileRecord, error) {
	var file FileRecord
	if err := r.db.WithContext(ctx).Where("file_hash = ? AND is_committed = true", hash).First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &file, nil
}

func (r *fileRepository) Update(ctx context.Context, file *FileRecord) error {
	return r.db.WithContext(ctx).Save(file).Error
}

func (r *fileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&FileRecord{}, "id = ?", id).Error
}
