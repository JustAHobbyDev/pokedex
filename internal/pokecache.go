package pokecache

import (
    "time"
)

type Cache struct {
    Entries map[string]cacheEntry
}

type cacheEntry struct {
    createdAt time.Time
    val       []byte
}
