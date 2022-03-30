package objects

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"../locate"
)

func del(w http.ResponseWriter, r *http.Request) {
	hash := strings.Split(r.URL.EscapedPath(), "/")[2]
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/" + hash + ".*")
	if len(files) != 1 {
		return
	}
	locate.Del(hash)
	os.Rename(files[0], os.Getenv("STORAGE_ROOT")+"/garbage/"+filepath.Base(files[0]))
}

// del函数根据对象散列值搜索对象文件，调用 1ocate.Del将该散列值移出对象定位缓存，
// 并调用os.Rename将对象文件移动到$STORAGE_ROOT/garbage/目录下。

// $STORAGE_ROOT/garbage/目录下的文件需要定期检查，在超过一定时间后可以彻底删除，
// 在彻底删除前还要再次确认元数据服务中不存在相关散列值,如果真的发生竞争，散列值存在，
// 我们还需要将该对象重新上传一次。本书没有实现能够二次检查并重新上传的工具软件,留给有兴趣的读者自行实现。
