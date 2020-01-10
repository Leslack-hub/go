package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/moego?charset=utf8")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()
	//查询数据，指定字段名，返回sql.Rows结果集
	var moego_code, nickname string

	idAry := []string{"1", "2", "3"}
	ids := strings.Join(idAry, "','")
	sqlRaw := fmt.Sprintf(`SELECT moego_code,nickname FROM user WHERE id IN ('%s')`, ids)
	rows, err := db.Query(sqlRaw)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		rows.Scan(&moego_code, &nickname)
		fmt.Println(moego_code, nickname)
	}
	defer rows.Close()

	rows, err = db.Query("select moego_code,nickname from user where id in (?,?)", 6836, 7430)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		rows.Scan(&moego_code, &nickname)
		fmt.Println(moego_code, nickname)
	}
	defer rows.Close()
}

func updateData(DB *sql.DB) {
	result, err := DB.Exec("UPDATE users set age=? where id=?", "30", 3)
	if err != nil {
		fmt.Printf("Insert failed,err:%v", err)
		return
	}
	rowsaffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Get RowsAffected failed,err:%v", err)
		return
	}
	fmt.Println("RowsAffected:", rowsaffected)
}

func insertData(DB *sql.DB) {
	result, err := DB.Exec("insert INTO users(name,age) values(?,?)", "YDZ", 23)
	if err != nil {
		fmt.Printf("Insert failed,err:%v", err)
		return
	}
	lastInsertID, err := result.LastInsertId() //插入数据的主键id
	if err != nil {
		fmt.Printf("Get lastInsertID failed,err:%v", err)
		return
	}
	fmt.Println("LastInsertID:", lastInsertID)
	rowsaffected, err := result.RowsAffected() //影响行数
	if err != nil {
		fmt.Printf("Get RowsAffected failed,err:%v", err)
		return
	}
	fmt.Println("RowsAffected:", rowsaffected)
}

func deleteData(DB *sql.DB) {
	result, err := DB.Exec("delete from users where id=?", 1)
	if err != nil {
		fmt.Printf("Insert failed,err:%v", err)
		return
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		fmt.Printf("Get lastInsertID failed,err:%v", err)
		return
	}
	fmt.Println("LastInsertID:", lastInsertID)
	rowsaffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Get RowsAffected failed,err:%v", err)
		return
	}
	fmt.Println("RowsAffected:", rowsaffected)
}
