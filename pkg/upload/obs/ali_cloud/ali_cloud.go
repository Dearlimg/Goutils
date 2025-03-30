package ali_cloud

import (
	"errors"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/joho/godotenv"
	"log"
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
	BucketDomain     string
}

var (
	NotAvatar         = "NotAvatar"
	AccountAvatarType = "AccountAvatarType"
	GroupAvatarType   = "GroupAvatarType"
)

type OSS struct {
	config Config
}

func Init(config Config) *OSS {
	return &OSS{config: config}
}

var ErrFileOpen = errors.New("文件打开失败")

func (o *OSS) UploadFile(file *multipart.FileHeader) (string, string, error) {
	// 1. 创建OSS客户端（修改方法调用）
	client, err := o.createOSSClient()
	if err != nil {
		return "", "", err
	}

	// 2. 生成对象键（保持原有逻辑）
	key := o.config.BasePath + time.Now().Format("2006-01-02-15:04:05.99") + path.Ext(file.Filename)
	if o.config.AvatarType == AccountAvatarType {
		key = o.config.AccountAvatarUrl + time.Now().Format("2006-01-02-15:04:05.99") + path.Ext(file.Filename)
	} else if o.config.AvatarType == GroupAvatarType {
		key = o.config.GroupAvatarUrl + time.Now().Format("2006-01-02-15:04:05.99") + path.Ext(file.Filename)
	}

	// 3. 获取存储桶实例
	bucket, err := client.Bucket(o.config.BucketName)
	if err != nil {
		return "", "", errors.New("get bucket failed: " + err.Error())
	}

	// 4. 读取文件并上传（修改上传方式）
	f, openError := file.Open()
	if openError != nil {
		return "", "", ErrFileOpen
	}
	defer f.Close()

	// 阿里云OSS上传接口差异
	err = bucket.PutObject(key, f)
	if err != nil {
		return "", "", errors.New("function bucket.PutObject failed: " + err.Error())
	}

	// 5. 返回访问URL（使用BucketDomain）
	return o.config.BucketUrl + "/" + key, key, nil
}

func (o *OSS) DeleteFile(keys ...string) (oss.DeleteObjectsResult, error) {
	client, err := o.createOSSClient()
	if err != nil {
		return oss.DeleteObjectsResult{}, err
	}

	bucket, err := client.Bucket(o.config.BucketName)
	if err != nil {
		return oss.DeleteObjectsResult{}, errors.New("get bucket failed: " + err.Error())
	}

	// 阿里云删除接口差异（直接传入字符串数组）
	delRes, err := bucket.DeleteObjects(keys)
	if err != nil {
		return oss.DeleteObjectsResult{}, err
	}
	return delRes, nil
}

// 修改客户端创建方法
func (o *OSS) createOSSClient() (*oss.Client, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("无法加载 .env 文件")
	}
	o.config.AccessKeyID = os.Getenv("ALIYUN_OSS_ACCESS_KEY_ID")
	o.config.SecretAccessKey = os.Getenv("ALIYUN_OSS_ACCESS_KEY_SECRET")

	// 创建客户端（注意参数顺序差异）
	client, err := oss.New(o.config.AccessKeyID, o.config.SecretAccessKey, o.config.Endpoint)
	if err != nil {
		return nil, errors.New("create OSS client failed: " + err.Error())
	}
	return client, nil
}
