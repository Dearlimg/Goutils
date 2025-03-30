package upload

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"mime/multipart"
)

type OBS interface {
	UploadFile(file *multipart.FileHeader, input *obs.PutObjectInput) (string, string, error)
	DeleteFile(keys ...string) (*obs.DeleteObjectsOutput, error)
}

// OSS 尝试oss
type OSS interface {
	UploadFile(file *multipart.FileHeader) (string, string, error)
	DeleteFile(keys ...string) (oss.DeleteObjectsResult, error)
}
