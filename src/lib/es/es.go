package es

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Metadata struct {
	Name    string
	Version int
	Size    int64
	Hash    string
}

type hit struct {
	Source Metadata `json:"_source"`
}

type searchResult struct {
	Hits struct {
		Total int
		Hits  []hit
	}
}

func getMetadata(name string, versionId int) (meta Metadata, e error) {
	url := fmt.Sprintf("http://%s/metadata/objects/%s_%d/_source",
		os.Getenv("ES_SERVER"), name, versionId)
	r, e := http.Get(url)
	if e != nil {
		return
	}
	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("fail to get %s_%d: %d", name, versionId, r.StatusCode)
		return
	}
	result, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(result, &meta)
	return
}

// getMetadata用于根据对象的名字和版本号来获取对象的元数据，其URL 中的服务器地址来自环境变量 ES_SERVER,
// 索引是metadata,，类型是 objects,文档的id由对象的名字和版本号拼接而成。
// 通过这种方式GET这个URL可以直接获取该对象的元数据，这样就免除了耗时的搜索操作。
// ES 返回的结果经过JSON解码后被es、SearchLatestVerson函数的实现见例3-8，保存在meta变量返回。
// meta 的类型是Metadata结构体，其结构和ES映射中定义的objects类型的属性一一对应，
// 同样是包含Name、Version、Size和 Hash。

func SearchLatestVersion(name string) (meta Metadata, e error) {
	url := fmt.Sprintf("http://%s/metadata/_search?q=name:%s&size=1&sort=version:desc",
		os.Getenv("ES_SERVER"), url.PathEscape(name))
	r, e := http.Get(url)
	if e != nil {
		return
	}
	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("fail to search latest metadata: %d", r.StatusCode)
		return
	}
	result, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	if len(sr.Hits.Hits) != 0 {
		meta = sr.Hits.Hits[0].Source
	}
	return
}

// 例3-8显示了es 包的SearchLatestVersion函数，它以对象的名字为参数，调用ES搜索API。
// 它在URL中指定了对象的名字，且版本号以降序排列只返回第一个结果。
// ES返回的结果被JSON解码到一个 searchResult 结构体，这个结构体和ES搜索API返回的结构保持一致，
// 以方便我们读取搜索到的元数据并赋值给meta返回。
// 如果ES返回的结果长度为0，说明没有搜到相对应的元数据，我们直接返回。此时meta中各
// 属性都为初始值:字符串为空字符串“”，整型为0。

func GetMetadata(name string, version int) (Metadata, error) {
	if version == 0 {
		return SearchLatestVersion(name)
	}
	return getMetadata(name, version)
}

// GetMetadata函数的功能类似getMetadata，输入对象的名字和版本号返回对象，
// 区别在于当version为0时，我们会调用SearchLatestVersion获取当前最新的版本。

func PutMetadata(name string, version int, size int64, hash string) error {
	doc := fmt.Sprintf(`{"name":"%s","version":%d,"size":%d,"hash":"%s"}`,
		name, version, size, hash)
	client := http.Client{}
	url := fmt.Sprintf("http://%s/metadata/objects/%s_%d?op_type=create",
		os.Getenv("ES_SERVER"), name, version)
	// 我们使用了op_type=create参数，如果向时有多个客户端上传同一个元数据,结果会发生冲突，
	// 只有第一个文档被成功创建。之后的PUT请求，ES会返回409 Conflict。
	// 此时，我们的函数会让版本号加1并递归调用自身继续上传。

	request, _ := http.NewRequest("PUT", url, strings.NewReader(doc))
	r, e := client.Do(request)
	if e != nil {
		return e
	}
	if r.StatusCode == http.StatusConflict {
		return PutMetadata(name, version+1, size, hash)
	}
	if r.StatusCode != http.StatusCreated {
		result, _ := ioutil.ReadAll(r.Body)
		return fmt.Errorf("fail to put metadata: %d %s", r.StatusCode, string(result))
	}
	return nil
}

// PutMetadata 函数用于向ES服务上传一个新的元数据。它的4个输入参数对应元数据的4个属性，
// 函数会将它们拼成一个ES文档，一个ES的文档相当于数据库的-条记录。
// 我们用PUT方法把这个文档上传到metadata索引的objects类型，且文档id由元数据的name和 version拼成，
// 方便我们GET。

func AddVersion(name, hash string, size int64) error {
	version, e := SearchLatestVersion(name)
	if e != nil {
		return e
	}
	return PutMetadata(name, version.Version+1, size, hash)
}

// Add Version函数首先调用SearchLatestVersion获取对象最新的版本，然后在版本号上加1调用PutMetadata。

func SearchAllVersions(name string, from, size int) ([]Metadata, error) {
	url := fmt.Sprintf("http://%s/metadata/_search?sort=name,version&from=%d&size=%d",
		os.Getenv("ES_SERVER"), from, size)
	if name != "" {
		url += "&q=name:" + name
	}
	r, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	metas := make([]Metadata, 0)
	result, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	for i := range sr.Hits.Hits {
		metas = append(metas, sr.Hits.Hits[i].Source)
	}
	return metas, nil
}

// SearchAllVersions函数用于搜索某个对象或所有对象的全部版本。它的输入参数name表示对象的名字，
// 如果name不为空字符串则搜索指定对象的所有版本，否则搜索所有对象的全部版本。
// 输入参数 from和 size指定分页的显示结果，其功能和ES搜索API的 from和 size参数一致。
// 搜索的结果按照对象的名字和版本号排序，并被保存在Metadata的数组里用于返回。

func DelMetadata(name string, version int) {
	client := http.Client{}
	url := fmt.Sprintf("http://%s/metadata/objects/%s_%d",
		os.Getenv("ES_SERVER"), name, version)
	request, _ := http.NewRequest("DELETE", url, nil)
	client.Do(request)
}

type Bucket struct {
	Key         string
	Doc_count   int
	Min_version struct {
		Value float32
	}
}

type aggregateResult struct {
	Aggregations struct {
		Group_by_name struct {
			Buckets []Bucket
		}
	}
}

func SearchVersionStatus(min_doc_count int) ([]Bucket, error) {
	client := http.Client{}
	url := fmt.Sprintf("http://%s/metadata/_search", os.Getenv("ES_SERVER"))
	body := fmt.Sprintf(`
        {
          "size": 0,
          "aggs": {
            "group_by_name": {
              "terms": {
                "field": "name",
                "min_doc_count": %d
              },
              "aggs": {
                "min_version": {
                  "min": {
                    "field": "version"
                  }
                }
              }
            }
          }
        }`, min_doc_count)
	request, _ := http.NewRequest("GET", url, strings.NewReader(body))
	r, e := client.Do(request)
	if e != nil {
		return nil, e
	}
	b, _ := ioutil.ReadAll(r.Body)
	var ar aggregateResult
	json.Unmarshal(b, &ar)
	return ar.Aggregations.Group_by_name.Buckets, nil
}

func HasHash(hash string) (bool, error) {
	url := fmt.Sprintf("http://%s/metadata/_search?q=hash:%s&size=0", os.Getenv("ES_SERVER"), hash)
	r, e := http.Get(url)
	if e != nil {
		return false, e
	}
	b, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(b, &sr)
	return sr.Hits.Total != 0, nil
}

func SearchHashSize(hash string) (size int64, e error) {
	url := fmt.Sprintf("http://%s/metadata/_search?q=hash:%s&size=1",
		os.Getenv("ES_SERVER"), hash)
	r, e := http.Get(url)
	if e != nil {
		return
	}
	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("fail to search hash size: %d", r.StatusCode)
		return
	}
	result, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	if len(sr.Hits.Hits) != 0 {
		size = sr.Hits.Hits[0].Source.Size
	}
	return
}
