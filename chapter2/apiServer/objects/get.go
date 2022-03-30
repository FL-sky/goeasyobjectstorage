package objects

import (
	"io"
	"log"
	"net/http"
	"strings"
)
)
// objects包的get函数用来处理 GET请求,
func get(w http.ResponseWriter, r *http.Request) {
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	stream, e := getStream(object)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	io.Copy(w, stream)
}
// 和put一样，get函数先获取<object_name>，然后以之为参数调用getStream 生成一个类型为 io.Reader的 stream,
// 如果出现错误，则打印 log 并返回404Not Found;否则用io.Copy将stream的内容写入HTTP响应的正文。
