package routes

import (
	"jewepe/controllers"
	"jewepe/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RouteAPI(r *gin.Engine, app *gorm.DB) {
	// controller
	server := controllers.NewHttpServer(app)

	// external library
	// use logger
	r.Use(gin.Logger())
	// use recovery
	r.Use(gin.Recovery())
	// serve static file
	r.Use(static.Serve("/static", static.LocalFile("./static", true)))
	// load html
	r.LoadHTMLGlob("templates/*")
	// use session
	r.Use(sessions.Sessions("session", cookie.NewStore([]byte("secret"))))

	// all routes jewepe
	// jewepe index
	jewepe := r.Group("/jewepe")

	// test
	jewepe.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// home
	jewepe.GET("/", server.ViewIndex)

	// login register
	jewepe.POST("/register", server.RegisterUsers)
	jewepe.GET("/login", server.ViewLogin)
	jewepe.POST("/login", server.LoginUsers)
	jewepe.GET("/logout", server.LogoutUser)

	// users pembaca
	jewepe.GET("/artikel/:id", server.ViewDetailArtikel)
	jewepe.GET("/users/artikel/:id/komentar", server.ViewCreateKomentar)
	jewepe.POST("/users/artikel/:id/komentar", server.CreateKomentar)

	// users admin
	admin := jewepe.Group("/dashboard")
	admin.Use(middleware.AuthRequired)
	admin.GET("/", server.ViewDashboard)
	admin.GET("/create", server.ViewCreateArtikel)
	admin.POST("/create", server.CreateArtikel)
	admin.DELETE("/delete/:id", server.DeleteArtikel)
	admin.DELETE("komentar/delete/:id", server.DeleteKomentar)
}
