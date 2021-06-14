package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"varchar(110;not null;unique"`
	Password  string `gorm:"size:255;not null"`
}

func main() {
	db := InitDB()

	//defer db.
	r := gin.Default()
	r.POST("/api/auth/register", func(ctx *gin.Context) {
		//获取参数
		name := ctx.PostForm("name")
		telephone := ctx.PostForm("telephone")
		password := ctx.PostForm("password")

		//数据验证
		if len(telephone) != 11 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"})
			return
		}

		if len(password) < 6 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码不能少于6位"})
			return
		}

		if len(name) == 0 {
			name = RandomString(10)
		}

		log.Println(name, telephone, password)
		//判断手机号是否存在

		if isTelephoneExist(db,telephone){
			ctx.JSON(http.StatusUnprocessableEntity,gin.H{"code":422,"msg":"手机已被注册"})
			return
		}else {
			newUser := User{
				Name: name,
				Telephone: telephone,
				Password: password,
			}
			db.Create(&newUser)
			ctx.JSON(200,gin.H{"msg":"注册成功"})
		}

		//创建用户

		//返回结果
		ctx.JSON(200, gin.H{
			"message": "注册成功",
		})
	})
	panic(r.Run()) // listen and serve on 0.0.0.0:8080
}

func RandomString(n int) string {
	var letters = []byte("zxcvbnmasdfghjklqwertyuiopASDFGHJKLZXCVBNMQWERTYUIOP")
	result := make([]byte, n)
	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func isTelephoneExist(db *gorm.DB,telephone string)bool  {
	var user User
	db.Where("telephone = ?",telephone).First(&user)
	if user.ID != 0{
		return true
	}else {
		return false
	}
}

func InitDB() *gorm.DB {
	host := "localhost"
	port := "3306"
	database := "ginessential"
	username := "root"
	password := "hxz000118"
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)

	print(args)
	db, err := gorm.Open(mysql.Open(args), &gorm.Config{})
	if err != nil {
		panic("failed to connect database ,err :" + err.Error())
	}
	db.AutoMigrate(&User{})
	return db
}
