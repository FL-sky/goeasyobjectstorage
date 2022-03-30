package locate

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"../../../src/lib/rabbitmq"
	"../../../src/lib/types"
)

// 由于在磁盘上保存的对象文件名格式发生了变化，我们的locate包也有相应的变化,

var objects = make(map[string]int)
var mutex sync.Mutex

func Locate(hash string) int {
	mutex.Lock()
	id, ok := objects[hash]
	mutex.Unlock()
	if !ok {
		return -1
	}
	return id
}

// 相比第4章，我们的Locate函数不仅要告知某个对象是否存在，
// 同时还需要告知本节点保存的是该对象哪个分片，所以我们返回一个整型，用于返回分片的id。
// 如果对象不存在，则返回-1。

func Add(hash string, id int) {
	mutex.Lock()
	objects[hash] = id
	mutex.Unlock()
}

// Add函数用于将对象及其分片id加入缓存。

func Del(hash string) {
	mutex.Lock()
	delete(objects, hash)
	mutex.Unlock()
}

// Del函数未发生变化,故未打印。

func StartLocate() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	q.Bind("dataServers")
	c := q.Consume()
	for msg := range c {
		hash, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		id := Locate(hash)
		if id != -1 {
			q.Send(msg.ReplyTo, types.LocateMessage{Addr: os.Getenv("LISTEN_ADDRESS"), Id: id})
		}
	}
}

// StartLocate函数读取来自接口服务需要定位的对象散列值hash后，调用Locate获得分片id，
// 如果id不为-1，则将自身的节点监听地址和id打包成一个types.LocateMessage结构体作为反馈消息发送。
// types.LocateMessage的定义比较简单，见例5-12。

func CollectObjects() {
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")
	for i := range files {
		file := strings.Split(filepath.Base(files[i]), ".")
		if len(file) != 3 {
			panic(files[i])
		}
		hash := file[0]
		id, e := strconv.Atoi(file[1])
		if e != nil {
			panic(e)
		}
		objects[hash] = id
	}
}

// CollectObjects函数调用filepath.Glob获取SSTORAGE_ROOT/objects/目录下所有文件，
// 并以‘.’分割其基本文件名，获得对象的散列值hash 以及分片id,加入定位缓存。
