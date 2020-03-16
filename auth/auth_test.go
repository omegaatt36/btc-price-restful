package auth

import (
	"BTC-price-restful/models"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"gopkg.in/go-playground/assert.v1"
)

func TestGenerateToken(t *testing.T) {
	user := models.User{
		UserName: "UserName",
		Password: "pwd",
	}
	claisms := jwt.MapClaims{"user_name": user.UserName}
	want, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claisms).SignedString(secret)
	got, _ := GenerateToken(&user)
	assert.Equal(t, want, got)
}
