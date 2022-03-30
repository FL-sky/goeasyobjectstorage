package rs

import (
	"fmt"
	"io"

	"../../lib/objectstream"
)

// RSPutStream结构体内嵌了一个encoder 结构体。

type RSPutStream struct {
	*encoder
}

// Go语言没有面向对象语言常见的继承机制，而是通过内嵌来连接对象之间的关系。
// 当结构体A包含了指向结构体B的无名指针时，我们就说A内嵌了B。
// A的使用者可以像访问A的方法或成员一样访问B的方法或成员。

// NewRSPutStream函数有3个输入参数，
// dataServers是一个字符串数组，用来保存6个数据服务节点的地址,
// hash和 size分别是需要PUT的对象的散列值和大小。
// 我们首先检查dataServers 的长度是否为6,如果不为6，则返回错误。
// 然后根据size计算出每个分片的大小perShard，也就是size/4再向上取整。
// 然后我们创建了一个长度为6的io.Writers 数组，其中每一个元素都是一个objectstream.TempPutStream，
// 用于上传一个分片对象。
// 最后我们调用NewEncoder 函数创建一个encoder结构体的指针enc并将其作为RSPutStream的内嵌结构体返回。

func NewRSPutStream(dataServers []string, hash string, size int64) (*RSPutStream, error) {
	if len(dataServers) != ALL_SHARDS {
		return nil, fmt.Errorf("dataServers number mismatch")
	}

	perShard := (size + DATA_SHARDS - 1) / DATA_SHARDS
	writers := make([]io.Writer, ALL_SHARDS)
	var e error
	for i := range writers {
		writers[i], e = objectstream.NewTempPutStream(dataServers[i],
			fmt.Sprintf("%s.%d", hash, i), perShard)
		if e != nil {
			return nil, e
		}
	}
	// NewEncoder函数调用reedsolomon.New生成了一个具有4个数据片加两个校验片的RS码编码器enc，
	// 并将输入参数writers和 enc作为生成的encoder结构体的成员返回。

	enc := NewEncoder(writers)

	return &RSPutStream{enc}, nil
	// encoder结构体包含了
	// 一个io.Writers 数组 writers，
	// 一个reedsolomon.Encoder 接口的enc以及
	// 一个用来做输入数据缓存的字节数组cache。

}

func (s *RSPutStream) Commit(success bool) {
	s.Flush()
	for i := range s.writers {
		s.writers[i].(*objectstream.TempPutStream).Commit(success)
	}
}

// Commit方法首先调用其内嵌结构体encoder的Flush方法将缓存中最后的数据写入,
// 然后对encoder的成员数组 writers中的元素调用Commit方法将6个临时对象依次转正或删除。
