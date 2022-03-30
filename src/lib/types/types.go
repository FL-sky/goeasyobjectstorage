package types

type LocateMessage struct {
	Addr string
	Id   int
}

// 由于该结构体需要同时被接口服务的数据服务引用，所以放在types包里。
