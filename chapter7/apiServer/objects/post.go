package objects

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"../../../src/lib/es"
	"../../../src/lib/rs"
	"../../../src/lib/utils"
	"../heartbeat"
	"../locate"
)

func post(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	size, e := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Println("missing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if locate.Exist(url.PathEscape(hash)) {
		e = es.AddVersion(name, hash, size)
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		return
	}
	ds := heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS, nil)
	if len(ds) != rs.ALL_SHARDS {
		log.Println("cannot find enough dataServer")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	stream, e := rs.NewRSResumablePutStream(ds, name, url.PathEscape(hash), size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("location", "/temp/"+url.PathEscape(stream.ToToken()))
	w.WriteHeader(http.StatusCreated)
}

// post 函数和 put函数的处理流程在前半段是一样的，都是从请求的URL中获得对象的名字，
// 从请求的相应头部获得对象的大小和散列值，然后对散列值进行定位。
// 如果该散列值已经存在，那么我们可以直接往元数据服务添加新版本并返回200 OK;
// 如果散列值不存在，那么随机选出6个数据节点,
// 然后调用rs.NewRSResumablePutStream生成数据流stream，
// 并调用其ToToken方法生成一个字符串token，放入 Location 响应头部，并返回HTTP代码201 Created。
