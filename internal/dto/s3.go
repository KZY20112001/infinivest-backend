package dto

type PresignedUploadUrlRequest struct {
	BucketName string `json:"bucket_name"`
	ObjectKey  string `json:"object_key"`
}
