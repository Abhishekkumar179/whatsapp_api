package crud

import (
	"context"
	"mime/multipart"
	models "whatsapp_api/model"

	"github.com/labstack/echo"
)

type Repository interface {
	Delete_AppUser(ctx context.Context, appUserId string, appId string) (*models.Response, error)
	Delete_AppUser_Profile(ctx context.Context, appId string, appUserId string) (*models.Response, error)
	GetAllMessageByAppUserId(ctx context.Context, appUserId string, appId string) ([]byte, error)
	GetAppUserDetails(ctx context.Context, appUserId string, appId string) ([]byte, error)
	Get_allId(ctx context.Context, domain_uuid string) (*models.Response, error)
	Get_Customer_by_agent_uuid(ctx context.Context, customer_id string) (*models.Response, error)
	App_user(ctx context.Context, body []byte) (*models.Response, error)
	Pre_createUser(ctx context.Context, appId string, id int64, userId string, surname string, givenName string) (*models.Response, error)
	Update_AppUser(tx context.Context, appUserId string, appId string, surname string, givenName string) (*models.Response, error)
	Add_Smooch_configuration(ctx context.Context, name string, domain_uuid string, appId string, appKey string, appSecret string) (*models.Response, error)
	Update_Smooch_configuration(ctx context.Context, id int64, name string, domain_uuid string, appId string, appKey string, appSecret string) (*models.Response, error)
	Delete_Smooch_configuration(ctx context.Context, id int64, domain_uuid string) (*models.Response, error)
	Get_Smooch_configuration(ctx context.Context, domain_uuid string) (*models.Response, error)
	Add_Whatsapp_configuration(ctx context.Context, td models.WhatsappConfigurations) (*models.Response, error)
	Get_Whatsapp_configuration(ctx context.Context, domain_uuid string) (*models.Response, error)
	Update_Whatsapp_configuration(ctx context.Context, id int64, domain_uuid string, td models.WhatsappConfigurations) (*models.Response, error)
	Delete_Whatsapp_configuration(ctx context.Context, id int64, domain_uuid string) (*models.Response, error)
	Add_Facebook_configuration(ctx context.Context, td models.FacebookConfigurations) (*models.Response, error)
	Get_Facebook_configuration(ctx context.Context, domain_uuid string) (*models.Response, error)
	Update_Facebook_configuration(ctx context.Context, id int64, domain_uuid string, td models.FacebookConfigurations) (*models.Response, error)
	Delete_Facebook_configuration(ctx context.Context, id int64, domain_uuid string) (*models.Response, error)
	Add_Twitter_configuration(ctx context.Context, td models.TwitterConfigurations) (*models.Response, error)
	Get_Twitter_configuration(ctx context.Context, domain_uuid string) (*models.Response, error)
	Update_Twitter_configuration(ctx context.Context, id int64, domain_uuid string, td models.TwitterConfigurations) (*models.Response, error)
	Delete_Twitter_configuration(ctx context.Context, id int64, domain_uuid string) (*models.Response, error)
	List_integration(ctx context.Context, appId string) ([]byte, error)
	DeleteAllMessage(ctx context.Context, appUserId string, appId string) (*models.Response, error)
	PostMessage(ctx context.Context, appId string, ConversationId string, p models.User) ([]byte, error)
	DeleteMessage(ctx context.Context, appId string, appUserId string, messageId string) (*models.Response, error)
	Create_Text_Template(ctx context.Context, appId string, p models.Comtemplate) ([]byte, error)
	Create_Carousel_Template(ctx context.Context, appId string, p models.Payload) ([]byte, error)
	Create_Compound_Template(ctx context.Context, appId string, p models.Comtemplate) ([]byte, error)
	Create_Quickreply_Template(ctx context.Context, appId string, p models.Comtemplate) ([]byte, error)
	Create_Request_Location_Template(ctx context.Context, appId string, p models.Comtemplate) ([]byte, error)
	Update_Text_Template(ctx context.Context, appId string, template_id string, p models.Payload) ([]byte, error)
	Get_template(ctx context.Context, appId string, template_id string) ([]byte, error)
	List_template(ctx context.Context, appId string) ([]byte, error)
	Delete_template(ctx context.Context, appId string, template_id string) (*models.Response, error)
	Send_Location(ctx context.Context, appId string, appUserId string, p models.Locations) ([]byte, error)
	Message_Action_Types(ctx context.Context, appId string, appUserId string, p models.User) ([]byte, error)
	Quickreply_Message(ctx context.Context, appId string, appUserId string, p models.User) ([]byte, error)
	Send_Carousel_Message(ctx context.Context, appId string, appUserId string, p models.User) ([]byte, error)
	Link_appUser_to_Channel(ctx context.Context, appId string, appUserId string, p models.Link) ([]byte, error)
	Unlink_appUser_to_Channel(ctx context.Context, appId string, appUserId string, channel string) ([]byte, error)
	Upload_Attachments(ctx context.Context, displayName string, AvatarURL string, appId string, conversationId string, Type string, Text string, IntegrationID string, Size int64, file multipart.File, handler *multipart.FileHeader) (*models.Response, error)
	TypingActivity(ctx context.Context, appId string, appUserId string, p models.User) ([]byte, error)
	Disable_AppUser(ctx context.Context, appUserId string) (*models.Response, error)
	Reset_Unread_Count(ctx context.Context, appId string, appUserId string) (*models.Response, error)
	Create_Queue(ctx context.Context, Id int64, Queue_uuid string, Map_with string, Name string, IntegrationID string, domain_uuid string) (*models.Response, error)
	Assign_Agent_To_Queue(ctx context.Context, agent_name string, agent_uuid string, queue_name string, tenant_domain_uuid string, queue_uuid string) (*models.Response, error)
	Remove_Agent_From_Queue(ctx context.Context, agent_uuid string) (*models.Response, error)
	Get_Assigned_Agent_list_From_Queue(ctx context.Context, queue_uuid string) (*models.Response, error)
	Get_Queue_List(ctx context.Context, domain_uuid string) (*models.Response, error)
	Get_Available_Agents_Queue_List(ctx context.Context, agent_uuid string, queue_uuid string) (*models.Response, error)
	Update_Queue(ctx context.Context, queue_uuid string, Name string, IntegrationID string, Map_with string, Domain_uuid string) (*models.Response, error)
	Delete_Queue(ctx context.Context, domain_uuid string) (*models.Response, error)
	Available_Agents(ctx context.Context, domain_uuid string, queue_uuid string) (*models.Response, error)
	Transfer_customer(ctx context.Context, agent_name string, conversation_id string, agent_uuid string, appUserId string) (*models.Response, error)
	Publish_Post_on_FB_Page(ctx context.Context, pageId string, message string, access_token string) ([]byte, error)
	Getall_Post_of_Page(ctx context.Context, pageId string, access_token string) ([]byte, error)
	Delete_Post_of_Page(ctx context.Context, page_postId string, access_token string) ([]byte, error)
	Update_Post_of_Page(ctx context.Context, page_postId string, message string, access_token string) ([]byte, error)
	Get_Comments_on_Post_of_Page(ctx context.Context, page_postId string, access_token string) ([]byte, error)
	Get_Likes_on_Post_of_Page(ctx context.Context, page_postId string, access_token string) ([]byte, error)
	Comment_on_Post_of_Page(ctx context.Context, page_postId string, message string, access_token string) ([]byte, error)
	UVoiceFacebookLogin(ctx context.Context, c echo.Context, client_id string, client_secret string, flac_uuid string) (*models.Response, error)
	UVoiceFacebookLoginCallback(ctx context.Context, c echo.Context) (*models.Response, error)
	Get_Page_ID(ctx context.Context, access_token string) ([]byte, error)
	Schedule_Post(ctx context.Context, pageId string, message string, scheduled_publish_time string, access_token string) ([]byte, error)
	AddFacebookApplication(ctx context.Context, domain_uuid string, app_id string, app_secret string, app_name string) (*models.Response, error)
	ShowFacebookApplication(ctx context.Context, domain_uuid string) (*models.Response, error)
	DeleteFacebookApplication(ctx context.Context, domain_uuid string, flac_uuid string) (*models.Response, error)
	Publish_link_with_message_on_Post(ctx context.Context, pageId string, message string, link string, access_token string) ([]byte, error)
	Upload_Photo_on_Post(ctx context.Context, pageId string, access_token string, message string, Type string, file multipart.File, handler *multipart.FileHeader) ([]byte, error)
	AssignAgentToFacebookApplication(ctx context.Context, domain_uuid string, flac_uuid string, agent_uuid string) (*models.Response, error)
	AgentListAssignedToFacebookApplication(ctx context.Context, flac_uuid string) (*models.Response, error)
	AgentListNotInFacebookApplication(ctx context.Context, flac_uuid string, domain_uuid string) (*models.Response, error)
	ShowAgentFacebookApplication(ctx context.Context, agent_uuid string) (*models.Response, error)
	Convert_Access_Token_into_Longlived_Token(ctx context.Context, clientId string, clientSecret string, exchange_token string, access_token string) ([]byte, error)
	RemoveAgentAssignedToFacebookApplication(ctx context.Context, agent_uuid string) (*models.Response, error)
	UpdateFacebookApplication(ctx context.Context, domain_uuid string, flac_uuid string, app_id string, app_secret string, app_name string) (*models.Response, error)
	Send_Private_Message(ctx context.Context, pageId string, postId string, message string, access_token string) ([]byte, error)
	Like_and_Unlike_Post_and_Comment(ctx context.Context, postId string, commentId string, access_token string, Type string) ([]byte, error)
	Delete_Tickets(ctx context.Context, ticket_uuid string) (*models.Response, error)
	GetAll_Tickets(ctx context.Context, domain_uuid string) (*models.Response, error)
	SaveTwitterAuth(ctx context.Context, id int64, domain_uuid string, api_key string, api_secret string, bearer_token string, access_token string, token_secret string) (*models.Response, error)
	UpdateTwitterAuth(ctx context.Context, id int64, domain_uuid string, api_key string, api_secret string, bearer_token string, access_token string, token_secret string) (*models.Response, error)
	GetTwitterAuth(ctx context.Context, domain_uuid string) (*models.Response, error)
	DeleteTwitterAuth(ctx context.Context, id int64, domain_uuid string) (*models.Response, error)
	Twitter_Apis(ctx context.Context, tweet_id string, screen_name string, api_key string, api_type string, author_id string, message string) ([]byte, error)
	AssignAgentToTwitter(ctx context.Context, twitter_uuid string, domain_uuid string, api_key string, agent_uuid string) (*models.Response, error)
	TwitterAssignAgentList(ctx context.Context, domain_uuid string, twitter_uuid string) (*models.Response, error)
	RemoveTwitterAssignAgent(ctx context.Context, agent_uuid string, twitter_uuid string) (*models.Response, error)
	Get_Quoted_Retweet_List(ctx context.Context, api_key string, tweet_id string) (*models.Response, error)
	AssigncustomerToAgent(ctx context.Context, domain_uuid string, agent_uuid string, app_user_id string) (*models.Response, error)
}
