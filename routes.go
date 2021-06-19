package main

import (
	"github.com/gin-gonic/gin"
	"silicon.com/ginessential/controller"
	"silicon.com/ginessential/middleware"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.POST("/api/auth/register",controller.Register)
	r.POST("/api/auth/login",controller.Login)
	r.GET("/api/auth/info",middleware.AuthMiddleware(),controller.Info)
	return r
}


