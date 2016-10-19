package main

import (
	"fmt"
	"hds/model"
	"log"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	db        *gorm.DB
	err       error
	redisCon0 redis.Conn
	redisCon1 redis.Conn
	counter   int
	jobStatus chan string
)

func LoadData() {

	db.LogMode(true)
	db.SetLogger(log.New(os.Stdout, "\r\n", 0))

	defer db.Close()

	custs := []model.Customer{}
	db.Debug().Limit(10000000).Find(&custs)

	fmt.Println(len(custs))
	fmt.Println(custs[100])
	ct := len(custs)
	/*ct = 10000*/
	//5个线程并发
	concurr_factor := 100
	jobStatus = make(chan string)
	fact := ct / concurr_factor
	for i := 0; i < concurr_factor; i++ {
		if i+1 < concurr_factor {
			go func(i2 int, fact2 int, custs2 []model.Customer, ct2 int) {
				toRedis(i2*fact2, i2*fact2+fact2, custs2)
			}(i, fact, custs, ct)
		} else {
			go func(i int, fact int, custs []model.Customer, ct int) {
				toRedis(i*fact, ct, custs)
			}(i, fact, custs, ct)
		}
	}
	for i := 0; i < concurr_factor; i++ {
		select {
		case msg := <-jobStatus:
			fmt.Println(msg)
		}
	}
	fmt.Println("total finished,counter", counter)
}

func toRedis(from int, to int, queue []model.Customer) {
	redisCon0 := getRedisConnection(0)
	redisCon1 := getRedisConnection(1)
	defer redisCon1.Close()
	defer redisCon0.Close()
	for i := from; i < to; i++ {
		counter++
		item := queue[i]
		/*fmt.Println(i)*/
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
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
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

func getRedisConnection(db int) redis.Conn {

	redisCon, err := redis.DialTimeout("tcp", "127.0.0.1:6379", 0, 1*time.Second, 1*time.Second)

	if err != nil {
		log.Fatal("redis connecting error!")
		panic("redis connecting error!!!!!")
	}

	_, err = redisCon.Do("SELECT", db)
	if err != nil {
		log.Fatal("redis db changing error!", db)
		panic(err)
	}
	return redisCon
}
