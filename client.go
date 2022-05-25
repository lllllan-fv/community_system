package main

import (
	"fmt"
	"io"
	"net"
	"os"
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

func (client *Client) Run() {
	for client.code != 0 {
	}
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
