package delivery

import "github.com/gin-gonic/gin"

func ErrorResponce(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"message": message})
}

func SuccessResponce(c *gin.Context, code int, data interface{}) {
	c.JSON(code, data)
}
