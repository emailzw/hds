package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func main() {
	fmt.Println(db.HasTable("distributor"))
}

func init() {
	var err error
	db, err = gorm.Open("mysql", "hds:hds@tcp(43.254.151.243:3307)/hds?charset=utf8&parseTime=True")
	if err != nil {
		panic(err)
	}

}
