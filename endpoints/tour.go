package endpoints

import "github.com/gin-gonic/gin"

func Tour(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "tour",
	})
}
