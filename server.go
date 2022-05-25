package main

import (
	"fmt"
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

func (this *Server) Handler(conn net.Conn) {
	user := NewUser(conn)

	// 在开启 server 的终端中打印信息
	fmt.Println("[" + user.Name + "]建立连接...")
	defer fmt.Println("[" + user.Name + "]关闭连接...")

	// 当前连接的终端打印信息
	conn.Write([]byte("建立连接...\n"))
	defer conn.Close()
	defer conn.Write([]byte("关闭连接...\n"))
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

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("listener accept err:", err)
			return
		}

		go this.Handler(conn)
	}
}
