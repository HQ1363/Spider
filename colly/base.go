package colly

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func init() {
	fmt.Println("此目录主要用于存放go-colly相关的爬虫")
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "spider_" + defaultTableName;
	}
}
