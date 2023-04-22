package apiserver

import "net"

var (
	userMgr *UserMgr
)

type UserProcess struct {
	//字段
	Conn net.Conn
	//增加一个字段，表示该Conn是哪个用户
	UserId int
}
type UserMgr struct {
	onlineUsers map[int]*UserProcess
}
