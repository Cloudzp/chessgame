package cg

import (
	"encoding/json"
	"errors"
	"game/src/ipc"
	"sync"
)

var _ ipc.Server = &CenterServer{}

type Message struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
}

type CenterServer struct {
	servers map[string]ipc.IpcServer
	player  []*Player
	mutex   sync.Mutex
}

func NewCenterServer() *CenterServer {
	servers := make(map[string]ipc.IpcServer)
	player := make([]*Player, 0)
	return &CenterServer{servers: servers, player: player}
}

func (server *CenterServer) addPlayer(params string) error {
	player := NewPlayer()
	err := json.Unmarshal([]byte(params), &player)
	if err != nil {
		return err
	}
	server.mutex.Lock()
	defer server.mutex.Unlock()

	server.player = append(server.player, player)
	return nil

}

func (server *CenterServer) removePlayer(params string) error {
	for i, v := range server.player {
		if v.Name == params {
			server.player[i] = nil
			return nil
		}
	}
	return errors.New("Player Not Exist!")
}

func (server *CenterServer) listPlayer() (string, error) {
	if len(server.player) > 0 {
		ps, err := json.Marshal(server.player)
		if err != nil {
			return "", err
		}
		return string(ps), nil
	}

	return "", errors.New("No player online!")
}

func (server *CenterServer) broadcast(params string) error {
	var massage Message
	json.Unmarshal([]byte(params), &massage)
	if len(server.player) > 0 {
		for _, player := range server.player {
			player.mq <- &massage
		}
	} else {
		return errors.New("No player online!")
	}
	return nil

}

func (server *CenterServer) Name() string {
	return "CenterServer"
}
func (server *CenterServer) Handle(method, params string) *ipc.Response {
	switch method {
	case "addPlayer":
		 server.addPlayer(params)
	case "removePlayer":
		server.removePlayer(params)
	case "listPlayer":
		server.listPlayer()
	case "broadcast":
		server.broadcast(params)
	default:
		return &ipc.Response{"400", method + ":" + params}
	}
	return &ipc.Response{"200",""}
}
