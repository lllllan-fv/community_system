package main

import "net"

type User struct {
	Name string
	Addr string
	conn net.Conn

	// 消息接收
	Message chan string
}

func NewUser(conn net.Conn) *User {
	addr := conn.RemoteAddr().String()

	user := &User{
		Name:    addr,
		Addr:    addr,
		conn:    conn,
		Message: make(chan string),
	}

	go user.ListenMessage()

	return user
}

func (this *User) ListenMessage() {
	for {
		msg := <-this.Message
		this.PrintMessage(msg)
	}
}

func (this *User) PrintMessage(msg string) {
	conn := this.conn
	conn.Write([]byte(msg + "\n"))
}

func (this *User) Online(server *Server) {
	server.mapLock.Lock()
	server.UserMap[this.Addr] = this
	server.mapLock.Unlock()

	server.BroadCast(this, "已上线")
}

func (this *User) Offline(server *Server) {
	server.mapLock.Lock()
	delete(server.UserMap, this.Addr)
	server.mapLock.Unlock()

	server.BroadCast(this, "已下线")
}
