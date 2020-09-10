package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Role      string  `json:"role,omitempty"`
	Type      string  `json:"type,omitempty"`
	Text      string  `json:"text,omitempty"`
	MediaType string  `json:"mediaType,omitempty"`
	MediaUrl  string  `json:"mediaUrl,omitempty"`
	Items     []Items `json:"items,omitempty"`
	// Coordinates Coordinates `json:"coordinates,omitempty"`
	// Location    Location    `json:"location,omitempty"`
	Action    []Actions `json:"actions,omitempty"`
	Name      string    `json:"name,omitempty"`
	AvatarURL string    `json:"avatarUrl,omitempty"`
}
type Locations struct {
	Role        string      `json:"role,omitempty"`
	Type        string      `json:"type,omitempty"`
	Text        string      `json:"text,omitempty"`
	Coordinates Coordinates `json:"coordinates,omitempty"`
	Location    Location    `json:"location,omitempty"`
}
type Coordinates struct {
	Lat  float64 `json:"lat,omitempty"`
	Long float64 `json:"long,omitempty"`
}
type Location struct {
	Address string `json:"address,omitempty"`
	Name    string `json:"name,omitempty"`
}
type Payload struct {
	Name    string  `json:"name,omitempty"`
	Message Message `json:"message"`
}
type Comtemplate struct {
	Name    string `json:"name,omitempty"`
	Message User   `json:"message"`
}
type Message struct {
	Role  string  `json:"role,omitempty"`
	Type  string  `json:"type,omitempty"`
	Text  string  `json:"text,omitempty"`
	Items []Items `json:"items,omitempty"`
}

type Items struct {
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	MediaType   string    `json:"mediaType,omitempty"`
	MediaUrl    string    `json:"mediaUrl,omitempty"`
	Actions     []Actions `json:"actions,omitempty"`
}
type Actions struct {
	Type     string `json:"type,omitempty"`
	Text     string `json:"text,omitempty"`
	Payload  string `json:"payload,omitempty"`
	URI      string `json:"uri,omitempty"`
	Amount   int64  `json:"amount,omitempty"`
	Size     string `json:"size,omitempty"`
	Fallback string `json:"fallback,omitempty"`
}

type Data struct {
	AppUser struct {
		Id                  string           `json:"_id,omitempty"`
		UserId              string           `json:"userId,omitempty"`
		Surname             string           `json:"surname,omitempty"`
		GivenName           string           `json:"givenName,omitempty"`
		SignedUpAt          time.Time        `json:"signedUpAt,omitempty" gorm:"type:timestamps;default:'2000-01-01 00:00:01';not null"`
		HasPaymentInfo      bool             `json:"hasPaymentInfo,omitempty"`
		ConversationStarted bool             `json:"conversationStarted,omitempty"`
		Clients             []Clients        `json:"clients,omitempty"`
		PendingClients      []PendingClients `json:"pendingClients,omitempty"`
		Properties          Properties       `json:"properties,omitempty"`
	} `json:"appUser,omitempty"`
}
type Appuser struct {
	Id                  string           `json:"_id,omitempty"`
	UserId              string           `json:"userId,omitempty"`
	Surname             string           `json:"surname,omitempty"`
	GivenName           string           `json:"givenName,omitempty"`
	SignedUpAt          time.Time        `gorm:"type:timestamp;" json:"signedUpAt,omitempty"`
	HasPaymentInfo      bool             `json:"hasPaymentInfo,omitempty"`
	ConversationStarted bool             `json:"conversationStarted,omitempty"`
	Clients             []Clients        `json:"clients,omitempty"`
	PendingClients      []PendingClients `json:"pending_clients,omitempty"`
	Properties          Properties       `json:"properties,omitempty"`
}
type Properties struct {
}
type PendingClients struct {
}
type Clients struct {
}
type Response struct {
	ResponseCode                  int                              `json:",omitempty"`
	Status                        string                           `json:",omitempty"`
	Msg                           string                           `json:",omitempty"`
	Appuser                       *Data                            `json:",omitempty"`
	Data                          []byte                           `json:",omitempty"`
	AppUserList                   []ReceiveUserDetails             `json:",omitempty"`
	Customer                      []Customer_Agents                `json:",omitempty"`
	Message                       *Payload                         `json:",omitempty"`
	Received                      *Received                        `json:",omitempty"`
	Tenant_details                *Tenant_details                  `json:",omitempty"`
	Tenant_list                   []Tenant_details                 `json:",omitempty"`
	WhatsappConfiguration         *WhatsappConfiguration           `json:",omitempty"`
	List                          []WhatsappConfiguration          `json:",omitempty"`
	Fb                            []FacebookConfiguration          `json:",omitempty"`
	AssignAgent                   []AgentQueue                     `json:",omitempty"`
	Queue                         []Queue                          `json:",omitempty"`
	Agent                         []V_call_center_agents           `json:",omitempty"`
	FacebookGetCode               *FacebookGetCode                 `json:",omitempty"`
	FacebookGetAuthInfo           *FacebookGetAuthInfo             `json:",omitempty"`
	FacebookLoginAppConfiguration *[]FacebookLoginAppConfiguration `json:",omitempty"`
}

type AutoGenerated struct {
	Destination Destination `json:"destination"`
	Author      Author      `json:"author"`
	Message     Messages    `json:"message"`
}
type Destination struct {
	IntegrationID string `json:"integrationId"`
	DestinationID string `json:"destinationId"`
}
type Author struct {
	Role string `json:"role"`
}
type Language struct {
	Policy string `json:"policy"`
	Code   string `json:"code"`
}
type Image struct {
	Link string `json:"link"`
}
type Parameters struct {
	Type  string `json:"type"`
	Image Image  `json:"image"`
}
type Components struct {
	Type       string       `json:"type"`
	Parameters []Parameters `json:"parameters"`
}
type Template struct {
	Namespace  string       `json:"namespace"`
	Name       string       `json:"name"`
	Language   Language     `json:"language"`
	Components []Components `json:"components"`
}
type Messages struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
type App struct {
	ID string `json:"id,omitempty"`
}
type Appusers struct {
	ID                  string     `json:"_id,omitempty"`
	Surname             string     `json:"surname,omitempty"`
	GivenName           string     `json:"givenName,omitempty"`
	SignedUpAt          time.Time  `json:"signedUpAt,omitempty"`
	Properties          Properties `json:"properties,omitempty"`
	ConversationStarted bool       `json:"conversationStarted,omitempty"`
}
type Conversations struct {
	ID string `json:"_id,omitempty"`
}
type Messages1 struct {
	Type     string  `json:"type,omitempty"`
	Text     string  `json:"text,omitempty"`
	Role     string  `json:"role,omitempty"`
	Received float64 `json:"received,omitempty"`
	Name     string  `json:"name,omitempty"`
	AuthorID string  `json:"authorId,omitempty"`
	ID       string  `json:"_id,omitempty"`
	Source   Source  `json:"source,omitempty"`
}
type Source struct {
	OriginalMessageID        string  `json:"originalMessageId,omitempty"`
	OriginalMessageTimestamp float64 `json:"originalMessageTimestamp,omitempty"`
	Type                     string  `json:"type,omitempty"`
	IntegrationID            string  `json:"integrationId,omitempty"`
}
type Received struct {
	Trigger      string        `json:"trigger,omitempty"`
	Version      string        `json:"version,omitempty"`
	App          App           `json:"app,omitempty"`
	AppUser      Appusers      `json:"appUser,omitempty"`
	Conversation Conversations `json:"conversation,omitempty"`
	Messages     []Messages1   `json:"messages,omitempty"`
}
type ReceiveUserDetails struct {
	Trigger                  string    `json:"trigger,omitempty"`
	Version                  string    `json:"version,omitempty"`
	AppId                    string    `json:"appId,omitempty"`
	AppUserId                string    `json:"appUserId,omitempty"`
	Surname                  string    `json:"surname,omitempty"`
	GivenName                string    `json:"givenName,omitempty"`
	SignedUpAt               time.Time `json:"signedUpAt,omitempty"`
	ConversationStarted      bool      `json:"conversationStarted,omitempty"`
	Conversation_id          string    `json:"conversation_id,omitempty"`
	Type                     string    `json:"type,omitempty"`
	Text                     string    `json:"text,omitempty"`
	Role                     string    `json:"role,omitempty"`
	Received                 float64   `json:"received,omitempty"`
	Name                     string    `json:"name,omitempty"`
	AuthorID                 string    `json:"authorId,omitempty"`
	Message_id               string    `json:"Message_id,omitempty"`
	OriginalMessageID        string    `json:"originalMessageId,omitempty"`
	OriginalMessageTimestamp float64   `json:"originalMessageTimestamp,omitempty"`
	Source_Type              string    `json:"source_type,omitempty"`
	IntegrationID            string    `json:"integrationId,omitempty"`
	Is_enabled               bool      `json:"is_enabled,omitempty"`
	UnreadCount              int64     `json:"unread_count,omitempty"`
}

type Tenant_details struct {
	Id                int64  `gorm:"PRIMARY_KEY" json:"id,omitempty"`
	ConfigurationName string `json:"configuration_name,omitempty"`
	Domain_uuid       string `json:"domain_uuid,omitempty" gorm:"type:uuid"`
	AppId             string `json:"app_id,omitempty"`
	AppKey            string `json:"app_key,omitempty"`
	AppSecret         string `json:"app_secret,omitempty"`
}

type WhatsappConfigurations struct {
	Id                    int64  `gorm:"PRIMARY_KEY" json:"id,omitempty"`
	ConfigurationName     string `json:"configuration_name,omitempty"`
	Domain_uuid           string `json:"domain_uuid,omitempty"`
	AppId                 string `json:"appId,omitempty"`
	AppKey                string `json:"appKey,omitempty"`
	AppSecret             string `json:"appSecret,omitempty"`
	Message               string `json:"message,omitempty"`
	Size                  string `json:"size,omitempty"`
	WhatsappIntegrationID string `json:"whatsapp_integration_id,omitempty"`
	Trigger               Trigger
	WorkingDays           []WorkingDays `json:"working_days,omitempty"`
}
type WorkingDays struct {
	Day                  string `json:"day,omitempty"`
	WorkingHourStartTime string `json:"workingHourStartTime,omitempty"`
	WorkingHourEndTime   string `json:"workingHourEndTime,omitempty"`
}
type WhatsappConfiguration struct {
	Id                    int64  `gorm:"PRIMARY_KEY" json:"id,omitempty"`
	ConfigurationName     string `json:"configuration_name,omitempty"`
	Domain_uuid           string `json:"domain_uuid,omitempty"`
	AppId                 string `json:"appId,omitempty"`
	AppKey                string `json:"appKey,omitempty"`
	AppSecret             string `json:"appSecret,omitempty"`
	Message               string `json:"message,omitempty"`
	Size                  string `json:"size,omitempty"`
	WhatsappIntegrationID string `json:"whatsapp_integration_id,omitempty"`
	Day1                  string `json:"day1,omitempty"`
	Day2                  string `json:"day2,omitempty"`
	Day3                  string `json:"day3,omitempty"`
	Day4                  string `json:"day4,omitempty"`
	Day5                  string `json:"day5,omitempty"`
	Day6                  string `json:"day6,omitempty"`
	Day7                  string `json:"day7,omitempty"`
	Workstart1            string `json:"workstart1,omitempty"`
	Workstart2            string `json:"workstart2,omitempty"`
	Workstart3            string `json:"workstart3,omitempty"`
	Workstart4            string `json:"workstart4,omitempty"`
	Workstart5            string `json:"workstart5,omitempty"`
	Workstart6            string `json:"workstart6,omitempty"`
	Workstart7            string `json:"workstart7,omitempty"`
	Workend1              string `json:"workend1,omitempty"`
	Workend2              string `json:"workend2,omitempty"`
	Workend3              string `json:"workend3,omitempty"`
	Workend4              string `json:"workend4,omitempty"`
	Workend5              string `json:"workend5,omitempty"`
	Workend6              string `json:"workend6,omitempty"`
	Workend7              string `json:"workend7,omitempty"`
	TriggerWhen           string `json:"trigger_when,omitempty"`
	TriggerName           string `json:"trigger_name,omitempty"`
	TriggerMessage        string `json:"trigger_message,omitempty"`
}
type FacebookConfigurations struct {
	Id                    int64  `gorm:"PRIMARY_KEY" json:"id,omitempty"`
	ConfigurationName     string `json:"configuration_name,omitempty"`
	Domain_uuid           string `json:"domain_uuid,omitempty"`
	AppId                 string `json:"appId,omitempty"`
	AppKey                string `json:"appKey,omitempty"`
	AppSecret             string `json:"appSecret,omitempty"`
	FacebookIntegrationID string `json:"facebook_integration_id,omitempty"`
	Message               string `json:"message,omitempty"`
	Size                  string `json:"size,omitempty"`
	Trigger               Trigger
	WorkingDays           []WorkingDays `json:"working_days,omitempty"`
}
type Trigger struct {
	When    string `json:"when,omitempty"`
	Name    string `json:"name,omitempty"`
	Message string `json:"message,omitempty"`
}
type FacebookConfiguration struct {
	Id                    int64  `gorm:"PRIMARY_KEY" json:"id,omitempty"`
	ConfigurationName     string `json:"configuration_name,omitempty"`
	Domain_uuid           string `json:"domain_uuid,omitempty"`
	AppId                 string `json:"appId,omitempty"`
	AppKey                string `json:"appKey,omitempty"`
	AppSecret             string `json:"appSecret,omitempty"`
	Message               string `json:"message,omitempty"`
	Size                  string `json:"size,omitempty"`
	FacebookIntegrationID string `json:"facebook_integration_id,omitempty"`
	Day1                  string `json:"day1,omitempty"`
	Day2                  string `json:"day2,omitempty"`
	Day3                  string `json:"day3,omitempty"`
	Day4                  string `json:"day4,omitempty"`
	Day5                  string `json:"day5,omitempty"`
	Day6                  string `json:"day6,omitempty"`
	Day7                  string `json:"day7,omitempty"`
	Workstart1            string `json:"workstart1,omitempty"`
	Workstart2            string `json:"workstart2,omitempty"`
	Workstart3            string `json:"workstart3,omitempty"`
	Workstart4            string `json:"workstart4,omitempty"`
	Workstart5            string `json:"workstart5,omitempty"`
	Workstart6            string `json:"workstart6,omitempty"`
	Workstart7            string `json:"workstart7,omitempty"`
	Workend1              string `json:"workend1,omitempty"`
	Workend2              string `json:"workend2,omitempty"`
	Workend3              string `json:"workend3,omitempty"`
	Workend4              string `json:"workend4,omitempty"`
	Workend5              string `json:"workend5,omitempty"`
	Workend6              string `json:"workend6,omitempty"`
	Workend7              string `json:"workend7,omitempty"`
	TriggerWhen           string `json:"trigger_when,omitempty"`
	TriggerName           string `json:"trigger_name,omitempty"`
	TriggerMessage        string `json:"trigger_message,omitempty"`
}
type Link struct {
	Type         string       `json:"type,omitempty"`
	Confirmation Confirmation `json:"confirmation,omitempty"`
	Address      string       `json:"address,omitempty"`
}
type Confirmation struct {
	Type string `json:"type,omitempty"`
}
type Queue struct {
	Id            int64  `gorm:"PRIMARY_KEY" json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	IntegrationID string `json:"integration_id,omitempty"`
	Queue_uuid    string `json:"queue_uuid,omitempty" gorm:"type:uuid"`
	Map_with      string `json:"map_with,omitempty"`
	Domain_uuid   string `json:"domain_uuid,omitempty" gorm:"type:uuid"`
}
type AgentQueue struct {
	QueueName          string `json:"queue_name,omitempty"`
	AgentName          string `json:"agent_name,omitempty"`
	Queue_uuid         string `json:"queue_uuid,omitempty" gorm:"type:uuid"`
	Agent_uuid         string `json:"agent_uuid,omitempty" gorm:"type:uuid"`
	Tenant_domain_uuid string `json:"tenant_domain_uuid,omitempty" gorm:"type:uuid"`
}
type V_call_center_agents struct {
	CallCenterAgentUUID uuid.UUID `json:"call_center_agent_uuid,omitempty" gorm:"type:uuid"`
	AgentName           string    `json:"agent_name,omitempty"`
}
type Customer_Agents struct {
	Domain_uuid string `json:"domain_uuid,omitempty" gorm:"type:uuid"`
	AppUserId   string `json:"appUserId,omitempty"`
	Agent_uuid  string `json:"agent_uuid,omitempty"`
}
type Count_Agent_customer struct {
	Agent_uuid         string `json:"agent_uuid,omitempty"`
	Count              int64  `json:"count,omitempty"`
	Tenant_domain_uuid string `json:"tenant_domain_uuid,omitempty"`
}
type FacebookGetCode struct {
	Code string
}
type FacebookGetAuthInfo struct {
	AccessToken string
	Id          string
	Name        string
}
type FacebookLoginAppConfiguration struct {
	FlacUUID   string `json:"flac_uuid" gorm:"type:uuid;default:uuid_generate_v4();"`
	DomainUUID string `json:"domain_uuid"`
	AppId      string `json:"app_id"`
	AppSecret  string `json:"app_secret"`
	AppName    string `json:"app_name"`
}

// 97a8b191dc46299a5eb349ea
// 6a88e08e1828ca95ea97f305
// 93b6e7aac00d360dec27cc1b
// fd97b30e00bef024d299f90d
// 866f9006f016e19489194921
// 0372777fd23e11b60f009d65
// 21faa13735fe544ba7396b5b
// 454b8b4a14d40e44326875bb
// 4a54c37fe46fd12409aa3e42
// d23095efa00166fff3af155a
// e3c249966d31f73b943fd168
// 19709560f99d553e86939ea6
// 96a4aee8726abacb8f9555b6
// 58a5627d25aca54f82b381b8
// update receive_user_details set unread_count = 0 where app_user_id = '58a5627d25aca54f82b381b8';
