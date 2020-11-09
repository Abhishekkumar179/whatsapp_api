package models

type WsResponse struct {
	MessageId       string              `json:"message_id,omitempty"`
	ResponseMessage string              `json:",omitempty"`
	ResponseCode    int                 `json:",omitempty"`
	WsMessage       *WsMessage          `json:",omitempty"`
	Customer        *ReceiveUserDetails `json:"customer,omitempty"`
}

type WsMessage struct {
	MessageId string `json:"message_id,omitempty"`
	UserId    string `json:"user_id,omitempty"`
	UserType  string `json:"user_type,omitempty"`
}
