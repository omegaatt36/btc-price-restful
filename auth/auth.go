package auth

import (
	"BTC-price-restful/models"

	jwt "github.com/dgrijalva/jwt-go"
)

// UserNameKey provises key to generate/decode token
const UserNameKey = "user_name"

var secret = []byte("secret")

// GenerateToken return a JWT by uasr name
func GenerateToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_name": user.UserName,
	})
	return token.SignedString(secret)
}
