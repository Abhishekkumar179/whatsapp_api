package usecase

import (
	"context"
	"fmt"
	"mime/multipart"
	"strconv"
	models "whatsapp_api/model"
	crud "whatsapp_api/whatsapp"

	"github.com/labstack/echo"
	//	"encoding/json"
)

type crudUsecase struct {
	repository crud.Repository
}

func NewcrudUsecase(repo crud.Repository) crud.Usecase {
	return &crudUsecase{
		repository: repo,
	}
}

/**************************************************Create User***************************************************/

func (r *crudUsecase) Delete_AppUser(ctx context.Context, appUserId string, appId string) (*models.Response, error) {

	return r.repository.Delete_AppUser(ctx, appUserId, appId)
}

/*************************************************Delete AppUser Profile*****************************************/
func (r *crudUsecase) Delete_AppUser_Profile(ctx context.Context, appId string, appUserId string) (*models.Response, error) {

	return r.repository.Delete_AppUser_Profile(ctx, appId, appUserId)
}

/**************************************************Get User***************************************************/

func (r *crudUsecase) GetAllMessageByAppUserId(ctx context.Context, appUserId string, appId string) ([]byte, error) {

	return r.repository.GetAllMessageByAppUserId(ctx, appUserId, appId)
}

/************************************************Get AppUser Details*******************************************/
func (r *crudUsecase) GetAppUserDetails(ctx context.Context, appUserId string, appId string) ([]byte, error) {

	return r.repository.GetAppUserDetails(ctx, appUserId, appId)
}

/**************************************************Update User***************************************************/

func (r *crudUsecase) Get_allId(ctx context.Context, domain_uuid string) (*models.Response, error) {

	return r.repository.Get_allId(ctx, domain_uuid)
}

/**************************************Get customer by appUserId ************************************************/
func (r *crudUsecase) Get_Customer_by_agent_uuid(ctx context.Context, agent_uuid string, customer_id string) (*models.Response, error) {

	return r.repository.Get_Customer_by_agent_uuid(ctx, agent_uuid, customer_id)
}

/**************************************************Delete User***************************************************/

func (r *crudUsecase) App_user(ctx context.Context, body []byte) (*models.Response, error) {

	return r.repository.App_user(ctx, body)
}

/*****************************************************Get By Id******************************************************/

func (r *crudUsecase) Pre_createUser(ctx context.Context, appId string, flow map[string]interface{}) (*models.Response, error) {
	Id := fmt.Sprintf("%v", flow["id"])
	id, _ := strconv.ParseInt(Id, 10, 64)
	userId := fmt.Sprintf("%v", flow["userId"])
	surname := fmt.Sprintf("%v", flow["surname"])
	givenName := fmt.Sprintf("%v", flow["givenName"])
	return r.repository.Pre_createUser(ctx, appId, id, userId, surname, givenName)
}

/*******************************************List integration***********************************************/
func (r *crudUsecase) List_integration(ctx context.Context, appId string) ([]byte, error) {

	return r.repository.List_integration(ctx, appId)
}

/****************************************************Update_AppUser*************************************************/
func (r *crudUsecase) Update_AppUser(ctx context.Context, appUserId string, appId string, flow map[string]interface{}) (*models.Response, error) {
	surname := fmt.Sprintf("%v", flow["surname"])
	givenName := fmt.Sprintf("%v", flow["givenName"])
	return r.repository.Update_AppUser(ctx, appUserId, appId, surname, givenName)
}

/***************************************Add SmoochConfiguration**************************************/
func (r *crudUsecase) Add_Smooch_configuration(ctx context.Context, flow map[string]interface{}) (*models.Response, error) {
	name := fmt.Sprintf("%v", flow["configuration_name"])
	domain_uuid := fmt.Sprintf("%v", flow["domain_uuid"])
	appId := fmt.Sprintf("%v", flow["appId"])
	appKey := fmt.Sprintf("%v", flow["appKey"])
	appSecret := fmt.Sprintf("%v", flow["appSecret"])

	return r.repository.Add_Smooch_configuration(ctx, name, domain_uuid, appId, appKey, appSecret)
}

/***************************************Update SmoochConfiguration*************************************/
func (r *crudUsecase) Update_Smooch_configuration(ctx context.Context, id int64, domain_uuid string, flow map[string]interface{}) (*models.Response, error) {
	name := fmt.Sprintf("%v", flow["configuration_name"])
	appId := fmt.Sprintf("%v", flow["appId"])
	appKey := fmt.Sprintf("%v", flow["appKey"])
	appSecret := fmt.Sprintf("%v", flow["appSecret"])

	return r.repository.Update_Smooch_configuration(ctx, id, name, domain_uuid, appId, appKey, appSecret)
}

/**************************************Delete smooch configuration************************************/
func (r *crudUsecase) Delete_Smooch_configuration(ctx context.Context, id int64, domain_uuid string) (*models.Response, error) {

	return r.repository.Delete_Smooch_configuration(ctx, id, domain_uuid)
}

/***************************************Get Smooch configuration****************************************/
func (r *crudUsecase) Get_Smooch_configuration(ctx context.Context, domain_uuid string) (*models.Response, error) {

	return r.repository.Get_Smooch_configuration(ctx, domain_uuid)
}

/***************************************save tenant details*******************************************/
func (r *crudUsecase) Add_Whatsapp_configuration(ctx context.Context, td models.WhatsappConfigurations) (*models.Response, error) {

	return r.repository.Add_Whatsapp_configuration(ctx, td)
}

/*******************************************Get Tenant AppId*********************************************/
func (r *crudUsecase) Get_Whatsapp_configuration(ctx context.Context, domain_uuid string) (*models.Response, error) {

	return r.repository.Get_Whatsapp_configuration(ctx, domain_uuid)
}

/***************************************Update_tenant_details********************************************/
func (r *crudUsecase) Update_Whatsapp_configuration(ctx context.Context, id int64, domain_uuid string, td models.WhatsappConfigurations) (*models.Response, error) {

	return r.repository.Update_Whatsapp_configuration(ctx, id, domain_uuid, td)
}

/***************************************Delete Tenant details******************************************/
func (r *crudUsecase) Delete_Whatsapp_configuration(ctx context.Context, id int64, domain_uuid string) (*models.Response, error) {

	return r.repository.Delete_Whatsapp_configuration(ctx, id, domain_uuid)
}

/***************************************Add facebook configuration************************************/
func (r *crudUsecase) Add_Facebook_configuration(ctx context.Context, td models.FacebookConfigurations) (*models.Response, error) {

	return r.repository.Add_Facebook_configuration(ctx, td)
}

/*******************************************Get Tenant AppId*********************************************/
func (r *crudUsecase) Get_Facebook_configuration(ctx context.Context, domain_uuid string) (*models.Response, error) {

	return r.repository.Get_Facebook_configuration(ctx, domain_uuid)
}

/***************************************Update_tenant_details********************************************/
func (r *crudUsecase) Update_Facebook_configuration(ctx context.Context, id int64, domain_uuid string, td models.FacebookConfigurations) (*models.Response, error) {

	return r.repository.Update_Facebook_configuration(ctx, id, domain_uuid, td)
}

/***************************************Delete Tenant details******************************************/
func (r *crudUsecase) Delete_Facebook_configuration(ctx context.Context, id int64, domain_uuid string) (*models.Response, error) {

	return r.repository.Delete_Facebook_configuration(ctx, id, domain_uuid)
}

/***************************************save tenant details*******************************************/
func (r *crudUsecase) Add_Twitter_configuration(ctx context.Context, td models.TwitterConfigurations) (*models.Response, error) {

	return r.repository.Add_Twitter_configuration(ctx, td)
}

/*******************************************Get Tenant AppId*********************************************/
func (r *crudUsecase) Get_Twitter_configuration(ctx context.Context, domain_uuid string) (*models.Response, error) {

	return r.repository.Get_Twitter_configuration(ctx, domain_uuid)
}

/***************************************Update_tenant_details********************************************/
func (r *crudUsecase) Update_Twitter_configuration(ctx context.Context, id int64, domain_uuid string, td models.TwitterConfigurations) (*models.Response, error) {

	return r.repository.Update_Twitter_configuration(ctx, id, domain_uuid, td)
}

/***************************************Delete Tenant details******************************************/
func (r *crudUsecase) Delete_Twitter_configuration(ctx context.Context, id int64, domain_uuid string) (*models.Response, error) {

	return r.repository.Delete_Twitter_configuration(ctx, id, domain_uuid)
}

/***************************************Delete All Message*******************************************/
func (r *crudUsecase) DeleteAllMessage(ctx context.Context, appUserId string, appId string) (*models.Response, error) {

	return r.repository.DeleteAllMessage(ctx, appUserId, appId)

}

/**************************************Delete Message by message id**********************************/
func (r *crudUsecase) DeleteMessage(ctx context.Context, appId string, appUserId string, messageId string) (*models.Response, error) {

	return r.repository.DeleteMessage(ctx, appId, appUserId, messageId)

}

/************************************************Post message******************************************/
func (r *crudUsecase) PostMessage(ctx context.Context, appId string, ConversationId string, p models.User) ([]byte, error) {

	return r.repository.PostMessage(ctx, appId, ConversationId, p)

}

/********************************************Create Text Template**************************************/
func (r *crudUsecase) Create_Text_Template(ctx context.Context, appId string, p models.Comtemplate) ([]byte, error) {

	return r.repository.Create_Text_Template(ctx, appId, p)

}

/*******************************************Create Carousel template***********************************/
func (r *crudUsecase) Create_Carousel_Template(ctx context.Context, appId string, p models.Payload) ([]byte, error) {

	return r.repository.Create_Carousel_Template(ctx, appId, p)

}

/********************************************Create Compound template*****************************************/
func (r *crudUsecase) Create_Compound_Template(ctx context.Context, appId string, p models.Comtemplate) ([]byte, error) {

	return r.repository.Create_Compound_Template(ctx, appId, p)

}

/************************************************Create Quickreply Template**********************************/
func (r *crudUsecase) Create_Quickreply_Template(ctx context.Context, appId string, p models.Comtemplate) ([]byte, error) {

	return r.repository.Create_Quickreply_Template(ctx, appId, p)

}

/*******************************************create location Template**************************************/
func (r *crudUsecase) Create_Request_Location_Template(ctx context.Context, appId string, p models.Comtemplate) ([]byte, error) {

	return r.repository.Create_Request_Location_Template(ctx, appId, p)

}

/*************************************************Update text template***********************************/
func (r *crudUsecase) Update_Text_Template(ctx context.Context, appId string, template_id string, p models.Payload) ([]byte, error) {

	return r.repository.Update_Text_Template(ctx, appId, template_id, p)

}

/******************************************************Get template Id***************************************/
func (r *crudUsecase) Get_template(ctx context.Context, appId string, template_id string) ([]byte, error) {

	return r.repository.Get_template(ctx, appId, template_id)

}

/******************************************************List template**************************************/
func (r *crudUsecase) List_template(ctx context.Context, appId string) ([]byte, error) {

	return r.repository.List_template(ctx, appId)

}

/*****************************************************delete template**********************************/
func (r *crudUsecase) Delete_template(ctx context.Context, appId string, template_id string) (*models.Response, error) {

	return r.repository.Delete_template(ctx, appId, template_id)

}

/**************************************************Send Location******************************************/
func (r *crudUsecase) Send_Location(ctx context.Context, appId string, appUserId string, p models.Locations) ([]byte, error) {

	return r.repository.Send_Location(ctx, appId, appUserId, p)

}

/**************************************************Send Message action**********************************/
func (r *crudUsecase) Message_Action_Types(ctx context.Context, appId string, appUserId string, p models.User) ([]byte, error) {

	return r.repository.Message_Action_Types(ctx, appId, appUserId, p)

}

/*******************************************Quickreply message*********************************************/
func (r *crudUsecase) Quickreply_Message(ctx context.Context, appId string, appUserId string, p models.User) ([]byte, error) {

	return r.repository.Quickreply_Message(ctx, appId, appUserId, p)

}

/************************************Send Carousel Message******************************************/
func (r *crudUsecase) Send_Carousel_Message(ctx context.Context, appId string, appUserId string, p models.User) ([]byte, error) {

	return r.repository.Send_Carousel_Message(ctx, appId, appUserId, p)

}

/*************************************Link appuser to channel*****************************************/
func (r *crudUsecase) Link_appUser_to_Channel(ctx context.Context, appId string, appUserId string, p models.Link) ([]byte, error) {

	return r.repository.Link_appUser_to_Channel(ctx, appId, appUserId, p)

}

/***************************************Unlink appuser to Channel**************************************/
func (r *crudUsecase) Unlink_appUser_to_Channel(ctx context.Context, appId string, appUserId string, channel string) ([]byte, error) {

	return r.repository.Unlink_appUser_to_Channel(ctx, appId, appUserId, channel)
}

/*************************************************Upload Attachments**********************************/
func (r *crudUsecase) Upload_Attachments(ctx context.Context, displayName string, AvatarURL string, appId string, conversationId string, Type string, Text string, IntegrationID string, Size int64, file multipart.File, handler *multipart.FileHeader) (*models.Response, error) {

	return r.repository.Upload_Attachments(ctx, displayName, AvatarURL, appId, conversationId, Type, Text, IntegrationID, Size, file, handler)
}

/***********************************TypingActivity**************************************************/
func (r *crudUsecase) TypingActivity(ctx context.Context, appId string, appUserId string, p models.User) ([]byte, error) {

	return r.repository.TypingActivity(ctx, appId, appUserId, p)

}

/**********************************************Disable_AppUser*************************************************/
func (r *crudUsecase) Disable_AppUser(ctx context.Context, appUserId string, domain_uuid string) (*models.Response, error) {

	return r.repository.Disable_AppUser(ctx, appUserId, domain_uuid)

}

/********************************************Reset Unread Count************************************************/
func (r *crudUsecase) Reset_Unread_Count(ctx context.Context, appId string, appUserId string) (*models.Response, error) {

	return r.repository.Reset_Unread_Count(ctx, appId, appUserId)
}

/**********************************************create Queue*************************************************/
func (r *crudUsecase) Create_Queue(ctx context.Context, flow map[string]interface{}) (*models.Response, error) {
	id1 := fmt.Sprintf("%v", flow["id"])
	Id, _ := strconv.ParseInt(id1, 10, 64)
	Queue_uuid := fmt.Sprintf("%v", flow["queue_uuid"])
	Map_with := fmt.Sprintf("%v", flow["map_with"])
	Name := fmt.Sprintf("%v", flow["name"])
	IntegrationID := fmt.Sprintf("%v", flow["integration_id"])
	Domain_uuid := fmt.Sprintf("%v", flow["domain_uuid"])

	return r.repository.Create_Queue(ctx, Id, Queue_uuid, Map_with, Name, IntegrationID, Domain_uuid)
}

/***************************************************Assign_Agent************************************************/
func (r *crudUsecase) Assign_Agent_To_Queue(ctx context.Context, flow map[string]interface{}) (*models.Response, error) {

	Agent_name := fmt.Sprintf("%v", flow["agent_name"])
	Agent_uuid := fmt.Sprintf("%v", flow["agent_uuid"])
	Queue_name := fmt.Sprintf("%v", flow["queue_name"])
	Tenant_domain_uuid := fmt.Sprintf("%v", flow["tenant_domain_uuid"])
	Queue_uuid := fmt.Sprintf("%v", flow["queue_uuid"])
	return r.repository.Assign_Agent_To_Queue(ctx, Agent_name, Agent_uuid, Queue_name, Tenant_domain_uuid, Queue_uuid)
}

/*************************************************Remove Agent From Queue****************************************/
func (r *crudUsecase) Remove_Agent_From_Queue(ctx context.Context, agent_uuid string) (*models.Response, error) {

	return r.repository.Remove_Agent_From_Queue(ctx, agent_uuid)
}

/**********************************************Get Assigned Agent list from Queue******************************/
func (r *crudUsecase) Get_Assigned_Agent_list_From_Queue(ctx context.Context, queue_uuid string) (*models.Response, error) {

	return r.repository.Get_Assigned_Agent_list_From_Queue(ctx, queue_uuid)
}

/***************************************Get Queue List****************************************************/
func (r *crudUsecase) Get_Queue_List(ctx context.Context, domain_uuid string) (*models.Response, error) {

	return r.repository.Get_Queue_List(ctx, domain_uuid)

}

/************************************Get Available Agents List In Queue**********************************/
func (r *crudUsecase) Get_Available_Agents_Queue_List(ctx context.Context, agent_uuid string, queue_uuid string) (*models.Response, error) {

	return r.repository.Get_Available_Agents_Queue_List(ctx, agent_uuid, queue_uuid)
}

/**********************************************Update_Queue*****************************************************/
func (r *crudUsecase) Update_Queue(ctx context.Context, queue_uuid string, flow map[string]interface{}) (*models.Response, error) {

	Name := fmt.Sprintf("%v", flow["name"])
	IntegrationID := fmt.Sprintf("%v", flow["integration_id"])
	Map_with := fmt.Sprintf("%v", flow["map_with"])
	Domain_uuid := fmt.Sprintf("%v", flow["domain_uuid"])
	return r.repository.Update_Queue(ctx, queue_uuid, Name, IntegrationID, Map_with, Domain_uuid)
}

/********************************************Delete Queue******************************************************/

func (r *crudUsecase) Delete_Queue(ctx context.Context, queue_uuid string) (*models.Response, error) {

	return r.repository.Delete_Queue(ctx, queue_uuid)

}

/**********************************************Available Agnets***********************************************/
func (r *crudUsecase) Available_Agents(ctx context.Context, domain_uuid string, queue_uuid string) (*models.Response, error) {

	return r.repository.Available_Agents(ctx, domain_uuid, queue_uuid)

}

/**********************************************Transfer customer*********************************************/
func (r *crudUsecase) Transfer_customer(ctx context.Context, flow map[string]interface{}) (*models.Response, error) {
	agent_name := fmt.Sprintf("%v", flow["agent_name"])
	conversation_id := fmt.Sprintf("%v", flow["conversation_id"])
	agent_uuid := fmt.Sprintf("%v", flow["agent_uuid"])
	appUserId := fmt.Sprintf("%v", flow["appUserId"])

	return r.repository.Transfer_customer(ctx, agent_name, conversation_id, agent_uuid, appUserId)
}

/************************************************Post on FB page**********************************************/
func (r *crudUsecase) Publish_Post_on_FB_Page(ctx context.Context, flow map[string]interface{}) ([]byte, error) {
	pageId := fmt.Sprintf("%v", flow["page_id"])
	message := fmt.Sprintf("%v", flow["message"])
	access_token := fmt.Sprintf("%v", flow["access_token"])
	// Post_type := fmt.Sprintf("%v", flow["post_type"])
	return r.repository.Publish_Post_on_FB_Page(ctx, pageId, message, access_token)

}

/*************************************************Get all Post of a Page***********************************/
func (r *crudUsecase) Getall_Post_of_Page(ctx context.Context, pageId string, access_token string) ([]byte, error) {

	return r.repository.Getall_Post_of_Page(ctx, pageId, access_token)
}

/*************************************************Delete Post of a page***************************************/
func (r *crudUsecase) Delete_Post_of_Page(ctx context.Context, page_postId string, access_token string) ([]byte, error) {

	return r.repository.Delete_Post_of_Page(ctx, page_postId, access_token)
}

/*************************************************Update Post of a Page***************************************/
func (r *crudUsecase) Update_Post_of_Page(ctx context.Context, flow map[string]interface{}) ([]byte, error) {
	page_postId := fmt.Sprintf("%v", flow["page_post_id"])
	message := fmt.Sprintf("%v", flow["message"])
	access_token := fmt.Sprintf("%v", flow["access_token"])
	return r.repository.Update_Post_of_Page(ctx, page_postId, message, access_token)
}

/******************************************Get comments of a page***********************************************/
func (r *crudUsecase) Get_Comments_on_Post_of_Page(ctx context.Context, page_postId string, access_token string) ([]byte, error) {

	return r.repository.Get_Comments_on_Post_of_Page(ctx, page_postId, access_token)
}

/*******************************************Get Likes on a Page*************************************************/
func (r *crudUsecase) Get_Likes_on_Post_of_Page(ctx context.Context, page_postId string, access_token string) ([]byte, error) {

	return r.repository.Get_Likes_on_Post_of_Page(ctx, page_postId, access_token)
}

/*******************************************Comment on post of a page******************************************/
func (r *crudUsecase) Comment_on_Post_of_Page(ctx context.Context, flow map[string]interface{}) ([]byte, error) {
	page_postId := fmt.Sprintf("%v", flow["page_post_id"])
	message := fmt.Sprintf("%v", flow["message"])
	access_token := fmt.Sprintf("%v", flow["access_token"])
	return r.repository.Comment_on_Post_of_Page(ctx, page_postId, message, access_token)
}

/********************************************Get Page Id******************************************************/
func (r *crudUsecase) Get_Page_ID(ctx context.Context, access_token string) ([]byte, error) {

	return r.repository.Get_Page_ID(ctx, access_token)
}

/*******************************************Schedule a Post*********************************************/
func (r *crudUsecase) Schedule_Post(ctx context.Context, flow map[string]interface{}) ([]byte, error) {
	pageId := fmt.Sprintf("%v", flow["page_id"])
	message := fmt.Sprintf("%v", flow["message"])
	scheduled_publish_time := fmt.Sprintf("%v", flow["schedule_publish_time"])
	access_token := fmt.Sprintf("%v", flow["access_token"])
	return r.repository.Schedule_Post(ctx, pageId, message, scheduled_publish_time, access_token)
}

/********************************************Publish link with message************************************/
func (r *crudUsecase) Publish_link_with_message_on_Post(ctx context.Context, flow map[string]interface{}) ([]byte, error) {
	pageId := fmt.Sprintf("%v", flow["page_id"])
	message := fmt.Sprintf("%v", flow["message"])
	link := fmt.Sprintf("%v", flow["link"])
	access_token := fmt.Sprintf("%v", flow["access_token"])
	return r.repository.Publish_link_with_message_on_Post(ctx, pageId, message, link, access_token)
}

/*******************************************Upload Photo with message on post**********************************/

func (r *crudUsecase) Upload_Photo_on_Post(ctx context.Context, pageId string, access_token string, message string, Type string, file multipart.File, handler *multipart.FileHeader) ([]byte, error) {

	return r.repository.Upload_Photo_on_Post(ctx, pageId, access_token, message, Type, file, handler)
}

/********************************************Facebook Login Api**********************************************/
func (r *crudUsecase) UVoiceFacebookLogin(ctx context.Context, c echo.Context, client_id string, client_secret string, flac_uuid string) (*models.Response, error) {
	return r.repository.UVoiceFacebookLogin(ctx, c, client_id, client_secret, flac_uuid)
}

/************************************************Facebook login callback***************************************/
func (r *crudUsecase) UVoiceFacebookLoginCallback(ctx context.Context, c echo.Context) (*models.Response, error) {
	return r.repository.UVoiceFacebookLoginCallback(ctx, c)
}

/**************************************************Add Facebook Application************************************/

func (r *crudUsecase) AddFacebookApplication(ctx context.Context, domain_uuid string, app_id string, app_secret string, app_name string) (*models.Response, error) {
	return r.repository.AddFacebookApplication(ctx, domain_uuid, app_id, app_secret, app_name)
}

/*************************************************Show facebook Application**************************************/
func (r *crudUsecase) ShowFacebookApplication(ctx context.Context, domain_uuid string) (*models.Response, error) {
	return r.repository.ShowFacebookApplication(ctx, domain_uuid)
}

/*************************************************Delete Facebook Application************************************/
func (r *crudUsecase) DeleteFacebookApplication(ctx context.Context, domain_uuid string, flac_uuid string) (*models.Response, error) {
	return r.repository.DeleteFacebookApplication(ctx, domain_uuid, flac_uuid)
}

/*************************************************Update Facebook Application*************************************/
func (r *crudUsecase) UpdateFacebookApplication(ctx context.Context, domain_uuid string, flac_uuid string, fb map[string]interface{}) (*models.Response, error) {
	app_id := fmt.Sprintf("%v", fb["app_id"])
	app_secret := fmt.Sprintf("%v", fb["app_secret"])
	app_name := fmt.Sprintf("%v", fb["app_name"])
	return r.repository.UpdateFacebookApplication(ctx, domain_uuid, flac_uuid, app_id, app_secret, app_name)
}

/*************************************************Assign Agent to Facebook Application******************************/
func (r *crudUsecase) AssignAgentToFacebookApplication(ctx context.Context, domain_uuid string, flac_uuid string, agent_uuid string) (*models.Response, error) {
	return r.repository.AssignAgentToFacebookApplication(ctx, domain_uuid, flac_uuid, agent_uuid)
}

/*************************************************Assign Agent list in Facebook Application******************************/
func (r *crudUsecase) AgentListAssignedToFacebookApplication(ctx context.Context, flac_uuid string) (*models.Response, error) {
	return r.repository.AgentListAssignedToFacebookApplication(ctx, flac_uuid)
}

/*****************************************Not Assigned Agent list in Facebook Application******************************/
func (r *crudUsecase) AgentListNotInFacebookApplication(ctx context.Context, flac_uuid string, domain_uuid string) (*models.Response, error) {
	return r.repository.AgentListNotInFacebookApplication(ctx, flac_uuid, domain_uuid)
}

/********************************************Show Assigned Agent of Facebook Application******************************/
func (r *crudUsecase) ShowAgentFacebookApplication(ctx context.Context, agent_uuid string) (*models.Response, error) {
	return r.repository.ShowAgentFacebookApplication(ctx, agent_uuid)
}

/********************************************Convert access token into longlived token******************************/
func (r *crudUsecase) Convert_Access_Token_into_Longlived_Token(ctx context.Context, clientId string, clientSecret string, exchange_token string, access_token string) ([]byte, error) {

	return r.repository.Convert_Access_Token_into_Longlived_Token(ctx, clientId, clientSecret, exchange_token, access_token)
}

/*********************************************Remove Assigned Agent from Facebook Application*************************/
func (r *crudUsecase) RemoveAgentAssignedToFacebookApplication(ctx context.Context, agent_uuid string) (*models.Response, error) {

	return r.repository.RemoveAgentAssignedToFacebookApplication(ctx, agent_uuid)
}

/********************************************Send Private Message***************************************************/
func (r *crudUsecase) Send_Private_Message(ctx context.Context, flow map[string]interface{}) ([]byte, error) {
	pageId := fmt.Sprintf("%v", flow["page_id"])
	postId := fmt.Sprintf("%v", flow["post_id"])
	message := fmt.Sprintf("%v", flow["message"])
	access_token := fmt.Sprintf("%v", flow["access_token"])
	return r.repository.Send_Private_Message(ctx, pageId, postId, message, access_token)
}

/*******************************************Likes and unlike post and comments***********************************/
func (r *crudUsecase) Like_and_Unlike_Post_and_Comment(ctx context.Context, flow map[string]interface{}) ([]byte, error) {
	postId := fmt.Sprintf("%v", flow["post_id"])
	commentId := fmt.Sprintf("%v", flow["comment_id"])
	access_token := fmt.Sprintf("%v", flow["access_token"])
	Type := fmt.Sprintf("%v", flow["type"])
	return r.repository.Like_and_Unlike_Post_and_Comment(ctx, postId, commentId, access_token, Type)
}

/*************************************************Delete Tickets***********************************************/
func (r *crudUsecase) Delete_Tickets(ctx context.Context, ticket_uuid string) (*models.Response, error) {

	return r.repository.Delete_Tickets(ctx, ticket_uuid)
}

/***************************************************Get All Tickets******************************************/
func (r *crudUsecase) GetAll_Tickets(ctx context.Context, domain_uuid string) (*models.Response, error) {

	return r.repository.GetAll_Tickets(ctx, domain_uuid)
}

/*************************************************Save Twitter Auth***************************************/
func (r *crudUsecase) SaveTwitterAuth(ctx context.Context, flow map[string]interface{}) (*models.Response, error) {
	Id := fmt.Sprintf("%v", flow["id"])
	id, _ := strconv.ParseInt(Id, 10, 64)
	domain_uuid := fmt.Sprintf("%v", flow["domain_uuid"])
	api_key := fmt.Sprintf("%v", flow["api_key"])
	api_secret := fmt.Sprintf("%v", flow["api_secret"])
	bearer_token := fmt.Sprintf("%v", flow["bearer_token"])
	access_token := fmt.Sprintf("%v", flow["access_token"])
	token_secret := fmt.Sprintf("%v", flow["token_secret"])
	return r.repository.SaveTwitterAuth(ctx, id, domain_uuid, api_key, api_secret, bearer_token, access_token, token_secret)
}

/*************************************************Update Twitter Auth************************************/
func (r *crudUsecase) UpdateTwitterAuth(ctx context.Context, id int64, domain_uuid string, flow map[string]interface{}) (*models.Response, error) {
	api_key := fmt.Sprintf("%v", flow["api_key"])
	api_secret := fmt.Sprintf("%v", flow["api_secret"])
	bearer_token := fmt.Sprintf("%v", flow["bearer_token"])
	access_token := fmt.Sprintf("%v", flow["access_token"])
	token_secret := fmt.Sprintf("%v", flow["token_secret"])
	return r.repository.UpdateTwitterAuth(ctx, id, domain_uuid, api_key, api_secret, bearer_token, access_token, token_secret)
}

/***********************************************Get Twitter Auth*****************************************/
func (r *crudUsecase) GetTwitterAuth(ctx context.Context, domain_uuid string) (*models.Response, error) {

	return r.repository.GetTwitterAuth(ctx, domain_uuid)
}

/***********************************************Delete Twitter Auth**************************************/
func (r *crudUsecase) DeleteTwitterAuth(ctx context.Context, id int64, domain_uuid string) (*models.Response, error) {

	return r.repository.DeleteTwitterAuth(ctx, id, domain_uuid)
}

/**********************************************Twitter Timeline*******************************************/
func (r *crudUsecase) Twitter_Apis(ctx context.Context, flow map[string]interface{}) ([]byte, error) {
	tweet_id := fmt.Sprintf("%v", flow["tweet_id"])
	screen_name := fmt.Sprintf("%v", flow["screen_name"])
	api_key := fmt.Sprintf("%v", flow["api_key"])
	api_type := fmt.Sprintf("%v", flow["api_type"])
	author_id := fmt.Sprintf("%v", flow["author_id"])
	message := fmt.Sprintf("%v", flow["message"])
	return r.repository.Twitter_Apis(ctx, tweet_id, screen_name, api_key, api_type, author_id, message)
}

/**********************************************Assign agent to twitter*************************************/
func (r *crudUsecase) AssignAgentToTwitter(ctx context.Context, flow map[string]interface{}) (*models.Response, error) {
	twitter_uuid := fmt.Sprintf("%v", flow["twitter_uuid"])
	domain_uuid := fmt.Sprintf("%v", flow["domain_uuid"])
	api_key := fmt.Sprintf("%v", flow["api_key"])
	agent_uuid := fmt.Sprintf("%v", flow["agent_uuid"])

	return r.repository.AssignAgentToTwitter(ctx, twitter_uuid, domain_uuid, api_key, agent_uuid)
}

/**************************************Twitter Assigned agent list**************************************/
func (r *crudUsecase) TwitterAssignAgentList(ctx context.Context, domain_uuid string, twitter_uuid string) (*models.Response, error) {

	return r.repository.TwitterAssignAgentList(ctx, domain_uuid, twitter_uuid)
}

/*****************************************Remove Twitter Assigned Agent************************************/

func (r *crudUsecase) RemoveTwitterAssignAgent(ctx context.Context, agent_uuid string, twitter_uuid string) (*models.Response, error) {

	return r.repository.RemoveTwitterAssignAgent(ctx, agent_uuid, twitter_uuid)
}

/****************************************Get Quoted Retweet List****************************************/
func (r *crudUsecase) Get_Quoted_Retweet_List(ctx context.Context, api_key string, tweet_id string) (*models.Response, error) {

	return r.repository.Get_Quoted_Retweet_List(ctx, api_key, tweet_id)
}

/*******************************************Assign Agent To customer**************************************/
func (r *crudUsecase) AssigncustomerToAgent(ctx context.Context, flow map[string]interface{}) (*models.Response, error) {
	domain_uuid := fmt.Sprintf("%v", flow["domain_uuid"])
	agent_uuid := fmt.Sprintf("%v", flow["agent_uuid"])
	app_user_id := fmt.Sprintf("%v", flow["app_user_id"])
	return r.repository.AssigncustomerToAgent(ctx, domain_uuid, agent_uuid, app_user_id)
}

/********************************************Real Time Like And Comments************************************/
func (r *crudUsecase) Webhook_verify(ctx context.Context, mode string, token string, challenge string, body []byte) (string, error) {

	return r.repository.Webhook_verify(ctx, mode, token, challenge, body)
}

/*******************************************Facebook Real time Like And comments********************************/
func (r *crudUsecase) FacebookLikeAndComments(ctx context.Context, body []byte) (*models.Response, error) {

	return r.repository.FacebookLikeAndComments(ctx, body)
}
