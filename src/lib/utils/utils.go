package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// utils.GetOffsetFromHeader 函数获取 HTTP 的 Range头部，Range头部的格式必须是“bytes=<first>-”开头，
// 我们调用strings.Split将<first>部分切取出来并调用strconv.ParseInt 将字符串转化成int64返回。

func GetOffsetFromHeader(h http.Header) int64 {
	byteRange := h.Get("range")
	if len(byteRange) < 7 {
		return 0
	}
	if byteRange[:6] != "bytes=" {
		return 0
	}
	bytePos := strings.Split(byteRange[6:], "-")
	offset, _ := strconv.ParseInt(bytePos[0], 0, 64)
	return offset
}

func GetHashFromHeader(h http.Header) string {
	digest := h.Get("digest")
	if len(digest) < 9 {
		return ""
	}
	if digest[:8] != "SHA-256=" {
		return ""
	}
	return digest[8:]
}

func GetSizeFromHeader(h http.Header) int64 {
	size, _ := strconv.ParseInt(h.Get("content-length"), 0, 64)
	return size
}

func CalculateHash(r io.Reader) string {
	h := sha256.New()
	io.Copy(h, r)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
