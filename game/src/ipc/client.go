package ipc

import (
	"encoding/json"
	"fmt"
)

type IpcClient struct {
	conn chan string
}
func NewIpcClient(server *IpcServer)*IpcClient{
	c := server.Connect()
	return &IpcClient{c}
}
func (clinent *IpcClient)Call(method,params string) (resp *Response,err error){
	req := &Request{method,params}

	response,err := json.Marshal(req)
	if err != nil {
		fmt.Println("Marshal Request Error!")
		return
	}
	clinent.conn <- string(response)
	str := <-clinent.conn
	err = json.Unmarshal([]byte(str),&resp)
	if err != nil {
		fmt.Println("Handle Error!")
	}
	return
}
func (clinent *IpcClient)Close(){
	clinent.conn <- "CLOSE"
}