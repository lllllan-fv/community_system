package main

import "fmt"

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

func (this *Server) start() {
	fmt.Println("Server running...")

	defer fmt.Println("Server over...")
}
