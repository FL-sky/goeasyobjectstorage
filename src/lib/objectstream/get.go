package objectstream

import (
	"fmt"
	"io"
	"net/http"
)

type GetStream struct {
	reader io.Reader
}

// GetStream 比 PutStream 简单很多，因为Go语言的 http包会返回一个io.Reader,我们可以直接从中读取响应的正文，
// 而不需要像PutStream那样使用管道来适配。所以我们的GetStream只需要一个成员reader用于记录http返回的io.Reader。

func newGetStream(url string) (*GetStream, error) {
	r, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dataServer return http code %d", r.StatusCode)
	}
	return &GetStream{r.Body}, nil
}

// newGetStream 函数的输入参数url是一个字符串，表示用于获取数据流的 HTTP服务地址。
// 我们调用http.Get发起一个GET请求，获取该地址的HTTP响应。
// http.Get返回的r的类型是*http.Response，其成员StatusCode是该HTTP 响应的错误代码，
// 成员 Body则是用于读取 HTTP响应正文的 io.Reader。我们将 r.Body作为新生成的GetStream的reader，并返回这个 GetStream。

func NewGetStream(server, object string) (*GetStream, error) {
	if server == "" || object == "" {
		return nil, fmt.Errorf("invalid server %s object %s", server, object)
	}
	return newGetStream("http://" + server + "/objects/" + object)
}

// NewGetStream是newGetStream 的封装函数。newGetStream首字母小写，说明该函数并没有暴露给objectstream包外部使用，
// NewGetStream 的函数签名只需要server和 object这两个字符串，它们会在函数内部拼成一个url传给newGetStream，
// 这样，对外就隐藏了url 的细节。使用者不需要知道具体的url，只需要提供数据服务节点地址和对象名就可以读取对象。

func (r *GetStream) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}

// GetStream.Read方法用于读取 reader成员。只要实现了该方法，我们的 GetStream结构体就实现了io.Reader接口。
// 这也是为什么NewGetStream函数第一个返回值的类型是*GetStream，
// 而例2-9中，objects.getStream函数第一个返回值的类型却是io.Reader 。
///（chapter 2 例2-9 接口服务的 objects.get相关函数 getStream）
