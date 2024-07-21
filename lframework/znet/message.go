package znet

//Message 消息
type Message struct {
	DataLen uint32 //消息的长度
	ID      uint32 //消息的ID
	Sub uint32		//消息类型
	Data    []byte //消息的内容
}

//NewMsgPackage 创建一个Message消息包
func NewMsgPackage(ID uint32,Sub uint32 , data []byte) *Message {
	return &Message{
		DataLen: uint32(len(data)),
		ID:      ID,
		Sub:	Sub,
		Data:    data,
	}
}

//GetDataLen 获取消息数据段长度
func (msg *Message) GetDataLen() uint32 {
	return msg.DataLen
}

//GetMsgID 获取消息ID
func (msg *Message) GetMsgID() uint32 {
	return msg.ID
}

//GetMsgSub 获取消息Sub
func (msg *Message) GetMsgSub() uint32 {
	return msg.Sub
}

//GetData 获取消息内容
func (msg *Message) GetData() []byte {
	return msg.Data
}

//SetDataLen 设置消息数据段长度
func (msg *Message) SetDataLen(len uint32) {
	msg.DataLen = len
}

//SetMsgID 设计消息ID
func (msg *Message) SetMsgID(msgID uint32) {
	msg.ID = msgID
}


//SetMsgSub 设计消息Sub
func (msg *Message) SetMsgSub(msgSub uint32) {
	msg.Sub = msgSub
}


//SetData 设计消息内容
func (msg *Message) SetData(data []byte) {
	msg.Data = data
}
