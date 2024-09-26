package lru

import (
	"testing"
)

type Integer int
type String string

func (i Integer) Size() int {
	return 4
}
func (s String) Size() int {
	return len(s)
}

func TestGet(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	cap := len(k1 + k2 + v1 + v2)
	lru := New(int64(cap), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	if _, ok := lru.Get("key1"); ok {
		t.Fatalf("Removeoldest key1 failed")
	}
}
