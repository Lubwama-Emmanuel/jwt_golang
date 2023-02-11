package routes

import (
	"github.com/Lubwama-Emmanuel/go_jwt_mongodb/middleware"
	"github.com/gin-gonic/gin"

	controller "github.com/Lubwama-Emmanuel/go_jwt_mongodb/controllers"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/users/:user_id", controller.GetUser())

}
