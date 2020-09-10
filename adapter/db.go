package db

import (
	"fmt"
	config "whatsapp_api/config"
	models "whatsapp_api/model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
)

func DB(config *config.Config) *gorm.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Database.Host, config.Database.Port, config.Database.User, config.Database.Pass, config.Database.DBName)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(models.Appuser{})
	db.AutoMigrate(models.ReceiveUserDetails{})
	db.AutoMigrate(models.Tenant_details{})
	db.AutoMigrate(models.WhatsappConfiguration{})
	db.AutoMigrate(models.FacebookConfiguration{})
	db.AutoMigrate(models.Queue{})
	db.AutoMigrate(models.AgentQueue{})
	db.AutoMigrate(models.Customer_Agents{})
	db.AutoMigrate(models.FacebookLoginAppConfiguration{})
	fmt.Println("Successfully connected!")
	return db
}

// SELECT agent_uuid, domain_uuid FROM agent_queues INNER JOIN queues
//         ON agent_queues.queue_uuid = queues.queue_uuid;
// select count(agent_uuid),aq.agent_uuid from customer_agents where agent_uuid in  (select aq.agent_uuid from agent_queues aq where aq.queue_uuid=(select queue_uuid from queues where integration_id='5edbb2eb385adf000f97cddb'));
//select count(aq.agent_uuid),aq.agent_uuid from customer_agents ca right join (select agent_uuid from agent_queues where queue_uuid=(select queue_uuid from queues where integration_id='5edbb2eb385adf000f97cddb')) aq on aq.agent_uuid=ca.agent_uuid group by aq.agent_uuid;
