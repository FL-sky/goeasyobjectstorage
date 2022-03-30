package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"../../src/lib/es"
	"../../src/lib/utils"
	"../apiServer/objects"
)

// 我们用于检查和修复对象数据的工具叫作 objectScanner,

// objectScanner也需要在数据服务节点上定期运行，它调用filepath.Glob 获取$STORAGE_ROOT/objects/目录下所有文件，
// 并在for 循环中遍历访问这些文件，从文件名中获得对象的散列值，并调用verify检查数据。

func main() {
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")

	for i := range files {
		hash := strings.Split(filepath.Base(files[i]), ".")[0]
		verify(hash)
		// verify函数调用es.SearchHashSize从元数据服务中获取该散列值对应的对象大小,
		// 然后以对象的散列值和大小为参数调用objects.GetStream创建一个对象数据流，
		// 并调用utils.Calculate Hash计算对象的散列值，检查是否一致。
		// 如果不一致，需要以log的形式报告错误。最后调用stream.Close 关闭数据对象流。

	}
}

func verify(hash string) {
	log.Println("verify", hash)
	size, e := es.SearchHashSize(hash)
	// es包的SearchHashSize函数的输入参数是对象的散列值 hash，
	// 它通过ES 的search API查询元数据属性中hash等于该散列值的文档的size属性，并返回这个size的值。

	if e != nil {
		log.Println(e)
		return
	}
	stream, e := objects.GetStream(hash, size)
	// 我们在第5章中已经介绍过 objects.GetStream，它会创建一个指向rs.RSGetStream结构体的指针。
	// 通过读取 rs.RSGetStream并在最后关闭它，底层的实现会自动完成数据的修复。
	// 如果数据已经损坏得不可修复，那么在计算散列值的时候就不可能匹配,
	// 我们可以根据log打印的报告发现。

	if e != nil {
		log.Println(e)
		return
	}
	d := utils.CalculateHash(stream)
	if d != hash {
		log.Printf("object hash mismatch, calculated=%s, requested=%s", d, hash)
	}
	stream.Close()
}
