package initialize

import (
	"fmt"
	"golf-booking-go/global"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func checkErrorPanic(err error, errString string) {
	if err != nil {
		global.Logger.Error(errString, zap.Error((err)))
		panic(err)
	}
}
func InitMysql() {
	config := global.Config.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True", config.Username, config.Password, config.Host, config.Port, config.Dbname)

	db, err := gorm.Open(mysql.Open(dsn))
	checkErrorPanic(err, "InitMysql initialization error")

	global.DB = db
	SetPool()
}

func SetPool() {
	m := global.Config.Mysql
	db, _ := global.DB.DB()
	if err := db.Ping(); err != nil {
		fmt.Printf("mysql error: %s::", err)
	}
	db.SetConnMaxIdleTime(time.Duration(m.ConnMaxLifetime))
	db.SetMaxOpenConns(m.MaxOpenConns)
	db.SetConnMaxLifetime(time.Duration(m.ConnMaxLifetime))
}
