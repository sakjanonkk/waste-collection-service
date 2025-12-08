package datasources

import (
	"context"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
)

type MinioClient struct {
	Client   *minio.Client
	Bucket   string
	Endpoint string
	UseSSL   bool
}

var MinioConn *MinioClient

// InitMinio initializes MinIO client connection
func InitMinio() error {
	endpoint := viper.GetString("minio.endpoint")
	accessKey := viper.GetString("minio.access_key")
	secretKey := viper.GetString("minio.secret_key")
	useSSL := viper.GetBool("minio.use_ssl")
	bucket := viper.GetString("minio.bucket")

	if endpoint == "" || accessKey == "" || secretKey == "" {
		return fmt.Errorf("MinIO configuration is incomplete")
	}

	// Initialize MinIO client
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return fmt.Errorf("failed to create MinIO client: %w", err)
	}

	// Create bucket if it doesn't exist
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		log.Printf("Bucket '%s' created successfully", bucket)

		// Set bucket policy to public read (for permanent URLs)
		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {"AWS": ["*"]},
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::%s/*"]
				}
			]
		}`, bucket)

		err = minioClient.SetBucketPolicy(ctx, bucket, policy)
		if err != nil {
			log.Printf("Warning: failed to set bucket policy: %v", err)
		}
	}

	MinioConn = &MinioClient{
		Client:   minioClient,
		Bucket:   bucket,
		Endpoint: endpoint,
		UseSSL:   useSSL,
	}

	log.Printf("MinIO client initialized successfully - Endpoint: %s, Bucket: %s", endpoint, bucket)
	return nil
}

// GetMinioClient returns the initialized MinIO client
func GetMinioClient() *MinioClient {
	return MinioConn
}
