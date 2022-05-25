package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	code       int // 当前模式
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		code:       999,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return nil
	}

	client.conn = conn

	return client
}

func (this *Client) Run() {
	for this.code != 0 {
		for this.menu() != true {
		}

		switch this.code {
		case Rename:
			this.Rename()
			break
		case PrivateChat:
			fmt.Println("私聊模式")
			break
		case PublicChat:
			fmt.Println("公聊模式")
			break
		}
	}
}

// 请求菜单
func (this *Client) menu() bool {
	var code int

	fmt.Println("1.更改用户名")
	fmt.Println("2.私聊模式")
	fmt.Println("3.公聊模式")
	fmt.Println("0.退出")

	fmt.Scanln(&code)

	if code >= 0 && code <= 3 {
		this.code = code
		return true
	} else {
		fmt.Println(">>>请输入合法范围内的数字...")
		return false
	}

}

// Rename 修改用户名
func (this *Client) Rename() {
	fmt.Println(">>>请输入用户名...")

	var newName string
	fmt.Scanln(&newName)

	sendMsg := "rename|" + newName + "\n"
	_, err := this.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return
	}

	// 让子弹飞，等待 server 的响应结果
	time.Sleep(time.Microsecond * 100)

}

// 处理 server 回应的数据，直接显示到便准输出即可
func (this *Client) DealResponse() {
	// 一旦 client.conn 有数据，就直接 copy 到 stdout 标准输出上
	// 永久阻塞监听
	io.Copy(os.Stdout, this.conn)

	//上下两种写法的效果等价
	//for {
	//	buf := make([]byte, 4096)
	//	client.conn.Read(buf)
	//	fmt.Println(buf)
	//}
}

func main() {
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println(">>>服务器连接失败...")
		return
	}

	// 单独开启一个 goroutine 去处理 server 的回执消息
	go client.DealResponse()
	client.Run()
}
