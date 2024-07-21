package core

import (
	"encoding/json"
	"fmt"
	"mdbc_server/lframework/utils"
	"mdbc_server/pb"
	"net"
	"sync"
	"time"
)

type SceneInfo struct {
	SceneID      string //当前场景id
	Mode         byte   //1-学练模式 2-考评模式 3-协同模式
	MainCtroller string //主控玩家
	Running      bool   //是否正在运行
	Players      map[string]int32
	TakeObjects  map[int]int32
	JHObjects    map[int]*JHObject
	Steps        map[int]*Step
	Questions    map[int]int32
	//GlobalStepId int
}

/*
当前游戏世界的总管理模块
*/
type WorldManager struct {
	AoiMgr    *AOIManager       //当前世界地图的AOI规划管理器
	Players   map[int32]*Player //当前在线的玩家集合
	PMs       map[string]int32  //记录登录的玩家，防止重复登录,map[UserName]Pid
	pLock     sync.RWMutex      //保护Players的互斥读写机制
	AppID     string            //用户登录标识，无此标识则不允许登录
	TUserName string            //老师账号
	TUid      int32             //老师UID
	MScene    *SceneInfo        //当前场景状态
}

// 提供一个对外的世界管理模块句柄
var WorldMgrObj *WorldManager

// 提供WorldManager 初始化方法
func init() {
	WorldMgrObj = &WorldManager{
		Players: make(map[int32]*Player),
		AoiMgr:  NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_CNTS_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNTS_Y),
		PMs:     make(map[string]int32),
		MScene:  &SceneInfo{SceneID: "", Mode: 1, MainCtroller: "", Running: false},
	}
}

// 提供添加一个玩家的的功能，将玩家添加进玩家信息表Players
func (wm *WorldManager) StartApp(appid string) {
	//将player添加到 世界管理器中
	wm.pLock.Lock()
	wm.AppID = appid
	wm.pLock.Unlock()

	//将player 添加到AOI网络规划中
	//wm.AoiMgr.AddToGrIDByPos(int(player.PID), player.X, player.Z)

	//启动udp用来广播自己的ip及端口
	go wm.BroadcastNet()

}

func (wm *WorldManager) BroadcastNet() {

	t1 := time.NewTimer(time.Millisecond * 1000) //1s
	//L:
	for {

		select {
		case <-t1.C:
			t1.Reset(time.Millisecond * 1000)
			SendUdpBroadcastToAll()
		}
	}
}

func SendUdpBroadcastToAll() {
	//fmt.Println("SendUdpBroadcast To Student tid = ",tid)
	//log.Println("SendUdpBroadcast To Student")
	//utils.GlobalObject.Host
	//laddrStu := net.UDPAddr{
	//	IP:   net.ParseIP(utils.GlobalObject.Host),
	//	Port: utils.GlobalObject.UdpPort,
	//}

	var sData pb.Sync_Hello
	sData.Ip = utils.GlobalObject.Host
	sData.Port = utils.GlobalObject.TCPPort

	// 这里设置接收者的IP地址为广播地址
	raddrStu := net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: utils.GlobalObject.UdpPort,
	}
	//fmt.Println("SendUdpBroadcast To Student raddrStu = ",laddrStu," raddrStu = ",raddrStu)
	connStu, err := net.DialUDP("udp", nil, &raddrStu) //&laddrStu
	if err != nil {
		println(err.Error())
		return
	}

	//fmt.Println("pBroadcast LanServer " ,raddrStu,"data = ", sData)
	x, _ := json.Marshal(sData)
	connStu.Write(x)
	connStu.Close()
}

// 提供添加一个玩家的的功能，将玩家添加进玩家信息表Players
func (wm *WorldManager) AddPlayer(player *Player) {
	//将player添加到 世界管理器中
	wm.pLock.Lock()
	wm.Players[player.PID] = player
	//wm.PMs[player.CID] = true
	wm.pLock.Unlock()

	//将player 添加到AOI网络规划中
	//wm.AoiMgr.AddToGrIDByPos(int(player.PID), player.X, player.Z)
}

// username 用户名
// accountType 账号类型 0-学生 1-老师
func (wm *WorldManager) Login(username string, accountType int, uid int32) {

	//将player添加到 世界管理器中
	wm.pLock.Lock()
	//wm.Players[player.PID] = player
	//wm.PMs[player.CID] = true
	if _, ok := wm.PMs[username]; ok {
		fmt.Println("[重复登录]username=", username, ",原pid=", wm.PMs[username], ",现pid=", uid)
		if wm.PMs[username] == uid {
			fmt.Println("pid重复了")
		} else {
			p := wm.GetPlayerByPID(uid)
			//旧的句柄释放掉
			if p != nil {
				//重复登录踢掉之前的
				fmt.Println("重复登录踢掉之前的")
				p.Kicked()
				wm.LostConnection(wm.PMs[username], username, 1)
				p.LostConnection()
			}
		}
		//rcode = false

	}
	if accountType == 1 {
		//记录老师状态
		wm.TUserName = username
		wm.TUid = uid
		//fmt.Println("[记录老师状态]")
	}
	//用户加到世界列表
	wm.PMs[username] = uid
	//检测作品成员列表
	if wm.MScene.Running {
		if _, ok := wm.MScene.Players[username]; ok {
			//delete(wm.MScene.Players,pID)
			wm.MScene.Players[username] = 1
		}
	}
	wm.pLock.Unlock()

	wm.PrintUserList()

}

func (wm *WorldManager) PrintUserList() {
	fmt.Println("[世界成员列表]")
	for _key, value := range wm.PMs {
		fmt.Println("->username=", _key, ",uid=", value)
	}

	if wm.MScene != nil && wm.MScene.Players != nil && len(wm.MScene.Players) > 0 {
		fmt.Println("[当前作品成员列表]")
		for _key, value := range wm.MScene.Players {
			fmt.Println("->username=", _key, ",状态(0-资源未加载，1-资源已加载)=", value)
		}
	} else {
		fmt.Println("[当前作品成员列表](无)")
	}
}

// 玩家掉线了
// code 0-掉线 1-被踢
func (wm *WorldManager) LostConnection(pID int32, userName string, code int) {

	wm.RemovePlayer(pID)
	wm.RemovePlayerUName(userName)
	wm.RemoveWorkPlayer(userName)

	wm.PrintUserList()

	if code == 0 {
		if wm.MScene.Running {
			if wm.MScene.MainCtroller == userName {
				fmt.Println("[离线重新指定主控]原主控=", userName)
				wm.ChangeMainCtroller(1)
			} else {
				wm.CheckOver(1)
			}
		}
	}
}

// code 0-用户中途离开指定主控 1-用户掉线指定主控
func (wm *WorldManager) ChangeMainCtroller(code int) bool {

	//重新指定一个主控
	for username, state := range wm.MScene.Players {

		if state == 1 {
			if uid, ok := wm.PMs[username]; ok {
				p := wm.GetPlayerByPID(uid)
				if p != nil {
					fmt.Println("[离线重新指定主控]新主控=", username)
					p.SendMsg(2, 10008, []byte("ok"))
				}
				return true
			}
		}
	}

	//没有人了，结束作品
	wm.WorkFinish(0)
	return false

}

// code 0-用户离开了检测结束 1-用户离线了检测结束
func (wm *WorldManager) CheckOver(code int) bool {

	//检测是否还有人
	for _, state := range wm.MScene.Players {

		if state == 1 {
			return false
		}
	}

	//没有找到有效的主控则结束课程
	wm.WorkFinish(0)

	return true
}

// 世界同步给所有参与人的消息
func (wm *WorldManager) WorldToa(msgID uint32, msgSub uint32, data []byte) {

	for _, player := range wm.Players {
		if player != nil {
			player.SendMsg(msgID, msgSub, data)
		}
	}
}

// 作品内同步给所有参与人的消息
func (wm *WorldManager) Toa(msgID uint32, msgSub uint32, data []byte) {

	//fmt.Println("Toa->消息ID=",msgID,",SubId=",msgSub)
	//wm.PrintUserList()
	if wm.MScene != nil {
		if wm.MScene.Running {
			players := wm.MScene.Players
			if players != nil && len(players) > 0 {
				for username, _ := range players {
					//if state == 1 {
					p := wm.GetPlayerByUserName(username)
					if p != nil {
						p.SendMsg(msgID, msgSub, data)
					}
					//}
				}
			}
		}
	}
}

// 作品内同步给其他参与人的消息
func (wm *WorldManager) Too(uid int32, msgID uint32, msgSub uint32, data []byte) {
	//fmt.Println("Too->消息ID=",msgID,",SubId=",msgSub)

	if wm.MScene != nil {
		if wm.MScene.Running {
			players := wm.MScene.Players
			if players != nil && len(players) > 0 {
				for username, _ := range players {
					p := wm.GetPlayerByUserName(username)
					if p != nil {
						if p.PID != uid {
							//if state == 1 {
							p.SendMsg(msgID, msgSub, data)
							//}
						}
					}
				}
			}
		}
	}
}

// 从玩家信息表中移除一个玩家
func (wm *WorldManager) RemovePlayer(pID int32) {
	wm.pLock.Lock()
	delete(wm.Players, pID)
	wm.pLock.Unlock()
}

func (wm *WorldManager) RemovePlayerUName(userName string) {
	wm.pLock.Lock()
	if _, ok := wm.PMs[userName]; ok {
		delete(wm.PMs, userName)
	}
	wm.pLock.Unlock()
}

// 移除作品中的玩家(假移除，可以让用户重连进作品)
func (wm *WorldManager) RemoveWorkPlayer(userName string) {
	wm.pLock.Lock()
	//移除参与玩家
	if wm.MScene.Running {
		if _, ok := wm.MScene.Players[userName]; ok {
			wm.MScene.Players[userName] = 2 //表示当前参与玩家离线了
		}
	}
	wm.pLock.Unlock()
}

func (wm *WorldManager) PullUserToCell() {
	if wm.Players != nil && len(wm.Players) > 0 {
		for username, _ := range wm.PMs {
			wm.MScene.Players[username] = 0
		}
	}
}

func (wm *WorldManager) PullListToCell(users []string) {

	for _, userName := range users {
		//fmt.Println("userName=",userName,wm.PMs)
		if _, ok := wm.PMs[userName]; ok {
			//fmt.Println("userName1=",userName)
			wm.MScene.Players[userName] = 1
		}
	}
}

func (wm *WorldManager) UpdateReadedStatus(pid int32) {
	if wm.MScene.Running {

		p := wm.GetPlayerByPID(pid)
		if p != nil {
			if _, ok := wm.MScene.Players[p.UserName]; ok {
				wm.MScene.Players[p.UserName] = 1 //表示已经准备好了的
				fmt.Println(p.UserName, "[准备完毕]")
			}
		}
	}

}

// 准备好的人数
func (wm *WorldManager) AllReaded() bool {

	if wm.MScene.Players != nil && len(wm.MScene.Players) > 0 {
		for _, readed := range wm.MScene.Players {
			if readed == 0 {
				fmt.Println("有未准备好的成员，全部人数：", len(wm.MScene.Players))
				return false
			}
		}
		return true
	}
	return false
}

// 通过玩家ID 获取对应玩家信息
func (wm *WorldManager) GetPlayerByPID(pID int32) *Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	return wm.Players[pID]
}

// 通过玩家ID 获取对应玩家信息
func (wm *WorldManager) GetPlayerByUserName(username string) *Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()
	if uid, ok := wm.PMs[username]; ok {
		return wm.Players[uid]
	}
	return nil
}

// 获取所有玩家的信息
func (wm *WorldManager) GetAllPlayers() []*Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	//创建返回的player集合切片
	players := make([]*Player, 0)

	//添加切片
	for _, v := range wm.Players {
		players = append(players, v)
	}

	//返回
	return players
}

// 获取指定gID中的所有player信息
func (wm *WorldManager) GetPlayersByGID(gID int) []*Player {
	//通过gID获取 对应 格子中的所有pID
	pIDs := wm.AoiMgr.grIDs[gID].GetPlyerIDs()

	//通过pID找到对应的player对象
	players := make([]*Player, 0, len(pIDs))
	wm.pLock.RLock()
	for _, pID := range pIDs {
		players = append(players, wm.Players[int32(pID)])
	}
	wm.pLock.RUnlock()

	return players
}

// 通过玩家ID 获取对应玩家信息
func (wm *WorldManager) HasLogined(cid string) bool {

	_, ok := wm.PMs[cid]

	return ok

}

// ==================业务====================
func (wm *WorldManager) RegisterObject(data *pb.TCP_RegisterObj) {

	if _, ok := wm.MScene.JHObjects[data.ObjId]; !ok {
		obj := &JHObject{
			ObjId:           data.ObjId,
			InteractiveType: data.InteractiveType,
			Tb:              data.Tb,
			Visiable:        data.Visiable,
			X:               data.X,
			Y:               data.Y,
			Z:               data.Z,
			RX:              data.RX,
			RY:              data.RY,
			RZ:              data.RZ,
		}
		wm.MScene.JHObjects[data.ObjId] = obj
	}

}

func (wm *WorldManager) UpdateObjectPos(data *pb.TCP_TbObj) {

	if _, ok := wm.MScene.JHObjects[data.ObjId]; ok {
		wm.MScene.JHObjects[data.ObjId].X = data.X
		wm.MScene.JHObjects[data.ObjId].Y = data.Y
		wm.MScene.JHObjects[data.ObjId].Z = data.Z
		wm.MScene.JHObjects[data.ObjId].RX = data.RX
		wm.MScene.JHObjects[data.ObjId].RY = data.RY
		wm.MScene.JHObjects[data.ObjId].RZ = data.RZ
	}

}

func (wm *WorldManager) UpdateObjectStatus(data *pb.Tcp_ObjectStatus) {

	if _, ok := wm.MScene.JHObjects[data.ObjId]; ok {
		wm.MScene.JHObjects[data.ObjId].Visiable = data.Status
	}

}

func (wm *WorldManager) RegisterStep(data *pb.Tcp_Step) {

	//_id := wm.MScene.GlobalStepId+1
	//wm.MScene.GlobalStepId = _id
	obj := &Step{
		StepId:    data.StepId,
		StepState: data.StepState,
		StepDate:  data.StepDate,
		UName:     data.UName,
	}
	wm.MScene.Steps[data.StepId] = obj

}

// code 0-服务器检测没人了结束作品 1-老师控制结束作品
func (wm *WorldManager) WorkFinish(code int) {
	fmt.Println("[作品结束了]code=", code)
	if wm.MScene.Running {

		wm.MScene.Running = false
	}
	wm.PrintUserList()
}

// 个人离开活动
// _state //0-中途离开 1-结束离开(结束离开如果是多人场景不重新指定主控)
func (wm *WorldManager) WorkLeave(pID int32, _state int) {

	//移除参与玩家
	if wm.MScene != nil && wm.MScene.Running {
		p := wm.GetPlayerByPID(pID)
		if p != nil {
			fmt.Println("[用户离开作品]pid=", pID, ",username=", p.UserName)
			if _, ok := wm.MScene.Players[p.UserName]; ok {
				delete(wm.MScene.Players, p.UserName)
			}
		}

		if _state != 1 {
			//中途离开的话需要检测
			if wm.MScene.MainCtroller == p.UserName {
				fmt.Println("[中途离开重新指定主控]原主控=", p.UserName)
				wm.ChangeMainCtroller(1)
			} else {
				wm.CheckOver(1)
			}
		} else {
			wm.CheckOver(1)
		}
	}

}
