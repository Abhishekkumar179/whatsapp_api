package controller

import (
	"io"
	"log"
	models "whatsapp_api/model"

	"golang.org/x/net/websocket"
	// "fmt"
)

type User struct {
	Id       int
	UType    string
	UName    string
	Ws       *websocket.Conn
	Ch       chan *models.WsMessage
	DoneCh   chan bool
	Server   *ServerUserList
	Syncsend bool
}

func (u *User) Listen() {
	go u.listenWrite()
	u.listenRead()
}
func (u *User) WriteData(msg *models.WsMessage) {
	select {
	case u.Ch <- msg:
	default:
		log.Printf("client %v is disconnected.", u.UName)
	}
}
func (u *User) sendDatabaseMessages(resultstoremsg []*models.WsMessage) {

	for i := 0; i < len(resultstoremsg); i++ {
		// log.Println(resultstoremsg[i].Body)
		for u.Syncsend {
		}
		u.Syncsend = true
		if resultstoremsg[i].File != nil {
			u.WriteFile(resultstoremsg[i])
		} else {
			select {
			case u.Ch <- resultstoremsg[i]:
				// log.Printf("database messages %v",value)
			default:
				log.Printf("client %v is disconnected.", u.UName)
			}
		}
	}
	// for _,value:=range resultstoremsg{
	// }
}
func (u *User) sendNewDatabaseMessages(resultstoremsg []*models.WsMessage) {

	for i := len(resultstoremsg) - 1; i >= 0; i-- {
		// log.Println(resultstoremsg[i].Body)
		for u.Syncsend {
		}
		u.Syncsend = true
		if resultstoremsg[i].File != nil {
			u.WriteFile(resultstoremsg[i])
		} else {
			select {
			case u.Ch <- resultstoremsg[i]:
				// log.Printf("database messages %v",value)
			default:
				log.Printf("client %v is disconnected.", u.UName)
			}
		}
	}
	// for _,value:=range resultstoremsg{
	// }
}
func (u *User) WriteMessageReciept(msg *models.WsMessage) {
	// updateMsgState:=models.WsMessage{MessageType:msg.MessageType,To:msg.To,From:msg.From,Body:msg.Body,File:msg.File,MessageId:msg.MessageId,MessageStageId:msg.MessageStageId}
	// log.Println("writing reciept into from user")

	updateMsgState := models.WsMessage{MessageType: msg.MessageType, To: msg.To, From: msg.From, Body: msg.Body, ReplyMsg: msg.ReplyMsg, File: msg.File, MessageId: msg.MessageId, MessageStageId: msg.MessageStageId}
	updateMsgState.MessageType = "0"
	if updateMsgState.MessageStageId == "0" {
		updateMsgState.MessageStageId = "1"
	}
	// else if updateMsgState.MessageStageId=="1"{
	//         updateMsgState.MessageStageId="2"
	// }else if updateMsgState.MessageStageId=="2"{
	//     updateMsgState.MessageStageId="3"
	//     log.Println("unknown stage ",updateMsgState);
	// }
	log.Println("message reciept")
	u.Server.updateMessageStoreToDb(&updateMsgState)
	if u.UName == msg.To {
		return
	}
	select {
	case u.Ch <- &updateMsgState:
		log.Println(u.UName, updateMsgState)
	default:
		log.Printf("client %v is disconnected.", u.UName)
	}
}

func (u *User) WriteMessage(msg *models.WsMessage) {
	updateMsgState := models.WsMessage{MessageType: msg.MessageType, To: msg.To, From: msg.From, Body: msg.Body, ReplyMsg: msg.ReplyMsg, File: msg.File, MessageDate: msg.MessageDate, MessageTime: msg.MessageTime, MessageId: msg.MessageId, MessageStageId: msg.MessageStageId}
	// //log.Println("writing into to user")
	if updateMsgState.MessageStageId == "0" {
		updateMsgState.MessageStageId = "1"
	} else if updateMsgState.MessageStageId == "1" {
		updateMsgState.MessageStageId = "2"
	} else {
		//log.Println("unknown stage ");
	}
	select {
	case u.Ch <- &updateMsgState:
		//log.Println(u.Name)
	default:
		log.Printf("client %v is disconnected.", u.UName)
	}
}
func (u *User) WriteFileReciept(msg *models.WsMessage) {
	updateMsgState := models.WsMessage{MessageType: msg.MessageType, To: msg.To, From: msg.From, MessageId: msg.MessageId, MessageStageId: msg.MessageStageId}
	updateMsgState.MessageType = "0"
	if updateMsgState.MessageStageId == "0" {
		updateMsgState.MessageStageId = "1"
	} else if updateMsgState.MessageStageId == "1" {
		updateMsgState.MessageStageId = "2"
	} else {
		log.Println("unknown stage ")
	}
	log.Println(u.UName, updateMsgState)
	select {
	case u.Ch <- &updateMsgState:
		log.Println(u.UName, updateMsgState)
	default:
		log.Printf("client %v is disconnected.", u.UName)
	}
}

func (u *User) WriteFile(msg *models.WsMessage) {
	// //log.Println("writing into to user")
	// updateMsgState.MessageType="4"
	// updateMsgState.MessageStageId="1"
	// updateMsgState.File.Filejourneystate="start"
	// updateMsgState.File.Sliceresponse="false"
	// if err := websocket.JSON.Send(u.Ws, &updateMsgState); err != nil {
	//             log.Println(err);
	//    }
	if msg.MessageStageId != "2" && u.UName == msg.To {
		msg.MessageType = "4"
	} else {
		msg.MessageType = "6"
	}
	var Fjrnst = "start"
	for i, value := range msg.File.Filedata {
		if i == (len(msg.File.Filedata) - 1) {
			Fjrnst = "end"
		}
		updateMsgState := models.WsMessage{
			MessageType: msg.MessageType,
			To:          msg.To,
			From:        msg.From,
			Body:        msg.Body,
			ReplyMsg:    msg.ReplyMsg,
			File: &models.File{
				Filename:      msg.File.Filename,
				Fileextension: msg.File.Fileextension,
				Filesize:      msg.File.Filesize,
				Filedata:      make([]string, 0),
				// Filesliceddata: value,
				Filejourneystate: Fjrnst,
			},
			MessageDate:    msg.MessageDate,
			MessageTime:    msg.MessageTime,
			MessageId:      msg.MessageId,
			MessageStageId: msg.MessageStageId,
		}
		updateMsgState.File.Filedata = append(updateMsgState.File.Filedata, value)
		if err := websocket.JSON.Send(u.Ws, &updateMsgState); err != nil {
			//log.Println(err);
		}
	}
	u.Syncsend = false
	log.Println("end write file")
	return
}
func (u *User) WriteInfoStatus(msg *models.WsMessage) {
	select {
	case u.Ch <- msg:
		// log.Println(u.UName)
	default:
		log.Printf("client %v is disconnected.", u.UName)
	}
}
func (u *User) listenWrite() {
	//log.Println("Listening write to client")
	for {
		select {
		case msg := <-u.Ch:
			websocket.JSON.Send(u.Ws, msg)
			u.Syncsend = false
			// log.Println(u.UName,msg)
		case <-u.DoneCh:
			u.Server.Del(u)
			u.DoneCh <- true
			return
		}
	}
}

func (u *User) listenRead() {
	//log.Println("Listening read from client")
	for {
		select {
		default:
			var msg models.WsMessage
			err := websocket.JSON.Receive(u.Ws, &msg)
			if err == io.EOF {
				log.Println("EOF err: ", err)
				u.DoneCh <- true
			} else if err != nil {
				log.Println("err: ", err)
			} else {
				// log.Printf("%#v",&msg)
				if msg.MessageType == "0" {
					// log.Println("got rcp type")
					u.Server.SendReceipt(&msg)
				} else if msg.MessageType == "1" {
					//log.Println("got msg type")
					u.Server.SendMessage(&msg)
				} else if msg.MessageType == "2" {
					// filerequest(u, &msg)
				} else if msg.MessageType == "3" {
					// fileslicing(u, &msg)
				} else if msg.MessageType == "5" {
					u.Server.getDatabaseMessages(&msg)
				} else if msg.MessageType == "77" {
					u.Server.sendUsersStatusToNewUser(&msg)
				} else if msg.MessageType == "8" {
					u.Server.usersSearchList(&msg)
				} else if msg.MessageType == "9" {
					u.Server.getFriendList(&msg)
				} else if msg.MessageType == "10" {
					u.Server.userTyping(&msg)
				} else if msg.MessageType == "11" {
					u.Server.friendrequest(&msg)
				}
			}
		}
	}
}

// var newfile *models.File

// func filerequest(u *User, msg *models.WsMessage) {
// 	if msg.File.Filemessagetype == "file" {
// 		msg.File.FileRequest = true
// 		newfile = new(models.File)
// 		newfile.Filename = msg.File.Filename
// 		newfile.RefFileId = msg.File.RefFileId
// 		newfile.Fileextension = msg.File.Fileextension
// 		newfile.Filesize = msg.File.Filesize
// 		newfile.Filedata = make([]string, 0)

// 		// log.Println("%v",msg.File);
// 		if err := websocket.JSON.Send(u.Ws, &msg); err != nil {
// 			//log.Println(err);
// 		}
// 		//log.Println("newfile");
// 	} else {
// 		//log.Println("not a file",msg);
// 	}
// }
// func fileslicing(u *User, msg *models.WsMessage) {
// 	if msg.File.Filejourneystate == "ini" {
// 		newfile.Filedata = append(newfile.Filedata, msg.File.Filesliceddata)
// 		msg.File.Sliceresponse = true
// 		if err := websocket.JSON.Send(u.Ws, &msg); err != nil {
// 			//log.Println("err:",err)
// 		}
// 		//log.Println("ini")
// 	} else if msg.File.Filejourneystate == "end" {
// 		newfile.Filedata = append(newfile.Filedata, msg.File.Filesliceddata)
// 		msg.File.Sliceresponse = true
// 		if err := websocket.JSON.Send(u.Ws, &msg); err != nil {
// 			//log.Println("err:",err)
// 		} else {
// 			//log.Println("sending last response")
// 			msg.File = newfile
// 			//log.Println("end")
// 			u.Server.SendFile(msg)
// 		}
// 	} else {
// 		//log.Println("s");
// 	}

// }
