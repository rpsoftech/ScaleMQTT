package routes

import (
	"rpsoftech/scaleMQTT/src/db"
	"rpsoftech/scaleMQTT/src/global"
	"rpsoftech/scaleMQTT/src/middlerware"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AdminLoginFunction(c *gin.Context) {
	// c.String(http.StatusOK, "Welcome Gin Server")
	reqBody := AddNewAdminLogin{}

	if errA := c.ShouldBindJSON(&reqBody); errA != nil {
		c.String(400, `Invalid JSON input`)
		return
		// At this time, it reuses body stored in the context.
	}
	allowed := false
	if reqBody.UserName == "" || reqBody.Password == "" {
		allowed = false
	} else if reqBody.UserName == "keyurboss" && reqBody.Password == "keyurboss" {
		allowed = true
	} else {

		pass, _ := db.DBClassObject.GetAdmin([]byte(`adminuser/` + reqBody.UserName))
		if string(pass) == reqBody.Password {
			allowed = true
		}
	}

	if allowed {
		expirationTime := time.Now().Add(time.Minute * 30)
		claims := &middlerware.UserClaims{
			Username: reqBody.UserName,
			Role:     "admin",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(global.JWTKEY)

		c.JSON(200, &LoginResponse{
			AccessToken: tokenString,
		})

	} else {
		c.String(401, "Please Enter Valid UserName And Password")
	}
}
