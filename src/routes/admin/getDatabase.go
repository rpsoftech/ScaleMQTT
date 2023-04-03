package routes

import (
	"rpsoftech/scaleMQTT/src/db"
	"rpsoftech/scaleMQTT/src/global"

	"github.com/gin-gonic/gin"
)

func GetAllDataFromDatabase(c *gin.Context) {
	c.JSON(200, db.TakeBackup())
}

func GetAllMqttStatus(c *gin.Context) {
	if val, exist := c.GetQuery("GetAll"); exist && val == "true" {
		c.JSON(200, global.MQTTConnectionStatusMap)
	} else if val, exist := c.GetQuery("username"); exist {
		if v, ok := global.MQTTConnectionStatusMap[val]; ok {
			c.JSON(200, v)
		} else {
			c.String(400, "Not Connected")
		}
	} else {

		c.String(400, "Please Pass Valid Params")
	}
}
