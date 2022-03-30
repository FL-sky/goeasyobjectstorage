package temp

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"../../../src/lib/es"
	"../../../src/lib/rs"
	"../../../src/lib/utils"
	"../locate"
)

func put(w http.ResponseWriter, r *http.Request) {
	// put函数首先从 URL 中获取<token>，然后调用rs.NewRSResumablePutStreamFromToken
	// 根据<token>中的内容创建RSResumablePutStream结构体并获得指向该结构体的指针 stream，
	// 然后调用CurrentSize方法获得token当前大小,如果大小为-1，则说明该token不存在。
	// 接下来我们调用utils.GetOffsetFromHeader 从 Range头部获得 offset。
	// 如果 offset和当前的大小不一致，则接口服务返回416 Range Not Satisfiable

	token := strings.Split(r.URL.EscapedPath(), "/")[2]
	stream, e := rs.NewRSResumablePutStreamFromToken(token)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	current := stream.CurrentSize()
	if current == -1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	offset := utils.GetOffsetFromHeader(r.Header)
	if current != offset {
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}
	// 如果offset和当前大小一致,我们在一个 for循环中以32000字节为长度读取HTTP请求的正文并写入stream。
	// 如果读到的总长度超出了对象的大小，说明客户端上传的数据有误，接口服务删除临时对象并返回 403 Forbidden。
	// 如果某次读取的长度不到32 000字节且读到的总长度不等于对象的大小，说明本次客户端上传结束，还有后续数据需要上传。
	// 此时接口服务会丢弃最后那次读取的长度不到32 000字节的数据。

	bytes := make([]byte, rs.BLOCK_SIZE)
	for {
		n, e := io.ReadFull(r.Body, bytes)
		if e != nil && e != io.EOF && e != io.ErrUnexpectedEOF {
			log.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		current += int64(n)
		if current > stream.Size {
			stream.Commit(false)
			log.Println("resumable put exceed size")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if n != rs.BLOCK_SIZE && current != stream.Size {
			return
		}

		// 为什么接口服务需要丢弃数据，而不是将这部分数据写入临时对象或缓存在接口服务的内存里?
		// 因为将这部分数据缓存在接口服务的内存里没有意义，下次客户端不一定还访问同一个接口服务节点。
		// 而如果我们将这部分数据直接写入临时对象，那么我们就破坏了每个数据片以8000字节为一个块写入的约定，在读取时就会发生错误。

		stream.Write(bytes[:n])

		// 最后如果读到的总长度等于对象的大小，说明客户端上传了对象的全部数据。
		// 我们调用 stream 的 Flush方法将剩余数据写入临时对象，
		// 然后调用rs.NewRSResumableGetStream生成一个临时对象读取流getStream，读取getStream中的数据并计算散列值。
		// 如果散列值不一致,则说明客户端上传的数据有误，接口服务删除临时对象并返回403Forbidden。
		// 如果散列值一致，则继续检查该散列值是否已经存在，如果存在，则删除临时对象;否则将临时对象转正。
		// 最后调用es.AddVersion添加新版本。

		if current == stream.Size {
			stream.Flush()
			getStream, e := rs.NewRSResumableGetStream(stream.Servers, stream.Uuids, stream.Size)
			hash := url.PathEscape(utils.CalculateHash(getStream))
			if hash != stream.Hash {
				stream.Commit(false)
				log.Println("resumable put done but hash mismatch")
				w.WriteHeader(http.StatusForbidden)
				return
			}
			if locate.Exist(url.PathEscape(hash)) {
				stream.Commit(false)
			} else {
				stream.Commit(true)
			}
			e = es.AddVersion(stream.Name, stream.Hash, stream.Size)
			if e != nil {
				log.Println(e)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
	}
}
