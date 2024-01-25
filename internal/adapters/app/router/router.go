package router

import (
	"github.com/Kosodaka/enricher-service/pkg/config"
	"github.com/gin-gonic/gin"
)

type personRouter interface {
	AddPerson(c *gin.Context)
	GetPerson(c *gin.Context)
	UpdatePerson(c *gin.Context)
	DeletePerson(c *gin.Context)
	GetPersons(c *gin.Context)
}

type Router struct {
	Server       *gin.Engine
	Port         string
	Env          string
	PersonRouter personRouter
}

func (r *Router) InitRoutes() {
	r.Server.POST("/persons", r.PersonRouter.AddPerson)
	r.Server.GET("/person/:id", r.PersonRouter.GetPerson)
	r.Server.GET("/persons", r.PersonRouter.GetPersons)
	r.Server.PATCH("/person", r.PersonRouter.UpdatePerson)
	r.Server.DELETE("/person", r.PersonRouter.DeletePerson)

}

func (r *Router) Run() error {
	return r.Server.Run(":" + r.Port)
}

func NewRouter(cfg config.Config, p personRouter) *Router {
	router := &Router{
		PersonRouter: p,
		Port:         cfg.HttpPort,
	}

	router.Server = gin.Default()
	router.Server.Use(gin.Recovery())

	router.InitRoutes()

	return router
}
