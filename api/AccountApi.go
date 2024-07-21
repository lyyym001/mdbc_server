package api

import (
	"encoding/json"
	"mdbc_server/core"
	"mdbc_server/pb"

	"fmt"
	//"github.com/fd/lframework/utils"
	"mdbc_server/lframework/ziface"
	"mdbc_server/lframework/znet"
	//"log"
)

type AccountApi struct {
	znet.BaseRouter
}

func (aa *AccountApi) Handle(request ziface.IRequest) {

	//1. 得到消息的Sub，用来细化业务实现
	sub := request.GetMsgSub()
	//fmt.Println("Account Api Do : msgID = " , request.GetMsgID() , " Sub = " , request.GetMsgSub() , " msgLength = " , len(request.GetData()))

	//2. 得知当前的消息是从哪个玩家传递来的,从连接属性pID中获取
	pID, err := request.GetConnection().GetProperty("pID")
	if err != nil {
		fmt.Println("GetProperty pID error", err)
		request.GetConnection().Stop()
		return
	}
	//3. 根据pID得到player对象
	player := core.WorldMgrObj.GetPlayerByPID(pID.(int32))
	if player == nil {
		return
	}
	//fmt.Println("[Receive Account Msg] : Player = " , player.PID )

	switch sub {

	case 10002: //登录
		aa.Handle_onRequest10002(player, request.GetData())
		break

	}

}

// /成员主动离开(离开后无法继续进入)
func (aa *AccountApi) Handle_onRequest10002(p *core.Player, data []byte) {
	workRuning := 0
	request_data := &pb.Tcp_Login{}
	json.Unmarshal(data, request_data)
	fmt.Println("[请求登录]用户名=", request_data.UserName, ",类型(0-学生，1老师)=", request_data.AccountType)

	core.WorldMgrObj.Login(request_data.UserName, request_data.AccountType, p.PID)
	if core.WorldMgrObj.MScene != nil && core.WorldMgrObj.MScene.Running {
		workRuning = 1
	}
	p.Login(request_data.AccountType, request_data.UserName, workRuning)

}
