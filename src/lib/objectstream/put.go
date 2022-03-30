package objectstream

import (
	"fmt"
	"io"
	"net/http"
)

type PutStream struct {
	writer *io.PipeWriter
	c      chan error
}

// PutStream是一个结构体,内含一个io.PipeWriter的指针 writer和一个error的channel c。
// writer用于实现Write方法，c用于把在一个goroutine传输数据的过程中发生的错误传回主线程。

func NewPutStream(server, object string) *PutStream {
	reader, writer := io.Pipe()
	c := make(chan error)
	go func() {
		request, _ := http.NewRequest("PUT", "http://"+server+"/objects/"+object, reader)
		client := http.Client{}
		r, e := client.Do(request)
		if e == nil && r.StatusCode != http.StatusOK {
			e = fmt.Errorf("dataServer return http code %d", r.StatusCode)
		}
		c <- e
	}()
	return &PutStream{writer, c}
}

// NewPutStream 函数用于生成一个PutStream 结构体。它用io.Pipe 创建了一对reader和 writer，
// 类型分别是*io.PipeReader和*io.PipeWriter。它们是管道互联的，写入 writer的内容可以从reader中读出来。
// 之所以要有这样的一对管道是因为我们希望能以写入数据流的方式操作HTTP的PUT请求。
// Go语言的http包在生成一个PUT请求时需要提供一个io.Reader作为http.NewRequest的参数，
// 由一个类型为 http.Client的client负责从中读取需要PUT 的内容。有了这对管道，我们就可以在满足http.NewRequest 的
// 参数要求的同时用写入 writer的方式实现PutStream的 Write方法。另外，由于管道的读写是阻塞的，
// 我们需要在一个 goroutine中调用 client.Do方法。该方法的返回值有两个:HTTP 响应的错误代码和error。
// 如果error 不等于空(nil)，说明出现了错误，我们需要把这些错误发送进 channel c。
// 如果error等于空，但是 HTTP错误代码不为200 OK，我们也需要把这种情况记录为一种错误，然后将这个错误发送进channel c。
// 之后这个错误会在PutStream.Close方法中被读取出来。

func (w *PutStream) Write(p []byte) (n int, err error) {
	return w.writer.Write(p)
}

// PutStream.Write方法用于写入 writer。只有实现了这个方法，我们的 PutStream才会被认为是实现了io.Write接口。

func (w *PutStream) Close() error {
	w.writer.Close()
	return <-w.c
}

// PutStream.Close方法用于关闭writer。这是为了让管道另一端的reader读到io.EOF,
// 否则在goroutine中运行的client.Do将始终阻塞无法返回。然后我们从c中读取发送自goroutine的错误并返回。
