package models

import (
	"time"

	"github.com/google/uuid"
)

type Account_details struct {
	Tenant_name string `json:"tenant_name,omitempty"`
	Domain_uuid string `json:"domain_uuid,omitempty" gorm:"type:uuid"`
}
type Sender struct {
	Attachment struct {
		MediaType string `json:"mediaType,omitempty"`
		MediaUrl  string `json:"mediaUrl,omitempty"`
	} `json:"attachment,omitempty"`
}
type SaveTwitterAuth struct {
	Id           int64  `json:"id,omitempty"`
	Twitter_uuid string `json:"twitter_uuid,omitempty"`
	Domain_uuid  string `json:"domain_uuid,omitempty"`
	Api_Key      string `json:"api_key,omitempty"`
	Api_Secret   string `json:"api_secret,omitempty"`
	Bearer_Token string `json:"bearer_token,omitempty"`
	Access_Token string `json:"access_token,omitempty"`
	Token_Secret string `json:"token_secret,omitempty"`
}
type TwitterAssignedAgents struct {
	Twitter_uuid string `json:"twitter_uuid,omitempty"`
	Domain_uuid  string `json:"domain_uuid,omitempty"`
	Agent_uuid   string `json:"agent_uuid,omitempty"`
	Api_Key      string `json:"api_key,omitempty"`
}
type Author struct {
	Type        string `json:"type,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	AvatarURL   string `json:"avatarUrl,omitempty"`
}
type Content struct {
	Type     string `json:"type,omitempty"`
	Text     string `json:"text,omitempty"`
	MediaUrl string `json:"mediaUrl,omitempty"`
}

type User struct {
	Role      string  `json:"role,omitempty"`
	Type      string  `json:"type,omitempty"`
	Text      string  `json:"text,omitempty"`
	MediaType string  `json:"mediaType,omitempty"`
	MediaUrl  string  `json:"mediaUrl,omitempty"`
	Items     []Items `json:"items,omitempty"`
	Author    Author  `json:"author,omitempty"`
	Content   Content `json:"content,omitempty"`
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
	ResponseCode                           int                                       `json:",omitempty"`
	Status                                 string                                    `json:",omitempty"`
	Msg                                    string                                    `json:",omitempty"`
	Appuser                                *Data                                     `json:",omitempty"`
	Data                                   []byte                                    `json:",omitempty"`
	AppUserList                            []ReceiveUserDetails                      `json:",omitempty"`
	Customer                               *ReceiveUserDetails                       `json:",omitempty"`
	Message                                *Payload                                  `json:",omitempty"`
	Received                               *Received                                 `json:",omitempty"`
	Tenant_details                         *Tenant_details                           `json:",omitempty"`
	Tenant_list                            []Tenant_details                          `json:",omitempty"`
	WhatsappConfiguration                  *WhatsappConfiguration                    `json:",omitempty"`
	List                                   []WhatsappConfiguration                   `json:",omitempty"`
	Fb                                     []FacebookConfiguration                   `json:",omitempty"`
	Twitter                                []TwitterConfiguration                    `json:",omitempty"`
	AssignAgent                            []AgentQueue                              `json:",omitempty"`
	Queue                                  []Queue                                   `json:",omitempty"`
	Agent                                  []V_call_center_agents                    `json:",omitempty"`
	FacebookGetCode                        *FacebookGetCode                          `json:",omitempty"`
	FacebookGetAuthInfo                    *FacebookGetAuthInfo                      `json:",omitempty"`
	FacebookLoginAppConfiguration          *[]FacebookLoginAppConfiguration          `json:",omitempty"`
	FacebookLoginAppConfigurationAgentList *[]FacebookLoginAppConfigurationAgentList `json:",omitempty"`
	AgentList                              *[]V_call_center_agents                   `json:",omitempty"`
	TicketList                             []SocialMediaTickets                      `json:",omitempty"`
	TwitterAuthList                        []SaveTwitterAuth                         `json:",omitempty"`
	TwitterAssignAgentList                 []TwitterAssignedAgents                   `json:",omitempty"`
	Quote_retweet_list                     []Result                                  `json:",omitempty"`
	FacebookLikesAndComments               *FacebookLikesAndComments                 `json:",omitempty"`
}

type AutoGenerated struct {
	Destination Destination `json:"destination"`
	Author      Authors     `json:"author"`
	Message     Messages    `json:"message"`
}
type Destination struct {
	IntegrationID string `json:"integrationId"`
	DestinationID string `json:"destinationId"`
}
type Authors struct {
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
	Day                      string    `json:"day,omitempty"`
	Date                     int       `json:"date,omitempty"`
	AfterOfficeTime          bool      `json:"after_office_time,omitempty"`
	Domain_uuid              string    `json:"domain_uuid,omitempty"`
	Agent_Request_uuid       string    `json:"agent_request_uuid,omitempty"`
	Agent_Request_Time       int64     `json:"agent_request_time,omitempty" gorm:"default:0"`
}
type SocialMediaTickets struct {
	Ticket_uuid     string  `json:"ticket_uuid,omitempty" gorm:"type:uuid"`
	Domain_uuid     string  `json:"domain_uuid,omitempty" gorm:"type:uuid"`
	Ticket_name     string  `json:"ticket_name,omitempty"`
	CustomerId      string  `json:"customer_id,omitempty"`
	CustomerName    string  `json:"customer_name,omitempty"`
	Message         string  `json:"message,omitempty"`
	MessageType     string  `json:"message_type,omitempty"`
	IntegrationID   string  `json:"integration_id,omitempty"`
	Source_type     string  `json:"source_type,omitempty"`
	Conversation_id string  `json:"conversation_id,omitempty"`
	Timestamp       float64 `json:"timestamp,omitempty"`
}
type AfterOfficeTime struct {
	Day       string `json:"day,omitempty"`
	Date      string `json:"date,omitempty"`
	AppUserId string `json:"app_user_id,omitempty"`
}
type Tenant_details struct {
	Id                int64  `gorm:"PRIMARY_KEY" json:"id,omitempty"`
	ConfigurationName string `json:"configuration_name,omitempty"`
	Domain_uuid       string `json:"domain_uuid,omitempty"`
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
type TwitterConfigurations struct {
	Id                   int64  `gorm:"PRIMARY_KEY" json:"id,omitempty"`
	ConfigurationName    string `json:"configuration_name,omitempty"`
	Domain_uuid          string `json:"domain_uuid,omitempty"`
	AppId                string `json:"appId,omitempty"`
	AppKey               string `json:"appKey,omitempty"`
	AppSecret            string `json:"appSecret,omitempty"`
	TwitterIntegrationID string `json:"twitter_integration_id,omitempty"`
	Message              string `json:"message,omitempty"`
	Size                 string `json:"size,omitempty"`
	Trigger              Trigger
	WorkingDays          []WorkingDays `json:"working_days,omitempty"`
}
type TwitterConfiguration struct {
	Id                   int64  `gorm:"PRIMARY_KEY" json:"id,omitempty"`
	ConfigurationName    string `json:"configuration_name,omitempty"`
	Domain_uuid          string `json:"domain_uuid,omitempty"`
	AppId                string `json:"appId,omitempty"`
	AppKey               string `json:"appKey,omitempty"`
	AppSecret            string `json:"appSecret,omitempty"`
	Message              string `json:"message,omitempty"`
	Size                 string `json:"size,omitempty"`
	TwitterIntegrationID string `json:"twitter_integration_id,omitempty"`
	Day1                 string `json:"day1,omitempty"`
	Day2                 string `json:"day2,omitempty"`
	Day3                 string `json:"day3,omitempty"`
	Day4                 string `json:"day4,omitempty"`
	Day5                 string `json:"day5,omitempty"`
	Day6                 string `json:"day6,omitempty"`
	Day7                 string `json:"day7,omitempty"`
	Workstart1           string `json:"workstart1,omitempty"`
	Workstart2           string `json:"workstart2,omitempty"`
	Workstart3           string `json:"workstart3,omitempty"`
	Workstart4           string `json:"workstart4,omitempty"`
	Workstart5           string `json:"workstart5,omitempty"`
	Workstart6           string `json:"workstart6,omitempty"`
	Workstart7           string `json:"workstart7,omitempty"`
	Workend1             string `json:"workend1,omitempty"`
	Workend2             string `json:"workend2,omitempty"`
	Workend3             string `json:"workend3,omitempty"`
	Workend4             string `json:"workend4,omitempty"`
	Workend5             string `json:"workend5,omitempty"`
	Workend6             string `json:"workend6,omitempty"`
	Workend7             string `json:"workend7,omitempty"`
	TriggerWhen          string `json:"trigger_when,omitempty"`
	TriggerName          string `json:"trigger_name,omitempty"`
	TriggerMessage       string `json:"trigger_message,omitempty"`
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
	Queue_uuid    string `json:"queue_uuid,omitempty"`
	Map_with      string `json:"map_with,omitempty"`
	Domain_uuid   string `json:"domain_uuid,omitempty"`
}
type AgentQueue struct {
	QueueName          string `json:"queue_name,omitempty"`
	AgentName          string `json:"agent_name,omitempty"`
	Queue_uuid         string `json:"queue_uuid,omitempty"`
	Agent_uuid         string `json:"agent_uuid,omitempty" gorm:"type:uuid"`
	Tenant_domain_uuid string `json:"tenant_domain_uuid,omitempty" gorm:"type:uuid"`
}
type V_call_center_agents struct {
	CallCenterAgentUUID uuid.UUID `json:"call_center_agent_uuid,omitempty" gorm:"type:uuid"`
	AgentName           string    `json:"agent_name,omitempty"`
	AgentStatus         string    `json:"agent_status,omitempty"`
	Domain_uuid         string    `json:"domain_uuid,omitempty" gorm:"type:uuid"`
}
type Customer_Agents struct {
	Domain_uuid string `json:"domain_uuid,omitempty" gorm:"type:uuid"`
	AppUserId   string `json:"appUserId,omitempty"`
	Agent_uuid  string `json:"agent_uuid,omitempty"`
}
type Count_Agent_queue struct {
	Count int64 `json:"count,omitempty"`
}
type Count_customer struct {
	Count int64 `json:"count,omitempty"`
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
type FacebookLoginAppConfigurationAgent struct {
	DomainUUID string
	FlacUUID   string
	AgentUUID  string
}
type FacebookLoginAppConfigurationAgentList struct {
	DomainUUID string
	FlacUUID   string
	AgentUUID  string
	AgentName  string
}
type AccessToken struct {
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   int64  `json:"expires_in,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
}
type QuotedTweet struct {
	Data     []Datas  `json:"data,omitempty"`
	Includes Includes `json:"includes,omitempty"`
}
type Datas struct {
	Text      string `json:"text,omitempty"`
	Author_id string `json:"author_id,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}
type Includes struct {
	Users []Users `json:"users,omitempty"`
}
type Users struct {
	Id       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Name     string `json:"name,omitempty"`
}
type Quoted struct {
	//  []Result
	// Res
}
type Result struct {
	Text      string `json:"text,omitempty"`
	Author_id string `json:"author_id,omitempty"`
	Username  string `json:"username,omitempty"`
	Name      string `json:"name,omitempty"`
	RetweetId string `json:"retweet_id,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}
type FacebookLikesAndComments struct {
	Object string  `json:"object,omitempty"`
	Entry  []Entry `json:"entry,omitempty"`
}
type Entry struct {
	Id      string    `json:"id,omitempty"`
	Time    int64     `json:"time,omitempty"`
	Changes []Changes `json:"changes,omitempty"`
}
type Changes struct {
	Value Value  `json:"value,omitempty"`
	Feild string `json:"feild,omitempty"`
}
type Value struct {
	From        From   `json:"from,omitempty"`
	Post        Post   `json:"post,omitempty"`
	Message     string `json:"message,omitempty"`
	PostId      string `json:"post_id,omitempty"`
	CommentId   string `json:"comment_id,omitempty"`
	CreatedTime int64  `json:"created_time,omitempty"`
	ParentId    string `json:"parent_id,omitempty"`
	Verb        string `json:"verb,omitempty"`
}
type From struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
type Post struct {
	StatusType      string `json:"status_type,omitempty"`
	IsPublished     string `json:"is_published,omitempty"`
	UpdateTime      string `json:"update_time,omitempty"`
	PermaLinkUrl    string `json:"perma_link_url,omitempty"`
	PromotionStatus string `json:"promotion_status,omitempty"`
}

// {
// 	"object":"page",
// "entry":
// [
// 	{
// 		"id":"104315744750790",
// 		"time":1604755296,
// 	"changes":
// 	[
// 		{
// 			"value":
// 			{
// 				"from":
// 			{
// 				"id":"104315744750790",
// 				"name":"Demo"
// 				},
// 				"post":
// 				{
// 					"status_type":"mobile_status_update",
// 					"is_published":true,
// 					"updated_time":"2020-11-07T13:21:33+0000",
// 					"permalink_url":"https://www.facebook.com/permalink.php?story_fbid=105776537938044&id=104315744750790",
// 					"promotion_status":"inactive","id":"104315744750790_105776537938044"
// 					},
// 					"message":"hi",
// 					"post_id":"104315744750790_105776537938044",
// 					"comment_id":"105776537938044_149415636907467",
// 					"created_time":1604755293,
// 					"item":"comment",
// 					"parent_id":"104315744750790_105776537938044",
// 					"verb":"add"
// 					},
// 					"field":"feed"
// 					}
// 					]
// 					}
// 					]
// 				}
