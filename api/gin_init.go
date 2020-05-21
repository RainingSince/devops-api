package api

import (
	"cicd/controller"
	"cicd/middleware"
	"github.com/gin-gonic/gin"
)

func ApiInit(port string) (err error) {
	router := gin.Default()
	router.Use(middleware.Cors())

	auth := router.Group("/")
	auth.Use(middleware.Author())
	auth.POST("/login", controller.Login)
	auth.GET("/user", controller.GetUserDetail)
	auth.POST("/pipe", controller.CreatePipeline)
	auth.POST("/repos/config", controller.SaveReposConfig)
	auth.GET("/repos/config", controller.GetReposConfig)
	auth.GET("/repos/porjects", controller.GetReposProjects)
	auth.PUT("/repos/config", controller.UpdateReposConfig)
	err = router.Run(":" + port)
	return err
}
