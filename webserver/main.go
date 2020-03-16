package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})
	router.POST("/account/new", NewAccount)
	router.POST("/account/sign", Sign)
	router.Run(":3000")
}
