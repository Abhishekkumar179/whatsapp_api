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
	"strconv"
	"strings"
	"time"
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

/**************************************************Delete AppUser Profile****************************************/
func (r *crudRepository) Delete_AppUser_Profile(ctx context.Context, appId string, appUserId string) (*models.Response, error) {

	td := models.Tenant_details{}
	db := r.DBConn.Where("app_id = ?", appId).Find(&td)
	if db.Error != nil {
		return nil, db.Error
	}
	res, err := http.NewRequest("DELETE", "https://api.smooch.io/v1.1/apps/"+appId+"/appusers/"+appUserId+"/profile", nil)
	res.Header.Set("Content-Type", "application/json")
	res.SetBasicAuth(td.AppKey, td.AppSecret)

	client := &http.Client{}
	response, err := client.Do(res)
	if err != nil {
		fmt.Printf("error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
		us := models.Appuser{}
		re := models.ReceiveUserDetails{}
		st := r.DBConn.Where("id = ?", appUserId).Delete(&us)
		if st.Error != nil {
		}
		db := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", appUserId).Find(&re)
		if db.Error != nil {
			return &models.Response{ResponseCode: 409, Status: "Error", Msg: "UserId Not Found"}, nil
		}
		del := r.DBConn.Where("app_user_id = ?", appUserId).Delete(&re)
		if del.RowsAffected == 0 {
			return &models.Response{ResponseCode: 404, Status: "Not Deleted", Msg: "User Profile  Not Deleted."}, nil
		}
		return &models.Response{ResponseCode: 201, Status: "OK", Msg: "User Profile Deleted successfully."}, nil

	}

	return &models.Response{ResponseCode: 404, Status: "Failed", Msg: "User NOT Found1."}, nil

}

/**************************************************Getall Id***************************************************/

func (r *crudRepository) Get_allId(ctx context.Context) (*models.Response, error) {
	list := make([]models.ReceiveUserDetails, 0)

	if rows, err := r.DBConn.Raw("select app_id, app_user_id, surname, given_name,type,text,role,name,author_id,message_id,original_message_id,integration_id,source_type, signed_up_at, conversation_started from receive_user_details where is_enabled = true").Rows(); err != nil {

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

/*************************************************Get AppUser Details*****************************************/
func (r *crudRepository) GetAppUserDetails(ctx context.Context, appUserId string, appId string) ([]byte, error) {
	td := models.Tenant_details{}
	db := r.DBConn.Where("app_id = ?", appId).Find(&td)
	if db.Error != nil {

	}
	res, err := http.NewRequest("GET", "https://api.smooch.io/v1.1/apps/"+appId+"/appusers/"+appUserId, nil)
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
	w := models.WhatsappConfiguration{}
	fb := models.FacebookConfiguration{}
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
		Is_enabled:               true,
	}
	s := int64(f.Messages[0].Received)
	myDate := time.Unix(s, 0)
	fmt.Println(myDate)
	hour := strconv.Itoa(myDate.Hour())
	err := r.DBConn.Where("app_user_id = ?", f.AppUser.ID).Find(&u)
	fmt.Println(err)
	if u.Is_enabled == false {
		update := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Update("is_enabled", true)
		fmt.Println(update, update.RowsAffected)
	}

	if f.Messages[0].Source.Type == "messenger" {
		db := r.DBConn.Where("facebook_integration_id = ?", f.Messages[0].Source.IntegrationID).Find(&fb)
		if db.Error != nil {
			fmt.Println("error")
		}
		if myDate.Weekday().String() == fb.Day1 {
			workstart1 := fb.Workstart1
			components := strings.Split(workstart1, ":")
			StartHour, _ := components[0], components[1]
			workend1 := fb.Workend1
			components1 := strings.Split(workend1, ":")
			EndHour, _ := components1[0], components1[1]
			if hour < StartHour || hour > EndHour {
				p := models.User{
					Role: "appMaker",
					Type: "text",
					Text: fb.Message,
				}
				r.PostMessage(ctx, fb.AppId, f.AppUser.ID, p)
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}

					return &models.Response{Received: &f}, nil
				}
				return &models.Response{Msg: "Userid already exist."}, nil
			} else {

				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}
					p := models.User{
						Role: "appMaker",
						Type: "text",
						Text: fb.TriggerMessage,
					}
					r.PostMessage(ctx, fb.AppId, f.AppUser.ID, p)
					return &models.Response{Received: &f}, nil
				}
				fmt.Println("appUserId already exist.")
				return &models.Response{Msg: "Userid already exist."}, nil

			}
		} else if myDate.Weekday().String() == fb.Day2 {
			workstart2 := fb.Workstart2
			components := strings.Split(workstart2, ":")
			StartHour, _ := components[0], components[1]
			workend2 := fb.Workend2
			components1 := strings.Split(workend2, ":")
			EndHour, _ := components1[0], components1[1]
			if hour < StartHour || hour > EndHour {
				p := models.User{
					Role: "appMaker",
					Type: "text",
					Text: fb.Message,
				}
				r.PostMessage(ctx, fb.AppId, f.AppUser.ID, p)
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}

					return &models.Response{Received: &f}, nil
				}
				return &models.Response{Msg: "Userid already exist."}, nil
			} else {

				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}
					p := models.User{
						Role: "appMaker",
						Type: "text",
						Text: fb.TriggerMessage,
					}
					r.PostMessage(ctx, fb.AppId, f.AppUser.ID, p)
					return &models.Response{Received: &f}, nil
				}
				fmt.Println("appUserId already exist.")
				return &models.Response{Msg: "Userid already exist."}, nil

			}
		} else if myDate.Weekday().String() == fb.Day3 {
			workstart3 := fb.Workstart3
			components := strings.Split(workstart3, ":")
			StartHour, _ := components[0], components[1]
			workend3 := fb.Workend3
			components1 := strings.Split(workend3, ":")
			EndHour, _ := components1[0], components1[1]
			if hour < StartHour || hour > EndHour {
				p := models.User{
					Role: "appMaker",
					Type: "text",
					Text: fb.Message,
				}
				r.PostMessage(ctx, fb.AppId, f.AppUser.ID, p)
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}

					return &models.Response{Received: &f}, nil
				}
				return &models.Response{Msg: "Userid already exist."}, nil
			} else {

				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}
					p := models.User{
						Role: "appMaker",
						Type: "text",
						Text: fb.TriggerMessage,
					}
					r.PostMessage(ctx, fb.AppId, f.AppUser.ID, p)
					return &models.Response{Received: &f}, nil
				}
				fmt.Println("appUserId already exist.")
				return &models.Response{Msg: "Userid already exist."}, nil

			}
		} else if myDate.Weekday().String() == fb.Day4 {
			workstart4 := fb.Workstart4
			components := strings.Split(workstart4, ":")
			StartHour, _ := components[0], components[1]
			workend4 := fb.Workend4
			components1 := strings.Split(workend4, ":")
			EndHour, _ := components1[0], components1[1]
			if hour < StartHour || hour > EndHour {
				p := models.User{
					Role: "appMaker",
					Type: "text",
					Text: fb.Message,
				}
				r.PostMessage(ctx, fb.AppId, f.AppUser.ID, p)
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}

					return &models.Response{Received: &f}, nil
				}
				return &models.Response{Msg: "Userid already exist."}, nil
			} else {

				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}
					p := models.User{
						Role: "appMaker",
						Type: "text",
						Text: fb.TriggerMessage,
					}
					r.PostMessage(ctx, fb.AppId, f.AppUser.ID, p)
					return &models.Response{Received: &f}, nil
				}
				fmt.Println("appUserId already exist.")
				return &models.Response{Msg: "Userid already exist."}, nil

			}
		} else if myDate.Weekday().String() == fb.Day5 {
			workstart5 := fb.Workstart5
			components := strings.Split(workstart5, ":")
			StartHour, _ := components[0], components[1]
			workend5 := fb.Workend5
			components1 := strings.Split(workend5, ":")
			EndHour, _ := components1[0], components1[1]
			if hour < StartHour || hour > EndHour {
				p := models.User{
					Role: "appMaker",
					Type: "text",
					Text: fb.Message,
				}
				r.PostMessage(ctx, fb.AppId, f.AppUser.ID, p)
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}

					return &models.Response{Received: &f}, nil
				}
				return &models.Response{Msg: "Userid already exist."}, nil
			} else {

				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}
					p := models.User{
						Role: "appMaker",
						Type: "text",
						Text: fb.TriggerMessage,
					}
					r.PostMessage(ctx, fb.AppId, f.AppUser.ID, p)
					return &models.Response{Received: &f}, nil
				}
				fmt.Println("appUserId already exist.")
				return &models.Response{Msg: "Userid already exist."}, nil

			}
		} else if myDate.Weekday().String() == fb.Day6 {
			workstart6 := fb.Workstart6
			components := strings.Split(workstart6, ":")
			StartHour, _ := components[0], components[1]
			workend6 := fb.Workend6
			components1 := strings.Split(workend6, ":")
			EndHour, _ := components1[0], components1[1]
			if hour < StartHour || hour > EndHour {
				p := models.User{
					Role: "appMaker",
					Type: "text",
					Text: fb.Message,
				}
				r.PostMessage(ctx, fb.AppId, f.AppUser.ID, p)
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}

					return &models.Response{Received: &f}, nil
				}
				return &models.Response{Msg: "Userid already exist."}, nil
			} else {

				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}
					p := models.User{
						Role: "appMaker",
						Type: "text",
						Text: fb.TriggerMessage,
					}
					r.PostMessage(ctx, fb.AppId, f.AppUser.ID, p)
					return &models.Response{Received: &f}, nil
				}
				fmt.Println("appUserId already exist.")
				return &models.Response{Msg: "Userid already exist."}, nil

			}
		} else if myDate.Weekday().String() == fb.Day7 {
			workstart7 := fb.Workstart7
			components := strings.Split(workstart7, ":")
			StartHour, _ := components[0], components[1]
			workend7 := fb.Workend7
			components1 := strings.Split(workend7, ":")
			EndHour, _ := components1[0], components1[1]
			if hour < StartHour || hour > EndHour {
				p := models.User{
					Role: "appMaker",
					Type: "text",
					Text: fb.Message,
				}
				r.PostMessage(ctx, fb.AppId, f.AppUser.ID, p)
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}

					return &models.Response{Received: &f}, nil
				}
				return &models.Response{Msg: "Userid already exist."}, nil
			} else {

				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}
					p := models.User{
						Role: "appMaker",
						Type: "text",
						Text: fb.TriggerMessage,
					}
					r.PostMessage(ctx, fb.AppId, f.AppUser.ID, p)
					return &models.Response{Received: &f}, nil
				}
				fmt.Println("appUserId already exist.")
				return &models.Response{Msg: "Userid already exist."}, nil

			}
		} else {
			fmt.Println("User Registered.")
		}
	} else if f.Messages[0].Source.Type == "whatsapp" {
		db := r.DBConn.Where("whatsapp_integration_id = ?", f.Messages[0].Source.IntegrationID).Find(&w)
		if db.Error != nil {
			fmt.Println("error")
		}
		if myDate.Weekday().String() == w.Day1 {
			workstart1 := w.Workstart1
			components := strings.Split(workstart1, ":")
			StartHour, _ := components[0], components[1]
			workend1 := w.Workend1
			components1 := strings.Split(workend1, ":")
			EndHour, _ := components1[0], components1[1]
			if hour < StartHour || hour > EndHour {
				p := models.User{
					Role: "appMaker",
					Type: "text",
					Text: w.Message,
				}
				r.PostMessage(ctx, w.AppId, f.AppUser.ID, p)
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}

					return &models.Response{Received: &f}, nil
				}
				fmt.Println("appUserId already exist.")
				return &models.Response{Msg: "Userid already exist."}, nil

			} else {
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}

					p := models.User{
						Role: "appMaker",
						Type: "text",
						Text: w.TriggerMessage,
					}
					r.PostMessage(ctx, w.AppId, f.AppUser.ID, p)
					return &models.Response{Received: &f}, nil
				}
				fmt.Println("appUserId already exist.")
				return &models.Response{Msg: "Userid already exist."}, nil

			}
		} else if myDate.Weekday().String() == w.Day2 {
			workstart2 := w.Workstart2
			components := strings.Split(workstart2, ":")
			StartHour, _ := components[0], components[1]
			workend2 := w.Workend2
			components1 := strings.Split(workend2, ":")
			EndHour, _ := components1[0], components1[1]
			if hour < StartHour || hour > EndHour {
				p := models.User{
					Role: "appMaker",
					Type: "text",
					Text: w.Message,
				}
				r.PostMessage(ctx, w.AppId, f.AppUser.ID, p)
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}

					return &models.Response{Received: &f}, nil
				}
				fmt.Println("appUserId already exist.")
				return &models.Response{Msg: "Userid already exist."}, nil

			} else {
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}

					p := models.User{
						Role: "appMaker",
						Type: "text",
						Text: w.TriggerMessage,
					}
					r.PostMessage(ctx, w.AppId, f.AppUser.ID, p)
					return &models.Response{Received: &f}, nil
				}
				fmt.Println("appUserId already exist.")
				return &models.Response{Msg: "Userid already exist."}, nil

			}
		} else if myDate.Weekday().String() == w.Day3 {
			workstart3 := w.Workstart3
			components := strings.Split(workstart3, ":")
			StartHour, _ := components[0], components[1]
			workend3 := w.Workend3
			components1 := strings.Split(workend3, ":")
			EndHour, _ := components1[0], components1[1]
			if hour < StartHour || hour > EndHour {
				p := models.User{
					Role: "appMaker",
					Type: "text",
					Text: w.Message,
				}
				r.PostMessage(ctx, w.AppId, f.AppUser.ID, p)
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}

					return &models.Response{Received: &f}, nil
				}
				fmt.Println("appUserId already exist.")
				return &models.Response{Msg: "Userid already exist."}, nil

			} else {
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}

					p := models.User{
						Role: "appMaker",
						Type: "text",
						Text: w.TriggerMessage,
					}
					r.PostMessage(ctx, w.AppId, f.AppUser.ID, p)
					return &models.Response{Received: &f}, nil
				}
				fmt.Println("appUserId already exist.")
				return &models.Response{Msg: "Userid already exist."}, nil

			}
		} else if myDate.Weekday().String() == w.Day4 {
			workstart4 := w.Workstart4
			components := strings.Split(workstart4, ":")
			StartHour, _ := components[0], components[1]
			workend4 := w.Workend4
			components1 := strings.Split(workend4, ":")
			EndHour, _ := components1[0], components1[1]
			if hour < StartHour || hour > EndHour {
				p := models.User{
					Role: "appMaker",
					Type: "text",
					Text: w.Message,
				}
				r.PostMessage(ctx, w.AppId, f.AppUser.ID, p)
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}
					return &models.Response{Received: &f}, nil
				}
				fmt.Println("appUserId already exist.")
				return &models.Response{Msg: "Userid already exist."}, nil

			} else {
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}

					p := models.User{
						Role: "appMaker",
						Type: "text",
						Text: w.TriggerMessage,
					}
					r.PostMessage(ctx, w.AppId, f.AppUser.ID, p)
					return &models.Response{Received: &f}, nil
				}
				fmt.Println("appUserId already exist.")
				return &models.Response{Msg: "Userid already exist."}, nil

			}
		} else if myDate.Weekday().String() == w.Day5 {
			workstart5 := w.Workstart5
			components := strings.Split(workstart5, ":")
			StartHour, _ := components[0], components[1]
			workend5 := w.Workend5
			components1 := strings.Split(workend5, ":")
			EndHour, _ := components1[0], components1[1]
			if hour < StartHour || hour > EndHour {
				p := models.User{
					Role: "appMaker",
					Type: "text",
					Text: w.Message,
				}
				r.PostMessage(ctx, w.AppId, f.AppUser.ID, p)
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}

					return &models.Response{Received: &f}, nil
				}
				fmt.Println("appUserId already exist.")
				return &models.Response{Msg: "Userid already exist."}, nil

			} else {
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}

					p := models.User{
						Role: "appMaker",
						Type: "text",
						Text: w.TriggerMessage,
					}
					r.PostMessage(ctx, w.AppId, f.AppUser.ID, p)
					return &models.Response{Received: &f}, nil
				}
				fmt.Println("appUserId already exist.")
				return &models.Response{Msg: "Userid already exist."}, nil

			}
		} else if myDate.Weekday().String() == w.Day6 {
			workstart6 := w.Workstart6
			components := strings.Split(workstart6, ":")
			StartHour, _ := components[0], components[1]
			workend6 := w.Workend6
			components1 := strings.Split(workend6, ":")
			EndHour, _ := components1[0], components1[1]
			if hour < StartHour || hour > EndHour {
				p := models.User{
					Role: "appMaker",
					Type: "text",
					Text: w.Message,
				}
				r.PostMessage(ctx, w.AppId, f.AppUser.ID, p)
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}
					return &models.Response{Received: &f}, nil
				}
				fmt.Println("appUserId already exist.")
				return &models.Response{Msg: "Userid already exist."}, nil

			} else {
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}

					p := models.User{
						Role: "appMaker",
						Type: "text",
						Text: w.TriggerMessage,
					}
					r.PostMessage(ctx, w.AppId, f.AppUser.ID, p)
					return &models.Response{Received: &f}, nil
				}
				fmt.Println("appUserId already exist.")
				return &models.Response{Msg: "Userid already exist."}, nil

			}
		} else if myDate.Weekday().String() == w.Day7 {
			workstart7 := w.Workstart7
			components := strings.Split(workstart7, ":")
			StartHour, _ := components[0], components[1]
			workend7 := w.Workend7
			components1 := strings.Split(workend7, ":")
			EndHour, _ := components1[0], components1[1]
			if hour < StartHour || hour > EndHour {
				p := models.User{
					Role: "appMaker",
					Type: "text",
					Text: w.Message,
				}
				r.PostMessage(ctx, w.AppId, f.AppUser.ID, p)
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}
					return &models.Response{Received: &f}, nil
				}
				fmt.Println("appUserId already exist.")
				return &models.Response{Msg: "Userid already exist."}, nil

			} else {
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}
					p := models.User{
						Role: "appMaker",
						Type: "text",
						Text: w.TriggerMessage,
					}
					r.PostMessage(ctx, w.AppId, f.AppUser.ID, p)
					return &models.Response{Received: &f}, nil
				}
				fmt.Println("appUserId already exist.")
				return &models.Response{Msg: "Userid already exist."}, nil

			}
		} else {
			fmt.Println("user Registered.")
		}

	} else {
		return &models.Response{Msg: "Userid already exist."}, nil
	}
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
	s := int64(1597047140)
	myDate := time.Unix(s, 0)

	fmt.Println(myDate.Hour(), myDate.Weekday())

	w := models.Tenant_details{}
	fmt.Println(myDate, myDate.Hour(), myDate.Day(), myDate.Weekday(), "vjhwvhj")
	list := make([]models.Tenant_details, 0)
	db := r.DBConn.Where("domain_uuid = ?", domain_uuid).Find(&w)
	if db.Error != nil {
		return &models.Response{Status: "0", Msg: "Record Not Found", ResponseCode: 401}, nil

	}
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
func (r crudRepository) Add_Whatsapp_configuration(ctx context.Context, td models.WhatsappConfigurations) (*models.Response, error) {

	ts := models.WhatsappConfiguration{
		Domain_uuid:           td.Domain_uuid,
		AppId:                 td.AppId,
		AppKey:                td.AppKey,
		AppSecret:             td.AppSecret,
		Message:               td.Message,
		Size:                  td.Size,
		TriggerWhen:           td.Trigger.When,
		TriggerName:           td.Trigger.Name,
		TriggerMessage:        td.Trigger.Message,
		WhatsappIntegrationID: td.WhatsappIntegrationID,
		Day1:                  td.WorkingDays[0].Day,
		Day2:                  td.WorkingDays[1].Day,
		Day3:                  td.WorkingDays[2].Day,
		Day4:                  td.WorkingDays[3].Day,
		Day5:                  td.WorkingDays[4].Day,
		Day6:                  td.WorkingDays[5].Day,
		Day7:                  td.WorkingDays[6].Day,
		Workstart1:            td.WorkingDays[0].WorkingHourStartTime,
		Workstart2:            td.WorkingDays[1].WorkingHourStartTime,
		Workstart3:            td.WorkingDays[2].WorkingHourStartTime,
		Workstart4:            td.WorkingDays[3].WorkingHourStartTime,
		Workstart5:            td.WorkingDays[4].WorkingHourStartTime,
		Workstart6:            td.WorkingDays[5].WorkingHourStartTime,
		Workstart7:            td.WorkingDays[6].WorkingHourStartTime,
		Workend1:              td.WorkingDays[0].WorkingHourEndTime,
		Workend2:              td.WorkingDays[1].WorkingHourEndTime,
		Workend3:              td.WorkingDays[2].WorkingHourEndTime,
		Workend4:              td.WorkingDays[3].WorkingHourEndTime,
		Workend5:              td.WorkingDays[4].WorkingHourEndTime,
		Workend6:              td.WorkingDays[5].WorkingHourEndTime,
		Workend7:              td.WorkingDays[6].WorkingHourEndTime,
	}
	if err := r.DBConn.Table("whatsapp_configurations").Where("app_id = ?", td.AppId).Find(&ts).Error; err != nil {
		if db := r.DBConn.Table("whatsapp_configurations").Where("whatsapp_integration_id = ?", td.WhatsappIntegrationID).Find(&ts).Error; db != nil {
			st := r.DBConn.Create(&ts)
			if st.Error != nil {
				return &models.Response{ResponseCode: 409, Status: "Error", Msg: "Not created"}, nil
			}
			return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Whatsapp configuration saved successfully."}, nil

		}
		return &models.Response{ResponseCode: 409, Status: "Error", Msg: "Whatsapp Integration Id  Already Exist."}, nil
	}
	return &models.Response{ResponseCode: 409, Status: "Error", Msg: "AppId Already Exist."}, nil
}

/**********************************************Get appID by tenant_domain_uuid******************************/
func (r crudRepository) Get_Whatsapp_configuration(ctx context.Context, domain_uuid string) (*models.Response, error) {
	w := models.WhatsappConfiguration{}
	list := make([]models.WhatsappConfiguration, 0)
	db := r.DBConn.Where("domain_uuid = ?", domain_uuid).Find(&w)
	if db.Error != nil {
		return &models.Response{Status: "0", Msg: "Record Not Found", ResponseCode: 401}, nil

	}

	row, err := r.DBConn.Raw("select id, domain_uuid, app_id, app_key, app_secret, whatsapp_integration_id,size,trigger_name,trigger_message,trigger_when, day1, day2, day3, day4, day5, day6, day7, workstart1, workstart2, workstart3, workstart4, workstart5, workstart6, workstart7, workend1, workend2, workend3, workend4, workend5, workend6, workend7 from whatsapp_configurations WHERE domain_uuid = ?", domain_uuid).Rows()
	if err != nil {
		return &models.Response{Status: "0", Msg: "Record Not Found", ResponseCode: 401}, nil
	}
	defer row.Close()
	for row.Next() {
		f := models.WhatsappConfiguration{}
		if err := row.Scan(&f.Id, &f.Domain_uuid, &f.AppId, &f.AppKey, &f.AppSecret, &f.WhatsappIntegrationID, &f.Size, &f.TriggerName, &f.TriggerMessage, &f.TriggerWhen, &f.Day1, &f.Day2, &f.Day3, &f.Day4, &f.Day5, &f.Day6, &f.Day7, &f.Workstart1, &f.Workstart2, &f.Workstart3, &f.Workstart4, &f.Workstart5, &f.Workstart6, &f.Workstart7, &f.Workend1, &f.Workend2, &f.Workend3, &f.Workend4, &f.Workend5, &f.Workend6, &f.Workend7); err != nil {

			return nil, err
		}

		list = append(list, f)
	}
	return &models.Response{Status: "1", Msg: "Record Found", ResponseCode: 200, List: list}, nil
}

/**********************************************Update Tenant details*************************************/
func (r crudRepository) Update_Whatsapp_configuration(ctx context.Context, id int64, domain_uuid string, td models.WhatsappConfigurations) (*models.Response, error) {
	w := models.WhatsappConfiguration{}
	db1 := r.DBConn.Where("domain_uuid = ? AND id = ?", domain_uuid, id).Find(&w)
	if db1.Error != nil {
		return &models.Response{Status: "Not Found", Msg: "Record Doesn't Exist", ResponseCode: 401}, nil
	}
	if td.WhatsappIntegrationID == w.WhatsappIntegrationID && td.AppId == w.AppId {
		fmt.Println("part1")
		if db := r.DBConn.Table("whatsapp_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"app_key": td.AppKey, "app_secret": td.AppSecret, "message": td.Message, "size": td.Size, "trigger_name": td.Trigger.Name, "trigger_message": td.Trigger.Message, "trigger_when": td.Trigger.When, "day1": td.WorkingDays[0].Day, "day2": td.WorkingDays[1].Day, "day3": td.WorkingDays[2].Day, "day4": td.WorkingDays[3].Day, "day5": td.WorkingDays[4].Day, "day6": td.WorkingDays[5].Day, "day7": td.WorkingDays[6].Day, "workstart1": td.WorkingDays[0].WorkingHourStartTime, "workstart2": td.WorkingDays[1].WorkingHourStartTime, "workstart3": td.WorkingDays[2].WorkingHourStartTime, "workstart4": td.WorkingDays[3].WorkingHourStartTime, "workstart5": td.WorkingDays[4].WorkingHourStartTime, "workstart6": td.WorkingDays[5].WorkingHourStartTime, "workstart7": td.WorkingDays[6].WorkingHourStartTime, "workend1": td.WorkingDays[0].WorkingHourEndTime, "workend2": td.WorkingDays[1].WorkingHourEndTime, "workend3": td.WorkingDays[2].WorkingHourEndTime, "workend4": td.WorkingDays[3].WorkingHourEndTime, "workend5": td.WorkingDays[4].WorkingHourEndTime, "workend6": td.WorkingDays[5].WorkingHourEndTime, "workend7": td.WorkingDays[6].WorkingHourEndTime}).Error; db != nil {
			return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
		}

		return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Whatsapp integration Updated successfully."}, nil
	} else if td.WhatsappIntegrationID != w.WhatsappIntegrationID && td.AppId != w.AppId {
		fmt.Println("part2")
		w := models.WhatsappConfiguration{}
		if row := r.DBConn.Table("whatsapp_configurations").Where("app_id = ?", td.AppId).Find(&w).Error; row != nil {
			if err := r.DBConn.Table("whatsapp_configurations").Where("whatsapp_integration_id = ?", td.WhatsappIntegrationID).Find(&w).Error; err != nil {
				if db := r.DBConn.Table("whatsapp_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"app_id": td.AppId, "app_key": td.AppKey, "app_secret": td.AppSecret, "message": td.Message, "whatsapp_integration_id": td.WhatsappIntegrationID, "size": td.Size, "trigger_name": td.Trigger.Name, "trigger_message": td.Trigger.Message, "trigger_when": td.Trigger.When, "day1": td.WorkingDays[0].Day, "day2": td.WorkingDays[1].Day, "day3": td.WorkingDays[2].Day, "day4": td.WorkingDays[3].Day, "day5": td.WorkingDays[4].Day, "day6": td.WorkingDays[5].Day, "day7": td.WorkingDays[6].Day, "workstart1": td.WorkingDays[0].WorkingHourStartTime, "workstart2": td.WorkingDays[1].WorkingHourStartTime, "workstart3": td.WorkingDays[2].WorkingHourStartTime, "workstart4": td.WorkingDays[3].WorkingHourStartTime, "workstart5": td.WorkingDays[4].WorkingHourStartTime, "workstart6": td.WorkingDays[5].WorkingHourStartTime, "workstart7": td.WorkingDays[6].WorkingHourStartTime, "workend1": td.WorkingDays[0].WorkingHourEndTime, "workend2": td.WorkingDays[1].WorkingHourEndTime, "workend3": td.WorkingDays[2].WorkingHourEndTime, "workend4": td.WorkingDays[3].WorkingHourEndTime, "workend5": td.WorkingDays[4].WorkingHourEndTime, "workend6": td.WorkingDays[5].WorkingHourEndTime, "workend7": td.WorkingDays[6].WorkingHourEndTime}).Error; db != nil {
					return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
				}

				return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Whatsapp integration Updated successfully."}, nil
			}
			return &models.Response{ResponseCode: 409, Status: "0", Msg: "Whatsapp id already exist."}, nil

		}
		return &models.Response{ResponseCode: 409, Status: "0", Msg: "AppId already exist."}, nil

	} else if td.WhatsappIntegrationID == w.WhatsappIntegrationID && td.AppId != w.AppId {
		fmt.Println("part3")
		w := models.WhatsappConfiguration{}
		if row := r.DBConn.Where("app_id = ?", td.AppId).Find(&w).Error; row != nil {
			if db := r.DBConn.Table("whatsapp_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"app_id": td.AppId, "app_key": td.AppKey, "app_secret": td.AppSecret, "message": td.Message, "size": td.Size, "trigger_name": td.Trigger.Name, "trigger_message": td.Trigger.Message, "trigger_when": td.Trigger.When, "day1": td.WorkingDays[0].Day, "day2": td.WorkingDays[1].Day, "day3": td.WorkingDays[2].Day, "day4": td.WorkingDays[3].Day, "day5": td.WorkingDays[4].Day, "day6": td.WorkingDays[5].Day, "day7": td.WorkingDays[6].Day, "workstart1": td.WorkingDays[0].WorkingHourStartTime, "workstart2": td.WorkingDays[1].WorkingHourStartTime, "workstart3": td.WorkingDays[2].WorkingHourStartTime, "workstart4": td.WorkingDays[3].WorkingHourStartTime, "workstart5": td.WorkingDays[4].WorkingHourStartTime, "workstart6": td.WorkingDays[5].WorkingHourStartTime, "workstart7": td.WorkingDays[6].WorkingHourStartTime, "workend1": td.WorkingDays[0].WorkingHourEndTime, "workend2": td.WorkingDays[1].WorkingHourEndTime, "workend3": td.WorkingDays[2].WorkingHourEndTime, "workend4": td.WorkingDays[3].WorkingHourEndTime, "workend5": td.WorkingDays[4].WorkingHourEndTime, "workend6": td.WorkingDays[5].WorkingHourEndTime, "workend7": td.WorkingDays[6].WorkingHourEndTime}).Error; db != nil {
				return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
			}

			return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Whatsapp integration Updated successfully."}, nil
		}
		return &models.Response{ResponseCode: 409, Status: "0", Msg: "AppId already exist."}, nil

	} else if td.WhatsappIntegrationID != w.WhatsappIntegrationID && td.AppId == w.AppId {
		fmt.Println("part4")
		w := models.WhatsappConfiguration{}
		if err := r.DBConn.Where("whatsapp_integration_id = ?", td.WhatsappIntegrationID).Find(&w).Error; err != nil {
			if db := r.DBConn.Table("whatsapp_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"app_key": td.AppKey, "app_secret": td.AppSecret, "message": td.Message, "whatsapp_integration_id": td.WhatsappIntegrationID, "size": td.Size, "trigger_name": td.Trigger.Name, "trigger_message": td.Trigger.Message, "trigger_when": td.Trigger.When, "day1": td.WorkingDays[0].Day, "day2": td.WorkingDays[1].Day, "day3": td.WorkingDays[2].Day, "day4": td.WorkingDays[3].Day, "day5": td.WorkingDays[4].Day, "day6": td.WorkingDays[5].Day, "day7": td.WorkingDays[6].Day, "workstart1": td.WorkingDays[0].WorkingHourStartTime, "workstart2": td.WorkingDays[1].WorkingHourStartTime, "workstart3": td.WorkingDays[2].WorkingHourStartTime, "workstart4": td.WorkingDays[3].WorkingHourStartTime, "workstart5": td.WorkingDays[4].WorkingHourStartTime, "workstart6": td.WorkingDays[5].WorkingHourStartTime, "workstart7": td.WorkingDays[6].WorkingHourStartTime, "workend1": td.WorkingDays[0].WorkingHourEndTime, "workend2": td.WorkingDays[1].WorkingHourEndTime, "workend3": td.WorkingDays[2].WorkingHourEndTime, "workend4": td.WorkingDays[3].WorkingHourEndTime, "workend5": td.WorkingDays[4].WorkingHourEndTime, "workend6": td.WorkingDays[5].WorkingHourEndTime, "workend7": td.WorkingDays[6].WorkingHourEndTime}).Error; db != nil {
				return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
			}

			return &models.Response{ResponseCode: 200, Status: "OK", Msg: "Whatsapp integration Updated successfully."}, nil
		}
		return &models.Response{ResponseCode: 409, Status: "0", Msg: "Whatsapp id already exist."}, nil
	} else {
		return &models.Response{ResponseCode: 200, Status: "OK", Msg: "Whatsapp integration Updated successfully."}, nil

	}

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
func (r crudRepository) Add_Facebook_configuration(ctx context.Context, td models.FacebookConfigurations) (*models.Response, error) {

	ts := models.FacebookConfiguration{
		Domain_uuid:           td.Domain_uuid,
		AppId:                 td.AppId,
		AppKey:                td.AppKey,
		AppSecret:             td.AppSecret,
		Message:               td.Message,
		FacebookIntegrationID: td.FacebookIntegrationID,
		Size:                  td.Size,
		TriggerWhen:           td.Trigger.When,
		TriggerName:           td.Trigger.Name,
		TriggerMessage:        td.Trigger.Message,
		Day1:                  td.WorkingDays[0].Day,
		Day2:                  td.WorkingDays[1].Day,
		Day3:                  td.WorkingDays[2].Day,
		Day4:                  td.WorkingDays[3].Day,
		Day5:                  td.WorkingDays[4].Day,
		Day6:                  td.WorkingDays[5].Day,
		Day7:                  td.WorkingDays[6].Day,
		Workstart1:            td.WorkingDays[0].WorkingHourStartTime,
		Workstart2:            td.WorkingDays[1].WorkingHourStartTime,
		Workstart3:            td.WorkingDays[2].WorkingHourStartTime,
		Workstart4:            td.WorkingDays[3].WorkingHourStartTime,
		Workstart5:            td.WorkingDays[4].WorkingHourStartTime,
		Workstart6:            td.WorkingDays[5].WorkingHourStartTime,
		Workstart7:            td.WorkingDays[6].WorkingHourStartTime,
		Workend1:              td.WorkingDays[0].WorkingHourEndTime,
		Workend2:              td.WorkingDays[1].WorkingHourEndTime,
		Workend3:              td.WorkingDays[2].WorkingHourEndTime,
		Workend4:              td.WorkingDays[3].WorkingHourEndTime,
		Workend5:              td.WorkingDays[4].WorkingHourEndTime,
		Workend6:              td.WorkingDays[5].WorkingHourEndTime,
		Workend7:              td.WorkingDays[6].WorkingHourEndTime,
	}
	if err := r.DBConn.Table("facebook_configurations").Where("app_id = ?", td.AppId).Find(&ts).Error; err != nil {
		if db := r.DBConn.Table("facebook_configurations").Where("facebook_integration_id = ?", td.FacebookIntegrationID).Find(&ts).Error; db != nil {
			st := r.DBConn.Create(&ts)
			if st.Error != nil {
				return &models.Response{ResponseCode: 409, Status: "Error", Msg: "Not created"}, nil
			}
			return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Facebook configuration saved successfully."}, nil

		}
		return &models.Response{ResponseCode: 409, Status: "Error", Msg: "Facebook Integration Id  Already Exist."}, nil
	}
	return &models.Response{ResponseCode: 409, Status: "Error", Msg: "AppId  Already Exist."}, nil
}

/**********************************************Get appID by tenant_domain_uuid******************************/
func (r crudRepository) Get_Facebook_configuration(ctx context.Context, domain_uuid string) (*models.Response, error) {
	list := make([]models.FacebookConfiguration, 0)
	w := models.FacebookConfiguration{}
	db := r.DBConn.Where("domain_uuid = ?", domain_uuid).Find(&w)
	if db.Error != nil {
		return &models.Response{Status: "0", Msg: "Record Not Found", ResponseCode: 401}, nil

	}
	row, err := r.DBConn.Raw("select id, domain_uuid, app_id, app_key, app_secret, facebook_integration_id, size, trigger_name, trigger_message, trigger_when, day1, day2, day3, day4, day5, day6, day7, workstart1, workstart2, workstart3, workstart4, workstart5, workstart6, workstart7, workend1, workend2, workend3, workend4, workend5, workend6, workend7 from facebook_configurations WHERE domain_uuid = ?", domain_uuid).Rows()
	if err != nil {
		return &models.Response{Status: "0", Msg: "Record Not Found", ResponseCode: 401}, nil
	}
	defer row.Close()
	for row.Next() {
		f := models.FacebookConfiguration{}
		if err := row.Scan(&f.Id, &f.Domain_uuid, &f.AppId, &f.AppKey, &f.AppSecret, &f.FacebookIntegrationID, &f.Size, &f.TriggerName, &f.TriggerMessage, &f.TriggerWhen, &f.Day1, &f.Day2, &f.Day3, &f.Day4, &f.Day5, &f.Day6, &f.Day7, &f.Workstart1, &f.Workstart2, &f.Workstart3, &f.Workstart4, &f.Workstart5, &f.Workstart6, &f.Workstart7, &f.Workend1, &f.Workend2, &f.Workend3, &f.Workend4, &f.Workend5, &f.Workend6, &f.Workend7); err != nil {

			return nil, err
		}

		list = append(list, f)
	}
	return &models.Response{Status: "1", Msg: "Record Found", ResponseCode: 200, Fb: list}, nil
}

/**********************************************Update Tenant details*************************************/
func (r crudRepository) Update_Facebook_configuration(ctx context.Context, id int64, domain_uuid string, td models.FacebookConfigurations) (*models.Response, error) {
	w := models.FacebookConfiguration{}
	db1 := r.DBConn.Where("domain_uuid = ? AND id = ?", domain_uuid, id).Find(&w)
	if db1.RowsAffected == 0 {
		return &models.Response{Status: "Not Found", Msg: "Record Doesn't Exist", ResponseCode: 401}, nil
	}
	if td.FacebookIntegrationID == w.FacebookIntegrationID && td.AppId == w.AppId {
		if db := r.DBConn.Table("facebook_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"app_key": td.AppKey, "app_secret": td.AppSecret, "size": td.Size, "trigger_name": td.Trigger.Name, "trigger_message": td.Trigger.Message, "trigger_when": td.Trigger.When, "day1": td.WorkingDays[0].Day, "day2": td.WorkingDays[1].Day, "day3": td.WorkingDays[2].Day, "day4": td.WorkingDays[3].Day, "day5": td.WorkingDays[4].Day, "day6": td.WorkingDays[5].Day, "day7": td.WorkingDays[6].Day, "workstart1": td.WorkingDays[0].WorkingHourStartTime, "workstart2": td.WorkingDays[1].WorkingHourStartTime, "workstart3": td.WorkingDays[2].WorkingHourStartTime, "workstart4": td.WorkingDays[3].WorkingHourStartTime, "workstart5": td.WorkingDays[4].WorkingHourStartTime, "workstart6": td.WorkingDays[5].WorkingHourStartTime, "workstart7": td.WorkingDays[6].WorkingHourStartTime, "workend1": td.WorkingDays[0].WorkingHourEndTime, "workend2": td.WorkingDays[1].WorkingHourEndTime, "workend3": td.WorkingDays[2].WorkingHourEndTime, "workend4": td.WorkingDays[3].WorkingHourEndTime, "workend5": td.WorkingDays[4].WorkingHourEndTime, "workend6": td.WorkingDays[5].WorkingHourEndTime, "workend7": td.WorkingDays[6].WorkingHourEndTime, "message": td.Message}).Error; db != nil {
			return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
		}

		return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Facebook integration Updated successfully."}, nil
	} else if td.FacebookIntegrationID != w.FacebookIntegrationID && td.AppId != w.AppId {
		w := models.FacebookConfiguration{}
		if row := r.DBConn.Where("app_id = ?", td.AppId).Find(&w).Error; row != nil {
			if err := r.DBConn.Where("facebook_integration_id = ?", td.FacebookIntegrationID).Find(&w).Error; err != nil {
				if db := r.DBConn.Table("facebook_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"app_id": td.AppId, "app_key": td.AppKey, "app_secret": td.AppSecret, "facebook_integration_id": td.FacebookIntegrationID, "size": td.Size, "trigger_name": td.Trigger.Name, "trigger_message": td.Trigger.Message, "trigger_when": td.Trigger.When, "day1": td.WorkingDays[0].Day, "day2": td.WorkingDays[1].Day, "day3": td.WorkingDays[2].Day, "day4": td.WorkingDays[3].Day, "day5": td.WorkingDays[4].Day, "day6": td.WorkingDays[5].Day, "day7": td.WorkingDays[6].Day, "workstart1": td.WorkingDays[0].WorkingHourStartTime, "workstart2": td.WorkingDays[1].WorkingHourStartTime, "workstart3": td.WorkingDays[2].WorkingHourStartTime, "workstart4": td.WorkingDays[3].WorkingHourStartTime, "workstart5": td.WorkingDays[4].WorkingHourStartTime, "workstart6": td.WorkingDays[5].WorkingHourStartTime, "workstart7": td.WorkingDays[6].WorkingHourStartTime, "workend1": td.WorkingDays[0].WorkingHourEndTime, "workend2": td.WorkingDays[1].WorkingHourEndTime, "workend3": td.WorkingDays[2].WorkingHourEndTime, "workend4": td.WorkingDays[3].WorkingHourEndTime, "workend5": td.WorkingDays[4].WorkingHourEndTime, "workend6": td.WorkingDays[5].WorkingHourEndTime, "workend7": td.WorkingDays[6].WorkingHourEndTime, "message": td.Message}).Error; db != nil {
					return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
				}

				return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Facebook integration Updated successfully."}, nil
			}
			return &models.Response{ResponseCode: 409, Status: "0", Msg: "Facebook id already exist."}, nil
		}
		return &models.Response{ResponseCode: 409, Status: "0", Msg: "AppId already exist."}, nil
	} else if td.FacebookIntegrationID == w.FacebookIntegrationID && td.AppId != w.AppId {
		w := models.FacebookConfiguration{}
		if err := r.DBConn.Where("app_id = ?", td.AppId).Find(&w).Error; err != nil {
			if db := r.DBConn.Table("facebook_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"app_id": td.AppId, "app_key": td.AppKey, "app_secret": td.AppSecret, "size": td.Size, "trigger_name": td.Trigger.Name, "trigger_message": td.Trigger.Message, "trigger_when": td.Trigger.When, "day1": td.WorkingDays[0].Day, "day2": td.WorkingDays[1].Day, "day3": td.WorkingDays[2].Day, "day4": td.WorkingDays[3].Day, "day5": td.WorkingDays[4].Day, "day6": td.WorkingDays[5].Day, "day7": td.WorkingDays[6].Day, "workstart1": td.WorkingDays[0].WorkingHourStartTime, "workstart2": td.WorkingDays[1].WorkingHourStartTime, "workstart3": td.WorkingDays[2].WorkingHourStartTime, "workstart4": td.WorkingDays[3].WorkingHourStartTime, "workstart5": td.WorkingDays[4].WorkingHourStartTime, "workstart6": td.WorkingDays[5].WorkingHourStartTime, "workstart7": td.WorkingDays[6].WorkingHourStartTime, "workend1": td.WorkingDays[0].WorkingHourEndTime, "workend2": td.WorkingDays[1].WorkingHourEndTime, "workend3": td.WorkingDays[2].WorkingHourEndTime, "workend4": td.WorkingDays[3].WorkingHourEndTime, "workend5": td.WorkingDays[4].WorkingHourEndTime, "workend6": td.WorkingDays[5].WorkingHourEndTime, "workend7": td.WorkingDays[6].WorkingHourEndTime, "message": td.Message}).Error; db != nil {
				return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
			}

			return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Facebook integration Updated successfully."}, nil
		}
		return &models.Response{ResponseCode: 409, Status: "0", Msg: "AppId already exist."}, nil
	} else if td.FacebookIntegrationID != w.FacebookIntegrationID && td.AppId == w.AppId {
		w := models.FacebookConfiguration{}
		if err := r.DBConn.Where("facebook_integration_id = ?", td.FacebookIntegrationID).Find(&w).Error; err != nil {
			if db := r.DBConn.Table("facebook_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"app_key": td.AppKey, "app_secret": td.AppSecret, "facebook_integration_id": td.FacebookIntegrationID, "size": td.Size, "trigger_name": td.Trigger.Name, "trigger_message": td.Trigger.Message, "trigger_when": td.Trigger.When, "day1": td.WorkingDays[0].Day, "day2": td.WorkingDays[1].Day, "day3": td.WorkingDays[2].Day, "day4": td.WorkingDays[3].Day, "day5": td.WorkingDays[4].Day, "day6": td.WorkingDays[5].Day, "day7": td.WorkingDays[6].Day, "workstart1": td.WorkingDays[0].WorkingHourStartTime, "workstart2": td.WorkingDays[1].WorkingHourStartTime, "workstart3": td.WorkingDays[2].WorkingHourStartTime, "workstart4": td.WorkingDays[3].WorkingHourStartTime, "workstart5": td.WorkingDays[4].WorkingHourStartTime, "workstart6": td.WorkingDays[5].WorkingHourStartTime, "workstart7": td.WorkingDays[6].WorkingHourStartTime, "workend1": td.WorkingDays[0].WorkingHourEndTime, "workend2": td.WorkingDays[1].WorkingHourEndTime, "workend3": td.WorkingDays[2].WorkingHourEndTime, "workend4": td.WorkingDays[3].WorkingHourEndTime, "workend5": td.WorkingDays[4].WorkingHourEndTime, "workend6": td.WorkingDays[5].WorkingHourEndTime, "workend7": td.WorkingDays[6].WorkingHourEndTime, "message": td.Message}).Error; db != nil {
				return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
			}

			return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Facebook integration Updated successfully."}, nil
		}
		return &models.Response{ResponseCode: 409, Status: "0", Msg: "Facebook id already exist."}, nil
	}
	return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Facebook integration Updated successfully."}, nil
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
func (r crudRepository) Upload_Attachments(ctx context.Context, appId string, appUserId string, Type string, IntegrationID string, Size int64, file multipart.File, handler *multipart.FileHeader) (*models.Response, error) {
	td := models.WhatsappConfiguration{}
	db := r.DBConn.Table("whatsapp_configurations").Where("app_id = ? AND whatsapp_integration_id = ?", appId, IntegrationID).Find(&td)
	if db.Error != nil {
		fb := models.FacebookConfiguration{}
		df := r.DBConn.Table("facebook_configurations").Where("app_id = ? AND facebook_integration_id = ?", appId, IntegrationID).Find(&fb)
		if df.Error != nil {
			return &models.Response{Status: "0", Msg: "FacebookId not found.", ResponseCode: 404}, nil
		}
		if IntegrationID == fb.FacebookIntegrationID {
			n, _ := strconv.ParseInt(fb.Size, 10, 64)
			Size = n
			fmt.Println(Size, "size", handler.Size/1048576)
			if handler.Size/1048576 > n {
				return &models.Response{Status: "0", Msg: "File size is large please choose small file.", ResponseCode: 400}, nil
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
			req.SetBasicAuth(fb.AppKey, fb.AppSecret)
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
				return &models.Response{Status: "1", Msg: "File is sent successfully.", ResponseCode: 200}, nil
			}
		}
	}
	if IntegrationID == td.WhatsappIntegrationID {
		n, _ := strconv.ParseInt(td.Size, 10, 64)
		Size = n
		if handler.Size/1048576 > n {
			return &models.Response{Status: "0", Msg: "File size is large please choose small file.", ResponseCode: 400}, nil
		}
		fmt.Println(td.AppSecret, td.AppKey, "bcbjbjbj")
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
			fmt.Println(jsonData, string(data), u.MediaUrl, u.MediaType, "ghvghvqv")
			p := models.User{
				Role:      "appMaker",
				Type:      Type,
				MediaType: u.MediaType,
				MediaUrl:  u.MediaUrl,
			}
			r.PostMessage(ctx, appId, appUserId, p)
			return &models.Response{Status: "1", Msg: "File is sent successfully.", ResponseCode: 200}, nil
		}

	} else {
		return &models.Response{Status: "0", Msg: "Please choose message type.", ResponseCode: 200}, nil
	}
}

/***********************************************TypingActivity***********************************************/
func (r crudRepository) TypingActivity(ctx context.Context, appId string, appUserId string, p models.User) ([]byte, error) {
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

/*************************************************Disable AppUser*************************************************/
func (r crudRepository) Disable_AppUser(ctx context.Context, appUserId string) (*models.Response, error) {
	u := models.ReceiveUserDetails{}
	err := r.DBConn.Where("app_user_id = ?", appUserId).Find(&u)
	if err != nil {
		return &models.Response{Status: "0", Msg: "AppUserId Not Found.", ResponseCode: 404}, nil

	}
	db := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", appUserId).Update("is_enabled", false)
	if db.Error != nil {
		return &models.Response{Status: "0", Msg: "AppUserId not Updated.", ResponseCode: 404}, nil
	}
	return &models.Response{Status: "0", Msg: "AppUserId Disabled Successfully.", ResponseCode: 404}, nil
}
