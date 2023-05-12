package routes

import (
	"rpsoftech/scaleMQTT/src/db"
	"rpsoftech/scaleMQTT/src/global"
	"rpsoftech/scaleMQTT/src/systypes"

	"github.com/gin-gonic/gin"
)

func ChangeDeviceID(c *gin.Context) {

	config := systypes.ChangeDeviceIdData{}
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(400, systypes.BaseResponseFormat{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	if config.DevID == config.OldDevID {
		c.JSON(400, systypes.BaseResponseFormat{
			Success: false,
			Error:   "Old and New Device Id Cannot be same",
		})
	}
	configObject, err := db.DBClassObject.GetScaleConfigData(config.DevID)
	if err != nil {
		c.JSON(400, systypes.BaseResponseFormat{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	byteConfig, err := configObject.JSON()
	if err != nil {
		c.JSON(400, systypes.BaseResponseFormat{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	if config.OldDevID != global.RootDefaultDeviceId {
		err = db.DBClassObject.SetChanegedDeviceId(config.OldDevID, configObject.DevID)
		if err != nil {
			c.JSON(400, systypes.BaseResponseFormat{
				Success: false,
				Error:   err.Error(),
			})
			return
		}
	}
	println(config.OldDevID + global.DefaultMQTTDeviceSubscribeTopicSuffix)
	global.MQTTserver.Publish(config.OldDevID+global.DefaultMQTTDeviceSubscribeTopicSuffix, []byte("{\"mqttcfg\":"+string(byteConfig)+"}"), true, 2)
	global.MQTTserver.Publish(config.DevID+global.DefaultMQTTDeviceSubscribeTopicSuffix, []byte("{\"mqttcfg\":"+string(byteConfig)+"}"), true, 2)
	c.JSON(200, config)
}
