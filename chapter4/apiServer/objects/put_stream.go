package objects

import (
	"fmt"

	"../../../src/lib/objectstream"
	"../heartbeat"
)

func putStream(hash string, size int64) (*objectstream.TempPutStream, error) {
	server := heartbeat.ChooseRandomDataServer()
	if server == "" {
		return nil, fmt.Errorf("cannot find any dataServer")
	}

	return objectstream.NewTempPutStream(server, hash, size)
}

// putStream 唯一的变化在于:
// 第2章的 putStream 调用objectstream.NewPutStream生成一个对象的写入流，
// 而本章的 putStream调用的则是 objectstream.NewTempPutStream,
// 这是因为数据服务的temp接口代替了原先的对象PUT接口。TempPutStream相关代码见例4-2。
