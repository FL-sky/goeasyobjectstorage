package locate

import (
	"encoding/json"
	"os"
	"time"

	"../../../src/lib/rabbitmq"
	"../../../src/lib/rs"
	"../../../src/lib/types"
)

// 为了实现RS码，接口服务的locate、heartbeat和 objects包都需要发生变化，
// 首先让我们来看一下接口服务locate包发生的改变。

func Locate(name string) (locateInfo map[int]string) {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	q.Publish("dataServers", name)
	c := q.Consume()
	go func() {
		time.Sleep(time.Second)
		q.Close()
	}()
	locateInfo = make(map[int]string)
	for i := 0; i < rs.ALL_SHARDS; i++ {
		msg := <-c
		if len(msg.Body) == 0 {
			return
		}
		var info types.LocateMessage
		json.Unmarshal(msg.Body, &info)
		locateInfo[info.Id] = info.Addr
	}
	return
}

// locate包的Handler函数和第2章相比没有发生变化，这里略过。
// 在第2章，我们的Locate函数从接收定位反馈消息的临时消息队列中只获取1条反馈消息，
// 现在我们需要一个for 循环来获取最多6条消息，
// 每条消息都包含了拥有某个分片的数据服务节点的地址和分片的id，并被放在输出参数的 locateInfo变量中返回。
// rs.ALL_SHARDS是rs包的常数6，代表一共有4+2个分片。
// locateInfo的类型是以int为键、string为值的map，它的键是分片的id，而值则是含有该分片的数据服务节点地址。
// 1s超时发生时，无论当前收到了多少条反馈消息都会立即返回。

func Exist(name string) bool {
	return len(Locate(name)) >= rs.DATA_SHARDS
}

// Exist函数判断收到的反馈消息数量是否大于等于4，
// 为true则说明对象存在，否则说明对象不存在(或者说就算存在我们也无法读取)。
