package main

import "net"

type User struct {
	// 名字
	Name string
	// 地址
	Addr string
	// 当前用户的channel
	C chan string
	// 当前用户的连接
	Conn net.Conn
}

// 监听当前User channel方法 一旦有消息 就直接发送给对端客户端
func (u *User) ListenMessage() {
	for {
		// 不断的从用户的管道中读数据
		msg := <-u.C
		// 一旦有消息 写入数据 这里用byte切片不太理解 怎么不用string 我懂了 这个函数直接接byte 这里中文cmd会乱码 操
		u.Conn.Write([]byte("have new msg" + msg))
	}
}

//创建一个用户的api
func NewUser(conn net.Conn) *User {

	// 默认用户名是客户端地址
	// 拿到地址
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		Conn: conn,
	}

	// 开线程去监听user channel消息
	go user.ListenMessage()

	return user
}
