package emitter

import "hermes/internal/datastructures"

type Registry interface {
	GetList(ID string) *datastructures.LinkedList
}

type Cache struct {
	registry Registry
	cache    map[string]Emitter
}

func NewCache(registry Registry) *Cache {
	return &Cache{
		registry: registry,
		cache:    make(map[string]Emitter),
	}
}

func (c *Cache) Fetch(ID string) Emitter {
	if list, ok := c.cache[ID]; ok {
		return list
	}

	list := c.registry.GetList(ID)
	reader := NewSubscriptionListReader(list)
	c.cache[ID] = reader
	return reader
}
