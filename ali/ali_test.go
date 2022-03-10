package ali

import (
	"log"
	"os"
	"testing"

	"github.com/xueqing/s3/common"

	"gopkg.in/yaml.v2"
)

func parseS3Info() (info common.S3Info, err error) {
	// read s3 config file
	f, _ := os.OpenFile("s3.yaml", os.O_RDONLY, os.ModePerm)
	defer f.Close()
	body := make([]byte, 2048)
	n, _ := f.Read(body)
	body = body[:n]

	if err = yaml.Unmarshal(body, &info); err != nil {
		log.Fatalf("parseS3Info: parse config file error(%v)", err)
	}

	return
}

func TestUpload(t *testing.T) {
	info, err := parseS3Info()
	if err != nil {
		t.Errorf("TestUpload: parse s3 info error(%v)", err)
		return
	}

	// no prefix slash, otherwise you will see error like "oss: service returned error: StatusCode=400,
	// ErrorCode=InvalidObjectName, ErrorMessage="The specified object is not valid.", RequestId=xxx"
	key := "putobject/00.m3u8"

	// read file to be uploaded
	f, _ := os.OpenFile("../resource/media.m3u8", os.O_RDONLY, os.ModePerm)
	defer f.Close()
	body := make([]byte, 2048)
	n, _ := f.Read(body)
	body = body[:n]

	// upload with right s3 info and key
	if err := Upload(body, &info, key); err != nil {
		t.Errorf("TestUpload: upload error(%v)", err)
		return
	}
}

func TestDownload(t *testing.T) {
	TestUpload(t)
	info, err := parseS3Info()
	if err != nil {
		t.Errorf("TestDownload: parse s3 info error(%v)", err)
		return
	}

	key := "putobject/00.m3u8"

	// save file downloaded from s3
	f, err := os.Create("download.m3u8")
	if err != nil {
		t.Errorf("TestDownload: create file error(%v)", err)
		return
	}
	defer f.Close()

	// upload with right s3 info and key
	r, err := Download(&info, key)
	if err != nil {
		t.Errorf("TestDownload: download error(%v)", err)
		return
	}
	body := make([]byte, 2048)
	for {
		n, err := r.Read(body)
		if n != 0 {
			t.Logf("TestDownload: read body len(%v)", n)
			f.Write(body[:n])
		}
		if err != nil {
			t.Logf("TestDownload: read body error(%v)", err)
			break
		}
	}
}

func TestDeleteObject(t *testing.T) {
	info, err := parseS3Info()
	if err != nil {
		t.Errorf("TestDeleteObject: parse s3 info error(%v)", err)
		return
	}

	key := "deleteobject/media.m3u8"

	// read file to be uploaded
	f, _ := os.OpenFile("media.m3u8", os.O_RDONLY, os.ModePerm)
	defer f.Close()
	body := make([]byte, 2048)
	n, _ := f.Read(body)
	body = body[:n]

	// upload with right s3 info and key
	if err := Upload(body, &info, key); err != nil {
		t.Errorf("TestDeleteObject: upload error(%v)", err)
		return
	}

	// delete with right s3 info and key
	if err := DeleteObject(&info, key); err != nil {
		t.Errorf("TestDeleteObject: delete error(%v)", err)
		return
	}
}

func TestDelete(t *testing.T) {
	info, err := parseS3Info()
	if err != nil {
		t.Errorf("TestDelete: parse s3 info error(%v)", err)
		return
	}

	keys := []string{
		"deleteobjects/00.m3u8",
		"deleteobjects/01.m3u8",
	}

	// read file to be uploaded
	f, _ := os.OpenFile("../resource/media.m3u8", os.O_RDONLY, os.ModePerm)
	defer f.Close()
	body := make([]byte, 2048)
	n, _ := f.Read(body)
	body = body[:n]

	// upload with right s3 info and key
	for _, key := range keys {
		if err := Upload(body, &info, key); err != nil {
			t.Errorf("TestDelete: upload error(%v)", err)
			return
		}
	}

	// delete with right s3 info and key
	if _, err := Delete(&info, keys); err != nil {
		t.Errorf("TestDelete: delete error(%v)", err)
		return
	}
}
