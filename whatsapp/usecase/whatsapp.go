package usecase

import (
	"context"
	"fmt"
	"mime/multipart"
	"strconv"
	models "whatsapp_api/model"
	crud "whatsapp_api/whatsapp"
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

func (r *crudUsecase) Get_allId(ctx context.Context) (*models.Response, error) {

	return r.repository.Get_allId(ctx)
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
	domain_uuid := fmt.Sprintf("%v", flow["domain_uuid"])
	appId := fmt.Sprintf("%v", flow["appId"])
	appKey := fmt.Sprintf("%v", flow["appKey"])
	appSecret := fmt.Sprintf("%v", flow["appSecret"])

	return r.repository.Add_Smooch_configuration(ctx, domain_uuid, appId, appKey, appSecret)
}

/***************************************Update SmoochConfiguration*************************************/
func (r *crudUsecase) Update_Smooch_configuration(ctx context.Context, id int64, domain_uuid string, flow map[string]interface{}) (*models.Response, error) {
	appId := fmt.Sprintf("%v", flow["appId"])
	appKey := fmt.Sprintf("%v", flow["appKey"])
	appSecret := fmt.Sprintf("%v", flow["appSecret"])

	return r.repository.Update_Smooch_configuration(ctx, id, domain_uuid, appId, appKey, appSecret)
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
	// appId := fmt.Sprintf("%v", flow["appId"])
	// appKey := fmt.Sprintf("%v", flow["appKey"])
	// appSecret := fmt.Sprintf("%v", flow["appSecret"])
	// WhatsappIntegrationID := fmt.Sprintf("%v", flow["whatsapp_integration_id"])
	// WorkingHourStartTime := fmt.Sprintf("%v", flow["WorkingHourStartTime"])
	// WorkingHourEndTime := fmt.Sprintf("%v", flow["WorkingHourEndTime"])
	// workingDays := fmt.Sprintf("%v", flow["working_days"])

	return r.repository.Update_Whatsapp_configuration(ctx, id, domain_uuid, td)
}

/***************************************Delete Tenant details******************************************/
func (r *crudUsecase) Delete_Whatsapp_configuration(ctx context.Context, id int64, domain_uuid string) (*models.Response, error) {

	return r.repository.Delete_Whatsapp_configuration(ctx, id, domain_uuid)
}

/***************************************Add facebook configuration************************************/
func (r *crudUsecase) Add_Facebook_configuration(ctx context.Context, td models.FacebookConfigurations) (*models.Response, error) {
	// domain_uuid := fmt.Sprintf("%v", flow["domain_uuid"])
	// appId := fmt.Sprintf("%v", flow["appId"])
	// appKey := fmt.Sprintf("%v", flow["appKey"])
	// appSecret := fmt.Sprintf("%v", flow["appSecret"])
	// FacebookIntegrationId := fmt.Sprintf("%v", flow["facebook_integration_id"])
	// WorkingHourStartTime := fmt.Sprintf("%v", flow["WorkingHourStartTime"])
	// WorkingHourEndTime := fmt.Sprintf("%v", flow["WorkingHourEndTime"])
	// workingDays := fmt.Sprintf("%v", flow["working_days"])

	return r.repository.Add_Facebook_configuration(ctx, td)
}

/*******************************************Get Tenant AppId*********************************************/
func (r *crudUsecase) Get_Facebook_configuration(ctx context.Context, domain_uuid string) (*models.Response, error) {

	return r.repository.Get_Facebook_configuration(ctx, domain_uuid)
}

/***************************************Update_tenant_details********************************************/
func (r *crudUsecase) Update_Facebook_configuration(ctx context.Context, id int64, domain_uuid string, td models.FacebookConfigurations) (*models.Response, error) {
	// appId := fmt.Sprintf("%v", flow["appId"])
	// appKey := fmt.Sprintf("%v", flow["appKey"])
	// appSecret := fmt.Sprintf("%v", flow["appSecret"])
	// FacebookIntegrationId := fmt.Sprintf("%v", flow["facebook_integration_id"])
	// WorkingHourStartTime := fmt.Sprintf("%v", flow["WorkingHourStartTime"])
	// WorkingHourEndTime := fmt.Sprintf("%v", flow["WorkingHourEndTime"])
	// workingDays := fmt.Sprintf("%v", flow["working_days"])

	return r.repository.Update_Facebook_configuration(ctx, id, domain_uuid, td)
}

/***************************************Delete Tenant details******************************************/
func (r *crudUsecase) Delete_Facebook_configuration(ctx context.Context, id int64, domain_uuid string) (*models.Response, error) {

	return r.repository.Delete_Facebook_configuration(ctx, id, domain_uuid)
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
func (r *crudUsecase) PostMessage(ctx context.Context, appId string, appUserId string, p models.User) ([]byte, error) {

	return r.repository.PostMessage(ctx, appId, appUserId, p)

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
func (r *crudUsecase) Upload_Attachments(ctx context.Context, appId string, appUserId string, Type string, IntegrationID string, Size int64, file multipart.File, handler *multipart.FileHeader) (*models.Response, error) {

	return r.repository.Upload_Attachments(ctx, appId, appUserId, Type, IntegrationID, Size, file, handler)
}

/***********************************TypingActivity**************************************************/
func (r *crudUsecase) TypingActivity(ctx context.Context, appId string, appUserId string, p models.User) ([]byte, error) {

	return r.repository.TypingActivity(ctx, appId, appUserId, p)

}

/**********************************************Disable_AppUser*************************************************/
func (r *crudUsecase) Disable_AppUser(ctx context.Context, appUserId string) (*models.Response, error) {

	return r.repository.Disable_AppUser(ctx, appUserId)

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
