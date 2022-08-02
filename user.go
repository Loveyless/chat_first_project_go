package main

import (
	"net"
)

type User struct {
	// 名字
	Name string
	// 地址
	Addr string
	// 当前用户的channel
	C chan string
	// 当前用户的连接
	Conn net.Conn
	//当前用户关联的Server
	Server *Server
}

//上线功能
func (u *User) Online() {
	u.Server.maplock.Lock() //先锁 不太明白 弹幕说这里的map不是线程安全的 不过查到好像加锁就是安全的了
	//2.添加用户到map
	u.Server.OnlienMap[u.Name] = u
	u.Server.maplock.Unlock() //解锁 关于安全map的说明 https://zhuanlan.zhihu.com/p/449078860

	//广播上线消息
	u.Server.BroadCast(u, "online")
}

//下线功能
func (u *User) Offline() {

	u.Server.maplock.Lock() //先锁 不太明白 弹幕说这里的map不是线程安全的 不过查到好像加锁就是安全的了
	//2.添加用户到map
	delete(u.Server.OnlienMap, u.Name)
	u.Server.maplock.Unlock() //解锁 关于安全map的说明 https://zhuanlan.zhihu.com/p/449078860

	//广播下线消息
	u.Server.BroadCast(u, "offline")
}

//处理消息
func (u *User) DoMessage(msg string) {
	//如果用户输入的是w 那么查询当前在线用户
	if msg == "who" {
		u.Server.maplock.Lock()
		for _, user := range u.Server.OnlienMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ": now online\n"
			u.C <- onlineMsg
		}
		u.Server.maplock.Unlock()
	} else {
		//输入的不是w就正常发布消息
		u.Server.BroadCast(u, msg)
	}
}

// 监听当前User channel方法 一旦有消息 就直接发送给对端客户端
func (u *User) ListenMessage() {
	for {
		// 不断的从用户的管道中读数据
		msg := <-u.C
		// 一旦有消息 写入数据 这里用byte切片不太理解 怎么不用string 我懂了 这个函数直接接byte 这里中文cmd会乱码
		u.Conn.Write([]byte(msg))
	}
}

//创建一个用户的api
func NewUser(conn net.Conn, server *Server) *User {

	// 默认用户名是客户端地址
	// 拿到地址
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		Conn:   conn,
		Server: server,
	}

	// 开线程去监听user channel消息
	go user.ListenMessage()

	return user
}
