package core

import (
	"encoding/json"
	"fmt"
	"mdbc_server/lframework/ziface"
	"mdbc_server/pb"
	"sync"
	"time"
)

// 玩家对象
type Player struct {

	//RID  int32				//房间ID
	//TID string				//老师ID
	//CID string				//当前Player的账号，有可能是学生账号，也有可能是老师账号 如果学生跟老师cid tid相同，则是老师

	//CLFlag bool				//这个用户是否是重复连接的用户
	//SNum string				//学号
	//PName string
	//CourseCore float32
	//CourseAbility string
	//CourseSetupDate int64     	//课件开始时间记录
	//CourseId        string    	//课程ID
	//CourseMode      string    	//课程模式

	//下面是new
	BeatTime int    //心跳检测时间
	UserName string //用户名
	//UName 			string		//真实姓名
	Conn        ziface.IConnection //当前玩家的连接
	AccountType int                //账号类型 0-学生 1-老师
	PID         int32              //玩家ID

	X float32 //平面x坐标
	Y float32 //高度
	Z float32 //平面y坐标 (注意不是Y)
	V float32 //旋转0-360度

}

/*
Player ID 生成器
*/
var PIDGen int32 = 1  //用来生成玩家ID的计数器
var IDLock sync.Mutex //保护PIDGen的互斥机制

// 创建一个玩家对象
func NewPlayer(conn ziface.IConnection) *Player {
	//生成一个PID
	IDLock.Lock()
	ID := PIDGen

	PIDGen++
	IDLock.Unlock()

	p := &Player{
		PID:  ID,
		Conn: conn,
		//CLFlag:false,
	}

	return p
}

// 通知错误信息 0-重复连接 1-授权未检测
func (p *Player) LoginError(flag int) {

	data, _ := json.Marshal(&pb.Tcp_Info{Flag: flag})
	p.SendMsg(1, 10002, data)
}

func (p *Player) Login(account_type int, uname string, workRuning int) {

	//fmt.Println("登录成功，" ,data)
	//1. player赋值
	//p.TID = data.TID
	//p.CID = data.CID
	p.UserName = uname
	//p.UName = uname
	//p.Number = number
	//p.Code = code
	p.AccountType = account_type
	p.BeatTime = int(time.Now().Unix())

	//回执
	data1, _ := json.Marshal(&pb.Tcp_Info{Flag: 1, UName: uname, WorkRuning: workRuning})
	p.SendMsg(1, 10002, data1)

	//2. Player加入房间(教师)
	//RoomMgrObj.AddPlayer(p)

	//如果是学生 ， 需要通知老师，学生连接上了
	//if p.CID != p.TID {
	//	player := RoomMgrObj.GetTPlayer(p.TID)
	//	if player != nil {
	//		data1,_ := json.Marshal(&pb.StudentInfo{StuUserName:p.CID,Flag:1})
	//		//
	//		////发送数据给客户端
	//		player.SendMsg(1,10007, data1)
	//	}
	//
	//
	//学生增加一个心跳
	//	go p.BroadcastPlayer()
	//
	//}else{
	//
	//	//老师进来，同步学生状态
	//	players := RoomMgrObj.GetAllPlayers(p.TID)
	//	if players != nil {
	//		for _,player := range players {
	//
	//			data1, _ := json.Marshal(&pb.StudentInfo{StuUserName: player.CID, Flag: 1})
	//			//
	//			////发送数据给客户端
	//			p.SendMsg(1, 10007, data1)
	//		}
	//	}
	//}

	//用户自己登录回调

	//Room := RoomMgrObj.GetRoom(p.TID)
	//if Room != nil {
	//	data1,_ := json.Marshal(&pb.SyncLoginB{CtrlFlag:Room.CtrlFlag,Code:"success"})
	//	////发送数据给客户端
	//	p.SendMsg(1,10002, data1)
	//}

	//3. 测试一下数据库===============
	//db := utils.GlobalObject.SqliteInst.GetDB()

	//// (1) QUERY
	//var username string
	//var isTeachClose int
	//rows, err := db.Query("select * FROM tb_user where userName = ?" , "t010001")
	//if err != nil {
	//	fmt.Println("Sqlite Test Query DB Err")
	//}else {
	//	defer rows.Close()
	//	for rows.Next() {
	//		if err := rows.Scan(&username, &isTeachClose); err != nil {
	//			fmt.Println("Sqlite Test GetData DB Err")
	//		}
	//	}
	//}
	//fmt.Println("username = " , username , " isteachclose = " , isTeachClose)
	//
	//
	//// (2) 更新
	//stmt , err := db.Prepare("UPDATE tb_user SET IsTeacherClose = ? WHERE username = ?")
	//if err != nil{
	//	fmt.Println("Sqlite Test Update DB Err")
	//}else {
	//	defer stmt.Close()
	//	result , err :=stmt.Exec(1,"t010001")
	//	affectNum, err := result.RowsAffected()
	//	if err != nil {
	//		fmt.Println("Sqlite Test affect DB Err")
	//	}
	//	fmt.Println("update affect rows is ", affectNum)
	//}

}

// 告知客户端被踢了
func (p *Player) Kicked() {

	p.SendMsg(1, 10003, []byte("ok"))
}

// 告知客户端pID,同步已经生成的玩家ID给客户端
func (p *Player) SyncPID() {

	////组建MsgID0 proto数据
	data1, _ := json.Marshal(&pb.SyncPID{PID: p.PID})
	//
	////发送数据给客户端
	fmt.Println("SendPID To Client", p.PID)
	p.SendMsg(1, 10001, data1)
}

// 广播玩家自己的出生地点
func (p *Player) BroadCastStartPosition() {

	//组建MsgID200 proto数据
	msg := &pb.BroadCast{
		PID: p.PID,
		Pos: pb.Position{
			X: p.X,
			Y: p.Y,
			Z: p.Z,
			V: p.V,
		},
	}
	data, _ := json.Marshal(msg)
	//发送数据给客户端
	p.SendMsg(1, 10004, data)
}

// 给当前玩家周边的(九宫格内)玩家广播自己的位置，让他们显示自己
func (p *Player) SyncSurrounding() {
	//1.获取玩家列表
	players := WorldMgrObj.GetAllPlayers()
	//for _, player := range players {
	//	player.SendMsg(2,10002 , data)
	//}
	//3.1 组建MsgID200 proto数据
	//组建MsgID200 proto数据
	msg := &pb.BroadCast{
		PID: p.PID,
		Pos: pb.Position{
			X: p.X,
			Y: p.Y,
			Z: p.Z,
			V: p.V,
		},
	}
	data, _ := json.Marshal(msg)
	//3.2 每个玩家分别给对应的客户端发送200消息，显示人物
	//告诉其他人我在哪
	for _, player := range players {
		if player.PID != p.PID {
			p.SendMsg(1, 10004, data)
		}
	}

	//4 让周围九宫格内的玩家出现在自己的视野中
	//4.1 制作Message SyncPlayers 数据
	var SyncPlayersMsg pb.SyncPlayers
	var datas []pb.BroadCast
	for _, player := range players {
		var p = pb.BroadCast{
			PID: player.PID,
			Pos: pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		datas = append(datas, p)
	}

	//4.2 封装SyncPlayer protobuf数据
	SyncPlayersMsg.Ps = datas
	data1, _ := json.Marshal(SyncPlayersMsg)
	//4.3 给当前玩家发送需要显示周围的全部玩家数据
	p.SendMsg(1, 10004, data1)
}

// 广播玩家位置移动
func (p *Player) UpdatePos(x float32, y float32, z float32, v float32) {

	//触发消失视野和添加视野业务
	//计算旧格子gID
	//oldGID := WorldMgrObj.AoiMgr.GetGIDByPos(p.X, p.Z)
	//计算新格子gID
	//newGID := WorldMgrObj.AoiMgr.GetGIDByPos(x, z)

	//更新玩家的位置信息
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v

	//if oldGID != newGID {
	//	//触发gird切换
	//	//把pID从就的aoi格子中删除
	//	WorldMgrObj.AoiMgr.RemovePIDFromGrID(int(p.PID), oldGID)
	//	//把pID添加到新的aoi格子中去
	//	WorldMgrObj.AoiMgr.AddPIDToGrID(int(p.PID), newGID)
	//
	//	_ = p.OnExchangeAoiGrID(oldGID, newGID)
	//}
	//
	////组装protobuf协议，发送位置给周围玩家
	//msg := &pb.BroadCast{
	//	PID: p.PID,
	//	Tp:  4, //4- 移动之后的坐标信息
	//	Data: &pb.BroadCast_P{
	//		P: &pb.Position{
	//			X: p.X,
	//			Y: p.Y,
	//			Z: p.Z,
	//			V: p.V,
	//		},
	//	},
	//}
	//
	////获取当前玩家周边全部玩家
	//players := p.GetSurroundingPlayers()
	////向周边的每个玩家发送MsgID:200消息，移动位置更新消息
	//for _, player := range players {
	//	player.SendMsg(200, msg)
	//}
}

////广播玩家聊天
//func (p *Player) Talk(content string) {
//
//	//拼接：
//	msg := "["+p.CID+"]：" + content
//
//	//1. 组建MsgID200 proto数据
//	data,_ := json.Marshal(&pb.Talk{Content:msg,})
//
//	//2. 得到当前世界所有的在线玩家
//	players := WorldMgrObj.GetAllPlayers()
//
//	fmt.Println("Talk -> players Length:" , len(players))
//
//	//3. 向所有的玩家发送MsgID:200消息
//	for _, player := range players {
//		player.SendMsg(2,10002 , data)
//	}
//}

// 玩家下线
func (p *Player) LostConnection() {

	//WorldMgrObj.LostConnection(p.PID,p.UserName)
	//pID, err := p.Conn.GetTCPConnection().GetProperty("pID")
	//if err != nil {
	//	fmt.Println("GetProperty pID error", err)
	//	request.GetConnection().Stop()
	//	return
	//}
	//_, err := p.Conn.GetTCPConnection(). .Write([]byte("After ping .....\n"))
	//if err != nil {
	//	fmt.Println("call back ping ping ping error")
	//	p.Conn.Stop()
	//}
	//if p.Conn.GetTCPConnection() {
	//
	//}
	p.Conn.Stop()
	p = nil
}

//广播玩家自己的出生地点
//func (p *Player) BroadCastStartPosition() {
//
//	//组建MsgID200 proto数据
//	msg := &pb.BroadCast{
//		PID: p.PID,
//		Tp:  2, //TP2 代表广播坐标
//		Data: &pb.BroadCast_P{
//			P: &pb.Position{
//				X: p.X,
//				Y: p.Y,
//				Z: p.Z,
//				V: p.V,
//			},
//		},
//	}
//
//	//发送数据给客户端
//	p.SendMsg(200, msg)
//}

//给当前玩家周边的(九宫格内)玩家广播自己的位置，让他们显示自己
//func (p *Player) SyncSurrounding() {
//	//1 根据自己的位置，获取周围九宫格内的玩家pID
//	pIDs := WorldMgrObj.AoiMgr.GetPIDsByPos(p.X, p.Z)
//	//2 根据pID得到所有玩家对象
//	players := make([]*Player, 0, len(pIDs))
//	//3 给这些玩家发送MsgID:200消息，让自己出现在对方视野中
//	for _, pID := range pIDs {
//		players = append(players, WorldMgrObj.GetPlayerByPID(int32(pID)))
//	}
//	//3.1 组建MsgID200 proto数据
//	msg := &pb.BroadCast{
//		PID: p.PID,
//		Tp:  2, //TP2 代表广播坐标
//		Data: &pb.BroadCast_P{
//			P: &pb.Position{
//				X: p.X,
//				Y: p.Y,
//				Z: p.Z,
//				V: p.V,
//			},
//		},
//	}
//	//3.2 每个玩家分别给对应的客户端发送200消息，显示人物
//	for _, player := range players {
//		player.SendMsg(200, msg)
//	}
//	//4 让周围九宫格内的玩家出现在自己的视野中
//	//4.1 制作Message SyncPlayers 数据
//	playersData := make([]*pb.Player, 0, len(players))
//	for _, player := range players {
//		p := &pb.Player{
//			PID: player.PID,
//			P: &pb.Position{
//				X: player.X,
//				Y: player.Y,
//				Z: player.Z,
//				V: player.V,
//			},
//		}
//		playersData = append(playersData, p)
//	}
//
//	//4.2 封装SyncPlayer protobuf数据
//	SyncPlayersMsg := &pb.SyncPlayers{
//		Ps: playersData[:],
//	}
//
//	//4.3 给当前玩家发送需要显示周围的全部玩家数据
//	p.SendMsg(202, SyncPlayersMsg)
//}

////广播玩家位置移动
//func (p *Player) UpdatePos(x float32, y float32, z float32, v float32) {
//
//	//触发消失视野和添加视野业务
//	//计算旧格子gID
//	oldGID := WorldMgrObj.AoiMgr.GetGIDByPos(p.X, p.Z)
//	//计算新格子gID
//	newGID := WorldMgrObj.AoiMgr.GetGIDByPos(x, z)
//
//	//更新玩家的位置信息
//	p.X = x
//	p.Y = y
//	p.Z = z
//	p.V = v
//
//	if oldGID != newGID {
//		//触发gird切换
//		//把pID从就的aoi格子中删除
//		WorldMgrObj.AoiMgr.RemovePIDFromGrID(int(p.PID), oldGID)
//		//把pID添加到新的aoi格子中去
//		WorldMgrObj.AoiMgr.AddPIDToGrID(int(p.PID), newGID)
//
//		_ = p.OnExchangeAoiGrID(oldGID, newGID)
//	}
//
//	//组装protobuf协议，发送位置给周围玩家
//	msg := &pb.BroadCast{
//		PID: p.PID,
//		Tp:  4, //4- 移动之后的坐标信息
//		Data: &pb.BroadCast_P{
//			P: &pb.Position{
//				X: p.X,
//				Y: p.Y,
//				Z: p.Z,
//				V: p.V,
//			},
//		},
//	}
//
//	//获取当前玩家周边全部玩家
//	players := p.GetSurroundingPlayers()
//	//向周边的每个玩家发送MsgID:200消息，移动位置更新消息
//	for _, player := range players {
//		player.SendMsg(200, msg)
//	}
//}

//func (p *Player) OnExchangeAoiGrID(oldGID, newGID int) error {
//	//获取就的九宫格成员
//	oldGrIDs := WorldMgrObj.AoiMgr.GetSurroundGrIDsByGID(oldGID)
//
//	//为旧的九宫格成员建立哈希表,用来快速查找
//	oldGrIDsMap := make(map[int]bool, len(oldGrIDs))
//	for _, grID := range oldGrIDs {
//		oldGrIDsMap[grID.GID] = true
//	}
//
//	//获取新的九宫格成员
//	newGrIDs := WorldMgrObj.AoiMgr.GetSurroundGrIDsByGID(newGID)
//	//为新的九宫格成员建立哈希表,用来快速查找
//	newGrIDsMap := make(map[int]bool, len(newGrIDs))
//	for _, grID := range newGrIDs {
//		newGrIDsMap[grID.GID] = true
//	}
//
//	//------ > 处理视野消失 <-------
//	offlineMsg := &pb.SyncPID{
//		PID: p.PID,
//	}
//
//	//找到在旧的九宫格中出现,但是在新的九宫格中没有出现的格子
//	leavingGrIDs := make([]*GrID, 0)
//	for _, grID := range oldGrIDs {
//		if _, ok := newGrIDsMap[grID.GID]; !ok {
//			leavingGrIDs = append(leavingGrIDs, grID)
//		}
//	}
//
//	//获取需要消失的格子中的全部玩家
//	for _, grID := range leavingGrIDs {
//		players := WorldMgrObj.GetPlayersByGID(grID.GID)
//		for _, player := range players {
//			//让自己在其他玩家的客户端中消失
//			player.SendMsg(201, offlineMsg)
//
//			//将其他玩家信息 在自己的客户端中消失
//			anotherOfflineMsg := &pb.SyncPID{
//				PID: player.PID,
//			}
//			p.SendMsg(201, anotherOfflineMsg)
//			time.Sleep(200 * time.Millisecond)
//		}
//	}
//
//	//------ > 处理视野出现 <-------
//
//	//找到在新的九宫格内出现,但是没有在就的九宫格内出现的格子
//	enteringGrIDs := make([]*GrID, 0)
//	for _, grID := range newGrIDs {
//		if _, ok := oldGrIDsMap[grID.GID]; !ok {
//			enteringGrIDs = append(enteringGrIDs, grID)
//		}
//	}
//
//	onlineMsg := &pb.BroadCast{
//		PID: p.PID,
//		Tp:  2,
//		Data: &pb.BroadCast_P{
//			P: &pb.Position{
//				X: p.X,
//				Y: p.Y,
//				Z: p.Z,
//				V: p.V,
//			},
//		},
//	}
//
//	//获取需要显示格子的全部玩家
//	for _, grID := range enteringGrIDs {
//		players := WorldMgrObj.GetPlayersByGID(grID.GID)
//
//		for _, player := range players {
//			//让自己出现在其他人视野中
//			player.SendMsg(200, onlineMsg)
//
//			//让其他人出现在自己的视野中
//			anotherOnlineMsg := &pb.BroadCast{
//				PID: player.PID,
//				Tp:  2,
//				Data: &pb.BroadCast_P{
//					P: &pb.Position{
//						X: player.X,
//						Y: player.Y,
//						Z: player.Z,
//						V: player.V,
//					},
//				},
//			}
//
//			time.Sleep(200 * time.Millisecond)
//			p.SendMsg(200, anotherOnlineMsg)
//		}
//	}
//
//	return nil
//}

//获得当前玩家的AOI周边玩家信息
//func (p *Player) GetSurroundingPlayers() []*Player {
//	//得到当前AOI区域的所有pID
//	pIDs := WorldMgrObj.AoiMgr.GetPIDsByPos(p.X, p.Z)
//
//	//将所有pID对应的Player放到Player切片中
//	players := make([]*Player, 0, len(pIDs))
//	for _, pID := range pIDs {
//		players = append(players, WorldMgrObj.GetPlayerByPID(int32(pID)))
//	}
//
//	return players
//}

/*
发送消息给客户端，
主要是将pb的protobuf数据序列化之后发送
*/
func (p *Player) SendMsg(msgID uint32, msgSub uint32, data []byte) {
	//fmt.Printf("before Marshal data = %+v\n", data)
	//将NetBody结构体序列化
	//创建一个存放bytes字节的缓冲

	//msg, err := proto.Marshal(data)
	//if err != nil {
	//	fmt.Println("marshal msg err: ", err)
	//	return
	//}
	////fmt.Printf("after Marshal data = %+v\n", msg)
	//
	if p.Conn == nil {
		fmt.Println("connection in player is nil")
		return
	}

	//调用Zinx框架的SendMsg发包
	if err := p.Conn.SendMsg(msgID, msgSub, data); err != nil {
		fmt.Println("Player SendMsg error !")
		return
	}
	if msgSub == 20006 || msgSub == 20009 {
		fmt.Println("[Send Msg] To Player = ", p.UserName, " MsgId = ", msgID, " MsgSub = ", msgSub, " dataLength = ", len(data), " NowDate = ", time.Now().Format("2006-01-02 15:04:05"))
	}

	return
}

//func (p *Player) BroadcastPlayer()  {
//
//
//
//
//	t1 := time.NewTimer(time.Millisecond * 5000) //5s
//L:
//	for {
//		if p == nil || p.TID == p.CID {
//			break L
//		}
//		//fmt.Println("len(cRoom.ClientHandle), ", len(cRoom.ClientHandle), "cRoom.TeacherCli,", cRoom.TeacherCli)
//		/*if len(cRoom.ClientHandle) == 0 && cRoom.TeacherCli == nil{
//			goto ForEnd
//		}*/
//		select {
//		case <-t1.C:
//			t1.Reset(time.Millisecond * 5000)
//			//SendUdpBroadcastToStudent(rm.TID)
//			if p == nil{
//				break L
//			}
//			_BeatTime := int(time.Now().Unix())
//			if _BeatTime - p.BeatTime > 22 {
//				//fmt.Println("学生端掉线了")
//				fmt.Println("学生掉线了 ： " ,p.PID,p.CID , p.TID)
//				//p.Conn.Stop()
//				p.LostConnection()
//				//p.Conn.Stop()
//
//				break L
//			}
//		}
//	}
//
//}
