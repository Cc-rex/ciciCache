package lru

import "container/list"

// TODO: implement the LRU(least recently used)

type OnEliminated func(key string, value Sizeable)

type Cache struct {
	maxMemory        int64                    // Maximum memory allowed
	usedMemory       int64                    // Currently used memory
	hashMap          map[string]*list.Element // value is a pointer to the corresponding node in list
	doubleLinkedList *list.List               // The head of the list is the least recently used data
	callBack         OnEliminated             // Handler function executed when key-value is eliminated
}

// Sizeable Get the size of the space occupied by the object itself
type Sizeable interface {
	Size() int // byte-based
}

// Value Define the objects stored in a double Linked List node
type Value struct {
	key   string
	value Sizeable
}

// New the Constructor of Cache
func New(maxMemory int64, callBack OnEliminated) *Cache {
	return &Cache{
		maxMemory:        maxMemory,
		doubleLinkedList: list.New(),
		hashMap:          make(map[string]*list.Element),
		callBack:         callBack,
	}
}

// Get Retrieve the value of the corresponding key
func (c *Cache) Get(key string) (value Sizeable, ok bool) {
	if element, exist := c.hashMap[key]; exist {
		c.doubleLinkedList.MoveToFront(element)
		entry := element.Value.(*Value)
		return entry.value, true
	}
	return
}

// Remove delete the least recently used node
func (c *Cache) Remove() {
	element := c.doubleLinkedList.Back()
	if element != nil {
		entry := element.Value.(*Value)
		k, v := entry.key, entry.value
		delete(c.hashMap, k)                            // remove mapping
		c.doubleLinkedList.Remove(element)              // remove cache
		c.usedMemory -= int64(len(k)) + int64(v.Size()) // update the memory situation
		// executing Handler function
		if c.callBack != nil {
			c.callBack(k, v)
		}
	}
}

// Add adds a value to the cache.
func (c *Cache) Add(k string, v Sizeable) {
	// check the memory
	kvSize := int64(len(k)) + int64(v.Size())
	for c.maxMemory != 0 && c.usedMemory+kvSize > c.maxMemory {
		c.Remove()
	}
	// update/create cache
	if element, exist := c.hashMap[k]; exist {
		// key exists, so update the value
		c.doubleLinkedList.MoveToFront(element)
		entry := element.Value.(*Value)
		// change the size of memory first, then change the value
		c.usedMemory += kvSize
		entry.value = v
	} else {
		// not exist, add the key and value
		entry := Value{
			key:   k,
			value: v,
		}
		element := c.doubleLinkedList.PushFront(&entry)
		c.hashMap[k] = element
		c.usedMemory += kvSize
	}
}
