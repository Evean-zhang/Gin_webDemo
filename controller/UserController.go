package controller

import (
	"Gin_webDemo/common"
	"Gin_webDemo/model"
	"Gin_webDemo/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

func Register(context *gin.Context) {

	DB := common.GetDB()

	//获取参数
	name := context.PostForm("name")
	phone := context.PostForm("phone")
	password := context.PostForm("password")

	//数据验证
	if len(phone) != 11 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "手机号必须为11位",
		})
		return
	}
	if len(password) < 6 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "密码至少为6位",
		})
		return
	}

	//如果名称位空，传一个随机字符串
	if len(name) == 0 {
		name = util.RandomString(10)
	}
	log.Panicln(name, phone, password)
	//判断手机号是否存在
	if PhoneExist(DB, phone) {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "用户已经注册",
		})
		return
	}
	//创建用户
	newUser := model.User{
		Name:     name,
		Phone:    phone,
		Password: password,
	}
	DB.Create(&newUser)

	//返回结果
	context.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
	})
}

func PhoneExist(db *gorm.DB, phone string) bool {
	var user model.User
	db.Where("phone = ?", phone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
