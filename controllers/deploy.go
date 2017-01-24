package controllers

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/gin-gonic/gin"
	"encoding/base64"
	"github.com/AndBobsYourUncle/docker-backbone/containers"
	"log"
	"math/rand"
	"os"
	"os/user"
)

var noAuth = docker.AuthConfiguration{}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// DeployerPayload ...
type DeployerPayload struct {
	ID           string              `json:"id"`
	Project      string              `json:"project"`
	Registry     Registry            `json:"registry"`
	ComposeFile  string              `json:"compose_file"`
	Extra        map[string]string   `json:"extra"`
}

// Registry ...
type Registry struct {
	URL          string              `json:"url"`
	Login        string              `json:"login"`
	Password     string              `json:"password"`
}

func homeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

func projectPath(job *DeployerPayload) string {
	path := fmt.Sprintf("/tmp/projects/%s/%s", job.ID, job.Project)
	return path
}

func randStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func raiseError(e error) {
	if e != nil {
		panic(e)
	}
}

func saveComposeFile(path string, composeFile string) {
	os.MkdirAll(path, 0777)
	sDec, _ := base64.StdEncoding.DecodeString(composeFile)

	filePath := fmt.Sprintf("%s/docker-compose.yml", path)

	f, err := os.Create(filePath)
	raiseError(err)
	defer f.Close()
	_, err = f.Write(sDec)
	f.Sync()
	raiseError(err)
}

// Deploy ...
func Deploy(c *gin.Context) {
	var data *DeployerPayload
	if c.BindJSON(&data) != nil {
		c.JSON(400, gin.H{"message": "Invalid payload", "form": data})
		c.Abort()
		return
	}
	data.ID = randStringBytesRmndr(10)

	path := projectPath(data)

	saveComposeFile(path, data.ComposeFile)

	reg := data.Registry
	containers.DockerLogin(reg.URL, reg.Login, reg.Password)

	env := createEnv(data.Extra)

	containers.DockerCompose(path, env, []string{"pull"})
	containers.DockerCompose(path, env, []string{"up", "-d"})
}

func createEnv(extra map[string]string) []string {
	c := make([]string, len(extra))
	i := 0
	for k, v := range extra {
		c[i] = fmt.Sprintf("%s=%s", k, v)
		i++
	}
	return c
}
