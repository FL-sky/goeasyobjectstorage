package objects

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"../../../src/lib/utils"

	"../locate"
)

func storeObject(r io.Reader, hash string, size int64) (int, error) {
	// storeObject函数首先调用locate.Exist定位对象的散列值，
	// 如果已经存在，则跳过后续上传操作直接返回200. OK;否则调用putStream 生成对象的写入流stream用于写入。
	// 注意，这里进行定位的散列值和之后作为参数调用putStream 的散列值都经过url.PathEscape的处理，
	// 原因在之前的章节已经讲过，是为了确保这个散列值可以被放在URL中使用。
	if locate.Exist(url.PathEscape(hash)) {
		return http.StatusOK, nil
	}

	stream, e := putStream(url.PathEscape(hash), size)
	if e != nil {
		return http.StatusInternalServerError, e
	}

	reader := io.TeeReader(r, stream)
	// io.TeeReader的功能类似Unix的tee命令。它有两个输入参数,
	// 分别是作为io.Reader的r和作为io. Writer 的stream，它返回的reader也是一个io.Reader。
	// 当reader被读取时，其实际的内容读取自r，同时会写入, stream。
	// 我们用utils.CalculateHash 从reader中读取数据的同时也写入了stream。

	d := utils.CalculateHash(reader)
	// utils.CalculateHash函数调用sha256.New生成的变量h，类型是 sha256.digest结构体，
	// 实现的接口则是hash.Hash。io.Copy 从参数r中读取数据并写入h，h会对写入的数据计算其散列值，
	// 这个散列值可以通过h.Sum方法读取。
	// 我们从h.Sum读取到的散列值是一个二进制的数据，
	// 还需要用 base64.StdEncoding.Encode ToString函数进行Base64编码,
	// 然后跟对象的散列值hash进行比较,
	// 如果不一致,则调用stream. Commit(false)删除临时对象，并返回400 Bad Request:;
	// 如果一致，则调用stream.Commit (true)将临时对象转正并返回200 OK。

	if d != hash {
		stream.Commit(false)
		return http.StatusBadRequest, fmt.Errorf("object hash mismatch, calculated=%s, requested=%s", d, hash)
	}
	stream.Commit(true)
	return http.StatusOK, nil
}
