package temp

// temp包的改动主要在处理临时对象转正时，也就是commitTempObject函数,

import (
	"net/url"
	"os"
	"strconv"
	"strings"

	"../../../src/lib/utils"
	"../locate"
)

func (t *tempInfo) hash() string {
	s := strings.Split(t.Name, ".")
	return s[0]
}

func (t *tempInfo) id() int {
	s := strings.Split(t.Name, ".")
	id, _ := strconv.Atoi(s[1])
	return id
}

func commitTempObject(datFile string, tempinfo *tempInfo) {
	f, _ := os.Open(datFile)
	d := url.PathEscape(utils.CalculateHash(f))
	f.Close()
	os.Rename(datFile, os.Getenv("STORAGE_ROOT")+"/objects/"+tempinfo.Name+"."+d)
	locate.Add(tempinfo.hash(), tempinfo.id())
}

// 我们回顾一下第4章，commitTempObject的实现非常简单，
// 只需要将临时对象的数据文件重命名为SSTORAGE_ROOT/objects/<hash>，<hash>是该对象的散列值。
// 而在本章，正式对象文件名是$STORAGE_ROOTlobjects/<hash>.X.<hash of shard X>。
// 所以在重命名时，commitTemp Object函数需要读取临时对象的数据并计算散列值<hash ofshard X>。
// 最后，我们调用locate.Add,以<hash>为键、分片的id为值添加进定位缓存。
