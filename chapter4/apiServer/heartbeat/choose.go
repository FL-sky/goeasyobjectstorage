package heartbeat

import (
	"math/rand"
)

func ChooseRandomDataServer() string {
	ds := GetDataServers()
	n := len(ds)
	if n == 0 {
		return ""
	}
	return ds[rand.Intn(n)]
}

// ChooseRandomDataServer函数会在当前所有的数据服务节点中随机选出一个节点并返回，如果当前数据服务节点为空，则返回空字符串。
