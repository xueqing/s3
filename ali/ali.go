package ali

import (
	"bytes"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/logger"
	"github.com/xueqing/s3/common"
)

// InitS3 ...
func InitS3(info *common.S3Info) (cli *oss.Client, err error) {
	cli, err = oss.New(info.Domain, info.ID, info.Secret)
	if err != nil {
		logger.Warningf("InitS3: New with info(%v) error(%v)", info, err)
		return
	}
	// logger.Infof("InitS3: New with info(%v) success", info)
	return
}

// Upload ...
func Upload(data []byte, info *common.S3Info, key string) (err error) {
	cli, err := InitS3(info)
	if err != nil {
		return
	}
	bucket, err := cli.Bucket(info.Bucket)
	if err != nil {
		logger.Warningf("Upload: get Bucket from s3(%v) error(%v)", info, err)
		return
	}

	err = bucket.PutObject(key, bytes.NewReader(data))
	if err != nil {
		logger.Warningf("Upload: PutObject s3(%v) key(%v) error(%v)", info, key, err)
		return
	}
	logger.Infof("Upload: upload data to s3 bucket(%v) key(%s) success", info.Bucket, key)
	return
}

// DeleteObject Delete one object
func DeleteObject(info *common.S3Info, key string) (err error) {
	cli, err := InitS3(info)
	if err != nil {
		return
	}
	bucket, err := cli.Bucket(info.Bucket)
	if err != nil {
		logger.Warningf("DeleteObject: get Bucket from s3(%v) error(%v)", info, err)
		return
	}

	err = bucket.DeleteObject(key)
	if err != nil {
		logger.Warningf("DeleteObject: DeleteObject s3(%v) key(%v) error(%v)", info, key, err)
		return
	}

	logger.Infof("DeleteObject: delete bucket(%v) key(%v) from s3 success", info.Bucket, key)
	return
}

// Delete Batch delete files
func Delete(info *common.S3Info, keys []string) (deleted []string, err error) {
	cli, err := InitS3(info)
	if err != nil {
		return
	}
	bucket, err := cli.Bucket(info.Bucket)
	if err != nil {
		logger.Warningf("Delete: get Bucket from s3(%v) error(%v)", info, err)
		return
	}

	output, err := bucket.DeleteObjects(keys)
	// save deleted keys
	deleted = make([]string, len(output.DeletedObjects))
	for idx, item := range output.DeletedObjects {
		deleted[idx] = item
	}

	if err != nil {
		logger.Warningf("Delete: DeleteObjects s3(%v) input(%v) output(%v) error(%v)",
			info, len(keys), len(deleted), err)
		return
	}

	logger.Infof("Delete: delete bucket(%v) items(%v) from s3 success", info.Bucket, len(deleted))
	return
}
