package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"junk-journal-board/internal/utils"

	"github.com/google/uuid"
)

type UploadService struct{}

func NewUploadService() *UploadService {
	return &UploadService{}
}

// ValidateFile validates the uploaded file type and size
func (s *UploadService) ValidateFile(fileHeader *multipart.FileHeader) error {
	// Check file size (10MB limit)
	const maxFileSize = 10 * 1024 * 1024 // 10MB in bytes
	if fileHeader.Size > maxFileSize {
		return utils.NewValidationError("File size exceeds 10MB limit")
	}

	// Check file type by extension
	filename := strings.ToLower(fileHeader.Filename)
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".gif"}

	isValidExtension := false
	for _, ext := range allowedExtensions {
		if strings.HasSuffix(filename, ext) {
			isValidExtension = true
			break
		}
	}

	if !isValidExtension {
		return utils.NewValidationError("File type not allowed. Only JPG, PNG, and GIF files are supported")
	}

	return nil
}

// SaveFile saves the uploaded file to the organized directory structure
func (s *UploadService) SaveFile(fileHeader *multipart.FileHeader, boardID uuid.UUID) (string, error) {
	// Create directory structure: uploads/boards/{boardID}/
	uploadDir := filepath.Join("uploads", "boards", boardID.String())
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Generate unique filename with UUID
	fileExt := filepath.Ext(fileHeader.Filename)
	fileName := fmt.Sprintf("%s%s", uuid.New().String(), fileExt)
	filePath := filepath.Join(uploadDir, fileName)

	// Open the uploaded file
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Create the destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy file contents
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// Return the full public URL with backend host
	backendURL := os.Getenv("BACKEND_URL")
	if backendURL == "" {
		backendURL = "http://localhost:8080" // Default backend URL
	}
	publicURL := fmt.Sprintf("%s/uploads/boards/%s/%s", backendURL, boardID.String(), fileName)
	return publicURL, nil
}

// DeleteFile removes a file from the storage
func (s *UploadService) DeleteFile(filePath string) error {
	// Convert public URL to file system path
	if strings.HasPrefix(filePath, "/uploads/") {
		filePath = strings.TrimPrefix(filePath, "/")
	}

	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// GetFileExtension returns the file extension from filename
func (s *UploadService) GetFileExtension(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}

// GetMimeType returns the MIME type based on file extension
func (s *UploadService) GetMimeType(filename string) string {
	ext := s.GetFileExtension(filename)
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	default:
		return "application/octet-stream"
	}
}
