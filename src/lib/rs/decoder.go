package rs

import (
	"io"

	"github.com/klauspost/reedsolomon"
)

// decoder结构体除了readers，writers两个数组以外还包含若干成员，
// enc的类型是reedsolomon.Encoder接口用于RS解码，size是对象的大小，
// cache和 cacheSize用于缓存读取的数据，total表示当前已经读取了多少字节。

type decoder struct {
	readers   []io.Reader
	writers   []io.Writer
	enc       reedsolomon.Encoder
	size      int64
	cache     []byte
	cacheSize int
	total     int64
}

func NewDecoder(readers []io.Reader, writers []io.Writer, size int64) *decoder {
	enc, _ := reedsolomon.New(DATA_SHARDS, PARITY_SHARDS)
	return &decoder{readers, writers, enc, size, nil, 0, 0}
}

// RSGetStream的 Read方法就是其内嵌结构体decoder 的Read方法。decoder的Read
// 方法当cache中没有更多数据时会调用getData方法获取数据，如果getData返回的e不为nil，
// 说明我们没能获取更多数据，则返回0和这个e通知调用方。
// length是Read方法输入参数p 的数组长度，
// 如果 length超出当前缓存的数据大小，我们令 length等于缓存的数据大小。
// 我们用copy函数将缓存中length长度的数据复制给输入参数p,然后调整缓存，仅保留剩下的部分。
// 最后Read方法返回length，通知调用方本次读取一共有多少数据被复制到p中。

func (d *decoder) Read(p []byte) (n int, err error) {
	if d.cacheSize == 0 {
		e := d.getData()
		if e != nil {
			return 0, e
		}
	}
	length := len(p)
	if d.cacheSize < length {
		length = d.cacheSize
	}
	d.cacheSize -= length
	copy(p, d.cache[:length])
	d.cache = d.cache[length:]
	return length, nil
}

// getData方法首先判断当前已经解码的数据大小是否等于对象原始大小，
// 如果已经相等，说明所有数据都已经被读取，我们返回 io.EOF;
// 如果还有数据需要读取，我们会创建一个长度为6的数组 shards，以及一个长度为0的整型数组repairIds。
// shards数组中每一个元素都是一个字节数组，用于保存相应分片中读取的数据。
// 我们在一个for循环中遍历6个shards，如果某个分片对应的reader是nil，说明该分片已丢失，
// 我们会在repairIds中添加该分片的id;如果对应的reader 不为nil，
// 那么对应的shards需要被初始化成一个长度为8000的字节数组，
// 然后调用io.ReadFull 从 reader 中完整读取8000字节的数据保存在shards 里;
// 如果发生了非 EOF 失败，该shards会被置为nil,
// 如果读取的数据长度n不到8000字节，我们将该shards实际的长度缩减为n。

func (d *decoder) getData() error {
	if d.total == d.size {
		return io.EOF
	}
	shards := make([][]byte, ALL_SHARDS)
	repairIds := make([]int, 0)
	for i := range shards {
		if d.readers[i] == nil {
			repairIds = append(repairIds, i)
		} else {
			shards[i] = make([]byte, BLOCK_PER_SHARD)
			n, e := io.ReadFull(d.readers[i], shards[i])
			if e != nil && e != io.EOF && e != io.ErrUnexpectedEOF {
				shards[i] = nil
			} else if n != BLOCK_PER_SHARD {
				shards[i] = shards[i][:n]
			}
		}
	}

	// 	遍历读取一轮之后，要么每个shards中保存了读取自对应分片的数据，要么因为分
	// 片丢失或读取错误，该shards被置为nil。
	// 我们调用成员enc的 Reconstruct方法尝试将被置为nil 的 shards 恢复出来，
	// 这一步如果返回错误，说明我们的对象已经遭到了不可修复的破坏，我们只能将错误原样返回给上层。
	// 如果修复成功，6个 shards中都保存了对应分片的正确数据，
	// 我们遍历repairIds，将需要恢复的分片的数据写入相应的writer。

	e := d.enc.Reconstruct(shards)
	if e != nil {
		return e
	}
	for i := range repairIds {
		id := repairIds[i]
		d.writers[id].Write(shards[id])
	}

	// 最后，我们遍历4个数据分片，将每个分片中的数据添加到缓存cache 中，
	// 修改缓存当前的大小cacheSize以及当前已经读取的全部数据的大小total。

	for i := 0; i < DATA_SHARDS; i++ {
		shardSize := int64(len(shards[i]))
		if d.total+shardSize > d.size {
			shardSize -= d.total + shardSize - d.size
		}
		d.cache = append(d.cache, shards[i][:shardSize]...)
		d.cacheSize += int(shardSize)
		d.total += shardSize
	}
	return nil
}

// 恢复分片的写入需要用到数据服务的temp接口，
// 所以 objects.get函数会在最后调用stream.Close方法将用于写入恢复分片的临时对象转正，
// 该方法的实现见例5-10。
// 例5-10 RSGetStream.Close方法
