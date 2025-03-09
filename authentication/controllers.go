package authentication

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartEmailAuth(c *gin.Context){
	var body StartEmailAuthRequest

	if c.Bind(body) != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "error reading request body",
		})
	}

	// Pass to service

	// Send response
}