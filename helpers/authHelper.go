package helpers

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func checkUserType(c *gin.Context, user_type string) error {
	role := c.GetString("user_type")
	if user_type != role {
		err := errors.New("Unauthorized user")
		return err
	}
	return nil
}

func MatchUserTypeToUid(c *gin.Context, user_id string) error {
	user_type := c.GetString("user_type")
	uid := c.GetString("uid")

	if user_type == "USER" && user_id != uid {
		err := errors.New("Unauthorized user")
		return err
	}
	
	err := checkUserType(c, user_type)

	return err
}