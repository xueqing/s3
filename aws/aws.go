package aws

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/logger"

	"github.com/xueqing/s3/common"
)

var (
	// DeleteRequestTimeout Set timeout for s3 delete request
	DeleteRequestTimeout = 10 * time.Second
)

// InitS3 ...
func InitS3(info *common.S3Info) (s *s3.S3, err error) {
	creds := credentials.NewStaticCredentials(info.ID, info.Secret, "")
	cfgs := &aws.Config{
		Credentials: creds,
		Endpoint:    aws.String(info.Domain),
		Region:      aws.String(info.Region),

		DisableSSL:       aws.Bool(info.DisableSSL),
		S3ForcePathStyle: aws.Bool(info.PathStyle),
		// S3DisableContentMD5Validation: aws.Bool(DisableContentMD5Validation),
	}
	sess, err := session.NewSession(cfgs)
	sess = session.Must(sess, err)

	if err != nil {
		logger.Warningf("InitS3: NewSession with info(%v) error(%v)", info, err)
		return
	}
	s = s3.New(sess)
	// logger.Infof("InitS3: NewSession with info(%v) success", info)
	return
}

// Upload ...
func Upload(data []byte, info *common.S3Info, key string) (err error) {
	s, err := InitS3(info)
	if err != nil {
		return
	}
	input := &s3.PutObjectInput{
		ACL:    aws.String("public-read"), // 对象设置为公共可读, http下载需要; 公共可读情况下生成的http下载url,一直是有效的，不受Presign() 过期时间的影响
		Body:   bytes.NewReader(data),
		Bucket: aws.String(info.Bucket),
		Key:    aws.String(key),
	}
	_, err = s.PutObject(input)
	if err != nil {
		logger.Warningf("Upload: PutObject s3(%v) key(%v) error(%v)", info, key, err)
		return
	}
	logger.Infof("Upload: upload data to s3 bucket(%v) key(%s) success", info.Bucket, key)
	return
}

// Download ...
func Download(info *common.S3Info, key string) (r io.ReadCloser, err error) {
	s, err := InitS3(info)
	if err != nil {
		return
	}
	input := &s3.GetObjectInput{
		Bucket: aws.String(info.Bucket),
		Key:    aws.String(key),
	}
	out, err := s.GetObject(input)
	if err != nil {
		logger.Warningf("Download: GetObject s3(%v) key(%v) error(%v)", info, key, err)
		return
	}
	logger.Infof("Download: read data from s3 bucket(%v) key(%s) length(%v) success", info.Bucket, key, *out.ContentLength)
	r = out.Body
	return
}

// DeleteObject Delete one object
func DeleteObject(info *common.S3Info, key string) (err error) {
	s, err := InitS3(info)
	if err != nil {
		return
	}

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(info.Bucket),
		Key:    aws.String(key),
	}

	var output *s3.DeleteObjectOutput
	output, err = s.DeleteObject(input)
	if err != nil {
		logger.Warningf("DeleteObject: DeleteObject s3(%v) input(%v) output(%v) error(%v)",
			info, input, output, err)
		return
	}

	logger.Infof("DeleteObject: delete bucket(%v) items(%v) from s3 success", info.Bucket, output)
	return
}

// Delete Batch delete files
func Delete(info *common.S3Info, keys []string) (deleted []string, err error) {
	s, err := InitS3(info)
	if err != nil {
		return
	}

	objects := make([]*s3.ObjectIdentifier, len(keys))
	for idx, key := range keys {
		objects[idx] = &s3.ObjectIdentifier{
			Key: aws.String(key),
		}
	}
	input := &s3.DeleteObjectsInput{
		Bucket: aws.String(info.Bucket),
		Delete: &s3.Delete{
			Objects: objects,
			// Quiet:   aws.Bool(false),
		},
	}

	// DeleteObjects won't exit when network failure, so use DeleteObjectsWithContext
	ctx, cancel := context.WithTimeout(context.Background(), DeleteRequestTimeout)
	defer cancel()
	output, err := s.DeleteObjectsWithContext(ctx, input)
	// save deleted keys
	deleted = make([]string, len(output.Deleted))
	for idx, item := range output.Deleted {
		deleted[idx] = *item.Key
	}

	if err != nil {
		logger.Warningf("Delete: DeleteObjects s3(%v) input(%v) output(%v) error(%v)",
			info, len(input.Delete.Objects), len(deleted), err)
		return
	}

	logger.Infof("Delete: delete bucket(%v) items(%v) from s3 success", info.Bucket, len(deleted))
	return
}
