package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func GetRedisConnection(db int) redis.Conn {

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

func GetMainDBConnection() *gorm.DB {
	db, err := gorm.Open("mysql", "hds:hds@/hds?charset=utf8&parseTime=True&loc=Local")
	CheckError(err)
	return db
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false

	}
	return exist

}

// 计算执行时长
func RuntimeCounter(name string, start time.Time) {
	dis := time.Now().Sub(start).Seconds()
	pubstr := fmt.Sprintf("%s 执行总用时 %f seconds", name, dis)
	fmt.Println(pubstr)

}
