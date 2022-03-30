package heartbeat

///数据服务的heartbeat包

import (
	"os"
	"time"

	"../../../src/lib/rabbitmq"
)

// heartbeat包的实现非常简单,只有一个 StartHeartbeat函数
//每5s向apiServers exchange发送一条消息——把本服务节点的监听地址发送出去。见例2-2。

func StartHeartbeat() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	for {
		q.Publish("apiServers", os.Getenv("LISTEN_ADDRESS"))
		time.Sleep(5 * time.Second)
	}
}

// heartbeat.StartHeartbeat调用rabbitmq.New创建了一个rabbitmq.RabbitMQ结构体，
// 并在一个无限循环中调用rabbitmq.RabbitMQ结构体的Publish方法向 apiServersexchange发送本节点的监听地址。
// 由于该函数在一个goroutine中执行,所以就算不返回也不会影响其他功能。
// rabbitmq 包封装了我们对消息队列的操作，本章后续会详细介绍rabbitmq包的各个函数
