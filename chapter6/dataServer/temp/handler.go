package temp

import "net/http"

// Handler函数相比第4章多了对HEAD/PUT 方法的处理。
// 如果接口服务以HEAD方式访问数据服务的temp接口，Handler 会调用head;
// 如果接口服务以GET方式访问数据服务的 temp接口，则Handler会调用get。

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodHead {
		head(w, r)
		return
	}
	if m == http.MethodGet {
		get(w, r)
		return
	}
	if m == http.MethodPut {
		put(w, r)
		return
	}
	if m == http.MethodPatch {
		patch(w, r)
		return
	}
	if m == http.MethodPost {
		post(w, r)
		return
	}
	if m == http.MethodDelete {
		del(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
