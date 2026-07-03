package delivery

import "github.com/gin-gonic/gin"

func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"message": message})
}

func SuccessResponse(c *gin.Context, code int, data interface{}) {
	c.JSON(code, data)
}
