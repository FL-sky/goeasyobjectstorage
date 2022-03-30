package main

import (
	"log"
	"net/http"
	"os"

	"./heartbeat"
	"./locate"
	"./objects"
	"./versions"
)

func main() {
	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/versions/", versions.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}

// 相比上一章，本章的接口服务main 函数多了一个用于处理/versions/的函数，名字叫versions.Handler。
// 读者现在应该已经对这样的写法很熟悉了，明白这是versions包的Handler函数。
