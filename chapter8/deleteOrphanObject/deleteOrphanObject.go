package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"../../src/lib/es"
)

// 我们用于删除对象数据的工具叫作 deleteOrphanObject,

// deleteOrphanObject程序需要在每一个数据服务节点上定期运行，它调用filepath.Glob
// 获取$STORAGE ROOT/objects/目录下所有文件，并在 for循环中遍历访问这些文件，
// 从文件名中获得对象的散列值，并调用es.HasHash检查元数据服务中是否存在该散列值。
// 如果不存在,则调用 del 删除散列值。

func main() {
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")

	for i := range files {
		hash := strings.Split(filepath.Base(files[i]), ".")[0]
		hashInMetadata, e := es.HasHash(hash)
		// es包的 HasHash 函数通过ES的search API搜索所有对象元数据中hash属性等于散列值的文档，
		// 如果满足条件的文档数量不为0，说明还存在对该散列值的引用，函数返回 true,否则返回 false。
		if e != nil {
			log.Println(e)
			return
		}
		if !hashInMetadata {
			del(hash)
			// del 函数访问数据服务的DELETE对象接口进行散列值的删除。

		}
	}
}

func del(hash string) {
	log.Println("delete", hash)
	url := "http://" + os.Getenv("LISTEN_ADDRESS") + "/objects/" + hash
	request, _ := http.NewRequest("DELETE", url, nil)
	client := http.Client{}
	client.Do(request)
}
