package objects

import (
	"fmt"

	"../../../src/lib/objectstream"
	"../heartbeat"
)

func putStream(object string) (*objectstream.PutStream, error) {
	server := heartbeat.ChooseRandomDataServer()
	if server == "" {
		return nil, fmt.Errorf("cannot find any dataServer")
	}

	return objectstream.NewPutStream(server, object), nil
}

// putStream函数首先调用heartbeat.ChooseRandomDataServer函数获得一个随机数据服务节点的地址server，
// 如果server为空字符串，则意味着当前没有可用的数据服务节点,
// 我们返回一个objectstream.PutStream的空指针和一个“cannot find any dataServer'的error。
// 此时storeObject 会返回 http.StatusServiceUnavailable,客户端会收到HTTP错误代码503 Service Unavailable。
// 如果server不为空，则以server和 object为参数调用 objectstream.NewPutStream 生成一个objectstream.PutStream的指针并返回。

// objectstream包是我们对Go语言 http包的一个封装，用来把一些http函数的调用转换成读写流的形式，方便我们处理。
// 其 PutStream相关代码见例2-8。
