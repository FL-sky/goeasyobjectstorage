package main

//接口服务除了提供对象的REST接口以外还需要提供locate 功能，其main函数见

import (
	"log"
	"net/http"
	"os"

	"./heartbeat"
	"./locate"
	"./objects"
)

// 接口服务的main函数用goroutine启动了一个线程来执行 heartbeat.ListenHeartbeat函数。
// 接口服务除了需要 objects.Handler 处理URL 以/objects/开头的对象请求以外，
// 还要有一个locate.Handler函数处理URL 以/locate/开头的定位请求

func main() {
	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}

// 注意:接口服务的objects/heartbeat/locate这3个包和数据服务的同名包有很大区别。
// 数据服务的objects包负责对象在本地磁盘上的存取;而接口服务的objects,包则负责将对象请求转发给数据服务。
// 数据服务的 heartbeat包用于发送心跳消息;而接口服务的heartbeat 包则用于接收数据服务节点的心跳消息。
// 数据服务的locate包用于接收定位消息、定位对象以及发送反馈消息;而接口服务的 locate包则用于发送定位消息并处理反馈消息。
