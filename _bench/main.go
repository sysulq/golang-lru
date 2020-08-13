package main

import (
	"flag"
	"fmt"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/allegro/bigcache"
	"github.com/coocood/freecache"
	"github.com/hnlq715/golang-lru/shardmap"
)

func gcPause() (int64, time.Duration) {
	runtime.GC()
	var stats debug.GCStats
	debug.ReadGCStats(&stats)
	return stats.NumGC, stats.PauseTotal
}

const (
	entries   = 20000000
	valueSize = 100
)

var typ = flag.String("type", "", "")

func main() {
	flag.Parse()
	debug.SetGCPercent(10)
	fmt.Println("Number of entries: ", entries)

	switch *typ {
	case "map":
		//------------------------------------------
		mapCache := make(map[string][]byte)
		for i := 0; i < entries; i++ {
			key, val := generateKeyValue(i, valueSize)
			mapCache[key] = val
		}
		num, total := gcPause()
		fmt.Println("GC pause for map: ", num, total)
	case "shardmap":
		//------------------------------------------
		shardmap := shardmap.New(entries)
		for i := 0; i < entries; i++ {
			key, val := generateKeyValue(i, valueSize)
			shardmap.Set(key, val)
		}
		num, total := gcPause()
		fmt.Println("GC pause for shardmap: ", num, total)
	case "bigcache":
		config := bigcache.Config{
			Shards:             256,
			LifeWindow:         100 * time.Minute,
			MaxEntriesInWindow: entries,
			MaxEntrySize:       200,
			Verbose:            true,
		}

		bigcache, _ := bigcache.NewBigCache(config)
		for i := 0; i < entries; i++ {
			key, val := generateKeyValue(i, valueSize)
			bigcache.Set(key, val)
		}

		num, total := gcPause()
		fmt.Println("GC pause for bigcache: ", num, total)
	case "freecache":
		freeCache := freecache.NewCache(entries * 200) //allocate entries * 200 bytes
		for i := 0; i < entries; i++ {
			key, val := generateKeyValue(i, valueSize)
			if err := freeCache.Set([]byte(key), val, 0); err != nil {
				fmt.Println("Error in set: ", err.Error())
			}
		}

		num, total := gcPause()
		fmt.Println("GC pause for freecache: ", num, total)
	}

}
