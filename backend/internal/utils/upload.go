package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	UploadDir   = "uploads"
	MaxFileSize = 5 << 20 // 5MB
)

var AllowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
}

// UploadFile handles file upload and returns the file path
func UploadFile(file *multipart.FileHeader, subDir string) (string, error) {
	// Check file size
	if file.Size > MaxFileSize {
		return "", fmt.Errorf("file size exceeds maximum limit of %d bytes", MaxFileSize)
	}

	// Check file type
	if !AllowedImageTypes[file.Header.Get("Content-Type")] {
		return "", fmt.Errorf("file type not allowed")
	}

	// Create upload directory if it doesn't exist
	uploadPath := filepath.Join(UploadDir, subDir)
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%s%s",
		time.Now().Format("20060102"),
		uuid.New().String(),
		ext,
	)

	// Save file
	dst := filepath.Join(uploadPath, filename)
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %v", err)
	}
	defer out.Close()

	// Copy the uploaded file to the destination file
	buffer := make([]byte, 1024)
	for {
		n, err := src.Read(buffer)
		if err != nil {
			break
		}
		out.Write(buffer[:n])
	}

	// Return the relative file path
	return filepath.Join(subDir, filename), nil
}

// GetFileURL returns the full URL for a file path
func GetFileURL(path string) string {
	if path == "" {
		return ""
	}
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}
	// TODO: 从配置中获取基础URL
	baseURL := "http://localhost:8080"
	return fmt.Sprintf("%s/uploads/%s", baseURL, path)
}

// DeleteFile deletes a file from the upload directory
func DeleteFile(path string) error {
	if path == "" {
		return nil
	}
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return nil
	}
	return os.Remove(filepath.Join(UploadDir, path))
}
