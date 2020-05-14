package oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	jsoniter "github.com/json-iterator/go"
	"log"
	"os"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func ListObject(bucketName string, prefix string, marker string, clients ...*oss.Client) oss.ListObjectsResult {

	var client = Client
	for _, cli := range clients {
		client = cli
	}

	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		log.Println("Error:", err)
		os.Exit(-1)
	}

	// 设置列举文件的最大个数，并列举文件。
	lsRes, err := bucket.ListObjects(oss.MaxKeys(200), oss.Prefix(prefix), oss.Marker(marker))
	if err != nil {
		log.Println("Error:", err)
		os.Exit(-1)
	}

	return lsRes
}

func DownloadObject(bucketName, objectName, downloadedFileName string, clients ...*oss.Client) {
	var client = Client
	for _, cli := range clients {
		client = cli
	}

	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		log.Println(err.Error())
	}
	// 下载文件。
	err = bucket.GetObjectToFile(objectName, downloadedFileName)
	if err != nil {
		log.Println(err.Error())
	}
}

func UploadObject(bucketName, objectName, localFileName string, clients ...*oss.Client) {

	var client = Client
	for _, cli := range clients {
		client = cli
	}

	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		log.Println(err.Error())
	}
	// 上传文件。
	err = bucket.PutObjectFromFile(objectName, localFileName)
	if err != nil {
		log.Println(err.Error())
	}
}
