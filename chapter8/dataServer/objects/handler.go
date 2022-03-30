package objects

import (
	"net/http"
)

// 为了支持对象的删除操作，数据服务的 objects包变化见例8-3。

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodGet {
		get(w, r)
		return
	}
	if m == http.MethodDelete {
		del(w, r)
		// 在Handler函数中，如果访问方式是DELETE，那么调用 del函数。

		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
