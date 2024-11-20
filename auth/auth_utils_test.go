package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	jwtSecret := "test_secret"
	username := "testuser"

	tokenString, err := GenerateJWT(username, jwtSecret)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	token, ok := ValidateJWT(tokenString, jwtSecret)
	assert.True(t, ok)
	assert.True(t, token.Valid)

	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, username, claims["username"])
	assert.WithinDuration(t, time.Now().Add(time.Hour*72), time.Unix(int64(claims["exp"].(float64)), 0), time.Minute)
}

func TestValidateJWT(t *testing.T) {
	jwtSecret := "test_secret"
	username := "testuser"

	tokenString, _ := GenerateJWT(username, jwtSecret)

	token, isValid := ValidateJWT(tokenString, jwtSecret)
	assert.True(t, isValid)
	assert.NotNil(t, token)

	token, isValid = ValidateJWT(tokenString, "wrong_secret")
	assert.False(t, isValid)
	assert.Nil(t, token)

	token, isValid = ValidateJWT("invalid.token.string", jwtSecret)
	assert.False(t, isValid)
	assert.Nil(t, token)
}

func TestHashPassword(t *testing.T) {
	password := "password123"

	hashedPassword, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)

	assert.True(t, VerifyPassword(password, hashedPassword))
}

func TestVerifyPassword(t *testing.T) {
	password := "password123"

	hashedPassword, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)

	isValid := VerifyPassword(password, hashedPassword)
	assert.True(t, isValid)

	isValid = VerifyPassword("wrongpassword", hashedPassword)
	assert.False(t, isValid)
}
