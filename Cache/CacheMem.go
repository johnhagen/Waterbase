package CacheMem

import (
	"fmt"
	"sync"
	"time"
)

var Cache CacheDB

type CacheDB struct {
	cacheData        map[string]CacheEntity
	maxEntries       int
	maxCacheLifeTime int
	m                sync.Mutex
}

type CacheEntity struct {
	data         []byte
	creationTime time.Time
}

func (d *CacheDB) Init(lifeTime int, maxEntries int) {
	d.cacheData = make(map[string]CacheEntity)
	d.maxCacheLifeTime = lifeTime
	d.maxEntries = maxEntries
}

func (c *CacheDB) Insert(name string, data []byte) *[]byte {
	c.m.Lock()
	if len(c.cacheData) >= c.maxEntries {
		fmt.Println("Max cache entries reached")
		c.m.Unlock()
		return nil
	}

	newCache := CacheEntity{}

	newCache.data = data
	newCache.creationTime = time.Now()

	c.cacheData[name] = newCache
	c.m.Unlock()
	return &newCache.data
}

func (c *CacheDB) Get(name string) *[]byte {
	c.m.Lock()
	if _, exist := c.cacheData[name]; !exist {
		c.m.Unlock()
		return nil
	}

	tNow := time.Now()

	if int(tNow.Sub(c.cacheData[name].creationTime).Seconds()) > c.maxCacheLifeTime {
		c.m.Unlock()
		fmt.Println("Cached data is to old: " + name)
		return nil
	}

	data := c.cacheData[name].data

	c.m.Unlock()
	return &data
}

func (c *CacheDB) Delete(name string) *[]byte {
	c.m.Lock()
	if _, exist := c.cacheData[name]; !exist {
		fmt.Println("Cached data does not exist")
		c.m.Unlock()
		return nil
	}

	data := c.cacheData[name].data
	delete(c.cacheData, name)

	c.m.Unlock()
	return &data
}

func (c *CacheDB) PurgeCacheWorker(checkTimer int) {
	for {
		time.Sleep(time.Duration(checkTimer) * time.Second)
		c.m.Lock()
		tNow := time.Now()
		for f, g := range c.cacheData {
			if int(tNow.Sub(g.creationTime).Seconds()) > c.maxCacheLifeTime {
				fmt.Println("CACHE Purging: " + f)
				delete(c.cacheData, f)
			}
		}
		c.m.Unlock()
	}
}
