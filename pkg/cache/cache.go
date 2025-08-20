package cache

import (
	"hash/fnv"
	"sync"
	"time"
)

type CacheItem struct {
	Value     interface{}
	ExpiresAt time.Time
}

func (item *CacheItem) IsExpired() bool {
	return time.Now().After(item.ExpiresAt)
}

type Cache struct {
	items   sync.Map
	janitor *janitor
}

type janitor struct {
	interval time.Duration
	stop     chan bool
}

func NewCache(cleanupInterval time.Duration) *Cache {
	c := &Cache{}

	if cleanupInterval > 0 {
		c.janitor = &janitor{
			interval: cleanupInterval,
			stop:     make(chan bool),
		}
		go c.janitor.run(c)
	}

	return c
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	item := &CacheItem{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
	c.items.Store(key, item)
}

func (c *Cache) Get(key string) (interface{}, bool) {
	value, exists := c.items.Load(key)
	if !exists {
		return nil, false
	}

	item := value.(*CacheItem)
	if item.IsExpired() {
		c.items.Delete(key)
		return nil, false
	}

	return item.Value, true
}

func (c *Cache) GetWithTTL(key string) (interface{}, time.Duration, bool) {
	value, exists := c.items.Load(key)
	if !exists {
		return nil, 0, false
	}

	item := value.(*CacheItem)
	if item.IsExpired() {
		c.items.Delete(key)
		return nil, 0, false
	}

	remaining := time.Until(item.ExpiresAt)
	return item.Value, remaining, true
}

func (c *Cache) Delete(key string) {
	c.items.Delete(key)
}

func (c *Cache) Exists(key string) bool {
	_, exists := c.Get(key)
	return exists
}

func (c *Cache) Clear() {
	c.items.Range(func(key, value interface{}) bool {
		c.items.Delete(key)
		return true
	})
}

func (c *Cache) Size() int {
	count := 0
	c.items.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

func (c *Cache) CleanExpired() int {
	var expiredKeys []interface{}

	c.items.Range(func(key, value interface{}) bool {
		item := value.(*CacheItem)
		if item.IsExpired() {
			expiredKeys = append(expiredKeys, key)
		}
		return true
	})

	for _, key := range expiredKeys {
		c.items.Delete(key)
	}

	return len(expiredKeys)
}

func (c *Cache) Stop() {
	if c.janitor != nil {
		c.janitor.stop <- true
	}
}

func (j *janitor) run(c *Cache) {
	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.CleanExpired()
		case <-j.stop:
			return
		}
	}
}

var DefaultCache = NewCache(5 * time.Second)

func Set(key string, value interface{}, ttl time.Duration) {
	DefaultCache.Set(key, value, ttl)
}

func Get(key string) (interface{}, bool) {
	return DefaultCache.Get(key)
}

func Delete(key string) {
	DefaultCache.Delete(key)
}

func Exists(key string) bool {
	return DefaultCache.Exists(key)
}

func Clear() {
	DefaultCache.Clear()
}

type ShardedCache struct {
	shards    []*Cache
	shardMask uint32
	janitor   *shardedJanitor
}

type shardedJanitor struct {
	interval time.Duration
	stop     chan bool
}

func isPowerOfTwo(n uint32) bool {
	return n > 0 && (n&(n-1)) == 0
}

func nextPowerOfTwo(n uint32) uint32 {
	if n <= 1 {
		return 1
	}
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n++
	return n
}

func (sc *ShardedCache) hash(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

func (sc *ShardedCache) getShard(key string) *Cache {
	return sc.shards[sc.hash(key)&sc.shardMask]
}

func NewShardedCache(shardCount uint32, cleanupInterval time.Duration) *ShardedCache {
	if shardCount == 0 {
		shardCount = 16
	}
	if !isPowerOfTwo(shardCount) {
		shardCount = nextPowerOfTwo(shardCount)
	}

	sc := &ShardedCache{
		shards:    make([]*Cache, shardCount),
		shardMask: shardCount - 1,
	}

	for i := uint32(0); i < shardCount; i++ {
		sc.shards[i] = NewCache(0)
	}

	if cleanupInterval > 0 {
		sc.janitor = &shardedJanitor{
			interval: cleanupInterval,
			stop:     make(chan bool),
		}
		go sc.janitor.run(sc)
	}

	return sc
}

func (sc *ShardedCache) Set(key string, value interface{}, ttl time.Duration) {
	shard := sc.getShard(key)
	shard.Set(key, value, ttl)
}

func (sc *ShardedCache) Get(key string) (interface{}, bool) {
	shard := sc.getShard(key)
	return shard.Get(key)
}

func (sc *ShardedCache) Delete(key string) {
	shard := sc.getShard(key)
	shard.Delete(key)
}

func (sc *ShardedCache) Exists(key string) bool {
	shard := sc.getShard(key)
	return shard.Exists(key)
}

func (sc *ShardedCache) CleanExpired() int {
	totalCleaned := 0
	for _, shard := range sc.shards {
		totalCleaned += shard.CleanExpired()
	}
	return totalCleaned
}

func (sc *ShardedCache) Stop() {
	if sc.janitor != nil {
		sc.janitor.stop <- true
	}
	for _, shard := range sc.shards {
		shard.Stop()
	}
}

func (sj *shardedJanitor) run(sc *ShardedCache) {
	ticker := time.NewTicker(sj.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			sc.CleanExpired()
		case <-sj.stop:
			return
		}
	}
}
