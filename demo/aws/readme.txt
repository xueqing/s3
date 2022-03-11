1. 修改配置文件 s3.yaml
   sinfo        S3 的信息
   uploadfile   本地要上传的文件绝对路径
   uploadkey    上传到 s3 的 key
   downloadfile 保存到本地的文件绝对路径
   downloadkey  和 uploadkey 保持一样即可，表示下载刚上传成功的文件
2. 执行 ./aws-demo 测试使用 aws-sdk-go 上传和下载文件