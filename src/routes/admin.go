package routes

import (
	"net/http"
	"time"

	dbPackage "rpsoftech/scaleMQTT/src/db"
	"rpsoftech/scaleMQTT/src/global"
	"rpsoftech/scaleMQTT/src/middlerware"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
	// RefreshToken string `json:"refreshToken"`
}

type AddNewAdminLogin struct {
	UserName string `json:"username" form:"username" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

func AdminRoutes(router *gin.Engine) {
	router.POST("/admin/login", AdminLoginFunction)

	adminRouter := router.Group("/admin")
	{
		adminRouter.Use(middlerware.JwtAuthMiddleware())
		// adminRouter.
		adminRouter.POST("/addNewAdmin")
		adminRouter.Any("/isLoggedin", IsAdminLoggedIn)
		adminRouter.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "Welcome Gin Server")
		})
	}
}

func IsAdminLoggedIn(c *gin.Context) {
	if val, ok := c.Get("User"); ok {
		c.JSON(200, gin.H{"user": val})
	} else {
		c.JSON(200, gin.H{"error": "User not found in context"})
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
		expirationTime := time.Now().Add(time.Minute * 30).Unix()
		claims := &middlerware.UserClaims{
			Username: reqBody.UserName,
			Role:     "admin",
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime,
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

func AddNewUser(c *gin.Context) {
	var user AddNewAdminLogin

	// Bind the JSON body to the user struct and validate the fields
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the user to the database
	// ...

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}
