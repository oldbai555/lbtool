package storage

import (
	"github.com/oldbai555/lb/log"
	"github.com/oldbai555/lb/utils"
	"io"
	"os"
)

var FileStorage FileStorageInterface

type FileStorageInterface interface {
	SignURL(objectKey string, method utils.HTTPMethod, expiredInSec int64) (signedURL string, err error)
	Get(objectKey string) (content io.ReadCloser, err error)
	Put(objectKey string, reader io.Reader) (err error)
	IsExist(objectKey string) (ok bool, err error)
	PutFromFile(objectKey string, filePath string) (err error)
	Delete(objectKeys ...string) (deletedObjects []string, err error)
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
