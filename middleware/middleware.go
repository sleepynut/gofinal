package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Auth(c *gin.Context) {
	log.Println("Start authentication")
	if c.GetHeader("Authorization") != "November 10, 2009" {
		c.JSON(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))

		// need to abort gin
		c.Abort()
		return
	}

	// move to next component behind AUTHEN middleware
	c.Next()
	log.Println("Exit authentication")
}
