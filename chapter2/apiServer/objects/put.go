package objects

import (
	"log"
	"net/http"
	"strings"
)

// 接口服务的objects包跟数据服务有很大区别，其put函数和get函数并不会访问本地磁盘上的对象，
// 而是将HTTP请求转发给数据服务。put函数负责处理对象PUT请求
func put(w http.ResponseWriter, r *http.Request) {
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	c, e := storeObject(r.Body, object)
	if e != nil {
		log.Println(e)
	}
	w.WriteHeader(c)
}

// put函数首先从 URL 中获取<object_name>部分赋值给object，然后将 r.Body和 object作为参数调用storeObject。
// storeObject函数会返回两个结果，第一个返回值是一一个int类型的变量，用来表示HTTP错误代码，
// 我们会使用w.WriteHeader方法把这个错误代码写入HTTP响应，
// 第二个返回值则是一个 error，如果该error不为nil，则我们需要把这个错误打印在 log 中。
