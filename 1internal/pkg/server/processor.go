package server

import (
	process2 "easy-go/1internal/pkg/process"
	"easy-go/1internal/pkg/util/transfer"
	"easy-go/2pkg/log"
	"easy-go/2pkg/model"
	"io"
	"net"
)

type Processor struct {
	Conn net.Conn
}

func (this *Processor) process2() (err error) {

	//循环的客户端发送的信息
	for {
		//这里我们将读取数据包，直接封装成一个函数readPkg(), 返回Message, Err
		//创建一个Transfer 实例完成读包任务
		tf := &transfer.Transfer{
			Conn: this.Conn,
		}
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				log.Debugf("5alive 客户端退出，服务器端也退出..")
				return err
			} else {
				log.Debugf("6alive readPkg err=", err)
				return err
			}

		}
		log.Debugf("7alive 开始处理数据包", err)
		err = this.serverProcessMes(&mes)
		if err != nil {
			return err
		}
	}

}
func (this *Processor) serverProcessMes(mes *model.Message) (err error) {

	//看看是否能接收到客户端发送的群发的消息
	log.Debugf("8alive mes=", mes)

	switch mes.Type {
	case model.LoginMesType:
		//处理登录登录
		//创建一个UserProcess实例
		up := &process2.UserProcess{
			Conn: this.Conn,
		}
		err = up.ServerProcessLogin(mes)
	case model.RegisterMesType:
		//处理注册
		up := &process2.UserProcess{
			Conn: this.Conn,
		}
		err = up.ServerProcessRegister(mes) // type : data
	case model.SmsMesType:
		//创建一个SmsProcess实例完成转发群聊消息.
		smsProcess := &process2.SmsProcess{}
		smsProcess.SendGroupMes(mes)
	default:
		log.Debugf("9 消息类型不存在，无法处理...")
	}
	return
}
