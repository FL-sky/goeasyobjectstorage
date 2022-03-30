package objects

// 接口服务只有objects.get函数发生了改变,

import (
	"compress/gzip"
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
	stream, e := GetStream(hash, meta.Size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	offset := utils.GetOffsetFromHeader(r.Header)
	if offset != 0 {
		stream.Seek(offset, io.SeekCurrent)
		w.Header().Set("content-range", fmt.Sprintf("bytes %d-%d/%d", offset, meta.Size-1, meta.Size))
		w.WriteHeader(http.StatusPartialContent)
	}
	// 相比第6章中介绍的，get函数多了一个对 Accept-Encoding请求头部的检查，
	// 如果该头部中含有gzip，则说明客户端可以接受 gzip 压缩数据，
	// 我们会设置Content-Encoding响应头部为gzip,同时以w为参数调用gzip.NewWriter创建一个指向gzip.Writer结构体的指针w2。
	// 我们将对象数据流stream的内容用io.Copy写入 w2，数据就会自动被压缩，然后写入W。

	acceptGzip := false
	encoding := r.Header["Accept-Encoding"]
	for i := range encoding {
		if encoding[i] == "gzip" {
			acceptGzip = true
			break
		}
	}
	if acceptGzip {
		w.Header().Set("content-encoding", "gzip")
		w2 := gzip.NewWriter(w)
		io.Copy(w2, stream)
		w2.Close()
	} else {
		io.Copy(w, stream)
	}
	stream.Close()
}
