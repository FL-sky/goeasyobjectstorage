package main

import (
	"log"

	"../../src/lib/es"
)

const MIN_VERSION_COUNT = 5

// 我们用于删除过期元数据的工具叫作 deleteOldMetadata,

// deleteOldMetadata的main函数很简单，
// 调用es.SearchVersionStatus将元数据服务中所有版本数量大于等于6的对象都搜索出来保存在 Bucket结构体的数组 buckets里。
// Bucket结构体的属性包括:字符串 Key，表示对象的名字;整型 Doc_count，表示该对象目前有多少个版本;Min_version结构体，
// 内含32位浮点数Value，表示对象当前最小的版本号。

func main() {
	// es包的 Search VersionStatus 函数的输入参数min_doc_count用于指示需要搜索对象的最小版本数量。
	// 它使用ElasticSearch的 aggregations search API搜索元数据，以对象的名字分组，搜索版本数量大于等于min_doc_count的对象并返回。
	// 本书不对ES 的各种API进行深入讲解，有兴趣的读者可以自行查阅ES在线文档。

	buckets, e := es.SearchVersionStatus(MIN_VERSION_COUNT + 1)
	if e != nil {
		log.Println(e)
		return
	}
	// main函数遍历buckets,并在一个for循环中调用es.DelMetadata,从该对象当前最小的版本号开始一一删除，直到最后还剩5个。

	for i := range buckets {
		bucket := buckets[i]
		for v := 0; v < bucket.Doc_count-MIN_VERSION_COUNT; v++ {
			es.DelMetadata(bucket.Key, v+int(bucket.Min_version.Value))
			// es 包的 DelMetadata 函数根据对象的名字name和版本号 version 删除相应的对象元数据。
		}
	}
}
