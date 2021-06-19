package common

import (
	"github.com/dgrijalva/jwt-go"
	"silicon.com/ginessential/model"
	"time"
)

var jwtkey = []byte("a_secret_crect")

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

func ReleaseToken(user model.User)(string,error)  {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)  // 失效时间偏移一星期，一星期失效

	claims := &Claims{
		UserId: user.ID,
		StandardClaims : jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),  //何时token失效
			IssuedAt: time.Now().Unix(),   //何时发放
			Issuer: "silicon" ,//谁发放的token
			Subject: "user token",
		},

	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)   //用claims生成,用es256加密
	tokenString ,err := token.SignedString(jwtkey)
	if err != nil{
		return "", err
	}

	return tokenString,nil
}


func ParseToken(tokenString string)(*jwt.Token,*Claims,error)  {
	claims := &Claims{}
	token,err  := jwt.ParseWithClaims(tokenString,claims, func(token *jwt.Token) (interface{}, error) {
		return jwtkey,nil
	})
	return token,claims,err
}