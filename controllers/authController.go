package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Sophinaz/go-jwt-project/database"
	"github.com/Sophinaz/go-jwt-project/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
)

var validate = validator.New()

func Signup() {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User

		if err := c.BindJSON(user); err != nil {
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
		



	}
}

func login() {

}