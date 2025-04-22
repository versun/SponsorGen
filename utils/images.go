package utils

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// DownloadImage downloads an image from a URL and returns it as a base64-encoded data URI
func DownloadImage(imageURL string, cacheDir string) (string, error) {
	// Use cached version if available
	cacheKey := getImageCacheKey(imageURL)
	cachePath := filepath.Join(cacheDir, cacheKey)
	
	// Check if cached image exists and is not too old (e.g., 24 hours)
	if cachedImage, err := os.ReadFile(cachePath); err == nil {
		fileInfo, err := os.Stat(cachePath)
		if err == nil && time.Since(fileInfo.ModTime()) < 24*time.Hour {
			return string(cachedImage), nil
		}
	}
	
	// If it's a data URI already, return it as is
	if strings.HasPrefix(imageURL, "data:") {
		return imageURL, nil
	}
	
	// If it's a local file, read it
	if strings.HasPrefix(imageURL, "./") || strings.HasPrefix(imageURL, "/") {
		data, err := os.ReadFile(imageURL)
		if err != nil {
			return "", fmt.Errorf("error reading local image: %w", err)
		}
		
		// Determine MIME type based on file extension
		mimeType := getMimeType(imageURL)
		dataURI := fmt.Sprintf("data:%s;base64,%s", mimeType, base64.StdEncoding.EncodeToString(data))
		
		// Cache the result
		if err := os.MkdirAll(cacheDir, 0755); err == nil {
			_ = os.WriteFile(cachePath, []byte(dataURI), 0644)
		}
		
		return dataURI, nil
	}
	
	// Download the image
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("error downloading image: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error downloading image, status code: %d", resp.StatusCode)
	}
	
	// Read the image data
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading image data: %w", err)
	}
	
	// Get content type from response
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = getMimeType(imageURL)
	}
	
	// Create data URI
	dataURI := fmt.Sprintf("data:%s;base64,%s", contentType, base64.StdEncoding.EncodeToString(data))
	
	// Cache the result
	if err := os.MkdirAll(cacheDir, 0755); err == nil {
		_ = os.WriteFile(cachePath, []byte(dataURI), 0644)
	}
	
	return dataURI, nil
}

// getImageCacheKey generates a cache key for an image URL
func getImageCacheKey(url string) string {
	// Replace non-alphanumeric characters with underscores
	clean := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return '_'
	}, url)
	
	return fmt.Sprintf("img_%s.txt", clean)
}

// getMimeType determines the MIME type based on file extension
func getMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	default:
		return "image/png" // Default to PNG
	}
}