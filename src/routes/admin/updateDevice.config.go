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
		NegativeChar:     "\\f",
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
	if val, ok := global.MQTTConnectionWithUidStatusMap[config.DevID]; ok {
		val.Config = config
	}
	topic := config.DevID + global.DefaultMQTTDeviceSubscribeTopicSuffix
	global.MQTTserver.Publish(topic, []byte("{\"devcfg\":"+string(byteConfig)+"}"), false, 2)
	global.MQTTserver.Publish(topic, []byte("{\"mqttcfg\":"+string(byteConfig)+"}"), false, 2)
	c.String(200, string(byteConfig))
}
