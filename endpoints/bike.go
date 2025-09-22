package endpoints

import "github.com/gin-gonic/gin"

func Bike(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "bike",
	})
}
