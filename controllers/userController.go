package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Sophinaz/go-jwt-project/database"
	"github.com/Sophinaz/go-jwt-project/helpers"
	"github.com/Sophinaz/go-jwt-project/models"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")


func GetUsers() gin.HandlerFunc {
	return func (c *gin.Context) {
		c.IndentedJSON(http.StatusOK, gin.H{"he": "llo"})
	}
}

func GetUser() gin.HandlerFunc {
	return func (c *gin.Context) {
		user_id := c.Param("id")

		err := helpers.MatchUserTypeToUid(c, user_id)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"erroe": err.Error()})
			return
		}

		var user models.User
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = userCollection.FindOne(ctx, bson.M{"user_id": user_id}).Decode(&user)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, user)
	}
}

