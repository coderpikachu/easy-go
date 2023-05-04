package process2

import (
	"easy-go/1internal/pkg/code"
	"easy-go/1internal/pkg/util/transfer"
	"easy-go/2pkg/model"
	"fmt"
	"net"

	"encoding/json"

	"github.com/marmotedu/errors"
)

type SmsProcess struct {
	//..[暂时不需字段]
}

//写方法转发消息
func (this *SmsProcess) SendGroupMes(mes *model.Message) (err error) {

	//遍历服务器端的onlineUsers map[int]*UserProcess,
	//将消息转发取出.
	//取出mes的内容 SmsMes
	var smsMes model.SmsMes
	err = json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return errors.WithCode(code.ErrJsonUnmarshal, "")
	}

	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return errors.WithCode(code.ErrJsonUnmarshal, "")
	}

	for id, up := range userMgr.onlineUsers {
		//这里，还需要过滤到自己,即不要再发给自己
		if id == smsMes.UserId {
			continue
		}
		this.SendMesToEachOnlineUser(data, up.Conn)
	}
	return nil
}
func (this *SmsProcess) SendMesToEachOnlineUser(data []byte, conn net.Conn) {

	//创建一个Transfer 实例，发送data
	tf := &transfer.Transfer{
		Conn: conn, //
	}
	err := tf.WritePkg(data)
	if err != nil {
		fmt.Println("转发消息失败 err=", err)
	}
}
