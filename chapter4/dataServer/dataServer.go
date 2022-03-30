package main

import (
	"log"
	"net/http"
	"os"

	"./heartbeat"
	"./locate"
	"./objects"
	"./temp"
)

// 和第⒉章相比我们的main函数多了一个locate.CollectObjects的函数调用并引入
// temp.Handler处理函数的注册。

func main() {
	locate.CollectObjects()
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}

// 数据服务的 locate包是用来对节点本地磁盘上的对象进行定位的。
// 在第⒉章，我们的定位通过调用os.Stat检查对象文件是否存在来实现。
// 这样的实现意味着每次定位请求都会导致一次磁盘访问。这会对整个系统带来很大的负担。
// 别忘了我们不止在PUT去重的时候需要进行一次定位，GET的时候也一样要做，
// 可以说定位是对象存储系统最频繁的操作。

// 为了减少对磁盘访问的次数，从而提高磁盘的性能，
// 本章的数据服务定位功能仅在程序启动的时候扫描一遍本地磁盘，并将磁盘中所有的对象散列值读入内存，之后
// 在定位的时候就不需要再次访问磁盘，只需要搜索内存就可以了。
