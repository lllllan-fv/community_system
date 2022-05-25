package main

import "strings"

type Msg struct {
	str  string
	code int
}

const (
	Rename = iota + 1
	PrivateChat
	PublicChat
)

func NewMsg(str string) *Msg {
	msg := &Msg{
		str:  str,
		code: calCode(str),
	}

	return msg
}

/*
	1 rename| 改名
	2 to|toUser| 私聊
	3 pub| 公聊
	0 <else> 其他格式暂时不管
*/
func calCode(str string) int {
	n := len(strings.Split(str, "|"))

	if len(str) > 4 && str[:4] == "pub|" && n == 2 {
		return PublicChat
	}

	if len(str) > 5 && str[:3] == "to|" && n == 3 {
		return PrivateChat
	}

	if len(str) > 7 && str[:7] == "rename|" && n == 2 {
		return Rename
	}

	return 0
}
