package main

import (
	"log"
	"net/http"
	"os"

	"./heartbeat"
	"./locate"
	"./objects"
	"./temp"
	"./versions"
)

// 接口服务的main 函数以及objects包发生了改变，且新增了temp包，versions/locate/heartbeat包没有变化。

func main() {
	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/versions/", versions.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}

// 相比第3章，main函数多了一个 temp.Handler 函数用于处理对/temp/的请求。
// 在深入temp包的实现之前,让我们先去看看objects包发生的改动。
