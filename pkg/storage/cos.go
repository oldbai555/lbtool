package storage

import (
	"fmt"
	"github.com/oldbai555/lb/pkg/exception"
	"github.com/oldbai555/lb/utils"
	"github.com/tencentyun/cos-go-sdk-v5"
	"golang.org/x/net/context"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

// https://cloud.tencent.com/document/product/436/65644 腾讯云存储

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

type COSStorage struct {
	Client *cos.Client
	Config Config
}

func NewCOS(conf Config) (storage COSStorage, err error) {
	u, err := url.Parse(conf.BucketURL)
	if err != nil {
		err = exception.NewErr(exception.ErrStorageOptErr, fmt.Sprintf("invalid BucketURL,err is %v", err))
		return
	}

	b := &cos.BaseURL{BucketURL: u}
	storage.Client = cos.NewClient(b, &http.Client{
		// 设置超时时间
		Timeout: 100 * time.Second,
		Transport: &cos.AuthorizationTransport{
			// 如实填写账号和密钥，也可以设置为环境变量
			SecretID:  conf.SecretID,
			SecretKey: conf.SecretKey,
		},
	})

	storage.Config = conf

	return
}

// SignURL 预签名 URL
// https://cloud.tencent.com/document/product/436/35059
func (o COSStorage) SignURL(objectKey string, method utils.HTTPMethod, expiredInSec int64) (signedURL string, err error) {
	contentType, err := GetContentType(objectKey)
	if err != nil {
		err = exception.NewErr(exception.ErrStorageOptErr, fmt.Sprintf("GetContentType failed,err is %v", err))
		return
	}

	opt := &cos.PresignedURLOptions{
		Header: &http.Header{},
	}
	opt.Header.Set("Content-Type", contentType)

	u, err := o.Client.Object.GetPresignedURL(
		context.Background(),
		string(method),
		objectKey,
		o.Config.SecretID,
		o.Config.SecretKey,
		time.Duration(expiredInSec)*time.Second,
		nil,
	)
	if err != nil {
		err = exception.NewErr(exception.ErrStorageOptErr, fmt.Sprintf("GetPresignedURL failed,err is %v", err))
		return
	}

	if o.Config.CdnURL != "" {
		cdnURL, cdnErr := url.Parse(o.Config.CdnURL)
		if cdnErr != nil {
			cdnErr = exception.NewErr(exception.ErrStorageOptErr, fmt.Sprintf("url.ParseLink failed,err is %v", cdnErr))
			return signedURL, cdnErr
		}

		u.Host = cdnURL.Host
		u.Scheme = cdnURL.Scheme
	}

	signedURL = u.String()

	return
}

func (o COSStorage) Get(objectKey string) (content io.ReadCloser, err error) {
	resp, err := o.Client.Object.Get(context.Background(), objectKey, nil)
	if err != nil {
		err = exception.NewErr(exception.ErrStorageOptErr, fmt.Sprintf("GetObject failed,err is %v", err))
		return
	}

	return resp.Body, nil
}

func (o COSStorage) Put(objectKey string, reader io.Reader) (err error) {
	contentType, err := GetContentType(objectKey)
	if err != nil {
		err = exception.NewErr(exception.ErrStorageOptErr, fmt.Sprintf("GetContentType failed,err is %v", err))
		return
	}

	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: contentType,
		},
		ACLHeaderOptions: &cos.ACLHeaderOptions{
			// 如果不是必要操作，建议上传文件时不要给单个文件设置权限，避免达到限制。若不设置默认继承桶的权限。
			XCosACL: "private",
		},
	}
	_, err = o.Client.Object.Put(context.Background(), objectKey, reader, opt)
	if err != nil {
		err = exception.NewErr(exception.ErrStorageOptErr, fmt.Sprintf("PutObject failed,%v", err))
		return
	}

	return
}

func (o COSStorage) IsExist(objectKey string) (ok bool, err error) {
	_, err = o.Client.Object.Head(context.Background(), objectKey, nil)
	if err != nil {
		err = exception.NewErr(exception.ErrStorageOptErr, fmt.Sprintf("Head failed,err is %v", err))
		return
	}
	return
}

func (o COSStorage) PutFromFile(objectKey string, filePath string) (err error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == "" {
		err = exception.NewErr(exception.ErrStorageOptErr, fmt.Sprintf("file ext is required,err is %v", err))
		return
	}

	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		err = exception.NewErr(exception.ErrStorageOptErr, fmt.Sprintf("invalid file ext,err is %v", err))
		return
	}

	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: contentType,
		},
		ACLHeaderOptions: &cos.ACLHeaderOptions{
			XCosACL: "private",
		},
	}

	_, err = o.Client.Object.PutFromFile(context.Background(), objectKey, filePath, opt)
	if err != nil {
		err = exception.NewErr(exception.ErrStorageOptErr, fmt.Sprintf("PutFromFile failed ,err is %v", err))
		return
	}

	return
}

func (o COSStorage) Delete(objectKeys ...string) (deletedObjects []string, err error) {
	objects := make([]cos.Object, 0)
	for _, key := range objectKeys {
		objects = append(objects, cos.Object{
			Key: key,
		})
	}
	opt := &cos.ObjectDeleteMultiOptions{
		Objects: objects,
	}

	result, _, err := o.Client.Object.DeleteMulti(context.Background(), opt)
	if err != nil {
		err = exception.NewErr(exception.ErrStorageOptErr, fmt.Sprintf("DeleteMulti failed,err is %v", err))
		return
	}

	for _, object := range result.DeletedObjects {
		deletedObjects = append(deletedObjects, object.Key)
	}

	return
}
