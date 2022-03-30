package objects

import (
	"log"
	"net/http"
	"strings"

	"../../../src/lib/es"
	"../../../src/lib/utils"
)

// 跟第3章的实现相比，put函数唯一的区别在于storeObject多了一个size参数。
// 这是因为我们新的PUT流程需要在一开始就确定临时对象的大小。

func put(w http.ResponseWriter, r *http.Request) {
	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Println("missing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	size := utils.GetSizeFromHeader(r.Header)
	c, e := storeObject(r.Body, hash, size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(c)
		return
	}
	if c != http.StatusOK {
		w.WriteHeader(c)
		return
	}

	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	e = es.AddVersion(name, hash, size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// objects包处理对象PUT相关的其他部分代码并没有发生改变，
// 但是在对象PUT过程中，我们写入对象数据的流调用的是rs.RSPutStream，
// 所以实际背后调用的 Write方法也不一样，见例5-5。
