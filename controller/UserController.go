package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"silicon.com/ginessential/common"
	"silicon.com/ginessential/dto"
	"silicon.com/ginessential/model"
	"silicon.com/ginessential/response"
	"silicon.com/ginessential/util"
)

func Register(ctx *gin.Context) {
	//获取参数
	DB := common.GetDB()
	name := ctx.PostForm("name")
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")

	//数据验证
	if len(telephone) != 11 {
		response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"手机号必须为11位")
		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"})
		return
	}

	if len(password) < 6 {
		response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"密码不能少于6位")

		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码不能少于6位"})
		return
	}

	if len(name) == 0 {
		name = util.RandomString(10)
	}

	log.Println(name, telephone, password)
	//判断手机号是否存在

	if isTelephoneExist(DB, telephone) {
		response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"手机号已经被注册")
		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机已被注册"})
		return
	}

	//创建用户
	// goto: ??
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx,http.StatusUnprocessableEntity,500,nil,"加密错误")
		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 500, "msg": "加密错误"})
		return
	}

	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasedPassword),
	}

	DB.Create(&newUser)

	//返回结果
	response.Response(ctx,http.StatusOK,200,nil,"注册成功")
	//ctx.JSON(200, gin.H{"code": 200, "message": "注册成功"})
}

func Login(ctx *gin.Context) {
	DB := common.GetDB()
	//获取参数
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")

	//数据验证
	if len(telephone) != 11 {
		response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"手机号必须为11位")
		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"})
		return
	}

	if len(password) < 6 {
		response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"密码不能少于6位")
		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码不能少于6位"})
		return
	}

	//手机号是否存在
	var user model.User
	DB.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"用户不存在")
		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户不存在"})
		return
	}
	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"密码错误")
		//ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户不存在"})
		return
	}
	//发放token
	token,err := common.ReleaseToken(user)
	if err != nil{
		response.Response(ctx,http.StatusUnprocessableEntity,500,nil,"token获取失败")

		//ctx.JSON(http.StatusUnprocessableEntity,gin.H{"code":500,"msg":"token获取失败"})
		log.Printf("token generate error : %v",err)
		return
	}
	//返回结果
	response.Response(ctx,http.StatusOK,200,gin.H{"token": token},"登录成功")
	//ctx.JSON(200, gin.H{
	//	"code": 200,
	//	"msg":  "登录成功",
	//	"data": gin.H{"token": token},
	//})
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	} else {
		return false
	}
}


func Info(ctx *gin.Context){
	user,_ := ctx.Get("user")
	response.Response(ctx,http.StatusOK,200,gin.H{"user":dto.ToUserDto(user.(model.User))},"")
	//ctx.JSON(http.StatusOK,gin.H{
	//	"code":200,
	//	"data":gin.H{"user":dto.ToUserDto(user.(model.User))},
	//	"msg":"成功",
	//})
}