package main

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
	if len(str) > 7 && str[:7] == "rename|" {
		return Rename
	}

	if len(str) > 5 && str[:3] == "to|" {
		return PrivateChat
	}

	if len(str) > 4 && str[:4] == "pub|" {
		return PublicChat
	}

	return 0
}
