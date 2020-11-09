package controller

import (
	"fmt"
	"log"
	"net/http"
	models "whatsapp_api/model"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"golang.org/x/net/websocket"
)

type ServerWsConnection struct {
	Clients           map[string]*Client
	CommunicationChan chan *Client
	DeleteChan        chan *Client
	Termination       chan bool
	Dbconn            *gorm.DB
	Count             int
}

func NewWsServer(dbconn *gorm.DB) *ServerWsConnection {
	clients := make(map[string]*Client)
	communicationChan := make(chan *Client)
	deleteChan := make(chan *Client)
	termination := make(chan bool)
	count := 0
	return &ServerWsConnection{
		clients,
		communicationChan,
		deleteChan,
		termination,
		dbconn,
		count,
	}
}

func (s *ServerWsConnection) Del(u *Client) {
	s.DeleteChan <- u
}

func (s *ServerWsConnection) NewWsClient(msg *models.WsResponse, ws *websocket.Conn) *Client {
	nc := new(Client)
	nc.Id = msg.WsMessage.UserId
	nc.UType = msg.WsMessage.UserType
	nc.Ws = ws
	nc.Ch = make(chan *models.WsResponse, 100)
	nc.Termination = make(chan bool)
	nc.Server = s
	nc.Syncsend = false
	s.CommunicationChan <- nc
	return nc
}

func (s *ServerWsConnection) userLogin(msg *models.WsResponse, ws *websocket.Conn) {
	u := s.NewWsClient(msg, ws)
	ReWriteResponse(msg, "1", http.StatusOK, "Session Started.")
	if err := websocket.JSON.Send(u.Ws, msg); err != nil {
		log.Println("Can't send", err)
	}
	u.Listen()
}
func ReWriteResponse(msg *models.WsResponse, msgId string, rspCode int, rspMsg string) {
	msg.MessageId = msgId
	msg.ResponseCode = rspCode
	msg.ResponseMessage = rspMsg
}
func (s *ServerWsConnection) VerifyClientConnection(msg *models.WsResponse) bool {
	validUserId := msg.WsMessage.UserId != ""
	if validUserId {
		for _, oldu := range s.Clients {
			if oldu.Id == msg.WsMessage.UserId {
				ReWriteResponse(msg, "2", http.StatusBadRequest, "user exists")
				return false
			}
		}
	} else {
		ReWriteResponse(msg, "0", http.StatusBadRequest, "invalid user")
		return false
	}
	return true
}

func (s *ServerWsConnection) WsRegister(c echo.Context) error {

	websocket.Handler(func(ws *websocket.Conn) {
		var msg models.WsResponse
		for {
			if err := websocket.JSON.Send(ws, "register"); err != nil {
				log.Println("Can't send", err)
			}
			if err := websocket.JSON.Receive(ws, &msg.WsMessage); err != nil {
				log.Printf("Can't receive %v ", err.Error)
				break
			}
			fmt.Printf("recieved %v \n", msg)
			if s.VerifyClientConnection(&msg) {
				s.userLogin(&msg, ws)
			} else {
				break
			}
		}
		if err := websocket.JSON.Send(ws, msg); err != nil {
			log.Println("Can't send", err)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func (s *ServerWsConnection) updateAgentStatus(id string, status string) {
	s.Dbconn.Table("v_call_center_agents").Where("call_center_agent_uuid=?", id).Update("agent_status", status)
	s.Dbconn.Table("agents").Where("name=?", id).Update("status", status)
}

func (s *ServerWsConnection) Controller(e *echo.Echo) {

	e.GET("/main", s.WsRegister)
	for {
		select {
		case u := <-s.CommunicationChan:
			s.Clients[u.Id] = u
			s.Count++
			log.Println("Now", s.Count, u.Id, "+1 users connected.")

		case u := <-s.DeleteChan:
			delete(s.Clients, u.Id)
			s.updateAgentStatus(u.Id, "Logged Out")
			s.Count--
			log.Println("Now", s.Count, u.Id, "-1 users connected.")
		}
	}
}
