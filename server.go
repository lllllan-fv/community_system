package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	// 在线用户列表
	UserMap map[string]*User
	mapLock sync.RWMutex

	// 广播消息
	Message chan string
}

func newServer(ip string, port int) *Server {
	server := &Server{
		Ip:      ip,
		Port:    port,
		UserMap: make(map[string]*User),
		Message: make(chan string),
	}

	return server
}

// Handler 处理每一个连接请求
func (this *Server) Handler(conn net.Conn) {
	user := NewUser(conn)

	// 当前连接的终端打印信息
	conn.Write([]byte("建立连接...\n"))

	// 上线加入 UserMap
	user.Online(this)

	// server 监听该用户的输入
	go this.ListenUserWrite(user)
}

// BroadCast
// - server 对所有用户进行广播消息
// - user 是发送消息的用户，可以为 nil
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.mapLock.Lock()
	for _, toUser := range this.UserMap {
		if user != toUser {
			toUser.Message <- sendMsg
		}
	}
	this.mapLock.Unlock()
}

// ListenUserWrite 监听用户的输入
func (this *Server) ListenUserWrite(user *User) {
	conn := user.conn
	buf := make([]byte, 4096)

	for {
		n, err := conn.Read(buf)

		// 用户下线，不再发送消息
		if n == 0 {
			user.Offline(this)
			return
		}

		if err != nil && err != io.EOF {
			fmt.Println("Conn Read err:", err)
			return
		}

		// 获取用户输入（去掉'\n'）
		msg := string(buf[:n-1])
		this.BroadCast(user, msg)
	}
}

func (this *Server) start() {
	fmt.Println("Server start...")
	defer fmt.Println("Server over...")

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	defer listener.Close()

	if err != nil {
		fmt.Println("net.listener err:", err)
		return
	}

	// 监听端口连接请求
	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("listener accept err:", err)
			return
		}

		// 处理每一个连接请求
		go this.Handler(conn)
	}
}
