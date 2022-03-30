package locate

// 接口服务的locate包

import (
	"os"
	"strconv"
	"time"

	"../../../src/lib/rabbitmq"
)

// 接口服务的 locate包有3个函数，用于向数据服务节点群发定位消息并接收反馈,

func Locate(name string) string {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	q.Publish("dataServers", name)
	c := q.Consume()
	go func() {
		time.Sleep(time.Second)
		q.Close()
	}()
	msg := <-c
	s, _ := strconv.Unquote(string(msg.Body))
	return s
}

// Locate函数接受一个string 类型的参数name，也就是需要定位的对象的名字。
// 它会创建一个新的消息队列，并向 dataServers exchange群发这个对象名字的定位信息,
// 随后用goroutine启动一个匿名函数，用于在1s后关闭这个临时消息队列。
// 这是为了设置一个超时机制，避免无止境的等待。
// 因为Locate函数随后就会阻塞等待数据服务节点向自己发送反馈消息，如果1s后还没有任何反馈，则消息队列关闭，
// 我们会收到一个长度为0的消息，此时我们需要返回一个空字符串;
// 如果在1s内有来自数据服务节点的消息，则返回该消息的正文内容，也就是该数据服务节点的监听地址。

func Exist(name string) bool {
	return Locate(name) != ""
}

// Exist函数通过检查Locate结果是否为空字符串来判定对象是否存在。
