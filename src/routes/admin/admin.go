package routes

import (
	"net/http"

	"rpsoftech/scaleMQTT/src/middlerware"

	"github.com/gin-gonic/gin"
)

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
	// RefreshToken string `json:"refreshToken"`
}

func AdminRoutes(router *gin.Engine) {
	router.POST("/admin/login", AdminLoginFunction)

	adminRouter := router.Group("/admin")
	{
		adminRouter.Use(middlerware.JwtAuthMiddleware())
		adminRouter.POST("/addNewAdmin", AddNewAdminUser)
		adminRouter.GET("/databaseSnapshot", GetAllDataFromDatabase)
		adminRouter.GET("/mqttStatus", GetAllMqttStatus)
		adminRouter.POST("/modifyMqttUserNamePassword", AddChangeMqttUserNameAndPassword)
		adminRouter.GET("/isLoggedin", IsAdminLoggedIn)
		adminRouter.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "Welcome Gin Server")
		})
	}
}

func IsAdminLoggedIn(c *gin.Context) {
	if val, ok := c.Get("User"); ok {
		c.JSON(200, gin.H{"user": val})
	} else {
		c.JSON(200, gin.H{"error": "User not found in context"})
	}
}
