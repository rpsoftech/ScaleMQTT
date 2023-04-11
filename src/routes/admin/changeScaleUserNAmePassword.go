package routes

import (
	"net/http"

	"rpsoftech/scaleMQTT/src/db"
	"rpsoftech/scaleMQTT/src/validator"

	"github.com/gin-gonic/gin"
)

func AddChangeMqttUserNameAndPassword(c *gin.Context) {
	var user AddNewAdminLogin

	// Bind the JSON body to the user struct and validate the fields
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validator.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.AddScaleUserNamePassword(user.UserName, []byte(user.Password)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"success": true})
}
