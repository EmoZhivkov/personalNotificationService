package auth

import (
	"net/http"
	"personalNotificationService/common"
	"personalNotificationService/repositories"

	"github.com/gin-gonic/gin"
)

type AuthenticationHandler struct {
	jwtSecret    string
	userDatabase repositories.UserDatabase
}

func NewAuthenticationHandler(jwtSecret string, userDatabase repositories.UserDatabase) *AuthenticationHandler {
	return &AuthenticationHandler{
		jwtSecret:    jwtSecret,
		userDatabase: userDatabase,
	}
}

// Authenticate godoc
// @Summary Authenticate user
// @Description Authenticates the user and returns a JWT token if the credentials are valid
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param AuthenticateRequest body AuthenticateRequest true "User credentials"
// @Success 200 {object} AuthenticateResponse
// @Failure 400 {object} common.ErrorResponse "Invalid Request"
// @Failure 401 {object} common.ErrorResponse "Unauthorized or User not found"
// @Failure 500 {object} common.ErrorResponse "Internal Server Error"
// @Router /api/v1/authenticate [post]
func (a *AuthenticationHandler) Authenticate(c *gin.Context) {
	var request AuthenticateRequest
	if err := c.Bind(&request); err != nil {
		common.RespondWithError(c, http.StatusBadRequest, "Invalid Request")
		return
	}

	user, err := a.userDatabase.GetUserByUsername(request.Username)
	if err != nil {
		common.RespondWithError(c, http.StatusUnauthorized, "There is no such user")
		return
	}

	if valid := VerifyPassword(request.Password, user.PasswordHash); !valid {
		common.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	token, err := GenerateJWT(request.Username, a.jwtSecret)
	if err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	c.JSON(http.StatusOK, AuthenticateResponse{Token: token})
}
