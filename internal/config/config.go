package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

var YamlConfig SdkConfig

type SdkConfig struct {
	App   APP    `yaml:"app"`
	Mysql MYSQL  `yaml:"mysql"`
	Conf  Config `yaml:"conf"`
}

type Config struct {
	Blood     int32             `yaml:"blood"`
	ReadyTime int32             `yaml:"readyTime"`
	HList     map[string]string `yaml:"hList"`
	LList     map[string]string `yaml:"lList"`
}

type APP struct {
	Version        string `yaml:"version"`
	Host           string `yaml:"host"`
	Name           string `yaml:"name"`
	Port           int    `yaml:"port"`
	GinPort        int    `yaml:"ginPort"`
	MaxConn        int    `yaml:"maxConn"`
	WorkerPoolSize uint32 `yaml:"workerPoolSize"`
	LogFile        string `yaml:"logFile"`
	MaxPacketSize  uint32 `yaml:"maxPacketSize"`
}

type MYSQL struct {
	Dns string `yaml:"dns"`
}

func Read() bool {

	dataBytes, err := os.ReadFile("./conf/app.yaml")
	if err != nil {
		fmt.Println("读取文件失败：", err)
		return false
	}
	//YamlConfig := SdkConfig{}
	fmt.Println("readYamData = ", len(dataBytes))
	err = yaml.Unmarshal(dataBytes, &YamlConfig)
	if err != nil {
		fmt.Println("解析 yaml 文件失败：", err)
		return false
	} else {
		//fmt.Println("yamlConfig:ApiKey=", YamlConfig.Conf.Blood)
		//fmt.Println("yamlConfig:ApiSecret=", YamlConfig.Livekit.ApiCert)
		//fmt.Println("yamlConfig.Host=", YamlConfig.Livekit.Host)
		fmt.Println("Conf,Host=", YamlConfig.App.Host)
		fmt.Println("Conf,红方=", YamlConfig.Conf.HList)
		fmt.Println("Conf,蓝方=", YamlConfig.Conf.LList)
		//fmt.Println("yamlConfig.Mysql.Dns=", YamlConfig.Mysql.Dns)
		return true
	}

	//viperConfig := viper.New()
	//viperConfig.AddConfigPath("./internal/config")
	////viperConfig.AddConfigPath(".")
	//viperConfig.SetConfigName("app")
	//viperConfig.SetConfigType("yaml")
	//// 读取解析
	//if err := viperConfig.ReadInConfig(); err != nil {
	//	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
	//		fmt.Printf("app.yaml not found!%v\n", err)
	//		return
	//	} else {
	//		fmt.Printf("app.yaml read error,%v\n", err)
	//		return
	//	}
	//}
	//// 映射到结构体
	//if err := viperConfig.Unmarshal(&YamlConfig); err != nil {
	//	//fmt.Printf("app.yaml Unmarshal error,%v\n", err)
	//	log.Fatal("unmarshal config error:", err)
	//} else {
	//	fmt.Println("yamlConfig:ApiKey=", YamlConfig.Livekit.ApiKey)
	//	fmt.Println("yamlConfig:ApiSecret=", YamlConfig.Livekit.ApiCert)
	//	fmt.Println("yamlConfig.Host=", YamlConfig.Livekit.Host)
	//}

}
