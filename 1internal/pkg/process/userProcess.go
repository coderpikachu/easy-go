package process2

import (
	"easy-go/1internal/pkg/code"
	"easy-go/1internal/pkg/util/transfer"
	"easy-go/2pkg/log"
	"easy-go/2pkg/model"
	"encoding/json"
	"net"

	"github.com/marmotedu/errors"
)

type UserProcess struct {
	//字段
	Conn net.Conn
	//增加一个字段，表示该Conn是哪个用户
	UserId int
}

func (this *UserProcess) NotifyOthersOnlineUser(userId int) {

	//遍历 onlineUsers, 然后一个一个的发送 NotifyUserStatusMes
	for id, up := range userMgr.onlineUsers {
		//过滤到自己
		if id == userId {
			continue
		}
		//开始通知【单独的写一个方法】
		up.NotifyMeOnline(userId)
	}
}

func (this *UserProcess) NotifyMeOnline(userId int) (err error) {

	//组装我们的NotifyUserStatusMes
	var mes model.Message
	mes.Type = model.NotifyUserStatusMesType

	var notifyUserStatusMes model.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = model.UserOnline

	//将notifyUserStatusMes序列化
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		log.Debugf("json.Marshal err=", err)
		return errors.WithCode(code.ErrJsonUnmarshal, "")
	}
	//将序列化后的notifyUserStatusMes赋值给 mes.Data
	mes.Data = string(data)

	//对mes再次序列化，准备发送.
	data, err = json.Marshal(mes)
	if err != nil {
		log.Debugf("json.Marshal err=", err)
		return errors.WithCode(code.ErrJsonUnmarshal, "")
	}

	//发送,创建我们Transfer实例，发送
	tf := &transfer.Transfer{
		Conn: this.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		log.Debugf("NotifyMeOnline err=", err)
		return
	}
	return nil
}
func (this *UserProcess) ServerProcessRegister(mes *model.Message) (err error) {

	//1.先从mes 中取出 mes.Data ，并直接反序列化成RegisterMes
	var registerMes model.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		log.Debugf("json.Unmarshal fail err=", err)
		return errors.WithCode(code.ErrJsonUnmarshal, "")
	}

	//1先声明一个 resMes
	var resMes model.Message
	resMes.Type = model.RegisterResMesType
	var registerResMes model.RegisterResMes

	//我们需要到redis数据库去完成注册.
	//1.使用model.MyUserDao 到redis去验证
	err = model.MyUserDao.Register(&registerMes.User)

	if err != nil {
		coder := errors.ParseCoder(err)
		registerResMes.Code = coder.Code()
		registerResMes.Error = coder.String()
	} else {
		registerResMes.Code = 200
	}

	data, err := json.Marshal(registerResMes)
	if err != nil {
		log.Debugf("json.Marshal fail", err)
		return errors.WithCode(code.ErrJsonUnmarshal, "")
	}

	//4. 将data 赋值给 resMes
	resMes.Data = string(data)

	//5. 对resMes 进行序列化，准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		log.Debugf("json.Marshal fail", err)
		return errors.WithCode(code.ErrJsonUnmarshal, "")
	}
	//6. 发送data, 我们将其封装到writePkg函数
	//因为使用分层模式(mvc), 我们先创建一个Transfer 实例，然后读取
	tf := &transfer.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	log.Debugf("alive:register end")
	return

}

//编写一个函数serverProcessLogin函数， 专门处理登录请求
func (this *UserProcess) ServerProcessLogin(mes *model.Message) (err error) {
	//核心代码...
	//1. 先从mes 中取出 mes.Data ，并直接反序列化成LoginMes
	var loginMes model.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		log.Debugf("json.Unmarshal fail err=", err)
		return errors.WithCode(code.ErrJsonUnmarshal, "")
	}
	//1先声明一个 resMes
	var resMes model.Message
	resMes.Type = model.LoginResMesType
	//2在声明一个 LoginResMes，并完成赋值
	var loginResMes model.LoginResMes

	//我们需要到redis数据库去完成验证.
	//1.使用model.MyUserDao 到redis去验证

	log.Debugf("alive,prepare redis")
	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)
	if err != nil {
		coder := errors.ParseCoder(err)
		loginResMes.Code = coder.Code()
		loginResMes.Error = coder.String()
		log.Debugf("%v %v %v", loginResMes.Code, coder.HTTPStatus(), loginResMes.Error)

	} else {
		loginResMes.Code = 200
		//这里，因为用户登录成功，我们就把该登录成功的用放入到userMgr中
		//将登录成功的用户的userId 赋给 this
		this.UserId = loginMes.UserId
		userMgr.AddOnlineUser(this)
		//通知其它的在线用户， 我上线了
		this.NotifyOthersOnlineUser(loginMes.UserId)
		//将当前在线用户的id 放入到loginResMes.UsersId
		//遍历 userMgr.onlineUsers
		for id, _ := range userMgr.onlineUsers {
			loginResMes.UsersId = append(loginResMes.UsersId, id)
		}
		log.Debugf(user.UserName, "登录成功")
	}

	// //如果用户id= 100， 密码=123456, 认为合法，否则不合法

	// if loginMes.UserId == 1 && loginMes.UserPwd == "1" {
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
		return errors.WithCode(code.ErrJsonUnmarshal, "")
	}

	//4. 将data 赋值给 resMes
	resMes.Data = string(data)

	//5. 对resMes 进行序列化，准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		log.Debugf("json.Marshal fail", err)
		return errors.WithCode(code.ErrJsonUnmarshal, "")
	}
	//6. 发送data, 我们将其封装到writePkg函数
	//因为使用分层模式(mvc), 我们先创建一个Transfer 实例，然后读取
	tf := &transfer.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	log.Debugf("alive end")
	return
}
