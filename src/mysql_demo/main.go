package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/moego?charset=utf8")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()
	//查询数据，指定字段名，返回sql.Rows结果集
	var moego_code, nickname string

	rows, err := db.Query("select moego_code,nickname from user where id in (?,?)", 6836,7430)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		rows.Scan(&moego_code, &nickname)
		fmt.Println(moego_code, nickname)
	}
	defer rows.Close()


}
