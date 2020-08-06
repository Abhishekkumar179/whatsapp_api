package repository

import (
	"bytes"
	"context"
	json "encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	models "whatsapp_api/model"
	crud "whatsapp_api/whatsapp"

	"github.com/jinzhu/gorm"
)

type crudRepository struct {
	DBConn *gorm.DB
}

func NewcrudRepository(conn *gorm.DB) crud.Repository {
	return &crudRepository{
		DBConn: conn,
	}
}

/******************************************Create_text_template**************************************/
func (r *crudRepository) Delete_AppUser(ctx context.Context, appUserId string, appId string) (*models.Response, error) {

	td := models.Tenant_details{}
	db := r.DBConn.Where("app_id = ?", appId).Find(&td)
	if db.Error != nil {
		return nil, db.Error
	}
	res, err := http.NewRequest("DELETE", "https://api.smooch.io/v1.1/apps/"+appId+"/appusers/"+appUserId, nil)
	res.Header.Set("Content-Type", "application/json")
	res.SetBasicAuth(td.AppKey, td.AppSecret)

	client := &http.Client{}
	response, err := client.Do(res)
	if err != nil {
		fmt.Printf("error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)

		fmt.Println(string(data))

		if string(data) == "{}" {
			us := models.Appuser{}
			re := models.ReceiveUserDetails{}
			st := r.DBConn.Where("id = ?", appUserId).Delete(&us)
			if st.Error != nil {
				return &models.Response{ResponseCode: 409, Status: "Error", Msg: "Not Deleted"}, nil
			}
			db := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", appUserId).Delete(&re)
			if db.Error != nil {
				return &models.Response{ResponseCode: 409, Status: "Error", Msg: "Not Deleted"}, nil
			}
			return &models.Response{ResponseCode: 201, Status: "OK", Msg: "User Deleted successfully."}, nil

		}
	}
	return &models.Response{ResponseCode: 404, Status: "Failed", Msg: "User NOT Found."}, nil

}

/**************************************************Getall Id***************************************************/

func (r *crudRepository) Get_allId(ctx context.Context) (*models.Response, error) {
	list := make([]models.ReceiveUserDetails, 0)

	if rows, err := r.DBConn.Raw("select app_id, app_user_id, surname, given_name,type,text,role,name,author_id,message_id,original_message_id,integration_id,source_type, signed_up_at, conversation_started from receive_user_details").Rows(); err != nil {

		return &models.Response{Status: "Not Found", Msg: "Record Not Found", ResponseCode: 204}, nil
	} else {
		defer rows.Close()
		for rows.Next() {
			f := models.ReceiveUserDetails{}
			if err := rows.Scan(&f.AppId, &f.AppUserId, &f.Surname, &f.GivenName, &f.Type, &f.Text, &f.Role, &f.Name, &f.AuthorID, &f.Message_id, &f.OriginalMessageID, &f.IntegrationID, &f.Source_Type, &f.SignedUpAt, &f.ConversationStarted); err != nil {

				return nil, err
			}

			list = append(list, f)
		}
		return &models.Response{Status: "OK", Msg: "Record Found", ResponseCode: 200, AppUserList: list}, nil
	}
}

/**************************************************Getall_messageByUserId***************************************************/

func (r *crudRepository) GetAllMessageByAppUserId(ctx context.Context, appUserId string, appId string) ([]byte, error) {
	td := models.Tenant_details{}
	db := r.DBConn.Where("app_id = ?", appId).Find(&td)
	if db.Error != nil {

	}
	res, err := http.NewRequest("GET", "https://api.smooch.io/v1.1/apps/"+appId+"/appusers/"+appUserId+"/messages", nil)
	res.Header.Set("Content-Type", "application/json")
	res.SetBasicAuth(td.AppKey, td.AppSecret)
	client := &http.Client{}
	response, err := client.Do(res)
	if err != nil {
		fmt.Printf("error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data), "enterrer")
		return data, nil
	}
	return nil, err
}

/**************************************************App User***************************************************/

func (r *crudRepository) App_user(ctx context.Context, body []byte) (*models.Response, error) {
	f := models.Received{}
	jsondata := json.Unmarshal(body, &f)
	fmt.Println(jsondata, f)

	u := models.ReceiveUserDetails{
		Trigger:                  f.Trigger,
		Version:                  f.Version,
		AppId:                    f.App.ID,
		AppUserId:                f.AppUser.ID,
		Surname:                  f.AppUser.Surname,
		GivenName:                f.AppUser.GivenName,
		SignedUpAt:               f.AppUser.SignedUpAt,
		ConversationStarted:      f.AppUser.ConversationStarted,
		Conversation_id:          f.Conversation.ID,
		Type:                     f.Messages[0].Type,
		Text:                     f.Messages[0].Text,
		Role:                     f.Messages[0].Role,
		Received:                 f.Messages[0].Received,
		Name:                     f.Messages[0].Name,
		AuthorID:                 f.Messages[0].AuthorID,
		Message_id:               f.Messages[0].ID,
		OriginalMessageID:        f.Messages[0].Source.OriginalMessageID,
		OriginalMessageTimestamp: f.Messages[0].Source.OriginalMessageTimestamp,
		Source_Type:              f.Messages[0].Source.Type,
		IntegrationID:            f.Messages[0].Source.IntegrationID,
	}
	if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
		db := r.DBConn.Create(&u)
		if db.Error != nil {

		}
		td := models.Tenant_details{}
		DB := r.DBConn.Where("whatsapp_integration_id = ?", f.Messages[0].Source.IntegrationID).Or("fb_integration_id = ?", f.Messages[0].Source.IntegrationID).Find(&td)
		fmt.Println(DB)
		p := models.User{
			Role: "appMaker",
			Type: "text",
			Text: "How can i help you.",
		}
		r.PostMessage(ctx, td.AppId, f.AppUser.ID, p)
		return &models.Response{Received: &f}, nil
	}
	fmt.Println("appUserId already exist.")
	return &models.Response{Msg: "Userid already exist."}, nil
}

/**************************************************Get By Id****************************************************/

func (r *crudRepository) Pre_createUser(ctx context.Context, appId string, id int64, userId string, surname string, givenName string) (*models.Response, error) {
	td := models.Tenant_details{}
	db := r.DBConn.Where("app_id = ?", appId).Find(&td)
	if db.Error != nil {
		return nil, db.Error
	}

	jsonData := map[string]string{"userId": userId, "surname": surname, "givenName": givenName}
	jsonValue, _ := json.MarshalIndent(jsonData, "", " ")
	req, _ := http.NewRequest("POST", "https://api.smooch.io/v1.1/apps/"+appId+"/appusers", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(td.AppKey, td.AppSecret)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	} else {

		var ds models.Data
		data, _ := ioutil.ReadAll(res.Body)
		jsondata := json.Unmarshal(data, &ds)
		fmt.Println(jsondata)
		us := models.Appuser{

			Id:                  ds.AppUser.Id,
			UserId:              ds.AppUser.UserId,
			Surname:             ds.AppUser.Surname,
			GivenName:           ds.AppUser.GivenName,
			Properties:          ds.AppUser.Properties,
			PendingClients:      ds.AppUser.PendingClients,
			ConversationStarted: ds.AppUser.ConversationStarted,
			Clients:             ds.AppUser.Clients,
			HasPaymentInfo:      ds.AppUser.HasPaymentInfo,
			SignedUpAt:          ds.AppUser.SignedUpAt,
		}

		if db := r.DBConn.Table("appusers").Where("user_id = ?", userId).Find(&us).Error; db != nil {
			st := r.DBConn.Create(&us)
			if st.Error != nil {
				return &models.Response{ResponseCode: 409, Status: "Error", Msg: "Not created"}, nil
			}
			return &models.Response{ResponseCode: 201, Status: "OK", Msg: "User registered successfully.", Appuser: &ds}, nil

		}
		return &models.Response{ResponseCode: 409, Status: "Error", Msg: "UserId Already Exist."}, nil
	}
}

/***********************************************update appuser******************************************/
func (r crudRepository) Update_AppUser(ctx context.Context, appUserId string, appId string, surname string, givenName string) (*models.Response, error) {
	td := models.Tenant_details{}
	db := r.DBConn.Where("app_id = ?", appId).Find(&td)
	if db.Error != nil {
		return nil, db.Error
	}
	jsonData := map[string]string{"surname": surname, "givenName": givenName}
	jsonValue, _ := json.MarshalIndent(jsonData, "", " ")

	req, _ := http.NewRequest("PUT", "https://api.smooch.io/v1.1/apps/"+appId+"/appusers"+appUserId, bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(td.AppKey, td.AppSecret)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	} else {
		var ds models.Data
		data, _ := ioutil.ReadAll(res.Body)
		jsondata := json.Unmarshal(data, &ds)
		fmt.Println(jsondata)

		if db := r.DBConn.Table("appusers").Where("id = ?", appUserId).Updates(map[string]string{"surname": surname, "given_name": givenName}).Error; db != nil {

			return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
		}
		if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", appUserId).Updates(map[string]interface{}{"surname": surname, "given_name": givenName}).Error; err != nil {
			return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
		}

		return &models.Response{ResponseCode: 201, Status: "OK", Msg: "User Updated successfully.", Appuser: &ds}, nil

	}
}

/****************************************SmoochConfiguration*********************************************/
func (r crudRepository) Add_Smooch_configuration(ctx context.Context, domain_uuid string, appId string, appKey string, appSecret string) (*models.Response, error) {
	td := models.Tenant_details{
		Domain_uuid: domain_uuid,
		AppId:       appId,
		AppKey:      appKey,
		AppSecret:   appSecret,
	}
	if db := r.DBConn.Table("tenant_details").Where("app_id = ?", appId).Find(&td).Error; db != nil {
		st := r.DBConn.Create(&td)
		if st.Error != nil {
			return &models.Response{ResponseCode: 409, Status: "Error", Msg: "Not created"}, nil
		}
		return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Smooch configuration saved successfully."}, nil

	}
	return &models.Response{ResponseCode: 409, Status: "Error", Msg: "AppId  Already Exist."}, nil
}

/***************************************Add smooch configuration*****************************************/
func (r crudRepository) Update_Smooch_configuration(ctx context.Context, id int64, domain_uuid string, appId string, appKey string, appSecret string) (*models.Response, error) {
	w := models.Tenant_details{}
	db := r.DBConn.Where("domain_uuid = ? AND id = ?", domain_uuid, id).Find(&w)
	if db.RowsAffected == 0 {
		return &models.Response{Status: "Not Found", Msg: "Record Doesn't Exist", ResponseCode: 401}, nil
	}
	if appId != w.AppId {

		if db := r.DBConn.Table("tenant_details").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"app_id": appId, "app_key": appKey, "app_secret": appSecret}).Error; db != nil {
			return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
		}

		return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Smooch configuration Updated successfully."}, nil

	} else {
		if db := r.DBConn.Table("tenant_details").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"app_id": appId, "app_key": appKey, "app_secret": appSecret}).Error; db != nil {
			return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
		}

		return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Smooch configuration Updated successfully."}, nil
	}

}

/***************************************Get smooch configuration*****************************************/
func (r crudRepository) Get_Smooch_configuration(ctx context.Context, domain_uuid string) (*models.Response, error) {
	list := make([]models.Tenant_details, 0)
	row, err := r.DBConn.Raw("select id, domain_uuid,app_id, app_key, app_secret from tenant_details WHERE domain_uuid = ?", domain_uuid).Rows()
	if err != nil {
		return &models.Response{Status: "0", Msg: "Record Not Found", ResponseCode: 401}, nil
	}
	defer row.Close()
	for row.Next() {
		f := models.Tenant_details{}
		if err := row.Scan(&f.Id, &f.Domain_uuid, &f.AppId, &f.AppKey, &f.AppSecret); err != nil {

			return nil, err
		}

		list = append(list, f)
	}
	return &models.Response{Status: "1", Msg: "Record Found", ResponseCode: 200, Tenant_list: list}, nil
}

/**************************************************Delete Smooch configuration***************************/
func (r crudRepository) Delete_Smooch_configuration(ctx context.Context, id int64, domain_uuid string) (*models.Response, error) {
	td := models.Tenant_details{}
	db := r.DBConn.Where("domain_uuid = ? AND id = ?", domain_uuid, id).Delete(&td)
	if db.Error != nil {
		return &models.Response{Status: "0", Msg: "Table not deleted", ResponseCode: 404}, nil
	}
	if db.RowsAffected == 0 {
		return &models.Response{Status: "Not Found", Msg: "Record Doesn't Exist", ResponseCode: 401}, nil
	}
	return &models.Response{Status: "1", Msg: "Smooch Configuration deleted.", ResponseCode: 200}, nil
}

/****************************************Save tenant details*********************************************/
func (r crudRepository) Add_Whatsapp_configuration(ctx context.Context, domain_uuid string, appId string, appKey string, appSecret string, WhatsappIntegrationID string) (*models.Response, error) {
	td := models.WhatsappConfiguration{
		Domain_uuid:           domain_uuid,
		AppId:                 appId,
		AppKey:                appKey,
		AppSecret:             appSecret,
		WhatsappIntegrationID: WhatsappIntegrationID,
	}
	if db := r.DBConn.Table("whatsapp_configurations").Where("whatsapp_integration_id = ?", WhatsappIntegrationID).Find(&td).Error; db != nil {
		st := r.DBConn.Create(&td)
		if st.Error != nil {
			return &models.Response{ResponseCode: 409, Status: "Error", Msg: "Not created"}, nil
		}
		return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Whatsapp configuration saved successfully."}, nil

	}
	return &models.Response{ResponseCode: 409, Status: "Error", Msg: "Whatsapp Integration Id  Already Exist."}, nil
}

/**********************************************Get appID by tenant_domain_uuid******************************/
func (r crudRepository) Get_Whatsapp_configuration(ctx context.Context, domain_uuid string) (*models.Response, error) {
	list := make([]models.WhatsappConfiguration, 0)
	row, err := r.DBConn.Raw("select id, domain_uuid, app_id, app_key, app_secret, whatsapp_integration_id from whatsapp_configurations WHERE domain_uuid = ?", domain_uuid).Rows()
	if err != nil {
		return &models.Response{Status: "0", Msg: "Record Not Found", ResponseCode: 401}, nil
	}
	defer row.Close()
	for row.Next() {
		f := models.WhatsappConfiguration{}
		if err := row.Scan(&f.Id, &f.Domain_uuid, &f.AppId, &f.AppKey, &f.AppSecret, &f.WhatsappIntegrationID); err != nil {

			return nil, err
		}

		list = append(list, f)
	}
	return &models.Response{Status: "1", Msg: "Record Found", ResponseCode: 200, List: list}, nil
}

/**********************************************Update Tenant details*************************************/
func (r crudRepository) Update_Whatsapp_configuration(ctx context.Context, id int64, domain_uuid string, appId string, appKey string, appSecret string, WhatsappIntegrationID string) (*models.Response, error) {
	w := models.WhatsappConfiguration{}
	db := r.DBConn.Where("domain_uuid = ? AND id = ?", domain_uuid, id).Find(&w)
	if db.RowsAffected == 0 {
		return &models.Response{Status: "Not Found", Msg: "Record Doesn't Exist", ResponseCode: 401}, nil
	}
	if err := r.DBConn.Where("whatsapp_integration_id = ?", WhatsappIntegrationID).Find(&w).Error; err != nil {
		if db := r.DBConn.Table("whatsapp_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"app_id": appId, "app_key": appKey, "app_secret": appSecret, "whatsapp_integration_id": WhatsappIntegrationID}).Error; db != nil {
			return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
		}

		return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Whatsapp integration Updated successfully."}, nil
	}
	return &models.Response{ResponseCode: 409, Status: "0", Msg: "Whatsapp id already exist."}, nil
}

/**********************************************Delete Tenant details*************************************/
func (r crudRepository) Delete_Whatsapp_configuration(ctx context.Context, id int64, domain_uuid string) (*models.Response, error) {
	td := models.WhatsappConfiguration{}
	db := r.DBConn.Where("domain_uuid = ? AND id = ?", domain_uuid, id).Delete(&td)
	if db.Error != nil {
		return &models.Response{Status: "0", Msg: "Table not deleted", ResponseCode: 404}, nil
	}
	if db.RowsAffected == 0 {
		return &models.Response{Status: "Not Found", Msg: "Record Doesn't Exist", ResponseCode: 401}, nil
	}
	return &models.Response{Status: "1", Msg: "Whatsapp Configuration deleted.", ResponseCode: 200}, nil
}

/*******************************************Add facebook configuration*******************************/
func (r crudRepository) Add_Facebook_configuration(ctx context.Context, domain_uuid string, appId string, appKey string, appSecret string, FacebookIntegrationId string) (*models.Response, error) {
	td := models.FacebookConfiguration{
		Domain_uuid:           domain_uuid,
		AppId:                 appId,
		AppKey:                appKey,
		AppSecret:             appSecret,
		FacebookIntegrationID: FacebookIntegrationId,
	}
	if db := r.DBConn.Table("facebook_configurations").Where("facebook_integration_id = ?", FacebookIntegrationId).Find(&td).Error; db != nil {
		st := r.DBConn.Create(&td)
		if st.Error != nil {
			return &models.Response{ResponseCode: 409, Status: "Error", Msg: "Not created"}, nil
		}
		return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Facebook configuration saved successfully."}, nil

	}
	return &models.Response{ResponseCode: 409, Status: "Error", Msg: "Facebook Integration Id  Already Exist."}, nil
}

/**********************************************Get appID by tenant_domain_uuid******************************/
func (r crudRepository) Get_Facebook_configuration(ctx context.Context, domain_uuid string) (*models.Response, error) {
	list := make([]models.FacebookConfiguration, 0)
	row, err := r.DBConn.Raw("select id, domain_uuid, app_id, app_key, app_secret, facebook_integration_id from facebook_configurations WHERE domain_uuid = ?", domain_uuid).Rows()
	if err != nil {
		return &models.Response{Status: "0", Msg: "Record Not Found", ResponseCode: 401}, nil
	}
	defer row.Close()
	for row.Next() {
		f := models.FacebookConfiguration{}
		if err := row.Scan(&f.Id, &f.Domain_uuid, &f.AppId, &f.AppKey, &f.AppSecret, &f.FacebookIntegrationID); err != nil {

			return nil, err
		}

		list = append(list, f)
	}
	return &models.Response{Status: "1", Msg: "Record Found", ResponseCode: 200, Fb: list}, nil
}

/**********************************************Update Tenant details*************************************/
func (r crudRepository) Update_Facebook_configuration(ctx context.Context, id int64, domain_uuid string, appId string, appKey string, appSecret string, FacebookIntegrationId string) (*models.Response, error) {
	w := models.FacebookConfiguration{}
	db := r.DBConn.Where("domain_uuid = ? AND id = ?", domain_uuid, id).Find(&w)
	if db.RowsAffected == 0 {
		return &models.Response{Status: "Not Found", Msg: "Record Doesn't Exist", ResponseCode: 401}, nil
	}
	if err := r.DBConn.Where("facebook_integration_id = ?", FacebookIntegrationId).Find(&w).Error; err != nil {
		if db := r.DBConn.Table("facebook_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"app_id": appId, "app_key": appKey, "app_secret": appSecret, "facebook_integration_id": FacebookIntegrationId}).Error; db != nil {
			return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
		}

		return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Facebook integration Updated successfully."}, nil
	}
	return &models.Response{ResponseCode: 409, Status: "0", Msg: "Facebook id already exist."}, nil
}

/**********************************************Delete Tenant details*************************************/
func (r crudRepository) Delete_Facebook_configuration(ctx context.Context, id int64, domain_uuid string) (*models.Response, error) {
	td := models.FacebookConfiguration{}
	db := r.DBConn.Where("domain_uuid = ? AND id = ?", domain_uuid, id).Delete(&td)
	if db.Error != nil {
		return &models.Response{Status: "0", Msg: "Table not deleted", ResponseCode: 404}, nil
	}
	if db.RowsAffected == 0 {
		return &models.Response{Status: "Not Found", Msg: "Record Doesn't Exist", ResponseCode: 401}, nil
	}
	return &models.Response{Status: "1", Msg: "Facebook Configuration deleted.", ResponseCode: 200}, nil
}

/**********************************************List integration**************************************/
func (r crudRepository) List_integration(ctx context.Context, appId string) ([]byte, error) {
	td := models.Tenant_details{}
	db := r.DBConn.Where("app_id = ?", appId).Find(&td)
	if db.Error != nil {
		return nil, db.Error
	}
	req, _ := http.NewRequest("GET", "https://api.smooch.io/v1.1/apps/"+appId+"/integrations", nil)
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(td.AppKey, td.AppSecret)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	} else {

		data, _ := ioutil.ReadAll(res.Body)
		return data, nil
	}
}

/***********************************************Delete All Messages*********************************/
func (r crudRepository) DeleteAllMessage(ctx context.Context, appUserId string, appId string) (*models.Response, error) {
	td := models.Tenant_details{}
	db := r.DBConn.Where("app_id = ?", appId).Find(&td)
	if db.Error != nil {
		return nil, db.Error
	}

	req, err := http.NewRequest("DELETE", "https://api.smooch.io/v1.1/apps/"+appId+"/appusers/"+appUserId+"/messages", nil)
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(td.AppKey, td.AppSecret)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return &models.Response{Status: "0", Msg: "Data not found", ResponseCode: 404}, nil
	} else {
		data, _ := ioutil.ReadAll(res.Body)
		if string(data) == "{}" {
			return &models.Response{Status: "1", Msg: "All Messages are Deleted.", ResponseCode: 200}, nil
		}
		return &models.Response{Msg: "No Message is found."}, nil

	}
}

/**************************************************Delete Message*****************************************/
func (r crudRepository) DeleteMessage(ctx context.Context, appId string, appUserId string, messageId string) (*models.Response, error) {
	td := models.Tenant_details{}
	db := r.DBConn.Where("app_id = ?", appId).Find(&td)
	if db.Error != nil {
		return &models.Response{Status: "0", Msg: "Not found", ResponseCode: 404}, nil
	}
	req, err := http.NewRequest("DELETE", "https://api.smooch.io/v1.1/apps/"+appId+"/appusers/"+appUserId+"/messages/"+messageId, nil)
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(td.AppKey, td.AppSecret)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return &models.Response{Status: "0", Msg: "Data not found", ResponseCode: 404}, nil
	} else {
		data, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(data))
		return &models.Response{Status: "1", Msg: "Message is Deleted.", ResponseCode: 200}, nil

	}
}

/**********************************************Create Text template*************************************/
func (r crudRepository) Create_Text_Template(ctx context.Context, appId string, p models.Comtemplate) ([]byte, error) {
	td := models.Tenant_details{}

	db := r.DBConn.Where("app_id = ?", appId).Find(&td)
	if db.Error != nil {
		return nil, db.Error
	}
	jsonValue, _ := json.Marshal(p)

	req, _ := http.NewRequest("POST", "https://api.smooch.io/v1.1/apps/"+appId+"/templates", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(td.AppKey, td.AppSecret)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	} else {
		data, _ := ioutil.ReadAll(res.Body)
		return data, nil
	}
}

/*************************************************Create carousel template********************************/
func (r crudRepository) Create_Carousel_Template(ctx context.Context, appId string, p models.Payload) ([]byte, error) {
	td := models.Tenant_details{}
	db := r.DBConn.Where("app_id = ?", appId).Find(&td)
	if db.Error != nil {
		return nil, db.Error
	}
	jsonValue, _ := json.Marshal(p)
	req, _ := http.NewRequest("POST", "https://api.smooch.io/v1.1/apps/"+appId+"/templates", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(td.AppKey, td.AppSecret)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	} else {
		data, _ := ioutil.ReadAll(res.Body)

		return data, nil
	}
}

/**************************************************Create Compound Template*******************************/
func (r crudRepository) Create_Compound_Template(ctx context.Context, appId string, p models.Comtemplate) ([]byte, error) {
	td := models.Tenant_details{}
	db := r.DBConn.Where("app_id = ?", appId).Find(&td)
	if db.Error != nil {
		return nil, db.Error
	}
	jsonValue, _ := json.Marshal(p)
	fmt.Println(p.Message.Type, p.Message.Text, p.Name, "bcjbvvwj")
	req, _ := http.NewRequest("POST", "https://api.smooch.io/v1.1/apps/"+appId+"/templates", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(td.AppKey, td.AppSecret)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	} else {
		data, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(data))
		return data, nil
	}
}

/**************************************************Post Message*******************************************/
func (r *crudRepository) PostMessage(ctx context.Context, appId string, appUserId string, p models.User) ([]byte, error) {
	td := models.Tenant_details{}
	db := r.DBConn.Where("app_id = ?", appId).Find(&td)
	if db.Error != nil {
		return nil, db.Error
	}

	jsonValue, _ := json.Marshal(p)
	req, _ := http.NewRequest("POST", "https://api.smooch.io/v1.1/apps/"+appId+"/appusers/"+appUserId+"/messages", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(td.AppKey, td.AppSecret)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	} else {
		data, _ := ioutil.ReadAll(res.Body)
		return data, nil
	}
}

/*********************************************create Quickreply template********************************/
func (r *crudRepository) Create_Quickreply_Template(ctx context.Context, appId string, p models.Comtemplate) ([]byte, error) {
	td := models.Tenant_details{}
	db := r.DBConn.Where("app_id = ?", appId).Find(&td)
	if db.Error != nil {
		return nil, db.Error
	}

	jsonValue, _ := json.Marshal(p)
	req, _ := http.NewRequest("POST", "https://api.smooch.io/v1.1/apps/"+appId+"/templates", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(td.AppKey, td.AppSecret)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	} else {
		data, _ := ioutil.ReadAll(res.Body)
		return data, nil
	}
}

/***************************************Create send location template****************************************/
func (r crudRepository) Create_Request_Location_Template(ctx context.Context, appId string, p models.Comtemplate) ([]byte, error) {
	td := models.Tenant_details{}
	db := r.DBConn.Where("app_id = ?", appId).Find(&td)
	if db.Error != nil {
		return nil, db.Error
	}
	jsonValue, _ := json.Marshal(p)
	req, _ := http.NewRequest("POST", "https://api.smooch.io/v1.1/apps/"+appId+"/templates", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(td.AppKey, td.AppSecret)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	} else {
		data, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(data))
		return data, nil
	}
}

/************************************************Update Text Template**********************************/
func (r crudRepository) Update_Text_Template(ctx context.Context, appId string, template_id string, p models.Payload) ([]byte, error) {
	td := models.Tenant_details{}
	db := r.DBConn.Where("app_id = ?", appId).Find(&td)
	if db.Error != nil {
		return nil, db.Error
	}
	jsonValue, _ := json.Marshal(p)

	req, _ := http.NewRequest("PUT", "https://api.smooch.io/v1.1/apps/"+appId+"/templates/"+template_id, bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(td.AppKey, td.AppSecret)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	} else {
		data, _ := ioutil.ReadAll(res.Body)
		return data, nil
	}
}

/*************************************************get tempalte id************************************/
func (r crudRepository) Get_template(ctx context.Context, appId string, template_id string) ([]byte, error) {
	td := models.Tenant_details{}
	db := r.DBConn.Where("app_id = ?", appId).Find(&td)
	if db.Error != nil {
		return nil, db.Error
	}
	res, err := http.NewRequest("GET", "https://api.smooch.io/v1.1/apps/"+appId+"/templates/"+template_id, nil)
	res.Header.Set("Content-Type", "application/json")
	res.SetBasicAuth(td.AppKey, td.AppSecret)

	client := &http.Client{}
	response, err := client.Do(res)
	if err != nil {
		fmt.Printf("error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		return data, nil
	}
	return nil, err
}

/***************************************************List tempalte*****************************************/
func (r crudRepository) List_template(ctx context.Context, appId string) ([]byte, error) {
	td := models.Tenant_details{}
	db := r.DBConn.Where("app_id = ?", appId).Find(&td)
	if db.Error != nil {
		return nil, db.Error
	}
	res, err := http.NewRequest("GET", "https://api.smooch.io/v1.1/apps/"+appId+"/templates/", nil)
	res.Header.Set("Content-Type", "application/json")
	res.SetBasicAuth(td.AppKey, td.AppSecret)

	client := &http.Client{}
	response, err := client.Do(res)
	if err != nil {
		fmt.Printf("error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		return data, nil
	}
	return nil, err
}

/************************************************Delete template*************************************/
func (r crudRepository) Delete_template(ctx context.Context, appId string, template_id string) (*models.Response, error) {
	td := models.Tenant_details{}
	db := r.DBConn.Where("app_id = ?", appId).Find(&td)
	if db.Error != nil {
		return nil, db.Error
	}
	res, err := http.NewRequest("DELETE", "https://api.smooch.io/v1.1/apps/"+appId+"/templates/"+template_id, nil)
	res.Header.Set("Content-Type", "application/json")
	res.SetBasicAuth(td.AppKey, td.AppSecret)

	client := &http.Client{}
	response, err := client.Do(res)
	if err != nil {
		fmt.Printf("error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(data)
		if string(data) == "{}" {
			return &models.Response{ResponseCode: 200, Status: "DELETED", Msg: "Template Deleted"}, nil
		}

		return &models.Response{Data: data}, nil
	}
	return &models.Response{Status: "0", ResponseCode: 404, Msg: "Error"}, nil
}

/**************************************************Send Location*************************************/
func (r crudRepository) Send_Location(ctx context.Context, appId string, appUserId string, p models.Locations) ([]byte, error) {
	td := models.Tenant_details{}
	db := r.DBConn.Where("app_id = ?", appId).Find(&td)
	if db.Error != nil {
		return nil, db.Error
	}

	jsonValue, _ := json.Marshal(p)
	req, _ := http.NewRequest("POST", "https://api.smooch.io/v1.1/apps/"+appId+"/appusers/"+appUserId+"/messages", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(td.AppKey, td.AppSecret)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	} else {
		data, _ := ioutil.ReadAll(res.Body)
		return data, nil
	}
}

/**************************************************send action message**************************************/
func (r crudRepository) Message_Action_Types(ctx context.Context, appId string, appUserId string, p models.User) ([]byte, error) {
	td := models.Tenant_details{}
	db := r.DBConn.Where("app_id = ?", appId).Find(&td)
	if db.Error != nil {
		return nil, db.Error
	}

	jsonValue, _ := json.Marshal(p)
	req, _ := http.NewRequest("POST", "https://api.smooch.io/v1.1/apps/"+appId+"/appusers/"+appUserId+"/messages", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(td.AppKey, td.AppSecret)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	} else {
		data, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(data))
		return data, nil
	}

}

/**********************************************Quickreply message************************************/
func (r crudRepository) Quickreply_Message(ctx context.Context, appId string, appUserId string, p models.User) ([]byte, error) {
	td := models.Tenant_details{}
	db := r.DBConn.Where("app_id = ?", appId).Find(&td)
	if db.Error != nil {
		return nil, db.Error
	}

	jsonValue, _ := json.Marshal(p)
	req, _ := http.NewRequest("POST", "https://api.smooch.io/v1.1/apps/"+appId+"/appusers/"+appUserId+"/messages", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(td.AppKey, td.AppSecret)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	} else {
		data, _ := ioutil.ReadAll(res.Body)
		return data, nil
	}
}

/************************************************Send carousel message**********************************/
func (r crudRepository) Send_Carousel_Message(ctx context.Context, appId string, appUserId string, p models.User) ([]byte, error) {

	td := models.Tenant_details{
		AppId: appId,
	}
	db := r.DBConn.Table("tenant_details").Where("app_id = ?", appId).Find(&td)
	if db.Error != nil {
		return nil, db.Error
	}

	jsonValue, _ := json.Marshal(p)
	req, _ := http.NewRequest("POST", "https://api.smooch.io/v1.1/apps/"+appId+"/appusers/"+appUserId+"/messages", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(td.AppKey, td.AppSecret)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	} else {
		data, _ := ioutil.ReadAll(res.Body)
		return data, nil
	}
}

/*********************************Link appuser to channel*********************************************/
func (r crudRepository) Link_appUser_to_Channel(ctx context.Context, appId string, appUserId string, p models.Link) ([]byte, error) {

	td := models.Tenant_details{
		AppId: appId,
	}
	db := r.DBConn.Table("tenant_details").Where("app_id = ?", appId).Find(&td)
	if db.Error != nil {
		return nil, db.Error
	}

	jsonValue, _ := json.Marshal(p)
	req, _ := http.NewRequest("POST", "https://api.smooch.io/v1.1/apps/"+appId+"/appusers/"+appUserId+"/channels", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(td.AppKey, td.AppSecret)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	} else {
		data, _ := ioutil.ReadAll(res.Body)
		return data, nil
	}
}

/*************************************Unlink appUser to channel*****************************************/
func (r crudRepository) Unlink_appUser_to_Channel(ctx context.Context, appId string, appUserId string, channel string) ([]byte, error) {

	td := models.Tenant_details{
		AppId: appId,
	}
	db := r.DBConn.Table("tenant_details").Where("app_id = ?", appId).Find(&td)
	if db.Error != nil {
		return nil, db.Error
	}
	req, _ := http.NewRequest("DELETE", "https://api.smooch.io/v1.1/apps/"+appId+"/appusers/"+appUserId+"/channels"+channel, nil)
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(td.AppKey, td.AppSecret)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	} else {
		data, _ := ioutil.ReadAll(res.Body)
		return data, nil
	}
}

/******************************************Upload Attachments********************************************/
func (r crudRepository) Upload_Attachments(ctx context.Context, appId string, appUserId string, Type string, file multipart.File, handler *multipart.FileHeader) ([]byte, error) {

	td := models.Tenant_details{
		AppId: appId,
	}
	db := r.DBConn.Table("tenant_details").Where("app_id = ?", appId).Find(&td)
	if db.Error != nil {
		return nil, db.Error
	}
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("source", handler.Filename)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	part.Write(fileBytes)
	writer.Close()
	req, _ := http.NewRequest("POST", "https://api.smooch.io/v1.1/apps/"+appId+"/attachments?access=public&for=message&appUserId="+appUserId, body)
	req.Header.Add("Content-Type", "multipart/form-data")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	req.SetBasicAuth(td.AppKey, td.AppSecret)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	} else {
		u := models.User{}
		data, _ := ioutil.ReadAll(res.Body)
		jsonData := json.Unmarshal(data, &u)
		fmt.Println(jsonData, u.MediaUrl, u.MediaType, "ghvghvqv")
		p := models.User{
			Role:      "appMaker",
			Type:      Type,
			MediaType: u.MediaType,
			MediaUrl:  u.MediaUrl,
		}
		r.PostMessage(ctx, appId, appUserId, p)
		return data, nil
	}

}

/***********************************************TypingActivity***********************************************/
func (r crudRepository) TypingActivity(ctx context.Context, appUserId string, appId string, p models.User) ([]byte, error) {
	td := models.Tenant_details{}

	db := r.DBConn.Where("app_id = ?", appId).Find(&td)
	if db.Error != nil {
		return nil, db.Error
	}
	jsonValue, _ := json.Marshal(p)
	req, _ := http.NewRequest("POST", "https://api.smooch.io/v1.1/apps/"+appId+"/appusers/"+appUserId+"/conversation/activity", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(td.AppKey, td.AppSecret)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	} else {
		data, _ := ioutil.ReadAll(res.Body)
		return data, nil
	}
}
