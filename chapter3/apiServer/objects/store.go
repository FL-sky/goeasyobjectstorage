package objects

import (
	"io"
	"net/http"
)

func storeObject(r io.Reader, object string) (int, error) {
	stream, e := putStream(object)
	if e != nil {
		return http.StatusServiceUnavailable, e
	}

	io.Copy(stream, r)
	e = stream.Close()
	if e != nil {
		return http.StatusInternalServerError, e
	}
	return http.StatusOK, nil
}

// storeObject函数以object为参数调用putStream生成stream，stream的类型是*objectstream.PutStream，
// 这是一个指向objectstream.PutStream结构体的指针，该结构体实现了Write方法，所以它是一个io.Write接口。
// 我们用io.Copy把 HTTP请求的正文写入这个stream，然后调用stream.Close()关闭这个流。
// objectstream.PutStream 的 Close方法返回一个error，用来通知调用者在数据传输过程中发生的错误，
// 如有错误，我们返回http.StatusInternalServerError，客户端会收到HTTP错误代码500 Internal Server Error。
