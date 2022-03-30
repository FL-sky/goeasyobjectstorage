package objects

// 数据服务的objects包

import "net/http"

// 数据服务除了新增temp包用于处理temp接口的请求以外，原来的objects包也需要进行改动，
// 第一个改动的地方是删除objects接口的PUT方法，见例4-10。

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodGet {
		get(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// 跟第⒉章的数据服务相比，本章的objects.Handler去除了处理PUT方法的put函数。
// 这是因为现在数据服务的对象上传完全依靠temp接口的临时对象转正，所以不再需要objects接口的PUT方法。

// 第二个改动则是在读取对象时进行一次数据校验，见例4-11。
