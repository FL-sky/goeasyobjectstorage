package objects

import (
	"log"
	"net/url"
	"os"

	"../../../src/lib/utils"
	"../locate"
)

func getFile(hash string) string {
	file := os.Getenv("STORAGE_ROOT") + "/objects/" + hash
	f, _ := os.Open(file)
	d := url.PathEscape(utils.CalculateHash(f))
	f.Close()
	if d != hash {
		log.Println("object hash mismatch, remove", file)
		locate.Del(hash)
		os.Remove(file)
		return ""
	}
	return file
}

// getFile函数的输入参数是对象的散列值<hash>,
// 它根据这个参数找到$STORAGEROOT/objects/<hash>对象文件，
// 然后对这个对象的内容计算SHA-256散列值，并用url.PathEscape转义，最后得到的就是可用于URL 的散列值字符串。
// 我们将该字符串和对象的散列值进行比较，
// 如果不一致则打印错误日志，并从缓存和磁盘上删除对象,返回空字符串;如果一致则返回对象的文件名。

//...
// 有读者可能要质疑这里的数据校验没有必要,因为在对象上传的时候已经在接口服务进行过数据校验了。
// 事实上这里的数据校验是用于防止存储系统的数据降解，
// 哪怕在上传时正确的数据也有可能随着时间的流逝而逐渐发生损坏，我们会在下一章介绍数据降解的成因。

// 对数据安全有要求的读者可能会进一步要求在临时文件转正时进行一次数据校验，
// 以此来确保从接口服务传输过来的数据没有发生损坏。然而这一步骤仅在本章可行。
// 我们在本章开头也讲过，随着我们的系统功能不断完善，
// 最终保存在数据服务节点上的对象数据和用户的对象数据可能截然不同,
// 我们无法根据用户对象的散列值校验数据服务节点上的对象数据。
