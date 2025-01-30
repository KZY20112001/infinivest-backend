package services

import (
	"context"

	"github.com/KZY20112001/infinivest-backend/internal/dto"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
)

type S3Service interface {
	GeneratePresignedUploadURL(ctx context.Context, dto dto.PresignedUploadUrlRequest) (string, error)
}

type s3ServiceImpl struct {
	repo repositories.S3Repository
}

func NewS3ServiceImpl(repo repositories.S3Repository) *s3ServiceImpl {
	return &s3ServiceImpl{repo: repo}
}

func (s *s3ServiceImpl) GeneratePresignedUploadURL(ctx context.Context, dto dto.PresignedUploadUrlRequest) (string, error) {
	return s.repo.GeneratePresignedUploadURL(ctx, dto.BucketName, dto.ObjectKey)
}
