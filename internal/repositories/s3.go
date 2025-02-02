package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Repository interface {
	GeneratePresignedUploadURL(ctx context.Context, bucketName, objectKey string) (string, error)
}

type s3RepositoryImpl struct {
	client *s3.PresignClient
}

func NewS3RepositoryImpl(client *s3.PresignClient) *s3RepositoryImpl {
	return &s3RepositoryImpl{client: client}
}

func (r *s3RepositoryImpl) GeneratePresignedUploadURL(ctx context.Context, bucketName, objectKey string) (string, error) {
	req, err := r.client.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(120 * int64(time.Second))
	})

	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return req.URL, nil
}
