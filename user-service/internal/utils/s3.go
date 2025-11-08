package utils

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func UploadToS3(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("S3_REGION")),
	)
	if err != nil {
		return "", fmt.Errorf("failed to load AWS config: %w", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	// Read the uploaded file into memory
	fileBytes := new(bytes.Buffer)
	_, err = fileBytes.ReadFrom(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	fileName := fmt.Sprintf("uploads/%d_%s", time.Now().Unix(), filepath.Base(fileHeader.Filename))

	bucket := os.Getenv("S3_BUCKET_NAME")
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        &bucket,
		Key:           &fileName,
		Body:          bytes.NewReader(fileBytes.Bytes()),
		ContentLength: aws.Int64(int64(fileBytes.Len())),
		ContentType:   aws.String(fileHeader.Header.Get("Content-Type")),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, os.Getenv("S3_REGION"), fileName)
	return url, nil
}
