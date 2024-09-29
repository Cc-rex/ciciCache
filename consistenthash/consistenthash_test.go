package consistenthash

import (
	"hash/crc32"
	"log"
	"sort"
	"testing"
)

func TestRegister(t *testing.T) {
	c := New(2, nil)
	c.Register("peer1", "peer2")
	// test the number of virtual node
	if len(c.Ring) != 4 {
		t.Errorf("the virtual node is wrong")
	}
	// test the hashValue
	hashValue := int(crc32.ChecksumIEEE([]byte("0peer2")))
	idx := sort.SearchInts(c.Ring, hashValue)
	if c.Ring[idx] != hashValue {
		t.Errorf("Actual: %d\tExpect: %d\n", c.Ring[idx], hashValue)
	}
}

func TestGetPeer(t *testing.T) {
	c := New(2, nil)
	c.Register("peer1", "peer2")
	key := "cc"
	keyHashValue := int(crc32.ChecksumIEEE([]byte(key)))
	log.Printf("key hash = %d\n", keyHashValue)
	for _, v := range c.Ring {
		log.Printf("%d -> %s\n", v, c.Hashmap[v])
	}
	peer := c.GetPeer(key)
	log.Printf("%s is getpeer", peer)
}
