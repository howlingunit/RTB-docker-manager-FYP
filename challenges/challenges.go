package challenges

import (
	"fmt"

	"github.com/gin-gonic/gin"
	dockerlib "github.com/howlingunit/RTB-docker-manager-FYP/dockerLib"
)

type runChallengeBody struct {
	Name      string `json:"name"`
	Team      string `json:"team"`
	Challenge string `json:"challenge"`
	Flag      string `json:"flag"`
}

func CreateChallenges(c *gin.Context) {
	var res []dockerlib.RunChallengeRes
	var body []runChallengeBody

	if err := c.BindJSON(&body); err != nil {
		c.String(500, "invalid body")
		return
	}

	for i := range body {
		ranChallenge, err := dockerlib.RunChallenge(body[i].Name, body[i].Team, body[i].Challenge, body[i].Flag)
		if err != nil {
			c.String(500, fmt.Sprint("error creating challenge", err))
			return
		}
		res = append(res, dockerlib.RunChallengeRes{
			Name: ranChallenge.Name,
			Team: ranChallenge.Team,
			Flag: ranChallenge.Flag,
			Ip:   ranChallenge.Ip,
		})
	}

	c.JSON(200, res)
}

func RemoveChallenges(c *gin.Context) {
	team := c.Param("team")
	res, err := dockerlib.RemoveContainers("Challenge", team)
	if err != nil {
		c.String(500, fmt.Sprint("Failed due to:", err))
	}
	if _, err := dockerlib.RemoveNetwork(team); err != nil {
		c.String(500, fmt.Sprint("Failed due to:", err))
	}

	c.JSON(200, res)
}

func GetChallenges(c *gin.Context) {
	challenges := dockerlib.ReadChallenges()

	c.IndentedJSON(200, challenges)

}
