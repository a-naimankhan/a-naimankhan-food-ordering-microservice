package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorResponse(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, gin.H{"status": http.StatusText(code), "message": message})
}

func SuccessResponse(c *gin.Context, code int, data interface{}) {
	c.JSON(code, data)
}
