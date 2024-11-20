package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"personalNotificationService/auth"
	"personalNotificationService/common"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRecoverFromPanicMiddleware(t *testing.T) {
	router := gin.New()
	router.Use(recoverFromPanicMiddleware)

	router.GET("/panic", func(c *gin.Context) {
		panic("something went wrong")
	})

	req, _ := http.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response common.ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, "Internal Server Error", response.Message)
}

func TestValidateJWTMiddleware_ValidToken(t *testing.T) {
	jwtSecret := "test_secret"

	tokenString, _ := auth.GenerateJWT("testuser", jwtSecret)

	router := gin.New()
	router.Use(validateJWTMiddleware(jwtSecret))

	router.GET("/protected", func(c *gin.Context) {
		username, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"username": username})
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "testuser")
}

func TestValidateJWTMiddleware_InvalidToken(t *testing.T) {
	jwtSecret := "test_secret"

	router := gin.New()
	router.Use(validateJWTMiddleware(jwtSecret))

	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "You should not see this"})
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+"invalid.token.string")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authorization header is invalid")
}

func TestValidateJWTMiddleware_MissingToken(t *testing.T) {
	jwtSecret := "test_secret"

	router := gin.New()
	router.Use(validateJWTMiddleware(jwtSecret))

	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "No token required"})
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "No token required")
}
