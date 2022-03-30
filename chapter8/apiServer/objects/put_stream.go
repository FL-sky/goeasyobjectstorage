package objects

// 接口服务的objects包

// 和之前一样，我们略过未改变的部分，只示例objects包中需要改变的函数。
// 首先是PUT对象时需要用到的 putStream函数。
// 它使用了新的 heartbeat.ChooseRandomDataServers函数获取随机数据服务节点地址，
// 并调用rs.NewRSPutStream来生成一个数据流，见例5-3。

import (
	"fmt"

	"../../../src/lib/rs"
	"../heartbeat"
)

func putStream(hash string, size int64) (*rs.RSPutStream, error) {
	servers := heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS, nil)
	if len(servers) != rs.ALL_SHARDS {
		return nil, fmt.Errorf("cannot find enough dataServer")
	}

	return rs.NewRSPutStream(servers, hash, size)
	// rs.NewRSPutStream返回的是一个指向rs.RSPutStream 结构体的指针。相关代码见例5-4。

}
