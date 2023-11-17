package v1

import (
	"github.com/gin-gonic/gin"
	"rule_engine/global"
	"rule_engine/management"
	"rule_engine/model"
)

var Config *_Config

type _Config struct {
}

// Post
// @Summary 创建waf运行模式
// @Tags 配置
// @Accept  json
// @Produce  json
// @Param raw body model.Config true "raw"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Router /api/v1/config [post]
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

// Put
// @Summary 修改waf运行模式
// @Tags 配置
// @Accept  json
// @Produce  json
// @Param id query string true "id"
// @Param raw body model.Config true "raw"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Router /api/v1/config [put]
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

// Get
// @Summary 获取waf当前运行模式
// @Tags 配置
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Router /api/v1/config [get]
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

// Enum
// @Summary 获取waf模式枚举对应关系
// @Tags 配置
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Router /api/v1/config/enum [get]
func (*_Config) Enum(ctx *gin.Context) {
	var (
		g   = model.Gin{Ctx: ctx}
		res = make(map[string][]map[string]interface{})
	)
	res["mode"] = append(res["mode"],
		map[string]interface{}{"key": global.Block.String(), "value": global.Block},
		map[string]interface{}{"key": global.Alert.String(), "value": global.Alert},
		map[string]interface{}{"key": global.Bypass.String(), "value": global.Bypass},
	)
	g.Success(res)

}
