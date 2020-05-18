package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

func main() {
	http.HandleFunc("/getOne", GetProduct)
	if err := http.ListenAndServe(":8084", nil); err != nil {
		log.Fatal("err:", err)
	}
}

func GetProduct(w http.ResponseWriter, req *http.Request) {
	if GetOneProduct() {
		w.Write([]byte("true"))
	} else {
		w.Write([]byte("false"))
	}
}

var (
	// 已下单抢购数量
	sum int64 = 0
	// 预存商品数量
	productNum int64 = 10000
	// 互斥锁
	mutex sync.Mutex
	// 限流计数
	count int64 = 0
)

// 获取秒杀商品
func GetOneProduct() bool {
	mutex.Lock()
	defer mutex.Unlock()
	count += 1
	// 限流，没100个值允许成功一个
	fmt.Printf("sum: %d  count: %d\n", sum, count)
	if count%100 == 0 {
		// 判断是否超过限制，防止超卖
		if sum < productNum {
			sum += 1
			fmt.Printf("sum: %d  count: %d\n", sum, count)
			return true
		}
	}
	return false
}
