package models

type WsMessage struct {
	MessageType    string     `json:"message_type,omitempty"`
	To             string     `json:"to,omitempty"`
	From           string     `json:"from,omitempty"`
	Body           string     `json:"body,omitempty"`
	ReplyMsg       *ReplyMsg  `json:"reply_msg,omitempty"`
	File           *File      `json:"file,omitempty"`
	MessageDate    string     `json:"message_date,omitempty"`
	MessageTime    string     `json:"message_time,omitempty"`
	MessageId      string     `json:"message_id,omitempty"`
	MessageStageId string     `json:"message_stage_id,omitempty"`
	UserStatus     string     `json:"user_status,omitempty"`
	NewFriend      string     `json:"new_friend,omitempty"`
	FriendName     string     `json:"friend_name,omitempty"`
	FriendId       string     `json:"friend_id,omitempty"`
	LastSeenDate   string     `json:"last_seen_date,omitempty"`
	LastSeenTime   string     `json:"last_seen_time,omitempty"`
	UserTyping     string     `json:"user_typing,omitempty"`
	UserSearchList []UserData `json:"user_data,omitempty"`
}
type UserData struct {
	UserId       string `json:"user_id,omitempty"`
	Username     string `json:"username,omitempty"`
	LastSeenDate string `json:"last_seen_date,omitempty"`
	LastSeenTime string `json:"last_seen_time,omitempty"`
}
type File struct {
	Filemessagetype  string   `json:"Filemessagetype,omitempty"`
	Filename         string   `json:"filename,omitempty"`
	RefFileId        string   `json:"ref_file_id,omitempty"`
	Fileextension    string   `json:"fileextension,omitempty"`
	Filesize         int      `json:"filesize,omitempty"`
	FileRequest      bool     `json:"filerequest,omitempty"`
	Filedata         []string `json:"filedata,omitempty"`
	Filesliceddata   string   `json:"filesliceddata,omitempty"`
	Filejourneystate string   `json:"filejourneystate,omitempty"`
	Sliceresponse    bool     `json:"sliceresponse,omitempty"`
}

type ReplyMsg struct {
	ReplyMsgId       string `json:"reply_msg_id,omitempty"`
	ReplyMsgOwner    string `json:"reply_msg_owner,omitempty"`
	ReplyMsgBody     string `json:"reply_msg_body,omitempty"`
	ReplyMsgFilename string `json:"reply_msg_filename,omitempty"`
	ReplyMsgFilelogo string `json:"reply_msg_filelogo,omitempty"`
}

type MyResponse struct {
	StatusCode    string    `json:"statuscode,omitempty"`
	StatusMessage string    `json:"statusmsg,omitempty"`
	Msg           WsMessage `json:"msg,omitempty"`
}
type Register struct {
	MessageType  string `json:"message_type,omitempty"`
	UName        string `json:"u_name,omitempty"`
	UType        string `json:"u_type,omitempty"`
	Name         string `json:"name,omitempty"`
	OnlineStatus bool   `json:"online_status,omitempty"`
	LastSeenDate string `json:"last_seen_date,omitempty"`
	LastSeenTime string `json:"last_seen_time,omitempty"`
}

type ProductItem struct {
	Status        string   `json:",omitempty"`
	Msg           string   `json:",omitempty"`
	ResponseCode  int      `json:",omitempty"`
	AccountId     string   `json:"account_id,omitempty"`
	TransactionId string   `json:"transaction_id,omitempty"`
	items         Mystruct `json:"items,omitempty"`
}
type Products struct {
	ProductId string `json:"product_id,omitempty"`
	PQuantity string `json:"p_quantity,omitempty"`
}
type Mystruct struct {
	inputs []*Products `json:"inputs,omitempty"`
}
