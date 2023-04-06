package routes

import (
	"net/http"
	"rpsoftech/scaleMQTT/src/db"
	"rpsoftech/scaleMQTT/src/global"

	"github.com/gin-gonic/gin"
)

type AddNewAdminLogin struct {
	UserName string `json:"username" form:"username" validate:"required" binding:"required"`
	Password string `json:"password" form:"password" validate:"required" binding:"required"`
}

func AddNewAdminUser(c *gin.Context) {
	var user AddNewAdminLogin

	// Bind the JSON body to the user struct and validate the fields
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := global.Validator.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := db.DBClassObject.AddNewAdmin([]byte(`adminuser/`+user.UserName), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}
