package upload

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"mime/multipart"
)

type OBS interface {
	UploadFile(file *multipart.FileHeader, input *obs.PutObjectInput) (string, string, error)
	DeleteFile(keys ...string) (*obs.DeleteObjectsOutput, error)
}
