package model

import "time"

type Customer struct {
	ID           int       `gorm:"column:CustomerProfileID"`
	CustomerName string    `gorm:"column:NameCn"`
	CustID       string    `gorm:"column:CustomerID"`
	CategoryType string    `gorm:"column:CustCategoryType"`
	CustType     string    `gorm:"column:CustType"`
	Upline       string    `gorm:"column:ReferralID"`
	dsappdate    time.Time `gorm:"column:DSAppDate"`
	spqualdate   time.Time `gorm:"column:SupQualDate"`
}

func (c Customer) TableName() string {
	return "customerprofile_cls"

}
