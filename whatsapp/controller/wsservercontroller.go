package controller

import (
	"log"
	"net/http"
	models "whatsapp_api/model"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"golang.org/x/net/websocket"
	// "time"
	// "strconv"
)

type WsController struct {
	slist *ServerUserList
}

func NewServerUser(dbconn *gorm.DB) *ServerUserList {
	Users := make(map[int]*User)
	uChan := make(chan *User)
	delch := make(chan *User)
	donech := make(chan bool)
	return &ServerUserList{
		Users,
		uChan,
		delch,
		donech,
		dbconn,
	}
}

var databaseLocation = "/home/startele/database"

type ServerUserList struct {
	Users    map[int]*User
	UserChan chan *User
	Delch    chan *User
	Donech   chan bool
	Dbconn   *gorm.DB
}

var (
	num            = 0
	channelBufSize = 100
)

func (s *ServerUserList) Del(u *User) {
	log.Printf("delete %#v", u)
	s.Delch <- u
}

func (s *ServerUserList) SendReceipt(msg *models.WsMessage) {
	log.Println("****************SENDING RECEIPTS***********")

	log.Printf("%#v ", msg)
	originusername := new(User)
	tousername := new(User)
	for _, u := range s.Users {
		if u.UName == msg.From {
			tousername = u
			go tousername.WriteMessageReciept(msg)
			break
		}
	}
	for _, u := range s.Users {
		if u.UName == msg.To {
			originusername = u
			go originusername.WriteMessageReciept(msg)
			break
		}
	}
	//log.Println("to user not found")
}

func (s *ServerUserList) SendMessage(msg *models.WsMessage) {
	log.Println("****************SEND Message***********")

	//log.Println("Sendallfc")
	log.Printf("%#v ", msg)
	s.MessageStoreToDb(msg)
	fromusername := new(User)
	tousername := new(User)
	for _, u := range s.Users {
		if u.UName == msg.From {
			fromusername = u
			//log.Println("from user found")
			break
		}
	}
	for _, u := range s.Users {
		if u.UName == msg.To {
			tousername = u
			rcp := msg
			go fromusername.WriteMessageReciept(rcp)
			tousername.WriteMessage(msg)
			return
		}
	}
	fromusername.WriteMessageReciept(msg)
	//log.Println("to user not found")
}

// func (s *ServerUserList) SendFile(msg *models.WsMessage) {
// 	//log.Println("Sendfile")
// 	// log.Printf("%#v ",msg)
// 	s.MessageStoreToDb(msg)
// 	fromusername := new(User)
// 	tousername := new(User)
// 	for _, u := range s.Users {
// 		if u.UName == msg.From {
// 			fromusername = u
// 			//log.Println("from user found")
// 			break
// 		}
// 	}
// 	for _, u := range s.Users {
// 		if u.UName == msg.To {
// 			tousername = u

// 			rcp := msg
// 			go fromusername.WriteMessageReciept(rcp)
// 			tousername.WriteFile(msg)
// 			return
// 		}
// 	}
// 	fromusername.WriteMessageReciept(msg)
// 	//log.Println("to user not found")
// }

func (s *ServerUserList) MessageStoreToDb(msg *models.WsMessage) {
	// log.Printf("%#v",msg)
	// blob := false
	// reply := false
	// if msg.File != nil {
	// 	blob = true
	// }
	// if msg.ReplyMsg != nil {
	// 	reply = true
	// }
	// var id int
	// err := s.Dbconn.QueryRow("insert into user_messages(user_to,user_from,message_number,message,message_date,message_time,message_stage_id,blob_file,reply_msg) values($1,$2,$3,$4,$5,$6,$7,$8,$9) returning message_id;", msg.To, msg.From, msg.MessageId, msg.Body, msg.MessageDate, msg.MessageTime, "1", blob, reply).Scan(&id)
	// if err != nil {
	// 	log.Println(err)
	// }
	// if blob {
	// 	userLocation := databaseLocation + "/" + msg.From
	// 	if err := os.MkdirAll(userLocation, 0755); err != nil {
	// 		log.Println(err)
	// 	} else {
	// 		filename := strings.SplitN(msg.File.Filename, ".", -1)
	// 		fileLocation := userLocation + "/" + filename[0] + ".txt"
	// 		if f, err := os.Create(fileLocation); err != nil {
	// 			log.Println(err)
	// 		} else {
	// 			defer f.Close()
	// 			for _, value := range msg.File.Filedata {
	// 				if _, err := f.WriteString(value + " "); err != nil {
	// 					log.Println(err)
	// 				}
	// 			}
	// 			if err := s.Dbconn.QueryRow("insert into user_blobs(mime_type,file_name,file_data_location,message_id) values($1,$2,$3,$4) ;", msg.File.Fileextension, msg.File.Filename, fileLocation, id); err != nil {
	// 				log.Println(err)
	// 			}
	// 		}
	// 	}
	// }
	// if reply {
	// 	if err := s.Dbconn.QueryRow("insert into reply_msgs(reply_msg_id,reply_msg_owner,reply_msg_body,reply_msg_filelogo,reply_msg_filename,message_id) values($1,$2,$3,$4,$5,$6);", msg.ReplyMsg.ReplyMsgId, msg.ReplyMsg.ReplyMsgOwner, msg.ReplyMsg.ReplyMsgBody, msg.ReplyMsg.ReplyMsgFilelogo, msg.ReplyMsg.ReplyMsgFilename, id); err != nil {
	// 		log.Println(err)
	// 	}
	// }
}
func (s *ServerUserList) updateMessageStoreToDb(msg *models.WsMessage) {
	// _, err := s.Dbconn.Exec("update user_messages set message_stage_id=$1 where user_to=$2 and user_from=$3 and message_number=$4", msg.MessageStageId, msg.To, msg.From, msg.MessageId)
	// if err != nil {
	// 	log.Println(err)
	// }
}
func (s *ServerUserList) getDatabaseMessages(msg *models.WsMessage) {
	// var rows *sql.Rows
	// var err error
	// if msg.MessageId == "0" {
	// 	rows, err = s.Dbconn.Query("select message_id,user_to,user_from,message,message_number,message_date,message_time,message_stage_id,blob_file,reply_msg from user_messages  where (user_to=$1 and user_from=$2) or (user_to=$2 and user_from=$1) order by message_number desc limit 10", msg.To, msg.From)
	// } else {
	// 	rows, err = s.Dbconn.Query("select message_id,user_to,user_from,message,message_number,message_date,message_time,message_stage_id,blob_file,reply_msg from user_messages  where ((user_to=$1 and user_from=$2) or (user_to=$2 and user_from=$1)) and message_number<$3 order by message_number desc limit 10", msg.To, msg.From, msg.MessageId)
	// }
	// if err != nil {
	// 	log.Println(err)
	// }
	// defer rows.Close()
	// fromusername := new(User)
	// var useronline bool
	// for _, u := range s.Users {
	// 	if u.UName == msg.From {
	// 		fromusername = u
	// 		//log.Println("from user found")
	// 		useronline = true
	// 		break
	// 	}
	// }
	// oldstoremsg := make([]*models.WsMessage, 0)
	// newstoremsg := make([]*models.WsMessage, 0)
	// for rows.Next() {
	// 	storemsg := new(models.WsMessage)
	// 	var blob bool
	// 	var reply bool
	// 	var id int
	// 	if err := rows.Scan(&id, &storemsg.To, &storemsg.From, &storemsg.Body, &storemsg.MessageId, &storemsg.MessageDate, &storemsg.MessageTime, &storemsg.MessageStageId, &blob, &reply); err != nil {
	// 		log.Println(err)
	// 	}

	// 	if reply {
	// 		storemsg.ReplyMsg = new(models.ReplyMsg)
	// 		if err := s.Dbconn.QueryRow("select reply_msg_id, reply_msg_owner, reply_msg_body, reply_msg_filelogo, reply_msg_filename from reply_msgs where message_id=$1 ", id).Scan(&storemsg.ReplyMsg.ReplyMsgId, &storemsg.ReplyMsg.ReplyMsgOwner, &storemsg.ReplyMsg.ReplyMsgBody, &storemsg.ReplyMsg.ReplyMsgFilelogo, &storemsg.ReplyMsg.ReplyMsgFilename); err != nil {
	// 			log.Println(err)
	// 		}
	// 	}
	// 	if blob {
	// 		var fileLocation string
	// 		storemsg.File = new(models.File)
	// 		if err := s.Dbconn.QueryRow("select mime_type,file_name,file_data_location from user_blobs  where message_id=$1 ", id).Scan(&storemsg.File.Fileextension, &storemsg.File.Filename, &fileLocation); err != nil {
	// 			log.Println(err)
	// 		}
	// 		content, err := ioutil.ReadFile(fileLocation)
	// 		if err != nil {
	// 			log.Println(err)
	// 		} else {
	// 			cno := string(content)
	// 			storemsg.File.Filedata = strings.Split(cno, " ")
	// 		}
	// 	}
	// 	if storemsg.To == msg.From && storemsg.MessageStageId == "1" {
	// 		storemsg.MessageType = "1"
	// 	} else {
	// 		storemsg.MessageType = msg.MessageType
	// 	}
	// 	// storemsg.MessageStageId="2"
	// 	if storemsg.To == msg.From && storemsg.MessageStageId == "1" {
	// 		newstoremsg = append(newstoremsg, storemsg)
	// 	} else {
	// 		oldstoremsg = append(oldstoremsg, storemsg)
	// 	}
	// 	// log.Println(len(newstoremsg1)
	// }
	// if useronline {
	// 	if len(newstoremsg) > 0 {
	// 		fromusername.sendNewDatabaseMessages(newstoremsg)
	// 	}
	// 	fromusername.sendDatabaseMessages(oldstoremsg)
	// }
}

func (s *ServerUserList) sendUsersStatusToNewUser(msg *models.WsMessage) {
	for _, u := range s.Users {
		if u.UName == msg.To {
			log.Println(msg)
			u.WriteInfoStatus(msg)
		}
	}
}
func (s *ServerUserList) userTyping(msg *models.WsMessage) {
	for _, u := range s.Users {
		if u.UName == msg.To {
			u.WriteData(msg)
			break
		}
	}
}
func (s *ServerUserList) sendUserStatus(name string, stt string) {
	for _, u := range s.Users {
		if u.UName != name {
			msg := models.WsMessage{MessageType: "7", To: u.UName, From: name, UserStatus: stt}
			// log.Println(msg)
			u.WriteInfoStatus(&msg)
		}
	}
}
func (s *ServerUserList) updateUserStatus(name string, id int, stt bool) {
	// for _, u := range s.Users {
	// 	if u.UName == name {
	// 		if stt {
	// 			_, err := s.Dbconn.Exec("update user_logins set user_map_id=$1, online_status=$2 where username=$3;", id, true, name)
	// 			if err != nil {
	// 				log.Println(err)
	// 			}
	// 		} else {
	// 			_, err := s.Dbconn.Exec("update user_logins set user_map_id=$1, online_status=$2 where username=$3;", 0, false, name)
	// 			if err != nil {
	// 				log.Println(err)
	// 			}
	// 		}
	// 	}
	// }
}

func (s *ServerUserList) getFriendList(msg *models.WsMessage) {
	// var rows *sql.Rows
	// var err error
	// rows, err = s.Dbconn.Query("select first_username,second_username from user_friends where (first_username=$1 or second_username=$1) order by id desc;", msg.From)
	// if err != nil {
	// 	log.Println(err)
	// }
	// defer rows.Close()
	// user := new(User)
	// for _, u := range s.Users {
	// 	if u.UName == msg.From {
	// 		user = u
	// 	}
	// }
	// for rows.Next() {
	// 	var friendpacket models.WsMessage

	// 	var firstid, firstname, secondid, secondname, updated_time string
	// 	if err := rows.Scan(&firstid, &secondid); err != nil {
	// 		log.Println(err)
	// 	} else {
	// 		if firstid == msg.From {
	// 			if err := s.Dbconn.QueryRow("select name from users where username=$1", secondid).Scan(&secondname); err != nil {
	// 				log.Println(err)
	// 				return
	// 			}
	// 			if err := s.Dbconn.QueryRow("select updated_time from user_logins where username=$1", secondid).Scan(&updated_time); err != nil {
	// 				log.Println(err)
	// 				return
	// 			}
	// 			s := strings.Split(updated_time, "T")
	// 			friendpacket = models.WsMessage{MessageType: msg.MessageType, To: msg.From, FriendName: secondname, FriendId: secondid, LastSeenDate: s[0], LastSeenTime: s[1]}
	// 		} else {
	// 			if err := s.Dbconn.QueryRow("select name from users where username=$1", firstid).Scan(&firstname); err != nil {
	// 				log.Println(err)
	// 				return
	// 			}
	// 			if err := s.Dbconn.QueryRow("select updated_time from user_logins where username=$1", firstid).Scan(&updated_time); err != nil {
	// 				log.Println(err)
	// 				return
	// 			}
	// 			s := strings.Split(updated_time, "T")
	// 			friendpacket = models.WsMessage{MessageType: msg.MessageType, To: msg.From, FriendName: firstname, FriendId: firstid, LastSeenDate: s[0], LastSeenTime: s[1]}
	// 		}
	// 		user.WriteData(&friendpacket)
	// 	}
	// }
}
func (s *ServerUserList) usersSearchList(msg *models.WsMessage) {
	// var rows *sql.Rows
	// var err error
	// // log.Println(msg.FriendName)
	// rows, err = s.Dbconn.Query("select username,name from users where username like $1 or name like $2 limit 25", "____"+msg.FriendName+"%", msg.FriendName+"%")
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// defer rows.Close()
	// user := new(User)
	// for _, u := range s.Users {
	// 	if u.UName == msg.From {
	// 		user = u
	// 	}
	// }
	// for rows.Next() {
	// 	var userid, name, updated_time string
	// 	if err := rows.Scan(&userid, &name); err != nil {
	// 		log.Println(err)
	// 	} else {
	// 		// log.Println(userid)
	// 		err := s.Dbconn.QueryRow("select updated_time from user_logins where username=$1", userid).Scan(&updated_time)
	// 		switch {
	// 		case err == sql.ErrNoRows:
	// 			log.Println(err)
	// 		case err != nil:
	// 			log.Println(err)
	// 		default:
	// 			s := strings.Split(updated_time, "T")
	// 			userdata := models.UserData{UserId: userid, Username: name, LastSeenDate: s[0], LastSeenTime: s[1]}
	// 			// log.Printf("%#v",userdata)
	// 			msg.UserSearchList = append(msg.UserSearchList, userdata)
	// 		}
	// 	}
	// }
	// user.WriteData(msg)
}

func (s *ServerUserList) friendrequest(msg *models.WsMessage) {
	// if result, err := s.Dbconn.Exec("select * from user_friends where (first_username=$1 and second_username=$2) or (first_username=$2 and second_username=$1);", msg.From, msg.NewFriend); err != nil {
	// 	log.Println(err)
	// } else {
	// 	if rows, err := result.RowsAffected(); err != nil {
	// 		log.Println(err)
	// 	} else if rows == 1 {
	// 		log.Println("friend pair exists")
	// 	} else if rows == 0 {
	// 		// log.Println("add friend pair");
	// 		if msg.From == msg.NewFriend {
	// 			// log.Println("wrong friend pair");
	// 		} else {
	// 			if result, err := s.Dbconn.Exec("insert into user_friends(first_username,second_username) values($1,$2);", msg.From, msg.NewFriend); err != nil {
	// 				log.Println(err)
	// 			} else if rows, err := result.RowsAffected(); err != nil {
	// 				log.Println(err)
	// 			} else if rows == 1 {
	// 				log.Println("added success")
	// 				user := new(User)
	// 				for _, u := range s.Users {
	// 					if u.UName == msg.From {
	// 						user = u
	// 					}
	// 				}

	// 				var friendName string
	// 				if err := s.Dbconn.QueryRow("select name from users where username=$1", msg.NewFriend).Scan(&friendName); err != nil {
	// 					log.Println(err)
	// 				} else {
	// 					var updated_time string
	// 					var friendpacket models.WsMessage
	// 					if err := s.Dbconn.QueryRow("select updated_time from user_logins where username=$1", msg.NewFriend).Scan(&updated_time); err != nil {
	// 						log.Println(err)
	// 						return
	// 					}
	// 					s := strings.Split(updated_time, "T")
	// 					friendpacket = models.WsMessage{MessageType: "9", To: msg.From, FriendName: friendName, FriendId: msg.NewFriend, LastSeenDate: s[0], LastSeenTime: s[1]}

	// 					go user.WriteData(&friendpacket)
	// 				}
	// 				frienduser := new(User)
	// 				for _, u := range s.Users {
	// 					if u.UName == msg.NewFriend {
	// 						frienduser = u
	// 					}
	// 				}

	// 				var fromName string
	// 				if err := s.Dbconn.QueryRow("select name from users where username=$1", msg.From).Scan(&fromName); err != nil {
	// 					log.Println(err)
	// 				} else {
	// 					var updated_time string
	// 					var friendpacket models.WsMessage
	// 					if err := s.Dbconn.QueryRow("select updated_time from user_logins where username=$1", msg.From).Scan(&updated_time); err != nil {
	// 						log.Println(err)
	// 						return
	// 					}
	// 					s := strings.Split(updated_time, "T")
	// 					friendpacket = models.WsMessage{MessageType: "9", To: msg.NewFriend, FriendName: fromName, FriendId: msg.From, LastSeenDate: s[0], LastSeenTime: s[1]}
	// 					go frienduser.WriteData(&friendpacket)
	// 					return
	// 				}
	// 			}
	// 		}
	// 	}
	// }
	// log.Println("some error connot add")
}

// func (s *ServerUserList) allUsers(c echo.Context) {
// 	//log.Println(s.Users)
// 	for k, v := range s.Users {
// 		fmt.Fprintf(c.Response, "%#v -> %#v\n", k, v)
// 	}
// }
func (s *ServerUserList) NewUser(username string, utype string, ws *websocket.Conn) *User {
	for _, u := range s.Users {
		if u.UName == username {
			u.Ws = ws
			return u
		}
	}
	nc := new(User)
	num++
	nc.Id = num
	nc.UName = username
	nc.UType = utype
	nc.Ws = ws
	nc.Ch = make(chan *models.WsMessage, channelBufSize)
	nc.DoneCh = make(chan bool)
	nc.Server = s
	nc.Syncsend = false
	s.UserChan <- nc
	return nc
}
func (s *ServerUserList) userLogin(msg map[string]interface{}, ws *websocket.Conn) {
	// var rowsAffected int64
	var u *User
	// if msg.OnlineStatus == true {
	u = s.NewUser(msg["user_id"].(string), msg["user_type"].(string), ws)
	// }
	// var username string
	// err := s.Dbconn.QueryRow("select username from user_logins where username=$1;", msg.UName).Scan(&username)
	// if err == sql.ErrNoRows {
	// result, err := s.Dbconn.Exec("insert into user_logins(username,user_map_id,online_status) values($1,$2,$3);", u.UName, u.Id, msg.OnlineStatus)
	// if err != nil {
	// 	log.Println(err)
	// }
	// rowsAffected, err = result.RowsAffected()
	// if err != nil {
	// 	log.Println(err)
	// }
	// } else if err != nil {
	// 	log.Println(err)
	// } else {
	// result, err := s.Dbconn.Exec("update user_logins set user_map_id=$1, online_status=$2  where username=$3;", u.Id, msg.OnlineStatus, username)
	// if err != nil {
	// 	log.Println(err)
	// }
	// 	rowsAffected, err = result.RowsAffected()
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// }
	// if msg.OnlineStatus == true {
	// var rsp models.MyResponse
	// if rowsAffected == 1 {
	msg["status_code"] = http.StatusOK
	msg["status_message"] = "Session Started."
	// rsp.StatusCode = "200"
	// rsp.StatusMessage = "valid user"
	if err := websocket.JSON.Send(u.Ws, msg); err != nil {
		log.Println("Can't send", err)
	}
	u.Listen()
	// } else {
	// rsp.StatusCode = "404"
	// rsp.StatusMessage = "invalid user"

	// if err := websocket.JSON.Send(u.Ws, rsp); err != nil {
	// 	log.Println("Can't send", err)
	// }

	// }
	// }
}
func (s *ServerUserList) UserRegister(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		// var msg models.Register
		var msg map[string]interface{}
		var rsp models.MyResponse
		for {
			if err := websocket.JSON.Receive(ws, &msg); err != nil {
				log.Println("Can't receive", err)
				break
			}

			if msg["user_id"].(string) != "" {
				log.Println("valid user", msg)
				for _, oldu := range s.Users {
					if oldu.UName == msg["user_id"] && oldu.UType == msg["user_type"] {
						log.Println("found user: ", oldu)
						msg["status_code"] = http.StatusBadRequest
						msg["status_message"] = "Session exist for user."
						if err := websocket.JSON.Send(ws, msg); err != nil {
							log.Println("Can't send", err)
						}
						s.Del(oldu)

					}
				}
				// var username string
				// err := s.Dbconn.QueryRow("select username from users where username=$1;", msg.UName).Scan(&username)
				// if err == sql.ErrNoRows {
				// 	result, err := s.Dbconn.Exec("insert into users(username,name) values($1,$2);", msg.UName, msg.Name)
				// 	if err != nil {
				// 		log.Println(err)
				// 		// s.Del(u)
				// 		break
				// 	}
				// 	if rows, err := result.RowsAffected(); err != nil {
				// 		log.Println(err)
				// 	} else if rows == 1 {
				// 		log.Println("registration success")
				// 		s.userLogin(&msg, ws)
				// 	} else {
				// 		break
				// 	}
				// } else if err != nil {
				// 	log.Println(err)
				// 	break
				// } else {
				// 	log.Println("registeration exists")
				s.userLogin(msg, ws)
				// }
			}
		}
		rsp.StatusCode = "404"
		rsp.StatusMessage = "invalid user"
		if err := websocket.JSON.Send(ws, rsp); err != nil {
			log.Println("Can't send", err)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func (s *ServerUserList) Controller(e *echo.Echo) {

	e.GET("/main", s.UserRegister)
	// e.GET("/allUsers", s.allUsers)
	//log.Println("controller created")
	for {
		select {
		case u := <-s.UserChan:
			s.Users[u.Id] = u
			// var count int
			// log.Printf("new user %#v", u)

			s.sendUserStatus(u.UName, "true")
			// s.updateUserStatus(u.UName,u.Id,true)
			// for _,_=range s.Users{
			// 	count++
			// }
			// log.Println("Now", count, "users connected.")

		case u := <-s.Delch:
			s.sendUserStatus(u.UName, "false")
			s.updateUserStatus(u.UName, u.Id, false)
			delete(s.Users, u.Id)
		}
	}
}
