package controller

import (
	"bytes"
	"context"
	json "encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	models "whatsapp_api/model"
	crud "whatsapp_api/whatsapp"

	"github.com/labstack/echo"
)

type CRUDController struct {
	usecase crud.Usecase
}

/**************************************************Delete All Message***************************************************/
func (r *CRUDController) DeleteAllMessage(c echo.Context) error {

	appUserId := c.Param("appUserId")
	appId := c.Param("appId")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.DeleteAllMessage(ctx, appUserId, appId)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/**************************************************Post Message***************************************************/
func (r *CRUDController) PostMessage(c echo.Context) error {
	appUserId := c.Param("appUserId")
	appId := c.Param("appId")
	u := models.User{}
	if err := c.Bind(&u); err != nil {
		return err
	}
	p := models.User{
		Role:      "appMaker",
		Type:      u.Type,
		Text:      u.Text,
		MediaType: u.MediaType,
		MediaUrl:  u.MediaUrl,
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.PostMessage(ctx, appId, appUserId, p)

	if authResponse == nil {
		return c.JSONBlob(http.StatusUnauthorized, authResponse)
	}
	return c.JSONBlob(http.StatusOK, authResponse)
}

/**************************************************Get All Message***************************************************/

func (r *CRUDController) GetAllMessageByAppUserId(c echo.Context) error {
	appUserId := c.Param("appUserId")
	appId := c.Param("appId")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.GetAllMessageByAppUserId(ctx, appUserId, appId)

	if authResponse == nil {
		return c.JSONBlob(http.StatusUnauthorized, authResponse)
	}
	return c.JSONBlob(http.StatusOK, authResponse)

}

/**********************************************Get AppUser Details***********************************************/
func (r *CRUDController) GetAppUserDetails(c echo.Context) error {
	appUserId := c.Param("appUserId")
	appId := c.Param("appId")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.GetAppUserDetails(ctx, appUserId, appId)

	if authResponse == nil {
		return c.JSONBlob(http.StatusUnauthorized, authResponse)
	}
	return c.JSONBlob(http.StatusOK, authResponse)

}

/**************************************************Delete Message***************************************************/
func (r *CRUDController) DeleteMessage(c echo.Context) error {
	appUserId := c.Param("appUserId")
	appId := c.Param("appId")
	messageId := c.Param("messageId")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.DeleteMessage(ctx, appId, appUserId, messageId)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)

}

/************************************************Pre-createUser**************************************************/

func (r *CRUDController) Pre_createUser(c echo.Context) error {
	appId := c.Param("appId")
	var crud map[string]interface{}
	err1 := json.NewDecoder(c.Request().Body).Decode(&crud)
	if err1 != nil {
		fmt.Println("err= ", err1)
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Pre_createUser(ctx, appId, crud)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/**********************************************Read appUser_Message**************************************/
func (r *CRUDController) App_user(c echo.Context) error {
	body, error := ioutil.ReadAll(c.Request().Body)
	if error != nil {
		return error
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.App_user(ctx, body)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)

}

/**********************************************get all_id**************************************************/
func (r *CRUDController) Get_allId(c echo.Context) error {
	domain_uuid := c.Param("domain_uuid")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Get_allId(ctx, domain_uuid)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)

}

/********************************************Get Customer Details by appUserId****************************/
func (r *CRUDController) Get_Customer_by_agent_uuid(c echo.Context) error {
	agent_uuid := c.Param("agent_uuid")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Get_Customer_by_agent_uuid(ctx, agent_uuid)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)

}

/**********************************************Create Template********************************************/
func (r *CRUDController) Create_Text_Template(c echo.Context) error {
	appId := c.Param("appId")
	u := models.Comtemplate{}
	if err := c.Bind(&u); err != nil {
		return err
	}
	p := models.Comtemplate{
		Name: u.Name,
		Message: models.User{
			Role:      u.Message.Role,
			Type:      u.Message.Type,
			Text:      u.Message.Text,
			MediaType: u.Message.MediaType,
			MediaUrl:  u.Message.MediaUrl,
		},
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Create_Text_Template(ctx, appId, p)

	if authResponse == nil {
		return c.JSONBlob(http.StatusUnauthorized, authResponse)
	}
	return c.JSONBlob(http.StatusOK, authResponse)

}

/***********************************************Create Carousel Template***********************************/
func (r *CRUDController) Create_Carousel_Template(c echo.Context) error {
	appId := c.Param("appId")
	u := models.Payload{}
	if err := c.Bind(&u); err != nil {
		return err
	}
	p := models.Payload{
		Name: u.Name,
		Message: models.Message{
			Role: u.Message.Role,
			Type: u.Message.Type,
			Items: []models.Items{
				{Title: u.Message.Items[0].Title, Description: u.Message.Items[0].Description, MediaType: u.Message.Items[0].MediaType, MediaUrl: u.Message.Items[0].MediaUrl, Actions: []models.Actions{
					{Type: u.Message.Items[0].Actions[0].Type, Text: u.Message.Items[0].Actions[0].Text, Payload: u.Message.Items[0].Actions[0].Payload},
					{Type: u.Message.Items[0].Actions[1].Type, Text: u.Message.Items[0].Actions[1].Text, Payload: u.Message.Items[0].Actions[1].Payload},
				}},
			},
		},
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Create_Carousel_Template(ctx, appId, p)

	if authResponse == nil {
		return c.JSONBlob(http.StatusUnauthorized, authResponse)
	}
	return c.JSONBlob(http.StatusOK, authResponse)

}

/********************************************create compound template*************************************/
func (r *CRUDController) Create_Compound_Template(c echo.Context) error {
	appId := c.Param("appId")
	u := models.Comtemplate{}
	if err := c.Bind(&u); err != nil {
		return err
	}
	p := models.Comtemplate{
		Name: u.Name,
		Message: models.User{
			Role:      u.Message.Role,
			Type:      u.Message.Type,
			Text:      u.Message.Text,
			MediaUrl:  u.Message.MediaUrl,
			MediaType: u.Message.MediaType,
			Action: []models.Actions{
				{Type: u.Message.Action[0].Type, Text: u.Message.Action[0].Text, Payload: u.Message.Action[0].Payload},
				{Type: u.Message.Action[1].Type, Text: u.Message.Action[1].Text, Payload: u.Message.Action[1].Payload},
			},
		},
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Create_Compound_Template(ctx, appId, p)

	if authResponse == nil {
		return c.JSONBlob(http.StatusUnauthorized, authResponse)
	}
	return c.JSONBlob(http.StatusOK, authResponse)
}

/*******************************************Create quick reply template***********************************/
func (r *CRUDController) Create_Quickreply_Template(c echo.Context) error {
	u := models.Comtemplate{}
	if err := c.Bind(&u); err != nil {
		return err
	}
	p := models.Comtemplate{
		Name: u.Name,
		Message: models.User{
			Role:      u.Message.Role,
			Type:      u.Message.Type,
			Text:      u.Message.Text,
			MediaUrl:  u.Message.MediaUrl,
			MediaType: u.Message.MediaType,
			Action: []models.Actions{
				{Type: u.Message.Action[0].Type, Text: u.Message.Action[0].Text, Payload: u.Message.Action[0].Payload},
				{Type: u.Message.Action[1].Type, Text: u.Message.Action[1].Text, Payload: u.Message.Action[1].Payload},
				{Type: u.Message.Action[2].Type, Text: u.Message.Action[2].Text, Payload: u.Message.Action[2].Payload},
				{Type: u.Message.Action[3].Type, Text: u.Message.Action[3].Text, Payload: u.Message.Action[3].Payload},
			},
		},
	}
	appId := c.Param("appId")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Create_Quickreply_Template(ctx, appId, p)

	if authResponse == nil {
		return c.JSONBlob(http.StatusUnauthorized, authResponse)
	}
	return c.JSONBlob(http.StatusOK, authResponse)

}

/********************************************Create request location template*****************************/
func (r *CRUDController) Create_Request_Location_Template(c echo.Context) error {

	appId := c.Param("appId")
	u := models.Comtemplate{}
	if err := c.Bind(&u); err != nil {
		return err
	}
	p := models.Comtemplate{
		Name: u.Name,
		Message: models.User{
			Role:      u.Message.Role,
			Type:      u.Message.Type,
			Text:      u.Message.Text,
			MediaUrl:  u.Message.MediaUrl,
			MediaType: u.Message.MediaType,
			Action: []models.Actions{
				{Type: u.Message.Action[0].Type, Text: u.Message.Action[0].Text},
			},
		},
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Create_Request_Location_Template(ctx, appId, p)

	if authResponse == nil {
		return c.JSONBlob(http.StatusUnauthorized, authResponse)
	}
	return c.JSONBlob(http.StatusOK, authResponse)

}

/*********************************************Update Template*********************************************/
func (r *CRUDController) Update_Text_Template(c echo.Context) error {
	appId := c.Param("appId")
	template_id := c.Param("template_id")
	u := models.Payload{}
	if err := c.Bind(&u); err != nil {
		return err
	}
	p := models.Payload{
		Name: u.Name,
		Message: models.Message{
			Role: u.Message.Role,
			Type: u.Message.Type,
			Text: u.Message.Text,
		},
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Update_Text_Template(ctx, appId, template_id, p)

	if authResponse == nil {
		return c.JSONBlob(http.StatusUnauthorized, authResponse)
	}
	return c.JSONBlob(http.StatusOK, authResponse)

}

/**********************************************Get Template by id******************************************/
func (r *CRUDController) Get_template(c echo.Context) error {
	appId := c.Param("appId")
	template_id := c.Param("template_id")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Get_template(ctx, appId, template_id)

	if authResponse == nil {
		return c.JSONBlob(http.StatusUnauthorized, authResponse)
	}
	return c.JSONBlob(http.StatusOK, authResponse)

}

/***********************************************List template**********************************************/
func (r *CRUDController) List_template(c echo.Context) error {
	appId := c.Param("appId")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.List_template(ctx, appId)

	if authResponse == nil {
		return c.JSONBlob(http.StatusUnauthorized, authResponse)
	}
	return c.JSONBlob(http.StatusOK, authResponse)

}

/***********************************************Delete template*******************************************/
func (r *CRUDController) Delete_template(c echo.Context) error {
	appId := c.Param("appId")
	template_id := c.Param("template_id")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Delete_template(ctx, appId, template_id)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)

}

/**********************************************Send Location***********************************************/
func (r *CRUDController) Send_Location(c echo.Context) error {
	appId := c.Param("appId")
	appUserId := c.Param("appUserId")
	u := models.Locations{}
	if err := c.Bind(&u); err != nil {
		return err
	}
	p := models.Locations{
		Role: "appMaker",
		Type: u.Type,
		Text: u.Text,
		Coordinates: models.Coordinates{
			Lat:  u.Coordinates.Lat,
			Long: u.Coordinates.Long,
		},
		Location: models.Location{
			Address: u.Location.Address,
			Name:    u.Location.Name,
		},
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Send_Location(ctx, appId, appUserId, p)

	if authResponse == nil {
		return c.JSONBlob(http.StatusUnauthorized, authResponse)
	}
	return c.JSONBlob(http.StatusOK, authResponse)

}

/***********************************************Send Notification*****************************************/
func (r *CRUDController) Send_Notification(c echo.Context) error {

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6ImFwcF81ZWY1ZGFjYTJlOTMwMDAwMGM2OTY5NzgifQ.eyJzY29wZSI6ImFwcCIsImlhdCI6MTU5MzY3NTI3NH0.vdyMVsCdFdvbPmpRL0hxGoRAvECGnroA5zz0PvdyngE"
	var bearer = "Bearer " + token
	u := models.AutoGenerated{}
	if err := c.Bind(&u); err != nil {
		return err
	}
	p := models.AutoGenerated{
		Destination: models.Destination{
			IntegrationID: "5ef6b7c0a5e6d2000cd6533c",
			DestinationID: u.Destination.DestinationID,
		},
		Author: models.Author{
			Role: "appMaker",
		},
		Message: models.Messages{
			Type: u.Message.Type,
			Text: u.Message.Text,
		},
	}

	jsonValue, _ := json.Marshal(p)
	req, _ := http.NewRequest("POST", "https://api.smooch.io/v1.1/apps/5ed5250711e9ad000f2ddd03/notifications", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearer)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	} else {
		data, _ := ioutil.ReadAll(res.Body)
		return c.JSONBlob(200, data)
	}
}

/**********************************************Delete AppUser*********************************************/
func (r *CRUDController) Delete_AppUser(c echo.Context) error {
	appUserId := c.Param("appUserId")
	appId := c.Param("appId")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Delete_AppUser(ctx, appUserId, appId)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/*********************************************Delete AppUser Profile****************************************/
func (r *CRUDController) Delete_AppUser_Profile(c echo.Context) error {
	appId := c.Param("appId")
	appUserId := c.Param("appUserId")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Delete_AppUser_Profile(ctx, appId, appUserId)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/**********************************************update Appuser********************************************/
func (r *CRUDController) Update_AppUser(c echo.Context) error {
	appUserId := c.Param("appUserId")
	appId := c.Param("appId")
	var update map[string]interface{}
	err1 := json.NewDecoder(c.Request().Body).Decode(&update)
	if err1 != nil {
		fmt.Println("err= ", err1)
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Update_AppUser(ctx, appUserId, appId, update)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/**********************************************Message Action*********************************************/
func (r *CRUDController) Message_Action_Types(c echo.Context) error {
	appUserId := c.Param("appUserId")
	appId := c.Param("appId")
	u := models.User{}
	if err := c.Bind(&u); err != nil {
		return err
	}
	p := models.User{
		Role:      u.Role,
		Type:      u.Type,
		Text:      u.Text,
		MediaUrl:  u.MediaUrl,
		MediaType: u.MediaType,
		Action: []models.Actions{
			{Type: u.Action[0].Type, Text: u.Action[0].Text, Amount: u.Action[0].Amount, Size: u.Action[0].Size, URI: u.Action[0].URI, Fallback: u.Action[0].Fallback},
		},
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Message_Action_Types(ctx, appId, appUserId, p)

	if authResponse == nil {
		return c.JSONBlob(http.StatusUnauthorized, authResponse)
	}
	return c.JSONBlob(http.StatusOK, authResponse)

}

/**********************************************Quickreply Message*****************************************/
func (r *CRUDController) Quickreply_Message(c echo.Context) error {
	appUserId := c.Param("appUserId")
	appId := c.Param("appId")
	u := models.User{}
	if err := c.Bind(&u); err != nil {
		return err
	}
	p := models.User{
		Role:      u.Role,
		Type:      u.Type,
		Text:      u.Text,
		MediaUrl:  u.MediaUrl,
		MediaType: u.MediaType,
		Action: []models.Actions{
			{Type: u.Action[0].Type, Text: u.Action[0].Text, Payload: u.Action[0].Payload},
			{Type: u.Action[1].Type, Text: u.Action[1].Text, Payload: u.Action[1].Payload},
			{Type: u.Action[2].Type, Text: u.Action[2].Text, Payload: u.Action[2].Payload},
			{Type: u.Action[3].Type, Text: u.Action[3].Text, Payload: u.Action[3].Payload},
		},
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Quickreply_Message(ctx, appId, appUserId, p)

	if authResponse == nil {
		return c.JSONBlob(http.StatusUnauthorized, authResponse)
	}
	return c.JSONBlob(http.StatusOK, authResponse)

}

/**********************************************Send Carousel message**************************************/
func (r *CRUDController) Send_Carousel_Message(c echo.Context) error {
	appUserId := c.Param("appUserId")
	appId := c.Param("appId")
	u := models.User{}
	if err := c.Bind(&u); err != nil {
		return err
	}
	p := models.User{
		Role: u.Role,
		Type: u.Type,
		Items: []models.Items{
			{Title: u.Items[0].Title, Description: u.Items[0].Description, MediaType: u.Items[0].MediaType, MediaUrl: u.Items[0].MediaUrl, Actions: []models.Actions{
				{Type: u.Items[0].Actions[0].Type, Text: u.Items[0].Actions[0].Text, Payload: u.Items[0].Actions[0].Payload, URI: u.Items[0].Actions[0].URI},
				{Type: u.Items[0].Actions[1].Type, Text: u.Items[0].Actions[1].Text, Payload: u.Items[0].Actions[1].Payload, URI: u.Items[0].Actions[0].URI},
			}},
		},
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Send_Carousel_Message(ctx, appId, appUserId, p)

	if authResponse == nil {
		return c.JSONBlob(http.StatusUnauthorized, authResponse)
	}
	return c.JSONBlob(http.StatusOK, authResponse)

}

/*****************************************Add Smooch AccountInfo******************************************/
func (r *CRUDController) Add_Smooch_configuration(c echo.Context) error {
	var Add_Smooch_configuration map[string]interface{}
	c.Bind(&Add_Smooch_configuration)
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Add_Smooch_configuration(ctx, Add_Smooch_configuration)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)

}

/**********************************************Update Smooch configuration***********************************/
func (r *CRUDController) Update_Smooch_configuration(c echo.Context) error {
	domain_uuid := c.Param("domain_uuid")
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	id := int64(idP)
	var Update_Smooch_configuration map[string]interface{}
	err1 := json.NewDecoder(c.Request().Body).Decode(&Update_Smooch_configuration)
	if err1 != nil {
		fmt.Println("err= ", err1)
	} else {
		fmt.Println("err= ", err1)
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Update_Smooch_configuration(ctx, id, domain_uuid, Update_Smooch_configuration)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)

}

/***********************************************Delete Smooch configuration*********************************/
func (r *CRUDController) Delete_Smooch_configuration(c echo.Context) error {
	domain_uuid := c.Param("domain_uuid")
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	id := int64(idP)
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Delete_Smooch_configuration(ctx, id, domain_uuid)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)

}

/***********************************************Get Smooch Configuration*************************************/
func (r *CRUDController) Get_Smooch_configuration(c echo.Context) error {
	domain_uuid := c.Param("domain_uuid")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Get_Smooch_configuration(ctx, domain_uuid)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)

}

/***********************************************Save wghatsappconfiguration***************************************/
func (r *CRUDController) Add_Whatsapp_configuration(c echo.Context) error {
	u := models.WhatsappConfigurations{}
	if err := c.Bind(&u); err != nil {
		return err
	}
	td := models.WhatsappConfigurations{
		Domain_uuid:           u.Domain_uuid,
		ConfigurationName:     u.ConfigurationName,
		AppId:                 u.AppId,
		AppKey:                u.AppKey,
		Message:               u.Message,
		AppSecret:             u.AppSecret,
		Size:                  u.Size,
		WhatsappIntegrationID: u.WhatsappIntegrationID,
		Trigger: models.Trigger{
			When:    u.Trigger.When,
			Name:    u.Trigger.Name,
			Message: u.Trigger.Message,
		},
		WorkingDays: []models.WorkingDays{
			{Day: u.WorkingDays[0].Day, WorkingHourStartTime: u.WorkingDays[0].WorkingHourStartTime, WorkingHourEndTime: u.WorkingDays[0].WorkingHourEndTime},
			{Day: u.WorkingDays[1].Day, WorkingHourStartTime: u.WorkingDays[1].WorkingHourStartTime, WorkingHourEndTime: u.WorkingDays[1].WorkingHourEndTime},
			{Day: u.WorkingDays[2].Day, WorkingHourStartTime: u.WorkingDays[2].WorkingHourStartTime, WorkingHourEndTime: u.WorkingDays[2].WorkingHourEndTime},
			{Day: u.WorkingDays[3].Day, WorkingHourStartTime: u.WorkingDays[3].WorkingHourStartTime, WorkingHourEndTime: u.WorkingDays[3].WorkingHourEndTime},
			{Day: u.WorkingDays[4].Day, WorkingHourStartTime: u.WorkingDays[4].WorkingHourStartTime, WorkingHourEndTime: u.WorkingDays[4].WorkingHourEndTime},
			{Day: u.WorkingDays[5].Day, WorkingHourStartTime: u.WorkingDays[5].WorkingHourStartTime, WorkingHourEndTime: u.WorkingDays[5].WorkingHourEndTime},
			{Day: u.WorkingDays[6].Day, WorkingHourStartTime: u.WorkingDays[6].WorkingHourStartTime, WorkingHourEndTime: u.WorkingDays[6].WorkingHourEndTime},
		},
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Add_Whatsapp_configuration(ctx, td)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)

}

/*********************************************** Get Whatsapp configuration******************************************/
func (r *CRUDController) Get_Whatsapp_configuration(c echo.Context) error {
	domain_uuid := c.Param("domain_uuid")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Get_Whatsapp_configuration(ctx, domain_uuid)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)

}

/**********************************************Update tenant details**************************************/
func (r *CRUDController) Update_Whatsapp_configuration(c echo.Context) error {
	domain_uuid := c.Param("domain_uuid")
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	id := int64(idP)
	u := models.WhatsappConfigurations{}
	if err := c.Bind(&u); err != nil {
		return err
	}
	td := models.WhatsappConfigurations{
		Domain_uuid:           u.Domain_uuid,
		ConfigurationName:     u.ConfigurationName,
		AppId:                 u.AppId,
		AppKey:                u.AppKey,
		AppSecret:             u.AppSecret,
		Size:                  u.Size,
		Message:               u.Message,
		WhatsappIntegrationID: u.WhatsappIntegrationID,
		Trigger: models.Trigger{
			When:    u.Trigger.When,
			Name:    u.Trigger.Name,
			Message: u.Trigger.Message,
		},
		WorkingDays: []models.WorkingDays{
			{Day: u.WorkingDays[0].Day, WorkingHourStartTime: u.WorkingDays[0].WorkingHourStartTime, WorkingHourEndTime: u.WorkingDays[0].WorkingHourEndTime},
			{Day: u.WorkingDays[1].Day, WorkingHourStartTime: u.WorkingDays[1].WorkingHourStartTime, WorkingHourEndTime: u.WorkingDays[1].WorkingHourEndTime},
			{Day: u.WorkingDays[2].Day, WorkingHourStartTime: u.WorkingDays[2].WorkingHourStartTime, WorkingHourEndTime: u.WorkingDays[2].WorkingHourEndTime},
			{Day: u.WorkingDays[3].Day, WorkingHourStartTime: u.WorkingDays[3].WorkingHourStartTime, WorkingHourEndTime: u.WorkingDays[3].WorkingHourEndTime},
			{Day: u.WorkingDays[4].Day, WorkingHourStartTime: u.WorkingDays[4].WorkingHourStartTime, WorkingHourEndTime: u.WorkingDays[4].WorkingHourEndTime},
			{Day: u.WorkingDays[5].Day, WorkingHourStartTime: u.WorkingDays[5].WorkingHourStartTime, WorkingHourEndTime: u.WorkingDays[5].WorkingHourEndTime},
			{Day: u.WorkingDays[6].Day, WorkingHourStartTime: u.WorkingDays[6].WorkingHourStartTime, WorkingHourEndTime: u.WorkingDays[6].WorkingHourEndTime},
		},
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Update_Whatsapp_configuration(ctx, id, domain_uuid, td)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)

}

/***********************************************Delete Tenant details**************************************/
func (r *CRUDController) Delete_Whatsapp_configuration(c echo.Context) error {
	domain_uuid := c.Param("domain_uuid")
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	id := int64(idP)
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Delete_Whatsapp_configuration(ctx, id, domain_uuid)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)

}

/*********************************************Add Facebook configuration**********************************/
func (r *CRUDController) Add_Facebook_configuration(c echo.Context) error {
	u := models.FacebookConfigurations{}
	if err := c.Bind(&u); err != nil {
		return err
	}
	td := models.FacebookConfigurations{
		Domain_uuid:           u.Domain_uuid,
		ConfigurationName:     u.ConfigurationName,
		AppId:                 u.AppId,
		AppKey:                u.AppKey,
		Message:               u.Message,
		Size:                  u.Size,
		AppSecret:             u.AppSecret,
		FacebookIntegrationID: u.FacebookIntegrationID,
		Trigger: models.Trigger{
			When:    u.Trigger.When,
			Name:    u.Trigger.Name,
			Message: u.Trigger.Message,
		},
		WorkingDays: []models.WorkingDays{
			{Day: u.WorkingDays[0].Day, WorkingHourStartTime: u.WorkingDays[0].WorkingHourStartTime, WorkingHourEndTime: u.WorkingDays[0].WorkingHourEndTime},
			{Day: u.WorkingDays[1].Day, WorkingHourStartTime: u.WorkingDays[1].WorkingHourStartTime, WorkingHourEndTime: u.WorkingDays[1].WorkingHourEndTime},
			{Day: u.WorkingDays[2].Day, WorkingHourStartTime: u.WorkingDays[2].WorkingHourStartTime, WorkingHourEndTime: u.WorkingDays[2].WorkingHourEndTime},
			{Day: u.WorkingDays[3].Day, WorkingHourStartTime: u.WorkingDays[3].WorkingHourStartTime, WorkingHourEndTime: u.WorkingDays[3].WorkingHourEndTime},
			{Day: u.WorkingDays[4].Day, WorkingHourStartTime: u.WorkingDays[4].WorkingHourStartTime, WorkingHourEndTime: u.WorkingDays[4].WorkingHourEndTime},
			{Day: u.WorkingDays[5].Day, WorkingHourStartTime: u.WorkingDays[5].WorkingHourStartTime, WorkingHourEndTime: u.WorkingDays[5].WorkingHourEndTime},
			{Day: u.WorkingDays[6].Day, WorkingHourStartTime: u.WorkingDays[6].WorkingHourStartTime, WorkingHourEndTime: u.WorkingDays[6].WorkingHourEndTime},
		},
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Add_Facebook_configuration(ctx, td)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)

}

/*********************************************** Get facebook configuration******************************************/
func (r *CRUDController) Get_Facebook_configuration(c echo.Context) error {
	domain_uuid := c.Param("domain_uuid")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Get_Facebook_configuration(ctx, domain_uuid)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)

}

/**********************************************Update Facebook configuration**************************************/
func (r *CRUDController) Update_Facebook_configuration(c echo.Context) error {
	domain_uuid := c.Param("domain_uuid")
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	id := int64(idP)
	u := models.FacebookConfigurations{}
	if err := c.Bind(&u); err != nil {
		return err
	}
	td := models.FacebookConfigurations{
		Domain_uuid:           u.Domain_uuid,
		ConfigurationName:     u.ConfigurationName,
		AppId:                 u.AppId,
		AppKey:                u.AppKey,
		AppSecret:             u.AppSecret,
		Message:               u.Message,
		Size:                  u.Size,
		FacebookIntegrationID: u.FacebookIntegrationID,
		Trigger: models.Trigger{
			When:    u.Trigger.When,
			Name:    u.Trigger.Name,
			Message: u.Trigger.Message,
		},
		WorkingDays: []models.WorkingDays{
			{Day: u.WorkingDays[0].Day, WorkingHourStartTime: u.WorkingDays[0].WorkingHourStartTime, WorkingHourEndTime: u.WorkingDays[0].WorkingHourEndTime},
			{Day: u.WorkingDays[1].Day, WorkingHourStartTime: u.WorkingDays[1].WorkingHourStartTime, WorkingHourEndTime: u.WorkingDays[1].WorkingHourEndTime},
			{Day: u.WorkingDays[2].Day, WorkingHourStartTime: u.WorkingDays[2].WorkingHourStartTime, WorkingHourEndTime: u.WorkingDays[2].WorkingHourEndTime},
			{Day: u.WorkingDays[3].Day, WorkingHourStartTime: u.WorkingDays[3].WorkingHourStartTime, WorkingHourEndTime: u.WorkingDays[3].WorkingHourEndTime},
			{Day: u.WorkingDays[4].Day, WorkingHourStartTime: u.WorkingDays[4].WorkingHourStartTime, WorkingHourEndTime: u.WorkingDays[4].WorkingHourEndTime},
			{Day: u.WorkingDays[5].Day, WorkingHourStartTime: u.WorkingDays[5].WorkingHourStartTime, WorkingHourEndTime: u.WorkingDays[5].WorkingHourEndTime},
			{Day: u.WorkingDays[6].Day, WorkingHourStartTime: u.WorkingDays[6].WorkingHourStartTime, WorkingHourEndTime: u.WorkingDays[6].WorkingHourEndTime},
		},
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Update_Facebook_configuration(ctx, id, domain_uuid, td)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)

}

/***********************************************Delete Facebook configuration*******************************/
func (r *CRUDController) Delete_Facebook_configuration(c echo.Context) error {
	domain_uuid := c.Param("domain_uuid")
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	id := int64(idP)
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Delete_Facebook_configuration(ctx, id, domain_uuid)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)

}

/**********************************************List integration*******************************************/
func (r *CRUDController) List_integration(c echo.Context) error {
	appId := c.Param("appId")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.List_integration(ctx, appId)

	if authResponse == nil {
		return c.JSONBlob(http.StatusUnauthorized, authResponse)
	}
	return c.JSONBlob(http.StatusOK, authResponse)

}

/*********************************************Link appUser to Channel**************************************/
func (r *CRUDController) Link_appUser_to_Channel(c echo.Context) error {
	appId := c.Param("appId")
	appUserId := c.Param("appUserId")
	u := models.Link{}
	if err := c.Bind(&u); err != nil {
		return err
	}

	p := models.Link{
		Type: u.Type,
		Confirmation: models.Confirmation{
			Type: u.Confirmation.Type,
		},
		Address: u.Address,
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Link_appUser_to_Channel(ctx, appId, appUserId, p)

	if authResponse == nil {
		return c.JSONBlob(http.StatusUnauthorized, authResponse)
	}
	return c.JSONBlob(http.StatusOK, authResponse)

}

/**********************************************unlink appuser to Channel**********************************/
func (r *CRUDController) Unlink_appUser_to_Channel(c echo.Context) error {
	appId := c.Param("appId")
	appUserId := c.Param("appUserId")
	channel := c.Param("channel")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Unlink_appUser_to_Channel(ctx, appId, appUserId, channel)

	if authResponse == nil {
		return c.JSONBlob(http.StatusUnauthorized, authResponse)
	}
	return c.JSONBlob(http.StatusOK, authResponse)

}

/**********************************************Upload Attachments******************************************/
func (r *CRUDController) Upload_Attachments(c echo.Context) error {
	appId := c.Param("appId")
	appUserId := c.Param("appUserId")
	Type := c.FormValue("type")
	IntegrationID := c.FormValue("integration_id")
	var Size int64
	err1 := json.NewDecoder(c.Request().Body).Decode(&Size)
	if err1 != nil {
		fmt.Println("err= ", err1)
	}

	err := c.Request().ParseMultipartForm(Size << 20) // 25Mb
	if err != nil {
		return err
	}
	file, handler, err := c.Request().FormFile("source")
	if err != nil {
		return err
	}
	defer file.Close()

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Upload_Attachments(ctx, appId, appUserId, Type, IntegrationID, Size, file, handler)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)

}

/**********************************************AppMaker Typing Activity************************************/
func (r *CRUDController) TypingActivity(c echo.Context) error {
	appId := c.Param("appId")
	appUserId := c.Param("appUserId")
	u := models.User{}
	if err := c.Bind(&u); err != nil {
		return err
	}
	p := models.User{
		Role:      "appMaker",
		Type:      "typing:" + u.Type,
		Name:      u.Name,
		AvatarURL: u.AvatarURL,
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.TypingActivity(ctx, appId, appUserId, p)

	if authResponse == nil {
		return c.JSONBlob(http.StatusUnauthorized, authResponse)
	}
	return c.JSONBlob(http.StatusOK, authResponse)

}

/**********************************************Disable AppUser**********************************************/
func (r *CRUDController) Disable_AppUser(c echo.Context) error {

	appUserId := c.Param("appUserId")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Disable_AppUser(ctx, appUserId)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)

}

/**********************************************Reset Unread count******************************************/
func (r *CRUDController) Reset_Unread_Count(c echo.Context) error {
	appId := c.Param("appId")
	appUserId := c.Param("appUserId")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Reset_Unread_Count(ctx, appId, appUserId)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)

}

/**********************************************Create Queue************************************************/
func (r *CRUDController) Create_Queue(c echo.Context) error {
	var create_queue map[string]interface{}
	err1 := json.NewDecoder(c.Request().Body).Decode(&create_queue)
	if err1 != nil {
		fmt.Println("err= ", err1)
	} else {
		fmt.Println("err= ", err1)
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Create_Queue(ctx, create_queue)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)

}

/***********************************************Assign Agent************************************************/
func (r *CRUDController) Assign_Agent_To_Queue(c echo.Context) error {
	var assign_agent map[string]interface{}
	err1 := json.NewDecoder(c.Request().Body).Decode(&assign_agent)
	if err1 != nil {
		fmt.Println("err= ", err1)
	} else {
		fmt.Println("err= ", err1)
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Assign_Agent_To_Queue(ctx, assign_agent)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)

}

/**********************************************Update Queue************************************************/
func (r *CRUDController) Update_Queue(c echo.Context) error {
	queue_uuid := c.Param("queue_uuid")
	var assign_agent map[string]interface{}
	err1 := json.NewDecoder(c.Request().Body).Decode(&assign_agent)
	if err1 != nil {
		fmt.Println("err= ", err1)
	} else {
		fmt.Println("err= ", err1)
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Update_Queue(ctx, queue_uuid, assign_agent)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)

}

/*******************************************Delete Queue****************************************************/
func (r *CRUDController) Delete_Queue(c echo.Context) error {
	queue_uuid := c.Param("queue_uuid")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Delete_Queue(ctx, queue_uuid)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)

}

/**********************************************Remove Agent From Queue**************************************/
func (r *CRUDController) Remove_Agent_From_Queue(c echo.Context) error {
	agent_uuid := c.Param("agent_uuid")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Remove_Agent_From_Queue(ctx, agent_uuid)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/*****************************************Get Assigned agents from queue List********************************/
func (r *CRUDController) Get_Assigned_Agent_list_From_Queue(c echo.Context) error {
	queue_uuid := c.Param("queue_uuid")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Get_Assigned_Agent_list_From_Queue(ctx, queue_uuid)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/*********************************************Get Queue List*************************************************/
func (r *CRUDController) Get_Queue_List(c echo.Context) error {
	domain_uuid := c.Param("domain_uuid")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Get_Queue_List(ctx, domain_uuid)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/*************************************************Available Agents************************************************/
func (r *CRUDController) Available_Agents(c echo.Context) error {
	domain_uuid := c.Param("domain_uuid")
	queue_uuid := c.Param("queue_uuid")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Available_Agents(ctx, domain_uuid, queue_uuid)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/**********************************************Transfer Customer******************************************/
func (r *CRUDController) Transfer_customer(c echo.Context) error {
	var transfer_customer map[string]interface{}
	c.Bind(&transfer_customer)
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Transfer_customer(ctx, transfer_customer)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/*********************************************Post on facebook Page*****************************************/
func (r *CRUDController) Publish_Post_on_FB_Page(c echo.Context) error {
	// pageId := c.Param("page_id")
	// message := c.Param("message")
	// access_token := c.Param("access_token")
	var fb map[string]interface{}
	err1 := json.NewDecoder(c.Request().Body).Decode(&fb)
	if err1 != nil {
		fmt.Println("err= ", err1)
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Publish_Post_on_FB_Page(ctx, fb)

	if authResponse == nil {
		return c.JSONBlob(http.StatusUnauthorized, authResponse)
	}
	return c.JSONBlob(http.StatusOK, authResponse)
}

/***********************************************Get all posts of a page****************************************/
func (r *CRUDController) Getall_Post_of_Page(c echo.Context) error {
	pageId := c.Param("page_id")
	access_token := c.Param("access_token")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Getall_Post_of_Page(ctx, pageId, access_token)

	if authResponse == nil {
		return c.JSONBlob(http.StatusUnauthorized, authResponse)
	}
	return c.JSONBlob(http.StatusOK, authResponse)
}

/********************************************Delete a Post of a Page**************************************/
func (r *CRUDController) Delete_Post_of_Page(c echo.Context) error {
	page_postId := c.Param("page_post_id")
	access_token := c.Param("access_token")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Delete_Post_of_Page(ctx, page_postId, access_token)

	if authResponse == nil {
		return c.JSONBlob(http.StatusUnauthorized, authResponse)
	}
	return c.JSONBlob(http.StatusOK, authResponse)
}

/**********************************************Update Post of a Page***************************************/
func (r *CRUDController) Update_Post_of_Page(c echo.Context) error {
	// page_postId := c.Param("page_post_id")
	// message := c.Param("message")
	// access_token := c.Param("access_token")
	var fb map[string]interface{}
	err1 := json.NewDecoder(c.Request().Body).Decode(&fb)
	if err1 != nil {
		fmt.Println("err= ", err1)
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Update_Post_of_Page(ctx, fb)

	if authResponse == nil {
		return c.JSONBlob(http.StatusUnauthorized, authResponse)
	}
	return c.JSONBlob(http.StatusOK, authResponse)
}

/**********************************************Get Comments of a Post of Page*********************************/
func (r *CRUDController) Get_Comments_on_Post_of_Page(c echo.Context) error {
	page_postId := c.Param("page_post_id")
	access_token := c.Param("access_token")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Get_Comments_on_Post_of_Page(ctx, page_postId, access_token)

	if authResponse == nil {
		return c.JSONBlob(http.StatusUnauthorized, authResponse)
	}
	return c.JSONBlob(http.StatusOK, authResponse)
}

/**********************************************Comment on post of a page**********************************/
func (r *CRUDController) Comment_on_Post_of_Page(c echo.Context) error {
	// page_postId := c.Param("page_post_id")
	// message := c.Param("message")
	// access_token := c.Param("access_token")
	var fb map[string]interface{}
	err1 := json.NewDecoder(c.Request().Body).Decode(&fb)
	if err1 != nil {
		fmt.Println("err= ", err1)
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Comment_on_Post_of_Page(ctx, fb)

	if authResponse == nil {
		return c.JSONBlob(http.StatusUnauthorized, authResponse)
	}
	return c.JSONBlob(http.StatusOK, authResponse)
}

/******************************************************************************************************/
func (r *CRUDController) UVoiceFacebookLogin(c echo.Context) error {

	client_id := c.Param("client_id")
	flac_uuid := c.Param("flac_uuid")
	client_secret := c.Param("client_secret")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.UVoiceFacebookLogin(ctx, c, client_id, client_secret, flac_uuid)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/************************************************/
func (r *CRUDController) UVoiceFacebookLoginCallback(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.UVoiceFacebookLoginCallback(ctx, c)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/***********************************************Get Page Id************************************************/
func (r *CRUDController) Get_Page_ID(c echo.Context) error {
	userId := c.Param("user_id")
	access_token := c.Param("access_token")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Get_Page_ID(ctx, userId, access_token)

	if authResponse == nil {
		return c.JSONBlob(http.StatusUnauthorized, authResponse)
	}
	return c.JSONBlob(http.StatusOK, authResponse)
}

/************************************************Publish Link with message***********************************/
func (r *CRUDController) Publish_link_with_message_on_Post(c echo.Context) error {
	var fb map[string]interface{}
	err1 := json.NewDecoder(c.Request().Body).Decode(&fb)
	if err1 != nil {
		fmt.Println("err= ", err1)
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Publish_link_with_message_on_Post(ctx, fb)

	if authResponse == nil {
		return c.JSONBlob(http.StatusUnauthorized, authResponse)
	}
	return c.JSONBlob(http.StatusOK, authResponse)
}

/************************************************Schedule a Post*********************************************/
func (r *CRUDController) Schedule_Post(c echo.Context) error {
	var fb map[string]interface{}
	err1 := json.NewDecoder(c.Request().Body).Decode(&fb)
	if err1 != nil {
		fmt.Println("err= ", err1)
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Schedule_Post(ctx, fb)

	if authResponse == nil {
		return c.JSONBlob(http.StatusUnauthorized, authResponse)
	}
	return c.JSONBlob(http.StatusOK, authResponse)
}
func (r *CRUDController) AddFacebookApplication(c echo.Context) error {
	var t map[string]interface{}
	c.Bind(&t)
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.AddFacebookApplication(ctx, t["domain_uuid"].(string), t["app_id"].(string), t["app_secret"].(string), t["app_name"].(string))

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}
func (r *CRUDController) ShowFacebookApplication(c echo.Context) error {
	domain_uuid := c.Param("domain_uuid")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.ShowFacebookApplication(ctx, domain_uuid)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}
func (r *CRUDController) DeleteFacebookApplication(c echo.Context) error {
	flac_uuid := c.Param("flac_uuid")
	domain_uuid := c.Param("domain_uuid")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.DeleteFacebookApplication(ctx, domain_uuid, flac_uuid)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/***********************************************Upload Photo on Post *************************************/
func (r *CRUDController) Upload_Photo_on_Post(c echo.Context) error {
	pageId := c.FormValue("page_id")
	access_token := c.FormValue("access_token")
	err := c.Request().ParseMultipartForm(10 << 20) // 25Mb
	if err != nil {
		return err
	}
	file, handler, err := c.Request().FormFile("source")
	if err != nil {
		return err
	}
	defer file.Close()

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Upload_Photo_on_Post(ctx, pageId, access_token, file, handler)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)

}

func UVoiceFacebookLoginStatus(c echo.Context) error {
	c.Response().Write([]byte(`<html><body><a href="javascript:window.open('','_self').close();">facebook login success!! please close this window</a></body></html>`))
	return nil
}

func (r *CRUDController) AssignAgentToFacebookApplication(c echo.Context) error {
	var t map[string]interface{}
	c.Bind(&t)
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.AssignAgentToFacebookApplication(ctx, t["domain_uuid"].(string), t["flac_uuid"].(string), t["agent_uuid"].(string))

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

func (r *CRUDController) AgentListAssignedToFacebookApplication(c echo.Context) error {
	flac_uuid := c.Param("flac_uuid")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.AgentListAssignedToFacebookApplication(ctx, flac_uuid)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}
func (r *CRUDController) AgentListNotInFacebookApplication(c echo.Context) error {
	flac_uuid := c.Param("flac_uuid")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.AgentListNotInFacebookApplication(ctx, flac_uuid)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}
func (r *CRUDController) ShowAgentFacebookApplication(c echo.Context) error {
	agent_uuid := c.Param("agent_uuid")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.ShowAgentFacebookApplication(ctx, agent_uuid)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/***********************************************Router*****************************************************/

func NewCRUDController(e *echo.Echo, crudusecase crud.Usecase) {
	handler := &CRUDController{
		usecase: crudusecase,
	}

	e.DELETE("delete_allMessage/:appId/:appUserId", handler.DeleteAllMessage)
	e.GET("get_allMessage_byAppUserId/:appId/:appUserId", handler.GetAllMessageByAppUserId)
	e.POST("pre-createUser/:appId", handler.Pre_createUser)
	e.POST("post_message/:appId/:appUserId", handler.PostMessage)
	e.DELETE("delete_message/:appId/:appUserId/:messageId", handler.DeleteMessage)
	e.POST("/messages", handler.App_user)
	e.GET("/getall_appUserId/:domain_uuid", handler.Get_allId)
	e.GET("getcustomerbyagent_uuid/:agent_uuid", handler.Get_Customer_by_agent_uuid)
	e.GET("/get_appUser_details/:appId/:appUserId", handler.GetAppUserDetails)
	e.POST("/create_text_template/:appId", handler.Create_Text_Template)
	e.POST("/create_compound_template/:appId", handler.Create_Compound_Template)
	e.POST("/create_carousel_template/:appId", handler.Create_Carousel_Template)
	e.POST("/create_quickreply_template/:appId", handler.Create_Quickreply_Template)
	e.POST("/create_request_location_template/:appId", handler.Create_Request_Location_Template)
	e.PUT("/update_text_template/:appId/:template_id", handler.Update_Text_Template)
	e.GET("get_template/:appId/:template_id", handler.Get_template)
	e.GET("list_template/:appId", handler.List_template)
	e.DELETE("delete_template/:appId/:template_id", handler.Delete_template)
	e.POST("send_location/:appId/:appUserId", handler.Send_Location)
	e.POST("send_notification/:appId", handler.Send_Notification)
	e.DELETE("delete_appUser/:appId/:appUserId", handler.Delete_AppUser)
	e.DELETE("delete_appUser_profile/:appId/:appUserId", handler.Delete_AppUser_Profile)
	e.PUT("update_appUser/:appId/:appUserId", handler.Update_AppUser)
	e.POST("send_messages/:appId/:appUserId", handler.Message_Action_Types)
	e.POST("quickreply_message/:appId/:appUserId", handler.Quickreply_Message)
	e.POST("send_carousel_message/:appId/:appUserId", handler.Send_Carousel_Message)
	e.GET("disable_appUser/:appUserId", handler.Disable_AppUser)
	e.GET("reset_unread_count/:appId/:appUserId", handler.Reset_Unread_Count)

	e.POST("add_smoochConfiguration", handler.Add_Smooch_configuration)
	e.GET("get_smoochConfiguration/:domain_uuid", handler.Get_Smooch_configuration)
	e.POST("update_smoochConfiguration/:id/:domain_uuid", handler.Update_Smooch_configuration)
	e.DELETE("delete_smoochConfiguration/:id/:domain_uuid", handler.Delete_Smooch_configuration)

	e.POST("add_whatsappConfiguration", handler.Add_Whatsapp_configuration)
	e.GET("get_whatsappConfiguration/:domain_uuid", handler.Get_Whatsapp_configuration)
	e.POST("update_whatsappConfiguration/:id/:domain_uuid", handler.Update_Whatsapp_configuration)
	e.DELETE("delete_whatsappConfiguration/:id/:domain_uuid", handler.Delete_Whatsapp_configuration)

	e.POST("add_facebookConfiguration", handler.Add_Facebook_configuration)
	e.GET("get_facebookConfiguration/:domain_uuid", handler.Get_Facebook_configuration)
	e.POST("update_facebookConfiguration/:id/:domain_uuid", handler.Update_Facebook_configuration)
	e.DELETE("delete_facebookConfiguration/:id/:domain_uuid", handler.Delete_Facebook_configuration)

	e.GET("list_integration/:appId", handler.List_integration)
	e.POST("link_appUser_to_channel/:appId/:appUserId", handler.Link_appUser_to_Channel)
	e.DELETE("unlink_appUser_to_channel/:appId/:appUserId/:channel", handler.Unlink_appUser_to_Channel)
	e.POST("upload_attachments/:appId/:appUserId", handler.Upload_Attachments)
	e.POST("typing_activity/:appId/:appUserId", handler.TypingActivity)

	e.POST("create_queue", handler.Create_Queue)
	e.POST("assign_agent_to_queue", handler.Assign_Agent_To_Queue)
	e.DELETE("remove_agent_from_queue/:agent_uuid", handler.Remove_Agent_From_Queue)
	e.GET("get_assigned_agent_list/:queue_uuid", handler.Get_Assigned_Agent_list_From_Queue)
	e.GET("get_queue_list/:domain_uuid", handler.Get_Queue_List)
	e.POST("update_queue/:queue_uuid", handler.Update_Queue)
	e.DELETE("delete_queue/:queue_uuid", handler.Delete_Queue)
	e.GET("available_agent_list/:domain_uuid/:queue_uuid", handler.Available_Agents)
	e.POST("transfer_customer", handler.Transfer_customer)

	e.POST("publish_post", handler.Publish_Post_on_FB_Page)
	e.GET("getall_post_of_page/:page_id/:access_token", handler.Getall_Post_of_Page)
	e.DELETE("delete_post_of_page/:page_post_id/:access_token", handler.Delete_Post_of_Page)
	e.POST("update_post/:page_post_id/:message/:access_token", handler.Update_Post_of_Page)
	e.GET("get_comments_of_post/:page_post_id/:access_token", handler.Get_Comments_on_Post_of_Page)
	e.POST("comment_on_post/:page_postId/:message/:access_token", handler.Comment_on_Post_of_Page)
	e.GET("get_page_id/:user_id/:access_token", handler.Get_Page_ID)
	e.POST("schedule_post", handler.Schedule_Post)
	e.POST("publish_link_with_message", handler.Publish_link_with_message_on_Post)
	e.POST("upload_photo_on_post", handler.Upload_Photo_on_Post)

	e.GET("/uvoice-facebook-login/:client_id/:client_secret/:flac_uuid", handler.UVoiceFacebookLogin)
	e.GET("/uvoice-facebook-login-callback", handler.UVoiceFacebookLoginCallback)
	e.GET("/uvoice-facebook-login-status", UVoiceFacebookLoginStatus)

	e.POST("/add-facebook-application", handler.AddFacebookApplication)
	e.GET("/show-facebook-application/:domain_uuid", handler.ShowFacebookApplication)
	e.DELETE("/delete-facebook-application/:flac_uuid/:domain_uuid", handler.DeleteFacebookApplication)

	e.POST("/assign-agent-to-facebook-application", handler.AssignAgentToFacebookApplication)
	e.GET("/agent-list-assigned-to-facebook-application/:flac_uuid", handler.AgentListAssignedToFacebookApplication)
	e.GET("/agent-list-not-in-facebook-application/:flac_uuid", handler.AgentListNotInFacebookApplication)
	e.GET("/show-agent-facebook-application/:agent_uuid", handler.ShowAgentFacebookApplication)
}
