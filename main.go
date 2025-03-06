package main

import (
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
)

func testGet(c *gin.Context) {
	c.String(200, "hi")
}

func main() {
	inter := flag.String("interface", "localhost", "The interface for the api")

	flag.Parse()

	router := gin.Default()

	router.GET("/test", testGet)

	router.Run(fmt.Sprintf("%s:8080", *inter))
}
