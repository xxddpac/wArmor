package v1

import (
	"github.com/gin-gonic/gin"
	"rule_engine/management"
	"rule_engine/model"
)

var Config *_Config

type _Config struct {
}

func (*_Config) Post(ctx *gin.Context) {
	var (
		g     = model.Gin{Ctx: ctx}
		param model.Config
	)
	if err := ctx.ShouldBindJSON(&param); err != nil {
		g.Fail(400, err)
		return
	}
	if err := management.ManagerConfig.Post(&param); err != nil {
		g.Fail(400, err)
		return
	}
	g.Success(nil)
}

func (*_Config) Put(ctx *gin.Context) {
	var (
		g     = model.Gin{Ctx: ctx}
		param model.Config
		query model.QueryID
	)
	if err := ctx.ShouldBindJSON(&param); err != nil {
		g.Fail(400, err)
		return
	}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		g.Fail(400, err)
		return
	}
	if err := management.ManagerConfig.Put(&query, &param); err != nil {
		g.Fail(400, err)
		return
	}
	g.Success(nil)
}

func (*_Config) Get(ctx *gin.Context) {
	var (
		g = model.Gin{Ctx: ctx}
	)
	if result, err := management.ManagerConfig.Get(); err != nil {
		g.Fail(400, err)
	} else {
		g.Success(result)
	}

}
