package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (c lruCache) Set(key Key, value interface{}) bool {
	item, ok := c.items[key]
	if ok {
		item.Value = value
		c.queue.MoveToFront(item)
		return true
	}
	if c.queue.Len() == c.capacity {
		item = c.queue.Back()
		for k, v := range c.items {
			if v == item {
				delete(c.items, k)
				break
			}
		}
		c.queue.Remove(item)
	}
	item = c.queue.PushFront(value)
	c.items[key] = item
	return false
}

func (c lruCache) Get(key Key) (interface{}, bool) {
	item, ok := c.items[key]
	if !ok {
		return nil, false
	}
	c.queue.MoveToFront(item)

	return item.Value, true
}

func (c lruCache) Clear() {
	for key, item := range c.items {
		c.queue.Remove(item)
		delete(c.items, key)
	}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
