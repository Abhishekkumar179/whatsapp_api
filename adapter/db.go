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
	fmt.Println("Successfully connected!")
	return db
}
