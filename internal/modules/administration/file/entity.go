package file

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FileRecord represents the adm.posm_files table
type FileRecord struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	FileHash      string         `gorm:"column:file_hash;not null;index"`
	FileName      string         `gorm:"column:file_name;not null"`
	FileExtension string         `gorm:"column:file_extension;not null"`
	FileSize      int64          `gorm:"column:file_size;not null"`
	MimeType      string         `gorm:"column:mime_type;not null"`
	Provider      string         `gorm:"column:provider;default:s3"`
	BucketName    string         `gorm:"column:bucket_name;not null"`
	FilePath      string         `gorm:"column:file_path;not null"`
	IsPublic      bool           `gorm:"column:is_public;default:false"`
	IsCommitted   bool           `gorm:"column:is_committed;default:false"`
	IsEncrypted   bool           `gorm:"column:is_encrypted;default:false"`
	S3VersionID   string         `gorm:"column:s3_version_id"`
	StorageClass  string         `gorm:"column:storage_class;default:STANDARD"`
	Metadata        map[string]any `gorm:"column:metadata;type:jsonb" json:"metadata"`
	CreationDate    time.Time      `gorm:"column:creation_date;default:CURRENT_TIMESTAMP" json:"creation_date"`
	CreationBy      string         `gorm:"column:creation_by"                             json:"creation_by"`
	LastUpdatedDate *time.Time     `gorm:"column:last_updated_date"                       json:"last_updated_date"`
	LastUpdatedBy   string         `gorm:"column:last_updated_by"                         json:"last_updated_by"`
	ProgramName     string         `gorm:"column:program_name"                           json:"program_name"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index"                         json:"-"`
}

// TableName sets the insert table name for this struct type
func (FileRecord) TableName() string {
	return "adm.posm_files"
}
