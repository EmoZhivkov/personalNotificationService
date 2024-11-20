package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"personalNotificationService/repositories"
	"personalNotificationService/repositories/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticationHandler_Authenticate_Success(t *testing.T) {
	jwtSecret := "test_secret"
	password := "password123"
	hashedPassword, _ := HashPassword(password)
	user := &repositories.User{
		Username:     "testuser",
		PasswordHash: hashedPassword,
	}

	mockUserDB := mocks.NewMockUserDatabase(t)
	mockUserDB.On("GetUserByUsername", user.Username).Return(user, nil)

	requestBody, _ := json.Marshal(AuthenticateRequest{
		Username: user.Username,
		Password: password,
	})

	req, _ := http.NewRequest("POST", "/authenticate", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler := NewAuthenticationHandler(jwtSecret, mockUserDB)
	handler.Authenticate(c)

	var response AuthenticateResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, response.Token)

	token, isValid := ValidateJWT(response.Token, jwtSecret)
	assert.True(t, isValid)

	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Contains(t, claims, "username")
	assert.Equal(t, user.Username, claims["username"])
}

func TestAuthenticationHandler_Authenticate_InvalidPassword(t *testing.T) {
	password := "password123"
	hashedPassword, _ := HashPassword(password)
	user := &repositories.User{
		Username:     "testuser",
		PasswordHash: hashedPassword,
	}

	mockUserDB := mocks.NewMockUserDatabase(t)
	mockUserDB.On("GetUserByUsername", user.Username).Return(user, nil)

	requestBody, _ := json.Marshal(AuthenticateRequest{
		Username: user.Username,
		Password: "wrongpassword",
	})

	req, _ := http.NewRequest("POST", "/authenticate", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler := NewAuthenticationHandler("test_secret", mockUserDB)
	handler.Authenticate(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthenticationHandler_Authenticate_UserNotFound(t *testing.T) {
	nonExistingUsername := "nonexistentuser"

	mockUserDB := mocks.NewMockUserDatabase(t)
	mockUserDB.On("GetUserByUsername", nonExistingUsername).Return(nil, errors.New("user not found"))

	requestBody, _ := json.Marshal(AuthenticateRequest{
		Username: nonExistingUsername,
		Password: "password123",
	})

	req, _ := http.NewRequest("POST", "/authenticate", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler := NewAuthenticationHandler("test_secret", mockUserDB)
	handler.Authenticate(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockUserDB.AssertExpectations(t)
}
