package routes

import (
	"rpsoftech/scaleMQTT/src/db"

	"github.com/gin-gonic/gin"
)

func GetAllDataFromDatabase(c *gin.Context) {
	// if val, err := json.Marshal(db.TakeBackup()); err == nil {
	// 	c.String(200, string(val))
	// } else {
	// 	c.String(500, err.Error())
	// }
	c.JSON(200, db.TakeBackup())
}
