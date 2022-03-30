package objects

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"../../../src/lib/es"
	"../../../src/lib/utils"
)

func put(w http.ResponseWriter, r *http.Request) {
	// GetHashFromHeader函数首先调用h.Get获取“digest”头部。
	// r的类型我们在第1章已经介绍过了，是一个指向:http.Request 的指针。
	// 它的 Header 成员类型则是一个http.Header，用于记录HTTP 的头部，
	// 其Get方法用于根据提供的参数获取相对应的头部的值。在这里，我们获取的就是 HTTP请求中 digest头部的值。
	// 我们检查该值的形式是否为“SHA-256=<hash>”，
	// 如果不是以“SHA-256=”开头我们返回空字符串，否则返回其后的部分。

	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Println("missing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c, e := storeObject(r.Body, url.PathEscape(hash))
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
	// 同样，GetSizeFromHeader也是调用h.Get获取“content-length”头部，
	// 并调用strconv.ParseInt将字符串转化为int64输出。
	// strconv.ParseInt和例3-6中 strconv.Atoi这两个函数的作用都是将一个字符串转换成一个数字。
	// 它们的区别在于ParseInt返回的类型是int64而Atoi返回的类型是int，
	// 且.ParseInt 的功能更加复杂，它额外的输入参数用于指定转换时的进制和结果的比特长度。
	// 比如说ParseInt可以将一个字符串“OxFF”以十六进制的方式转换为整数255，
	// 而Atoi则只能将字符串“255”转换为整数255。

	size := utils.GetSizeFromHeader(r.Header)
	e = es.AddVersion(name, hash, size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// 在第2章中，我们以<object_name>为参数调用storeObject。
// 而本章我们首先调用utils.GetHashFromHeader 从 HTTP请求头部获取对象的散列值，
// 然后以散列值为参数调用storeObject。
// 之后我们从URL中获取对象的名字并且调用utils.GetSizeFromHeader从HTTP请求头部获取对象的大小，
// 然后以对象的名字、散列值和大小为参数调用es.AddVersion给该对象添加新版本。
