package main

import (
	"database/sql"
	"fmt"
	"mdbc_server/api"
	"mdbc_server/core"
	"mdbc_server/gameutils"
	"mdbc_server/lframework/utils"
	"mdbc_server/lframework/ziface"
	"mdbc_server/lframework/znet"
	"os"
	"strconv"
	"strings"
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

	//for i, v := range os.Args {
	//	fmt.Printf("args[%v]=%v\n", i, v)
	//}
	var AppID string = "1_1"

	if len(os.Args) > 1 {
		//命名中带有参数“本地服务ip$TcpPort$UdpPort”
		//var res []string
		res := strings.Split(os.Args[1], "$")
		utils.GlobalObject.Host = res[0]
		utils.GlobalObject.TCPPort, _ = strconv.Atoi(res[1])
		utils.GlobalObject.UdpPort, _ = strconv.Atoi(res[2])
		AppID = res[3]
		fmt.Println("os.Args = ", os.Args)
	}

	fmt.Println("Host    = ", utils.GlobalObject.Host)
	fmt.Println("TCPPort = ", utils.GlobalObject.TCPPort)
	fmt.Println("UdpPort = ", utils.GlobalObject.UdpPort)
	fmt.Println("LogDir  = ", utils.GlobalObject.LogDir)
	fmt.Println("LogFile = ", utils.GlobalObject.LogFile)
	fmt.Println("sqlite  = ", utils.GlobalObject.SqlitePath)
	//fmt.Println("LanServer AppID = ",AppID)
	//zlog.Debug("Hello")
	//从世界启动app
	core.WorldMgrObj.StartApp(AppID)
	//创建服务器句柄
	s := znet.NewServer()

	//注册客户端连接建立和丢失函数
	s.SetOnConnStart(OnConnecionAdd)
	s.SetOnConnStop(OnConnectionLost)

	//启动本地老师端
	//print("Start Up TClient\n")
	//_,err := os.StartProcess("StartTClient.bat",nil, &os.ProcAttr{Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}})
	//if err != nil {
	//	print("TClient Started Error\n")
	//}else {
	//	print("TClient Started Succ\n")
	//}
	test()
	//注册路由

	//登录路由
	s.AddRouter(1, &api.AccountApi{})
	//作品路由
	s.AddRouter(2, &api.WorkApi{})
	////聊天路由
	//s.AddRouter(2,&api.RoomApi{})
	////课程路由
	//s.AddRouter(3,&api.CourseApi{})
	//启动服务
	s.Serve()

}

func test() {

	//RegisterWorkRecord()

}

func RegisterWorkRecord() int {
	var maxUniqueId sql.NullInt32

	db := utils.GlobalObject.SqliteInst.GetDB()
	rows, _ := db.Query("select max(uniqueid) from tb_work; ")
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&maxUniqueId); err == nil {
			//allStudent.Students = append(allStudent.Students, studentInfo)
		} else {
			fmt.Println("RegisterWorkRecord,", err)
		}
	}

	fmt.Println("maxUniqueId = ", maxUniqueId.Valid)
	return 1
}
