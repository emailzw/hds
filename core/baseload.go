package core

import (
	"fmt"
	"hds/model"
	"hds/utils"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	counter   int
	jobStatus chan string
)

func LoadData() chan string {
	db := utils.GetMainDBConnection()
	db.LogMode(false)
	db.SetLogger(log.New(os.Stdout, "\r\n", 0))

	defer db.Close()

	custs := []model.Customer{}
	fmt.Println("loading data from mysql...")
	db.Limit(1000).Find(&custs)

	fmt.Printf("total load items:%d\n", len(custs))
	ct := len(custs)

	//5个线程并发
	concurr_factor := 5
	jobStatus = make(chan string)
	fact := ct / concurr_factor
	for i := 0; i < concurr_factor; i++ {
		if i+1 < concurr_factor {
			go toRedis(i*fact, i*fact+fact, custs)
		} else {
			go toRedis(i*fact, ct, custs)
		}
	}
	for i := 0; i < concurr_factor; i++ {
		select {
		case msg := <-jobStatus:
			fmt.Println(msg)
		}
	}
	//fmt.Println("total finished,counter", counter)
	rtn_chan := make(chan string, 1)
	rtn_chan <- fmt.Sprintf("total load %d records into redis", counter)
	return rtn_chan
}

func toRedis(from int, to int, queue []model.Customer) {
	redisCon0 := utils.GetRedisConnection(0)
	redisCon1 := utils.GetRedisConnection(1)
	defer redisCon1.Close()
	defer redisCon0.Close()
	for i := from; i < to; i++ {
		counter++
		item := queue[i]
		redisCon0.Do("SADD", item.Upline, item.CustID)
		_, err := redisCon1.Do("HMSET", "CUST:"+item.CustID, "ID", item.CustID)
		if err != nil {
			fmt.Println(err)
		}
		_, err = redisCon1.Do("HMSET", "CUST:"+item.CustID, "CustID", item.CustID)
		if err != nil {
			fmt.Println(err)
		}
		_, err = redisCon1.Do("HMSET", "CUST:"+item.CustID, "CustomerName", item.CustomerName)
		if err != nil {
			fmt.Println(err)
		}
		_, err = redisCon1.Do("HMSET", "CUST:"+item.CustID, "CustType", item.CustType)
		if err != nil {
			fmt.Println(err)
		}
	}
	jobStatus <- fmt.Sprintf("job:%d-%d finished", from, to)
}
