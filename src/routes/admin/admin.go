package routes

import (
	"net/http"
	"strings"

	"rpsoftech/scaleMQTT/src/middlerware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
	// RefreshToken string `json:"refreshToken"`
}

func AdminRoutes(router *gin.Engine) {
	router.POST("/admin/login", AdminLoginFunction)

	adminRouter := router.Group("/admin")
	{
		adminRouter.Use(middlerware.JwtAuthMiddleware(true))
		adminRouter.POST("/addNewAdmin", AddNewAdminUser)
		adminRouter.GET("/databaseSnapshot", GetAllDataFromDatabase)
		adminRouter.GET("/mqttStatus", GetAllMqttStatus)
		adminRouter.POST("/changeDeviceID", ChangeDeviceID)
		adminRouter.GET("/isLoggedin", IsAdminLoggedIn)
		adminRouter.POST("/updateConfig", UpdateDeviceConfig)
		adminRouter.GET("/genNewDevId", func(c *gin.Context) {
			id := uuid.New()
			c.JSON(http.StatusOK, gin.H{"id": strings.ReplaceAll(id.String(), "-", "")})
		})
		adminRouter.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "Welcome Gin Server")
		})
	}

	router.GET("/scaleData", middlerware.JwtAuthMiddleware(false), GetScaleData)
}

func IsAdminLoggedIn(c *gin.Context) {
	if val, ok := c.Get("User"); ok {
		c.JSON(200, gin.H{"user": val})
	} else {
		c.JSON(200, gin.H{"error": "User not found in context"})
	}
}
