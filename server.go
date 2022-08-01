package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

func (s *Server) Hander(conn net.Conn) {
	//当前连接的业务
	fmt.Println("连接建立成功")
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
		Ip:   ip,
		Port: port,
	}

	return server
}
