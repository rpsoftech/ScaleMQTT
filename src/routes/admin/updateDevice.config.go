package routes

import (
	"rpsoftech/scaleMQTT/src/systypes"

	"github.com/gin-gonic/gin"
)

func UpdateDeviceConfig(c *gin.Context) {

	config := systypes.ScaleConfigData{
		// DivideMultiplyBy: 1,
	}
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(400, systypes.BaseResponseFormat{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	if _, err := config.Validate(); err != nil {
		c.JSON(400, systypes.BaseResponseFormat{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	c.JSON(200, config)
}
