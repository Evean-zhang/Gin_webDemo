package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(20);not null"`
	Phone    string `gorm:"type:varchar(11);not null;unique"`
	Password string `gorm:"size:255;not null"`
}

func main() {

	db := InitDB()
	defer db.Close()

	r := gin.Default()
	r.POST("/api/register", func(context *gin.Context) {
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
		if len(password) != 6 {
			context.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "密码至少为6位",
			})
			return
		}

		//如果名称位空，传一个随机字符串
		if len(name) == 0 {
			name = RandomString(10)
		}
		log.Panicln(name, phone, password)
		//判断手机号是否存在
		if PhoneExist(db, phone) {
			context.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "用户已经注册",
			})
			return
		}
		//创建用户
		newUser := User{
			Name:     name,
			Phone:    phone,
			Password: password,
		}
		db.Create(&newUser)

		//返回结果
		context.JSON(http.StatusOK, gin.H{
			"message": "注册成功",
		})
	})
	r.Run()
}

func RandomString(n int) string {
	var letters = []byte("asdawdadasfadsfadfsadf")
	result := make([]byte, n)
	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

//开启连接池
func InitDB() *gorm.DB {
	db, err := gorm.Open("mysql", "root:123456@(127.0.0.1:3306)/sql_test?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database,err" + err.Error())
	}
	//自动创建表 自动迁移 把结构体和数据表进行对应
	db.AutoMigrate(&User{}) //引用传递（指针）

	return db
}

func PhoneExist(db *gorm.DB, phone string) bool {
	var user User
	db.Where("phone = ?", phone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
