package routes

import (
	"github.com/Lubwama-Emmanuel/go_jwt_mongodb/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("users/signUp", controllers.SignUp())
	incomingRoutes.POST("users/logIn", controllers.LogIn())
}
