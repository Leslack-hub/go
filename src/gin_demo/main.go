package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:type:varchar(30);not null`
	Phone    string `gorm:type:varchar(11);notnull;unique`
	PassWard string `gorm:size:255;not null`
}

func main() {
	r := gin.Default()
	r.POST("/api/auth/register", func(ctx *gin.Context) {
		name := ctx.PostForm("name")
		phone := ctx.PostForm("phone")
		password := ctx.PostForm("passward")
		if !name || !phone || !password {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "参数错误"})
		}

		ctx.JSON(200, gin.H{
			"message": "nihao",
		})
	})

	r.Run()
}

func InitDB() *gorm.DB {
	driveName := "mysql"
	host := "localhost"
	port := "3306"
	database := "moego"
	username := "root"
	password := "root"
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=ture",
		username, password, host, port, database, charset)
	db, err := gorm.Open(driveName, args)
	if err != nil {
		log.Println("数据库链接失败" + err.Error())
	}
	return db
}
