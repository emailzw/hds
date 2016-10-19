package main

import (
	"bufio"
	"fmt"
	"hds/tree"
	"log"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var ()

func main() {
	db.LogMode(true)
	db.SetLogger(log.New(os.Stdout, "\r\n", 0))

	defer db.Close()
	//LoadData()
	cid := "10881474"

	//生成写入文件
	var filename = "./" + cid + ".txt"
	var f *os.File
	if checkFileIsExist(filename) {
		os.Remove(filename)
		fmt.Println("文件", filename, "已存在，删除文件...")
	}
	f, err = os.Create(filename)
	checkError(err)
	defer f.Close()
	buffile := bufio.NewWriter(f)

	custName, _ := redis.String(redisCon1.Do("HGET", "CUST:"+cid, "CustomerName"))
	formatMsg := fmt.Sprintf("市场查询%s[%s],", custName, cid)
	defer RuntimeCounter(formatMsg, time.Now())
	buildTreeBeginTime := time.Now()
	root := BuildTree(cid, 0)
	fmt.Printf("Buildtree 耗时:%f\n", time.Now().Sub(buildTreeBeginTime).Seconds())
	tree.TravelTree(root, func(node *tree.TreeNode) {
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
	tree.TravelTree(root, func(node *tree.TreeNode) {
		if node.Deep > maxDeepSize {
			maxDeepSize = node.Deep
			maxDeepCust = node.Value
		}
	})
	fmt.Println("maxdeepsize", maxDeepSize)
	fmt.Println("maxdeepcust", maxDeepCust)
}

func BuildTree(rootid string, deep int) *tree.TreeNode {
	node := tree.NewTreeNode(rootid)
	node.Deep = deep
	r, e := redis.Strings(redisCon0.Do("SMEMBERS", rootid))
	checkError(e)
	sort.Strings(r)
	for _, v := range r {
		child := BuildTree(v, deep+1)
		child.Parent = node
		node.Leafs = append(node.Leafs, child)
	}
	return node
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false

	}
	return exist

}

func init() {
	db, err = gorm.Open("mysql", "hds:hds@/hds?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal("mysql connecting error!", err)
		panic("redis connecting error!!!!!")
	}

	redisCon0, err = redis.DialTimeout("tcp", "127.0.0.1:6379", 0, 1*time.Second, 1*time.Second)

	if err != nil {
		log.Fatal("redis connecting error!")
		panic("redis connecting error!!!!!")
	}

	redisCon1, err = redis.DialTimeout("tcp", "127.0.0.1:6379", 0, 1*time.Second, 1*time.Second)
	if err != nil {
		log.Fatal("redis connecting error!")
		panic("redis connecting error!!!!!")
	}
	_, err := redisCon1.Do("SELECT", "1")
	if err != nil {
		log.Fatal("redis db1 changing error!")
		panic(err)
	}
}

// 计算执行时长
func RuntimeCounter(name string, start time.Time) {
	dis := time.Now().Sub(start).Seconds()
	pubstr := fmt.Sprintf("%s 执行总用时 %f seconds", name, dis)
	fmt.Println(pubstr)

}

//func GetCustomerFileds(custid string, fileds ...string) Map[string]string{

//custName, _ := redis.String(redisCon1.Do("HGET", "CUST:"+cid, "CustomerName"))
//redisCon1.Do("HGET","CUST:"+custid,"")
//}
