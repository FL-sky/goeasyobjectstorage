package objects

import (
	"net/http"
	"strings"
)

// 第二个改动则是在读取对象时进行一次数据校验，见例4-11。
func get(w http.ResponseWriter, r *http.Request) {
	file := getFile(strings.Split(r.URL.EscapedPath(), "/")[2])
	if file == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	sendFile(w, file)
	// sendFile有两个输入参数，分别是用于写入对象数据的w和对象的文件名file。
	// 它调用os.Open打开对象文件，并用io.Copy将文件的内容写入w。

}

// get函数首先从URL中获取对象的散列值，然后以散列值为参数调用getFile获得对象的文件名file，
// 如果file为空字符串则返回404 Not Found;否则调用sendFile将该对象文件的内容输出到HTTP响应。
