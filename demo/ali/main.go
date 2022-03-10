package main

import (
	"os"

	"github.com/google/logger"
	"gopkg.in/yaml.v2"

	"github.com/xueqing/s3/ali"
	"github.com/xueqing/s3/common"
)

// Info Saves sth about s3
type Info struct {
	SInfo common.S3Info `yaml:"sinfo"`

	UploadFile string `yaml:"uploadfile"`
	UploadKey  string `yaml:"uploadkey"`

	DownloadFile string `yaml:"downloadfile"`
	DownloadKey  string `yaml:"downloadkey"`
}

func parseS3Info() (info Info, err error) {
	// read s3 config file
	f, _ := os.OpenFile("s3.yaml", os.O_RDONLY, os.ModePerm)
	defer f.Close()
	body := make([]byte, 2048)
	n, _ := f.Read(body)
	body = body[:n]

	if err = yaml.Unmarshal(body, &info); err != nil {
		logger.Fatalf("parseS3Info: parse config file error(%v)", err)
	}

	return
}

func testUpload(info common.S3Info, file, key string) {
	logger.Infof("testUpload: info(%v) file(%v) key(%v)", info, file, key)

	// read file to be uploaded
	f, err := os.OpenFile(file, os.O_RDONLY, os.ModePerm)
	if err != nil {
		logger.Warningf("testUpload: open file(%v) error(%v)", file, err)
		return
	}
	defer f.Close()
	body := make([]byte, 2048)
	n, _ := f.Read(body)
	body = body[:n]

	// upload with right s3 info and key
	if err := ali.Upload(body, &info, key); err != nil {
		logger.Warningf("testUpload: upload error(%v)", err)
		return
	}
}

func testDownload(info common.S3Info, file, key string) {
	logger.Infof("testDownload: info(%v) file(%v) key(%v)", info, file, key)

	// save file downloaded from s3
	f, err := os.Create(file)
	if err != nil {
		logger.Warningf("testDownload: create file(%v) error(%v)", file, err)
		return
	}
	defer f.Close()

	// upload with right s3 info and key
	r, err := ali.Download(&info, key)
	if err != nil {
		logger.Warningf("testDownload: download error(%v)", err)
		return
	}
	body := make([]byte, 2048)
	for {
		n, err := r.Read(body)
		if n != 0 {
			logger.Infof("testDownload: read body len(%v)", n)
			f.Write(body[:n])
		}
		if err != nil {
			logger.Warningf("testDownload: read body error(%v)", err)
			break
		}
	}
}

func main() {
	info, err := parseS3Info()
	if err != nil {
		logger.Errorf("main: parse s3 info error(%v)", err)
		return
	}

	testUpload(info.SInfo, info.UploadFile, info.UploadKey)

	testDownload(info.SInfo, info.DownloadFile, info.DownloadKey)
	return
}
