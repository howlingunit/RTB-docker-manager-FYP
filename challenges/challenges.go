package challenges

import (
	"fmt"

	"github.com/gin-gonic/gin"
	dockerlib "github.com/howlingunit/RTB-docker-manager-FYP/dockerLib"
)

type runChallengeBody struct {
	Name string `json:"name"`
	Flag string `json:"flag"`
}

func CreateChallenges(c *gin.Context) {
	var res []dockerlib.RunChallengeRes
	var body []runChallengeBody

	if err := c.BindJSON(&body); err != nil {
		c.String(500, "invalid body")
		return
	}

	for i := range body {
		ranChallenge, err := dockerlib.RunChallenge(body[i].Name, body[i].Flag)
		if err != nil {
			c.String(500, fmt.Sprint("error creating challenge", err))
		}
		res = append(res, dockerlib.RunChallengeRes{
			Name: ranChallenge.Name,
			Flag: ranChallenge.Flag,
			Ip:   ranChallenge.Ip,
		})
	}

	c.JSON(200, res)
}

func RemoveChallenges(c *gin.Context) {
	c.String(200, "Ran Remove Challenges")
}

func GetChallenges(c *gin.Context) {
	challenges := dockerlib.ReadChallenges()

	c.IndentedJSON(200, challenges)

}
