package api

import (
	"log"
	"net/http"
	"personalNotificationService/auth"
	"personalNotificationService/common"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func recoverFromPanicMiddleware(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic recovered: %v", err)
			common.RespondWithError(c, http.StatusInternalServerError, "Internal Server Error")
		}
	}()
	c.Next()
}

func validateJWTMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			common.RespondWithError(c, http.StatusUnauthorized, "Authorization header format must be Bearer <token>")
			return
		}

		tokenString := parts[1]

		token, valid := auth.ValidateJWT(tokenString, jwtSecret)
		if !valid {
			common.RespondWithError(c, http.StatusUnauthorized, "Authorization header is invalid")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			common.RespondWithError(c, http.StatusUnauthorized, "Authorization header is invalid")
			return
		}

		c.Set("username", claims["username"])
		c.Next()
	}
}
