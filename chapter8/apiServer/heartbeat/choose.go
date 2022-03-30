package heartbeat

import (
	"math/rand"
)

// 接口服务的heartbeat包也需要进行改动，
// 将原来的返回一个随机数据服务节点的ChooseRandomDataServer函数
// 改为能够返回多个随机数据服务节点的 ChooseRandomDataServers函数，见例5-2。

func ChooseRandomDataServers(n int, exclude map[int]string) (ds []string) {
	candidates := make([]string, 0)
	reverseExcludeMap := make(map[string]int)
	for id, addr := range exclude {
		reverseExcludeMap[addr] = id
	}
	servers := GetDataServers()
	for i := range servers {
		s := servers[i]
		_, excluded := reverseExcludeMap[s]
		if !excluded {
			candidates = append(candidates, s)
		}
	}
	length := len(candidates)
	if length < n {
		return
	}
	p := rand.Perm(length) // 打散[0,1,...,length-1]
	for i := 0; i < n; i++ {
		ds = append(ds, candidates[p[i]])
	}
	return
}

// ChooseRandomDataServers函数有两个输入参数，
// 整型n表明我们需要多少个随机数据服务节点，
// exclude参数的作用是要求返回的随机数据服务节点不能包含哪些节点。
// 这是因为当我们的定位完成后，实际收到的反馈消息有可能不足6个，
// 此时我们需要进行数据修复，根据目前已有的分片将丢失的分片复原出来并再次上传到数据服务，
// 所以我们需要调用ChooseRandomDataServers函数来获取用于上传复原分片的随机数据服务节点。
// 很显然，目前已有的分片所在的数据服务节点需要被排除。

// exclude的类型和之前locate.Locate 的输出参数locateInfo一致，其map 的值是数据服务节点地址。
// 但是当我们实现查找算法时这个数据使用起来并不方便，所以需要进行一下键值转换。
// 转换后的reverseExcludeMap 以地址为键，
// 这样我们在后面遍历当前所有数据节点时就可以更容易检查某个数据节点是否需要被排除，
// 不需要被排除的加入candidates数组。如果最后得到的candidates 数组长度length小于n，
// 那么我们无法满足要求的n个数据服务节点,返回一个空数组，
// 否则调用rand.Perm将0到 length-1 的所有整数乱序排列返回一个数组,
// 取前n个作为candidates数组的下标取数据节点地址返回。
