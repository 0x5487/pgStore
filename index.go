package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	build_data()

	router := gin.Default()

	EnableApi(router)

	router.GET("/", func(c *gin.Context) {
		c.String(200, "hello world")
	})

	router.NoRoute(func(c *gin.Context) {
		logInfo("Not Found")
		c.JSON(404, gin.H{"message": "Not Found"})
	})

	// Listen and server on 0.0.0.0:8080
	router.Run(":8080")
}
