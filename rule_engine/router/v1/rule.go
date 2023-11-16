package v1

import (
	"github.com/gin-gonic/gin"
	"rule_engine/management"
	"rule_engine/model"
)

var Rule *_Rule

type _Rule struct {
}

func (*_Rule) Get(ctx *gin.Context) {
	var (
		g     = model.Gin{Ctx: ctx}
		param = model.RuleQueryFunc()
	)
	if err := ctx.ShouldBindQuery(&param); err != nil {
		g.Fail(400, err)
		return
	}
	result, err := management.ManagerRule.Get(param)
	if err != nil {
		g.Fail(400, err)
		return
	}
	g.Success(result)
}

func (*_Rule) Post(ctx *gin.Context) {
	var (
		g     = model.Gin{Ctx: ctx}
		param model.Rule
	)
	if err := ctx.ShouldBindJSON(&param); err != nil {
		g.Fail(400, err)
		return
	}
	if err := management.ManagerRule.Post(&param); err != nil {
		g.Fail(400, err)
		return
	}
	g.Success(nil)
}

func (*_Rule) Put(ctx *gin.Context) {
	var (
		g     = model.Gin{Ctx: ctx}
		param model.Rule
		query model.QueryID
	)
	if err := ctx.ShouldBindQuery(&query); err != nil {
		g.Fail(400, err)
		return
	}
	if err := ctx.ShouldBindJSON(&param); err != nil {
		g.Fail(400, err)
		return
	}
	if err := management.ManagerRule.Put(&query, &param); err != nil {
		g.Fail(400, err)
		return
	}
	g.Success(nil)
}

func (*_Rule) Delete(ctx *gin.Context) {
	var (
		g     = model.Gin{Ctx: ctx}
		query model.QueryID
	)
	if err := ctx.ShouldBindQuery(&query); err != nil {
		g.Fail(400, err)
		return
	}
	if err := management.ManagerRule.Delete(&query); err != nil {
		g.Fail(400, err)
		return
	}
	g.Success(nil)
}
