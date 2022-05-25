package main

import "net"

type User struct {
	Name string
	Addr string
	conn net.Conn

	// 消息接收
	message chan string
}

func NewUser(conn net.Conn) *User {
	addr := conn.RemoteAddr().String()

	user := &User{
		Name:    addr,
		Addr:    addr,
		conn:    conn,
		message: make(chan string),
	}

	return user
}
