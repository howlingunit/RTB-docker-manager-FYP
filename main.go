package main

import (
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/howlingunit/RTB-docker-manager-FYP/challenges"
)

func testGet(c *gin.Context) {
	c.String(200, "hi")
}

func main() {
	inter := flag.String("interface", "localhost", "The interface for the api")
	flag.Parse()

	router := gin.Default()

	router.GET("/test", testGet)
	router.GET("/get-challenges", challenges.GetChallenges)

	router.Run(fmt.Sprintf("%s:8080", *inter))
}
