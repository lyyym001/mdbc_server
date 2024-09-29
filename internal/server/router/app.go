package router

import (
	"github.com/gin-gonic/gin"
	"mdbc_server/internal/middlewares"
	"mdbc_server/internal/server/service"
)

func Router() *gin.Engine {

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//跨域
	r.Use(middlewares.Cors())
	/*
		注册
	*/
	r.POST("/user/register", service.UserRegister)
	/*
		登录
	*/
	r.POST("/user/login", service.UserLogin)
	/*
		鉴权
	*/
	auth := r.Group("/auth", middlewares.Auth())

	//读取用户配置数据
	auth.POST("/user/info", service.GetUserConfig)

	return r
}
