package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

func newServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}

	return server
}

func (this *Server) Handler(conn net.Conn) {
	// 在开启 server 的终端中打印信息
	fmt.Println("建立连接...")
	defer fmt.Println("关闭连接...")

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
