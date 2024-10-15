package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pion/turn/v3"
	"io"
	"log"
	"mdbc_server/api"
	"mdbc_server/core"
	"mdbc_server/internal/config"
	"mdbc_server/internal/models_sqlite"
	"mdbc_server/internal/server/router"
	"mdbc_server/lframework/utils"
	"mdbc_server/lframework/ziface"
	"mdbc_server/lframework/zlog"
	"mdbc_server/lframework/znet"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
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

	//fmt.Println(gameutils.GlobalScene)
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

	var AppID string = "MDBC"
	readed := config.Read()
	if readed {
		utils.GlobalObject.Host = config.YamlConfig.App.Host
		utils.GlobalObject.Name = config.YamlConfig.App.Name
		utils.GlobalObject.TCPPort = config.YamlConfig.App.Port
		utils.GlobalObject.MaxConn = config.YamlConfig.App.MaxConn
		utils.GlobalObject.WorkerPoolSize = config.YamlConfig.App.WorkerPoolSize
		utils.GlobalObject.LogFile = config.YamlConfig.App.LogFile
		utils.GlobalObject.MaxPacketSize = config.YamlConfig.App.MaxPacketSize
		utils.GlobalObject.UdpPort = config.YamlConfig.App.UdpPort
		fmt.Println("server setup ", utils.GlobalObject.Name)
		zlog.Debugf("Broadcast UdpPort = %d", utils.GlobalObject.UdpPort)
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
		zlog.Debug("1.db init")
		models_sqlite.NewDB()
		//models.NewDB()
		//zlog.Debug("2.db inited")
		//启动ginServices
		//zlog.Debug("3.License GinServer Start")
		zlog.Debug("2.SetupGinServer")
		go GinServices()
		//
		zlog.Debug("3.SetupStunServer")
		//go StunServer(config.YamlConfig.App.Host, config.YamlConfig.App.StunPort)
		//启动InteractionServices
		zlog.Debug("4.License InteractionServices Start")
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
	s.AddRouter(2, &api.WorkApi{})    //作品
	s.AddRouter(3, &api.DeviceApi{})  //设备
	s.AddRouter(4, &api.RoomApi{})    //房间
	//启动服务
	zlog.Debugf("6.InteractionServices run in %s , port = %d \n", utils.GlobalObject.Host, utils.GlobalObject.TCPPort)
	s.Serve()

}

// http - services
func GinServices() {

	gin.DisableConsoleColor()
	f, _ := os.Create("./log/gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

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

func StunServer(ip string, port int) {

	fmt.Println("Stun Server Ready To Start")

	publicIP := ip //flag.String("public-ip", "192.168.0.22", "IP Address that STUN can be contacted by.")
	//port           //flag.Int("port", 3478, "Listening port.")
	flag.Parse()

	if len(publicIP) == 0 {
		fmt.Println("'public-ip' is required")
		log.Fatalf("'public-ip' is required")
	}

	// Create a UDP listener to pass into pion/turn
	// pion/turn itself doesn't allocate any UDP sockets, but lets the user pass them in
	// this allows us to add logging, storage or modify inbound/outbound traffic
	udpListener, err := net.ListenPacket("udp4", "0.0.0.0:"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("Failed to create STUN server listener:", err)
		log.Panicf("Failed to create STUN server listener: %s", err)
	}
	fmt.Println("Stun Server Ready To NewServer")
	s, err := turn.NewServer(turn.ServerConfig{
		// PacketConnConfigs is a list of UDP Listeners and the configuration around them
		PacketConnConfigs: []turn.PacketConnConfig{
			{
				PacketConn: udpListener,
			},
		},
	})
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Stun Server Ready To ServerUped")
	fmt.Println("Listener Ip = " + publicIP)
	fmt.Println("Listener Port " + strconv.Itoa(port))
	// Block until user sends SIGINT or SIGTERM
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	if err = s.Close(); err != nil {
		log.Panic(err)
	}

}
