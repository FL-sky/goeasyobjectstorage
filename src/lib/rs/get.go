package rs

import (
	"fmt"
	"io"

	"../../lib/objectstream"
)

type RSGetStream struct {
	*decoder
}

// RSGetStream结构体内嵌decoder结构体。
// NewRSGetStream函数首先检查locateInfo和 dataServers的总数是否为6，满足4+2RS码的需求。
// 如果不满足，则返回错误。
// 然后我们需要创建一个长度为6的io.Reader 数组readers，用于读取6个分片的数据。
// 我们用一个 for循环遍历6个分片的id，在 locateInfo中查找该分片所在的数据服务节点地址，
// 如果某个分片a相对的数据服务节点地址为空，说明该分片丢失，我们需要取一个随机数据服务节点补上;
// 如果数据服务节点存在，我们调用objectstream.NewGetStream打开一个对象读取流用于读取该分片数据，
// 打开的流被保存在readers 数组相应的元素中。

func NewRSGetStream(locateInfo map[int]string, dataServers []string, hash string, size int64) (*RSGetStream, error) {
	if len(locateInfo)+len(dataServers) != ALL_SHARDS {
		return nil, fmt.Errorf("dataServers number mismatch")
	}

	readers := make([]io.Reader, ALL_SHARDS)
	for i := 0; i < ALL_SHARDS; i++ {
		server := locateInfo[i]
		if server == "" {
			locateInfo[i] = dataServers[0]
			dataServers = dataServers[1:]
			continue
		}
		reader, e := objectstream.NewGetStream(server, fmt.Sprintf("%s.%d", hash, i))
		if e == nil {
			readers[i] = reader
		}
	}

	// readers第一次遍历处理完毕后，有两种情况会导致readers 数组中某个元素为nil,
	// 一种是该分片数据服务节点地址为空;而另一种则是数据服务节点存在但打开流失败。
	// 我们用for循环再次遍历readers，如果某个元素为nil，则调用objectstream.NewTemp
	// PutStream 创建相应的临时对象写入流用于恢复分片。打开的流被保存在 writers 数组相应的元素中。

	writers := make([]io.Writer, ALL_SHARDS)
	perShard := (size + DATA_SHARDS - 1) / DATA_SHARDS
	var e error
	for i := range readers {
		if readers[i] == nil {
			writers[i], e = objectstream.NewTempPutStream(locateInfo[i], fmt.Sprintf("%s.%d", hash, i), perShard)
			if e != nil {
				return nil, e
			}
		}
	}

	// 处理完成后，readers和 writers 数组形成互补的关系，对于某个分片id，
	// 要么在readers中存在相应的读取流，要么在 writers中存在相应的写入流。
	// 我们将这样的两个数组以及对象的大小size作为参数调用NewDecoder 生成 decoder结构体的指针dec,
	// 并将其作为RSGetStream的内嵌结构体返回。

	// NewDecoder函数调用reedsolomon.New创建4+2RS 码的解码器enc，
	// 并设置decoder结构体中相应的属性后返回。

	dec := NewDecoder(readers, writers, size)
	return &RSGetStream{dec}, nil
}

func (s *RSGetStream) Close() {
	for i := range s.writers {
		if s.writers[i] != nil {
			s.writers[i].(*objectstream.TempPutStream).Commit(true)
		}
	}
}

// Close方法遍历 writers成员，如果某个分片的 writer不为nil，
// 则调用其 Commit方法，参数为 true，意味着临时对象将被转正。
// objectstream.TempPutStream的详细实现见第4章。

func (s *RSGetStream) Seek(offset int64, whence int) (int64, error) {
	if whence != io.SeekCurrent {
		panic("only support SeekCurrent")
	}
	if offset < 0 {
		panic("only support forward seek")
	}
	for offset != 0 {
		length := int64(BLOCK_SIZE)
		if offset < length {
			length = offset
		}
		buf := make([]byte, length)
		io.ReadFull(s, buf)
		offset -= length
	}
	return offset, nil
}

// RSGetStream.Seek方法有两个输入参数，offset表示需要跳过多少字节whence表示起跳点。
// 我们的方法只支持从当前位置(io.SeekCurrent）起跳，且跳过的字节数不能为负。
// 	我们在一个for循环中每次读取32000字节并丢弃，直到读到offset位置为止。
