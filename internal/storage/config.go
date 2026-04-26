package storage

const (
	BlobDriverLocal = "local"
	BlobDriverS3    = "s3"
)

type Config struct {
	Driver    string
	LocalPath string
	S3Bucket  string
	S3Region  string
}
