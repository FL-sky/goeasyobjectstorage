package rs

// reedsolomon包是一个RS编解码的开源库

import (
	"io"

	"github.com/klauspost/reedsolomon"
)

type encoder struct {
	writers []io.Writer
	enc     reedsolomon.Encoder
	cache   []byte
}

func NewEncoder(writers []io.Writer) *encoder {
	enc, _ := reedsolomon.New(DATA_SHARDS, PARITY_SHARDS)
	return &encoder{writers, enc, nil}
}

// RSPutStream本身并没有实现 Write方法，所以实现时函数会直接调用其内嵌结构体encoder的 Write方法。

func (e *encoder) Write(p []byte) (n int, err error) {
	length := len(p)
	current := 0
	for length != 0 {
		next := BLOCK_SIZE - len(e.cache)
		if next > length {
			next = length
		}
		e.cache = append(e.cache, p[current:current+next]...)
		if len(e.cache) == BLOCK_SIZE {
			e.Flush()
		}
		current += next
		length -= next
	}
	return len(p), nil
}

// encoder 的 Write方法在 for循环里将p中待写入的数据以块的形式放入缓存，
// 如果缓存已满就调用Flush方法将缓存实际写入 writers。
// 缓存的上限是每个数据片8000字节,4个数据片共32 000字节。
// 如果缓存里剩余的数据不满32000字节就暂不刷新,等待Write方法下一次被调用。

func (e *encoder) Flush() {
	if len(e.cache) == 0 {
		return
	}
	shards, _ := e.enc.Split(e.cache)
	e.enc.Encode(shards)
	for i := range shards {
		e.writers[i].Write(shards[i])
	}
	e.cache = []byte{}
}

// Flush方法首先调用encoder的成员变量enc的Split方法将缓存的数据切成4个数据片，
// 然后调用enc的 Encode方法生成两个校验片，最后在for循环中将6个片的数据依次写入writers并清空缓存。
