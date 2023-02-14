package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/Lubwama-Emmanuel/go_jwt_mongodb/controllers"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	// incomingRoutes.Use(middleware.Authenticate())
	// incomingRoutes.GET("/users", controllers.GetUsers())
	incomingRoutes.GET("/users/:user_id", controllers.GetUser())
}
