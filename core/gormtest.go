package main

import (
	"fmt"
	"hds/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var ()

func main() {
	db, err := gorm.Open("mysql", "hds:hds@/hds?charset=utf8&parseTime=True&loc=Local")
	checkError(err)
	cust := &model.Customer{}
	custs := []model.Customer{}
	defer db.Close()
	db.First(&cust)
	db.Find(&custs)
	fmt.Println(cust)
	fmt.Println(len(custs))
}

func init() {

}
