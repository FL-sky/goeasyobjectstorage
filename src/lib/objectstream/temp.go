package objectstream

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type TempPutStream struct {
	Server string
	Uuid   string
}

// TempPutStream结构体包含Server和 Uuid字符串。
// NewTempPutStream函数的输入参数分别是 server，hash和 size。
// server 参数表明了数据服务的节点地址，hash和size分别是对象的散列值和大小。
// 我们根据这些信息以 POST方法访问数据服务的temp接口从而获得uuid，
// 并将server 和 uuid保存在TempPutStream结构体的相应属性中返回。

func NewTempPutStream(server, object string, size int64) (*TempPutStream, error) {
	request, e := http.NewRequest("POST", "http://"+server+"/temp/"+object, nil)
	if e != nil {
		return nil, e
	}
	request.Header.Set("size", fmt.Sprintf("%d", size))
	client := http.Client{}
	response, e := client.Do(request)
	if e != nil {
		return nil, e
	}
	uuid, e := ioutil.ReadAll(response.Body)
	if e != nil {
		return nil, e
	}
	return &TempPutStream{server, string(uuid)}, nil
}

// TempPutStream.Write方法根据Server 和 Uuid属性的值，
// 以PATCH 方法访问数据服务的temp接口，将需要写入的数据上传。

func (w *TempPutStream) Write(p []byte) (n int, err error) {
	request, e := http.NewRequest("PATCH", "http://"+w.Server+"/temp/"+w.Uuid, strings.NewReader(string(p)))
	if e != nil {
		return 0, e
	}
	client := http.Client{}
	r, e := client.Do(request)
	if e != nil {
		return 0, e
	}
	if r.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("dataServer return http code %d", r.StatusCode)
	}
	return len(p), nil
}

// TempPutStream.Commit方法根据输入参数good决定用PUT还是DELETE 方法访问数据服务的temp接口。

func (w *TempPutStream) Commit(good bool) {
	method := "DELETE"
	if good {
		method = "PUT"
	}
	request, _ := http.NewRequest(method, "http://"+w.Server+"/temp/"+w.Uuid, nil)
	client := http.Client{}
	client.Do(request)
}

func NewTempGetStream(server, uuid string) (*GetStream, error) {
	return newGetStream("http://" + server + "/temp/" + uuid)
}
