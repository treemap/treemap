package main

import (
	"log"
)

type Cache struct {
	cache map[string]interface{}
}

func (c *Cache) Get(cacheKey string, runner func() interface{}) interface{} {
	if cacheKey == "" {
		return nil
	}

	if b, ok := c.cache[cacheKey]; ok {
		log.Println("Hitting cache:", cacheKey)
		return b
	}

	c.cache[cacheKey] = runner()

	return c.cache[cacheKey]
}

func NewCache() (c Cache) {
	c.cache = make(map[string]interface{})
	return
}
