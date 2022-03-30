package objects

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"../../../src/lib/es"
	"../../../src/lib/utils"
)

// objects 包除了新增post 函数以外，还修改了get函数,

func get(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	versionId := r.URL.Query()["version"]
	version := 0
	var e error
	if len(versionId) != 0 {
		version, e = strconv.Atoi(versionId[0])
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	meta, e := es.GetMetadata(name, version)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if meta.Hash == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	hash := url.PathEscape(meta.Hash)
	// 和第5章相比，本章的objects.get函数在调用GetStream 生成stream之后,
	// 还调用utils.GetOffsetFromHeader函数从 HTTP请求的 Range头部获得客户端要求的偏移量offset，
	// 如果 offset不为0，那么需要调用stream 的 Seek 方法跳到offset位置，
	// 设置Content-Range响应头部以及HTTP代码206 Partial Content。
	// 然后继续通过io.Copy输出数据。

	stream, e := GetStream(hash, meta.Size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// utils.GetOffsetFromHeader 函数获取 HTTP 的 Range头部，Range头部的格式必须是“bytes=<first>-”开头，
	// 我们调用strings.Split将<first>部分切取出来并调用strconv.ParseInt 将字符串转化成int64返回。
	offset := utils.GetOffsetFromHeader(r.Header)
	if offset != 0 {
		// RSGetStream.Seek方法有两个输入参数，offset表示需要跳过多少字节whence表示起跳点。
		// 我们的方法只支持从当前位置(io.SeekCurrent）起跳，且跳过的字节数不能为负。
		// 	我们在一个for循环中每次读取32000字节并丢弃，直到读到offset位置为止。

		stream.Seek(offset, io.SeekCurrent)
		w.Header().Set("content-range", fmt.Sprintf("bytes %d-%d/%d", offset, meta.Size-1, meta.Size))
		w.WriteHeader(http.StatusPartialContent)
	}
	io.Copy(w, stream)
	stream.Close()
}
