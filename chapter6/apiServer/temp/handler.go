package temp

// temp包一共有3个函数，Handler 用于注册HTTP处理函数，head 和 put分别处理相应的访问方法。
// 首先让我们看看temp.Handler函数，见例6-5。

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodHead {
		head(w, r)
		return
	}
	if m == http.MethodPut {
		put(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// temp.Handler函数首先检查访问方式，如果是HEAD则调用head函数，
// 如果是PUT则调用put函数，否则返回405 Method Not Allowed。put相关函数见例6-6
