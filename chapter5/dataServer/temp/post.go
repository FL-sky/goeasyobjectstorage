package temp

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type tempInfo struct {
	Uuid string
	Name string
	Size int64
}

// 结构体tempInfo用于记录临时对象的uuid、名字和大小。
// post函数用于处理HTTP请求，它会生成一个随机的 uuid，从请求的URL 获取对象的名字，也是散列值。
// 从Size头部读取对象的大小,然后拼成一个tempInfo结构体,
// 调用tempInfo的 writeToFile 方法将该结构体的内容写入磁盘上的文件。
// 然后它还会在$STORAGE_ROOT/temp/目录里创建一个名为<uuid>.dat的文件
// (<uuid>为实际生成的uuid的值)，用于保存该临时对象的内容，最后将该uuid通过HTTP响应返回给接口服务。

func post(w http.ResponseWriter, r *http.Request) {
	output, _ := exec.Command("uuidgen").Output()
	uuid := strings.TrimSuffix(string(output), "\n")
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	size, e := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t := tempInfo{uuid, name, size}
	e = t.writeToFile()
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	os.Create(os.Getenv("STORAGE_ROOT") + "/temp/" + t.Uuid + ".dat")
	w.Write([]byte(uuid))
}

// tempInfo的 writeToFile方法会在$STORAGE_ROOT/temp/目录里创建一个名为<uuid>的文件，
// 并将自身的内容经过JSON 编码后写入该文件。
// 注意，这个文件是用于保存临时对象信息的,跟用于保存对象内容的<uuid>.dat是不同的两个文件。

func (t *tempInfo) writeToFile() error {
	f, e := os.Create(os.Getenv("STORAGE_ROOT") + "/temp/" + t.Uuid)
	if e != nil {
		return e
	}
	defer f.Close()
	b, _ := json.Marshal(t)
	f.Write(b)
	return nil
}

// 接口服务在调用了POST 方法之后会从数据服务获得一个uuid，这意味着数据服务已经为这个临时对象做好了准备。
// 之后接口服务还需要继续调用PATCH方法将数据上传，PATCH方法相关函数见例4-7。
