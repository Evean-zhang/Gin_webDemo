package controller

import (
	"Gin_webDemo/common"
	"Gin_webDemo/model"
	"Gin_webDemo/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
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
	//创建用户   密码加密
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "加密码错误",
		})
	}
	newUser := model.User{
		Name:     name,
		Phone:    phone,
		Password: string(hasedPassword),
	}
	DB.Create(&newUser)

	//返回结果
	context.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "注册成功",
	})
}

func Login(context *gin.Context) {
	DB := common.GetDB()

	//获取参数
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

	//判断手机号是否存在
	var user model.User
	DB.Where("phone = ?", phone).First(&user)
	if len(password) == 0 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "用户不存在",
		})
		return
	}
	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "密码错误",
		})
		return
	}

	//发放Token给前端
	token := "11"

	//返回结果
	context.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    gin.H{"token": token},
		"message": "登陆成功",
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
