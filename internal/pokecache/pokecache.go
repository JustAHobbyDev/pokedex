package pokecache

import (
    "fmt"
    "time"
    "sync"
    "os"
    "os/signal"
    "golang.org/x/sys/unix"
)

type Cache struct {
    Entries map[string]cacheEntry
    mutex sync.Mutex
    interval time.Duration
}

type cacheEntry struct {
    createdAt time.Time
    val       []byte
}

func (c *Cache) Add(key string, val []byte) {
    c.mutex.Lock()

    newCacheEntry := &cacheEntry{
        createdAt: time.Now(),
        val: val,
    }
    c.Entries[key] = *newCacheEntry

    c.mutex.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
    c.mutex.Lock()
    entry, ok := c.Entries[key]
    c.mutex.Unlock()

    if !ok {
        return []byte{}, false
    }

    return entry.val, true
}

func (c *Cache) ReapLoop() {
	ticker := time.NewTicker(c.interval * time.Second)
	defer ticker.Stop()

    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, unix.SIGINT, unix.SIGTERM)
	done := make(chan bool)

	go func () {
        for {
            select {
            case <-ticker.C:
                now := time.Now()
                for key, entry := range c.Entries {
                    if entry.createdAt.Before(now) {
                        c.mutex.Lock()
                        delete(c.Entries, key)
                        c.mutex.Unlock()
                    }
                }
            case sig := <-sigs:
                fmt.Println("Shutting down gracefully...")
                fmt.Println("Received signal:", sig)
                close(done)
                return
            case <-done:
                return
            }
        }
    }()
}

func NewCache(interval time.Duration) Cache {
    return Cache {
        Entries: make(map[string]cacheEntry),
        interval: interval,
    }
}
