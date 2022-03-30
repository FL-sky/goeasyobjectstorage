package objects

import (
	"fmt"

	"../../../src/lib/rs"
	"../heartbeat"
	"../locate"
)

func GetStream(hash string, size int64) (*rs.RSGetStream, error) {
	locateInfo := locate.Locate(hash)
	if len(locateInfo) < rs.DATA_SHARDS {
		return nil, fmt.Errorf("object %s locate fail, result %v", hash, locateInfo)
	}
	dataServers := make([]string, 0)
	if len(locateInfo) != rs.ALL_SHARDS {
		dataServers = heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS-len(locateInfo), locateInfo)
	}
	return rs.NewRSGetStream(locateInfo, dataServers, hash, size)
}

// objects包处理对象GET相关的其他部分代码并没有发生改变,
// 但是在对象GET过程中，我们读取对象数据的流调用的是 rs.RSGetStream，
// 所以实际背后调用的Read方法也不一样，见例5-9。
