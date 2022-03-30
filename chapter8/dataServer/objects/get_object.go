package objects

// 数据服务的 objects包
// objects包发生变化的只有一个getFile函数,

import (
	"crypto/sha256"
	"encoding/base64"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"../locate"
)

func getFile(name string) string {
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/" + name + ".*")
	if len(files) != 1 {
		return ""
	}
	file := files[0]
	h := sha256.New()
	sendFile(h, file)
	d := url.PathEscape(base64.StdEncoding.EncodeToString(h.Sum(nil)))
	hash := strings.Split(file, ".")[2]
	if d != hash {
		log.Println("object hash mismatch, remove", file)
		locate.Del(hash)
		os.Remove(file)
		return ""
	}
	return file
}

// 数据服务对象接口使用的对象名的格式是<hash>.X，
// getFile函数需要在$STORAGEROOT/objects/目录下查找所有以<hash>.X开头的文件，
// 如果找不到则返回空字符串。找到之后计算其散列值，
// 如果跟<hash of shard X>的值不匹配则删除该对象并返回空字符串,否则返回该对象的文件名。
