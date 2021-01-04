package ecs

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/buzhiyun/aliyun-extend/config"
	"log"
)

//得到所有ECS 主机
func GetInstances() (res []ecs.Instance, err error) {
	//fmt.Println("access Info :",config.Conf.Global.Regionid, config.Conf.Global.Key)
	var instances []ecs.Instance

	pageNum := 1 // 先查第一页的
	maxPage := 1 //默认最大页数就是1
	for pageNum <= maxPage {
		client, err := ecs.NewClientWithAccessKey(config.Conf.Global.Regionid, config.Conf.Global.Key, config.Conf.Global.Secret)
		if err != nil {
			log.Print(err.Error())
			return instances, err
		}

		request := ecs.CreateDescribeInstancesRequest()
		request.PageSize = requests.NewInteger(100)
		request.PageNumber = requests.NewInteger(pageNum)
		response, err := client.DescribeInstances(request)

		if err != nil {
			log.Print(err.Error())
			return instances, err
		}

		if response != nil {
			maxPage = ((response.TotalCount - 1) / 100) + 1

			instances = append(instances, response.Instances.Instance...)
		}
		//增加页码，准备取下一页
		pageNum++
	}

	return instances, nil
}

//返回ecs的所有IP信息
func GetInstanceIp(instance ecs.Instance) (ipAddresses []string) {

	ipAddresses = append(ipAddresses, instance.InnerIpAddress.IpAddress...)
	ipAddresses = append(ipAddresses, instance.PublicIpAddress.IpAddress...)
	ipAddresses = append(ipAddresses, instance.VpcAttributes.PrivateIpAddress.IpAddress...)

	return ipAddresses
}
