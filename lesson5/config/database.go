package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	dsn := "root:520405@tcp(127.0.0.1:3306)/school?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("数据库连接失败：", err)
	}
	err = DB.Ping()
	if err != nil {
		log.Fatal("数据库ping失败", err)
	}
	fmt.Println("数据库连接成功")
	DB.SetMaxOpenConns(10)
	DB.SetConnMaxIdleTime(5)
}
