package huawei_cloud

import (
	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"mime/multipart"
	"os"
	"path"
	"time"
)

type Config struct {
	//前两个字段推荐从 系统环境变量中获取
	AccessKeyID      string //访问 OBS 所需的密钥 ID
	SecretAccessKey  string //访问 OBS 所需的密钥密钥
	Location         string //存储桶所在区域，必须和传入 Endpoint 中 Region 保持一致
	BucketName       string //存储桶名称
	BucketUrl        string //存储桶 URL."https://your-bucket-name.obs.cn-north-4.myhuaweicloud.com"
	Endpoint         string //OBS 服务的 Endpoint，用与访问 OBS 的 API "https://obs.cn-north-4.myhuaweicloud.com"
	BasePath         string //上传文件时，文件在存储桶中的基础路径
	AvatarType       string
	AccountAvatarUrl string
	GroupAvatarUrl   string
}

var (
	NotAvatar         = "NotAvatar"
	AccountAvatarType = "AccountAvatarType"
	GroupAvatarType   = "GroupAvatarType"
)

type OBS struct {
	config Config
}

func Init(config Config) *OBS {
	return &OBS{config: config}
}

var ErrFileOpen = errors.New("文件打开失败")

// UploadFile 流式上传对象。返回访问地址，文件 key，error
// 传入要上传的文件，以及上传文件时的参数（比如文件的存储类型，文件的访问权限等）
func (o *OBS) UploadFile(file *multipart.FileHeader, input *obs.PutObjectInput) (string, string, error) {
	// 获取 obsClient 实例
	obsClient, err := o.createObsClient()
	if err != nil {
		return "", "", err
	}

	//指定存储桶名称
	input.Bucket = o.config.BucketName

	//对象名。对象名是对象在存储桶中的唯一标识。对象名是对象在桶中的完整路径，路径中不包含桶名。
	//例如，您对象的访问地址为examplebucket.obs.cn-north-4.myhuaweicloud.com/folder/test.txt 中，对象名为folder/test.txt。
	key := o.config.BasePath + time.Now().Format("2006-01-02-15:04:05.99") + path.Ext(file.Filename)
	if o.config.AvatarType == AccountAvatarType {
		key = o.config.AccountAvatarUrl + time.Now().Format("2006-01-02-15:04:05.99") + path.Ext(file.Filename)
	} else if o.config.AvatarType == GroupAvatarType {
		key = o.config.GroupAvatarUrl + time.Now().Format("2006-01-02-15:04:05.99") + path.Ext(file.Filename)
	}
	//指定上传对象
	input.Key = key

	//读取本地文件
	f, openError := file.Open()
	if openError != nil {
		return "", "", ErrFileOpen
	}
	defer f.Close() //创建文件 defer 关闭
	input.Body = f

	//流式上传本地文件
	_, err = obsClient.PutObject(input)
	if err != nil {
		return "", "", errors.New("function obsClient.PutObject, err:" + err.Error())
	}

	return o.config.BucketUrl + "/" + key, key, nil
}

// DeleteFile 通过 key 删除对应文件
func (o *OBS) DeleteFile(keys ...string) (*obs.DeleteObjectsOutput, error) {
	// 获取 obsClient 实例
	obsClient, err := o.createObsClient()
	if err != nil {
		return &obs.DeleteObjectsOutput{}, err
	}

	input := &obs.DeleteObjectsInput{}
	//指定存储桶名称
	input.Bucket = o.config.BucketName
	//指定删除对象列表
	var objectsToDelete []obs.ObjectToDelete
	for _, key := range keys {
		objectsToDelete = append(objectsToDelete, obs.ObjectToDelete{Key: key})
	}
	input.Objects = objectsToDelete

	//批量删除对象
	output, err := obsClient.DeleteObjects(input)
	if err != nil {
		return &obs.DeleteObjectsOutput{}, err
	}
	return output, nil
}

// createObsClient 创建 obsClient 实例
func (o *OBS) createObsClient() (*obs.ObsClient, error) {
	//推荐通过环境变量获取 AK SK，这里也可以使用其他外部引入方式传入，如果使用硬编码可能会存在泄露风险。
	o.config.AccessKeyID = os.Getenv("AccessKeyID")
	o.config.SecretAccessKey = os.Getenv("SecretAccessKey")
	// 创建 OBSClient 实例
	obsClient, err := obs.New(o.config.AccessKeyID, o.config.SecretAccessKey, o.config.Endpoint)
	if err != nil {
		return nil, err
	}

	return obsClient, nil
}
