package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int
	//在线用户的列表 key为字符串 value是User对象 但是这里为什么传地址？可能是值传递不拷贝
	OnlienMap map[string]*User
	//一个锁 不太理解
	maplock sync.RWMutex

	//消息广播的channel
	Message chan string
}

//广播上线消息
func (s *Server) BroadCast(user *User, msg string) {
	//定义消息
	sendMsg := "[" + user.Addr + "]" + user.Name + msg + "\n"
	// 发送到广播管道中
	s.Message <- sendMsg
}

//监听当前消息channel的goroutine
func (s *Server) ListenMessager() {
	for {
		//不断的从message中读数据 这个会阻塞 当能读到数据后
		msg := <-s.Message

		//发送给全部的在线user 这里也要加锁 不太懂 弹幕说sync.RWMutex这个可以多用户读？或者不加锁？或者用Rlock
		s.maplock.Lock()
		for _, v := range s.OnlienMap {
			//穿给用户的channel
			v.C <- msg
		}
		s.maplock.Unlock()
	}
}

func (s *Server) Hander(conn net.Conn) {
	//当前连接的业务
	// fmt.Println("连接建立成功")

	//用户上线了 将用户加入到OnlineMap中

	//1.首先创建user
	user := NewUser(conn)

	s.maplock.Lock() //先锁 不太明白 弹幕说这里的map不是线程安全的 不过查到好像加锁就是安全的了
	//2.添加用户到map
	s.OnlienMap[user.Name] = user
	s.maplock.Unlock() //解锁 关于安全map的说明 https://zhuanlan.zhihu.com/p/449078860

	//广播上线消息
	s.BroadCast(user, "online")

	//当前handler阻塞 但是这里好像不写也没事 日后再看
	select {}
}

//启动服务器的接口
func (s *Server) Start() {

	//socket listen  //这里net.Listenr第二个参数需要拼接 因为要自定义
	Listenr, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("服务器启动失败", err)
		return
	}

	//close listen socket
	defer Listenr.Close()

	// 启动监听Message的go程
	go s.ListenMessager()

	for {
		//accept 监听请求 //这个会阻塞
		conn, err := Listenr.Accept()
		if err != nil {
			fmt.Println("listenr accept err:", err)
			//连接错误就跳出循环 开启下一次循环
			continue
		}

		//do handler //这个开一个goroutine 然后for循环会回去再监听请求 有连接的话就开一个携程去执行
		go s.Hander(conn)
	}

}

//创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlienMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}
