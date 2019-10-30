package main

import (
	"bytes"
	"fmt"
	aliyunecs "github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/buzhiyun/aliyun-extend/config"
	"github.com/buzhiyun/aliyun-extend/ecs"
	jsoniter "github.com/json-iterator/go"
	"os"
	"strings"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

//显示结果
func printEcsInfo(searchResult []aliyunecs.Instance) {
	buffer := bytes.Buffer{}
	for _, ecs := range searchResult {
		buffer.WriteString("id: ")
		buffer.WriteString(ecs.InstanceId)
		buffer.WriteString("\t主机名：")
		buffer.WriteString(ecs.InstanceName)
		buffer.WriteString("\t")
		if len(ecs.VpcAttributes.PrivateIpAddress.IpAddress) > 0 {
			buffer.WriteString("VPC IP: ")
			buffer.WriteString(ecs.VpcAttributes.PrivateIpAddress.IpAddress[0])
			buffer.WriteString("\t")
		}
		if len(ecs.InnerIpAddress.IpAddress) > 0 {
			buffer.WriteString("私网IP: ")
			buffer.WriteString(ecs.InnerIpAddress.IpAddress[0])
			buffer.WriteString("\t")
		}
		if len(ecs.PublicIpAddress.IpAddress) > 0 {
			buffer.WriteString("公网IP: ")
			buffer.WriteString(ecs.PublicIpAddress.IpAddress[0])
		}
		buffer.WriteString("\n")
	}
	fmt.Printf(buffer.String())

}

func search() {
	instances, _ := ecs.GetInstances()

	//jsonData , _ :=json.MarshalIndent(&instances,"","  ")
	//fmt.Println(string(jsonData))

	for i, x := range instances {
		jsonData, _ := json.MarshalIndent(&x, "", "  ")
		//jsonData , _ :=json.Marshal(&x)
		fmt.Printf("第 %d 个实例是 : %d \n", i, string(jsonData))
	}

}

//返回内网地址
func getInternalIp(ecs aliyunecs.Instance) string {
	if len(ecs.VpcAttributes.PrivateIpAddress.IpAddress) > 0 {
		return ecs.VpcAttributes.PrivateIpAddress.IpAddress[0]
	} else {
		return ecs.InnerIpAddress.IpAddress[0]
	}
}

//搜索IP
func SearchIP(ipStr string) {
	instances, _ := ecs.GetInstances()
	var searchResult []aliyunecs.Instance
	hitSearch := false

	for _, x := range instances {
		//jsonData , _ :=json.MarshalIndent(&x,"","  ")
		//fmt.Printf("第 %d 个实例是 : %d \n" ,i ,string(jsonData))
		hitSearch = false
		//查找 ip
		for _, ip := range x.VpcAttributes.PrivateIpAddress.IpAddress {
			if strings.Contains(ip, ipStr) {
				hitSearch = true
			}
		}
		for _, ip := range x.PublicIpAddress.IpAddress {
			if strings.Contains(ip, ipStr) {
				hitSearch = true
			}
		}
		for _, ip := range x.InnerIpAddress.IpAddress {
			if strings.Contains(ip, ipStr) {
				hitSearch = true
			}
		}

		//如果有命中就追加
		if hitSearch {
			searchResult = append(searchResult, x)
		}
	}
	printEcsInfo(searchResult)

}

//按名字搜索
func SearchName(ecsName string) {
	instances := ecs.GetInstances()
	var searchResult []aliyunecs.Instance

	//对每个实例搜索
	for _, instance := range instances {
		if strings.Contains(instance.InstanceName, ecsName) {
			searchResult = append(searchResult, instance)
		}
	}
	printEcsInfo(searchResult)

}

func main() {
	//检查参数是否足够
	if len(os.Args) < 3 {
		fmt.Println("search [ name | ip ]  [ searchValue ]")
		os.Exit(0)
	}

	//初始化配置
	config.GetConf()

	switch strings.ToLower(os.Args[1]) {
	case "ip":
		SearchIP(os.Args[2])
	case "name":
		SearchName(os.Args[2])
	}

}
