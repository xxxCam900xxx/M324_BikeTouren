package main

import (
	"biketouren/endpoints"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/bike", endpoints.Bike)
	router.GET("/tour", endpoints.Tour)
	router.Run() // listens on 0.0.0.0:8080 by default
}
