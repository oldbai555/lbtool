package storage

import (
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
	"io"
	"os"
	"time"
)

var FileStorage FileStorageInterface

// Config 存储配制
type Config struct {
	// Type 存储类型, 可配置 aliyun, qcloud；分别对应阿里云OSS, 腾讯云COS
	Type string `validate:"required,oneof=aliyun qcloud"`
	// CdnURL CDN绑定域名，可选配置，本地存储必填
	CdnURL string `validate:"omitempty,url"`

	// 阿里云OSS相关配置，请使用子账户凭据，且仅授权oss访问权限
	AccessKeyId     string `validate:"required_if=Type aliyun"`
	AccessKeySecret string `validate:"required_if=Type aliyun"`
	EndPoint        string `validate:"required_if=Type aliyun"`
	Bucket          string `validate:"required_if=Type aliyun"`

	// 腾讯云OSS相关配置，请使用子账户凭据，且仅授权cos访问权限
	SecretID  string `validate:"required_if=Type qcloud"`
	SecretKey string `validate:"required_if=Type qcloud"`
	// 格式 "https://bucket.cos.region.myqcloud.com" bucket 替换自己的桶，region 替换自己的区域
	BucketURL string `validate:"required_if=Type qcloud"`

	// 本地存储相关配置
	// LocalRootPath 本地存储文件的根目录，必须是绝对路径
	LocalRootPath string `validate:"required_if=Type local"`
	// ServerRootPath 文件服务的根目录，http服务中的文件根目录，相对路径，用于识别文件服务请求的路径标识
	ServerRootPath string `validate:"required_if=Type local"`
}

type Credentials struct {
	SecretID     string `json:"secret_id"`
	SecretKey    string `json:"secret_key"`
	SessionToken string `json:"session_token"`
}

type FileStorageInterface interface {
	SignURL(objectKey string, method utils.HTTPMethod, expiredInSec int64) (signedURL string, err error)
	Get(objectKey string) (content io.ReadCloser, err error)
	Put(objectKey string, reader io.Reader) (err error)
	IsExist(objectKey string) (ok bool, err error)
	PutFromFile(objectKey string, filePath string) (err error)
	Delete(objectKeys ...string) (deletedObjects []string, err error)
	GetCredentials() (*Credentials, error)
	GetSignature(httpMethod, name, ak, sk string, expired time.Duration) string
}

func Setup(conf Config) {
	var err error

	if conf.Type == string(utils.AliyunStorage) {
		FileStorage, err = NewOSS(conf)
		if err != nil {
			log.Errorf("NewOSS failed,err is %v", err)
			os.Exit(1)
			return
		}
	}

	if conf.Type == string(utils.QcloudStorage) {
		FileStorage, err = NewCOS(conf)
		if err != nil {
			log.Errorf("NewCOS failed,err is %v", err)
			os.Exit(1)
			return
		}
	}

	return
}
