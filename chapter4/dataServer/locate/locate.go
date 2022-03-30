package locate

import (
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"../../../src/lib/rabbitmq"
)

var objects = make(map[string]int)
var mutex sync.Mutex

// 例4-4显示了数据服务的locate包的实现，
// 函数中的包变量objects是一个以字符串为键，整型为值的map，它用于缓存所有对象。
// mutex互斥锁用于保护对 objects 的读写操作。
// Locate函数利用Go语言的map操作判断某个散列值是否存在于objects 中，如果存在返回true，否则返回false。

func Locate(hash string) bool {
	mutex.Lock()
	_, ok := objects[hash]
	mutex.Unlock()
	return ok
}

// Add函数用于将一个散列值加入缓存，其输入参数hash作为存入map 的键，值为1。

func Add(hash string) {
	mutex.Lock()
	objects[hash] = 1
	mutex.Unlock()
}

// Del函数则相反，用于将一个散列值移出缓存。

func Del(hash string) {
	mutex.Lock()
	delete(objects, hash)
	mutex.Unlock()
}

// StartLocate函数大半部分和第2章一样，
// 第2章的 StartLocate函数需要拼出完整的文件名作为Locate的参数，
// 本章则直接将RabbitMQ消息队列里收到的对象散列值作为Locate参数。

func StartLocate() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	q.Bind("dataServers")
	c := q.Consume()
	for msg := range c {
		hash, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		exist := Locate(hash)
		if exist {
			q.Send(msg.ReplyTo, os.Getenv("LISTEN_ADDRESS"))
		}
	}
}

// CollectObjects函数首先调用 filepath.Glob读取$STORAGE_ROOT/objects/目录里
// 的所有文件，对这些文件一一调用filepath.Base获取其基本文件名，
// 也就是对象的散列值，将它们加入 objects缓存
func CollectObjects() {
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")
	for i := range files {
		hash := filepath.Base(files[i])
		objects[hash] = 1
	}
}
