package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

// SaveUploadedFile saves the uploaded file to the public/images directory and returns the relative path
func SaveUploadedFile(file multipart.File, header *multipart.FileHeader, publicDir string) (string, error) {
	defer file.Close()
	// Ensure the images directory exists
	imagesDir := filepath.Join(publicDir, "images")
	if err := os.MkdirAll(imagesDir, os.ModePerm); err != nil {
		return "", err
	}
	// Generate a unique filename
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dstPath := filepath.Join(imagesDir, filename)

	out, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		return "", err
	}
	// Return the relative path (for storing in DB)
	relPath := filepath.ToSlash(filepath.Join("images", filename))
	return relPath, nil
}
