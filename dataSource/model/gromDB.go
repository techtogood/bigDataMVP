package model



import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"dataSource/configs"
)

func GormOpenDB() (db *gorm.DB, err error) {
	fmt.Println("GormOpenDB...",configs.Config.MysqlDB.Dsn())
	db, err = gorm.Open("mysql", configs.Config.MysqlDB.Dsn())
	return
}