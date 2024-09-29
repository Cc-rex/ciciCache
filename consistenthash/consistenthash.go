package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// HashFunc 定义hash函数的输入输出
type HashFunc func(data []byte) uint32

// Consistency 维护peer与其hash值的关联
type Consistency struct {
	Hash        HashFunc       // hash
	VirtualNode int            // the number of virtual node, preventing the data skew
	Ring        []int          // hash ring
	Hashmap     map[int]string // hashValue -> peerName
}

// Register 将各个peer注册到哈希环上
func (c *Consistency) Register(peers ...string) {
	if peers == nil {
		return
	}
	for _, peer := range peers {
		for i := 0; i < c.VirtualNode; i++ {
			hashValue := int(c.Hash([]byte(strconv.Itoa(i) + peer)))
			c.Ring = append(c.Ring, hashValue)
			c.Hashmap[hashValue] = peer
		}
	}
	sort.Ints(c.Ring)
}

// GetPeer 计算key应该被缓存到哪个peer上
func (c *Consistency) GetPeer(key string) string {
	if len(c.Ring) == 0 {
		return ""
	}
	keyHash := int(c.Hash([]byte(key)))
	index := sort.Search(len(c.Ring), func(i int) bool {
		return c.Ring[i] > keyHash
	})
	return c.Hashmap[c.Ring[index%len(c.Ring)]] // 保证index超过ring长度可以回到数组起点
}

// New the constructor of the consistency
func New(vNode int, fn HashFunc) *Consistency {
	c := &Consistency{
		Hash:        fn,
		VirtualNode: vNode,
		Ring:        []int{},
		Hashmap:     map[int]string{},
	}
	if c.Hash == nil {
		c.Hash = crc32.ChecksumIEEE
	}
	return c
}
