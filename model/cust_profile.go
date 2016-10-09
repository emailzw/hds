package model

import "time"

type Customer struct {
	ID           int       `gorm:"column:customerprofileid"`
	CustomerName string    `gorm:"column:namecn"`
	CustID       string    `gorm:"column:customerid"`
	CategoryType string    `gorm:"column:custcategorytype"`
	CustType     string    `gorm:"column:custtype"`
	Upline       string    `gorm:"column:referralid"`
	dsappdate    time.Time `gorm:"column:dsappdate"`
	spqualdate   time.Time `gorm:"column:supqualdate"`
}

func (c Customer) TableName() string {
	return "customerprofile_cls"

}
