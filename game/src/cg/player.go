package cg

import "fmt"

type Player struct {
	Name string
	Level int
	Exp int
	Room int

	mq chan *Message
}
func NewPlayer()*Player{
	m := make(chan *Message,1024)
	player := &Player{mq:m}

	go func (p *Player){
		for {
			msg := <-p.mq
			fmt.Println(p.Name,"recived message:",msg)
		}
	}(player)
	return player
}