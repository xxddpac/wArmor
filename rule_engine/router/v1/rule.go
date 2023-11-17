package v1

import (
	"github.com/gin-gonic/gin"
	"rule_engine/global"
	"rule_engine/management"
	"rule_engine/model"
)

var Rule *_Rule

type _Rule struct {
}

// Get
// @Summary 获取规则列表
// @Tags 规则
// @Accept  json
// @Produce  json
// @Param keyword query string false "模糊查询"
// @Param status query int false "规则状态【0=false,1=true】"
// @Param rule_variable query int false "规则变量【枚举值】"
// @Param rule_type query int false "规则类型【枚举值】"
// @Param rule_action query int false "规则动作【枚举值】"
// @Param severity query int false "规则级别【枚举值】"
// @Param page query string false "当前页数,默认值:1"
// @Param size query string false "当前条数,默认值:10"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Router /api/v1/rule [get]
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

// Post
// @Summary 创建规则
// @Tags 规则
// @Accept  json
// @Produce  json
// @Param raw body model.Rule true "raw"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Router /api/v1/rule [post]
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

// Put
// @Summary 修改规则
// @Tags 规则
// @Accept  json
// @Produce  json
// @Param id query string true "id"
// @Param raw body model.Rule true "raw"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Router /api/v1/rule [put]
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

// Delete
// @Summary 删除规则
// @Tags 规则
// @Accept  json
// @Produce  json
// @Param id query string true "id"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Router /api/v1/rule [delete]
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

// Enum
// @Summary 获取规则枚举对应关系
// @Tags 规则
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Router /api/v1/rule/enum [get]
func (*_Rule) Enum(ctx *gin.Context) {
	var (
		g   = model.Gin{Ctx: ctx}
		res = make(map[string][]map[string]interface{})
	)
	res["rule_variable"] = append(res["rule_variable"],
		map[string]interface{}{"key": global.HttpRequestArgs.String(), "value": global.HttpRequestArgs},
		map[string]interface{}{"key": global.HttpRequestMethod.String(), "value": global.HttpRequestMethod},
		map[string]interface{}{"key": global.HttpRequestURI.String(), "value": global.HttpRequestURI},
		map[string]interface{}{"key": global.HttpRequestBody.String(), "value": global.HttpRequestBody},
		map[string]interface{}{"key": global.HttpRequestHeaders.String(), "value": global.HttpRequestHeaders},
		map[string]interface{}{"key": global.Response.String(), "value": global.Response},
	)
	res["rule_type"] = append(res["rule_type"],
		map[string]interface{}{"key": global.Xss.String(), "value": global.Xss},
		map[string]interface{}{"key": global.WebShell.String(), "value": global.WebShell},
		map[string]interface{}{"key": global.SQLInjection.String(), "value": global.SQLInjection},
		map[string]interface{}{"key": global.PathTraversal.String(), "value": global.PathTraversal},
		map[string]interface{}{"key": global.UnsafeHttpMethod.String(), "value": global.UnsafeHttpMethod},
		map[string]interface{}{"key": global.WhiteURL.String(), "value": global.WhiteURL},
		map[string]interface{}{"key": global.BlackURL.String(), "value": global.BlackURL},
		map[string]interface{}{"key": global.Spider.String(), "value": global.Spider},
		map[string]interface{}{"key": global.SensitiveInformationMonitor.String(), "value": global.SensitiveInformationMonitor},
		map[string]interface{}{"key": global.CSRF.String(), "value": global.CSRF},
		map[string]interface{}{"key": global.CommandInjection.String(), "value": global.CommandInjection},
		map[string]interface{}{"key": global.DoSDenialOfService.String(), "value": global.DoSDenialOfService},
		map[string]interface{}{"key": global.AuthenticationBypass.String(), "value": global.AuthenticationBypass},
		map[string]interface{}{"key": global.LogicFlaw.String(), "value": global.LogicFlaw},
		map[string]interface{}{"key": global.Other.String(), "value": global.Other},
	)
	res["rule_action"] = append(res["rule_action"],
		map[string]interface{}{"key": global.Deny.String(), "value": global.Deny},
		map[string]interface{}{"key": global.Redirect.String(), "value": global.Redirect},
		map[string]interface{}{"key": global.Pass.String(), "value": global.Pass},
	)
	res["severity"] = append(res["severity"],
		map[string]interface{}{"key": global.Serious.String(), "value": global.Serious},
		map[string]interface{}{"key": global.High.String(), "value": global.High},
		map[string]interface{}{"key": global.Medium.String(), "value": global.Medium},
		map[string]interface{}{"key": global.Low.String(), "value": global.Low},
		map[string]interface{}{"key": global.Info.String(), "value": global.Info},
	)
	g.Success(res)
}
