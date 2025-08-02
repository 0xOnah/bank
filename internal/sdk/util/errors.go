package util

import "github.com/gin-gonic/gin"

func ErrorResponse(err error) gin.H {
	return gin.H{"error": map[string]any{
		"message": err.Error(),
	}}
}
