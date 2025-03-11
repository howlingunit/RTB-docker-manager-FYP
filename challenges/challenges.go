package challenges

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type ChallengeInfo struct {
	Name       string `json:"name"`
	Difficulty string `json:"difficulty"`
}

func readChallenges() []ChallengeInfo {
	challenges := "./vulnDockers"

	dir := os.DirFS(challenges)

	var ChallengeFiles []string

	fs.WalkDir(dir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if strings.HasSuffix(path, ".json") {
			ChallengeFiles = append(ChallengeFiles, path)
		}
		return nil
	})

	var res []ChallengeInfo
	for i := 0; i < len(ChallengeFiles); i++ {
		file, err := os.Open(fmt.Sprint("./vulnDockers/", ChallengeFiles[i]))
		if err != nil {
			log.Fatal()
		}
		defer file.Close()

		var info ChallengeInfo
		if err := json.NewDecoder(file).Decode(&info); err != nil {
			log.Fatal()
		}

		res = append(res, info)

	}

	return res
}

func CreateChallenges(c *gin.Context) {
	c.String(200, "Ran Create Challenges")
}

func RemoveChallenges(c *gin.Context) {
	c.String(200, "Ran Remove Challenges")
}

func GetChallenges(c *gin.Context) {
	challenges := readChallenges()

	c.IndentedJSON(200, challenges)

}
