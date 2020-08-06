package crud

import (
	"context"
	"mime/multipart"
	models "whatsapp_api/model"
)

type Repository interface {
	Delete_AppUser(ctx context.Context, appUserId string, appId string) (*models.Response, error)
	GetAllMessageByAppUserId(ctx context.Context, appUserId string, appId string) ([]byte, error)
	Get_allId(ctx context.Context) (*models.Response, error)
	App_user(ctx context.Context, body []byte) (*models.Response, error)
	Pre_createUser(ctx context.Context, appId string, id int64, userId string, surname string, givenName string) (*models.Response, error)
	Update_AppUser(tx context.Context, appUserId string, appId string, surname string, givenName string) (*models.Response, error)
	Add_Smooch_configuration(ctx context.Context, domain_uuid string, appId string, appKey string, appSecret string) (*models.Response, error)
	Update_Smooch_configuration(ctx context.Context, id int64, domain_uuid string, appId string, appKey string, appSecret string) (*models.Response, error)
	Delete_Smooch_configuration(ctx context.Context, id int64, domain_uuid string) (*models.Response, error)
	Get_Smooch_configuration(ctx context.Context, domain_uuid string) (*models.Response, error)
	Add_Whatsapp_configuration(ctx context.Context, domain_uuid string, appId string, appKey string, appSecret string, WhatsappIntegrationID string) (*models.Response, error)
	Get_Whatsapp_configuration(ctx context.Context, domain_uuid string) (*models.Response, error)
	Update_Whatsapp_configuration(ctx context.Context, id int64, domain_uuid string, appId string, appKey string, appSecret string, WhatsappIntegrationID string) (*models.Response, error)
	Delete_Whatsapp_configuration(ctx context.Context, id int64, domain_uuid string) (*models.Response, error)
	Add_Facebook_configuration(ctx context.Context, domain_uuid string, appId string, appKey string, appSecret string, FacebookIntegrationId string) (*models.Response, error)
	Get_Facebook_configuration(ctx context.Context, domain_uuid string) (*models.Response, error)
	Update_Facebook_configuration(ctx context.Context, id int64, domain_uuid string, appId string, appKey string, appSecret string, FacebookIntegrationId string) (*models.Response, error)
	Delete_Facebook_configuration(ctx context.Context, id int64, domain_uuid string) (*models.Response, error)
	List_integration(ctx context.Context, appId string) ([]byte, error)
	DeleteAllMessage(ctx context.Context, appUserId string, appId string) (*models.Response, error)
	PostMessage(ctx context.Context, appId string, appUserId string, p models.User) ([]byte, error)
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
	Upload_Attachments(ctx context.Context, appId string, appUserId string, Type string, file multipart.File, handler *multipart.FileHeader) ([]byte, error)
	TypingActivity(ctx context.Context, appId string, appUserId string, p models.User) ([]byte, error)
}
