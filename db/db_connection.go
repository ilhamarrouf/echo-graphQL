package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func CreateConnection() *gorm.DB {
	dbms := "mysql"
	user := "root"
	pass := "root"
	protocol := "tcp(0.0.0.0:3306)"
	database := "echo_graphql"

	connect := user+":"+pass+"@"+protocol+"/"+database

	db, err := gorm.Open(dbms, connect)
	if err != nil {
		panic(err.Error())
	}

	return db
}