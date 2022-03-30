package heartbeat

//接口服务的heartbeat包

import (
	"os"
	"strconv"
	"sync"
	"time"

	"../../../src/lib/rabbitmq"
)

///接口服务的heartbeat包有4个函数、用于接收和处理来自数据服务节点的心跳消息,

var dataServers = make(map[string]time.Time)

// 包变量 dataServers 的类型是map[string]time.Time。其中，键的类型是string，值的类型则是time.Time结构体，
// 它在整个包内可见，用于缓存所有的数据服务节点。

var mutex sync.Mutex

func ListenHeartbeat() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	q.Bind("apiServers")
	c := q.Consume()
	go removeExpiredDataServer()
	for msg := range c {
		dataServer, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		mutex.Lock()
		dataServers[dataServer] = time.Now()
		mutex.Unlock()
	}
}

// ListenHeartbeat函数创建消息队列来绑定apiServers exchange并通过go channel监听每一个来自数据服务节点的心跳消息，
// 将该消息的正文内容，也就是数据服务节点的监听地址作为map的键，收到消息的时间作为值存入dataServers .

func removeExpiredDataServer() {
	for {
		time.Sleep(5 * time.Second)
		mutex.Lock()
		for s, t := range dataServers {
			if t.Add(10 * time.Second).Before(time.Now()) {
				delete(dataServers, s)
			}
		}
		mutex.Unlock()
	}
}

// removeExpiredDataServer函数在一个goroutine中运行，每隔5s扫描一遍dataServers,
// 并清除其中超过10s没收到心跳消息的数据服务节点。

func GetDataServers() []string {
	mutex.Lock()
	defer mutex.Unlock()
	ds := make([]string, 0)
	for s, _ := range dataServers {
		ds = append(ds, s)
	}
	return ds
}

// GetDataServers遍历dataServers并返回当前所有的数据服务节点。
// 注意，这里对dataServers的读写全部都需要mutex的保护，以防止多个goroutine并发读写map造成错误。
// Go语言的map可以支持多个goroutine同时读，但不能支持多个goroutine同时写或同时既读又写,
// 所以我们在这里使用一个互斥锁mutex 保护map的并发读写。mutex.的类型是sync.Mutex，
// sync是 Go语言自带的一个标准包，它提供了包括Mutex, RWMutex在内的多种互斥锁的实现。
// 本书使用了较为简单的互斥锁Mutex，无论读写都只允许一个goroutine操作 map。
// 一个更具有效率的方法是使用读写锁RWMutex，因为读写锁可以允许多个goroutine同时读。有兴趣的读者可自行实现。
