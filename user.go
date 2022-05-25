package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

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

// ListenMessage 接收 server 发来的消息
func (this *User) ListenMessage() {
	for {
		msg := <-this.Message
		this.PrintMessage(msg)
	}
}

// PrintMessage 终端打印消息
func (this *User) PrintMessage(msg string) {
	conn := this.conn
	conn.Write([]byte(msg + "\n"))
}

// ListenWrite 监听用户输入
func (this *User) ListenWrite(server *Server) {
	buf := make([]byte, 4096)

	for {
		n, err := this.conn.Read(buf)

		// 用户下线，不再发送消息
		if n == 0 {
			this.Offline(server)
			return
		}

		if err != nil && err != io.EOF {
			fmt.Println("Conn Read err:", err)
			return
		}

		// 获取用户输入（去掉'\n'）
		msg := NewMsg(string(buf[:n-1]))
		switch msg.code {
		case Rename:
			str := strings.Split(msg.str, "|")[1]
			newName := strings.TrimSpace(str)
			this.Rename(server, newName)
			break
		case PrivateChat:
			break
		case PublicChat:
			break
		case 0:
			this.PrintMessage("[server]: 我不理解")
			break
		}
	}
}

// Rename 修改用户名
func (this *User) Rename(server *Server, newName string) {
	if newName == "" {
		this.PrintMessage("[修改失败]: 用户名不能为空")
		return
	}

	_, ok := server.UserMap[newName]
	if ok {
		this.PrintMessage("[修改失败]: 当前用户名已存在 ")
		return
	}

	oldName := this.Name

	server.mapLock.Lock()
	delete(server.UserMap, this.Name)
	server.UserMap[newName] = this
	server.mapLock.Unlock()

	this.Name = newName
	fmt.Println("[", oldName, "] rename to", "[", newName, "]")
	this.PrintMessage("[修改成功]: " + newName)
}

// Online 用户上线
// - 用户加入 server.UserMap
// - 对所有用户进行广播提示
func (this *User) Online(server *Server) {
	server.mapLock.Lock()
	server.UserMap[this.Name] = this
	server.mapLock.Unlock()

	fmt.Println("[", this.Name, "]", "已上线")
	server.BroadCast(this, "已上线")
}

// Offline 用户下线
// - 用户移出 server.UserMap
// - 对所有用户进行广播提示
func (this *User) Offline(server *Server) {
	server.mapLock.Lock()
	delete(server.UserMap, this.Name)
	server.mapLock.Unlock()

	fmt.Println("[", this.Name, "]", "已下线")
	server.BroadCast(this, "已下线")
}
