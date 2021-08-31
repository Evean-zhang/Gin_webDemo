package common

import (
	"Gin_webDemo/model"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

//开启连接池
func InitDB() *gorm.DB {
	db, err := gorm.Open("mysql", "root:123456@(127.0.0.1:3306)/gin_webDemo?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database,err" + err.Error())
	}
	//自动创建表 自动迁移 把结构体和数据表进行对应
	db.AutoMigrate(&model.User{}) //引用传递（指针）
	return DB
}

func GetDB() *gorm.DB {
	return DB
}
