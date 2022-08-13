package main

import (
	"flag"
	"fmt"
	"net"
)

//建立连接类
type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	//连接句柄 conn是i一个面向流的网络连接
	conn net.Conn
	//客户端输入
	flag int
}

//菜单
func (client *Client) menu() bool {
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")
	//接收输入
	fmt.Scanln(&flag)

	//合法
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else { //非法
		fmt.Println(">>>inport err<<<")
		return false
	}

}

func (client *Client) Run() {

	//判断是否是0
	for client.flag != 0 {

		for client.menu() != true { //这里之前写错了 != 写成了 == 然后就退出不了了 ps:不知道为啥

		}

		switch client.flag {
		case 1:
			//公聊
		case 2:
			//私聊
		case 3:
			//更新用户名

		}
	}
	fmt.Println("结束", client.flag)
}

//构造函数
func NewClient(serverIp string, serverPort int) *Client {
	//创建客户端对象 名字和conn都是之后在赋值
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999, //默认999
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

var serverIp string

var serverPort int

//解析命令行 所以解析要在mian之前 so 写一个init函数
func init() {
	// 会解析-ip之后的数据 ./client -ip 127.0.0.1
	// 指针 解析前缀 默认值 help的说明 如：go run ./client.go -h
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "set server ip")
	flag.IntVar(&serverPort, "port", 8000, "set server port")
}

func main() {

	//命令行解析 解析命令行数据
	flag.Parse()

	//构造连接
	client := NewClient(serverIp, serverPort)

	if client == nil {
		fmt.Println(">>>>>>>client server err......")
		return
	}

	fmt.Println(">>>>>>>client server success......")

	//启动客户端的业务
	client.Run()
}
