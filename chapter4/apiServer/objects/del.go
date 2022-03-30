package objects

import (
	"log"
	"net/http"
	"strings"

	"../../../src/lib/es"
)

func del(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	version, e := es.SearchLatestVersion(name)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	e = es.PutMetadata(name, version.Version+1, 0, "")
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// del函数用同样的方式从URL中获取<object_name>并赋值给name。
// 然后它以name为参数调用es.SearchLatestVersion，搜索该对象最新的版本。
// 接下来函数调用es.PutMetadata插入一条新的元数据。
// es.PutMetadata 接受4个输入参数，分别是元数据的name、version、size和 hash。
// 我们可以看到，函数参数中name就是对象的名字,version就是该对象最新版本号加1，size为0，hash为空字符串，
// 以此表示这是一个删除标记。
