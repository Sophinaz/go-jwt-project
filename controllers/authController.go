package controllers

import (
	"context"
	"log"
	"net/http"
	"time"
	"github.com/Sophinaz/go-jwt-project/helpers"
	"github.com/Sophinaz/go-jwt-project/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate = validator.New()

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if validationErr := validate.Struct(user); validationErr != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"validation Error": validationErr.Error()})
			return
		}

		count1, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error while checking the email"})
			return
		}

		count2, err := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error while checking the phone"})
			return
		}

		if count1 + count2 > 0 {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "email or phone already exixts"})
			return
		}
		
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, *&user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken

		insertionNumber, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "couldn't create user"})
			return
		}
		defer cancel()

		c.IndentedJSON(http.StatusOK, insertionNumber)
	}
}

func Login() gin.HandlerFunc{
	return func(c *gin.Context) {

	}
}