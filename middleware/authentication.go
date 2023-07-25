package middleware

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	admin := session.Get("adminkey")
	if admin == nil {
		log.Println("Admin not logged in")
		c.Redirect(http.StatusMovedPermanently, "/jewepe/login")
		c.Abort()
		return
	}
	c.Next()
}
