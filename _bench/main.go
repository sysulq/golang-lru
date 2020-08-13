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

func memStatus() runtime.MemStats {
	var memStatus runtime.MemStats
	runtime.ReadMemStats(&memStatus)
	return memStatus
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
		fmt.Printf("mem stats: %+v\n", memStatus())
	case "shardmap":
		//------------------------------------------
		shardmap := shardmap.New(entries)
		for i := 0; i < entries; i++ {
			key, val := generateKeyValue(i, valueSize)
			shardmap.Set(key, val)
		}
		num, total := gcPause()
		fmt.Println("GC pause for shardmap: ", num, total)
		fmt.Printf("mem stats: %+v\n", memStatus())
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
		fmt.Printf("mem stats: %+v\n", memStatus())
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
		fmt.Printf("mem stats: %+v\n", memStatus())
	}

}

func generateKeyValue(index int, valSize int) (string, []byte) {
	key := fmt.Sprintf("key-%010d", index)
	fixedNumber := []byte(fmt.Sprintf("%010d", index))
	val := append(make([]byte, valSize-10), fixedNumber...)

	return key, val
}
