package platforms

import (
	"fmt"

	"github.com/gin-gonic/gin"
	dockerlib "github.com/howlingunit/RTB-docker-manager-FYP/dockerLib"
)

type CreatePlatformResponse struct {
	User string `json:"user"`
	Ip   string `json:"ip"`
}

func CreatePlatforms(c *gin.Context) {
	var users []string
	var res []CreatePlatformResponse

	if err := c.BindJSON(&users); err != nil {
		c.String(500, "invalid body")
	}

	for i := range users {
		runPlatform, err := dockerlib.RunPlatform(users[i])
		if err != nil {
			c.String(500, fmt.Sprint("error creating platform", err))
		}

		res = append(res, CreatePlatformResponse{User: users[i], Ip: runPlatform})
	}

	c.JSON(200, res)
}

type GetPlatformResponse struct {
	User string `json:"user"`
	Ip   string `json:"ip"`
}

func GetPlatform(c *gin.Context) {
	user := c.Param("user")

	res, err := dockerlib.DockerInfo("Platform", user)
	if err != nil {
		c.String(500, fmt.Sprint("could not get platform error:", err))
		return
	}
	c.JSON(200, res)

}

func RemovePlatforms(c *gin.Context) {
	res, err := dockerlib.RemoveContainers("Platform")

	if err != nil {
		c.String(500, fmt.Sprint("Failed due to:", err))
	}

	c.JSON(200, res)
}
