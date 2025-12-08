package utils

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/zercle/gofiber-skelton/internal/datasources"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/valyala/fastjson"
)

var (
	// Session storage
	SessStore *session.Store
	// parser pool
	JsonParserPool *fastjson.ParserPool
)

func init() {
	if JsonParserPool == nil {
		JsonParserPool = new(fastjson.ParserPool)
	}
}

// UploadFileToMinio uploads a file to MinIO and returns the permanent public URL
// Returns: (url string, error)
func UploadFileToMinio(ctx context.Context, file *multipart.FileHeader) (string, error) {
	minioClient := datasources.GetMinioClient()
	if minioClient == nil {
		return "", fmt.Errorf("MinIO client not initialized")
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Generate unique filename with original extension
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), uuid.New().String(), ext)

	// Get content type
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Upload to MinIO
	_, err = minioClient.Client.PutObject(
		ctx,
		minioClient.Bucket,
		filename,
		src,
		file.Size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload to MinIO: %w", err)
	}

	url := fmt.Sprintf("https://s3.mysterchat.com/%s/%s", minioClient.Bucket, filename)

	return url, nil
}

// UploadFileStreamToMinio uploads a file from io.Reader to MinIO
// Returns: (url string, error)
func UploadFileStreamToMinio(ctx context.Context, reader io.Reader, filename string, contentType string, size int64) (string, error) {
	minioClient := datasources.GetMinioClient()
	if minioClient == nil {
		return "", fmt.Errorf("MinIO client not initialized")
	}

	// Generate unique filename if not provided
	if filename == "" {
		filename = fmt.Sprintf("%d_%s", time.Now().Unix(), uuid.New().String())
	} else {
		// Add unique prefix to avoid name collision
		ext := filepath.Ext(filename)
		name := strings.TrimSuffix(filename, ext)
		filename = fmt.Sprintf("%d_%s%s", time.Now().Unix(), name, ext)
	}

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Upload to MinIO
	_, err := minioClient.Client.PutObject(
		ctx,
		minioClient.Bucket,
		filename,
		reader,
		size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload to MinIO: %w", err)
	}

	url := fmt.Sprintf("https://s3.mysterchat.com/%s/%s", minioClient.Bucket, filename)

	return url, nil
}

// DeleteFileFromMinio deletes a file from MinIO using its URL or filename
func DeleteFileFromMinio(ctx context.Context, urlOrFilename string) error {
	minioClient := datasources.GetMinioClient()
	if minioClient == nil {
		return fmt.Errorf("MinIO client not initialized")
	}

	// Extract filename from URL if necessary
	filename := urlOrFilename
	if strings.Contains(urlOrFilename, "/") {
		parts := strings.Split(urlOrFilename, "/")
		filename = parts[len(parts)-1]
	}

	err := minioClient.Client.RemoveObject(ctx, minioClient.Bucket, filename, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete from MinIO: %w", err)
	}

	return nil
}

// GetFileInfoFromMinio gets file information from MinIO
func GetFileInfoFromMinio(ctx context.Context, filename string) (minio.ObjectInfo, error) {
	minioClient := datasources.GetMinioClient()
	if minioClient == nil {
		return minio.ObjectInfo{}, fmt.Errorf("MinIO client not initialized")
	}

	info, err := minioClient.Client.StatObject(ctx, minioClient.Bucket, filename, minio.StatObjectOptions{})
	if err != nil {
		return minio.ObjectInfo{}, fmt.Errorf("failed to get file info: %w", err)
	}

	return info, nil
}
