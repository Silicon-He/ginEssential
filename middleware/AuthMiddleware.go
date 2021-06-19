package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"silicon.com/ginessential/common"
	"silicon.com/ginessential/model"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//获取authorization header
		tokenString := ctx.GetHeader("Authorization")

		//validate token formate
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 402, "msg": "权限不足"})
			ctx.Abort()
			return
		}
		tokenString = tokenString[7:] //因为前缀"Bearer "一共有7位，所以从第八位截取
		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "解析失败"})
			ctx.Abort()
			return
		}

		//通过验证，获取user.id，提供服务
		userId := claims.UserId
		DB := common.GetDB()
		var user model.User
		DB.First(&user, userId)

		//验证用户是否存在
		if user.ID == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort()
			return
		}

		//用户存在 将user信息写入上下文
		ctx.Set("user", user)
		ctx.Next()

	}
}
