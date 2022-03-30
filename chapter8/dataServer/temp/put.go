package temp

import (
	"log"
	"net/http"
	"os"
	"strings"
)

// 和patch函数类似，put函数一开始也是获取uuid，打开信息文件读取tempInfo结构体，
// 打开数据文件读取临时象大小并进行比较，如果大小一致，则调用commitTempObject将临时对象转正。

func put(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	tempinfo, e := readFromFile(uuid)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	infoFile := os.Getenv("STORAGE_ROOT") + "/temp/" + uuid
	datFile := infoFile + ".dat"
	f, e := os.Open(datFile)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	info, e := f.Stat()
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	actual := info.Size()
	os.Remove(infoFile)
	if actual != tempinfo.Size {
		os.Remove(datFile)
		log.Println("actual size mismatch, expect", tempinfo.Size, "actual", actual)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	commitTempObject(datFile, tempinfo)
}

// commitTempObject函数会调用os.Rename 将临时对象的数据文件改名为$STORAGEROOT/objects/<hash>。
// <hash>是对象的名字，也是散列值。之后还会调用 locate.Add将<hash>加入数据服务的对象定位缓存。

// DELETE方法相关函数见例4-9。
// 例4-9temp.del函数
