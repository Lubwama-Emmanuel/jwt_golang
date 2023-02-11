package controllers

import (
	"github.com/Lubwama-Emmanuel/go_jwt_mongodb/database"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")

func Hashpassword()

func VerifyPassword()

func SignUp()

func Login()

func GetUsers()

func GetUser()
