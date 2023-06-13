package routes

import (
	"rpsoftech/scaleMQTT/src/db"
	"rpsoftech/scaleMQTT/src/global"
	"rpsoftech/scaleMQTT/src/middlerware"

	"github.com/gin-gonic/gin"
)

func GetAllDataFromDatabase(c *gin.Context) {
	c.JSON(200, db.DBClassObject.TakeBackup())
}

func GetAllMqttStatus(c *gin.Context) {
	if val, exist := c.GetQuery("GetAll"); exist && val == "true" {
		// global.MQTTConnectionWithUidStatusMap
		for _, v := range global.MQTTConnectionWithUidStatusMap {
			if v.Cl.Closed() {
				v.Connected = false
			}
		}
		c.JSON(200, global.MQTTConnectionWithUidStatusMap)
	} else if val, exist := c.GetQuery("username"); exist {
		if v, ok := global.MQTTConnectionWithUidStatusMap[val]; ok {
			if v.Cl.Closed() {
				v.Connected = false
			}
			c.JSON(200, v)
		} else {
			c.String(400, "Not Connected")
		}
	} else {

		c.String(400, "Please Pass Valid Params")
	}
}

func GetScaleData(c *gin.Context) {
	var UserConfig *middlerware.UserClaims
	if val, ok := c.Get("User"); !ok {
		c.String(500, "Something went wrong User Data Not Found")
	} else {
		UserConfig = val.(*middlerware.UserClaims)
	}
	if val, ok := global.MQTTConnectionWithUidStatusMap[UserConfig.Username]; ok {
		if !val.Connected {
			c.String(400, "Not Connected")
			return
		}
		if val.Cl.Closed() {
			val.Connected = false
		}
		c.JSON(200, val)
	} else {
		c.String(400, "Not Connected")
		return
	}
}
