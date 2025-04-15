package main

import (
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/howlingunit/RTB-docker-manager-FYP/challenges"
	dockerlib "github.com/howlingunit/RTB-docker-manager-FYP/dockerLib"
	"github.com/howlingunit/RTB-docker-manager-FYP/platforms.go"
)

func testGet(c *gin.Context) {
	c.String(200, "hi")
}

func main() {
	inter := flag.String("interface", "localhost", "The interface for the api")
	flag.Parse()

	dockerlib.InitDocker()

	router := gin.Default()

	router.GET("/test", testGet)
	router.GET("/get-challenges", challenges.GetChallenges)
	router.GET("/get-platform/:user", platforms.GetPlatform)
	router.POST("/create-challenges", challenges.CreateChallenges)
	router.POST("/create-platforms", platforms.CreatePlatforms)
	router.DELETE("/remove-challenges/:team", challenges.RemoveChallenges)
	router.DELETE("/remove-platforms", platforms.RemovePlatforms)

	router.Run(fmt.Sprintf("%s:8080", *inter))
}
