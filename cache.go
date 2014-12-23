package main

import (
	"errors"
	"image"
	"sync"
)

type RGBACache struct {
	mu     *sync.RWMutex
	cache  []*image.RGBA
	size   int
	curr   int
	keymap map[string]int
}

func NewRGBACache(size int) *RGBACache {
	return &RGBACache{
		cache:  make([]*image.RGBA, size),
		size:   size,
		keymap: make(map[string]int),
		mu:     &sync.RWMutex{},
	}
}

func (c *RGBACache) Get(key string) (im *image.RGBA, exists bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, exists := c.keymap[key]
	return c.cache[v], exists
}

func (c *RGBACache) Put(key string, im *image.RGBA) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.keymap[key]; exists {
		return errors.New("key already exists")
	}
	if len(c.keymap) >= c.size {
		rmidx := (c.curr + c.size - 1) % c.size
		rmkey := ""
		for k, idx := range c.keymap {
			if idx == rmidx {
				rmkey = k
				break
			}
		}
		delete(c.keymap, rmkey)
	}
	c.curr = (c.curr + 1) % c.size
	c.cache[c.curr] = im
	c.keymap[key] = c.curr
	return nil
}
