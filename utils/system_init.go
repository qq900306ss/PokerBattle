package utils

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func InitMySQL() { //連接mysql
	DB, err = gorm.Open(mysql.Open("root:900306@tcp(localhost:3306)/poker?charset=utf8mb4&parseTime=True&loc=Local")) // parse time 把他轉go 的time.time 格式 loc 是當地時間
	if err != nil {
		fmt.Println("連接mysql失敗", err)
	}
}
