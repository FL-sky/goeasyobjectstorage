package temp

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// patch函数首先获取请求URL的<uuid>部分,然后从相关信息文件中读取 tempInfo结构体，
// 如果找不到相关的信息文件，我们就返回404 Not Found;
// 如果相关信息文件存在，则用os.OpenFile打开临时对象的数据文件，并用io.Copy将请求的正文写入数据文件。
// 写完后调用f.Stat方法获取数据文件的信息 info ，用info.Size获取数据文件当前的大小，
// 如果超出了tempInfo中记录的大小，我们就删除信息文件和数据文件并返回500 Internal Server Error。

func patch(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	tempinfo, e := readFromFile(uuid)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	infoFile := os.Getenv("STORAGE_ROOT") + "/temp/" + uuid
	datFile := infoFile + ".dat"
	f, e := os.OpenFile(datFile, os.O_WRONLY|os.O_APPEND, 0)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	_, e = io.Copy(f, r.Body)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	info, e := f.Stat()
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	actual := info.Size()
	if actual > tempinfo.Size {
		os.Remove(datFile)
		os.Remove(infoFile)
		log.Println("actual size", actual, "exceeds", tempinfo.Size)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// readFromFile函数的输入参数是uuid，它用os.Open打开$STORAGE_ROOT/temp/<uuid>文件，
// 读取其全部内容并经过JSON解码成一个tempInfo结构体返回。

func readFromFile(uuid string) (*tempInfo, error) {
	f, e := os.Open(os.Getenv("STORAGE_ROOT") + "/temp/" + uuid)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	b, _ := ioutil.ReadAll(f)
	var info tempInfo
	json.Unmarshal(b, &info)
	return &info, nil
}

// 接口服务调用PATCH方法将整个临时对象上传完毕后,自己也已经完成了数据校验的工作，
// 根据数据校验的结果决定是调用 PUT 方法将临时文件转正还是调用DELETE方法删除临时文件，
// PUT方法相关函数见例4-8。
