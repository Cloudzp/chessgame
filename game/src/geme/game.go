package main

import (
	"fmt"
	"game/src/cg"
	"game/src/ipc"
	"strconv"
	"strings"
	"bufio"
	"os"
)

var centerClient *cg.CenterClient

func startCenterService() error {

	server := ipc.NewIpcServer(&cg.CenterServer{})
	client := ipc.NewIpcClient(server)
	centerClient = &cg.CenterClient{client}

	return nil
}

func Help(args []string) int {
	fmt.Println(`
	Commands:
	     login <username>
	     send <message>
	     listplayer
	     quit(q)
	     help(h)
	`)
	return 0
}

func Quit(args []string) int {
	return 1
}

func Logout(args []string) int {
	if len(args) != 2 {
		fmt.Println("USAGE: logout <username>")
		return 0
	}
	centerClient.RemovePlayer(args[1])
	return 0
}

func Login(args []string) int {
	if len(args) != 4 {
		fmt.Println("USAGE:login <username><level><exp>")
		return 0
	}

	level, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("Invalid Parameter:<level>")
		return 0
	}
	exp, err := strconv.Atoi(args[3])
	if err != nil {
		fmt.Println("Invalid Parameter:<exp>")
		return 0
	}

	player := cg.NewPlayer()
	player.Name = args[1]
	player.Exp = exp
	player.Level = level

	err = centerClient.AddPlayer(player)
	if err != nil {
		fmt.Println("Failed adding player", err)
	}

	return 0

}
func ListPlayer(args []string) int {

	ps, err := centerClient.ListPlayer("")
	if err != nil {
		fmt.Println("Failed.", err)
	} else {
		for i, v := range ps {
			fmt.Println(i+1, ":", v)
		}
	}
	return 0
}

func Send(args []string) int {
	message := strings.Join(args[1:], "")

	err := centerClient.Broadcast(message)

	if err != nil {
		fmt.Println("Failed.", err)
	}

	return 0
}

func GetCommandHandlers() map[string]func(args []string) int {
	return map[string]func([]string) int{
		"help":       Help,
		"h":          Help,
		"quit":       Quit,
		"q":          Quit,
		"login":      Login,
		"logout":     Logout,
		"listplayer": ListPlayer,
		"send":       Send,
	}
}

func main(){
	fmt.Println("Casual Game Server Solution")

	Help(nil)

	r := bufio.NewReader(os.Stdin)

	handlers := GetCommandHandlers()

	for {
		fmt.Print("Command>>>")
		b,_,_ := r.ReadLine()
		line := string(b)

		token := strings.Split(line," ")

		if handler, ok := handlers[token[0]];ok {
			ret := handler(token)
			if ret != 0 {
				break
			}
		}else {
			fmt.Println("Unknown command:",token[0])
		}
	}
}