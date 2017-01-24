package routing

import (
	"github.com/AndBobsYourUncle/docker-backbone/configuration"
	"github.com/AndBobsYourUncle/docker-backbone/controllers"
	"github.com/gin-gonic/gin"
)

// GetRouter ...
func GetRouter(conf configuration.Configuration) *gin.Engine {
	router := gin.Default()

	router.Use(conf.Authenticate())

	router.POST("/deploy", func(c *gin.Context) {
		controllers.Deploy(c)
	})

	return router
}
