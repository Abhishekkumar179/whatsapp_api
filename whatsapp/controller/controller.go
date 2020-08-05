package controller

import (
	"bytes"
	"context"
	json "encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Get_allId(ctx)

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
	appId := ("appId")
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

/***********************************************Save Tenant_Details***************************************/
func (r *CRUDController) Save_Tenant_details(c echo.Context) error {
	var tenantdetails map[string]interface{}
	c.Bind(&tenantdetails)
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Save_Tenant_details(ctx, tenantdetails)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)

}

/*********************************************** Get Tenant appId******************************************/
func (r *CRUDController) Get_Tenant_appId(c echo.Context) error {
	domain_uuid := c.Param("domain_uuid")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Get_Tenant_appId(ctx, domain_uuid)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)

}

/**********************************************Update tenant details**************************************/
func (r *CRUDController) Update_Tenant_details(c echo.Context) error {
	domain_uuid := c.Param("domain_uuid")
	var tenantdetails map[string]interface{}
	err1 := json.NewDecoder(c.Request().Body).Decode(&tenantdetails)
	if err1 != nil {
		fmt.Println("err= ", err1)
	} else {
		fmt.Println("err= ", err1)
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Update_Tenant_details(ctx, domain_uuid, tenantdetails)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)

}

/***********************************************Delete Tenant details**************************************/
func (r *CRUDController) Delete_Tenant_details(c echo.Context) error {
	domain_uuid := c.Param("domain_uuid")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.usecase.Delete_Tenant_details(ctx, domain_uuid)

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
	err := c.Request().ParseMultipartForm(25 << 20) // 25Mb
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
	authResponse, _ := r.usecase.Upload_Attachments(ctx, appId, appUserId, Type, file, handler)

	if authResponse == nil {
		return c.JSONBlob(http.StatusUnauthorized, authResponse)
	}
	return c.JSONBlob(http.StatusOK, authResponse)

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
	e.GET("/getall_appUserId/:appId", handler.Get_allId)
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
	e.PUT("update_appUser/:appId/:appUserId", handler.Update_AppUser)
	e.POST("send_messages/:appId/:appUserId", handler.Message_Action_Types)
	e.POST("quickreply_message/:appId/:appUserId", handler.Quickreply_Message)
	e.POST("send_carousel_message/:appId/:appUserId", handler.Send_Carousel_Message)
	e.POST("save_tenant_details", handler.Save_Tenant_details)
	e.GET("get_tenant_appId/:domain_uuid", handler.Get_Tenant_appId)
	e.POST("update_tenant_details/:domain_uuid", handler.Update_Tenant_details)
	e.DELETE("delete_tenant_details/:domain_uuid", handler.Delete_Tenant_details)
	e.GET("list_integration/:appId", handler.List_integration)
	e.POST("link_appUser_to_channel/:appId/:appUserId", handler.Link_appUser_to_Channel)
	e.DELETE("unlink_appUser_to_channel/:appId/:appUserId/:channel", handler.Unlink_appUser_to_Channel)
	e.POST("upload_attachments/:appId/:appUserId", handler.Upload_Attachments)
	e.POST("typing_activity/:appId/:appUserId", handler.TypingActivity)
}
