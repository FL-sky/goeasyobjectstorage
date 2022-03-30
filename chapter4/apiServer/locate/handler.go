package locate

import (
	"encoding/json"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	info := Locate(strings.Split(r.URL.EscapedPath(), "/")[2])
	if len(info) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, _ := json.Marshal(info)
	w.Write(b)
}

// Handler函数用于处理HTTP请求,如果请求方法不为GET,则返回405 Method NotAllowed;
// 如果请求方法为GET，我们将<object_name>作为Locate函数的参数进行定位。
// 如果Locate函数返回的字符串长度为空,说明该对象 locate失败,我们返回404 NotFound;
// 如果不为空，则是拥有该对象的一个数据服务节点的地址，我们将该地址作为HTTP响应的正文输出。
