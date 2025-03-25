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
