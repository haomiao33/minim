package ossclient

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"log"
)

type AliyunOssClient struct {
	client     *oss.Client
	BucketName string
	Host       string
}

func NewAliyunOssClient(endpoint, accessKeyId, accessKeySecret, bucketName, region, host string) *AliyunOssClient {
	// 创建OSSClient实例，并使用V4签名。
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		log.Fatalf("oss error:%v", err)
	}
	// 输出客户端信息。
	log.Printf("Client: %#v\n", client)

	return &AliyunOssClient{
		client:     client,
		BucketName: bucketName,
		Host:       host,
	}
}

func (a *AliyunOssClient) Upload(key string, pathFile string) (string, error) {
	log.Printf("Upload file:%s -> %s", pathFile, key)
	bucketName := a.BucketName
	bucket, err := a.client.Bucket(bucketName)
	if err != nil {
		log.Printf("Error init bucket:%v", err)
		return "", err
	}

	err = bucket.PutObjectFromFile(key, pathFile)
	if err != nil {
		log.Printf("Error upload 2:%v", err)
		return "", err
	}
	log.Printf("Upload file:%s -> %s", key, a.Host+"/"+key)
	return a.Host + "/" + key, nil
}

func (a *AliyunOssClient) UploadByReader(key string, reader io.Reader) (string, error) {
	log.Printf("UploadByReader file:%s ", key)
	bucketName := a.BucketName
	bucket, err := a.client.Bucket(bucketName)
	if err != nil {
		log.Printf("Error init bucket:%v", err)
		return "", err
	}

	err = bucket.PutObject(key, reader)
	if err != nil {
		log.Printf("Error upload 2:%v", err)
		return "", err
	}
	log.Printf("Upload file:%s -> %s", key, a.Host+"/"+key)
	return a.Host + "/" + key, nil
}
