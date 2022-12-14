package main

import (
	"fmt"
	"net"
	"strings"
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
	//添加用户到map
	u.Server.OnlienMap[u.Name] = u
	fmt.Println(u.Name + " Online !")
	u.Server.maplock.Unlock() //解锁 关于安全map的说明 https://zhuanlan.zhihu.com/p/449078860

	//广播上线消息
	u.Server.BroadCast(u, "online")
}

//下线功能
func (u *User) Offline() {

	u.Server.maplock.Lock() //先锁 不太明白 弹幕说这里的map不是线程安全的 不过查到好像加锁就是安全的了
	//删除用户
	delete(u.Server.OnlienMap, u.Name)
	fmt.Println(u.Name + " Offline !")
	u.Server.maplock.Unlock() //解锁 关于安全map的说明 https://zhuanlan.zhihu.com/p/449078860

	//广播下线消息
	u.Server.BroadCast(u, "offline")
}

//处理消息
func (u *User) DoMessage(msg string) {

	//如果用户输入的是who 那么查询当前在线用户
	if msg == "who" {
		u.Server.maplock.Lock()
		for _, user := range u.Server.OnlienMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ": now online\n"
			u.C <- onlineMsg
		}
		u.Server.maplock.Unlock()

		//如果长度大于7 并且前面7个字符为rename
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		//取到要改的名称
		newName := strings.Split(msg, "|")[1] //split方法 和 js类似
		//判断是否被占用
		_, ok := u.Server.OnlienMap[newName] //拿不拿得到？我去 这样循环都不用写
		if ok {
			u.Server.BroadCast(u, "username is occupying")
		} else {
			u.Server.maplock.Lock()
			// u.Server.OnlienMap[u.Name].Name = newName //不能这样？ 他们说map的key改不了 好像这里的key是用名字的 怪不得不能这样
			delete(u.Server.OnlienMap, u.Name)
			u.Server.OnlienMap[newName] = u //这里先赋值获取 但是当前u姓名还没改 以后改也没事 因为是传指针
			u.Server.maplock.Unlock()

			u.Name = newName
			u.Server.BroadCast(u, "update name over!")
		}

		//私聊消息
	} else if len(msg) > 4 && msg[:3] == "to|" {
		//消息格式 to|张三|你好xxxxx
		//1.获取对方用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			u.C <- "import username please like to|username|message \n"
			return
		}
		//2.根据用户名得到对方的user对象
		sideUser, ok := u.Server.OnlienMap[remoteName]
		if !ok {
			u.C <- "user is not undefined"
			return
		}
		//3.获取消息内容 通过对方的user对象 发送内容
		// if msg = strings.Split(msg, "|")[2]; msg == "" {
		// 	u.C <- "you not import message \n"
		// 	return
		// }
		context := strings.Split(msg, "|")[2] //上面注释掉的写法不行 应该是因为if不能那样写
		if context == "" {
			u.C <- "you not import message \n"
			return
		}
		message := u.Name + ":to you message:" + strings.Split(msg, "|")[2] + "\n"
		sideUser.C <- message
	} else {
		//输入的不是w或者改名就正常发布消息
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
