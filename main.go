package main

import (
	"github.com/gin-gonic/gin"
)

func testGet(c *gin.Context) {
	c.String(200, "hi")
}



func main() {
	router := gin.Default()

	router.GET("/test", testGet)

	router.Run("localhost:8080")
}
