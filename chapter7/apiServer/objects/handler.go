package objects

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodPost {
		post(w, r)
		return
	}
	if m == http.MethodPut {
		put(w, r)
		return
	}
	if m == http.MethodGet {
		get(w, r)
		return
	}
	if m == http.MethodDelete {
		del(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// 相比第3章,我们可以看到本章的Handler多了一个对 POST 方法的处理函数post,相关函数见例6-3。
// 例6-3 objects.post 相关函数
