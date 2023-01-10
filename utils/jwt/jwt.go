package jwt

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/erfanshekari/go-talk/internal/global"
	"github.com/erfanshekari/go-talk/models"
	jwt "github.com/golang-jwt/jwt/v4"
)

func GetToken(accessToken string) (*jwt.Token, error) {
	return jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(global.GetInstance(nil).SecretKey), nil
	})
}

func GetUserFromJWT(t *jwt.Token) (*models.User, error) {
	err := errors.New("Invalid JWT Token...")
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}
	userID, ok := (claims["user_id"]).(float64)
	if !ok {
		return nil, err
	}
	return &models.User{UserID: strconv.FormatInt(int64(int(userID)), 10)}, nil
}
