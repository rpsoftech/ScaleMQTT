package routes

import (
	"net/http"
	"rpsoftech/scaleMQTT/src/global"
	"rpsoftech/scaleMQTT/src/middlerware"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type TokenGenerateBody struct {
	DevId          string `json:"dev_id"`
	ExpirationTime string `json:"expiration_time,omitempty"`
}

func GenerateToken(c *gin.Context) {
	reqBody := TokenGenerateBody{}

	if errA := c.ShouldBindJSON(&reqBody); errA != nil {
		c.String(400, `Invalid JSON input`)
		return
		// At this time, it reuses body stored in the context.
	}
	claims := &middlerware.UserClaims{
		Username: reqBody.DevId,
		Role:     "User",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(global.JWTKEY)

	c.JSON(http.StatusOK, gin.H{"token": tokenString, "dev_id": reqBody.DevId})
}
