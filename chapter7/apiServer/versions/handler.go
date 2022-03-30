package versions

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"../../../src/lib/es"
)

// 接口服务的versions包

// versions包比较简单，只有Handler函数，其主要工作都是调用es包的函数来完成

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	from := 0
	size := 1000
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	for {
		metas, e := es.SearchAllVersions(name, from, size)
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for i := range metas {
			b, _ := json.Marshal(metas[i])
			w.Write(b)
			w.Write([]byte("\n"))
		}
		if len(metas) != size {
			return
		}
		from += size
	}
}

// 这个函数首先检查HTTP方法是否为GET,如果不为GET,则返回405 Method NotAllowed;
// 如果方法为GET，则获取 URL中<object_name>的部分，获取的方式跟之前一样，
// 调用strings.Split函数将URL 以“/”为分隔符切成数组并取第三个元素赋值给name。
// 这里要注意的是，如果客户端的HTTP 请求的 URL 是“/versions/”而不含<object_name>部分,
// 那么name 就是空字符串。

// 接下来我们会在一个无限for循环中调用es包的SearchAllVersions函数并将name,from和 size作为参数传递给该函数。
// from 从О开始，size 则固定为1000。es.SearchAllVersions函数会返回一个元数据的数组，
// 我们遍历该数组，将元数据—一写入HTTP响应的正文。
// 如果返回的数组长度不等于size，说明元数据服务中没有更多的数据了，此时我们让函数返回;
// 否则我们就把from的值增加1000进行下一个迭代。

// es包封装了我们访问元数据服务的各种API的操作，本章后续会有详细介绍。
