package configuration

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"os"
)

// Configuration ...
type Configuration struct {
	Token string
}

// Load ...
func Load() Configuration {
	conf := Configuration{}

	viper.AutomaticEnv()
	viper.SetConfigType("yaml")
	viper.SetConfigName("default")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()

	if err != nil {
		panic(err)
	}

	conf.Token = viper.GetString("token")

	return conf
}

// Authenticate ...
func (conf *Configuration) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Auth-Token") != conf.Token {
			c.JSON(401, gin.H{})
			c.Abort()
		}
	}
}

// CheckGenerateToken ...
func CheckGenerateToken() {
	if _, err := os.Stat("./config/default.yml"); os.IsNotExist(err) {
		fmt.Println("Token not present. Generating new token...")

		f, err := os.Create("./config/default.yml")
		raiseError(err)

		defer f.Close()

		token, err := GenerateRandomString(32)
		raiseError(err)

		_, err = f.WriteString("---\ntoken: " + token)
		raiseError(err)

		f.Sync()
	}
}

func raiseError(e error) {
	if e != nil {
		panic(e)
	}
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}
