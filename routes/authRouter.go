package routes

import (
	controller "github.com/Lubwama-Emmanuel/go_jwt_mongodb/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("users/signUp", controller.SignUp)
	incomingRoutes.POST("users/logIn", controller.LogIn)
}
