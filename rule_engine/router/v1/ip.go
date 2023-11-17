package v1

import (
	"github.com/gin-gonic/gin"
	"rule_engine/global"
	"rule_engine/management"
	"rule_engine/model"
)

var Ip *_Ip

type _Ip struct {
}

// Post
// @Summary 创建IP黑白名单
// @Tags IP黑白名单
// @Accept  json
// @Produce  json
// @Param raw body model.Ip true "raw"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Router /api/v1/ip [post]
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

// Get
// @Summary 获取IP黑白名单列表
// @Tags IP黑白名单
// @Accept  json
// @Produce  json
// @Param keyword query string false "模糊查询"
// @Param block_type query int false "封禁类型【枚举值】"
// @Param ip_type query int false "ip类型【枚举值】"
// @Param page query string false "当前页数,默认值:1"
// @Param size query string false "当前条数,默认值:10"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Router /api/v1/ip [get]
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

// Remove
// @Summary 定时任务检测黑名单IP是否过期,若过期则删除
// @Tags IP黑白名单
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Router /api/v1/ip/remove [delete]
func (*_Ip) Remove(ctx *gin.Context) {
	var (
		g = model.Gin{Ctx: ctx}
	)
	go management.ManagerIp.Remove()
	g.Success(nil)
}

// Delete
// @Summary 删除IP黑白名单
// @Tags IP黑白名单
// @Accept  json
// @Produce  json
// @Param id query string true "id"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Router /api/v1/ip [delete]
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

// Enum
// @Summary 获取IP黑白名单枚举对应关系
// @Tags IP黑白名单
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Router /api/v1/ip/enum [get]
func (*_Ip) Enum(ctx *gin.Context) {
	var (
		g   = model.Gin{Ctx: ctx}
		res = make(map[string][]map[string]interface{})
	)
	res["block_type"] = append(res["block_type"],
		map[string]interface{}{"key": global.Permanent.String(), "value": global.Permanent},
		map[string]interface{}{"key": global.Temporary.String(), "value": global.Temporary},
	)
	res["ip_type"] = append(res["ip_type"],
		map[string]interface{}{"key": global.WhiteList.String(), "value": global.WhiteList},
		map[string]interface{}{"key": global.BlackList.String(), "value": global.BlackList},
	)
	g.Success(res)
}
