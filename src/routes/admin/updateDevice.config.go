package routes

import (
	"rpsoftech/scaleMQTT/src/db"
	"rpsoftech/scaleMQTT/src/global"
	"rpsoftech/scaleMQTT/src/systypes"

	"github.com/gin-gonic/gin"
)

func UpdateDeviceConfig(c *gin.Context) {

	config := systypes.ScaleConfigData{
		DivideMultiplyBy: 1,
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
	// config.
	byteConfig, err := config.JSON()
	if err != nil {
		c.JSON(400, systypes.BaseResponseFormat{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	err = db.DBClassObject.SetScaleConfigData(config.DevID, config)
	if err != nil {
		c.JSON(400, systypes.BaseResponseFormat{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	global.MQTTserver.Publish(config.DevID+global.DefaultMQTTDeviceSubscribeTopicSuffix, byteConfig, true, 2)
	c.JSON(200, gin.H{"s": string(byteConfig)})
}
