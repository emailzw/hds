package model

import (
	"fmt"
	"time"
)

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
	return "customerprofile"

}

func (c Customer) String() string {
	return fmt.Sprintf("%s-%s-%s", c.CustID, c.CustomerName, c.CustType)
}

type ByCustID []Customer

func (a ByCustID) Len() int           { return len(a) }
func (a ByCustID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCustID) Less(i, j int) bool { return a[i].CustID < a[j].CustID }
