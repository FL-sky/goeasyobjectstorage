package objects

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"../../../src/lib/es"
)

func get(w http.ResponseWriter, r *http.Request) {
	// 跟第2章相比,本章objects.get函数从URL获取了对象的名字之后还需从URL的查询参数中获取“version”参数的值。
	// r.URL 的类型是*url.URL，它是指向url.URL结构体的指针。
	// url.URL结构体的Query方法会返回一个保存URL所有查询参数的map,该map 的键是查询参数的名字，
	// 而值则是一个字符串数组，这是因为 HTTP的 URL查询参数允许存在多个值。
	// 以“version”为key就可以得到URL 中该查询参数的所有值，然后赋值给versionld变量。
	// 如果URL中并没有“version”这个查询参数，versionId变量则是空数组。
	// 由于我们不考虑多个“version”查询参数的情况，
	// 所以我们始终以versionId数组的第1个元素作为客户端提供的版本号，
	// 并将其从字符串转换为整型赋值给version变量。
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	versionId := r.URL.Query()["version"]
	version := 0
	var e error
	if len(versionId) != 0 {
		version, e = strconv.Atoi(versionId[0])
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	// 然后我们以对象的名字和版本号为参数调用es.GetMetadata，得到对象的元数据meta。
	// meta.Hash 就是对象的散列值。如果散列值为空字符串说明该对象该版本是一个删除标记，
	// 我们返回404 Not Found;否则以散列值为对象名从数据服务层获取对象并输出。
	meta, e := es.GetMetadata(name, version)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if meta.Hash == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	object := url.PathEscape(meta.Hash)
	stream, e := getStream(object)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	io.Copy(w, stream)
}
