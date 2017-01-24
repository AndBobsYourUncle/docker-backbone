package containers

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"log"
	"os/user"
)

func homeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

func dockerClient() *docker.Client {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)

	return client
}

func pullImage(repo string, tag string) {
	auth := docker.AuthConfiguration{}

	fmt.Println("Pulling docker image", repo, tag)

	client := dockerClient()
	options := docker.PullImageOptions{
		Repository: repo,
		Tag:        tag,
	}
	err := client.PullImage(options, auth)

	if err != nil {
		fmt.Println(err.Error())
	}
}

func waitForContainer(container string) {
	client := dockerClient()
	code, err := client.WaitContainer(container)

	if err != nil {
		fmt.Println("Container waiting error:", err.Error())
	} else {
		fmt.Println("Container finished:", code)
		if code == 0 {
			fmt.Println("Removing container:", container)
			removeContainer(container)
		} else {
			fmt.Println("Container exited with error:", container)
		}
	}
}

func dockerRunContainer(conf docker.Config, hostConf docker.HostConfig) string {
	client := dockerClient()

	opts := docker.CreateContainerOptions{Config: &conf, HostConfig: &hostConf}
	cont, err := client.CreateContainer(opts)

	if err != nil {
		fmt.Println("Error creating container:", err.Error())
	} else {
		fmt.Println("Container created: ", cont.ID)
	}

	err = client.StartContainer(cont.ID, &docker.HostConfig{})

	return cont.ID
}

func removeContainer(container string) {
	client := dockerClient()

	opts := docker.RemoveContainerOptions{ID: container}
	err := client.RemoveContainer(opts)

	if err != nil {
		fmt.Println("Error removing container:", err.Error())
	}
}

// DockerLogin ...
func DockerLogin(url string, user string, pass string) {
	pullImage("docker", "1.12.0")

	contID := dockerRunContainer(
		docker.Config{
			Image: "docker:1.12.0",
			Cmd: []string{
				"login",
				"-u", user,
				"-p", pass,
				url,
			},
		},
		docker.HostConfig{
			Binds: []string{
				"/var/run/docker.sock:/var/run/docker.sock",
				homeDir() + "/.docker:/root/.docker",
			},
		},
	)

	waitForContainer(contID)
}

// DockerCompose ...
func DockerCompose(composePath string, extra []string, command []string) {
	pullImage("docker/compose", "1.8.0")

	contID := dockerRunContainer(
		docker.Config{
			Image: "docker/compose:1.8.0",
			Cmd:   command,
			Volumes: map[string]struct{}{
				"/tmp/projects:/tmp/projects": {},
			},
			Env:        extra,
			WorkingDir: composePath,
		},
		docker.HostConfig{
			Binds: []string{
				"/tmp/projects:/tmp/projects",
				"/var/run/docker.sock:/var/run/docker.sock",
				homeDir() + "/.docker:/root/.docker",
			},
		},
	)

	waitForContainer(contID)
}
