package routing

type EmitterCache struct {
	factory EmitterFetcher
	cache   map[string]Emitter
}

func NewEmitterCache(factory EmitterFetcher) *EmitterCache {
	return &EmitterCache{
		factory: factory,
		cache:   make(map[string]Emitter),
	}
}

func (c *EmitterCache) Fetch(id string) Emitter {
	if emitter, ok := c.cache[id]; ok {
		return emitter
	}

	emitter := c.factory.Fetch(id)
	c.cache[id] = emitter
	return emitter
}
