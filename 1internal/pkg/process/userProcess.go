package process2

import (
	"easy-go/1internal/pkg/util/transfer"
	"easy-go/2pkg/log"
	"easy-go/2pkg/model"
	"encoding/json"
	"net"
)

type UserProcess struct {
	//字段
	Conn net.Conn
	//增加一个字段，表示该Conn是哪个用户
	UserId int
}

//编写一个函数serverProcessLogin函数， 专门处理登录请求
func (this *UserProcess) ServerProcessLogin(mes *model.Message) (err error) {
	//核心代码...
	//1. 先从mes 中取出 mes.Data ，并直接反序列化成LoginMes
	var loginMes model.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		log.Debugf("json.Unmarshal fail err=", err)
		return
	}
	//1先声明一个 resMes
	var resMes model.Message
	resMes.Type = model.LoginResMesType
	//2在声明一个 LoginResMes，并完成赋值
	var loginResMes model.LoginResMes

	//我们需要到redis数据库去完成验证.
	//1.使用model.MyUserDao 到redis去验证

	log.Debugf("alive,prepare redis")
	//user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)
	if err != nil {

		// if err == model.ERROR_USER_NOTEXISTS {
		// 	loginResMes.Code = 500
		// 	loginResMes.Error = err.Error()
		// } else if err == model.ERROR_USER_PWD {
		// 	loginResMes.Code = 403
		// 	loginResMes.Error = err.Error()
		// } else {
		// 	loginResMes.Code = 505
		// 	loginResMes.Error = "服务器内部错误..."
		// }

	} else {
		loginResMes.Code = 200
		//这里，因为用户登录成功，我们就把该登录成功的用放入到userMgr中
		//将登录成功的用户的userId 赋给 this
		this.UserId = loginMes.UserId
		// userMgr.AddOnlineUser(this)
		// //通知其它的在线用户， 我上线了
		// this.NotifyOthersOnlineUser(loginMes.UserId)
		// //将当前在线用户的id 放入到loginResMes.UsersId
		// //遍历 userMgr.onlineUsers
		// for id, _ := range userMgr.onlineUsers {
		// 	loginResMes.UsersId = append(loginResMes.UsersId, id)
		// }
		//log.Debugf(user, "登录成功")
	}
	// //如果用户id= 100， 密码=123456, 认为合法，否则不合法

	// if loginMes.UserId == 100 && loginMes.UserPwd == "123456" {
	// 	//合法
	// 	loginResMes.Code = 200

	// } else {
	// 	//不合法
	// 	loginResMes.Code = 500 // 500 状态码，表示该用户不存在
	// 	loginResMes.Error = "该用户不存在, 请注册再使用..."
	// }

	//3将 loginResMes 序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		log.Debugf("json.Marshal fail", err)
		return
	}

	//4. 将data 赋值给 resMes
	resMes.Data = string(data)

	//5. 对resMes 进行序列化，准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		log.Debugf("json.Marshal fail", err)
		return
	}
	//6. 发送data, 我们将其封装到writePkg函数
	//因为使用分层模式(mvc), 我们先创建一个Transfer 实例，然后读取
	tf := &transfer.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return
}
