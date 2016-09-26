package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
)

var (
	db        *sql.DB
	err       error
	redisCon0 redis.Conn
	redisCon1 redis.Conn
)

func main() {
	defer db.Close()
	defer redisCon0.Close()
	defer redisCon1.Close()
	var rows *sql.Rows
	var records []map[string]string
	rows, err = db.Query("SELECT * FROM distributor ")
	checkError(err)
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var tmp_map map[string]string
	for rows.Next() {
		tmp_map = make(map[string]string)
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}
		for i, col := range values {
			if col == nil {
			} else {
				tmp_map[columns[i]] = string(col.([]byte))
			}
		}
		records = append(records, tmp_map)
	}
	if err = rows.Err(); err != nil {
		panic(err.Error())
	}

	for _, item := range records {
		custno := item["custno"]
		upline := item["upline"]
		redisCon0.Do("SADD", upline, custno)

		for k, v := range item {
			redisCon1.Do("HMSET", "CUST:"+custno, k, v)
		}
	}

	//toredis()
}

func toredis() {
	size, err := redisCon0.Do("DBSIZE")
	if err != nil {
		log.Fatal("redis operator error")
	}
	fmt.Printf("redis  size is %d \n", size)

	redisCon0.Do("SET", "name", "jerry")
	res, _ := redisCon0.Do("GET", "name")
	fmt.Printf("the key of name is %s \n", string(res.([]byte)))

}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func init() {
	db, err = sql.Open("mysql", "hds:hds@tcp(43.254.151.243:3307)/hds?charset=utf8")
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	if err = db.Ping(); err != nil {
		log.Fatal("db connection error!!!")
	}

	redisCon0, err = redis.DialTimeout("tcp", "43.254.151.243:6379", 0, 1*time.Second, 1*time.Second)

	if err != nil {
		log.Fatal("redis connecting error!")
		panic("redis connecting error!!!!!")
	}

	redisCon1, err = redis.DialTimeout("tcp", "43.254.151.243:6379", 0, 1*time.Second, 1*time.Second)
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
