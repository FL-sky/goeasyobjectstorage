package temp

import (
	"compress/gzip"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"

	"../../../src/lib/utils"
	"../locate"
)

// 数据服务有两个地方发生了改动，首先是用于将临时对象转正的commitTempObject函数

func (t *tempInfo) hash() string {
	s := strings.Split(t.Name, ".")
	return s[0]
}

func (t *tempInfo) id() int {
	s := strings.Split(t.Name, ".")
	id, _ := strconv.Atoi(s[1])
	return id
}

// 在本章之前的实现代码中，commitTempObject使用os.Rename将临时对象文件重命名为正式对象文件。
// 本章的实现首先用os.Create创建正式对象文件w，
// 然后以w为参数调用gzip.NewWriter 创建w2，
// 然后将临时对象文件f中的数据复制进w2，最后删除临时对象文件并添加对象定位缓存。

func commitTempObject(datFile string, tempinfo *tempInfo) {
	f, _ := os.Open(datFile)
	defer f.Close()
	d := url.PathEscape(utils.CalculateHash(f))
	f.Seek(0, io.SeekStart)
	w, _ := os.Create(os.Getenv("STORAGE_ROOT") + "/objects/" + tempinfo.Name + "." + d)
	w2 := gzip.NewWriter(w)
	io.Copy(w2, f)
	w2.Close()
	os.Remove(datFile)
	locate.Add(tempinfo.hash(), tempinfo.id())
}
