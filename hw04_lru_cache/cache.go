package hw04lrucache

import "fmt"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type cacheItem struct {
	key   Key
	value interface{}
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (c lruCache) Set(key Key, value interface{}) bool {
	item, ok := c.items[key]
	if ok {
		item.Value = cacheItem{key: key, value: value}
		c.queue.MoveToFront(item)
		return true
	}
	if c.queue.Len() == c.capacity {
		var itemValue cacheItem
		itemValue, ok = c.queue.Back().Value.(cacheItem)
		if !ok {
			err := fmt.Errorf("cache is broken on element %v", c.queue.Back())
			fmt.Println(err.Error())
		}
		delete(c.items, itemValue.key)
		c.queue.Remove(c.queue.Back())
	}
	item = c.queue.PushFront(cacheItem{key: key, value: value})
	c.items[key] = item
	return false
}

func (c lruCache) Get(key Key) (interface{}, bool) {
	item, ok := c.items[key]
	if !ok {
		return nil, false
	}
	c.queue.MoveToFront(item)

	var itemValue cacheItem
	itemValue, ok = item.Value.(cacheItem)
	if !ok {
		return nil, false
	}
	return itemValue.value, true
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
