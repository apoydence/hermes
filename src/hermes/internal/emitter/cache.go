package emitter

type Cache struct {
	factory EmitterFetcher
	cache   map[string]Emitter
}

func NewCache(factory EmitterFetcher) *Cache {
	return &Cache{
		factory: factory,
		cache:   make(map[string]Emitter),
	}
}

func (c *Cache) Fetch(id string) Emitter {
	if emitter, ok := c.cache[id]; ok {
		return emitter
	}

	emitter := c.factory.Fetch(id)
	c.cache[id] = emitter
	return emitter
}
