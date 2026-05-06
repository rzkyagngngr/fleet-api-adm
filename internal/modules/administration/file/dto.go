package file

import "github.com/google/uuid"

type UploadSignatureRequest struct {
	FileName string `json:"file_name" binding:"required"`
	FileSize int64  `json:"file_size" binding:"required"`
	FileHash string `json:"file_hash" binding:"required"`
	MimeType string `json:"mime_type" binding:"required"`
}

type UploadSignatureResponse struct {
	FileID      uuid.UUID `json:"file_id"`
	UploadURL   string    `json:"upload_url,omitempty"`
	IsDuplicate bool      `json:"is_duplicate"`
}

type FileResponse struct {
	ID            uuid.UUID `json:"id"`
	FileName      string    `json:"file_name"`
	FileExtension string    `json:"file_extension"`
	FileSize      int64     `json:"file_size"`
	MimeType      string    `json:"mime_type"`
	CreatedAt     string    `json:"created_at"`
	DownloadURL   string    `json:"download_url,omitempty"`
}

// FileAttachment is a reusable struct for embedding file info into other module responses.
type FileAttachment struct {
	FileID        string `json:"file_id"`
	URL           string `json:"url,omitempty"`
	FileName      string `json:"file_name,omitempty"`
	FileSize      int64  `json:"file_size,omitempty"`
	MimeType      string `json:"mime_type,omitempty"`
	FileExtension string `json:"file_extension,omitempty"`
}
