package objects

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"../../../src/lib/es"
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
	// 我们可以看到，原本调用getStream的地方变成调用GetStream且其参数多了一个size。
	// 大小写的变化是为了将该函数导出给包外部使用，我们在之前已经多次提到Go语言这一个特性了。
	// 增加 size参数是因为RS码的实现要求每一个数据片的长度完全一样，
	// 在编码时如果对象长度不能被4整除，函数会对最后一个数据片进行填充。
	// 因此在解码时必须提供对象的准确长度，防止填充数据被当成原始对象数据返回。

	stream, e := GetStream(hash, meta.Size)
	// GetStream函数首先根据对象散列值 hash定位对象，
	// 如果反馈的定位结果locateInfo数组长度小于4，则返回错误;
	// 如果 locateInfo数组的长度不为6，说明该对象有部分分片丢失，
	// 我们调用heartbeat.ChooseRandomDataServers随机选取用于接收恢复分片的数据服务节点，
	// 以数组的形式保存在dataServers里。
	// 最后我们以locateInfo、dataServers.hash 以及对象的大小size为参数
	// 调用rs.NewRSGetStream函数创建rs.RSGetStream,相关函数见例5-8。

	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// 在调用io.Copy将对象数据流写入HTTP响应时，如果返回错误，说明对象数据在RS解码过程中发生了错误，
	// 这意味着该对象已经无法被读取，我们返回404 Not Found,
	// 如果没有返回错误，我们需要在 get函数最后调用stream.Close方法。
	// GetStream返回的stream的类型是一个指向rs.RSGetStream结构体的指针，
	// 我们在GET对象时会对缺失的分片进行即时修复，修复的过程也使用数据服务的 temp接口，
	// RSGetStream的Close方法用于在流关闭时将临时对象转正。

	_, e = io.Copy(w, stream)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	stream.Close()
}
