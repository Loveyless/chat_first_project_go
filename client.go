package main

import (
	"fmt"
	"net"
)

//建立连接
type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	//连接句柄 conn是i一个面向流的网络连接
	conn net.Conn
}

//构造函数
func NewClient(serverIp string, serverPort int) *Client {
	//创建客户端对象 名字和conn都是之后在赋值
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}
	//尝试连接服务器 连接server Dial连接指定网络上的地址
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial err", err)
		return nil
	}
	client.conn = conn
	//返回对象

	return client
}

func main() {

	client := NewClient("127.0.0.1", 8000)

	if client == nil {
		fmt.Println(">>>>>>>client server err......")
		return
	}

	fmt.Println(">>>>>>>client server success......")

	//启动客户端的业务
	select {}
}
