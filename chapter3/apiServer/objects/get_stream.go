package objects

import (
	"fmt"
	"io"

	"../../../src/lib/objectstream"
	"../locate"
	// "go-implement-your-object-storage/chapter2/locate"
)

func getStream(object string) (io.Reader, error) {
	server := locate.Locate(object)
	if server == "" {
		return nil, fmt.Errorf("object %s locate fail", object)
	}
	return objectstream.NewGetStream(server, object)
}

// getStream函数的参数object是一个字符串，它代表对象的名字。我们首先调用locate.Locate定位这个对象，
// 如果返回的数据服务节点地址为空字符串，则返回定位失败的错误;objectstream.NewGetStream 并返回其结果
