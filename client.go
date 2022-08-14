package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
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
			client.PublicChat()
		case 2:
			//私聊
			client.PrivateChat()
		case 3:
			//更新用户名
			client.UpdateName()
		}
	}
	fmt.Println("结束", client.flag)
}

//处理server回应的消息 直接显示到标准输出
func (client *Client) DealResponse() {
	//简写方法 io.Copy会阻塞 等待conn
	//一旦有消息 就直接copy到标准输出上 永久阻塞监听(永久阻塞)
	io.Copy(os.Stdout, client.conn)
	//等价于
	// for {
	// 	buf := make([]byte, 4096)
	// 	client.conn.Read(buf)
	// 	fmt.Println(string(buf))
	// }
}

//查询在线用户
func (client *Client) SelectUsers() {
	msg := "who\n"
	_, err := client.conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("conn write err:", err)
	}
}

//私聊模式 需要先查询在线用户SelectUsers
func (client *Client) PrivateChat() {
	//查询当前在线用户
	client.SelectUsers()
	var remoteName string
	var chatMessage string
	fmt.Println(">>>>please import username , import exit quit")
	fmt.Scanln(&remoteName)

	//如果成功输入姓名
	for remoteName != "exit" {

		fmt.Println("please import message , import exit quit")
		fmt.Scanln(&chatMessage)

		//成功输入内容
		for chatMessage != "exit" {

			//判断是否为空
			if len(chatMessage) != 0 {

				chatMessage = "to|" + remoteName + "|" + chatMessage + "\n"
				_, err := client.conn.Write([]byte(chatMessage))
				if err != nil {
					fmt.Println("write err:", err)
					break
				}

			}

			//发完在提示用户继续输入消息 !!!但是这里exit退出不了 但是公聊可以 因为这个外面套了一个for来选对谁私聊
			chatMessage = ""
			// fmt.Println("please import message , import exit quit")
			fmt.Scanln(&chatMessage)

		}

	}

	//不发了就重新选对象
	remoteName = ""
	fmt.Println(">>>>please import username , import exit quit")
	fmt.Scanln(&remoteName)

}

func (client *Client) PublicChat() {
	//提示用户输入消息
	var msg string
	fmt.Println(">>>>please import message , import exit quit")
	fmt.Scanln(&msg)
	//发给服务器
	for msg != "exit" {

		if len(msg) != 0 {
			msg = msg + "\n"
			_, err := client.conn.Write([]byte(msg))
			if err != nil {
				fmt.Println("write err:", err)
				break
			}
		}

		//因为要一直接收 直到输入exit
		msg = ""
		// fmt.Println(">>>>please import message , import exit quit")
		fmt.Scanln(&msg)

	}
}

func (client *Client) UpdateName() bool {
	fmt.Println(">>> please import New name")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	//write将数据写入连接
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}
	//如果改变成功 服务器还会输出更改成功 应该用一个go程来接收服务器的消息
	return true
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

	//单独开启goroutine 接收消息
	go client.DealResponse()

	fmt.Println(">>>>>>>client server success......")

	//启动客户端的业务
	client.Run()
}
