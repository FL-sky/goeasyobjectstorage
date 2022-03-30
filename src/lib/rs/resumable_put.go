package rs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"../../lib/objectstream"
	"../../lib/utils"
)

type resumableToken struct {
	Name    string
	Size    int64
	Hash    string
	Servers []string
	Uuids   []string
}

type RSResumablePutStream struct {
	*RSPutStream
	*resumableToken
}

// rs.NewRSResumablePutStream 创建的stream 的类型是一个指向RSResumablePutStream结构体的指针。
// 该结构体内嵌了RSPutStream和 resumableToken。RSPutStream我们在上一章已经讲述过了。
// resumableToken中保存了对象的名字、大小、散列值，
// 另外还有6个分片所在的数据服务节点地址和 uuid，分别以数组的形式保存。

func NewRSResumablePutStream(dataServers []string, name, hash string, size int64) (*RSResumablePutStream, error) {
	putStream, e := NewRSPutStream(dataServers, hash, size)
	if e != nil {
		return nil, e
	}
	uuids := make([]string, ALL_SHARDS)
	for i := range uuids {
		uuids[i] = putStream.writers[i].(*objectstream.TempPutStream).Uuid
	}
	token := &resumableToken{name, size, hash, dataServers, uuids}
	return &RSResumablePutStream{putStream, token}, nil
}

// NewRSResumablePutStreamFromToken 函数对token进行Base64解码，
// 然后将JSON数据编出形成resumableToken结构体t,
// t的Servers和 Uuids 数组中保存了当初创建的6个分片临时对象所在的数据服务节点地址和 uuid，
// 我们根据这些信息创建6个objectstream. TempPutStream保存在 writers 数组，
// 以writers 数组为参数创建encoder结构体enc，以enc为内嵌结构体创建RSPutStream，
// 并最终以 RSPutStream和t为内嵌结构体创建 RSResumablePutStream返回。

func NewRSResumablePutStreamFromToken(token string) (*RSResumablePutStream, error) {
	b, e := base64.StdEncoding.DecodeString(token)
	if e != nil {
		return nil, e
	}

	var t resumableToken
	e = json.Unmarshal(b, &t)
	if e != nil {
		return nil, e
	}

	writers := make([]io.Writer, ALL_SHARDS)
	for i := range writers {
		writers[i] = &objectstream.TempPutStream{t.Servers[i], t.Uuids[i]}
	}
	enc := NewEncoder(writers)
	return &RSResumablePutStream{&RSPutStream{enc}, &t}, nil
}

// RSResumablePutStream.ToToken方法将自身数据以JSON格式编入;然后返回经过Base64编码后的字符串。

func (s *RSResumablePutStream) ToToken() string {
	b, _ := json.Marshal(s)
	return base64.StdEncoding.EncodeToString(b)
}

// 注意，
// 任何人都可以将Base64编码的字符串解码，本书的实现代码并未对token加密，
// 任何人都可以轻易从接口服务返回的响应头部中获取RSResumablePutStream 结构体的内部信息。
// 这是一个很大的信息泄露。本书旨在介绍和实现对象存储的各种功能，而信息安全不属于本书的范畴。
// 对信息安全有要求的读者需要自行实现对token的加密和解密操作。

func (s *RSResumablePutStream) CurrentSize() int64 {
	r, e := http.Head(fmt.Sprintf("http://%s/temp/%s", s.Servers[0], s.Uuids[0]))
	if e != nil {
		log.Println(e)
		return -1
	}
	if r.StatusCode != http.StatusOK {
		log.Println(r.StatusCode)
		return -1
	}
	size := utils.GetSizeFromHeader(r.Header) * DATA_SHARDS
	if size > s.Size {
		size = s.Size
	}
	return size
}

// RSResumablePutStream.CurrentSize 以HEAD方法获取第一个分片临时对象的大小并乘以4作为size返回。
// 如果size超出了对象的大小，则返回对象的大小。
