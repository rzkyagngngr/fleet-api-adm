package file

import (
	"context"
	"fmt"
	"omniport-api/internal/config"
	"omniport-api/internal/helper"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type FileService interface {
	GetUploadSignature(ctx context.Context, req UploadSignatureRequest, userID string) (UploadSignatureResponse, error)
	CommitUpload(ctx context.Context, id uuid.UUID) error
	GetDownloadURL(ctx context.Context, id uuid.UUID) (string, error)
	GetFileDetail(ctx context.Context, id uuid.UUID) (*FileResponse, error)
	EnrichAttachments(ctx context.Context, attachments []*FileAttachment) error
}

type fileService struct {
	repo    FileRepository
	storage helper.StorageProvider
	cfg     config.StorageConfig
}

func NewFileService(repo FileRepository, storage helper.StorageProvider, cfg config.StorageConfig) FileService {
	return &fileService{repo: repo, storage: storage, cfg: cfg}
}

func (s *fileService) GetUploadSignature(ctx context.Context, req UploadSignatureRequest, userID string) (UploadSignatureResponse, error) {
	// 1. Check for deduplication
	existing, err := s.repo.FindByHash(ctx, req.FileHash)
	if err != nil {
		return UploadSignatureResponse{}, err
	}

	if existing != nil {
		return UploadSignatureResponse{
			FileID:      existing.ID,
			IsDuplicate: true,
		}, nil
	}

	// 2. Generate new record and pre-signed URL
	fileID := uuid.New()
	ext := strings.TrimPrefix(filepath.Ext(req.FileName), ".")
	bucket := s.cfg.S3Bucket
	key := fmt.Sprintf("%s/%s/%s.%s", time.Now().Format("2006/01"), req.FileHash[:8], fileID.String(), ext)

	record := &FileRecord{
		ID:            fileID,
		FileHash:      req.FileHash,
		FileName:      req.FileName,
		FileExtension: ext,
		FileSize:      req.FileSize,
		MimeType:      req.MimeType,
		BucketName:    bucket,
		FilePath:      key,
		CreationBy:      userID,
		ProgramName:     "ADM_SERVICE",
		IsCommitted:   false,
	}

	if err := s.repo.Create(ctx, record); err != nil {
		return UploadSignatureResponse{}, err
	}

	uploadURL, err := s.storage.GeneratePresignedPutURL(ctx, bucket, key, req.MimeType, 15*time.Minute)
	if err != nil {
		return UploadSignatureResponse{}, err
	}

	return UploadSignatureResponse{
		FileID:    fileID,
		UploadURL: uploadURL,
	}, nil
}

func (s *fileService) CommitUpload(ctx context.Context, id uuid.UUID) error {
	file, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	file.IsCommitted = true
	now := time.Now()
	file.LastUpdatedDate = &now

	return s.repo.Update(ctx, file)
}

func (s *fileService) GetDownloadURL(ctx context.Context, id uuid.UUID) (string, error) {
	file, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return "", err
	}

	return s.storage.GeneratePresignedGetURL(ctx, file.BucketName, file.FilePath, 30*time.Minute)
}

func (s *fileService) GetFileDetail(ctx context.Context, id uuid.UUID) (*FileResponse, error) {
	file, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	url, _ := s.GetDownloadURL(ctx, id)

	return &FileResponse{
		ID:            file.ID,
		FileName:      file.FileName,
		FileExtension: file.FileExtension,
		FileSize:      file.FileSize,
		MimeType:      file.MimeType,
		CreatedAt:     file.CreationDate.Format(time.RFC3339),
		DownloadURL:   url,
	}, nil
}

func (s *fileService) EnrichAttachments(ctx context.Context, attachments []*FileAttachment) error {
	if len(attachments) == 0 {
		return nil
	}

	// 1. Collect unique IDs
	var ids []uuid.UUID
	idMap := make(map[string]bool)
	for _, a := range attachments {
		if a == nil || a.FileID == "" {
			continue
		}
		if uid, err := uuid.Parse(a.FileID); err == nil {
			if !idMap[a.FileID] {
				ids = append(ids, uid)
				idMap[a.FileID] = true
			}
		}
	}

	if len(ids) == 0 {
		return nil
	}

	// 2. Fetch metadata from DB
	records, err := s.repo.FindByIDs(ctx, ids)
	if err != nil {
		return err
	}

	metaMap := make(map[string]FileRecord)
	for _, r := range records {
		metaMap[r.ID.String()] = r
	}

	// 3. Populate fields and URLs
	for _, a := range attachments {
		if a == nil {
			continue
		}
		if meta, ok := metaMap[a.FileID]; ok {
			a.FileName = meta.FileName
			a.FileSize = meta.FileSize
			a.MimeType = meta.MimeType
			a.FileExtension = meta.FileExtension

			// Use 1 hour expiry for pre-signed GET URLs
			url, err := s.storage.GeneratePresignedGetURL(ctx, meta.BucketName, meta.FilePath, 1*time.Hour)
			if err == nil {
				a.URL = url
			}
		}
	}

	return nil
}
