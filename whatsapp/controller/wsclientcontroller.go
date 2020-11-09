package controller

import (
	"io"
	"log"
	"net/http"
	models "whatsapp_api/model"

	"golang.org/x/net/websocket"
	// "fmt"
)

type Client struct {
	Id          string
	UType       string
	Ws          *websocket.Conn
	Ch          chan *models.WsResponse
	Termination chan bool
	Server      *ServerWsConnection
	Syncsend    bool
}

func (u *Client) DeleteClient(msg *models.WsResponse) {
	ReWriteResponse(msg, "3", http.StatusOK, "log out")
	u.Ch <- msg
	u.Termination <- true
}

func (u *Client) Listen() {
	go u.listenWrite()
	u.listenRead()
}

func (u *Client) listenWrite() {
	for {
		select {
		case msg := <-u.Ch:
			websocket.JSON.Send(u.Ws, msg)
		case <-u.Termination:
			u.Server.Del(u)
			return
		}
	}
}

func (u *Client) listenRead() {
	// for {
	select {
	default:
		var msg models.WsResponse
		err := websocket.JSON.Receive(u.Ws, &msg.WsMessage)
		log.Println(msg)
		if err == io.EOF {
			log.Println("EOF err: ", err)
			u.DeleteClient(&msg)
		} else if err != nil {
			log.Println("err: ", err)
		} else {
			if msg.WsMessage.MessageId == "3" {
				// log.Println("log out user")
				u.DeleteClient(&msg)
			}
		}
	}
	// }
}
