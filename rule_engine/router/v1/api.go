package v1

import "github.com/gin-gonic/gin"

func Register(v1 *gin.RouterGroup) {
	//规则
	rule := v1.Group("/rule")
	{
		rule.GET("", Rule.Get)
		rule.POST("", Rule.Post)
		rule.PUT("", Rule.Put)
		rule.DELETE("", Rule.Delete)
		rule.GET("enum", Rule.Enum)
	}
	//配置
	config := v1.Group("/config")
	{
		config.POST("", Config.Post)
		config.PUT("", Config.Put)
		config.GET("", Config.Get)
		config.GET("enum", Config.Enum)
	}
	//黑白名单
	ip := v1.Group("/ip")
	{
		ip.POST("", Ip.Post)
		ip.GET("", Ip.Get)
		ip.DELETE("/remove", Ip.Remove) //任务调度检测黑名单IP是否过期,若过期则删除
		ip.DELETE("", Ip.Delete)
		ip.GET("enum", Ip.Enum)
	}
}
