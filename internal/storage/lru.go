package storage

import "container/list"

type LRU struct {
	cap int
	l   *list.List
	m   map[string]*list.Element
}

type Pair struct {
	key   string
	value string
}

func NewLRU(capacity int) *LRU {
	return &LRU{
		cap: capacity,
		l:   new(list.List),
		m:   make(map[string]*list.Element, capacity),
	}
}

func (c *LRU) Get(key string) string {
	if node, ok := c.m[key]; ok {
		val := node.Value.(*list.Element).Value.(Pair).value
		c.l.MoveToFront(node)
		return val
	}
	return ""
}

func (c *LRU) Put(key string, value string) (string, string) {

	removedKey := ""
	removedValue := ""

	if node, ok := c.m[key]; ok {
		c.l.MoveToFront(node)
		node.Value.(*list.Element).Value = Pair{key: key, value: value}
	} else {
		if c.l.Len() == c.cap {
			idx := c.l.Back().Value.(*list.Element).Value.(Pair).key
			removedKey = idx
			removedValue = c.m[idx].Value.(Pair).value
			delete(c.m, idx)
			c.l.Remove(c.l.Back())
		}

		node := &list.Element{
			Value: Pair{
				key:   key,
				value: value,
			},
		}
		ptr := c.l.PushFront(node)
		c.m[key] = ptr
	}

	return removedKey, removedValue
}
