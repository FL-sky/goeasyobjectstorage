package objects

import (
	"compress/gzip"
	"io"
	"log"
	"os"
)

func sendFile(w io.Writer, file string) {
	f, e := os.Open(file)
	if e != nil {
		log.Println(e)
		return
	}
	defer f.Close()
	gzipStream, e := gzip.NewReader(f)
	if e != nil {
		log.Println(e)
		return
	}
	io.Copy(w, gzipStream)
	gzipStream.Close()
}

// 本章不再直接用io.Copy 读取对象文件，而是先在对象文件上用 gzip.NewReader
// 创建一个指向gzip.Reader 结构体的指针gzipStream,然后读取 gzipStream中的数据。
// 通过这种方式，对象文件f中的数据会先被gzip解压，然后才被读取出来
