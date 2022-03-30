package objects

//接口服务的objects包

import "net/http"

// 本章加入元数据服务以后,接口服务的objects包与上一章相比发生了较大的变化,
// 除了多了一个对象的 DELETE方法以外，对象的PUT和 GET方法也都有所改变，
// 它们需要跟元数据服务互动。首先让我们从Handler 函数的改变开始看起,

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
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

// 可以看到，跟上一章相比，在Handler里多出了对 DELETE方法的处理函数del。
