package helpers

import (
	"log"
	"os"
	"time"
	"github.com/Sophinaz/go-jwt-project/database"
	jwt "github.com/dgrijalva/jwt-go"
)

type signedUser struct {
	Email string
	First_name string
	Last_name string
	Uid string
	User_type string
	jwt.StandardClaims
} 

var userCollection = database.OpenCollection(database.Client, "user")


func GenerateAllTokens(email string, first_name string, last_name string, user_type string, user_id string) (string, string, error) {
	claims := &signedUser{
		Email: email,
		First_name: first_name,
		Last_name: last_name,
		Uid: user_id,
		User_type: user_type,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &signedUser{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("SECRET_KEY")))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		log.Panic(err)
	}

	return token, refreshToken, err
}