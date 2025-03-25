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
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
		return "", err
	}
	return string(hashedPassword), nil
}

func validatePassword(password string, userPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(userPassword))
	result := true
	if err != nil {
		result = false
	}
	return result
}

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

		if count1+count2 > 0 {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "email or phone already exixts"})
			return
		}
		hashedPassword, err := hashPassword(*user.Password)

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "problem with the password"})
		}

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken
		user.Password = &hashedPassword

		insertionNumber, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "couldn't create user"})
			return
		}
		defer cancel()

		c.IndentedJSON(http.StatusOK, insertionNumber)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		var foundUser models.User
		err := c.BindJSON(&user)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid email or password"})
		}

		err = userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Cannot find user"})
			return
		}

		samePerson := validatePassword(*user.Password, *foundUser.Password)

		if !samePerson {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "wrong email or password"})
			return
		}

		token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, user.User_id)
		helpers.UpdateAllTokens(token, refreshToken, foundUser.User_id)
		userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)

		c.IndentedJSON(http.StatusOK, foundUser)

	}
}
