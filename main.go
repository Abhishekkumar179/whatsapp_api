package main

import (
	"fmt"
	"sync"
	adapterdatabase "whatsapp_api/adapter"
	config "whatsapp_api/config"
	logging "whatsapp_api/logger"
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
		logger := logging.NewLogger(config.Log.LogFile, config.Log.LogLevel)
		e.Use(middleware.Logger())
		e.Use(middleware.Recover())

		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		}))

		db := adapterdatabase.DB(config)
		crudRepo := crudRepo.NewcrudRepository(db)
		crudUc := crudUsecase.NewcrudUsecase(crudRepo)
		crudController.NewCRUDController(e, crudUc)

		if err := e.Start(":10000"); err != nil {
			fmt.Println("not connected")
			logger.WithError(err).Fatal("avb")
		}

	})
}
