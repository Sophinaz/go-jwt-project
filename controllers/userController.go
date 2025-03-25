package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"
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
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		users := []models.User{}

		err := helpers.CheckUserType(c, "ADMIN")

		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user"})
		}

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))

		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err := strconv.Atoi(c.Query("page"))

		if err != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}

		groupStage := bson.D{{"$group", bson.D{
			{"_id", bson.D{{"_id", "null"}}}, 
			{"total_count", bson.D{{"$sum", 1}}}, 
			{"data", bson.D{{"$push", "$$ROOT"}}},
		}}}

		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
			}},
		}



		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{ matchStage, groupStage, projectStage})
		defer cancel()
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error occured while list all users"})
		}

		if err = result.All(ctx, &users); err != nil{
			log.Fatal(err)
		}

		c.IndentedJSON(http.StatusOK, users)
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
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
