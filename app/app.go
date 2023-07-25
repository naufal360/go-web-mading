package app

import (
	"fmt"
	"jewepe/config"
	"jewepe/routes"

	"github.com/gin-gonic/gin"
)

var router = gin.New()

func StartApp() {
	db := config.InitGorm()
	routes.RouteAPI(router, db)

	port := "8000"

	// run gin
	router.Run(fmt.Sprintf(":%s", port))
}
