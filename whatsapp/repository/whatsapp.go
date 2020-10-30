package repository

import (
	"bytes"
	"context"
	json "encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	config "whatsapp_api/config"
	models "whatsapp_api/model"
	crud "whatsapp_api/whatsapp"
	controller "whatsapp_api/whatsapp/controller"

	oauth1 "github.com/klaidas/go-oauth1"

	"github.com/google/uuid"
	myNewUUID "github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"golang.org/x/net/websocket"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

var UserOs string

var HTTPSERVERHOST string

const HTTPSECURE = "https://"
const PORT = "30707"

type crudRepository struct {
	DBConn *gorm.DB
	SList  *controller.ServerUserList
}

func NewcrudRepository(conn *gorm.DB, slist *controller.ServerUserList, conf *config.Config) crud.Repository {
	UserOs = conf.Server.OsUser
	HTTPSERVERHOST = conf.HttpConfig.HTTPSERVERHOST
	return &crudRepository{
		DBConn: conn,
		SList:  slist,
	}
}
func getServerOs() string {
	return UserOs
}
func getserverhost() string {
	return HTTPSERVERHOST
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
			re := models.ReceiveUserDetails{}
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

func (r *crudRepository) Get_allId(ctx context.Context, domain_uuid string) (*models.Response, error) {
	td := models.Tenant_details{}
	list := make([]models.ReceiveUserDetails, 0)
	if db := r.DBConn.Where("domain_uuid = ?", domain_uuid).Find(&td).Error; db != nil {

		return &models.Response{Status: "0", Msg: "Contact list is not available", ResponseCode: 404}, nil
	}
	if rows, err := r.DBConn.Raw("select app_id, app_user_id, surname, given_name,type,text,role,name,author_id,message_id,original_message_id,integration_id,source_type, signed_up_at, conversation_started, unread_count from receive_user_details where is_enabled = true").Rows(); err != nil {

		return &models.Response{Status: "Not Found", Msg: "Record Not Found", ResponseCode: 204}, nil
	} else {
		defer rows.Close()
		for rows.Next() {
			f := models.ReceiveUserDetails{}
			if err := rows.Scan(&f.AppId, &f.AppUserId, &f.Surname, &f.GivenName, &f.Type, &f.Text, &f.Role, &f.Name, &f.AuthorID, &f.Message_id, &f.OriginalMessageID, &f.IntegrationID, &f.Source_Type, &f.SignedUpAt, &f.ConversationStarted, &f.UnreadCount); err != nil {

				return nil, err
			}

			list = append(list, f)
		}

		return &models.Response{Status: "OK", Msg: "Record Found", ResponseCode: 200, AppUserList: list}, nil
	}

}

/**********************************************Get customer by appUserId*********************************************/
func (r *crudRepository) Get_Customer_by_agent_uuid(ctx context.Context, customer_id string) (*models.Response, error) {
	customer := models.ReceiveUserDetails{
		AppUserId: customer_id,
	}
	//list := make([]models.Customer_Agents, 0)
	if db := r.DBConn.Where("app_user_id = ?", customer_id).Find(&customer).Error; db != nil {

		return &models.Response{Status: "0", Msg: "Contact list is not available", ResponseCode: 404}, nil
	}
	// if rows, err := r.DBConn.Raw("select domain_uuid, agent_uuid,app_user_id, surname, given_name,type,text,role,name,author_id,conversation_id,received,message_id,integration_id,source_type, signed_up_at, unread_count from customer_agents where app_user_id = ?", customer_id).Rows(); err != nil {

	// 	return &models.Response{Status: "Not Found", Msg: "Record Not Found", ResponseCode: 204}, nil
	// } else {
	// 	defer rows.Close()
	// 	for rows.Next() {
	// 		f := models.Customer_Agents{}
	// 		if err := rows.Scan(&f.Domain_uuid, &f.Agent_uuid, &f.AppUserId, &f.Surname, &f.GivenName, &f.Type, &f.Text, &f.Role, &f.Name, &f.AuthorID, &f.Conversation_id, &f.Received, &f.Message_id, &f.IntegrationID, &f.Source_Type, &f.SignedUpAt, &f.UnreadCount); err != nil {

	// 			return nil, err
	// 		}

	// 		list = append(list, f)
	// 	}

	return &models.Response{Status: "OK", Msg: "Record Found", ResponseCode: 200, Customer: &customer}, nil
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
	T := models.Account_details{}
	f := models.Received{}
	w := models.WhatsappConfiguration{}
	tw := models.TwitterConfiguration{}
	fb := models.FacebookConfiguration{}
	jsondata := json.Unmarshal(body, &f)
	fmt.Println(jsondata)
	s := int64(f.Messages[0].Received)
	myDate := time.Unix(s, 0)
	_, _, date := myDate.Date()
	var dom_uuid string
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
		UnreadCount:              0,
		Day:                      myDate.Weekday().String(),
		Date:                     date,
		AfterOfficeTime:          false,
		Domain_uuid:              dom_uuid,
	}

	fmt.Println(f.Messages[0].Source.Type, f.Messages[0].Source.IntegrationID, "values.......")
	//queue := models.Queue{}
	cou := []models.Count_Agent_customer{}
	agent := models.AgentQueue{}
	db := r.DBConn.Table("customer_agents ca").Select("count(ca.agent_uuid),aq.agent_uuid, aq.tenant_domain_uuid").Joins("right join (select agent_uuid,tenant_domain_uuid from agent_queues inner join v_call_center_agents on agent_queues.agent_uuid=v_call_center_agents.call_center_agent_uuid where agent_status='Available' and queue_uuid=(select queue_uuid from queues where integration_id='" + f.Messages[0].Source.IntegrationID + "')) aq on aq.agent_uuid::text=ca.agent_uuid group by aq.agent_uuid,aq.tenant_domain_uuid").Find(&cou)
	if db.Error != nil {
		fmt.Println(db.Error)
	}

	fmt.Println(cou, len(cou))
	var min int64 = 0
	var max int64 = 0
	lengtharray := len(cou)
	if lengtharray == 0 {
		if f.Messages[0].Source.Type == "messenger" {
			db := r.DBConn.Where("facebook_integration_id = ?", f.Messages[0].Source.IntegrationID).Find(&fb)
			if db.Error != nil {
				fmt.Println("error")
			}
			fmt.Println("enterrrrrrr.....")
			tenant := r.DBConn.Table("account_details").Where("domain_uuid = ?", fb.Domain_uuid).Find(&T)
			if tenant.Error != nil {
				fmt.Println("error")
			}

			uuid1, _ := myNewUUID.NewUUID()
			uuid := uuid1.String()
			tic := models.SocialMediaTickets{
				Ticket_uuid:     uuid,
				Domain_uuid:     T.Domain_uuid,
				Ticket_name:     f.Messages[0].Source.Type,
				CustomerId:      f.AppUser.ID,
				CustomerName:    f.Messages[0].Name,
				Message:         f.Messages[0].Text,
				MessageType:     f.Messages[0].Type,
				IntegrationID:   f.Messages[0].Source.IntegrationID,
				Source_type:     f.Messages[0].Source.Type,
				Conversation_id: f.Conversation.ID,
				Timestamp:       time.Now(),
			}
			if db := r.DBConn.Table("social_media_tickets").Where("customer_id = ?", f.AppUser.ID).Find(&tic).Error; db != nil {
				fmt.Println("enterrrrrrr.....1111111")
				row := r.DBConn.Create(&tic)
				if row.Error != nil {
					fmt.Println("Ticket not created.")
				}
				p := models.User{
					Author: models.Author{
						Type:        "business",
						DisplayName: T.Tenant_name,
						AvatarURL:   "https://www.gravatar.com/image.jpg",
					},
					Content: models.Content{
						Type: "text",
						Text: "Sorry " + f.Messages[0].Name + "! for the inconvenience our agents are not available at this time we will reach out to you when our agents are available. Thanks for contacting us.",
					},
				}
				r.PostMessage(ctx, fb.AppId, f.Conversation.ID, p)
			}
		} else if f.Messages[0].Source.Type == "whatsapp" {
			db := r.DBConn.Where("whatsapp_integration_id = ?", f.Messages[0].Source.IntegrationID).Find(&w)
			if db.Error != nil {
				fmt.Println("error")
			}
			tenant := r.DBConn.Table("account_details").Where("domain_uuid = ?", w.Domain_uuid).Find(&T)
			if tenant.Error != nil {
				fmt.Println("error")
			}
			uuid1, _ := myNewUUID.NewUUID()
			uuid := uuid1.String()
			tic := models.SocialMediaTickets{
				Ticket_uuid:     uuid,
				Domain_uuid:     T.Domain_uuid,
				Ticket_name:     f.Messages[0].Source.Type,
				CustomerId:      f.AppUser.ID,
				CustomerName:    f.Messages[0].Name,
				Message:         f.Messages[0].Text,
				MessageType:     f.Messages[0].Type,
				IntegrationID:   f.Messages[0].Source.IntegrationID,
				Source_type:     f.Messages[0].Source.Type,
				Conversation_id: f.Conversation.ID,
				Timestamp:       time.Now(),
			}
			if db := r.DBConn.Table("social_media_tickets").Where("customer_id = ?", f.AppUser.ID).Find(&tic).Error; db != nil {
				row := r.DBConn.Create(&tic)
				if row.Error != nil {
					fmt.Println("Ticket not created.")
				}
				p := models.User{
					Author: models.Author{
						Type:        "business",
						DisplayName: T.Tenant_name,
						AvatarURL:   "https://www.gravatar.com/image.jpg",
					},
					Content: models.Content{
						Type: "text",
						Text: "Sorry " + f.Messages[0].Name + "! for the inconvenience our agents are not available at this time we will reach out to you when our agents are available. Thanks for contacting us.",
					},
				}
				r.PostMessage(ctx, fb.AppId, f.Conversation.ID, p)

			}

		} else if f.Messages[0].Source.Type == "twitter" {
			db := r.DBConn.Where("twitter_integration_id = ?", f.Messages[0].Source.IntegrationID).Find(&tw)
			if db.Error != nil {
				fmt.Println("error")
			}
			tenant := r.DBConn.Table("account_details").Where("domain_uuid = ?", tw.Domain_uuid).Find(&T)
			if tenant.Error != nil {
				fmt.Println("error")
			}
			uuid1, _ := myNewUUID.NewUUID()
			uuid := uuid1.String()
			tic := models.SocialMediaTickets{
				Ticket_uuid:     uuid,
				Domain_uuid:     T.Domain_uuid,
				Ticket_name:     f.Messages[0].Source.Type,
				CustomerId:      f.AppUser.ID,
				CustomerName:    f.Messages[0].Name,
				Message:         f.Messages[0].Text,
				MessageType:     f.Messages[0].Type,
				IntegrationID:   f.Messages[0].Source.IntegrationID,
				Source_type:     f.Messages[0].Source.Type,
				Conversation_id: f.Conversation.ID,
				Timestamp:       time.Now(),
			}
			if db := r.DBConn.Table("social_media_tickets").Where("customer_id = ?", f.AppUser.ID).Find(&tic).Error; db != nil {
				row := r.DBConn.Create(&tic)
				if row.Error != nil {
					fmt.Println("Ticket not created.")
				}
				p := models.User{
					Author: models.Author{
						Type:        "business",
						DisplayName: T.Tenant_name,
						AvatarURL:   "https://www.gravatar.com/image.jpg",
					},
					Content: models.Content{
						Type: "text",
						Text: "Sorry " + f.Messages[0].Name + "! for the inconvenience our agents are not available at this time we will reach out to you when our agents are available. Thanks for contacting us.",
					},
				}
				r.PostMessage(ctx, tw.AppId, f.Conversation.ID, p)

			}
		}
	}
	//max_uuid := cou[0].Agent_uuid
	min_uuid := cou[0].Agent_uuid
	for k, v := range cou {
		fmt.Println(k, v)

		if v.Count >= max {
			max = v.Count
			//max_uuid = v.Agent_uuid
		} else {
			min = v.Count
			min_uuid = v.Agent_uuid
		}
		//fmt.Println(" min= ", min, " max= ", max, " min_uuid= ", min_uuid)
	}
	if min == max && min == 0 {

		agent.Agent_uuid = min_uuid
		agent.Tenant_domain_uuid = cou[0].Tenant_domain_uuid
	} else {

		agent.Agent_uuid = min_uuid
		agent.Tenant_domain_uuid = cou[0].Tenant_domain_uuid
	}
	//fmt.Println("2 min= ", min, " max= ", max, " min_uuid= ", min_uuid)
	customer := models.Customer_Agents{
		Domain_uuid:     agent.Tenant_domain_uuid,
		AppUserId:       f.AppUser.ID,
		Agent_uuid:      agent.Agent_uuid,
		Surname:         f.AppUser.Surname,
		GivenName:       f.AppUser.GivenName,
		SignedUpAt:      f.AppUser.SignedUpAt,
		Conversation_id: f.Conversation.ID,
		Type:            f.Messages[0].Type,
		Text:            f.Messages[0].Text,
		Role:            f.Messages[0].Role,
		Received:        f.Messages[0].Received,
		Name:            f.Messages[0].Name,
		AuthorID:        f.Messages[0].AuthorID,
		Message_id:      f.Messages[0].ID,
		Source_Type:     f.Messages[0].Source.Type,
		IntegrationID:   f.Messages[0].Source.IntegrationID,
		UnreadCount:     0,
	}
	//fmt.Println(agent, "cus ", customer)

	if cust := r.DBConn.Where("app_user_id = ?", f.AppUser.ID).Find(&customer).Error; cust != nil {
		db := r.DBConn.Create(&customer)
		if db.Error != nil {
			fmt.Println(db.Error)
		}

		for _, oldu := range r.SList.Users {
			if oldu.UName == customer.Agent_uuid {
				msg := map[string]interface{}{"message_id": "5", "customer_id": customer.AppUserId, "surname": f.AppUser.Surname, "given_name": f.AppUser.GivenName, "signed_up_at": f.AppUser.SignedUpAt, "conversation_id": f.Conversation.ID, "type": f.Messages[0].Type, "text": f.Messages[0].Text, "role": f.Messages[0].Role, "received": f.Messages[0].Received, "name": f.Messages[0].Name, "author_id": f.Messages[0].AuthorID, "messageId": f.Messages[0].ID, "source_type": f.Messages[0].Source.Type, "integration_id": f.Messages[0].Source.IntegrationID, "unread_count": u.UnreadCount, "user_id": customer.Agent_uuid}
				if err := websocket.JSON.Send(oldu.Ws, msg); err != nil {
					log.Println("Can't send", err)
				}
			}
		}
	}
	cou_cus := models.Count_customer{}
	cou_cust := r.DBConn.Table("customer_agents").Select("count(customer_id)").Where("agent_uuid = ?", agent.Agent_uuid).Find(&cou_cus)
	if cou_cust.Error != nil {
		fmt.Println(cou_cust.Error)
	}
	g := 5 * (cou_cus.Count / 1)
	errs := r.DBConn.Where("app_user_id = ?", f.AppUser.ID).Find(&u)
	fmt.Println(errs.Error, g)
	if f.Messages[0].Role == "appUser" {
		count := r.DBConn.Table("receive_user_details").Where("conversation_id = ? AND app_user_id = ?", f.Conversation.ID, f.AppUser.ID).Update("unread_count", u.UnreadCount+1)
		if count.Error != nil {
			fmt.Println(count.Error)
		}
	} else {
		fmt.Println("error")
	}
	if u.Is_enabled == false {
		update := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Update("is_enabled", true)
		fmt.Println(update, update.RowsAffected)
	}

	if f.Messages[0].Source.Type == "messenger" {
		db := r.DBConn.Where("facebook_integration_id = ?", f.Messages[0].Source.IntegrationID).Find(&fb)
		if db.Error != nil {
			fmt.Println("error")
		}
		tenant := r.DBConn.Table("account_details").Where("domain_uuid = ?", fb.Domain_uuid).Find(&T)
		if tenant.Error != nil {
			fmt.Println("error")
		}
		if myDate.Weekday().String() == fb.Day1 {
			workstart1 := fb.Workstart1
			components := strings.Split(workstart1, ":")
			StartHour, _ := components[0], components[1]
			workend1 := fb.Workend1
			components1 := strings.Split(workend1, ":")
			EndHour, _ := components1[0], components1[1]
			startHour, _ := strconv.Atoi(StartHour)
			endHour, _ := strconv.Atoi(EndHour)
			if myDate.Hour() < startHour || myDate.Hour() > endHour {
				if aot := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; aot != nil {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: fb.Message,
						},
					}
					r.PostMessage(ctx, fb.AppId, f.Conversation.ID, p)
					db := r.DBConn.Create(&u).Where("app_user_id = ?", f.AppUser.ID).Update("after_office_time", true)
					if db.Error != nil {
						fmt.Println(db.Error)
					}
				}

				if myDate.Weekday().String() == u.Day && date == u.Date {
					if u.AfterOfficeTime == true {
						fmt.Println("message already sent.")
					} else if u.AfterOfficeTime == false {
						p := models.User{
							Author: models.Author{
								Type:        "business",
								DisplayName: T.Tenant_name,
								AvatarURL:   "https://www.gravatar.com/image.jpg",
							},
							Content: models.Content{
								Type: "text",
								Text: fb.Message,
							},
						}
						r.PostMessage(ctx, fb.AppId, f.Conversation.ID, p)
						err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
						if err.RowsAffected == 0 {
							fmt.Println("rows not updated.")
						}
					}

				} else {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: fb.Message,
						},
					}
					r.PostMessage(ctx, fb.AppId, f.Conversation.ID, p)
					err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
					if err.RowsAffected == 0 {
						fmt.Println("rows not updated.")
					}
				}

			} else {

				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: fb.TriggerMessage,
						},
					}
					r.PostMessage(ctx, fb.AppId, f.Conversation.ID, p)
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
			startHour, _ := strconv.Atoi(StartHour)
			endHour, _ := strconv.Atoi(EndHour)
			if myDate.Hour() < startHour || myDate.Hour() > endHour {
				if aot := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; aot != nil {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: fb.Message,
						},
					}
					r.PostMessage(ctx, fb.AppId, f.Conversation.ID, p)
					db := r.DBConn.Create(&u).Where("app_user_id = ?", f.AppUser.ID).Update("after_office_time", true)
					if db.Error != nil {
						fmt.Println(db.Error)
					}
				}

				if myDate.Weekday().String() == u.Day && date == u.Date {
					if u.AfterOfficeTime == true {
						fmt.Println("message already sent.")
					} else if u.AfterOfficeTime == false {
						p := models.User{
							Author: models.Author{
								Type:        "business",
								DisplayName: T.Tenant_name,
								AvatarURL:   "https://www.gravatar.com/image.jpg",
							},
							Content: models.Content{
								Type: "text",
								Text: fb.Message,
							},
						}
						r.PostMessage(ctx, fb.AppId, f.Conversation.ID, p)
						err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
						if err.RowsAffected == 0 {
							fmt.Println("rows not updated.")
						}
					}

				} else {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: fb.Message,
						},
					}
					r.PostMessage(ctx, fb.AppId, f.Conversation.ID, p)
					err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
					if err.RowsAffected == 0 {
						fmt.Println("rows not updated.")
					}

				}

			} else {

				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: fb.TriggerMessage,
						},
					}
					r.PostMessage(ctx, fb.AppId, f.Conversation.ID, p)
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
			startHour, _ := strconv.Atoi(StartHour)
			endHour, _ := strconv.Atoi(EndHour)
			if myDate.Hour() < startHour || myDate.Hour() > endHour {
				if aot := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; aot != nil {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: fb.Message,
						},
					}
					r.PostMessage(ctx, fb.AppId, f.Conversation.ID, p)
					db := r.DBConn.Create(&u).Where("app_user_id = ?", f.AppUser.ID).Update("after_office_time", true)
					if db.Error != nil {
						fmt.Println(db.Error)
					}
				}

				if myDate.Weekday().String() == u.Day && date == u.Date {
					if u.AfterOfficeTime == true {
						fmt.Println("message already sent.")
					} else if u.AfterOfficeTime == false {
						p := models.User{
							Author: models.Author{
								Type:        "business",
								DisplayName: T.Tenant_name,
								AvatarURL:   "https://www.gravatar.com/image.jpg",
							},
							Content: models.Content{
								Type: "text",
								Text: fb.Message,
							},
						}
						r.PostMessage(ctx, fb.AppId, f.Conversation.ID, p)
						err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
						if err.RowsAffected == 0 {
							fmt.Println("rows not updated.")
						}
					}

				} else {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: fb.Message,
						},
					}
					r.PostMessage(ctx, fb.AppId, f.Conversation.ID, p)
					err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
					if err.RowsAffected == 0 {
						fmt.Println("rows not updated.")
					}

				}

			} else {

				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: fb.TriggerMessage,
						},
					}
					r.PostMessage(ctx, fb.AppId, f.Conversation.ID, p)
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
			startHour, _ := strconv.Atoi(StartHour)
			endHour, _ := strconv.Atoi(EndHour)
			if myDate.Hour() < startHour || myDate.Hour() > endHour {
				if aot := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; aot != nil {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: fb.Message,
						},
					}
					r.PostMessage(ctx, fb.AppId, f.Conversation.ID, p)
					db := r.DBConn.Create(&u).Where("app_user_id = ?", f.AppUser.ID).Update("after_office_time", true)
					if db.Error != nil {
						fmt.Println(db.Error)
					}
				}

				if myDate.Weekday().String() == u.Day && date == u.Date {
					if u.AfterOfficeTime == true {
						fmt.Println("message already sent.")
					} else if u.AfterOfficeTime == false {
						p := models.User{
							Author: models.Author{
								Type:        "business",
								DisplayName: T.Tenant_name,
								AvatarURL:   "https://www.gravatar.com/image.jpg",
							},
							Content: models.Content{
								Type: "text",
								Text: fb.Message,
							},
						}
						r.PostMessage(ctx, fb.AppId, f.Conversation.ID, p)
						err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
						if err.RowsAffected == 0 {
							fmt.Println("rows not updated.")
						}
					}

				} else {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: fb.Message,
						},
					}
					r.PostMessage(ctx, fb.AppId, f.Conversation.ID, p)
					err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
					if err.RowsAffected == 0 {
						fmt.Println("rows not updated.")
					}

				}

			} else {

				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: fb.TriggerMessage,
						},
					}
					r.PostMessage(ctx, fb.AppId, f.Conversation.ID, p)
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
			startHour, _ := strconv.Atoi(StartHour)
			endHour, _ := strconv.Atoi(EndHour)
			if myDate.Hour() < startHour || myDate.Hour() > endHour {
				if aot := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; aot != nil {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: fb.Message,
						},
					}
					r.PostMessage(ctx, fb.AppId, f.Conversation.ID, p)
					db := r.DBConn.Create(&u).Where("app_user_id = ?", f.AppUser.ID).Update("after_office_time", true)
					if db.Error != nil {
						fmt.Println(db.Error)
					}
				}

				if myDate.Weekday().String() == u.Day && date == u.Date {
					if u.AfterOfficeTime == true {
						fmt.Println("message already sent.")
					} else if u.AfterOfficeTime == false {
						p := models.User{
							Author: models.Author{
								Type:        "business",
								DisplayName: T.Tenant_name,
								AvatarURL:   "https://www.gravatar.com/image.jpg",
							},
							Content: models.Content{
								Type: "text",
								Text: fb.Message,
							},
						}
						r.PostMessage(ctx, fb.AppId, f.Conversation.ID, p)
						err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
						if err.RowsAffected == 0 {
							fmt.Println("rows not updated.")
						}
					}

				} else {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: fb.Message,
						},
					}
					r.PostMessage(ctx, fb.AppId, f.Conversation.ID, p)
					err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
					if err.RowsAffected == 0 {
						fmt.Println("rows not updated.")
					}

				}
			} else {

				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: fb.TriggerMessage,
						},
					}
					r.PostMessage(ctx, fb.AppId, f.Conversation.ID, p)
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
			startHour, _ := strconv.Atoi(StartHour)
			endHour, _ := strconv.Atoi(EndHour)
			if myDate.Hour() < startHour || myDate.Hour() > endHour {
				if aot := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; aot != nil {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: fb.Message,
						},
					}
					r.PostMessage(ctx, fb.AppId, f.Conversation.ID, p)
					db := r.DBConn.Create(&u).Where("app_user_id = ?", f.AppUser.ID).Update("after_office_time", true)
					if db.Error != nil {
						fmt.Println(db.Error)
					}
				}

				if myDate.Weekday().String() == u.Day && date == u.Date {
					if u.AfterOfficeTime == true {
						fmt.Println("message already sent.")
					} else if u.AfterOfficeTime == false {
						p := models.User{
							Author: models.Author{
								Type:        "business",
								DisplayName: T.Tenant_name,
								AvatarURL:   "https://www.gravatar.com/image.jpg",
							},
							Content: models.Content{
								Type: "text",
								Text: fb.Message,
							},
						}
						r.PostMessage(ctx, fb.AppId, f.Conversation.ID, p)
						err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
						if err.RowsAffected == 0 {
							fmt.Println("rows not updated.")
						}
					}

				} else {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: fb.Message,
						},
					}
					r.PostMessage(ctx, fb.AppId, f.Conversation.ID, p)
					err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
					if err.RowsAffected == 0 {
						fmt.Println("rows not updated.")
					}

				}

			} else {

				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: fb.TriggerMessage,
						},
					}
					r.PostMessage(ctx, fb.AppId, f.Conversation.ID, p)
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
			startHour, _ := strconv.Atoi(StartHour)
			endHour, _ := strconv.Atoi(EndHour)

			if myDate.Hour() < startHour || myDate.Hour() > endHour {
				if aot := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; aot != nil {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: fb.Message,
						},
					}
					r.PostMessage(ctx, fb.AppId, f.Conversation.ID, p)
					db := r.DBConn.Create(&u).Where("app_user_id = ?", f.AppUser.ID).Update("after_office_time", true)
					if db.Error != nil {
						fmt.Println(db.Error)
					}
				}

				if myDate.Weekday().String() == u.Day && date == u.Date {
					if u.AfterOfficeTime == true {
						fmt.Println("message already sent.")
					} else if u.AfterOfficeTime == false {
						p := models.User{
							Author: models.Author{
								Type:        "business",
								DisplayName: T.Tenant_name,
								AvatarURL:   "https://www.gravatar.com/image.jpg",
							},
							Content: models.Content{
								Type: "text",
								Text: fb.Message,
							},
						}
						r.PostMessage(ctx, fb.AppId, f.Conversation.ID, p)
						err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
						if err.RowsAffected == 0 {
							fmt.Println("rows not updated.")
						}
					}

				} else {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: fb.Message,
						},
					}
					r.PostMessage(ctx, fb.AppId, f.Conversation.ID, p)
					err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
					if err.RowsAffected == 0 {
						fmt.Println("rows not updated.")
					}

				}

			} else {

				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: fb.TriggerMessage,
						},
					}
					r.PostMessage(ctx, fb.AppId, f.Conversation.ID, p)
					return &models.Response{Received: &f}, nil
				}
				fmt.Println("appUserId already exist.")
				return &models.Response{Msg: "Userid already exist."}, nil

			}
		}
		// else {
		// 	fmt.Println("User Not Registered.")
		// }
	} else if f.Messages[0].Source.Type == "whatsapp" {
		db := r.DBConn.Where("whatsapp_integration_id = ?", f.Messages[0].Source.IntegrationID).Find(&w)
		if db.Error != nil {
			fmt.Println("error")
		}
		tenant := r.DBConn.Table("account_details").Where("domain_uuid = ?", w.Domain_uuid).Find(&T)
		if tenant.Error != nil {
			fmt.Println("error")
		}
		if myDate.Weekday().String() == w.Day1 {
			workstart1 := w.Workstart1
			components := strings.Split(workstart1, ":")
			StartHour, _ := components[0], components[1]
			workend1 := w.Workend1
			components1 := strings.Split(workend1, ":")
			EndHour, _ := components1[0], components1[1]
			startHour, _ := strconv.Atoi(StartHour)
			endHour, _ := strconv.Atoi(EndHour)
			if myDate.Hour() < startHour || myDate.Hour() > endHour {
				if aot := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; aot != nil {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: w.Message,
						},
					}
					r.PostMessage(ctx, w.AppId, f.Conversation.ID, p)
					db := r.DBConn.Create(&u).Where("app_user_id = ?", f.AppUser.ID).Update("after_office_time", true)
					if db.Error != nil {
						fmt.Println(db.Error)
					}
				}

				if myDate.Weekday().String() == u.Day && date == u.Date {
					if u.AfterOfficeTime == true {
						fmt.Println("message already sent.")
					} else if u.AfterOfficeTime == false {
						p := models.User{
							Author: models.Author{
								Type:        "business",
								DisplayName: T.Tenant_name,
								AvatarURL:   "https://www.gravatar.com/image.jpg",
							},
							Content: models.Content{
								Type: "text",
								Text: w.Message,
							},
						}
						r.PostMessage(ctx, w.AppId, f.Conversation.ID, p)
						err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
						if err.RowsAffected == 0 {
							fmt.Println("rows not updated.")
						}
					}

				} else {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: w.Message,
						},
					}
					r.PostMessage(ctx, w.AppId, f.Conversation.ID, p)
					err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
					if err.RowsAffected == 0 {
						fmt.Println("rows not updated.")
					}

				}

			} else {
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}

					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: w.TriggerMessage,
						},
					}
					r.PostMessage(ctx, w.AppId, f.Conversation.ID, p)
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
			startHour, _ := strconv.Atoi(StartHour)
			endHour, _ := strconv.Atoi(EndHour)
			if myDate.Hour() < startHour || myDate.Hour() > endHour {
				if aot := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; aot != nil {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: w.Message,
						},
					}
					r.PostMessage(ctx, w.AppId, f.Conversation.ID, p)
					db := r.DBConn.Create(&u).Where("app_user_id = ?", f.AppUser.ID).Update("after_office_time", true)
					if db.Error != nil {
						fmt.Println(db.Error)
					}
				}

				if myDate.Weekday().String() == u.Day && date == u.Date {
					if u.AfterOfficeTime == true {
						fmt.Println("message already sent.")
					} else if u.AfterOfficeTime == false {
						p := models.User{
							Author: models.Author{
								Type:        "business",
								DisplayName: T.Tenant_name,
								AvatarURL:   "https://www.gravatar.com/image.jpg",
							},
							Content: models.Content{
								Type: "text",
								Text: w.Message,
							},
						}
						r.PostMessage(ctx, w.AppId, f.Conversation.ID, p)
						err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
						if err.RowsAffected == 0 {
							fmt.Println("rows not updated.")
						}
					}

				} else {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: w.Message,
						},
					}
					r.PostMessage(ctx, w.AppId, f.Conversation.ID, p)
					err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
					if err.RowsAffected == 0 {
						fmt.Println("rows not updated.")
					}

				}

			} else {
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}

					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: w.TriggerMessage,
						},
					}
					r.PostMessage(ctx, w.AppId, f.Conversation.ID, p)
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
			startHour, _ := strconv.Atoi(StartHour)
			endHour, _ := strconv.Atoi(EndHour)
			if myDate.Hour() < startHour || myDate.Hour() > endHour {
				if aot := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; aot != nil {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: w.Message,
						},
					}
					r.PostMessage(ctx, w.AppId, f.Conversation.ID, p)
					db := r.DBConn.Create(&u).Where("app_user_id = ?", f.AppUser.ID).Update("after_office_time", true)
					if db.Error != nil {
						fmt.Println(db.Error)
					}
				}

				if myDate.Weekday().String() == u.Day && date == u.Date {
					if u.AfterOfficeTime == true {
						fmt.Println("message already sent.")
					} else if u.AfterOfficeTime == false {
						p := models.User{
							Author: models.Author{
								Type:        "business",
								DisplayName: T.Tenant_name,
								AvatarURL:   "https://www.gravatar.com/image.jpg",
							},
							Content: models.Content{
								Type: "text",
								Text: w.Message,
							},
						}
						r.PostMessage(ctx, w.AppId, f.Conversation.ID, p)
						err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
						if err.RowsAffected == 0 {
							fmt.Println("rows not updated.")
						}
					}

				} else {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: w.Message,
						},
					}
					r.PostMessage(ctx, w.AppId, f.Conversation.ID, p)
					err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
					if err.RowsAffected == 0 {
						fmt.Println("rows not updated.")
					}

				}

			} else {
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}

					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: w.TriggerMessage,
						},
					}
					r.PostMessage(ctx, w.AppId, f.Conversation.ID, p)
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
			startHour, _ := strconv.Atoi(StartHour)
			endHour, _ := strconv.Atoi(EndHour)
			if myDate.Hour() < startHour || myDate.Hour() > endHour {
				if aot := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; aot != nil {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: w.Message,
						},
					}
					r.PostMessage(ctx, w.AppId, f.Conversation.ID, p)
					db := r.DBConn.Create(&u).Where("app_user_id = ?", f.AppUser.ID).Update("after_office_time", true)
					if db.Error != nil {
						fmt.Println(db.Error)
					}
				}

				if myDate.Weekday().String() == u.Day && date == u.Date {
					if u.AfterOfficeTime == true {
						fmt.Println("message already sent.")
					} else if u.AfterOfficeTime == false {
						p := models.User{
							Author: models.Author{
								Type:        "business",
								DisplayName: T.Tenant_name,
								AvatarURL:   "https://www.gravatar.com/image.jpg",
							},
							Content: models.Content{
								Type: "text",
								Text: w.Message,
							},
						}
						r.PostMessage(ctx, w.AppId, f.Conversation.ID, p)
						err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
						if err.RowsAffected == 0 {
							fmt.Println("rows not updated.")
						}
					}

				} else {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: w.Message,
						},
					}
					r.PostMessage(ctx, w.AppId, f.Conversation.ID, p)
					err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
					if err.RowsAffected == 0 {
						fmt.Println("rows not updated.")
					}

				}

			} else {
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}

					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: w.TriggerMessage,
						},
					}
					r.PostMessage(ctx, w.AppId, f.Conversation.ID, p)
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
			startHour, _ := strconv.Atoi(StartHour)
			endHour, _ := strconv.Atoi(EndHour)
			if myDate.Hour() < startHour || myDate.Hour() > endHour {
				if aot := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; aot != nil {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: w.Message,
						},
					}
					r.PostMessage(ctx, w.AppId, f.Conversation.ID, p)
					db := r.DBConn.Create(&u).Where("app_user_id = ?", f.AppUser.ID).Update("after_office_time", true)
					if db.Error != nil {
						fmt.Println(db.Error)
					}
				}

				if myDate.Weekday().String() == u.Day && date == u.Date {
					if u.AfterOfficeTime == true {
						fmt.Println("message already sent.")
					} else if u.AfterOfficeTime == false {
						p := models.User{
							Author: models.Author{
								Type:        "business",
								DisplayName: T.Tenant_name,
								AvatarURL:   "https://www.gravatar.com/image.jpg",
							},
							Content: models.Content{
								Type: "text",
								Text: w.Message,
							},
						}
						r.PostMessage(ctx, w.AppId, f.Conversation.ID, p)
						err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
						if err.RowsAffected == 0 {
							fmt.Println("rows not updated.")
						}
					}

				} else {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: w.Message,
						},
					}
					r.PostMessage(ctx, w.AppId, f.Conversation.ID, p)
					err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
					if err.RowsAffected == 0 {
						fmt.Println("rows not updated.")
					}

				}

			} else {
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}

					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: w.TriggerMessage,
						},
					}
					r.PostMessage(ctx, w.AppId, f.Conversation.ID, p)
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
			startHour, _ := strconv.Atoi(StartHour)
			endHour, _ := strconv.Atoi(EndHour)
			if myDate.Hour() < startHour || myDate.Hour() > endHour {
				if aot := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; aot != nil {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: w.Message,
						},
					}
					r.PostMessage(ctx, w.AppId, f.Conversation.ID, p)
					db := r.DBConn.Create(&u).Where("app_user_id = ?", f.AppUser.ID).Update("after_office_time", true)
					if db.Error != nil {
						fmt.Println(db.Error)
					}
				}

				if myDate.Weekday().String() == u.Day && date == u.Date {
					if u.AfterOfficeTime == true {
						fmt.Println("message already sent.")
					} else if u.AfterOfficeTime == false {
						p := models.User{
							Author: models.Author{
								Type:        "business",
								DisplayName: T.Tenant_name,
								AvatarURL:   "https://www.gravatar.com/image.jpg",
							},
							Content: models.Content{
								Type: "text",
								Text: w.Message,
							},
						}
						r.PostMessage(ctx, w.AppId, f.Conversation.ID, p)
						err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
						if err.RowsAffected == 0 {
							fmt.Println("rows not updated.")
						}
					}

				} else {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: w.Message,
						},
					}
					r.PostMessage(ctx, w.AppId, f.Conversation.ID, p)
					err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
					if err.RowsAffected == 0 {
						fmt.Println("rows not updated.")
					}

				}

			} else {
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}

					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: w.TriggerMessage,
						},
					}
					r.PostMessage(ctx, w.AppId, f.Conversation.ID, p)
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
			startHour, _ := strconv.Atoi(StartHour)
			endHour, _ := strconv.Atoi(EndHour)
			if myDate.Hour() < startHour || myDate.Hour() > endHour {
				if aot := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; aot != nil {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: w.Message,
						},
					}
					r.PostMessage(ctx, w.AppId, f.Conversation.ID, p)
					db := r.DBConn.Create(&u).Where("app_user_id = ?", f.AppUser.ID).Update("after_office_time", true)
					if db.Error != nil {
						fmt.Println(db.Error)
					}
				}

				if myDate.Weekday().String() == u.Day && date == u.Date {
					if u.AfterOfficeTime == true {
						fmt.Println("message already sent.")
					} else if u.AfterOfficeTime == false {
						p := models.User{
							Author: models.Author{
								Type:        "business",
								DisplayName: T.Tenant_name,
								AvatarURL:   "https://www.gravatar.com/image.jpg",
							},
							Content: models.Content{
								Type: "text",
								Text: w.Message,
							},
						}
						r.PostMessage(ctx, w.AppId, f.Conversation.ID, p)
						err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
						if err.RowsAffected == 0 {
							fmt.Println("rows not updated.")
						}
					}

				} else {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: w.Message,
						},
					}
					r.PostMessage(ctx, w.AppId, f.Conversation.ID, p)
					err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
					if err.RowsAffected == 0 {
						fmt.Println("rows not updated.")
					}

				}

			} else {
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: w.TriggerMessage,
						},
					}
					r.PostMessage(ctx, w.AppId, f.Conversation.ID, p)
					return &models.Response{Received: &f}, nil
				}
				fmt.Println("appUserId already exist.")
				return &models.Response{Msg: "Userid already exist."}, nil

			}
		} else {
			fmt.Println("user NOT Registered.")
		}

	} else if f.Messages[0].Source.Type == "twitter" {
		db := r.DBConn.Where("twitter_integration_id = ?", f.Messages[0].Source.IntegrationID).Find(&tw)
		if db.Error != nil {
			fmt.Println("error")
		}
		tenant := r.DBConn.Table("account_details").Where("domain_uuid = ?", tw.Domain_uuid).Find(&T)
		if tenant.Error != nil {
			fmt.Println("error")
		}
		if myDate.Weekday().String() == tw.Day1 {
			workstart1 := tw.Workstart1
			components := strings.Split(workstart1, ":")
			StartHour, _ := components[0], components[1]
			workend1 := tw.Workend1
			components1 := strings.Split(workend1, ":")
			EndHour, _ := components1[0], components1[1]
			startHour, _ := strconv.Atoi(StartHour)
			endHour, _ := strconv.Atoi(EndHour)
			if myDate.Hour() < startHour || myDate.Hour() > endHour {
				if aot := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; aot != nil {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: tw.Message,
						},
					}
					r.PostMessage(ctx, tw.AppId, f.Conversation.ID, p)
					db := r.DBConn.Create(&u).Where("app_user_id = ?", f.AppUser.ID).Update("after_office_time", true)
					if db.Error != nil {
						fmt.Println(db.Error)
					}
				}

				if myDate.Weekday().String() == u.Day && date == u.Date {
					if u.AfterOfficeTime == true {
						fmt.Println("message already sent.")
					} else if u.AfterOfficeTime == false {
						p := models.User{
							Author: models.Author{
								Type:        "business",
								DisplayName: T.Tenant_name,
								AvatarURL:   "https://www.gravatar.com/image.jpg",
							},
							Content: models.Content{
								Type: "text",
								Text: tw.Message,
							},
						}
						r.PostMessage(ctx, tw.AppId, f.Conversation.ID, p)
						err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
						if err.RowsAffected == 0 {
							fmt.Println("rows not updated.")
						}
					}

				} else {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: tw.Message,
						},
					}
					r.PostMessage(ctx, tw.AppId, f.Conversation.ID, p)
					err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
					if err.RowsAffected == 0 {
						fmt.Println("rows not updated.")
					}

				}

			} else {
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}

					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: tw.TriggerMessage,
						},
					}
					r.PostMessage(ctx, tw.AppId, f.Conversation.ID, p)
					return &models.Response{Received: &f}, nil
				}
				fmt.Println("appUserId already exist.")
				return &models.Response{Msg: "Userid already exist."}, nil

			}
		} else if myDate.Weekday().String() == tw.Day2 {
			workstart2 := tw.Workstart2
			components := strings.Split(workstart2, ":")
			StartHour, _ := components[0], components[1]
			workend2 := tw.Workend2
			components1 := strings.Split(workend2, ":")
			EndHour, _ := components1[0], components1[1]
			startHour, _ := strconv.Atoi(StartHour)
			endHour, _ := strconv.Atoi(EndHour)
			if myDate.Hour() < startHour || myDate.Hour() > endHour {
				if aot := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; aot != nil {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: tw.Message,
						},
					}
					r.PostMessage(ctx, tw.AppId, f.Conversation.ID, p)
					db := r.DBConn.Create(&u).Where("app_user_id = ?", f.AppUser.ID).Update("after_office_time", true)
					if db.Error != nil {
						fmt.Println(db.Error)
					}
				}

				if myDate.Weekday().String() == u.Day && date == u.Date {
					if u.AfterOfficeTime == true {
						fmt.Println("message already sent.")
					} else if u.AfterOfficeTime == false {
						p := models.User{
							Author: models.Author{
								Type:        "business",
								DisplayName: T.Tenant_name,
								AvatarURL:   "https://www.gravatar.com/image.jpg",
							},
							Content: models.Content{
								Type: "text",
								Text: tw.Message,
							},
						}
						r.PostMessage(ctx, tw.AppId, f.Conversation.ID, p)
						err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
						if err.RowsAffected == 0 {
							fmt.Println("rows not updated.")
						}
					}

				} else {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: tw.Message,
						},
					}
					r.PostMessage(ctx, tw.AppId, f.Conversation.ID, p)
					err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
					if err.RowsAffected == 0 {
						fmt.Println("rows not updated.")
					}

				}

			} else {
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}

					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: tw.TriggerMessage,
						},
					}
					r.PostMessage(ctx, tw.AppId, f.Conversation.ID, p)
					return &models.Response{Received: &f}, nil
				}
				fmt.Println("appUserId already exist.")
				return &models.Response{Msg: "Userid already exist."}, nil

			}
		} else if myDate.Weekday().String() == tw.Day3 {
			workstart3 := tw.Workstart3
			components := strings.Split(workstart3, ":")
			StartHour, _ := components[0], components[1]
			workend3 := tw.Workend3
			components1 := strings.Split(workend3, ":")
			EndHour, _ := components1[0], components1[1]
			startHour, _ := strconv.Atoi(StartHour)
			endHour, _ := strconv.Atoi(EndHour)
			if myDate.Hour() < startHour || myDate.Hour() > endHour {
				if aot := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; aot != nil {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: tw.Message,
						},
					}
					r.PostMessage(ctx, tw.AppId, f.Conversation.ID, p)
					db := r.DBConn.Create(&u).Where("app_user_id = ?", f.AppUser.ID).Update("after_office_time", true)
					if db.Error != nil {
						fmt.Println(db.Error)
					}
				}

				if myDate.Weekday().String() == u.Day && date == u.Date {
					if u.AfterOfficeTime == true {
						fmt.Println("message already sent.")
					} else if u.AfterOfficeTime == false {
						p := models.User{
							Author: models.Author{
								Type:        "business",
								DisplayName: T.Tenant_name,
								AvatarURL:   "https://www.gravatar.com/image.jpg",
							},
							Content: models.Content{
								Type: "text",
								Text: tw.Message,
							},
						}
						r.PostMessage(ctx, tw.AppId, f.Conversation.ID, p)
						err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
						if err.RowsAffected == 0 {
							fmt.Println("rows not updated.")
						}
					}

				} else {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: tw.Message,
						},
					}
					r.PostMessage(ctx, tw.AppId, f.Conversation.ID, p)
					err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
					if err.RowsAffected == 0 {
						fmt.Println("rows not updated.")
					}

				}

			} else {
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}

					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: tw.TriggerMessage,
						},
					}
					r.PostMessage(ctx, tw.AppId, f.Conversation.ID, p)
					return &models.Response{Received: &f}, nil
				}
				fmt.Println("appUserId already exist.")
				return &models.Response{Msg: "Userid already exist."}, nil

			}
		} else if myDate.Weekday().String() == tw.Day4 {
			workstart4 := tw.Workstart4
			components := strings.Split(workstart4, ":")
			StartHour, _ := components[0], components[1]
			workend4 := tw.Workend4
			components1 := strings.Split(workend4, ":")
			EndHour, _ := components1[0], components1[1]
			startHour, _ := strconv.Atoi(StartHour)
			endHour, _ := strconv.Atoi(EndHour)
			if myDate.Hour() < startHour || myDate.Hour() > endHour {
				if aot := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; aot != nil {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: tw.Message,
						},
					}
					r.PostMessage(ctx, tw.AppId, f.Conversation.ID, p)
					db := r.DBConn.Create(&u).Where("app_user_id = ?", f.AppUser.ID).Update("after_office_time", true)
					if db.Error != nil {
						fmt.Println(db.Error)
					}
				}

				if myDate.Weekday().String() == u.Day && date == u.Date {
					if u.AfterOfficeTime == true {
						fmt.Println("message already sent.")
					} else if u.AfterOfficeTime == false {
						p := models.User{
							Author: models.Author{
								Type:        "business",
								DisplayName: T.Tenant_name,
								AvatarURL:   "https://www.gravatar.com/image.jpg",
							},
							Content: models.Content{
								Type: "text",
								Text: tw.Message,
							},
						}
						r.PostMessage(ctx, tw.AppId, f.Conversation.ID, p)
						err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
						if err.RowsAffected == 0 {
							fmt.Println("rows not updated.")
						}
					}

				} else {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: tw.Message,
						},
					}
					r.PostMessage(ctx, tw.AppId, f.Conversation.ID, p)
					err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
					if err.RowsAffected == 0 {
						fmt.Println("rows not updated.")
					}

				}

			} else {
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}

					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: tw.TriggerMessage,
						},
					}
					r.PostMessage(ctx, tw.AppId, f.Conversation.ID, p)
					return &models.Response{Received: &f}, nil
				}
				fmt.Println("appUserId already exist.")
				return &models.Response{Msg: "Userid already exist."}, nil

			}
		} else if myDate.Weekday().String() == tw.Day5 {
			workstart5 := tw.Workstart5
			components := strings.Split(workstart5, ":")
			StartHour, _ := components[0], components[1]
			workend5 := tw.Workend5
			components1 := strings.Split(workend5, ":")
			EndHour, _ := components1[0], components1[1]
			startHour, _ := strconv.Atoi(StartHour)
			endHour, _ := strconv.Atoi(EndHour)
			if myDate.Hour() < startHour || myDate.Hour() > endHour {
				if aot := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; aot != nil {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: tw.Message,
						},
					}
					r.PostMessage(ctx, tw.AppId, f.Conversation.ID, p)
					db := r.DBConn.Create(&u).Where("app_user_id = ?", f.AppUser.ID).Update("after_office_time", true)
					if db.Error != nil {
						fmt.Println(db.Error)
					}
				}

				if myDate.Weekday().String() == u.Day && date == u.Date {
					if u.AfterOfficeTime == true {
						fmt.Println("message already sent.")
					} else if u.AfterOfficeTime == false {
						p := models.User{
							Author: models.Author{
								Type:        "business",
								DisplayName: T.Tenant_name,
								AvatarURL:   "https://www.gravatar.com/image.jpg",
							},
							Content: models.Content{
								Type: "text",
								Text: tw.Message,
							},
						}
						r.PostMessage(ctx, tw.AppId, f.Conversation.ID, p)
						err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
						if err.RowsAffected == 0 {
							fmt.Println("rows not updated.")
						}
					}

				} else {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: tw.Message,
						},
					}
					r.PostMessage(ctx, tw.AppId, f.Conversation.ID, p)
					err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
					if err.RowsAffected == 0 {
						fmt.Println("rows not updated.")
					}

				}

			} else {
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}

					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: tw.TriggerMessage,
						},
					}
					r.PostMessage(ctx, tw.AppId, f.Conversation.ID, p)
					return &models.Response{Received: &f}, nil
				}
				fmt.Println("appUserId already exist.")
				return &models.Response{Msg: "Userid already exist."}, nil

			}
		} else if myDate.Weekday().String() == tw.Day6 {
			workstart6 := tw.Workstart6
			components := strings.Split(workstart6, ":")
			StartHour, _ := components[0], components[1]
			workend6 := tw.Workend6
			components1 := strings.Split(workend6, ":")
			EndHour, _ := components1[0], components1[1]
			startHour, _ := strconv.Atoi(StartHour)
			endHour, _ := strconv.Atoi(EndHour)
			if myDate.Hour() < startHour || myDate.Hour() > endHour {
				if aot := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; aot != nil {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: tw.Message,
						},
					}
					r.PostMessage(ctx, tw.AppId, f.Conversation.ID, p)
					db := r.DBConn.Create(&u).Where("app_user_id = ?", f.AppUser.ID).Update("after_office_time", true)
					if db.Error != nil {
						fmt.Println(db.Error)
					}
				}

				if myDate.Weekday().String() == u.Day && date == u.Date {
					if u.AfterOfficeTime == true {
						fmt.Println("message already sent.")
					} else if u.AfterOfficeTime == false {
						p := models.User{
							Author: models.Author{
								Type:        "business",
								DisplayName: T.Tenant_name,
								AvatarURL:   "https://www.gravatar.com/image.jpg",
							},
							Content: models.Content{
								Type: "text",
								Text: tw.Message,
							},
						}
						r.PostMessage(ctx, tw.AppId, f.Conversation.ID, p)
						err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
						if err.RowsAffected == 0 {
							fmt.Println("rows not updated.")
						}
					}

				} else {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: tw.Message,
						},
					}
					r.PostMessage(ctx, tw.AppId, f.Conversation.ID, p)
					err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
					if err.RowsAffected == 0 {
						fmt.Println("rows not updated.")
					}

				}

			} else {
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}

					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: tw.TriggerMessage,
						},
					}
					r.PostMessage(ctx, tw.AppId, f.Conversation.ID, p)
					return &models.Response{Received: &f}, nil
				}
				fmt.Println("appUserId already exist.")
				return &models.Response{Msg: "Userid already exist."}, nil

			}
		} else if myDate.Weekday().String() == tw.Day7 {
			workstart7 := tw.Workstart7
			components := strings.Split(workstart7, ":")
			StartHour, _ := components[0], components[1]
			workend7 := tw.Workend7
			components1 := strings.Split(workend7, ":")
			EndHour, _ := components1[0], components1[1]
			startHour, _ := strconv.Atoi(StartHour)
			endHour, _ := strconv.Atoi(EndHour)
			if myDate.Hour() < startHour || myDate.Hour() > endHour {
				if aot := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; aot != nil {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: tw.Message,
						},
					}
					r.PostMessage(ctx, tw.AppId, f.Conversation.ID, p)
					db := r.DBConn.Create(&u).Where("app_user_id = ?", f.AppUser.ID).Update("after_office_time", true)
					if db.Error != nil {
						fmt.Println(db.Error)
					}
				}

				if myDate.Weekday().String() == u.Day && date == u.Date {
					if u.AfterOfficeTime == true {
						fmt.Println("message already sent.")
					} else if u.AfterOfficeTime == false {
						p := models.User{
							Author: models.Author{
								Type:        "business",
								DisplayName: T.Tenant_name,
								AvatarURL:   "https://www.gravatar.com/image.jpg",
							},
							Content: models.Content{
								Type: "text",
								Text: tw.Message,
							},
						}
						r.PostMessage(ctx, tw.AppId, f.Conversation.ID, p)
						err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
						if err.RowsAffected == 0 {
							fmt.Println("rows not updated.")
						}
					}

				} else {
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: tw.Message,
						},
					}
					r.PostMessage(ctx, tw.AppId, f.Conversation.ID, p)
					err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Updates(map[string]interface{}{"day": myDate.Weekday().String(), "date": date, "after_office_time": true})
					if err.RowsAffected == 0 {
						fmt.Println("rows not updated.")
					}

				}

			} else {
				if err := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", f.AppUser.ID).Find(&u).Error; err != nil {
					db := r.DBConn.Create(&u)
					if db.Error != nil {

					}
					p := models.User{
						Author: models.Author{
							Type:        "business",
							DisplayName: T.Tenant_name,
							AvatarURL:   "https://www.gravatar.com/image.jpg",
						},
						Content: models.Content{
							Type: "text",
							Text: tw.TriggerMessage,
						},
					}
					r.PostMessage(ctx, tw.AppId, f.Conversation.ID, p)
					return &models.Response{Received: &f}, nil
				}
				fmt.Println("appUserId already exist.")
				return &models.Response{Msg: "Userid already exist."}, nil

			}
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
func (r *crudRepository) Update_AppUser(ctx context.Context, appUserId string, appId string, surname string, givenName string) (*models.Response, error) {
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
func (r *crudRepository) Add_Smooch_configuration(ctx context.Context, name string, domain_uuid string, appId string, appKey string, appSecret string) (*models.Response, error) {

	td := models.Tenant_details{
		ConfigurationName: name,
		Domain_uuid:       domain_uuid,
		AppId:             appId,
		AppKey:            appKey,
		AppSecret:         appSecret,
	}

	if db := r.DBConn.Table("tenant_details").Where("app_id = ?", appId).Or("configuration_name = ?", name).Find(&td).Error; db != nil {
		st := r.DBConn.Create(&td)
		if st.Error != nil {
			return &models.Response{ResponseCode: 409, Status: "Error", Msg: "Not created"}, nil
		}
		return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Smooch configuration saved successfully."}, nil

	}
	if td.AppId == appId {
		return &models.Response{ResponseCode: 409, Status: "Error", Msg: "AppId Already Exist."}, nil
	} else if td.ConfigurationName == name {
		return &models.Response{ResponseCode: 409, Status: "Error", Msg: "Configuration name Already Exist."}, nil
	}
	return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Smooch configuration saved successfully."}, nil

}

/***************************************Add smooch configuration*****************************************/
func (r *crudRepository) Update_Smooch_configuration(ctx context.Context, id int64, name string, domain_uuid string, appId string, appKey string, appSecret string) (*models.Response, error) {
	w := models.Tenant_details{}
	db := r.DBConn.Where("domain_uuid = ? AND id = ?", domain_uuid, id).Find(&w)
	if db.RowsAffected == 0 {
		return &models.Response{Status: "Not Found", Msg: "Record Doesn't Exist", ResponseCode: 401}, nil
	}
	if appId != w.AppId && name != w.ConfigurationName {
		fmt.Println("part1")
		if db := r.DBConn.Table("tenant_details").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"app_id": appId, "app_key": appKey, "app_secret": appSecret, "configuration_name": name}).Error; db != nil {
			return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
		}

		return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Smooch configuration Updated successfully."}, nil

	} else if appId == w.AppId && name == w.ConfigurationName {
		fmt.Println("part2")
		if db := r.DBConn.Table("tenant_details").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"app_key": appKey, "app_secret": appSecret}).Error; db != nil {
			return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
		}

		return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Smooch configuration Updated successfully."}, nil
	} else if appId != w.AppId && name == w.ConfigurationName {
		fmt.Println("part3")
		if db := r.DBConn.Table("tenant_details").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"app_id": appId, "app_key": appKey, "app_secret": appSecret}).Error; db != nil {
			return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
		}

		return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Smooch configuration Updated successfully."}, nil

	} else if appId == w.AppId && name != w.ConfigurationName {
		fmt.Println("part4")
		if db := r.DBConn.Table("tenant_details").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"app_key": appKey, "app_secret": appSecret, "configuration_name": name}).Error; db != nil {
			return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
		}

		return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Smooch configuration Updated successfully."}, nil

	}
	return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Smooch configuration Updated successfully."}, nil
}

/***************************************Get smooch configuration*****************************************/
func (r *crudRepository) Get_Smooch_configuration(ctx context.Context, domain_uuid string) (*models.Response, error) {
	w := models.Tenant_details{}
	list := make([]models.Tenant_details, 0)
	db := r.DBConn.Where("domain_uuid = ?", domain_uuid).Find(&w)
	if db.Error != nil {
		return &models.Response{Status: "0", Msg: "Record Not Found", ResponseCode: 401}, nil

	}
	row, err := r.DBConn.Raw("select id, configuration_name, domain_uuid,app_id, app_key, app_secret from tenant_details WHERE domain_uuid = ?", domain_uuid).Rows()
	if err != nil {
		return &models.Response{Status: "0", Msg: "Record Not Found", ResponseCode: 401}, nil
	}
	defer row.Close()
	for row.Next() {
		f := models.Tenant_details{}
		if err := row.Scan(&f.Id, &f.ConfigurationName, &f.Domain_uuid, &f.AppId, &f.AppKey, &f.AppSecret); err != nil {

			return nil, err
		}

		list = append(list, f)
	}
	return &models.Response{Status: "1", Msg: "Record Found", ResponseCode: 200, Tenant_list: list}, nil
}

/**************************************************Delete Smooch configuration***************************/
func (r *crudRepository) Delete_Smooch_configuration(ctx context.Context, id int64, domain_uuid string) (*models.Response, error) {
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
func (r *crudRepository) Add_Whatsapp_configuration(ctx context.Context, td models.WhatsappConfigurations) (*models.Response, error) {

	ts := models.WhatsappConfiguration{
		Domain_uuid:           td.Domain_uuid,
		ConfigurationName:     td.ConfigurationName,
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
	if err := r.DBConn.Table("whatsapp_configurations").Where("app_id = ?", td.AppId).Or("whatsapp_integration_id = ?", td.WhatsappIntegrationID).Or("configuration_name = ?", td.ConfigurationName).Find(&ts).Error; err != nil {
		st := r.DBConn.Create(&ts)
		if st.Error != nil {
			return &models.Response{ResponseCode: 409, Status: "Error", Msg: "Not created"}, nil
		}
		return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Whatsapp configuration saved successfully."}, nil
	}
	if ts.AppId == td.AppId {
		return &models.Response{ResponseCode: 409, Status: "Error", Msg: "AppId Already Exist."}, nil
	} else if ts.WhatsappIntegrationID == td.WhatsappIntegrationID {
		return &models.Response{ResponseCode: 409, Status: "Error", Msg: "Whatsapp Integration Id  Already Exist."}, nil
	} else if ts.ConfigurationName == td.ConfigurationName {
		return &models.Response{ResponseCode: 409, Status: "Error", Msg: "configuration name Already Exist."}, nil
	}
	return &models.Response{ResponseCode: 201, Status: "Ok", Msg: "Whatsapp configuration saved successfully."}, nil
}

/**********************************************Get appID by tenant_domain_uuid******************************/
func (r *crudRepository) Get_Whatsapp_configuration(ctx context.Context, domain_uuid string) (*models.Response, error) {
	w := models.WhatsappConfiguration{}
	list := make([]models.WhatsappConfiguration, 0)
	db := r.DBConn.Where("domain_uuid = ?", domain_uuid).Find(&w)
	if db.Error != nil {
		return &models.Response{Status: "0", Msg: "Record Not Found", ResponseCode: 401}, nil

	}

	row, err := r.DBConn.Raw("select * from whatsapp_configurations WHERE domain_uuid = ?", domain_uuid).Rows()
	if err != nil {
		return &models.Response{Status: "0", Msg: "Record Not Found", ResponseCode: 401}, nil
	}
	defer row.Close()
	for row.Next() {
		f := models.WhatsappConfiguration{}
		if err := row.Scan(&f.Id, &f.ConfigurationName, &f.Domain_uuid, &f.AppId, &f.AppKey, &f.AppSecret, &f.Message, &f.Size, &f.WhatsappIntegrationID, &f.Day1, &f.Day2, &f.Day3, &f.Day4, &f.Day5, &f.Day6, &f.Day7, &f.Workstart1, &f.Workstart2, &f.Workstart3, &f.Workstart4, &f.Workstart5, &f.Workstart6, &f.Workstart7, &f.Workend1, &f.Workend2, &f.Workend3, &f.Workend4, &f.Workend5, &f.Workend6, &f.Workend7, &f.TriggerWhen, &f.TriggerName, &f.TriggerMessage); err != nil {

			return nil, err
		}
		list = append(list, f)
	}
	return &models.Response{Status: "1", Msg: "Record Found", ResponseCode: 200, List: list}, nil
}

/**********************************************Update Tenant details*************************************/
func (r *crudRepository) Update_Whatsapp_configuration(ctx context.Context, id int64, domain_uuid string, td models.WhatsappConfigurations) (*models.Response, error) {
	w := models.WhatsappConfiguration{}
	db1 := r.DBConn.Where("domain_uuid = ? AND id = ?", domain_uuid, id).Find(&w)
	if db1.Error != nil {
		return &models.Response{Status: "Not Found", Msg: "Record Doesn't Exist", ResponseCode: 401}, nil
	}

	if td.WhatsappIntegrationID == w.WhatsappIntegrationID && td.AppId == w.AppId && td.ConfigurationName == w.ConfigurationName {
		fmt.Println("part1")

		if db := r.DBConn.Table("whatsapp_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"app_key": td.AppKey, "app_secret": td.AppSecret, "message": td.Message, "size": td.Size, "trigger_name": td.Trigger.Name, "trigger_message": td.Trigger.Message, "trigger_when": td.Trigger.When, "day1": td.WorkingDays[0].Day, "day2": td.WorkingDays[1].Day, "day3": td.WorkingDays[2].Day, "day4": td.WorkingDays[3].Day, "day5": td.WorkingDays[4].Day, "day6": td.WorkingDays[5].Day, "day7": td.WorkingDays[6].Day, "workstart1": td.WorkingDays[0].WorkingHourStartTime, "workstart2": td.WorkingDays[1].WorkingHourStartTime, "workstart3": td.WorkingDays[2].WorkingHourStartTime, "workstart4": td.WorkingDays[3].WorkingHourStartTime, "workstart5": td.WorkingDays[4].WorkingHourStartTime, "workstart6": td.WorkingDays[5].WorkingHourStartTime, "workstart7": td.WorkingDays[6].WorkingHourStartTime, "workend1": td.WorkingDays[0].WorkingHourEndTime, "workend2": td.WorkingDays[1].WorkingHourEndTime, "workend3": td.WorkingDays[2].WorkingHourEndTime, "workend4": td.WorkingDays[3].WorkingHourEndTime, "workend5": td.WorkingDays[4].WorkingHourEndTime, "workend6": td.WorkingDays[5].WorkingHourEndTime, "workend7": td.WorkingDays[6].WorkingHourEndTime}).Error; db != nil {
			return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
		}

		return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Whatsapp integration Updated successfully."}, nil
	} else if td.WhatsappIntegrationID != w.WhatsappIntegrationID && td.AppId != w.AppId && td.ConfigurationName != w.ConfigurationName {
		fmt.Println("part2")
		w := models.WhatsappConfiguration{}
		if row := r.DBConn.Table("whatsapp_configurations").Where("app_id = ?", td.AppId).Or("whatsapp_integration_id = ?", td.WhatsappIntegrationID).Or("configuration_name = ?", td.ConfigurationName).Find(&w).Error; row != nil {
			//if err := r.DBConn.Table("whatsapp_configurations").Where("whatsapp_integration_id = ?", td.WhatsappIntegrationID).Find(&w).Error; err != nil {
			if db := r.DBConn.Table("whatsapp_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"configuration_name": td.ConfigurationName, "app_id": td.AppId, "app_key": td.AppKey, "app_secret": td.AppSecret, "message": td.Message, "whatsapp_integration_id": td.WhatsappIntegrationID, "size": td.Size, "trigger_name": td.Trigger.Name, "trigger_message": td.Trigger.Message, "trigger_when": td.Trigger.When, "day1": td.WorkingDays[0].Day, "day2": td.WorkingDays[1].Day, "day3": td.WorkingDays[2].Day, "day4": td.WorkingDays[3].Day, "day5": td.WorkingDays[4].Day, "day6": td.WorkingDays[5].Day, "day7": td.WorkingDays[6].Day, "workstart1": td.WorkingDays[0].WorkingHourStartTime, "workstart2": td.WorkingDays[1].WorkingHourStartTime, "workstart3": td.WorkingDays[2].WorkingHourStartTime, "workstart4": td.WorkingDays[3].WorkingHourStartTime, "workstart5": td.WorkingDays[4].WorkingHourStartTime, "workstart6": td.WorkingDays[5].WorkingHourStartTime, "workstart7": td.WorkingDays[6].WorkingHourStartTime, "workend1": td.WorkingDays[0].WorkingHourEndTime, "workend2": td.WorkingDays[1].WorkingHourEndTime, "workend3": td.WorkingDays[2].WorkingHourEndTime, "workend4": td.WorkingDays[3].WorkingHourEndTime, "workend5": td.WorkingDays[4].WorkingHourEndTime, "workend6": td.WorkingDays[5].WorkingHourEndTime, "workend7": td.WorkingDays[6].WorkingHourEndTime}).Error; db != nil {
				return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
			}

			return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Whatsapp integration Updated successfully."}, nil
		}
		if w.WhatsappIntegrationID == td.WhatsappIntegrationID {
			return &models.Response{ResponseCode: 409, Status: "0", Msg: "Whatsapp id already exist."}, nil

		} else if w.AppId == td.AppId {
			return &models.Response{ResponseCode: 409, Status: "0", Msg: "AppId already exist."}, nil
		} else if w.ConfigurationName == td.ConfigurationName {
			return &models.Response{ResponseCode: 409, Status: "0", Msg: "Configuration name already exist."}, nil
		}
	} else if td.WhatsappIntegrationID == w.WhatsappIntegrationID && td.AppId != w.AppId && td.ConfigurationName != w.ConfigurationName {
		fmt.Println("part3")
		w := models.WhatsappConfiguration{}
		if row := r.DBConn.Where("app_id = ?", td.AppId).Or("configuration_name = ?", td.ConfigurationName).Find(&w).Error; row != nil {
			if db := r.DBConn.Table("whatsapp_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"configuration_name": td.ConfigurationName, "app_id": td.AppId, "app_key": td.AppKey, "app_secret": td.AppSecret, "message": td.Message, "size": td.Size, "trigger_name": td.Trigger.Name, "trigger_message": td.Trigger.Message, "trigger_when": td.Trigger.When, "day1": td.WorkingDays[0].Day, "day2": td.WorkingDays[1].Day, "day3": td.WorkingDays[2].Day, "day4": td.WorkingDays[3].Day, "day5": td.WorkingDays[4].Day, "day6": td.WorkingDays[5].Day, "day7": td.WorkingDays[6].Day, "workstart1": td.WorkingDays[0].WorkingHourStartTime, "workstart2": td.WorkingDays[1].WorkingHourStartTime, "workstart3": td.WorkingDays[2].WorkingHourStartTime, "workstart4": td.WorkingDays[3].WorkingHourStartTime, "workstart5": td.WorkingDays[4].WorkingHourStartTime, "workstart6": td.WorkingDays[5].WorkingHourStartTime, "workstart7": td.WorkingDays[6].WorkingHourStartTime, "workend1": td.WorkingDays[0].WorkingHourEndTime, "workend2": td.WorkingDays[1].WorkingHourEndTime, "workend3": td.WorkingDays[2].WorkingHourEndTime, "workend4": td.WorkingDays[3].WorkingHourEndTime, "workend5": td.WorkingDays[4].WorkingHourEndTime, "workend6": td.WorkingDays[5].WorkingHourEndTime, "workend7": td.WorkingDays[6].WorkingHourEndTime}).Error; db != nil {
				return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
			}

			return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Whatsapp integration Updated successfully."}, nil
		}
		if w.AppId == td.AppId {
			return &models.Response{ResponseCode: 409, Status: "0", Msg: "AppId already exist."}, nil
		} else if w.ConfigurationName == td.ConfigurationName {
			return &models.Response{ResponseCode: 409, Status: "0", Msg: "Configuration name already exist."}, nil
		}
	} else if td.WhatsappIntegrationID == w.WhatsappIntegrationID && td.AppId == w.AppId && td.ConfigurationName != w.ConfigurationName {
		fmt.Println("part4")
		w := models.WhatsappConfiguration{}
		if err := r.DBConn.Where("configuration_name = ?", td.ConfigurationName).Find(&w).Error; err != nil {
			if db := r.DBConn.Table("whatsapp_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"configuration_name": td.ConfigurationName, "app_key": td.AppKey, "app_secret": td.AppSecret, "message": td.Message, "size": td.Size, "trigger_name": td.Trigger.Name, "trigger_message": td.Trigger.Message, "trigger_when": td.Trigger.When, "day1": td.WorkingDays[0].Day, "day2": td.WorkingDays[1].Day, "day3": td.WorkingDays[2].Day, "day4": td.WorkingDays[3].Day, "day5": td.WorkingDays[4].Day, "day6": td.WorkingDays[5].Day, "day7": td.WorkingDays[6].Day, "workstart1": td.WorkingDays[0].WorkingHourStartTime, "workstart2": td.WorkingDays[1].WorkingHourStartTime, "workstart3": td.WorkingDays[2].WorkingHourStartTime, "workstart4": td.WorkingDays[3].WorkingHourStartTime, "workstart5": td.WorkingDays[4].WorkingHourStartTime, "workstart6": td.WorkingDays[5].WorkingHourStartTime, "workstart7": td.WorkingDays[6].WorkingHourStartTime, "workend1": td.WorkingDays[0].WorkingHourEndTime, "workend2": td.WorkingDays[1].WorkingHourEndTime, "workend3": td.WorkingDays[2].WorkingHourEndTime, "workend4": td.WorkingDays[3].WorkingHourEndTime, "workend5": td.WorkingDays[4].WorkingHourEndTime, "workend6": td.WorkingDays[5].WorkingHourEndTime, "workend7": td.WorkingDays[6].WorkingHourEndTime}).Error; db != nil {
				return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
			}

			return &models.Response{ResponseCode: 200, Status: "OK", Msg: "Whatsapp integration Updated successfully."}, nil
		}
		return &models.Response{ResponseCode: 409, Status: "0", Msg: "Configuration name already exist."}, nil
	} else if td.WhatsappIntegrationID != w.WhatsappIntegrationID && td.AppId != w.AppId && td.ConfigurationName == w.ConfigurationName {
		fmt.Println("part5")
		w := models.WhatsappConfiguration{}
		if err := r.DBConn.Where("whatsapp_integration_id = ?", td.WhatsappIntegrationID).Or("app_id = ?", td.AppId).Find(&w).Error; err != nil {
			if db := r.DBConn.Table("whatsapp_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"app_id": td.AppId, "app_key": td.AppKey, "app_secret": td.AppSecret, "message": td.Message, "whatsapp_integration_id": td.WhatsappIntegrationID, "size": td.Size, "trigger_name": td.Trigger.Name, "trigger_message": td.Trigger.Message, "trigger_when": td.Trigger.When, "day1": td.WorkingDays[0].Day, "day2": td.WorkingDays[1].Day, "day3": td.WorkingDays[2].Day, "day4": td.WorkingDays[3].Day, "day5": td.WorkingDays[4].Day, "day6": td.WorkingDays[5].Day, "day7": td.WorkingDays[6].Day, "workstart1": td.WorkingDays[0].WorkingHourStartTime, "workstart2": td.WorkingDays[1].WorkingHourStartTime, "workstart3": td.WorkingDays[2].WorkingHourStartTime, "workstart4": td.WorkingDays[3].WorkingHourStartTime, "workstart5": td.WorkingDays[4].WorkingHourStartTime, "workstart6": td.WorkingDays[5].WorkingHourStartTime, "workstart7": td.WorkingDays[6].WorkingHourStartTime, "workend1": td.WorkingDays[0].WorkingHourEndTime, "workend2": td.WorkingDays[1].WorkingHourEndTime, "workend3": td.WorkingDays[2].WorkingHourEndTime, "workend4": td.WorkingDays[3].WorkingHourEndTime, "workend5": td.WorkingDays[4].WorkingHourEndTime, "workend6": td.WorkingDays[5].WorkingHourEndTime, "workend7": td.WorkingDays[6].WorkingHourEndTime}).Error; db != nil {
				return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
			}

			return &models.Response{ResponseCode: 200, Status: "OK", Msg: "Whatsapp integration Updated successfully."}, nil
		}
		if w.WhatsappIntegrationID == td.WhatsappIntegrationID {
			return &models.Response{ResponseCode: 409, Status: "0", Msg: "Whatsapp id already exist."}, nil
		} else if w.AppId == td.AppId {
			return &models.Response{ResponseCode: 409, Status: "0", Msg: "AppId already exist."}, nil
		}
	} else if td.WhatsappIntegrationID != w.WhatsappIntegrationID && td.AppId == w.AppId && td.ConfigurationName == w.ConfigurationName {
		fmt.Println("part6")
		w := models.WhatsappConfiguration{}
		if err := r.DBConn.Where("whatsapp_integration_id = ?", td.WhatsappIntegrationID).Find(&w).Error; err != nil {
			if db := r.DBConn.Table("whatsapp_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"app_key": td.AppKey, "app_secret": td.AppSecret, "message": td.Message, "whatsapp_integration_id": td.WhatsappIntegrationID, "size": td.Size, "trigger_name": td.Trigger.Name, "trigger_message": td.Trigger.Message, "trigger_when": td.Trigger.When, "day1": td.WorkingDays[0].Day, "day2": td.WorkingDays[1].Day, "day3": td.WorkingDays[2].Day, "day4": td.WorkingDays[3].Day, "day5": td.WorkingDays[4].Day, "day6": td.WorkingDays[5].Day, "day7": td.WorkingDays[6].Day, "workstart1": td.WorkingDays[0].WorkingHourStartTime, "workstart2": td.WorkingDays[1].WorkingHourStartTime, "workstart3": td.WorkingDays[2].WorkingHourStartTime, "workstart4": td.WorkingDays[3].WorkingHourStartTime, "workstart5": td.WorkingDays[4].WorkingHourStartTime, "workstart6": td.WorkingDays[5].WorkingHourStartTime, "workstart7": td.WorkingDays[6].WorkingHourStartTime, "workend1": td.WorkingDays[0].WorkingHourEndTime, "workend2": td.WorkingDays[1].WorkingHourEndTime, "workend3": td.WorkingDays[2].WorkingHourEndTime, "workend4": td.WorkingDays[3].WorkingHourEndTime, "workend5": td.WorkingDays[4].WorkingHourEndTime, "workend6": td.WorkingDays[5].WorkingHourEndTime, "workend7": td.WorkingDays[6].WorkingHourEndTime}).Error; db != nil {
				return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
			}

			return &models.Response{ResponseCode: 200, Status: "OK", Msg: "Whatsapp integration Updated successfully."}, nil
		}
		return &models.Response{ResponseCode: 409, Status: "0", Msg: "Whatsapp id already exist."}, nil
	} else if td.WhatsappIntegrationID != w.WhatsappIntegrationID && td.AppId == w.AppId && td.ConfigurationName != w.ConfigurationName {
		fmt.Println("part7")
		w := models.WhatsappConfiguration{}
		if err := r.DBConn.Where("whatsapp_integration_id = ?", td.WhatsappIntegrationID).Or("configuration_name = ?", td.ConfigurationName).Find(&w).Error; err != nil {
			if db := r.DBConn.Table("whatsapp_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"configuration_name": td.ConfigurationName, "app_key": td.AppKey, "app_secret": td.AppSecret, "message": td.Message, "whatsapp_integration_id": td.WhatsappIntegrationID, "size": td.Size, "trigger_name": td.Trigger.Name, "trigger_message": td.Trigger.Message, "trigger_when": td.Trigger.When, "day1": td.WorkingDays[0].Day, "day2": td.WorkingDays[1].Day, "day3": td.WorkingDays[2].Day, "day4": td.WorkingDays[3].Day, "day5": td.WorkingDays[4].Day, "day6": td.WorkingDays[5].Day, "day7": td.WorkingDays[6].Day, "workstart1": td.WorkingDays[0].WorkingHourStartTime, "workstart2": td.WorkingDays[1].WorkingHourStartTime, "workstart3": td.WorkingDays[2].WorkingHourStartTime, "workstart4": td.WorkingDays[3].WorkingHourStartTime, "workstart5": td.WorkingDays[4].WorkingHourStartTime, "workstart6": td.WorkingDays[5].WorkingHourStartTime, "workstart7": td.WorkingDays[6].WorkingHourStartTime, "workend1": td.WorkingDays[0].WorkingHourEndTime, "workend2": td.WorkingDays[1].WorkingHourEndTime, "workend3": td.WorkingDays[2].WorkingHourEndTime, "workend4": td.WorkingDays[3].WorkingHourEndTime, "workend5": td.WorkingDays[4].WorkingHourEndTime, "workend6": td.WorkingDays[5].WorkingHourEndTime, "workend7": td.WorkingDays[6].WorkingHourEndTime}).Error; db != nil {
				return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
			}

			return &models.Response{ResponseCode: 200, Status: "OK", Msg: "Whatsapp integration Updated successfully."}, nil
		}
		if w.WhatsappIntegrationID == td.WhatsappIntegrationID {
			return &models.Response{ResponseCode: 409, Status: "0", Msg: "Whatsapp id already exist."}, nil
		} else if w.ConfigurationName == td.ConfigurationName {
			return &models.Response{ResponseCode: 409, Status: "0", Msg: "Configuration name already exist."}, nil
		}
	} else if td.WhatsappIntegrationID == w.WhatsappIntegrationID && td.AppId != w.AppId && td.ConfigurationName == w.ConfigurationName {
		fmt.Println("part8")
		w := models.WhatsappConfiguration{}
		if err := r.DBConn.Where("app_id = ?", td.AppId).Find(&w).Error; err != nil {
			if db := r.DBConn.Table("whatsapp_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"app_id": td.AppId, "app_key": td.AppKey, "app_secret": td.AppSecret, "message": td.Message, "size": td.Size, "trigger_name": td.Trigger.Name, "trigger_message": td.Trigger.Message, "trigger_when": td.Trigger.When, "day1": td.WorkingDays[0].Day, "day2": td.WorkingDays[1].Day, "day3": td.WorkingDays[2].Day, "day4": td.WorkingDays[3].Day, "day5": td.WorkingDays[4].Day, "day6": td.WorkingDays[5].Day, "day7": td.WorkingDays[6].Day, "workstart1": td.WorkingDays[0].WorkingHourStartTime, "workstart2": td.WorkingDays[1].WorkingHourStartTime, "workstart3": td.WorkingDays[2].WorkingHourStartTime, "workstart4": td.WorkingDays[3].WorkingHourStartTime, "workstart5": td.WorkingDays[4].WorkingHourStartTime, "workstart6": td.WorkingDays[5].WorkingHourStartTime, "workstart7": td.WorkingDays[6].WorkingHourStartTime, "workend1": td.WorkingDays[0].WorkingHourEndTime, "workend2": td.WorkingDays[1].WorkingHourEndTime, "workend3": td.WorkingDays[2].WorkingHourEndTime, "workend4": td.WorkingDays[3].WorkingHourEndTime, "workend5": td.WorkingDays[4].WorkingHourEndTime, "workend6": td.WorkingDays[5].WorkingHourEndTime, "workend7": td.WorkingDays[6].WorkingHourEndTime}).Error; db != nil {
				return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
			}

			return &models.Response{ResponseCode: 200, Status: "OK", Msg: "Whatsapp integration Updated successfully."}, nil
		}
		return &models.Response{ResponseCode: 409, Status: "0", Msg: "AppId already exist."}, nil
	} else {
		return &models.Response{ResponseCode: 401, Status: "Not Ok", Msg: "Whatsapp integration Not Updated."}, nil
	}
	return &models.Response{ResponseCode: 204, Status: "OK", Msg: "Whatsapp integration Updated successfully."}, nil

}

/**********************************************Delete Tenant details*************************************/
func (r *crudRepository) Delete_Whatsapp_configuration(ctx context.Context, id int64, domain_uuid string) (*models.Response, error) {
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
func (r *crudRepository) Add_Facebook_configuration(ctx context.Context, td models.FacebookConfigurations) (*models.Response, error) {

	ts := models.FacebookConfiguration{
		Domain_uuid:           td.Domain_uuid,
		ConfigurationName:     td.ConfigurationName,
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
	if err := r.DBConn.Table("facebook_configurations").Where("app_id = ?", td.AppId).Or("facebook_integration_id = ?", td.FacebookIntegrationID).Or("configuration_name = ?", td.ConfigurationName).Find(&ts).Error; err != nil {
		st := r.DBConn.Create(&ts)
		if st.Error != nil {
			return &models.Response{ResponseCode: 409, Status: "Error", Msg: "Not created"}, nil
		}
		return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Facebook configuration saved successfully."}, nil

	}
	if ts.FacebookIntegrationID == td.FacebookIntegrationID {
		return &models.Response{ResponseCode: 409, Status: "Error", Msg: "Facebook Integration Id  Already Exist."}, nil
	} else if ts.AppId == td.AppId {
		return &models.Response{ResponseCode: 409, Status: "Error", Msg: "AppId  Already Exist."}, nil
	} else if ts.ConfigurationName == td.ConfigurationName {
		return &models.Response{ResponseCode: 409, Status: "Error", Msg: "Configuration name Already Exist."}, nil
	}
	return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Facebook configuration saved successfully."}, nil
}

/**********************************************Get appID by tenant_domain_uuid******************************/
func (r *crudRepository) Get_Facebook_configuration(ctx context.Context, domain_uuid string) (*models.Response, error) {
	list := make([]models.FacebookConfiguration, 0)
	w := models.FacebookConfiguration{}
	db := r.DBConn.Where("domain_uuid = ?", domain_uuid).Find(&w)
	if db.Error != nil {
		return &models.Response{Status: "0", Msg: "Record Not Found", ResponseCode: 401}, nil

	}
	row, err := r.DBConn.Raw("select * from facebook_configurations WHERE domain_uuid = ?", domain_uuid).Rows()
	if err != nil {
		return &models.Response{Status: "0", Msg: "Record Not Found", ResponseCode: 401}, nil
	}
	defer row.Close()
	for row.Next() {
		f := models.FacebookConfiguration{}
		if err := row.Scan(&f.Id, &f.ConfigurationName, &f.Domain_uuid, &f.AppId, &f.AppKey, &f.AppSecret, &f.Message, &f.Size, &f.FacebookIntegrationID, &f.Day1, &f.Day2, &f.Day3, &f.Day4, &f.Day5, &f.Day6, &f.Day7, &f.Workstart1, &f.Workstart2, &f.Workstart3, &f.Workstart4, &f.Workstart5, &f.Workstart6, &f.Workstart7, &f.Workend1, &f.Workend2, &f.Workend3, &f.Workend4, &f.Workend5, &f.Workend6, &f.Workend7, &f.TriggerWhen, &f.TriggerName, &f.TriggerMessage); err != nil {

			return nil, err
		}

		list = append(list, f)
	}
	return &models.Response{Status: "1", Msg: "Record Found", ResponseCode: 200, Fb: list}, nil
}

/**********************************************Update Tenant details*************************************/
func (r *crudRepository) Update_Facebook_configuration(ctx context.Context, id int64, domain_uuid string, td models.FacebookConfigurations) (*models.Response, error) {
	w := models.FacebookConfiguration{}
	db1 := r.DBConn.Where("domain_uuid = ? AND id = ?", domain_uuid, id).Find(&w)
	if db1.RowsAffected == 0 {
		return &models.Response{Status: "Not Found", Msg: "Record Doesn't Exist", ResponseCode: 401}, nil
	}
	if td.FacebookIntegrationID == w.FacebookIntegrationID && td.AppId == w.AppId && td.ConfigurationName == w.ConfigurationName {
		fmt.Println("part1")
		if db := r.DBConn.Table("facebook_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"app_key": td.AppKey, "app_secret": td.AppSecret, "size": td.Size, "trigger_name": td.Trigger.Name, "trigger_message": td.Trigger.Message, "trigger_when": td.Trigger.When, "day1": td.WorkingDays[0].Day, "day2": td.WorkingDays[1].Day, "day3": td.WorkingDays[2].Day, "day4": td.WorkingDays[3].Day, "day5": td.WorkingDays[4].Day, "day6": td.WorkingDays[5].Day, "day7": td.WorkingDays[6].Day, "workstart1": td.WorkingDays[0].WorkingHourStartTime, "workstart2": td.WorkingDays[1].WorkingHourStartTime, "workstart3": td.WorkingDays[2].WorkingHourStartTime, "workstart4": td.WorkingDays[3].WorkingHourStartTime, "workstart5": td.WorkingDays[4].WorkingHourStartTime, "workstart6": td.WorkingDays[5].WorkingHourStartTime, "workstart7": td.WorkingDays[6].WorkingHourStartTime, "workend1": td.WorkingDays[0].WorkingHourEndTime, "workend2": td.WorkingDays[1].WorkingHourEndTime, "workend3": td.WorkingDays[2].WorkingHourEndTime, "workend4": td.WorkingDays[3].WorkingHourEndTime, "workend5": td.WorkingDays[4].WorkingHourEndTime, "workend6": td.WorkingDays[5].WorkingHourEndTime, "workend7": td.WorkingDays[6].WorkingHourEndTime, "message": td.Message}).Error; db != nil {
			return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
		}

		return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Facebook integration Updated successfully."}, nil
	} else if td.FacebookIntegrationID != w.FacebookIntegrationID && td.AppId != w.AppId && td.ConfigurationName != w.ConfigurationName {
		fmt.Println("part2")
		w := models.FacebookConfiguration{}
		if row := r.DBConn.Where("app_id = ?", td.AppId).Or("facebook_integration_id = ?", td.FacebookIntegrationID).Or("configuration_name = ?", td.ConfigurationName).Find(&w).Error; row != nil {
			//if err := r.DBConn.Where("facebook_integration_id = ?", td.FacebookIntegrationID).Find(&w).Error; err != nil {
			if db := r.DBConn.Table("facebook_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"configuration_name": td.ConfigurationName, "app_id": td.AppId, "app_key": td.AppKey, "app_secret": td.AppSecret, "facebook_integration_id": td.FacebookIntegrationID, "size": td.Size, "trigger_name": td.Trigger.Name, "trigger_message": td.Trigger.Message, "trigger_when": td.Trigger.When, "day1": td.WorkingDays[0].Day, "day2": td.WorkingDays[1].Day, "day3": td.WorkingDays[2].Day, "day4": td.WorkingDays[3].Day, "day5": td.WorkingDays[4].Day, "day6": td.WorkingDays[5].Day, "day7": td.WorkingDays[6].Day, "workstart1": td.WorkingDays[0].WorkingHourStartTime, "workstart2": td.WorkingDays[1].WorkingHourStartTime, "workstart3": td.WorkingDays[2].WorkingHourStartTime, "workstart4": td.WorkingDays[3].WorkingHourStartTime, "workstart5": td.WorkingDays[4].WorkingHourStartTime, "workstart6": td.WorkingDays[5].WorkingHourStartTime, "workstart7": td.WorkingDays[6].WorkingHourStartTime, "workend1": td.WorkingDays[0].WorkingHourEndTime, "workend2": td.WorkingDays[1].WorkingHourEndTime, "workend3": td.WorkingDays[2].WorkingHourEndTime, "workend4": td.WorkingDays[3].WorkingHourEndTime, "workend5": td.WorkingDays[4].WorkingHourEndTime, "workend6": td.WorkingDays[5].WorkingHourEndTime, "workend7": td.WorkingDays[6].WorkingHourEndTime, "message": td.Message}).Error; db != nil {
				return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
			}

			return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Facebook integration Updated successfully."}, nil
		}
		if td.FacebookIntegrationID == w.FacebookIntegrationID {
			return &models.Response{ResponseCode: 409, Status: "0", Msg: "Facebook id already exist."}, nil
		} else if td.AppId == w.AppId {
			return &models.Response{ResponseCode: 409, Status: "0", Msg: "AppId already exist."}, nil
		} else if td.ConfigurationName == w.ConfigurationName {
			return &models.Response{ResponseCode: 409, Status: "0", Msg: "Configuration name already exist."}, nil
		}
	} else if td.FacebookIntegrationID == w.FacebookIntegrationID && td.AppId != w.AppId && td.ConfigurationName != w.ConfigurationName {
		fmt.Println("part3")
		w := models.FacebookConfiguration{}
		if err := r.DBConn.Where("app_id = ?", td.AppId).Or("configuration_name = ?", td.ConfigurationName).Find(&w).Error; err != nil {
			if db := r.DBConn.Table("facebook_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"configuration_name": td.ConfigurationName, "app_id": td.AppId, "app_key": td.AppKey, "app_secret": td.AppSecret, "size": td.Size, "trigger_name": td.Trigger.Name, "trigger_message": td.Trigger.Message, "trigger_when": td.Trigger.When, "day1": td.WorkingDays[0].Day, "day2": td.WorkingDays[1].Day, "day3": td.WorkingDays[2].Day, "day4": td.WorkingDays[3].Day, "day5": td.WorkingDays[4].Day, "day6": td.WorkingDays[5].Day, "day7": td.WorkingDays[6].Day, "workstart1": td.WorkingDays[0].WorkingHourStartTime, "workstart2": td.WorkingDays[1].WorkingHourStartTime, "workstart3": td.WorkingDays[2].WorkingHourStartTime, "workstart4": td.WorkingDays[3].WorkingHourStartTime, "workstart5": td.WorkingDays[4].WorkingHourStartTime, "workstart6": td.WorkingDays[5].WorkingHourStartTime, "workstart7": td.WorkingDays[6].WorkingHourStartTime, "workend1": td.WorkingDays[0].WorkingHourEndTime, "workend2": td.WorkingDays[1].WorkingHourEndTime, "workend3": td.WorkingDays[2].WorkingHourEndTime, "workend4": td.WorkingDays[3].WorkingHourEndTime, "workend5": td.WorkingDays[4].WorkingHourEndTime, "workend6": td.WorkingDays[5].WorkingHourEndTime, "workend7": td.WorkingDays[6].WorkingHourEndTime, "message": td.Message}).Error; db != nil {
				return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
			}

			return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Facebook integration Updated successfully."}, nil
		}
		if td.AppId == w.AppId {
			return &models.Response{ResponseCode: 409, Status: "0", Msg: "AppId already exist."}, nil
		} else if td.ConfigurationName == w.ConfigurationName {
			return &models.Response{ResponseCode: 409, Status: "0", Msg: "Configuration name already exist."}, nil
		}
	} else if td.FacebookIntegrationID == w.FacebookIntegrationID && td.AppId == w.AppId && td.ConfigurationName != w.ConfigurationName {
		fmt.Println("part4")
		w := models.FacebookConfiguration{}
		if err := r.DBConn.Where("configuration_name = ?", td.ConfigurationName).Find(&w).Error; err != nil {
			if db := r.DBConn.Table("facebook_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"configuration_name": td.ConfigurationName, "app_key": td.AppKey, "app_secret": td.AppSecret, "facebook_integration_id": td.FacebookIntegrationID, "size": td.Size, "trigger_name": td.Trigger.Name, "trigger_message": td.Trigger.Message, "trigger_when": td.Trigger.When, "day1": td.WorkingDays[0].Day, "day2": td.WorkingDays[1].Day, "day3": td.WorkingDays[2].Day, "day4": td.WorkingDays[3].Day, "day5": td.WorkingDays[4].Day, "day6": td.WorkingDays[5].Day, "day7": td.WorkingDays[6].Day, "workstart1": td.WorkingDays[0].WorkingHourStartTime, "workstart2": td.WorkingDays[1].WorkingHourStartTime, "workstart3": td.WorkingDays[2].WorkingHourStartTime, "workstart4": td.WorkingDays[3].WorkingHourStartTime, "workstart5": td.WorkingDays[4].WorkingHourStartTime, "workstart6": td.WorkingDays[5].WorkingHourStartTime, "workstart7": td.WorkingDays[6].WorkingHourStartTime, "workend1": td.WorkingDays[0].WorkingHourEndTime, "workend2": td.WorkingDays[1].WorkingHourEndTime, "workend3": td.WorkingDays[2].WorkingHourEndTime, "workend4": td.WorkingDays[3].WorkingHourEndTime, "workend5": td.WorkingDays[4].WorkingHourEndTime, "workend6": td.WorkingDays[5].WorkingHourEndTime, "workend7": td.WorkingDays[6].WorkingHourEndTime, "message": td.Message}).Error; db != nil {
				return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
			}

			return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Facebook integration Updated successfully."}, nil
		}
		return &models.Response{ResponseCode: 409, Status: "0", Msg: "Configuration name already exist."}, nil
	} else if td.FacebookIntegrationID != w.FacebookIntegrationID && td.AppId != w.AppId && td.ConfigurationName == w.ConfigurationName {
		fmt.Println("part5")
		w := models.FacebookConfiguration{}
		if err := r.DBConn.Where("facebook_integration_id = ?", td.FacebookIntegrationID).Or("app_id = ?", td.AppId).Find(&w).Error; err != nil {
			if db := r.DBConn.Table("facebook_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"app_id": td.AppId, "app_key": td.AppKey, "app_secret": td.AppSecret, "facebook_integration_id": td.FacebookIntegrationID, "size": td.Size, "trigger_name": td.Trigger.Name, "trigger_message": td.Trigger.Message, "trigger_when": td.Trigger.When, "day1": td.WorkingDays[0].Day, "day2": td.WorkingDays[1].Day, "day3": td.WorkingDays[2].Day, "day4": td.WorkingDays[3].Day, "day5": td.WorkingDays[4].Day, "day6": td.WorkingDays[5].Day, "day7": td.WorkingDays[6].Day, "workstart1": td.WorkingDays[0].WorkingHourStartTime, "workstart2": td.WorkingDays[1].WorkingHourStartTime, "workstart3": td.WorkingDays[2].WorkingHourStartTime, "workstart4": td.WorkingDays[3].WorkingHourStartTime, "workstart5": td.WorkingDays[4].WorkingHourStartTime, "workstart6": td.WorkingDays[5].WorkingHourStartTime, "workstart7": td.WorkingDays[6].WorkingHourStartTime, "workend1": td.WorkingDays[0].WorkingHourEndTime, "workend2": td.WorkingDays[1].WorkingHourEndTime, "workend3": td.WorkingDays[2].WorkingHourEndTime, "workend4": td.WorkingDays[3].WorkingHourEndTime, "workend5": td.WorkingDays[4].WorkingHourEndTime, "workend6": td.WorkingDays[5].WorkingHourEndTime, "workend7": td.WorkingDays[6].WorkingHourEndTime, "message": td.Message}).Error; db != nil {
				return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
			}

			return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Facebook integration Updated successfully."}, nil
		}
		if td.FacebookIntegrationID == w.FacebookIntegrationID {
			return &models.Response{ResponseCode: 409, Status: "0", Msg: "Facebook id already exist."}, nil
		} else if td.AppId == w.AppId {
			return &models.Response{ResponseCode: 409, Status: "0", Msg: "AppId already exist."}, nil
		}
	} else if td.FacebookIntegrationID != w.FacebookIntegrationID && td.AppId == w.AppId && td.ConfigurationName == w.ConfigurationName {
		fmt.Println("part6")
		w := models.FacebookConfiguration{}
		if err := r.DBConn.Where("facebook_integration_id = ?", td.FacebookIntegrationID).Find(&w).Error; err != nil {
			if db := r.DBConn.Table("facebook_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"app_key": td.AppKey, "app_secret": td.AppSecret, "facebook_integration_id": td.FacebookIntegrationID, "size": td.Size, "trigger_name": td.Trigger.Name, "trigger_message": td.Trigger.Message, "trigger_when": td.Trigger.When, "day1": td.WorkingDays[0].Day, "day2": td.WorkingDays[1].Day, "day3": td.WorkingDays[2].Day, "day4": td.WorkingDays[3].Day, "day5": td.WorkingDays[4].Day, "day6": td.WorkingDays[5].Day, "day7": td.WorkingDays[6].Day, "workstart1": td.WorkingDays[0].WorkingHourStartTime, "workstart2": td.WorkingDays[1].WorkingHourStartTime, "workstart3": td.WorkingDays[2].WorkingHourStartTime, "workstart4": td.WorkingDays[3].WorkingHourStartTime, "workstart5": td.WorkingDays[4].WorkingHourStartTime, "workstart6": td.WorkingDays[5].WorkingHourStartTime, "workstart7": td.WorkingDays[6].WorkingHourStartTime, "workend1": td.WorkingDays[0].WorkingHourEndTime, "workend2": td.WorkingDays[1].WorkingHourEndTime, "workend3": td.WorkingDays[2].WorkingHourEndTime, "workend4": td.WorkingDays[3].WorkingHourEndTime, "workend5": td.WorkingDays[4].WorkingHourEndTime, "workend6": td.WorkingDays[5].WorkingHourEndTime, "workend7": td.WorkingDays[6].WorkingHourEndTime, "message": td.Message}).Error; db != nil {
				return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
			}

			return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Facebook integration Updated successfully."}, nil
		}
		return &models.Response{ResponseCode: 409, Status: "0", Msg: "Facebook id already exist."}, nil
	} else if td.FacebookIntegrationID != w.FacebookIntegrationID && td.AppId == w.AppId && td.ConfigurationName != w.ConfigurationName {
		fmt.Println("part7")
		w := models.FacebookConfiguration{}
		if err := r.DBConn.Where("facebook_integration_id = ?", td.FacebookIntegrationID).Or("configuration_name = ?", td.ConfigurationName).Find(&w).Error; err != nil {
			if db := r.DBConn.Table("facebook_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"configuration_name": td.ConfigurationName, "app_key": td.AppKey, "app_secret": td.AppSecret, "facebook_integration_id": td.FacebookIntegrationID, "size": td.Size, "trigger_name": td.Trigger.Name, "trigger_message": td.Trigger.Message, "trigger_when": td.Trigger.When, "day1": td.WorkingDays[0].Day, "day2": td.WorkingDays[1].Day, "day3": td.WorkingDays[2].Day, "day4": td.WorkingDays[3].Day, "day5": td.WorkingDays[4].Day, "day6": td.WorkingDays[5].Day, "day7": td.WorkingDays[6].Day, "workstart1": td.WorkingDays[0].WorkingHourStartTime, "workstart2": td.WorkingDays[1].WorkingHourStartTime, "workstart3": td.WorkingDays[2].WorkingHourStartTime, "workstart4": td.WorkingDays[3].WorkingHourStartTime, "workstart5": td.WorkingDays[4].WorkingHourStartTime, "workstart6": td.WorkingDays[5].WorkingHourStartTime, "workstart7": td.WorkingDays[6].WorkingHourStartTime, "workend1": td.WorkingDays[0].WorkingHourEndTime, "workend2": td.WorkingDays[1].WorkingHourEndTime, "workend3": td.WorkingDays[2].WorkingHourEndTime, "workend4": td.WorkingDays[3].WorkingHourEndTime, "workend5": td.WorkingDays[4].WorkingHourEndTime, "workend6": td.WorkingDays[5].WorkingHourEndTime, "workend7": td.WorkingDays[6].WorkingHourEndTime, "message": td.Message}).Error; db != nil {
				return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
			}

			return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Facebook integration Updated successfully."}, nil
		}
		if td.FacebookIntegrationID == w.FacebookIntegrationID {
			return &models.Response{ResponseCode: 409, Status: "0", Msg: "Facebook id already exist."}, nil
		} else if td.ConfigurationName == w.ConfigurationName {
			return &models.Response{ResponseCode: 409, Status: "0", Msg: "Configuration name already exist."}, nil
		}
		return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Facebook integration Updated successfully."}, nil
	} else if td.FacebookIntegrationID == w.FacebookIntegrationID && td.AppId != w.AppId && td.ConfigurationName == w.ConfigurationName {
		fmt.Println("part8")
		w := models.FacebookConfiguration{}
		if err := r.DBConn.Where("app_id = ?", td.AppId).Find(&w).Error; err != nil {
			if db := r.DBConn.Table("facebook_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"app_id": td.AppId, "app_key": td.AppKey, "app_secret": td.AppSecret, "size": td.Size, "trigger_name": td.Trigger.Name, "trigger_message": td.Trigger.Message, "trigger_when": td.Trigger.When, "day1": td.WorkingDays[0].Day, "day2": td.WorkingDays[1].Day, "day3": td.WorkingDays[2].Day, "day4": td.WorkingDays[3].Day, "day5": td.WorkingDays[4].Day, "day6": td.WorkingDays[5].Day, "day7": td.WorkingDays[6].Day, "workstart1": td.WorkingDays[0].WorkingHourStartTime, "workstart2": td.WorkingDays[1].WorkingHourStartTime, "workstart3": td.WorkingDays[2].WorkingHourStartTime, "workstart4": td.WorkingDays[3].WorkingHourStartTime, "workstart5": td.WorkingDays[4].WorkingHourStartTime, "workstart6": td.WorkingDays[5].WorkingHourStartTime, "workstart7": td.WorkingDays[6].WorkingHourStartTime, "workend1": td.WorkingDays[0].WorkingHourEndTime, "workend2": td.WorkingDays[1].WorkingHourEndTime, "workend3": td.WorkingDays[2].WorkingHourEndTime, "workend4": td.WorkingDays[3].WorkingHourEndTime, "workend5": td.WorkingDays[4].WorkingHourEndTime, "workend6": td.WorkingDays[5].WorkingHourEndTime, "workend7": td.WorkingDays[6].WorkingHourEndTime, "message": td.Message}).Error; db != nil {
				return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
			}

			return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Facebook integration Updated successfully."}, nil
		}

		return &models.Response{ResponseCode: 409, Status: "0", Msg: "AppId already exist."}, nil
	}
	return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Facebook integration Updated successfully."}, nil
}

/**********************************************Delete Tenant details*************************************/
func (r *crudRepository) Delete_Facebook_configuration(ctx context.Context, id int64, domain_uuid string) (*models.Response, error) {
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

/****************************************Save tenant details*********************************************/
func (r *crudRepository) Add_Twitter_configuration(ctx context.Context, td models.TwitterConfigurations) (*models.Response, error) {

	ts := models.TwitterConfiguration{
		Domain_uuid:          td.Domain_uuid,
		ConfigurationName:    td.ConfigurationName,
		AppId:                td.AppId,
		AppKey:               td.AppKey,
		AppSecret:            td.AppSecret,
		Message:              td.Message,
		Size:                 td.Size,
		TriggerWhen:          td.Trigger.When,
		TriggerName:          td.Trigger.Name,
		TriggerMessage:       td.Trigger.Message,
		TwitterIntegrationID: td.TwitterIntegrationID,
		Day1:                 td.WorkingDays[0].Day,
		Day2:                 td.WorkingDays[1].Day,
		Day3:                 td.WorkingDays[2].Day,
		Day4:                 td.WorkingDays[3].Day,
		Day5:                 td.WorkingDays[4].Day,
		Day6:                 td.WorkingDays[5].Day,
		Day7:                 td.WorkingDays[6].Day,
		Workstart1:           td.WorkingDays[0].WorkingHourStartTime,
		Workstart2:           td.WorkingDays[1].WorkingHourStartTime,
		Workstart3:           td.WorkingDays[2].WorkingHourStartTime,
		Workstart4:           td.WorkingDays[3].WorkingHourStartTime,
		Workstart5:           td.WorkingDays[4].WorkingHourStartTime,
		Workstart6:           td.WorkingDays[5].WorkingHourStartTime,
		Workstart7:           td.WorkingDays[6].WorkingHourStartTime,
		Workend1:             td.WorkingDays[0].WorkingHourEndTime,
		Workend2:             td.WorkingDays[1].WorkingHourEndTime,
		Workend3:             td.WorkingDays[2].WorkingHourEndTime,
		Workend4:             td.WorkingDays[3].WorkingHourEndTime,
		Workend5:             td.WorkingDays[4].WorkingHourEndTime,
		Workend6:             td.WorkingDays[5].WorkingHourEndTime,
		Workend7:             td.WorkingDays[6].WorkingHourEndTime,
	}
	if err := r.DBConn.Table("twitter_configurations").Where("app_id = ?", td.AppId).Or("twitter_integration_id = ?", td.TwitterIntegrationID).Or("configuration_name = ?", td.ConfigurationName).Find(&ts).Error; err != nil {
		st := r.DBConn.Create(&ts)
		if st.Error != nil {
			return &models.Response{ResponseCode: 409, Status: "Error", Msg: "Not created"}, nil
		}
		return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Twitter configuration saved successfully."}, nil
	}
	if ts.AppId == td.AppId {
		return &models.Response{ResponseCode: 409, Status: "Error", Msg: "AppId Already Exist."}, nil
	} else if ts.TwitterIntegrationID == td.TwitterIntegrationID {
		return &models.Response{ResponseCode: 409, Status: "Error", Msg: "Twitter Integration Id  Already Exist."}, nil
	} else if ts.ConfigurationName == td.ConfigurationName {
		return &models.Response{ResponseCode: 409, Status: "Error", Msg: "configuration name Already Exist."}, nil
	}
	return &models.Response{ResponseCode: 201, Status: "Ok", Msg: "Twitter configuration saved successfully."}, nil
}

/**********************************************Get appID by tenant_domain_uuid******************************/
func (r *crudRepository) Get_Twitter_configuration(ctx context.Context, domain_uuid string) (*models.Response, error) {
	w := models.TwitterConfiguration{}
	list := make([]models.TwitterConfiguration, 0)
	db := r.DBConn.Where("domain_uuid = ?", domain_uuid).Find(&w)
	if db.Error != nil {
		return &models.Response{Status: "0", Msg: "Record Not Found", ResponseCode: 401}, nil

	}

	row, err := r.DBConn.Raw("select * from twitter_configurations WHERE domain_uuid = ?", domain_uuid).Rows()
	if err != nil {
		return &models.Response{Status: "0", Msg: "Record Not Found", ResponseCode: 401}, nil
	}
	defer row.Close()
	for row.Next() {
		f := models.TwitterConfiguration{}
		if err := row.Scan(&f.Id, &f.ConfigurationName, &f.Domain_uuid, &f.AppId, &f.AppKey, &f.AppSecret, &f.Message, &f.Size, &f.TwitterIntegrationID, &f.Day1, &f.Day2, &f.Day3, &f.Day4, &f.Day5, &f.Day6, &f.Day7, &f.Workstart1, &f.Workstart2, &f.Workstart3, &f.Workstart4, &f.Workstart5, &f.Workstart6, &f.Workstart7, &f.Workend1, &f.Workend2, &f.Workend3, &f.Workend4, &f.Workend5, &f.Workend6, &f.Workend7, &f.TriggerWhen, &f.TriggerName, &f.TriggerMessage); err != nil {

			return nil, err
		}
		list = append(list, f)
	}
	return &models.Response{Status: "1", Msg: "Record Found", ResponseCode: 200, Twitter: list}, nil
}

/**********************************************Update Tenant details*************************************/
func (r *crudRepository) Update_Twitter_configuration(ctx context.Context, id int64, domain_uuid string, td models.TwitterConfigurations) (*models.Response, error) {
	w := models.TwitterConfiguration{}
	db1 := r.DBConn.Where("domain_uuid = ? AND id = ?", domain_uuid, id).Find(&w)
	if db1.Error != nil {
		return &models.Response{Status: "Not Found", Msg: "Record Doesn't Exist", ResponseCode: 401}, nil
	}

	if td.TwitterIntegrationID == w.TwitterIntegrationID && td.AppId == w.AppId && td.ConfigurationName == w.ConfigurationName {
		fmt.Println("part1")

		if db := r.DBConn.Table("twitter_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"app_key": td.AppKey, "app_secret": td.AppSecret, "message": td.Message, "size": td.Size, "trigger_name": td.Trigger.Name, "trigger_message": td.Trigger.Message, "trigger_when": td.Trigger.When, "day1": td.WorkingDays[0].Day, "day2": td.WorkingDays[1].Day, "day3": td.WorkingDays[2].Day, "day4": td.WorkingDays[3].Day, "day5": td.WorkingDays[4].Day, "day6": td.WorkingDays[5].Day, "day7": td.WorkingDays[6].Day, "workstart1": td.WorkingDays[0].WorkingHourStartTime, "workstart2": td.WorkingDays[1].WorkingHourStartTime, "workstart3": td.WorkingDays[2].WorkingHourStartTime, "workstart4": td.WorkingDays[3].WorkingHourStartTime, "workstart5": td.WorkingDays[4].WorkingHourStartTime, "workstart6": td.WorkingDays[5].WorkingHourStartTime, "workstart7": td.WorkingDays[6].WorkingHourStartTime, "workend1": td.WorkingDays[0].WorkingHourEndTime, "workend2": td.WorkingDays[1].WorkingHourEndTime, "workend3": td.WorkingDays[2].WorkingHourEndTime, "workend4": td.WorkingDays[3].WorkingHourEndTime, "workend5": td.WorkingDays[4].WorkingHourEndTime, "workend6": td.WorkingDays[5].WorkingHourEndTime, "workend7": td.WorkingDays[6].WorkingHourEndTime}).Error; db != nil {
			return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
		}

		return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Twitter integration Updated successfully."}, nil
	} else if td.TwitterIntegrationID != w.TwitterIntegrationID && td.AppId != w.AppId && td.ConfigurationName != w.ConfigurationName {
		fmt.Println("part2")
		w := models.TwitterConfiguration{}
		if row := r.DBConn.Table("twitter_configurations").Where("app_id = ?", td.AppId).Or("twitter_integration_id = ?", td.TwitterIntegrationID).Or("configuration_name = ?", td.ConfigurationName).Find(&w).Error; row != nil {
			//if err := r.DBConn.Table("whatsapp_configurations").Where("whatsapp_integration_id = ?", td.WhatsappIntegrationID).Find(&w).Error; err != nil {
			if db := r.DBConn.Table("twitter_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"configuration_name": td.ConfigurationName, "app_id": td.AppId, "app_key": td.AppKey, "app_secret": td.AppSecret, "message": td.Message, "twitter_integration_id": td.TwitterIntegrationID, "size": td.Size, "trigger_name": td.Trigger.Name, "trigger_message": td.Trigger.Message, "trigger_when": td.Trigger.When, "day1": td.WorkingDays[0].Day, "day2": td.WorkingDays[1].Day, "day3": td.WorkingDays[2].Day, "day4": td.WorkingDays[3].Day, "day5": td.WorkingDays[4].Day, "day6": td.WorkingDays[5].Day, "day7": td.WorkingDays[6].Day, "workstart1": td.WorkingDays[0].WorkingHourStartTime, "workstart2": td.WorkingDays[1].WorkingHourStartTime, "workstart3": td.WorkingDays[2].WorkingHourStartTime, "workstart4": td.WorkingDays[3].WorkingHourStartTime, "workstart5": td.WorkingDays[4].WorkingHourStartTime, "workstart6": td.WorkingDays[5].WorkingHourStartTime, "workstart7": td.WorkingDays[6].WorkingHourStartTime, "workend1": td.WorkingDays[0].WorkingHourEndTime, "workend2": td.WorkingDays[1].WorkingHourEndTime, "workend3": td.WorkingDays[2].WorkingHourEndTime, "workend4": td.WorkingDays[3].WorkingHourEndTime, "workend5": td.WorkingDays[4].WorkingHourEndTime, "workend6": td.WorkingDays[5].WorkingHourEndTime, "workend7": td.WorkingDays[6].WorkingHourEndTime}).Error; db != nil {
				return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
			}

			return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Twitter integration Updated successfully."}, nil
		}
		if w.TwitterIntegrationID == td.TwitterIntegrationID {
			return &models.Response{ResponseCode: 409, Status: "0", Msg: "Twitter id already exist."}, nil

		} else if w.AppId == td.AppId {
			return &models.Response{ResponseCode: 409, Status: "0", Msg: "AppId already exist."}, nil
		} else if w.ConfigurationName == td.ConfigurationName {
			return &models.Response{ResponseCode: 409, Status: "0", Msg: "Configuration name already exist."}, nil
		}
	} else if td.TwitterIntegrationID == w.TwitterIntegrationID && td.AppId != w.AppId && td.ConfigurationName != w.ConfigurationName {
		fmt.Println("part3")
		w := models.TwitterConfiguration{}
		if row := r.DBConn.Where("app_id = ?", td.AppId).Or("configuration_name = ?", td.ConfigurationName).Find(&w).Error; row != nil {
			if db := r.DBConn.Table("twitter_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"configuration_name": td.ConfigurationName, "app_id": td.AppId, "app_key": td.AppKey, "app_secret": td.AppSecret, "message": td.Message, "size": td.Size, "trigger_name": td.Trigger.Name, "trigger_message": td.Trigger.Message, "trigger_when": td.Trigger.When, "day1": td.WorkingDays[0].Day, "day2": td.WorkingDays[1].Day, "day3": td.WorkingDays[2].Day, "day4": td.WorkingDays[3].Day, "day5": td.WorkingDays[4].Day, "day6": td.WorkingDays[5].Day, "day7": td.WorkingDays[6].Day, "workstart1": td.WorkingDays[0].WorkingHourStartTime, "workstart2": td.WorkingDays[1].WorkingHourStartTime, "workstart3": td.WorkingDays[2].WorkingHourStartTime, "workstart4": td.WorkingDays[3].WorkingHourStartTime, "workstart5": td.WorkingDays[4].WorkingHourStartTime, "workstart6": td.WorkingDays[5].WorkingHourStartTime, "workstart7": td.WorkingDays[6].WorkingHourStartTime, "workend1": td.WorkingDays[0].WorkingHourEndTime, "workend2": td.WorkingDays[1].WorkingHourEndTime, "workend3": td.WorkingDays[2].WorkingHourEndTime, "workend4": td.WorkingDays[3].WorkingHourEndTime, "workend5": td.WorkingDays[4].WorkingHourEndTime, "workend6": td.WorkingDays[5].WorkingHourEndTime, "workend7": td.WorkingDays[6].WorkingHourEndTime}).Error; db != nil {
				return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
			}

			return &models.Response{ResponseCode: 201, Status: "OK", Msg: "Twitter integration Updated successfully."}, nil
		}
		if w.AppId == td.AppId {
			return &models.Response{ResponseCode: 409, Status: "0", Msg: "AppId already exist."}, nil
		} else if w.ConfigurationName == td.ConfigurationName {
			return &models.Response{ResponseCode: 409, Status: "0", Msg: "Configuration name already exist."}, nil
		}
	} else if td.TwitterIntegrationID == w.TwitterIntegrationID && td.AppId == w.AppId && td.ConfigurationName != w.ConfigurationName {
		fmt.Println("part4")
		w := models.TwitterConfiguration{}
		if err := r.DBConn.Where("configuration_name = ?", td.ConfigurationName).Find(&w).Error; err != nil {
			if db := r.DBConn.Table("twitter_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"configuration_name": td.ConfigurationName, "app_key": td.AppKey, "app_secret": td.AppSecret, "message": td.Message, "size": td.Size, "trigger_name": td.Trigger.Name, "trigger_message": td.Trigger.Message, "trigger_when": td.Trigger.When, "day1": td.WorkingDays[0].Day, "day2": td.WorkingDays[1].Day, "day3": td.WorkingDays[2].Day, "day4": td.WorkingDays[3].Day, "day5": td.WorkingDays[4].Day, "day6": td.WorkingDays[5].Day, "day7": td.WorkingDays[6].Day, "workstart1": td.WorkingDays[0].WorkingHourStartTime, "workstart2": td.WorkingDays[1].WorkingHourStartTime, "workstart3": td.WorkingDays[2].WorkingHourStartTime, "workstart4": td.WorkingDays[3].WorkingHourStartTime, "workstart5": td.WorkingDays[4].WorkingHourStartTime, "workstart6": td.WorkingDays[5].WorkingHourStartTime, "workstart7": td.WorkingDays[6].WorkingHourStartTime, "workend1": td.WorkingDays[0].WorkingHourEndTime, "workend2": td.WorkingDays[1].WorkingHourEndTime, "workend3": td.WorkingDays[2].WorkingHourEndTime, "workend4": td.WorkingDays[3].WorkingHourEndTime, "workend5": td.WorkingDays[4].WorkingHourEndTime, "workend6": td.WorkingDays[5].WorkingHourEndTime, "workend7": td.WorkingDays[6].WorkingHourEndTime}).Error; db != nil {
				return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
			}

			return &models.Response{ResponseCode: 200, Status: "OK", Msg: "Twitter integration Updated successfully."}, nil
		}
		return &models.Response{ResponseCode: 409, Status: "0", Msg: "Configuration name already exist."}, nil
	} else if td.TwitterIntegrationID != w.TwitterIntegrationID && td.AppId != w.AppId && td.ConfigurationName == w.ConfigurationName {
		fmt.Println("part5")
		w := models.TwitterConfiguration{}
		if err := r.DBConn.Where("twitter_integration_id = ?", td.TwitterIntegrationID).Or("app_id = ?", td.AppId).Find(&w).Error; err != nil {
			if db := r.DBConn.Table("twitter_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"app_id": td.AppId, "app_key": td.AppKey, "app_secret": td.AppSecret, "message": td.Message, "twitter_integration_id": td.TwitterIntegrationID, "size": td.Size, "trigger_name": td.Trigger.Name, "trigger_message": td.Trigger.Message, "trigger_when": td.Trigger.When, "day1": td.WorkingDays[0].Day, "day2": td.WorkingDays[1].Day, "day3": td.WorkingDays[2].Day, "day4": td.WorkingDays[3].Day, "day5": td.WorkingDays[4].Day, "day6": td.WorkingDays[5].Day, "day7": td.WorkingDays[6].Day, "workstart1": td.WorkingDays[0].WorkingHourStartTime, "workstart2": td.WorkingDays[1].WorkingHourStartTime, "workstart3": td.WorkingDays[2].WorkingHourStartTime, "workstart4": td.WorkingDays[3].WorkingHourStartTime, "workstart5": td.WorkingDays[4].WorkingHourStartTime, "workstart6": td.WorkingDays[5].WorkingHourStartTime, "workstart7": td.WorkingDays[6].WorkingHourStartTime, "workend1": td.WorkingDays[0].WorkingHourEndTime, "workend2": td.WorkingDays[1].WorkingHourEndTime, "workend3": td.WorkingDays[2].WorkingHourEndTime, "workend4": td.WorkingDays[3].WorkingHourEndTime, "workend5": td.WorkingDays[4].WorkingHourEndTime, "workend6": td.WorkingDays[5].WorkingHourEndTime, "workend7": td.WorkingDays[6].WorkingHourEndTime}).Error; db != nil {
				return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
			}

			return &models.Response{ResponseCode: 200, Status: "OK", Msg: "Twitter integration Updated successfully."}, nil
		}
		if w.TwitterIntegrationID == td.TwitterIntegrationID {
			return &models.Response{ResponseCode: 409, Status: "0", Msg: "Twitter id already exist."}, nil
		} else if w.AppId == td.AppId {
			return &models.Response{ResponseCode: 409, Status: "0", Msg: "AppId already exist."}, nil
		}
	} else if td.TwitterIntegrationID != w.TwitterIntegrationID && td.AppId == w.AppId && td.ConfigurationName == w.ConfigurationName {
		fmt.Println("part6")
		w := models.TwitterConfiguration{}
		if err := r.DBConn.Where("twitter_integration_id = ?", td.TwitterIntegrationID).Find(&w).Error; err != nil {
			if db := r.DBConn.Table("twitter_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"app_key": td.AppKey, "app_secret": td.AppSecret, "message": td.Message, "twitter_integration_id": td.TwitterIntegrationID, "size": td.Size, "trigger_name": td.Trigger.Name, "trigger_message": td.Trigger.Message, "trigger_when": td.Trigger.When, "day1": td.WorkingDays[0].Day, "day2": td.WorkingDays[1].Day, "day3": td.WorkingDays[2].Day, "day4": td.WorkingDays[3].Day, "day5": td.WorkingDays[4].Day, "day6": td.WorkingDays[5].Day, "day7": td.WorkingDays[6].Day, "workstart1": td.WorkingDays[0].WorkingHourStartTime, "workstart2": td.WorkingDays[1].WorkingHourStartTime, "workstart3": td.WorkingDays[2].WorkingHourStartTime, "workstart4": td.WorkingDays[3].WorkingHourStartTime, "workstart5": td.WorkingDays[4].WorkingHourStartTime, "workstart6": td.WorkingDays[5].WorkingHourStartTime, "workstart7": td.WorkingDays[6].WorkingHourStartTime, "workend1": td.WorkingDays[0].WorkingHourEndTime, "workend2": td.WorkingDays[1].WorkingHourEndTime, "workend3": td.WorkingDays[2].WorkingHourEndTime, "workend4": td.WorkingDays[3].WorkingHourEndTime, "workend5": td.WorkingDays[4].WorkingHourEndTime, "workend6": td.WorkingDays[5].WorkingHourEndTime, "workend7": td.WorkingDays[6].WorkingHourEndTime}).Error; db != nil {
				return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
			}

			return &models.Response{ResponseCode: 200, Status: "OK", Msg: "Twitter integration Updated successfully."}, nil
		}
		return &models.Response{ResponseCode: 409, Status: "0", Msg: "Twitter id already exist."}, nil
	} else if td.TwitterIntegrationID != w.TwitterIntegrationID && td.AppId == w.AppId && td.ConfigurationName != w.ConfigurationName {
		fmt.Println("part7")
		w := models.TwitterConfiguration{}
		if err := r.DBConn.Where("twitter_integration_id = ?", td.TwitterIntegrationID).Or("configuration_name = ?", td.ConfigurationName).Find(&w).Error; err != nil {
			if db := r.DBConn.Table("twitter_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"configuration_name": td.ConfigurationName, "app_key": td.AppKey, "app_secret": td.AppSecret, "message": td.Message, "twitter_integration_id": td.TwitterIntegrationID, "size": td.Size, "trigger_name": td.Trigger.Name, "trigger_message": td.Trigger.Message, "trigger_when": td.Trigger.When, "day1": td.WorkingDays[0].Day, "day2": td.WorkingDays[1].Day, "day3": td.WorkingDays[2].Day, "day4": td.WorkingDays[3].Day, "day5": td.WorkingDays[4].Day, "day6": td.WorkingDays[5].Day, "day7": td.WorkingDays[6].Day, "workstart1": td.WorkingDays[0].WorkingHourStartTime, "workstart2": td.WorkingDays[1].WorkingHourStartTime, "workstart3": td.WorkingDays[2].WorkingHourStartTime, "workstart4": td.WorkingDays[3].WorkingHourStartTime, "workstart5": td.WorkingDays[4].WorkingHourStartTime, "workstart6": td.WorkingDays[5].WorkingHourStartTime, "workstart7": td.WorkingDays[6].WorkingHourStartTime, "workend1": td.WorkingDays[0].WorkingHourEndTime, "workend2": td.WorkingDays[1].WorkingHourEndTime, "workend3": td.WorkingDays[2].WorkingHourEndTime, "workend4": td.WorkingDays[3].WorkingHourEndTime, "workend5": td.WorkingDays[4].WorkingHourEndTime, "workend6": td.WorkingDays[5].WorkingHourEndTime, "workend7": td.WorkingDays[6].WorkingHourEndTime}).Error; db != nil {
				return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
			}

			return &models.Response{ResponseCode: 200, Status: "OK", Msg: "Twitter integration Updated successfully."}, nil
		}
		if w.TwitterIntegrationID == td.TwitterIntegrationID {
			return &models.Response{ResponseCode: 409, Status: "0", Msg: "Twitter id already exist."}, nil
		} else if w.ConfigurationName == td.ConfigurationName {
			return &models.Response{ResponseCode: 409, Status: "0", Msg: "Configuration name already exist."}, nil
		}
	} else if td.TwitterIntegrationID == w.TwitterIntegrationID && td.AppId != w.AppId && td.ConfigurationName == w.ConfigurationName {
		fmt.Println("part8")
		w := models.TwitterConfiguration{}
		if err := r.DBConn.Where("app_id = ?", td.AppId).Find(&w).Error; err != nil {
			if db := r.DBConn.Table("twitter_configurations").Where("domain_uuid = ? AND id = ?", domain_uuid, id).Updates(map[string]interface{}{"app_id": td.AppId, "app_key": td.AppKey, "app_secret": td.AppSecret, "message": td.Message, "size": td.Size, "trigger_name": td.Trigger.Name, "trigger_message": td.Trigger.Message, "trigger_when": td.Trigger.When, "day1": td.WorkingDays[0].Day, "day2": td.WorkingDays[1].Day, "day3": td.WorkingDays[2].Day, "day4": td.WorkingDays[3].Day, "day5": td.WorkingDays[4].Day, "day6": td.WorkingDays[5].Day, "day7": td.WorkingDays[6].Day, "workstart1": td.WorkingDays[0].WorkingHourStartTime, "workstart2": td.WorkingDays[1].WorkingHourStartTime, "workstart3": td.WorkingDays[2].WorkingHourStartTime, "workstart4": td.WorkingDays[3].WorkingHourStartTime, "workstart5": td.WorkingDays[4].WorkingHourStartTime, "workstart6": td.WorkingDays[5].WorkingHourStartTime, "workstart7": td.WorkingDays[6].WorkingHourStartTime, "workend1": td.WorkingDays[0].WorkingHourEndTime, "workend2": td.WorkingDays[1].WorkingHourEndTime, "workend3": td.WorkingDays[2].WorkingHourEndTime, "workend4": td.WorkingDays[3].WorkingHourEndTime, "workend5": td.WorkingDays[4].WorkingHourEndTime, "workend6": td.WorkingDays[5].WorkingHourEndTime, "workend7": td.WorkingDays[6].WorkingHourEndTime}).Error; db != nil {
				return &models.Response{Status: "0", Msg: "Oops! There is some problem! Try again.", ResponseCode: http.StatusBadRequest}, nil
			}

			return &models.Response{ResponseCode: 200, Status: "OK", Msg: "Twitter integration Updated successfully."}, nil
		}
		return &models.Response{ResponseCode: 409, Status: "0", Msg: "AppId already exist."}, nil
	} else {
		return &models.Response{ResponseCode: 401, Status: "Not Ok", Msg: "Twitter integration Not Updated."}, nil
	}
	return &models.Response{ResponseCode: 204, Status: "OK", Msg: "Twitter integration Updated successfully."}, nil

}

/**********************************************Delete Tenant details*************************************/
func (r *crudRepository) Delete_Twitter_configuration(ctx context.Context, id int64, domain_uuid string) (*models.Response, error) {
	td := models.TwitterConfiguration{}
	db := r.DBConn.Where("domain_uuid = ? AND id = ?", domain_uuid, id).Delete(&td)
	if db.Error != nil {
		return &models.Response{Status: "0", Msg: "Table not deleted", ResponseCode: 404}, nil
	}
	if db.RowsAffected == 0 {
		return &models.Response{Status: "Not Found", Msg: "Record Doesn't Exist", ResponseCode: 401}, nil
	}
	return &models.Response{Status: "1", Msg: "Twitter Configuration deleted.", ResponseCode: 200}, nil
}

/**********************************************List integration**************************************/
func (r *crudRepository) List_integration(ctx context.Context, appId string) ([]byte, error) {
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
func (r *crudRepository) DeleteAllMessage(ctx context.Context, appUserId string, appId string) (*models.Response, error) {
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
func (r *crudRepository) DeleteMessage(ctx context.Context, appId string, appUserId string, messageId string) (*models.Response, error) {
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
func (r *crudRepository) Create_Text_Template(ctx context.Context, appId string, p models.Comtemplate) ([]byte, error) {
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
func (r *crudRepository) Create_Carousel_Template(ctx context.Context, appId string, p models.Payload) ([]byte, error) {
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
func (r *crudRepository) Create_Compound_Template(ctx context.Context, appId string, p models.Comtemplate) ([]byte, error) {
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
func (r *crudRepository) PostMessage(ctx context.Context, appId string, ConversationId string, p models.User) ([]byte, error) {
	td := models.Tenant_details{}
	db := r.DBConn.Where("app_id = ?", appId).Find(&td)
	if db.Error != nil {
		return nil, db.Error
	}

	jsonValue, _ := json.Marshal(p)
	req, _ := http.NewRequest("POST", "https://api.smooch.io/v2/apps/"+appId+"/conversations/"+ConversationId+"/messages", bytes.NewBuffer(jsonValue))
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
func (r *crudRepository) Create_Request_Location_Template(ctx context.Context, appId string, p models.Comtemplate) ([]byte, error) {
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
func (r *crudRepository) Update_Text_Template(ctx context.Context, appId string, template_id string, p models.Payload) ([]byte, error) {
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
func (r *crudRepository) Get_template(ctx context.Context, appId string, template_id string) ([]byte, error) {
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
func (r *crudRepository) List_template(ctx context.Context, appId string) ([]byte, error) {
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
func (r *crudRepository) Delete_template(ctx context.Context, appId string, template_id string) (*models.Response, error) {
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
func (r *crudRepository) Send_Location(ctx context.Context, appId string, appUserId string, p models.Locations) ([]byte, error) {
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
func (r *crudRepository) Message_Action_Types(ctx context.Context, appId string, appUserId string, p models.User) ([]byte, error) {
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
func (r *crudRepository) Quickreply_Message(ctx context.Context, appId string, appUserId string, p models.User) ([]byte, error) {
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
func (r *crudRepository) Send_Carousel_Message(ctx context.Context, appId string, appUserId string, p models.User) ([]byte, error) {

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
func (r *crudRepository) Link_appUser_to_Channel(ctx context.Context, appId string, appUserId string, p models.Link) ([]byte, error) {

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
func (r *crudRepository) Unlink_appUser_to_Channel(ctx context.Context, appId string, appUserId string, channel string) ([]byte, error) {

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
func (r *crudRepository) Upload_Attachments(ctx context.Context, displayName string, AvatarURL string, appId string, conversationId string, Type string, Text string, IntegrationID string, Size int64, file multipart.File, handler *multipart.FileHeader) (*models.Response, error) {
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
			req, _ := http.NewRequest("POST", "https://api.smooch.io/v2/apps/"+appId+"/attachments?access=public&for=message&conversationId="+conversationId, body)
			req.Header.Add("Content-Type", "multipart/form-data")
			req.Header.Set("Content-Type", writer.FormDataContentType())
			req.SetBasicAuth(fb.AppKey, fb.AppSecret)
			client := &http.Client{}
			res, err := client.Do(req)
			if err != nil {
				return nil, err
			} else {
				u := models.Sender{}
				data, _ := ioutil.ReadAll(res.Body)
				jsonData := json.Unmarshal(data, &u)
				fmt.Println(jsonData, u.Attachment.MediaUrl, u.Attachment.MediaType, "ghvghvqv")
				p := models.User{
					Author: models.Author{
						Type:        "business",
						DisplayName: displayName,
						AvatarURL:   AvatarURL,
					},
					Content: models.Content{
						Type:     Type,
						Text:     Text,
						MediaUrl: u.Attachment.MediaUrl,
					},
				}
				r.PostMessage(ctx, appId, conversationId, p)
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
		req, _ := http.NewRequest("POST", "https://api.smooch.io/v2/apps/"+appId+"/attachments?access=public&for=message&conversationId="+conversationId, body)
		req.Header.Add("Content-Type", "multipart/form-data")
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.SetBasicAuth(td.AppKey, td.AppSecret)
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		} else {
			u := models.Sender{}
			data, _ := ioutil.ReadAll(res.Body)
			jsonData := json.Unmarshal(data, &u)
			fmt.Println(jsonData, u.Attachment.MediaUrl, u.Attachment.MediaType, "ghvghvqv")
			p := models.User{
				Author: models.Author{
					Type:        "business",
					DisplayName: displayName,
					AvatarURL:   AvatarURL,
				},
				Content: models.Content{
					Type:     Type,
					Text:     Text,
					MediaUrl: u.Attachment.MediaUrl,
				},
			}
			r.PostMessage(ctx, appId, conversationId, p)
			return &models.Response{Status: "1", Msg: "File is sent successfully.", ResponseCode: 200}, nil
		}

	} else {
		return &models.Response{Status: "0", Msg: "Please choose message type.", ResponseCode: 200}, nil
	}
}

/***********************************************TypingActivity***********************************************/
func (r *crudRepository) TypingActivity(ctx context.Context, appId string, appUserId string, p models.User) ([]byte, error) {
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
func (r *crudRepository) Disable_AppUser(ctx context.Context, appUserId string) (*models.Response, error) {
	u := models.ReceiveUserDetails{}
	customer := models.Customer_Agents{
		AppUserId: appUserId,
	}
	err := r.DBConn.Where("app_user_id = ?", appUserId).Find(&u)
	if err.Error != nil {
		return &models.Response{Status: "0", Msg: "AppUserId Not Found.", ResponseCode: 404}, nil

	}
	db := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", appUserId).Delete(&u)
	if db.Error != nil {
		return &models.Response{Status: "0", Msg: "Customer not Disabled.", ResponseCode: 404}, nil
	}
	cust := r.DBConn.Where("app_user_id = ?", appUserId).Delete(&customer)
	if cust.Error != nil {
		return &models.Response{Status: "0", Msg: "Customer not Removed from queue.", ResponseCode: 404}, nil
	}
	return &models.Response{Status: "1", Msg: "Customer Disabled Successfully.", ResponseCode: 200}, nil
}

/****************************************Reset Unread Count*******************************************/
func (r *crudRepository) Reset_Unread_Count(ctx context.Context, appId string, appUserId string) (*models.Response, error) {
	u := models.ReceiveUserDetails{}
	td := models.Tenant_details{
		AppId: appId,
	}
	db := r.DBConn.Table("tenant_details").Where("app_id = ?", appId).Find(&td)
	if db.Error != nil {
		return nil, db.Error
	}
	req, _ := http.NewRequest("POST", "https://api.smooch.io/v1.1/apps/"+appId+"/appusers/"+appUserId+"/conversation/read", nil)
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(td.AppKey, td.AppSecret)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	} else {
		data, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(data))
		err := r.DBConn.Where("app_user_id = ?", appUserId).Find(&u)
		if err != nil {
			return &models.Response{Status: "0", Msg: "AppUserId Not Found.", ResponseCode: 404}, nil

		}
		db := r.DBConn.Table("receive_user_details").Where("app_user_id = ?", appUserId).Update("unread_count", 0)
		if db.Error != nil {
			return &models.Response{Status: "0", Msg: "AppUserId unread count not Updated.", ResponseCode: 404}, nil
		}
		return &models.Response{Status: "1", Msg: "Unread count reset Successfully.", ResponseCode: 200}, nil
	}
}

/************************************************Create Queue***********************************************/
func (r *crudRepository) Create_Queue(ctx context.Context, Id int64, Queue_uuid string, Map_with string, Name string, IntegrationID string, Domain_uuid string) (*models.Response, error) {
	uuid1, _ := myNewUUID.NewUUID()
	uuid := uuid1.String()
	u := models.Queue{
		Id:            Id,
		Name:          Name,
		IntegrationID: IntegrationID,
		Queue_uuid:    uuid,
		Map_with:      Map_with,
		Domain_uuid:   Domain_uuid,
	}
	if err := r.DBConn.Where("name = ?", Name).Find(&u).Error; err != nil {
		Queue := r.DBConn.Create(&u)
		if Queue.RowsAffected == 0 {
			return &models.Response{Status: "0", Msg: "Queue not created.", ResponseCode: 400}, nil
		} else {
			return &models.Response{Status: "1", Msg: "Queue created successfully.", ResponseCode: 200}, nil
		}

	}
	return &models.Response{Status: "0", Msg: "Queue name Already exists.", ResponseCode: 404}, nil
}

/***************************************************Assign_Agent********************************************/
func (r *crudRepository) Assign_Agent_To_Queue(ctx context.Context, Agent_name string, Agent_uuid string, Queue_name string, tenant_domain_uuid string, Queue_uuid string) (*models.Response, error) {
	u := models.AgentQueue{
		AgentName:          Agent_name,
		Agent_uuid:         Agent_uuid,
		QueueName:          Queue_name,
		Tenant_domain_uuid: tenant_domain_uuid,
		Queue_uuid:         Queue_uuid,
	}
	if err := r.DBConn.Where("agent_uuid = ? AND queue_name = ?", Agent_uuid, Queue_name).Find(&u).Error; err != nil {
		Queue := r.DBConn.Create(&u)
		if Queue.RowsAffected == 0 {
			return &models.Response{Status: "0", Msg: "Agent not Assigned.", ResponseCode: 400}, nil
		} else {
			return &models.Response{Status: "1", Msg: "Agent Assigned successfully.", ResponseCode: 200}, nil
		}

	}
	return &models.Response{Status: "0", Msg: "Agent Already In Queue.", ResponseCode: 404}, nil

}

/*****************************************************Remove Agent From Queue*********************************/
func (r *crudRepository) Remove_Agent_From_Queue(ctx context.Context, Agent_uuid string) (*models.Response, error) {
	u := models.AgentQueue{}
	err := r.DBConn.Where("agent_uuid = ?", Agent_uuid).Find(&u)
	if err.Error != nil {
		return &models.Response{Status: "0", Msg: "Agent not found in Queue.", ResponseCode: 404}, nil
	}
	db := r.DBConn.Where("agent_uuid = ?", Agent_uuid).Delete(&u)
	if db.RowsAffected == 0 {
		return &models.Response{Status: "0", Msg: "Agent not removed Queue..", ResponseCode: 404}, nil
	}

	return &models.Response{Status: "1", Msg: "Agent Removed From Queue.", ResponseCode: 200}, nil

}

/**************************************************Get Assigned Agent list From Queue***************************/
func (r *crudRepository) Get_Assigned_Agent_list_From_Queue(ctx context.Context, queue_uuid string) (*models.Response, error) {
	list := make([]models.AgentQueue, 0)
	a := models.AgentQueue{}
	err := r.DBConn.Where("queue_uuid = ?", queue_uuid).Find(&a)
	if err.Error != nil {
		return &models.Response{Status: "Not Found", Msg: "Queue Not Found", ResponseCode: 404}, nil
	}
	if rows, err := r.DBConn.Raw("select queue_name, agent_name, agent_uuid from agent_queues where queue_uuid = ?", queue_uuid).Rows(); err != nil {

		return &models.Response{Status: "Not Found", Msg: "Record Not Found", ResponseCode: 204}, nil
	} else {
		defer rows.Close()
		for rows.Next() {
			f := models.AgentQueue{}
			if err := rows.Scan(&f.QueueName, &f.AgentName, &f.Agent_uuid); err != nil {

				return nil, err
			}

			list = append(list, f)
		}

		return &models.Response{Status: "OK", Msg: "Queue Found", ResponseCode: 200, AssignAgent: list}, nil
	}
}

/*********************************************Get Queue List****************************************************/
func (r *crudRepository) Get_Queue_List(ctx context.Context, domain_uuid string) (*models.Response, error) {
	list := make([]models.Queue, 0)
	a := models.Queue{}
	err := r.DBConn.Where("domain_uuid = ?", domain_uuid).Find(&a)
	if err.Error != nil {
		return &models.Response{Status: "Not Found", Msg: "Queue Not Found", ResponseCode: 404}, nil
	}
	if rows, err := r.DBConn.Raw("select id, name, integration_id, queue_uuid, map_with, domain_uuid from queues where domain_uuid = ?", domain_uuid).Rows(); err != nil {

		return &models.Response{Status: "Not Found", Msg: "Record Not Found", ResponseCode: 404}, nil
	} else {
		defer rows.Close()
		for rows.Next() {
			f := models.Queue{}
			if err := rows.Scan(&f.Id, &f.Name, &f.IntegrationID, &f.Queue_uuid, &f.Map_with, &f.Domain_uuid); err != nil {

				return nil, err
			}

			list = append(list, f)
		}

		return &models.Response{Status: "OK", Msg: "Record Found", ResponseCode: 200, Queue: list}, nil
	}
}

/*********************************************Update Queue*************************************************/
func (r *crudRepository) Update_Queue(ctx context.Context, queue_uuid string, Name string, IntegrationID string, Map_with string, Domain_uuid string) (*models.Response, error) {
	u := models.Queue{}
	err := r.DBConn.Where("queue_uuid = ?", queue_uuid).Find(&u)
	if err.Error != nil {
		return &models.Response{Status: "0", Msg: "Queue not found.", ResponseCode: 404}, nil
	}
	Queue := r.DBConn.Table("queues").Where("queue_uuid = ?", queue_uuid).Updates(map[string]interface{}{"map_with": Map_with, "name": Name, "integration_id": IntegrationID, "domain_uuid": Domain_uuid})
	if Queue.Error != nil {
		return &models.Response{Status: "0", Msg: "Queue not Updated.", ResponseCode: 400}, nil
	}
	return &models.Response{Status: "1", Msg: "Queue Updated successfully.", ResponseCode: 200}, nil

}

/************************************************Delete Queue*******************************************/
func (r *crudRepository) Delete_Queue(ctx context.Context, queue_uuid string) (*models.Response, error) {
	u := models.Queue{}
	a := models.AgentQueue{}
	err := r.DBConn.Where("queue_uuid = ?", queue_uuid).Find(&u)
	if err.Error != nil {
		return &models.Response{Status: "0", Msg: "Queue not found.", ResponseCode: 404}, nil
	}
	assignqueue := r.DBConn.Where("queue_uuid = ?", queue_uuid).Find(&a)
	if assignqueue.Error != nil {
		// return &models.Response{Status: "0", Msg: "Queue not found.", ResponseCode: 404}, nil
	}
	db := r.DBConn.Where("queue_uuid = ?", queue_uuid).Delete(&u)
	if db.RowsAffected == 0 {
		return &models.Response{Status: "0", Msg: "Queue not Deleted.", ResponseCode: 404}, nil
	}
	que := r.DBConn.Where("queue_uuid = ?", queue_uuid).Delete(&a)
	if que.RowsAffected == 0 {
		//return &models.Response{Status: "0", Msg: "Queue not Deleted.", ResponseCode: 404}, nil
	}
	return &models.Response{Status: "1", Msg: "Queue Delete Successfully.", ResponseCode: 200}, nil

}

/************************************************Avaiable Agents*********************************************/
func (r *crudRepository) Available_Agents(ctx context.Context, domain_uuid string, queue_uuid string) (*models.Response, error) {
	list := make([]models.V_call_center_agents, 0)
	a := models.Queue{}
	b := models.V_call_center_agents{}
	err := r.DBConn.Where("queue_uuid = ?", queue_uuid).Find(&a)
	if err.Error != nil {
		return &models.Response{Status: "Not Found", Msg: "Queue Not Found", ResponseCode: 404}, nil
	}
	db := r.DBConn.Where("domain_uuid = ?", domain_uuid).Find(&b)
	if db.Error != nil {
		return &models.Response{Status: "Not Found", Msg: "Agent Not Found", ResponseCode: 404}, nil
	}
	if rows, err := r.DBConn.Raw("select call_center_agent_uuid, agent_name from v_call_center_agents where domain_uuid = ? EXCEPT select agent_queues.agent_uuid, agent_queues.agent_name from agent_queues where queue_uuid = ?", domain_uuid, queue_uuid).Rows(); err != nil {

		return &models.Response{Status: "Not Found", Msg: "Record Not Found", ResponseCode: 404}, nil
	} else {
		defer rows.Close()
		for rows.Next() {
			f := models.V_call_center_agents{}
			if err := rows.Scan(&f.CallCenterAgentUUID, &f.AgentName); err != nil {
				return nil, err
			}
			list = append(list, f)
		}
		if len(list) == 0 {
			return &models.Response{Status: "0", Msg: "Agents Not Found", ResponseCode: 404}, nil
		}
		return &models.Response{Status: "OK", Msg: "Record Found", ResponseCode: 200, Agent: list}, nil
	}
}

/*************************************************Available Agents in Queue********************************/
func (r *crudRepository) Get_Available_Agents_Queue_List(ctx context.Context, agent_uuid string, queue_uuid string) (*models.Response, error) {
	list := make([]models.AgentQueue, 0)
	a := models.AgentQueue{}
	err := r.DBConn.Where("queue_uuid = ?", queue_uuid).Find(&a)
	if err.Error != nil {
		return &models.Response{Status: "Not Found", Msg: "Agents Not Found", ResponseCode: 404}, nil
	}
	if rows, err := r.DBConn.Raw("select agent_name,agent_uuid from agent_queues where queue_uuid = ? And agent_uuid != ?", queue_uuid, agent_uuid).Rows(); err != nil {

		return &models.Response{Status: "Not Found", Msg: "Record Not Found", ResponseCode: 404}, nil
	} else {
		defer rows.Close()
		for rows.Next() {
			f := models.AgentQueue{}
			if err := rows.Scan(&f.AgentName, &f.Agent_uuid); err != nil {

				return nil, err
			}

			list = append(list, f)
		}

		return &models.Response{Status: "OK", Msg: "Record Found", ResponseCode: 200, AssignAgent: list}, nil
	}
}

/************************************************Transfer customer**********************************************/
func (r *crudRepository) Transfer_customer(ctx context.Context, agent_name string, conversation_id string, agent_uuid string, appUserId string) (*models.Response, error) {
	v_call_agents := models.V_call_center_agents{}
	//app := models.Customer_Agents{
	//	AppUserId: appUserId,
	//}
	td := models.Tenant_details{}
	agent_uuid1, _ := uuid.Parse(agent_uuid)

	agent := r.DBConn.Table("v_call_center_agents").Where("call_center_agent_uuid = ?", agent_uuid1).Find(&v_call_agents)
	if agent.Error != nil {
		return &models.Response{Status: "0", Msg: "Agent not Found.", ResponseCode: 404}, nil
	}
	smooch := r.DBConn.Table("tenant_details").Where("domain_uuid = ?", v_call_agents.Domain_uuid).Find(&td)
	if smooch.Error != nil {
		return &models.Response{Status: "0", Msg: "Domain not Found.", ResponseCode: 404}, nil
	}
	fmt.Println(td.AppId, "app_id", conversation_id)
	if v_call_agents.AgentStatus == "Available" {

		// err := r.DBConn.Where("app_user_id = ?", appUserId).Find(&app)
		// if err.Error != nil {
		// 	return &models.Response{Status: "0", Msg: "Customer not assigned to agent.", ResponseCode: 404}, nil
		// }
		db := r.DBConn.Table("customer_agents").Where("app_user_id = ?", appUserId).Update("agent_uuid", agent_uuid)
		if db.Error != nil {
			return &models.Response{Status: "0", Msg: "Customer not assigned to agent.", ResponseCode: 404}, nil
		}
		p := models.User{
			Author: models.Author{
				Type:        "business",
				DisplayName: "system",
				AvatarURL:   "https://www.gravatar.com/image.jpg",
			},
			Content: models.Content{
				Type: "text",
				Text: agent_name + " is transferring your chat to " + v_call_agents.AgentName,
			},
		}
		r.PostMessage(ctx, td.AppId, conversation_id, p)
		msg := map[string]interface{}{"message_id": "5", "customer_id": appUserId, "user_id": agent_uuid, "user_type": "agent"}

		for _, oldu := range r.SList.Users {
			if oldu.UName == agent_uuid {
				log.Println("found user: ", oldu)
				if err := websocket.JSON.Send(oldu.Ws, msg); err != nil {
					log.Println("Can't send", err)
				}
				fmt.Println("sucessfully sent ", oldu, appUserId)

			}
		}

		return &models.Response{Status: "1", Msg: "Customer assigned to agent successfully.", ResponseCode: 200}, nil
	} else {
		return &models.Response{Status: "0", Msg: "Agent is not Available.", ResponseCode: 409}, nil
	}

}

/*****************************************************Post page on Fb*****************************************/
func (r *crudRepository) Publish_Post_on_FB_Page(ctx context.Context, pageId string, message string, access_token string) ([]byte, error) {
	fmt.Println(pageId, access_token, message)
	message = strings.ReplaceAll(message, " ", "%20")
	res, err := http.NewRequest("POST", "https://graph.facebook.com/"+pageId+"/feed?message="+message+"&access_token="+access_token, nil)
	res.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(res)
	if err != nil {
		fmt.Printf("error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data), "enterrer")
		return data, nil
	}

	defer res.Body.Close()
	return nil, err
}

/**************************************Get all Post of a Page*****************************************/
func (r *crudRepository) Getall_Post_of_Page(ctx context.Context, pageId string, access_token string) ([]byte, error) {
	res, err := http.NewRequest("GET", "https://graph.facebook.com/"+pageId+"?fields=id,name,feed{created_time,message,attachments}&access_token="+access_token, nil)
	res.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(res)
	if err != nil {
		fmt.Printf("error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data), "enterrer")
		return data, nil
	}

	defer res.Body.Close()
	return nil, err
}

/*********************************************Delete Post Of a Page****************************************/
func (r *crudRepository) Delete_Post_of_Page(ctx context.Context, page_postId string, access_token string) ([]byte, error) {
	res, err := http.NewRequest("DELETE", "https://graph.facebook.com/"+page_postId+"?access_token="+access_token, nil)
	res.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(res)
	if err != nil {
		fmt.Printf("error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data), "enterrer")
		return data, nil
	}
	defer res.Body.Close()
	return nil, err
}

/**********************************************Update Post of Page******************************************/
func (r *crudRepository) Update_Post_of_Page(ctx context.Context, page_postId string, message string, access_token string) ([]byte, error) {
	message = strings.ReplaceAll(message, " ", "%20")
	res, err := http.NewRequest("POST", "https://graph.facebook.com/"+page_postId+"?message="+message+"&access_token="+access_token, nil)
	res.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(res)
	if err != nil {
		fmt.Printf("error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data), "enterrer")
		return data, nil
	}
	defer res.Body.Close()
	return nil, err
}

/********************************************Get Comments on Page********************************************/
func (r *crudRepository) Get_Comments_on_Post_of_Page(ctx context.Context, page_postId string, access_token string) ([]byte, error) {
	res, err := http.NewRequest("GET", "https://graph.facebook.com/"+page_postId+"/comments?limit=100&summary=total_count&access_token="+access_token, nil)
	res.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(res)
	if err != nil {
		fmt.Printf("error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data), "enterrer")
		return data, nil
	}
	defer res.Body.Close()
	return nil, err
}

/************************************************Get Likes of a page*******************************************/
func (r *crudRepository) Get_Likes_on_Post_of_Page(ctx context.Context, page_postId string, access_token string) ([]byte, error) {
	res, err := http.NewRequest("GET", "https://graph.facebook.com/"+page_postId+"/likes?fields=name,pic&summary=total_count&limit=100&access_token="+access_token, nil)
	res.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(res)
	if err != nil {
		fmt.Printf("error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data), "enterrer")
		return data, nil
	}
	defer res.Body.Close()
	return nil, err
}

/********************************************Comment on Post of Page******************************************/
func (r *crudRepository) Comment_on_Post_of_Page(ctx context.Context, page_postId string, message string, access_token string) ([]byte, error) {
	message = strings.ReplaceAll(message, " ", "%20")
	res, err := http.NewRequest("POST", "https://graph.facebook.com/"+page_postId+"/comments?message="+message+"&access_token="+access_token, nil)
	res.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(res)
	if err != nil {
		fmt.Printf("error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data), "enterrer")
		return data, nil
	}
	defer res.Body.Close()
	return nil, err
}

/************************************************Get Page Id *************************************************/
func (r *crudRepository) Get_Page_ID(ctx context.Context, access_token string) ([]byte, error) {
	res, err := http.NewRequest("GET", "https://graph.facebook.com/me/accounts?fields=redirect,access_token,picture,name&access_token="+access_token, nil)
	res.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(res)
	if err != nil {
		fmt.Printf("error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data), "enterrer")
		return data, nil
	}
	defer res.Body.Close()
	return nil, err
}

/************************************************Schedule Post************************************************/
func (r *crudRepository) Schedule_Post(ctx context.Context, pageId string, message string, scheduled_publish_time string, access_token string) ([]byte, error) {
	message = strings.ReplaceAll(message, " ", "%20")
	scheduled_publish_time = strings.ReplaceAll(scheduled_publish_time, ":", "%3A")
	res, err := http.NewRequest("POST", "https://graph.facebook.com/"+pageId+"/feed?published=false&message="+message+"&scheduled_publish_time="+scheduled_publish_time+"&access_token="+access_token, nil)
	res.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(res)
	if err != nil {
		fmt.Printf("error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data), "enterrer")
		return data, nil
	}
	defer res.Body.Close()
	return nil, err
}

/*******************************************Publish link with message on Post***************************************/
func (r *crudRepository) Publish_link_with_message_on_Post(ctx context.Context, pageId string, message string, link string, access_token string) ([]byte, error) {
	message = strings.ReplaceAll(message, " ", "%20")
	res, err := http.NewRequest("POST", "https://graph.facebook.com/"+pageId+"/feed?message="+message+"&link="+link+"&access_token="+access_token, nil)
	res.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(res)
	if err != nil {
		fmt.Printf("error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data), "enterrer")
		return data, nil
	}
	defer res.Body.Close()
	return nil, err
}

/******************************************Upload Photo on Post**************************************************/
func (r *crudRepository) Upload_Photo_on_Post(ctx context.Context, pageId string, access_token string, message string, Type string, file multipart.File, handler *multipart.FileHeader) ([]byte, error) {
	message = strings.ReplaceAll(message, " ", "%20")
	fmt.Println(Type, "type")
	if Type == "image" {
		fmt.Println("image")
		IMAGE_DIR := "/home/" + getServerOs() + "/Downloads/temp_images/"
		dir_location := IMAGE_DIR
		getFileName := handler.Filename

		fb_image_path := dir_location + getFileName

		if err := os.MkdirAll(dir_location, os.FileMode(0777)); err != nil {
			fmt.Println(err)
		}
		f, err := os.OpenFile(fb_image_path, os.O_WRONLY|os.O_CREATE, os.FileMode(0777))
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()
		io.Copy(f, file)

		res, err := http.NewRequest("POST", "https://graph.facebook.com/"+pageId+"/photos?url=http://"+getserverhost()+fb_image_path+"&message="+message+"&access_token="+access_token, nil)
		res.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		response, err := client.Do(res)
		if err != nil {
			fmt.Printf("error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			fmt.Println(string(data))
			err := os.Remove(fb_image_path)
			if err != nil {
				fmt.Println("errror", err)
			}
			return data, nil
		}
		return nil, err
	} else if Type == "video" {
		fmt.Println("video")
		VIDEO_DIR := "/home/" + getServerOs() + "/Downloads/temp_images/"
		dir_location := VIDEO_DIR
		getFileName := handler.Filename

		fb_video_path := dir_location + getFileName

		if err := os.MkdirAll(dir_location, os.FileMode(0777)); err != nil {
			fmt.Println(err)
		}
		f, err := os.OpenFile(fb_video_path, os.O_WRONLY|os.O_CREATE, os.FileMode(0777))
		if err != nil {
			fmt.Println(err)
		}

		defer f.Close()
		io.Copy(f, file)

		res, err := http.NewRequest("POST", "https://graph.facebook.com/"+pageId+"/videos?file_url=http://"+getserverhost()+fb_video_path+"&message="+message+"&access_token="+access_token, nil)
		//res.Header.Set("Content-Type", "application/json")
		res.Header.Add("Content-Type", "multipart/form-data")

		client := &http.Client{}
		response, err := client.Do(res)
		if err != nil {
			fmt.Printf("error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			fmt.Println(string(data))
			err := os.Remove(fb_video_path)
			if err != nil {
				fmt.Println("error", err)
			}
			return data, nil
		}
		return nil, err
	}
	return nil, nil
}

/**************************************************Facebook login Api********************************************/
func (r *crudRepository) UVoiceFacebookLogin(ctx context.Context, c echo.Context, client_id string, client_secret string, flac_uuid string) (*models.Response, error) {
	fmt.Println(c.Request)
	fmt.Println(c.Response)
	fmt.Println(client_id, client_secret)
	oauthConf := &oauth2.Config{
		ClientID:     client_id,
		ClientSecret: client_secret,
		RedirectURL:  HTTPSECURE + getserverhost() + ":" + PORT + "/uvoice-facebook-login-callback",
		Scopes:       []string{"public_profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.facebook.com/v8.0/dialog/oauth",
			TokenURL: facebook.Endpoint.TokenURL,
		},
	}
	oauthStateString := flac_uuid
	Url, err := url.Parse(oauthConf.Endpoint.AuthURL)
	if err != nil {
		log.Fatal("Parse: ", err)
	}
	parameters := url.Values{}
	parameters.Add("client_id", oauthConf.ClientID)
	parameters.Add("scope", strings.Join(oauthConf.Scopes, " "))
	parameters.Add("redirect_uri", oauthConf.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", oauthStateString)
	Url.RawQuery = parameters.Encode()
	url := Url.String()
	fmt.Println(url)
	// b := fmt.Sprintf("%#v", oauthConf)
	// newrq := c.Request()
	// newrq.Body = ioutil.NopCloser(bytes.NewBufferString(b))
	// c.SetRequest(newrq)
	c.Redirect(http.StatusMovedPermanently, url)
	return nil, nil
}

/******************************************Facebook login Api callback****************************************/
func (r *crudRepository) UVoiceFacebookLoginCallback(ctx context.Context, c echo.Context) (*models.Response, error) {
	code := c.FormValue("code")
	state := c.FormValue("state")
	var t models.FacebookLoginAppConfiguration
	if err := r.DBConn.Where("flac_uuid=?", state).Find(&t).Error; err != nil {
		return &models.Response{Status: "Error", Msg: "Failed", ResponseCode: http.StatusBadRequest}, nil
	}
	if state != t.FlacUUID {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", t.FlacUUID, state)
		// http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return &models.Response{Status: "Error", Msg: "Failed", ResponseCode: http.StatusBadRequest}, nil
	}
	if code == "" {
		return &models.Response{Status: "Error", Msg: "Failed", ResponseCode: http.StatusBadRequest}, nil
	}
	oauthConf := &oauth2.Config{
		ClientID:     t.AppId,
		ClientSecret: t.AppSecret,
		RedirectURL:  HTTPSECURE + getserverhost() + ":" + PORT + "/uvoice-facebook-login-callback",
		Scopes:       []string{"public_profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.facebook.com/v8.0/dialog/oauth",
			TokenURL: facebook.Endpoint.TokenURL,
		},
	}
	token, err := oauthConf.Exchange(ctx, code)
	if err != nil {
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		return &models.Response{Status: "Error", Msg: "Failed", ResponseCode: http.StatusBadRequest}, nil
	}
	fmt.Printf("%v \n", token)
	c.Response().Header().Set("access_token", token.AccessToken)
	c.SetCookie(&http.Cookie{Name: "uvoice_facebook_access_token", Value: token.AccessToken})
	c.Redirect(http.StatusTemporaryRedirect, HTTPSECURE+getserverhost()+":"+PORT+"/uvoice-facebook-login-status")
	// return &models.Response{Status: "OK", Msg: "Success1", ResponseCode: http.StatusOK, FacebookGetAuthInfo: &info}, nil
	return nil, nil
}

/************************************Add facebook Application**************************************************/
func (r *crudRepository) AddFacebookApplication(ctx context.Context, domain_uuid string, app_id string, app_secret string, app_name string) (*models.Response, error) {
	if err := r.DBConn.Create(&models.FacebookLoginAppConfiguration{
		DomainUUID: domain_uuid,
		AppId:      app_id,
		AppSecret:  app_secret,
		AppName:    app_name,
	}).Error; err != nil {
		return &models.Response{Status: "0", Msg: "Failed", ResponseCode: http.StatusBadRequest}, nil
	}
	return &models.Response{Status: "1", Msg: "Success", ResponseCode: http.StatusOK}, nil
}

/*******************************************************Show Facebook Application*******************************/
func (r *crudRepository) ShowFacebookApplication(ctx context.Context, domain_uuid string) (*models.Response, error) {
	t := []models.FacebookLoginAppConfiguration{}
	if err := r.DBConn.Model(&models.FacebookLoginAppConfiguration{}).Where("domain_uuid=?", domain_uuid).Find(&t).Error; err != nil {
		return &models.Response{Status: "0", Msg: "Failed", ResponseCode: http.StatusBadRequest}, nil
	}
	if len(t) == 0 {
		return &models.Response{Status: "0", Msg: "Empty", ResponseCode: http.StatusOK}, nil
	}
	return &models.Response{Status: "1", Msg: "Success", ResponseCode: http.StatusOK, FacebookLoginAppConfiguration: &t}, nil
}

/****************************************************Delete Facebook Application*********************************/
func (r *crudRepository) DeleteFacebookApplication(ctx context.Context, domain_uuid string, flac_uuid string) (*models.Response, error) {
	if err := r.DBConn.Where("domain_uuid=? and flac_uuid=? ", domain_uuid, flac_uuid).Delete(&models.FacebookLoginAppConfiguration{}).Error; err != nil {
		return &models.Response{Status: "0", Msg: "Failed", ResponseCode: http.StatusBadRequest}, nil
	}
	return &models.Response{Status: "1", Msg: "Deleted", ResponseCode: http.StatusOK}, nil
}

/******************************************Assign Agent to Facebook Application******************************/
func (r *crudRepository) AssignAgentToFacebookApplication(ctx context.Context, domain_uuid string, flac_uuid string, agent_uuid string) (*models.Response, error) {
	var l int64
	l = 0
	r.DBConn.Model(&models.FacebookLoginAppConfigurationAgent{}).Where("flac_uuid=? and agent_uuid=?", flac_uuid, agent_uuid).Count(&l)
	if l == 0 {
		if err := r.DBConn.Create(&models.FacebookLoginAppConfigurationAgent{
			DomainUUID: domain_uuid,
			AgentUUID:  agent_uuid,
			FlacUUID:   flac_uuid,
		}).Error; err != nil {
			return &models.Response{Status: "0", Msg: "Failed", ResponseCode: http.StatusBadRequest}, nil
		}
		return &models.Response{Status: "1", Msg: "Assigned agent facebook account.", ResponseCode: http.StatusOK}, nil
	}
	return &models.Response{Status: "0", Msg: "Agent already assigned facebook account.", ResponseCode: http.StatusBadRequest}, nil
}

/******************************************Agent List in facebook Application************************************/
func (r *crudRepository) AgentListAssignedToFacebookApplication(ctx context.Context, flac_uuid string) (*models.Response, error) {
	t := []models.FacebookLoginAppConfigurationAgentList{}
	if err := r.DBConn.Table("facebook_login_app_configuration_agents fa").Select("fa.*,va.username as agent_name").Joins(" inner join v_call_center_agents va on va.call_center_agent_uuid::text=fa.agent_uuid").Where("flac_uuid=?", flac_uuid).Find(&t).Error; err != nil {
		return &models.Response{Status: "0", Msg: "Failed", ResponseCode: http.StatusBadRequest}, nil
	}
	if len(t) == 0 {
		return &models.Response{Status: "0", Msg: "Empty", ResponseCode: http.StatusNoContent}, nil
	}
	return &models.Response{Status: "1", Msg: "List", ResponseCode: http.StatusOK, FacebookLoginAppConfigurationAgentList: &t}, nil
}

/*************************************Agent List Not in Facaebook Application*************************************/
func (r *crudRepository) AgentListNotInFacebookApplication(ctx context.Context, flac_uuid string, domain_uuid string) (*models.Response, error) {
	t := []models.V_call_center_agents{}
	if err := r.DBConn.Table("v_call_center_agents va").Select("va.call_center_agent_uuid,va.username as agent_name").Where("va.call_center_agent_uuid::text not in (select agent_uuid from facebook_login_app_configuration_agents where flac_uuid=? ) and domain_uuid = ? ", flac_uuid, domain_uuid).Find(&t).Error; err != nil {
		return &models.Response{Status: "0", Msg: "Failed", ResponseCode: http.StatusBadRequest}, nil
	}
	if len(t) == 0 {
		return &models.Response{Status: "0", Msg: "Empty", ResponseCode: http.StatusOK}, nil
	}
	return &models.Response{Status: "1", Msg: "Assigned agent facebook account.", ResponseCode: http.StatusOK, AgentList: &t}, nil
}

/*******************************************Show Agent facebook Application****************************************/
func (r *crudRepository) ShowAgentFacebookApplication(ctx context.Context, agent_uuid string) (*models.Response, error) {
	t := []models.FacebookLoginAppConfiguration{}
	if err := r.DBConn.Model(&models.FacebookLoginAppConfiguration{}).Where("flac_uuid::text in (select flac_uuid from facebook_login_app_configuration_agents where agent_uuid=? )", agent_uuid).Find(&t).Error; err != nil {
		return &models.Response{Status: "0", Msg: "Failed", ResponseCode: http.StatusBadRequest}, nil
	}
	if len(t) == 0 {
		return &models.Response{Status: "0", Msg: "Empty", ResponseCode: http.StatusOK}, nil
	}
	return &models.Response{Status: "1", Msg: "Success", ResponseCode: http.StatusOK, FacebookLoginAppConfiguration: &t}, nil
}

/***************************************Convert access token in to longlived token*********************************/
func (r *crudRepository) Convert_Access_Token_into_Longlived_Token(ctx context.Context, clientId string, clientSecret string, exchange_token string, access_token string) ([]byte, error) {

	res, err := http.NewRequest("GET", "https://graph.facebook.com/oauth/access_token?grant_type=fb_exchange_token&client_id="+clientId+"&client_secret="+clientSecret+"&fb_exchange_token="+exchange_token+"&access_token="+access_token, nil)
	res.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(res)
	if err != nil {
		fmt.Printf("error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data), "enterrer")
		return data, nil
	}

	defer res.Body.Close()
	return nil, err
}

/**************************************Remove Assigned Agent From Facebook Application******************************/
func (r *crudRepository) RemoveAgentAssignedToFacebookApplication(ctx context.Context, agent_uuid string) (*models.Response, error) {
	al := models.FacebookLoginAppConfigurationAgent{
		AgentUUID: agent_uuid,
	}
	db := r.DBConn.Where("agent_uuid = ?", agent_uuid).Find(&al).Delete(&al)
	if db.Error != nil {
		return &models.Response{Status: "0", Msg: "Agent is not removed.", ResponseCode: 404}, nil
	}
	return &models.Response{Status: "1", Msg: "Agent successfully Removed.", ResponseCode: 200}, nil
}

/*****************************************Update Facebook Application******************************************/
func (r *crudRepository) UpdateFacebookApplication(ctx context.Context, domain_uuid string, flac_uuid string, app_id string, app_secret string, app_name string) (*models.Response, error) {
	tl := &models.FacebookLoginAppConfiguration{}
	if err := r.DBConn.Table("facebook_login_app_configurations").Where("flac_uuid = ? AND domain_uuid = ?", flac_uuid, domain_uuid).Find(&tl).Updates(map[string]interface{}{"app_id": app_id, "app_name": app_name, "app_secret": app_secret}).Error; err != nil {
		return &models.Response{Status: "0", Msg: "Failed", ResponseCode: http.StatusBadRequest}, nil
	}
	return &models.Response{Status: "1", Msg: "Updated", ResponseCode: http.StatusOK}, nil

}

/******************************************Send_private_Message**********************************************/
func (r *crudRepository) Send_Private_Message(ctx context.Context, pageId string, postId string, message string, access_token string) ([]byte, error) {
	message = strings.ReplaceAll(message, " ", "%20")

	res, err := http.NewRequest("POST", "https://graph.facebook.com/"+pageId+"/messages?recipient="+postId+"&message="+message+"&message_type=RESPONSE&access_token="+access_token, nil)
	res.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(res)
	if err != nil {
		fmt.Printf("error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data), "enterrer")
		return data, nil
	}
	defer res.Body.Close()
	return nil, err
}

/**********************************************Likes and unlike Post ans comments*******************************/
func (r *crudRepository) Like_and_Unlike_Post_and_Comment(ctx context.Context, postId string, commentId string, access_token string, Type string) ([]byte, error) {
	if Type == "Like_Post" {
		res, err := http.NewRequest("POST", "https://graph.facebook.com/"+postId+"/likes?access_token="+access_token, nil)
		res.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		response, err := client.Do(res)
		if err != nil {
			fmt.Printf("error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			fmt.Println(string(data), "enterrer")
			return data, nil
		}
		defer res.Body.Close()
		return nil, err
	} else if Type == "Like_Comment" {
		res, err := http.NewRequest("POST", "https://graph.facebook.com/"+commentId+"/likes?access_token="+access_token, nil)
		res.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		response, err := client.Do(res)
		if err != nil {
			fmt.Printf("error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			fmt.Println(string(data), "enterrer")
			return data, nil
		}
		defer res.Body.Close()
		return nil, err
	} else if Type == "Unlike_Post" {
		res, err := http.NewRequest("DELETE", "https://graph.facebook.com/"+postId+"/likes?access_token="+access_token, nil)
		res.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		response, err := client.Do(res)
		if err != nil {
			fmt.Printf("error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			fmt.Println(string(data), "enterrer")
			return data, nil
		}
		defer res.Body.Close()
		return nil, err
	} else if Type == "Unlike_Comment" {
		res, err := http.NewRequest("DELETE", "https://graph.facebook.com/"+commentId+"/likes?access_token="+access_token, nil)
		res.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		response, err := client.Do(res)
		if err != nil {
			fmt.Printf("error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			fmt.Println(string(data), "enterrer")
			return data, nil
		}
		defer res.Body.Close()

		return nil, err
	}
	return nil, nil
}

/************************************************Delete Tickets********************************************/
func (r *crudRepository) Delete_Tickets(ctx context.Context, ticket_uuid string) (*models.Response, error) {
	u := models.SocialMediaTickets{}

	err := r.DBConn.Where("ticket_uuid = ?", ticket_uuid).Find(&u)
	if err.Error != nil {
		return &models.Response{Status: "0", Msg: "Ticket not found.", ResponseCode: 404}, nil
	}
	db := r.DBConn.Where("ticket_uuid = ?", ticket_uuid).Delete(&u)
	if db.RowsAffected == 0 {
		return &models.Response{Status: "0", Msg: "Ticket not Deleted.", ResponseCode: 404}, nil
	}
	return &models.Response{Status: "1", Msg: "Ticket Delete Successfully.", ResponseCode: 200}, nil

}

/********************************************Get All Tickets**********************************************/
func (r *crudRepository) GetAll_Tickets(ctx context.Context, domain_uuid string) (*models.Response, error) {
	td := models.SocialMediaTickets{}
	list := make([]models.SocialMediaTickets, 0)
	if db := r.DBConn.Where("domain_uuid = ?", domain_uuid).Find(&td).Error; db != nil {

		return &models.Response{Status: "0", Msg: "Tickets list is not available", ResponseCode: 404}, nil
	}

	if rows, err := r.DBConn.Raw("select ticket_uuid,domain_uuid,ticket_name,customer_id,customer_name,message,message_type,integration_id,source_type,conversation_id,timestamp from social_media_tickets where domain_uuid = ?", domain_uuid).Rows(); err != nil {

		return &models.Response{Status: "Not Found", Msg: "Record Not Found", ResponseCode: 204}, nil
	} else {
		defer rows.Close()
		for rows.Next() {
			f := models.SocialMediaTickets{}
			if err := rows.Scan(&f.Ticket_uuid, &f.Domain_uuid, &f.Ticket_name, &f.CustomerId, &f.CustomerName, &f.Message, &f.MessageType, &f.IntegrationID, &f.Source_type, &f.Conversation_id, &f.Timestamp); err != nil {

				return nil, err
			}

			list = append(list, f)
		}

		return &models.Response{Status: "OK", Msg: "Record Found", ResponseCode: 200, TicketList: list}, nil
	}
}

/*******************************************Save Twitter Auth******************************************/
func (r *crudRepository) SaveTwitterAuth(ctx context.Context, id int64, domain_uuid string, api_key string, api_secret string, bearer_token string, access_token string, token_secret string) (*models.Response, error) {
	uuid1, _ := myNewUUID.NewUUID()
	uuid := uuid1.String()
	au := models.SaveTwitterAuth{
		Id:           id,
		Twitter_uuid: uuid,
		Domain_uuid:  domain_uuid,
		Api_Key:      api_key,
		Api_Secret:   api_secret,
		Bearer_Token: bearer_token,
		Access_Token: access_token,
		Token_Secret: token_secret,
	}
	if db := r.DBConn.Where("api_key = ?", api_key).Find(&au).Error; db != nil {
		err := r.DBConn.Create(&au)
		if err.Error != nil {
			return &models.Response{Status: "0", Msg: "Twitter Auth is not saved.", ResponseCode: 409}, nil
		}
		return &models.Response{Status: "1", Msg: "Twitter Auth is saved successfully.", ResponseCode: 200}, nil
	}
	if api_key == au.Api_Key {
		return &models.Response{Status: "0", Msg: "Twitter Api key is already exist.", ResponseCode: 409}, nil
	}

	return &models.Response{Status: "1", Msg: "Twitter Auth is saved successfully.", ResponseCode: 200}, nil
}

/********************************************Update Twitter Auth********************************************/
func (r *crudRepository) UpdateTwitterAuth(ctx context.Context, id int64, domain_uuid string, api_key string, api_secret string, bearer_token string, access_token string, token_secret string) (*models.Response, error) {
	au := models.SaveTwitterAuth{}

	db := r.DBConn.Where("domain_uuid = ? and id = ?", domain_uuid, id).Find(&au)
	if db.Error != nil {
		return &models.Response{Status: "0", Msg: "Twitter Auth is not found.", ResponseCode: 204}, nil
	}
	if api_key == au.Api_Key {
		up := r.DBConn.Table("save_twitter_auths").Where("domain_uuid = ? and id = ?", domain_uuid, id).Updates(map[string]interface{}{"api_secret": api_secret, "bearer_token": bearer_token, "access_token": access_token, "token_secret": token_secret})
		if up.RowsAffected == 0 {
			return &models.Response{Status: "0", Msg: "Twitter Auth is not Updated.", ResponseCode: 204}, nil
		}
		return &models.Response{Status: "1", Msg: "Twitter Auth is Updated successfully.", ResponseCode: 200}, nil
	} else if api_key != au.Api_Key {
		up := r.DBConn.Table("save_twitter_auths").Where("domain_uuid = ? and id = ?", domain_uuid, id).Updates(map[string]interface{}{"api_key": api_key, "api_secret": api_secret, "bearer_token": bearer_token, "access_token": access_token, "token_secret": token_secret})
		if up.RowsAffected == 0 {
			return &models.Response{Status: "0", Msg: "Twitter Auth is not Updated.", ResponseCode: 204}, nil
		}
		return &models.Response{Status: "1", Msg: "Twitter Auth is Updated successfully.", ResponseCode: 200}, nil
	}

	return &models.Response{Status: "1", Msg: "Twitter Auth is Updated successfully.", ResponseCode: 200}, nil
}

/************************************************Get Twitter Auth********************************************/
func (r *crudRepository) GetTwitterAuth(ctx context.Context, domain_uuid string) (*models.Response, error) {
	td := models.SaveTwitterAuth{}
	list := make([]models.SaveTwitterAuth, 0)
	if db := r.DBConn.Where("domain_uuid = ?", domain_uuid).Find(&td).Error; db != nil {

		return &models.Response{Status: "0", Msg: "Twitter Auth list is not available", ResponseCode: 404}, nil
	}
	if rows, err := r.DBConn.Raw("select id,domain_uuid,twitter_uuid,api_key, api_secret,bearer_token,access_token,token_secret from save_twitter_auths where domain_uuid = ?", domain_uuid).Rows(); err != nil {

		return &models.Response{Status: "Not Found", Msg: "Record Not Found", ResponseCode: 204}, nil
	} else {
		defer rows.Close()
		for rows.Next() {
			f := models.SaveTwitterAuth{}
			if err := rows.Scan(&f.Id, &f.Domain_uuid, &f.Twitter_uuid, &f.Api_Key, &f.Api_Secret, &f.Bearer_Token, &f.Access_Token, &f.Token_Secret); err != nil {

				return nil, err
			}

			list = append(list, f)
		}

		return &models.Response{Status: "OK", Msg: "Record Found", ResponseCode: 200, TwitterAuthList: list}, nil
	}

}

/**********************************************Delete Twitter Auth******************************************/
func (r *crudRepository) DeleteTwitterAuth(ctx context.Context, id int64, domain_uuid string) (*models.Response, error) {
	au := models.SaveTwitterAuth{}
	db := r.DBConn.Where("domain_uuid = ? and id = ?", domain_uuid, id).Find(&au).Delete(&au)
	if db.Error != nil {
		return &models.Response{Status: "0", Msg: "Twiter Auth not found.", ResponseCode: 204}, nil
	}
	if db.RowsAffected == 0 {
		return &models.Response{Status: "0", Msg: "Twiter Auth not deleted.", ResponseCode: 204}, nil
	}
	return &models.Response{Status: "1", Msg: "Twiter Auth deleted successfully.", ResponseCode: 200}, nil
}

/******************************************Get Timeline******************************************************/
func (r *crudRepository) Twitter_Apis(ctx context.Context, tweet_id string, screen_name string, api_key string, api_type string, author_id string, message string) ([]byte, error) {
	value := models.SaveTwitterAuth{}
	db := r.DBConn.Where("api_key = ?", api_key).Find(&value)
	if db.Error != nil {
		fmt.Println("error")
	}
	if api_type == "getall_tweets" {
		fmt.Println("part1")
		auth := oauth1.OAuth1{
			ConsumerKey:    value.Api_Key,
			ConsumerSecret: value.Api_Secret,
			AccessToken:    value.Access_Token,
			AccessSecret:   value.Token_Secret,
		}
		fmt.Println(value.Api_Key, value.Api_Secret, value.Access_Token, value.Token_Secret, "credential...s")
		authHeader := auth.BuildOAuth1Header("GET", "https://api.twitter.com/1.1/statuses/user_timeline.json?screen_name="+screen_name, map[string]string{
			"screen_name": screen_name,
		})

		req, _ := http.NewRequest("GET", "https://api.twitter.com/1.1/statuses/user_timeline.json?screen_name="+screen_name, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authHeader)
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		} else {
			data, _ := ioutil.ReadAll(res.Body)
			return data, nil
		}
	} else if api_type == "getsingle_tweet" {
		fmt.Println("part2")
		auth := oauth1.OAuth1{
			ConsumerKey:    value.Api_Key,
			ConsumerSecret: value.Api_Secret,
			AccessToken:    value.Access_Token,
			AccessSecret:   value.Token_Secret,
		}

		authHeader := auth.BuildOAuth1Header("GET", "https://api.twitter.com/1.1/statuses/show.json?id="+tweet_id, map[string]string{
			"id": tweet_id,
		})
		req, _ := http.NewRequest("GET", "https://api.twitter.com/1.1/statuses/show.json?id="+tweet_id, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authHeader)
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		} else {
			data, _ := ioutil.ReadAll(res.Body)
			return data, nil
		}
	} else if api_type == "retweet" {
		fmt.Println("part3")
		auth := oauth1.OAuth1{
			ConsumerKey:    value.Api_Key,
			ConsumerSecret: value.Api_Secret,
			AccessToken:    value.Access_Token,
			AccessSecret:   value.Token_Secret,
		}

		authHeader := auth.BuildOAuth1Header("POST", "https://api.twitter.com/1.1/statuses/retweet/"+tweet_id+".json", nil)

		req, _ := http.NewRequest("POST", "https://api.twitter.com/1.1/statuses/retweet/"+tweet_id+".json", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authHeader)
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		} else {
			data, _ := ioutil.ReadAll(res.Body)
			return data, nil
		}
	} else if api_type == "un-retweet" {
		fmt.Println("part4")
		auth := oauth1.OAuth1{
			ConsumerKey:    value.Api_Key,
			ConsumerSecret: value.Api_Secret,
			AccessToken:    value.Access_Token,
			AccessSecret:   value.Token_Secret,
		}

		authHeader := auth.BuildOAuth1Header("POST", "https://api.twitter.com/1.1/statuses/unretweet/"+tweet_id+".json", nil)

		req, _ := http.NewRequest("POST", "https://api.twitter.com/1.1/statuses/unretweet/"+tweet_id+".json", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authHeader)
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		} else {
			data, _ := ioutil.ReadAll(res.Body)
			return data, nil
		}
	} else if api_type == "retweet_status" {
		fmt.Println("part5")
		auth := oauth1.OAuth1{
			ConsumerKey:    value.Api_Key,
			ConsumerSecret: value.Api_Secret,
			AccessToken:    value.Access_Token,
			AccessSecret:   value.Token_Secret,
		}

		authHeader := auth.BuildOAuth1Header("GET", "https://api.twitter.com/1.1/statuses/retweets/"+tweet_id+".json", nil)

		req, _ := http.NewRequest("GET", "https://api.twitter.com/1.1/statuses/retweets/"+tweet_id+".json", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authHeader)
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		} else {
			data, _ := ioutil.ReadAll(res.Body)
			return data, nil
		}

	} else if api_type == "like_tweet" {
		fmt.Println("part6")
		auth := oauth1.OAuth1{
			ConsumerKey:    value.Api_Key,
			ConsumerSecret: value.Api_Secret,
			AccessToken:    value.Access_Token,
			AccessSecret:   value.Token_Secret,
		}

		authHeader := auth.BuildOAuth1Header("POST", "https://api.twitter.com/1.1/favorites/create.json?id="+tweet_id, map[string]string{
			"id": tweet_id,
		})

		req, _ := http.NewRequest("POST", "https://api.twitter.com/1.1/favorites/create.json?id="+tweet_id, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authHeader)
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		} else {
			data, _ := ioutil.ReadAll(res.Body)
			return data, nil
		}
	} else if api_type == "unlike_tweet" {
		fmt.Println("part7")
		auth := oauth1.OAuth1{
			ConsumerKey:    value.Api_Key,
			ConsumerSecret: value.Api_Secret,
			AccessToken:    value.Access_Token,
			AccessSecret:   value.Token_Secret,
		}

		authHeader := auth.BuildOAuth1Header("POST", "https://api.twitter.com/1.1/favorites/destroy.json?id="+tweet_id, map[string]string{
			"id": tweet_id,
		})

		req, _ := http.NewRequest("POST", "https://api.twitter.com/1.1/favorites/destroy.json?id="+tweet_id, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authHeader)
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		} else {
			data, _ := ioutil.ReadAll(res.Body)
			return data, nil
		}
	} else if api_type == "get_user" {
		fmt.Println("part8")
		auth := oauth1.OAuth1{
			ConsumerKey:    value.Api_Key,
			ConsumerSecret: value.Api_Secret,
			AccessToken:    value.Access_Token,
			AccessSecret:   value.Token_Secret,
		}

		authHeader := auth.BuildOAuth1Header("GET", "https://api.twitter.com/2/users/"+author_id, nil)
		req, _ := http.NewRequest("GET", "https://api.twitter.com/2/users/"+author_id, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authHeader)
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		} else {
			data, _ := ioutil.ReadAll(res.Body)
			return data, nil
		}
	} else if api_type == "getreply_list" {
		fmt.Println("part9")
		auth := oauth1.OAuth1{
			ConsumerKey:    value.Api_Key,
			ConsumerSecret: value.Api_Secret,
			AccessToken:    value.Access_Token,
			AccessSecret:   value.Token_Secret,
		}

		authHeader := auth.BuildOAuth1Header("GET", "https://api.twitter.com/1.1/statuses/mentions_timeline.json?count=10", map[string]string{
			"count": "10",
		})
		req, _ := http.NewRequest("GET", "https://api.twitter.com/1.1/statuses/mentions_timeline.json?count=10", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authHeader)
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		} else {
			data, _ := ioutil.ReadAll(res.Body)
			return data, nil
		}
	} else if api_type == "getquoted_retweet_list" {
		fmt.Println("part10")
		auth := oauth1.OAuth1{
			ConsumerKey:    value.Api_Key,
			ConsumerSecret: value.Api_Secret,
			AccessToken:    value.Access_Token,
			AccessSecret:   value.Token_Secret,
		}

		authHeader := auth.BuildOAuth1Header("GET", "https://api.twitter.com/2/tweets/search/recent?query="+tweet_id+"&expansions=attachments.media_keys,author_id,entities.mentions.username,in_reply_to_user_id,referenced_tweets.id,referenced_tweets.id.author_id", map[string]string{
			"query":      tweet_id,
			"expansions": "attachments.media_keys,author_id,entities.mentions.username,in_reply_to_user_id,referenced_tweets.id,referenced_tweets.id.author_id",
		})
		req, _ := http.NewRequest("GET", "https://api.twitter.com/2/tweets/search/recent?query="+tweet_id+"&expansions=attachments.media_keys,author_id,entities.mentions.username,in_reply_to_user_id,referenced_tweets.id,referenced_tweets.id.author_id", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authHeader)
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		} else {
			data, _ := ioutil.ReadAll(res.Body)
			return data, nil
		}
	} else if api_type == "tweet_reply" {
		fmt.Println("part11")
		auth := oauth1.OAuth1{
			ConsumerKey:    value.Api_Key,
			ConsumerSecret: value.Api_Secret,
			AccessToken:    value.Access_Token,
			AccessSecret:   value.Token_Secret,
		}

		authHeader := auth.BuildOAuth1Header("POST", "https://api.twitter.com/1.1/statuses/update.json?status= welcome&in_reply_to_status_id="+tweet_id+"&auto_populate_reply_metadata=true", map[string]string{

			"status":                       message,
			"in_reply_to_status_id":        tweet_id,
			"auto_populate_reply_metadata": "true",
		})
		req, _ := http.NewRequest("POST", "https://api.twitter.com/1.1/statuses/update.json?status="+message+"&in_reply_to_status_id="+tweet_id+"&auto_populate_reply_metadata=true", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authHeader)
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		} else {
			data, _ := ioutil.ReadAll(res.Body)
			return data, nil
		}
	} else if api_type == "delete_reply" {
		fmt.Println("part12")
		auth := oauth1.OAuth1{
			ConsumerKey:    value.Api_Key,
			ConsumerSecret: value.Api_Secret,
			AccessToken:    value.Access_Token,
			AccessSecret:   value.Token_Secret,
		}

		authHeader := auth.BuildOAuth1Header("POST", "https://api.twitter.com/1.1/statuses/destroy/"+tweet_id+".json", nil)
		req, _ := http.NewRequest("POST", "https://api.twitter.com/1.1/statuses/destroy/"+tweet_id+".json", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authHeader)
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		} else {
			data, _ := ioutil.ReadAll(res.Body)
			return data, nil
		}
	}
	return nil, nil
}

/*******************************************Assign agent To Twitter**************************************/
func (r *crudRepository) AssignAgentToTwitter(ctx context.Context, twitter_uuid string, domain_uuid string, api_key string, agent_uuid string) (*models.Response, error) {
	asgn := models.TwitterAssignedAgents{
		Twitter_uuid: twitter_uuid,
		Domain_uuid:  domain_uuid,
		Agent_uuid:   agent_uuid,
		Api_Key:      api_key,
	}
	if db := r.DBConn.Table("twitter_assigned_agents").Where("agent_uuid = ? and twitter_uuid = ?", agent_uuid, twitter_uuid).Find(&asgn).Error; db != nil {
		create := r.DBConn.Create(&asgn)
		if create.Error != nil {
			return &models.Response{Status: "0", Msg: "Agent is not assigned.", ResponseCode: 204}, nil
		}
		return &models.Response{Status: "1", Msg: "Agent Assigned Successfully.", ResponseCode: 200}, nil
	}
	return &models.Response{Status: "0", Msg: "Agent Already Assigned.", ResponseCode: 204}, nil

}

/*****************************************Assign agent Twitter List***************************************/
func (r *crudRepository) TwitterAssignAgentList(ctx context.Context, domain_uuid string, twitter_uuid string) (*models.Response, error) {
	td := models.TwitterAssignedAgents{}
	list := make([]models.TwitterAssignedAgents, 0)
	if db := r.DBConn.Where("domain_uuid = ? and twitter_uuid = ?", domain_uuid, twitter_uuid).Find(&td).Error; db != nil {

		return &models.Response{Status: "0", Msg: "Agent list is not available", ResponseCode: 404}, nil
	}
	if rows, err := r.DBConn.Raw("select domain_uuid, agent_uuid, twitter_uuid, api_key from twitter_assigned_agents where domain_uuid = ? and twitter_uuid = ?", domain_uuid, twitter_uuid).Rows(); err != nil {

		return &models.Response{Status: "Not Found", Msg: "Record Not Found", ResponseCode: 204}, nil
	} else {
		defer rows.Close()
		for rows.Next() {
			f := models.TwitterAssignedAgents{}
			if err := rows.Scan(&f.Domain_uuid, &f.Agent_uuid, &f.Twitter_uuid, &f.Api_Key); err != nil {

				return nil, err
			}

			list = append(list, f)
		}

		return &models.Response{Status: "OK", Msg: "Record Found", ResponseCode: 200, TwitterAssignAgentList: list}, nil
	}

}

/************************************************Remove Assigned Agent Twitter***********************************/
func (r *crudRepository) RemoveTwitterAssignAgent(ctx context.Context, agent_uuid string, twitter_uuid string) (*models.Response, error) {
	agt := models.TwitterAssignedAgents{}

	if db := r.DBConn.Table("twitter_assigned_agents").Where("agent_uuid = ? and twitter_uuid = ?", agent_uuid, twitter_uuid).Find(&agt).Error; db != nil {
		return &models.Response{Status: "0", Msg: "Agent Not removed.", ResponseCode: 204}, nil
	}
	del := r.DBConn.Delete(&agt)
	if del.RowsAffected == 0 {
		return &models.Response{Status: "0", Msg: "Agent not Found.", ResponseCode: 204}, nil
	}
	return &models.Response{Status: "1", Msg: "Agent Removed successfully.", ResponseCode: 200}, nil
}

/**********************************************Get Quoted  Retweet List************************************/
func (r *crudRepository) Get_Quoted_Retweet_List(ctx context.Context, api_key string, tweet_id string) (*models.Response, error) {
	value := models.SaveTwitterAuth{}
	db := r.DBConn.Where("api_key = ?", api_key).Find(&value)
	if db.Error != nil {
		fmt.Println("error")
	}

	val := models.QuotedTweet{}
	datas := []models.QuotedTweet{}
	vak := make([]models.Result, 0)
	auth := oauth1.OAuth1{
		ConsumerKey:    value.Api_Key,
		ConsumerSecret: value.Api_Secret,
		AccessToken:    value.Access_Token,
		AccessSecret:   value.Token_Secret,
	}

	authHeader := auth.BuildOAuth1Header("GET", "https://api.twitter.com/2/tweets/search/recent?query="+tweet_id+"&expansions=attachments.media_keys,author_id,entities.mentions.username,in_reply_to_user_id,referenced_tweets.id,referenced_tweets.id.author_id&tweet.fields=created_at", map[string]string{
		"query":        tweet_id,
		"expansions":   "attachments.media_keys,author_id,entities.mentions.username,in_reply_to_user_id,referenced_tweets.id,referenced_tweets.id.author_id",
		"tweet.fields": "created_at",
	})
	req, _ := http.NewRequest("GET", "https://api.twitter.com/2/tweets/search/recent?query="+tweet_id+"&expansions=attachments.media_keys,author_id,entities.mentions.username,in_reply_to_user_id,referenced_tweets.id,referenced_tweets.id.author_id&tweet.fields=created_at", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return &models.Response{Status: "0", Msg: "Quote Retweet List not found.", ResponseCode: 204}, nil
	} else {
		data, _ := ioutil.ReadAll(res.Body)
		json := json.Unmarshal(data, &val)
		datas = append(datas, val)
		fmt.Println(json, getServerOs())
		j := 0
		for i := 0; i < len(val.Data) && j < len(val.Includes.Users); i++ {

			result := models.Result{
				Text:      val.Data[i].Text,
				Author_id: val.Includes.Users[j].Id,
				Username:  val.Includes.Users[j].Username,
				Name:      val.Includes.Users[j].Name,
				RetweetId: val.Data[i].Author_id,
				CreatedAt: val.Data[i].CreatedAt,
			}
			vak = append(vak, result)

		}

		return &models.Response{Status: "1", Msg: "Quote Retweet list", ResponseCode: 200, Quote_retweet_list: vak}, nil

	}

}
