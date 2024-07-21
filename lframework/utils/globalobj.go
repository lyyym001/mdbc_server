// Package utils 提供zinx相关工具类函数
// 包括:
//
//	全局配置
//	配置文件加载
//
// 当前文件描述:
// @Title  globalobj.go
// @Description  相关配置文件定义及加载方式
// @Author  Aceld - Thu Mar 11 10:32:29 CST 2019
package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"mdbc_server/lframework/ziface"
	"mdbc_server/lframework/zlog"
)

/*
存储一切有关Zinx框架的全局参数，供其他模块使用
一些参数也可以通过 用户根据 zinx.json来配置
*/
type GlobalObj struct {
	/*
		Server
	*/
	TCPServer ziface.IServer //当前Zinx的全局Server对象

	Host    string //当前服务器主机IP
	TCPPort int    //当前服务器主机监听端口号
	Name    string //当前服务器名称

	/*
		Zinx
	*/
	Version          string //当前Zinx版本号
	MaxPacketSize    uint32 //都需数据包的最大值
	MaxConn          int    //当前服务器主机允许的最大链接个数
	WorkerPoolSize   uint32 //业务工作Worker池的数量
	MaxWorkerTaskLen uint32 //业务工作Worker对应负责的任务队列最大任务存储数量
	MaxMsgChanLen    uint32 //SendBuffMsg发送消息的缓冲最大长度

	/*
		config file path
	*/
	ConfFilePath string

	/*
		Sqlite3
	*/
	SqliteUse  bool
	SqlitePath string
	SqliteInst ziface.ISqliteHandle //sqliteHandle
	/*
		logger
	*/
	LogDir        string //日志所在文件夹 默认"./log"
	LogFile       string //日志文件名称   默认""  --如果没有设置日志文件，打印信息将打印至stderr
	LogDebugClose bool   //是否关闭Debug日志级别调试信息 默认false  -- 默认打开debug信息

	/*
		udp 端口
	*/
	UdpPort    int //udp端口
	UdpPortDir int //第三方接口UDP协议
}

/*
定义一个全局的对象
*/
var GlobalObject *GlobalObj

// PathExists 判断一个文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Reload 读取用户的配置文件
func (g *GlobalObj) Reload() {

	//fmt.Println("configPath = ",g.ConfFilePath)
	if confFileExists, _ := PathExists(g.ConfFilePath); confFileExists != true {
		fmt.Println("Config File ", g.ConfFilePath, " is not exist!!")
		return
	}
	data, err := ioutil.ReadFile(g.ConfFilePath)
	if err != nil {
		panic(err)
	}
	//将json数据解析到struct中
	err = json.Unmarshal(data, g)
	if err != nil {
		panic(err)
	}

	//Logger 设置
	if g.LogFile != "" {
		fmt.Println("log - dir = ", g.LogDir, " LogFile = ", g.LogFile)
		zlog.SetLogFile(g.LogDir, g.LogFile)
	}
	if g.LogDebugClose == true {
		zlog.CloseDebug()
	}
	fmt.Println("1.ConfigInited")
}

/*
提供init方法，默认加载
*/
func init() {
	pwd, err := os.Getwd()
	if err != nil {
		pwd = "."
	}

	fmt.Println("ServerPath=", pwd)

	//初始化GlobalObject变量，设置一些默认值
	GlobalObject = &GlobalObj{
		Name:             "FDLFrameworkServer",
		Version:          "V0.11",
		TCPPort:          10106,
		Host:             LocalIp().String(),
		MaxConn:          12000,
		MaxPacketSize:    10240,
		ConfFilePath:     pwd + "/conf/config.json", //conf/config.json   //测试 lyyym/gameframework/conf/config.json
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		MaxMsgChanLen:    1024,
		LogDir:           pwd + "/log",
		LogFile:          "",
		LogDebugClose:    false,
		SqliteUse:        true,
		UdpPortDir:       10108,
	}
	//GlobalObject.Host = LocalIp().String()
	//NOTE: 从配置文件中加载一些用户配置的参数
	//fmt.Println("config Inited",GlobalObject)
	GlobalObject.Reload()
}

func LocalIp() net.IP {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var ip net.IP = nil
	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						ip = ipnet.IP
						break
						//fmt.Println(ip)
					}
				}
			}
		}
		if ip != nil {
			break
		}
	}

	return ip
}
