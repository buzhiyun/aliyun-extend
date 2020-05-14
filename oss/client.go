package oss

import (
	aliyunoss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/buzhiyun/aliyun-extend/config"
	"log"
)

var Client *aliyunoss.Client

func GetClient(endpoint, accessKeyId, accessKeySecret string) *aliyunoss.Client {
	// Endpoint以杭州为例，其它Region请按实际情况填写。
	//endpoint := "http://oss-cn-hangzhou.aliyuncs.com"
	if len(endpoint) == 0 {
		endpoint = config.Conf.Oss.EndPoint["default"]
	}

	// 阿里云主账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM账号进行API访问或日常运维，请登录 https://ram.console.aliyun.com 创建RAM账号。
	if len(accessKeyId) == 0 || len(accessKeySecret) == 0 {
		accessKeyId = config.Conf.Global.Key
		accessKeySecret = config.Conf.Global.Secret
	}

	// 创建OSSClient实例。
	client, err := aliyunoss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		log.Fatal(err.Error())
	}
	return client
}
