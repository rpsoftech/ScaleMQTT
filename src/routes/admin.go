package routes

import (
	"net/http"

	dbPackage "rpsoftech/scaleMQTT/src/db"

	"github.com/gin-gonic/gin"
)

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type AddNewAdminLogin struct {
	UserName string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

func AdminRoutes(router *gin.Engine) {
	router.GET("/admin/login", AdminLoginFunction)

	adminRouter := router.Group("/admin")
	{

		// adminRouter.
		adminRouter.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "Welcome Gin Server")
		})
	}
}

func AdminLoginFunction(c *gin.Context) {
	// c.String(http.StatusOK, "Welcome Gin Server")
	reqBody := AddNewAdminLogin{}

	if errA := c.ShouldBindQuery(&reqBody); errA != nil {
		c.String(400, `Invalid JSON input`)
		return
		// At this time, it reuses body stored in the context.
	}
	allowed := false
	if reqBody.UserName == "keyurboss" && reqBody.Password == "keyurboss" {
		allowed = true
	}
	pass, _ := dbPackage.DbConnection.Get([]byte(`adminuser/` + reqBody.UserName))
	if string(pass) == reqBody.Password {
		allowed = true
	}

	if allowed {
		c.JSON(200, map[string]interface{}{
			"AccessToken":  string(pass),
			"RefreshToken": "SASasaSAaaa",
		})

	} else {
		c.String(401, "Please Enter Valid UserName And Password")
	}
}
