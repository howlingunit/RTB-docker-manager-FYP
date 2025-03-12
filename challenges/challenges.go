package challenges

import (
	"github.com/gin-gonic/gin"
	dockerlib "github.com/howlingunit/RTB-docker-manager-FYP/dockerLib"
)

func CreateChallenges(c *gin.Context) {
	c.String(200, "Ran Create Challenges")
}

func RemoveChallenges(c *gin.Context) {
	c.String(200, "Ran Remove Challenges")
}

func GetChallenges(c *gin.Context) {
	challenges := dockerlib.ReadChallenges()

	c.IndentedJSON(200, challenges)

}
