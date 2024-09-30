package ciciCache

import (
	"fmt"
	"log"
	"sync"
)

// Retriever Implementing the ability of an object to fetch data from a data source
type Retriever interface {
	retrieve(string) ([]byte, error)
}

type RetrieverFunc func(key string) ([]byte, error)

// 定义一个函数类型 F，并且实现接口 A 的方法，然后在这个方法中调用自己。这是 Go 语言中将其他函数（参数返回值定义与 F 一致）转换为接口 A 的常用技巧。
func (f RetrieverFunc) retrieve(key string) ([]byte, error) {
	return f(key)
}

// A Group is a cache namespace and associated data loaded spread over
type Group struct {
	name      string
	retriever Retriever
	cache     *cache
	server    Picker
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup create a new instance of Group
func NewGroup(name string, capacity int64, retriever Retriever) *Group {
	if retriever == nil {
		panic("Group retriever must be existed!")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		retriever: retriever,
		cache:     newCache(capacity),
	}
	groups[name] = g
	return g
}

// RegisterSvr 为 Group 注册 Server
func (g *Group) RegisterSvr(p Picker) {
	if g.server != nil {
		panic("group had been registered server")
	}
	g.server = p
}

func DestroyGroup(name string) {
	g := GetGroup(name)
	if g != nil {
		svr := g.server.(*server)
		svr.Stop()
		delete(groups, name)
		log.Printf("Destroy cache [%s %s]", name, svr.addr)
	}
}

func GetGroup(name string) *Group {
	mu.RLock() // only read lock
	g := groups[name]
	mu.RUnlock()
	return g
}

// Get value for a key from cache
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key required")
	}
	if byteView, ok := g.cache.Get(key); ok {
		log.Println("cache hit")
		return byteView, nil
	}
	// cache missing, get it another way
	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	return g.getLocally(key)
}

// getLocally 本地向Retriever取回数据并填充缓存
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.retriever.retrieve(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: copyByte(bytes)}
	g.populateCache(key, value)
	return value, nil
}

// populateCache 提供填充缓存的能力
func (g *Group) populateCache(key string, value ByteView) {
	g.cache.Add(key, value)
}
