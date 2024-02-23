package bootstrap

import (
	"checkoutProject/pkg/common/database"
	"checkoutProject/pkg/common/env"
	"checkoutProject/pkg/common/logger"
	"checkoutProject/pkg/handlers/cart"
	"checkoutProject/pkg/handlers/item"
	"github.com/gin-gonic/gin"
)

func Initialize() error {
	log, err := logger.Initialize()
	if err != nil {
		return err
	}

	err = env.Load()
	if err != nil {
		return err
	}
	log.Info("read environment variables")

	err = database.Initialize()
	if err != nil {
		return err
	}
	log.Info("connected to the database")

	return nil
}

func RegisterRouters(r *gin.Engine) {
	apiRouter := r.Group("/api/cart")
	item.NewDefaultItemRouter().Register(apiRouter)
	item.NewDefaultVasItemRouter().Register(apiRouter)
	cart.NewDefaultCartRouter().Register(apiRouter)
}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	RegisterRouters(r)
	return r
}
