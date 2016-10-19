package tree

import (
	"fmt"
	"hds/model"
	"testing"
)

func Test_New(t *testing.T) {
	cust := &model.Customer{CustomerName: "张三", CustID: "CN000011", CustType: "SP"}
	root := NewTreeNode(cust)
	t.Fatal(root.Value.CustID)
	fmt.Println(root.Value.CustID)
	t.Log("success")
}
