package main

import (
	"Gin_webDemo/common"
	"Gin_webDemo/routes"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	db := common.InitDB()
	defer db.Close()

	r := gin.Default()
	r = routes.CollectRoute(r)
	r.Run()
}
