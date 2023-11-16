package v1

import (
	"github.com/gin-gonic/gin"
	"rule_engine/management"
	"rule_engine/model"
)

var Ip *_Ip

type _Ip struct {
}

func (*_Ip) Post(ctx *gin.Context) {
	var (
		g     = model.Gin{Ctx: ctx}
		param model.Ip
	)
	if err := ctx.ShouldBindJSON(&param); err != nil {
		g.Fail(400, err)
		return
	}
	if err := management.ManagerIp.Post(&param); err != nil {
		g.Fail(400, err)
		return
	}
	g.Success(nil)
}

func (*_Ip) Get(ctx *gin.Context) {
	var (
		g     = model.Gin{Ctx: ctx}
		param = model.IpQueryFunc()
	)
	if err := ctx.ShouldBindQuery(&param); err != nil {
		g.Fail(400, err)
		return
	}
	if result, err := management.ManagerIp.Get(param); err != nil {
		g.Fail(400, err)
		return
	} else {
		g.Success(result)
	}
}

func (*_Ip) Remove(ctx *gin.Context) {
	var (
		g = model.Gin{Ctx: ctx}
	)
	go management.ManagerIp.Remove()
	g.Success(nil)
}

func (*_Ip) Delete(ctx *gin.Context) {
	var (
		g     = model.Gin{Ctx: ctx}
		query model.QueryID
	)
	if err := ctx.ShouldBindQuery(&query); err != nil {
		g.Fail(400, err)
		return
	}
	if err := management.ManagerIp.Delete(query.ID); err != nil {
		g.Fail(400, err)
		return
	}
	g.Success(nil)
}
