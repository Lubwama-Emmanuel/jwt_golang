package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Lubwama-Emmanuel/go_jwt_mongodb/database"
	"github.com/Lubwama-Emmanuel/go_jwt_mongodb/helpers"
	"github.com/Lubwama-Emmanuel/go_jwt_mongodb/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")
var validate = validator.New()

func Hashpassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 20)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintln("Email or password is incorrect")
		check = false
	}
	return check, msg

}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": "An error occured at binding",
			})
			return
		}
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": validationErr,
			})
			return
		}

		// count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		// defer cancel()
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"Error": "error occured while checking for email address"})
		// 	log.Panic()
		// }

		// count, err := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		// defer cancel()
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"Error": "error occured while checking for phone number"})
		// 	log.Panic()
		// }
		hashedPassword := Hashpassword(user.Password)
		user.Password = hashedPassword

		// if count > 0 {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"Error": "Email or Phonenumber already exists, try logging in"})
		// }

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Id = primitive.NewObjectID()
		user.User_id = user.Id.Hex()
		token, refreshToken, _ := helpers.GenerateAllTokens(user.Email, user.First_name, user.Last_name, user.User_id, user.User_type)
		user.Token = token
		user.Refresh_token = refreshToken
		resultInsertNumber, insertErr := userCollection.InsertOne(ctx, user)

		if insertErr != nil {
			msg, _ := fmt.Println("User item not inserted")
			c.JSON(http.StatusInternalServerError, gin.H{"Error": msg})
			return
		}
		c.JSON(http.StatusOK, resultInsertNumber)
	}
}

func LogIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": err.Error(),
			})
			defer cancel()
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": "email or password is incorrect",
			})
			return
		}
		passwordIsValid, msg := VerifyPassword(foundUser.Password, user.Password)

		if !passwordIsValid {
			c.JSON(http.StatusBadRequest, gin.H{"Error": msg})
			return
		}

		token, refreshToken, _ := helpers.GenerateAllTokens(foundUser.Email, foundUser.First_name, foundUser.First_name, foundUser.User_id, foundUser.User_type)

		helpers.UpdateAllTokens(token, refreshToken, foundUser.User_id)
	}
}

// func GetUsers() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
// 			return
// 		}
// 		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
// 	}
// }

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		if err := helpers.MatchTypeOfUser(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}
