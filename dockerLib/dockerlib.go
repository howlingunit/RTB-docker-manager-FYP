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

	// build platform
	fmt.Println(buildDockerImage("plat", "./platformDocker/."))

	// create the network

	fmt.Println(createCTFNetwork())

}

type RunChallengeRes struct {
	Name string
	Flag string
	Ip   string
}

func RunChallenge(name string, flag string) (RunChallengeRes, error) {
	blank := RunChallengeRes{
		Name: "",
		Flag: "",
		Ip:   "",
	}
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	labels := map[string]string{
		"type": "Challenge",
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:    name,
		Hostname: name,
		Labels:   labels,
	}, &container.HostConfig{
		NetworkMode: "ctf-network",
	}, nil, nil, name)
	if err != nil {
		return blank, fmt.Errorf("failed to create container")
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return blank, fmt.Errorf("failed to start container")
	}

	containerInfo, err := cli.ContainerInspect(context.Background(), resp.ID)
	if err != nil {
		return blank, fmt.Errorf("failed to inspect container")
	}

	IP := containerInfo.NetworkSettings.Networks["ctf-network"].IPAddress

	execID, err := cli.ContainerExecCreate(ctx, resp.ID, container.ExecOptions{
		Cmd: []string{"sh", "-c", fmt.Sprintf("echo %s > /root/flag.txt", flag)},
	})
	if err != nil {
		return blank, fmt.Errorf("failecd to make exec")
	}

	err = cli.ContainerExecStart(ctx, execID.ID, container.ExecStartOptions{
		Tty: false,
	})
	if err != nil {
		return blank, fmt.Errorf("failecd to run exec")
	}

	return RunChallengeRes{
		Name: name,
		Flag: flag,
		Ip:   IP,
	}, nil
}

func RemoveChallenges() error {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	// list docker containers

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		panic(err)
	}

	// seperate out challenges
	var challenges []container.Summary

	for i := 0; i < len(containers); i++ {
		if containers[i].Labels["type"] == "Challenge" {
			challenges = append(challenges, containers[i])
		}
	}

	// kill challenges
	for i := 0; i < len(challenges); i++ {
		if challenges[i].State == "running" {
			if err := cli.ContainerStop(ctx, challenges[i].ID, container.StopOptions{}); err != nil {
				return err
			}
		}
	}
	// rm challenges

	for i := 0; i < len(challenges); i++ {
		if err := cli.ContainerRemove(ctx, challenges[i].ID, container.RemoveOptions{}); err != nil {
			return err
		}
	}
	return nil
}

func RunPlatform(user string) (string, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	labels := map[string]string{
		"type": "platform",
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:    "plat",
		Hostname: user,
		Labels:   labels,
	}, &container.HostConfig{
		NetworkMode: "ctf-network",
	}, nil, nil, user)
	if err != nil {
		return "", fmt.Errorf("failed to create container")
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start container")
	}

	containerInfo, err := cli.ContainerInspect(context.Background(), resp.ID)
	if err != nil {
		return "", fmt.Errorf("failed to inspect container")
	}

	IP := containerInfo.NetworkSettings.Networks["ctf-network"].IPAddress

	return IP, nil
}
