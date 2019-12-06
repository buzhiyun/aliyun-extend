package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var Conf config

func init() {
	//初始化配置
	GetConf()
}

type config struct {
	// golbal 配置阿里云的key 和secret
	Global struct {
		Key      string `yaml:key`
		Secret   string `yaml:secret`
		Regionid string `yaml:"regionId"`
	} `yaml:global`

	//oss 配置部分
	Oss struct {
		EndPoint map[string]string `yaml:endpoint`
	} `yaml:oss`
}

//获取当前目录
func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

//加载配置
func GetConf() {
	//filePath := "/aliyun3/lib/.aliyun.yaml"
	filePath := "/Volumes/data/DEV/go/jobs/src/aliyun/config/.aliyun.yaml"

	//把yaml形式的字符串解析成struct类型
	content, _ := ioutil.ReadFile(filePath)

	err := yaml.Unmarshal(content, &Conf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	//转换成yaml字符串类型
	//d, err := yaml.Marshal(&Conf)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Printf("", string(d))
}
