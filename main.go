package main

import (
	"fmt"
	"sync"
	adapterdatabase "whatsapp_api/adapter"
	config "whatsapp_api/config"

	// logging "whatsapp_api/logger"
	crudController "whatsapp_api/whatsapp/controller"
	crudRepo "whatsapp_api/whatsapp/repository"
	crudUsecase "whatsapp_api/whatsapp/usecase"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var onceRest sync.Once

func main() {
	onceRest.Do(func() {
		e := echo.New()
		//Setting up the config
		config := config.GetConfig()
		//Setting up the Logger
		//logger := logging.NewLogger(config.Log.LogFile, config.Log.LogLevel)
		e.Use(middleware.Logger())
		e.Use(middleware.Recover())

		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		}))

		db := adapterdatabase.DB(config)
		newServerUser := crudController.NewServerUser(db)
		crudRepo := crudRepo.NewcrudRepository(db, newServerUser, config)
		crudUc := crudUsecase.NewcrudUsecase(crudRepo)
		crudController.NewCRUDController(e, crudUc)

		go newServerUser.Controller(e)

		if err := e.StartTLS(config.HttpConfig.HostPort, config.HttpConfig.HostCert, config.HttpConfig.HostKey); err != nil {
			fmt.Println("not connected")
			//logger.WithError(err).Fatal("Unable to start the callCenter service")
		}

	})
}
