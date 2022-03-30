package locate

///数据服务的 locate包

import (
	"os"
	"strconv"

	"../../../src/lib/rabbitmq"
)

//Locate函数用os.Stat访问磁盘上对应的文件名，用os.IsNotExist判断文件名是否存在，
//如果存在则定位成功返回 true，否则定位失败返回false。

func Locate(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func StartLocate() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	q.Bind("dataServers")
	c := q.Consume()
	for msg := range c {
		object, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		if Locate(os.Getenv("STORAGE_ROOT") + "/objects/" + object) {
			q.Send(msg.ReplyTo, os.Getenv("LISTEN_ADDRESS"))
		}
	}
}

// StartLocate 函数会创建一个rabbitmq.RabbitMQ结构体，并调用其 Bind方法绑定dataServers exchange。
// rabbitmq.RabbitMQ结构体的Consume方法会返回一个Go语言的channel，用range遍历这个 channel可以接收消息。
// 消息的正文内容是接口服务发送过来的需要做定位的对象名字。
// 由于经过JSON 编码，所以对象名字上有一对双引号(JSON是 JavaScript Object Notation的缩写，
// 	是一种语言独立的数据格式。虽然它起源自JavaScript，但是目前很多编程语言都包含处理JSON 格式数据的代码)。
// 	strconv.Unquote函数的作用就是将输入的字符串前后的双引号去除并作为结果返回。
// 	我们在对象名字前加上相应的存储目录并以此作为文件名，然后调用Locate函数检查文件是否存在，
// 	如果存在，则调用Send方法向消息的发送方返回本服务节点的监听地址，表示该对象存在于本服务节点上。
