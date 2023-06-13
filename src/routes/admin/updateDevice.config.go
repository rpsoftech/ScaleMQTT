package routes

import (
	"rpsoftech/scaleMQTT/src/db"
	"rpsoftech/scaleMQTT/src/global"
	"rpsoftech/scaleMQTT/src/systypes"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type DevIdStruct struct {
	DevId string `json:"dev_id" validate:"required" `
}

func ValidateDeviceConfigAndSave(c *gin.Context, config *systypes.ScaleConfigData) {
	if err := c.ShouldBindBodyWith(config, binding.JSON); err != nil {
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
	go func(topic string, byteConfig string) {
		time.Sleep(100 * time.Millisecond)
		global.MQTTserver.Publish(topic, []byte("{\"devcfg\":"+byteConfig+"}"), false, 2)
		time.Sleep(100 * time.Millisecond)
		global.MQTTserver.Publish(topic, []byte("{\"mqttcfg\":"+byteConfig+"}"), false, 2)
	}(topic, string(byteConfig))
	c.String(200, string(byteConfig))
}

func UpdateDeviceConfig(c *gin.Context) {
	dev := DevIdStruct{}
	if err := c.ShouldBindBodyWith(&dev, binding.JSON); err != nil {
		c.JSON(400, systypes.BaseResponseFormat{
			Success: false,
			Error:   "Please Pass Valid Device Id",
		})
		return
	}
	config, err := db.DBClassObject.GetScaleConfigData(dev.DevId)
	if err != nil {
		c.JSON(400, systypes.BaseResponseFormat{
			Success: false,
			Error:   "Please Use Valid Device Id,Data Does not exitst",
		})
		return
	}
	ValidateDeviceConfigAndSave(c, &config)
}

func AddDeviceConfig(c *gin.Context) {
	dev := DevIdStruct{}
	if err := c.ShouldBindBodyWith(&dev, binding.JSON); err != nil {
		c.JSON(400, systypes.BaseResponseFormat{
			Success: false,
			Error:   "Invalid JSON Data " + err.Error(),
		})
		return
	}
	config := &systypes.ScaleConfigData{
		DivideMultiplyBy: 1,
		NegativeChar:     "\f",
		DevcfgForMqtt: systypes.DevcfgForMqtt{
			LogEnable: false,
		},
	}
	if configFromDb, err := db.DBClassObject.GetScaleConfigData(dev.DevId); err == nil {
		config = &configFromDb
	}
	ValidateDeviceConfigAndSave(c, config)
}
