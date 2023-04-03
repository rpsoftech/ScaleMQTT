package routes

import (
	adminRoutes "rpsoftech/scaleMQTT/src/routes/admin"

	"github.com/gin-gonic/gin"
)

func AddRoutes(router *gin.Engine) {
	adminRoutes.AdminRoutes(router)
}
