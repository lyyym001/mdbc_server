package main

import (
	"fmt"
	"log"
	"mdbc_server/api"
	"mdbc_server/core"
	"mdbc_server/gameutils"
	"mdbc_server/internal/config"
	"mdbc_server/internal/server/router"
	"mdbc_server/lframework/utils"
	"mdbc_server/lframework/ziface"
	"mdbc_server/lframework/zlog"
	"mdbc_server/lframework/znet"
)

//业务Api 这里定义跟客户都安通信的业务关联
//1	-	登录账号相关
//2 - 	房间业务

// 当客户端建立连接的时候的hook函数
func OnConnecionAdd(conn ziface.IConnection) {
	//创建一个玩家
	player := core.NewPlayer(conn)

	//同步当前玩家的初始化坐标信息给客户端，走MsgID:200消息
	//player.BroadCastStartPosition()

	//将当前新上线玩家添加到worldManager中
	core.WorldMgrObj.AddPlayer(player)

	//将该连接绑定属性PID
	conn.SetProperty("pID", player.PID)

	//同步周边玩家上线信息，与现实周边玩家信息
	//player.SyncSurrounding()

	//同步当前的PlayerID给客户端， 走MsgID:1 消息 这里需要客户端回执 登录信息
	player.SyncPID()

	//==============同步周边玩家上线信息，与现实周边玩家信息========
	//player.SyncSurrounding()

	fmt.Println("=====> Player pIDID = ", player.PID, " arrived ====")

	fmt.Println(gameutils.GlobalScene)
}

// 当客户端断开连接的时候的hook函数
func OnConnectionLost(conn ziface.IConnection) {

	fmt.Println("有客户端断开了连接")
	//获取当前连接的PID属性
	pID, _ := conn.GetProperty("pID")
	//fmt.Println("pID = " , pID)

	//根据pID获取对应的玩家对象
	player := core.WorldMgrObj.GetPlayerByPID(pID.(int32))
	if player != nil {
		//fmt.Println(player)
		fmt.Println("[断开连接用户]pid= ", player.PID, ",UserName =", player.UserName)
		//触发玩家下线业务
		core.WorldMgrObj.LostConnection(player.PID, player.UserName, 0)
		player.LostConnection()
	}

	//fmt.Println("====> Player ", pID, " left =====")
	//fmt.Println("123")
}

func main() {

	var AppID string = "MultiAttDemo"
	readed := config.Read()
	if readed {
		utils.GlobalObject.Host = config.YamlConfig.App.Host
		utils.GlobalObject.Name = config.YamlConfig.App.Name
		utils.GlobalObject.TCPPort = config.YamlConfig.App.Port
		utils.GlobalObject.MaxConn = config.YamlConfig.App.MaxConn
		utils.GlobalObject.WorkerPoolSize = config.YamlConfig.App.WorkerPoolSize
		utils.GlobalObject.LogFile = config.YamlConfig.App.LogFile
		utils.GlobalObject.MaxPacketSize = config.YamlConfig.App.MaxPacketSize
		fmt.Println("ready to setup server=", utils.GlobalObject.Name)

		//开启日志
		if utils.GlobalObject.LogFile != "" {
			fmt.Println("log-> dir= ", utils.GlobalObject.LogDir, " file= ", utils.GlobalObject.LogFile)
			zlog.SetLogFile(utils.GlobalObject.LogDir, utils.GlobalObject.LogFile)
		}
		//if g.LogDebugClose == true {
		//	zlog.CloseDebug()
		//}

		//
		//livekitclient.NewRoomClient()
		//zlog.Debug("1.db init")
		//models.NewDB()
		//zlog.Debug("2.db inited")
		//启动ginServices
		//zlog.Debug("3.License GinServer Start")
		zlog.Debug("1.SetupGinServer")
		go GinServices()
		//启动InteractionServices
		zlog.Debug("5.License InteractionServices Start")
		InteractionServices(AppID)
	}

}

func InteractionServices(AppID string) {

	//从世界启动app
	core.WorldMgrObj.StartApp(AppID)
	//创建服务器句柄
	s := znet.NewServer()
	//注册客户端连接建立和丢失函数
	s.SetOnConnStart(OnConnecionAdd)
	s.SetOnConnStop(OnConnectionLost)
	//注册路由
	s.AddRouter(1, &api.AccountApi{}) //登录路由
	s.AddRouter(2, &api.WorkApi{})    //房间
	//启动服务
	zlog.Debugf("6.InteractionServices run in %s , port = %d \n", utils.GlobalObject.Host, utils.GlobalObject.TCPPort)
	s.Serve()

}

// http - services
func GinServices() {

	e := router.Router()
	//e.Use(TlsHandler(8080))
	ginHost := fmt.Sprintf("%s:%d", config.YamlConfig.App.Host, config.YamlConfig.App.GinPort)
	fmt.Println("ginHost = ", ginHost)
	zlog.Debugf(" GinServices run in = %s\n", ginHost)
	err := e.Run(ginHost)
	//err := e.RunTLS(config.YamlConfig.App.Host, "./cert/server.pem", "./cert/server.key")
	if err != nil {
		log.Fatalln("run err.", err)
		return
	}
	//zlog.Debug("4.GinServices run in ", ginHost)

	fmt.Println("GinServices run in ", ginHost)
}
