package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
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

//广播消息
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

//业务
func (s *Server) Hander(conn net.Conn) {

	//用户上线了 将用户加入到OnlineMap中

	//1.首先创建user
	user := NewUser(conn, s) //这里传入当前Server的地址
	//上线 把user传入OnlineMap
	user.Online()

	// 监听用户是否活跃的channel
	isLive := make(chan bool)

	//接收用户传递的消息
	go func() {

		buf := make([]byte, 4096)

		for {
			//允许从当前连接读消息 成功返回字节数 失败err
			n, err := conn.Read(buf)
			//如果读取的是0那么说明用户关闭了 也就是下线
			if n == 0 {
				user.Offline()
				return
			}
			//发生错误  这个io.EOF不太懂 表示文件末尾？
			if err != nil && err != io.EOF {
				fmt.Println(user, err)
				return
			}

			//提取用户的消息 去除("\n") 因为BroadCast会加\n 这里n是否和len(buf)一个意思呢？
			// msg := string(buf[:n-2]) //这里-2才是取消换行符 //这win比较特殊 很麻烦
			//装了nc好像是len-1就能取消换行 目前不太懂这个n是什么和len有什么区别 nc好像就是测socket的流的 我写的好像就是socket
			msg := string(buf[:n-1])

			//处理+广播消息
			user.DoMessage(msg)

			//发送消息后 向用户活跃channel发送一个消息
			isLive <- true

		}

	}()

	//当前handler阻塞 但是这里好像不写也没事 日后再看
	for {
		select {
		case <-isLive:
			//说明当前用户活跃
			//这里有个技巧 如果isLive执行 那么下面的<-time.After(time.Second * 10):也会执行
			//time.After()会重置定时器 活跃就重置 不一定能进去 但是会执行case的语句
		case <-time.After(time.Second * 10): //go中的定时器

			//已经超时 将当前User强制关闭
			user.C <- "You have been kicked"

			//关闭管道 销毁资源 关闭之前等一会 要不然用户可能收不到消息
			time.Sleep(time.Second * 1)
			close(user.C)

			//列表中删除用户 但是前面读取的时候 如果读的消息是0那么自动下线 所以下面关闭连接后会自动踢
			s.maplock.Lock()
			delete(s.OnlienMap, user.Name)
			s.maplock.Unlock()

			//关闭连接
			conn.Close()

			//退出当前Handler
			return //或者runtime.Goexit() //这个函数没学
		}
	}
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
		//accept 监听用户连接请求 //这个会阻塞
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
