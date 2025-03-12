package dockerlib

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

type ChallengeInfo struct {
	Name       string `json:"name"`
	Difficulty string `json:"difficulty"`
}

func ReadChallenges() []ChallengeInfo {
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

func buildDockerImage(imageTag string, dir string) (string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer cli.Close()

	ctx := context.Background()
	buildContext, err := archive.TarWithOptions(dir, &archive.TarOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to create tar archive from context: %w", err)
	}

	buildOptions := types.ImageBuildOptions{
		Tags:   []string{imageTag},
		Remove: true,
	}

	resp, err := cli.ImageBuild(ctx, buildContext, buildOptions)
	if err != nil {
		return "", fmt.Errorf("failed to build image: %w", err)
	}
	defer resp.Body.Close()

	_, err = os.Stdout.ReadFrom(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read build response: %w", err)
	}
	return fmt.Sprint("Built tag:", imageTag), nil
}

func RunChallenge(name string, tag string) (string, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:    name,
		Hostname: name,
	}, nil, nil, nil, name)
	if err != nil {
		return "", fmt.Errorf("failed to create container")
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start container")
	}

	return fmt.Sprint("running:", name), nil
}

func createCTFNetwork() (string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer cli.Close()

	ctx := context.Background()
	netName := "ctf-network"

	networks, err := cli.NetworkList(ctx, network.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to list networks: %w", err)
	}

	// Return existing network ID if already created
	for _, net := range networks {
		if net.Name == netName {
			return net.ID, nil
		}
	}

	resp, err := cli.NetworkCreate(ctx, netName, network.CreateOptions{
		Driver: "bridge",
	})
	if err != nil {
		return "", fmt.Errorf("failed to create network: %w", err)
	}

	return fmt.Sprint("Network:", resp.ID), nil
}

func InitDocker() {
	challengeData := ReadChallenges()
	challengeDir := os.DirFS("./vulnDockers")

	var challengeFolders []string

	fs.WalkDir(challengeDir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if !strings.Contains(path, "/") && !strings.Contains(path, ".") {
			challengeFolders = append(challengeFolders, path)
		}
		return nil
	})

	// now build the containers!

	for i := 0; i < len(challengeFolders); i++ {

		fmt.Println(buildDockerImage(challengeData[i].Name, fmt.Sprint("./vulnDockers/", challengeFolders[i], "/.")))
	}

	// create the network

	fmt.Println(createCTFNetwork())

}
