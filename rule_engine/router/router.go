package router

import (
	e "errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"rule_engine/model"
	v1 "rule_engine/router/v1"
)

func NewHttpRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.NoRoute(func(c *gin.Context) {
		resp := model.Gin{Ctx: c}
		resp.Fail(http.StatusNotFound, e.New("not found route"))
	})
	router.NoMethod(func(c *gin.Context) {
		resp := model.Gin{Ctx: c}
		resp.Fail(http.StatusNotFound, e.New("not found method"))
	})
	router.GET("/ping", func(c *gin.Context) {
		resp := model.Gin{Ctx: c}
		resp.Success("pong")
	})
	v1.Register(router.Group("/api/v1"))
	return router
}
