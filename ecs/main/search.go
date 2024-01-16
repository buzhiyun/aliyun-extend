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

func cnStatus(status string) string {
	//Pending：创建中。
	//Running：运行中。
	//Starting：启动中。
	//Stopping：停止中。
	//Stopped：已停止。
	switch status {
	case "Pending":
		return colorfulString("创建中", YellowFont, BlueBG)
	case "Running":
		return colorfulString("运行中", GreenFont, BlueBG)
	case "Starting":
		return colorfulString("启动中", WhiteFont, BlueBG)
	case "Stopping":
		return colorfulString("停止中", RedFont, BlueBG)
	case "Stopped":
		return colorfulString("已停止", RedFont, BlueBG)
	default:
		return status
	}

}

// 前景 背景 颜色
// ---------------------------------------
// 30  40  黑色
// 31  41  红色
// 32  42  绿色
// 33  43  黄色
// 34  44  蓝色
// 35  45  紫红色
// 36  46  青蓝色
// 37  47  白色
//
// 代码 意义
// -------------------------
//  0  终端默认设置
//  1  高亮显示
//  4  使用下划线
//  5  闪烁
//  7  反白显示
//  8  不可见

type ColorFont int
type ColorBG int

const (
	// 字体本身
	BlackFont  ColorFont = 30 + iota // 黑色
	RedFont                          // 红色
	GreenFont                        // 绿色
	YellowFont                       // 黄色
	BlueFont                         // 蓝色
	PurpleFont                       // 紫色
	CyanFont                         // 青色
	WhiteFont                        // 白色

	// 背景色
	BlackBG  ColorBG = 40 + iota
	RedBG            // 红色
	GreenBG          // 绿色
	YellowBG         // 黄色
	BlueBG           // 蓝色
	PurpleBG         // 紫色
	CyanBG           // 青色
	WhiteBG          // 白色
)

func colorfulString(str string, fontColor ColorFont, backgroudColor ColorBG) string {
	return fmt.Sprintf("\033[1;%v;%vm%s\033[0m", fontColor, backgroudColor, str)
}

func colorfulPrint(str string, fontColor ColorFont, backgroudColor ColorBG) {
	fmt.Printf("\033[1;%v;%vm%s\033[0m", fontColor, backgroudColor, str)
}

func inSlice(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

// 显示结果
func printEcsInfo(searchResult []aliyunecs.Instance, args ...string) {
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
		buffer.WriteString("\t状态：")
		buffer.WriteString(cnStatus(ecs.Status))

		if inSlice(args, "--type") {
			buffer.WriteString("\t规格：")
			buffer.WriteString(ecs.InstanceType)
			buffer.WriteString(fmt.Sprintf(" %vC%vG", ecs.Cpu, ecs.Memory/1024))
		}

		if inSlice(args, "--os") {
			buffer.WriteString("\t操作系统：")
			buffer.WriteString(ecs.OSName)
		}

		if inSlice(args, "--expire") {
			buffer.WriteString("\t过期时间：")
			buffer.WriteString(ecs.ExpiredTime)
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

// 返回内网地址
func getInternalIp(ecs aliyunecs.Instance) string {
	if len(ecs.VpcAttributes.PrivateIpAddress.IpAddress) > 0 {
		return ecs.VpcAttributes.PrivateIpAddress.IpAddress[0]
	} else {
		return ecs.InnerIpAddress.IpAddress[0]
	}
}

// 搜索IP
func SearchIP(ipStr string, args ...string) {
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
	printEcsInfo(searchResult, args...)

}

// 按名字搜索
func SearchName(ecsName string, args ...string) {
	instances, err := ecs.GetInstances()
	if err != nil {
		fmt.Println(err.Error())
	}

	var searchResult []aliyunecs.Instance

	//对每个实例搜索
	for _, instance := range instances {
		if strings.Contains(instance.InstanceName, ecsName) {
			searchResult = append(searchResult, instance)
		}
	}
	printEcsInfo(searchResult, args...)

}

func main() {
	//检查参数是否足够
	if len(os.Args) < 3 {
		fmt.Println("search [ name | ip ]  [ searchValue ]  [ --expire | --type | --os ]")
		fmt.Println("eg. search name sdcf_v3")
		os.Exit(0)
	}

	//初始化配置
	config.GetConf()

	switch strings.ToLower(os.Args[1]) {
	case "ip":
		if len(os.Args) > 3 {
			SearchIP(os.Args[2], os.Args[3:]...)
		} else {
			SearchIP(os.Args[2])
		}

	case "name":
		if len(os.Args) > 3 {
			SearchName(os.Args[2], os.Args[3:]...)
		} else {
			SearchName(os.Args[2])
		}
	}

}
