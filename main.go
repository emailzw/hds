package main

import (
	"bufio"
	"fmt"
	"hds/core"
	"hds/tree"
	"hds/utils"
	"os"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var counter int

func main() {
	//redisCon0 := utils.GetRedisConnection(0)
	redisCon1 := utils.GetRedisConnection(1)

	//从mysql加载数据至redis
	chanRtn := <-core.LoadData()
	fmt.Println(chanRtn)

	cid := "CN004317"

	//生成写入文件
	var filename = "./" + cid + ".txt"
	if utils.CheckFileIsExist(filename) {
		os.Remove(filename)
		fmt.Println("文件", filename, "已存在，删除文件...")
	}
	f, err := os.Create(filename)
	utils.CheckError(err)
	defer f.Close()
	buffile := bufio.NewWriter(f)

	custName, _ := redis.String(redisCon1.Do("HGET", "CUST:"+cid, "CustomerName"))
	formatMsg := fmt.Sprintf("市场查询%s[%s],", custName, cid)
	defer utils.RuntimeCounter(formatMsg, time.Now())
	buildTreeBeginTime := time.Now()
	root := core.BuildTree(cid, 0)
	fmt.Printf("Buildtree 耗时:%f\n", time.Now().Sub(buildTreeBeginTime).Seconds())
	root.TravelTree(func(node *tree.TreeNode) {
		if node.Parent == nil {
			//fmtStr := fmt.Sprintf("%s \t %s",node.Value)
			buffile.WriteString(node.Value + "\t;")
		} else {
			buffile.WriteString(node.Value + "\t" + node.Parent.Value + ",Deep " + strconv.Itoa(node.Deep) + ";")
		}
		counter++
		buffile.WriteString(strconv.Itoa(counter) + "\n")
	})
	buffile.Flush()

	maxDeepSize := 0
	var maxDeepCust string
	root.TravelTree(func(node *tree.TreeNode) {
		if node.Deep > maxDeepSize {
			maxDeepSize = node.Deep
			maxDeepCust = node.Value
		}
	})
	fmt.Println("maxdeepsize", maxDeepSize)
	fmt.Println("maxdeepcust", maxDeepCust)
}

func init() {
}

//func GetCustomerFileds(custid string, fileds ...string) Map[string]string{

//custName, _ := redis.String(redisCon1.Do("HGET", "CUST:"+cid, "CustomerName"))
//redisCon1.Do("HGET","CUST:"+custid,"")
//}
