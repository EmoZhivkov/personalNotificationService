package common

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func RespondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, ErrorResponse{Code: code, Message: message})
	c.Abort()
}
