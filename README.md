golang-lru
==========
[![Build Status](https://travis-ci.org/hnlq715/golang-lru.svg?branch=master)](https://travis-ci.org/hnlq715/golang-lru)
[![Coverage](https://codecov.io/gh/hnlq715/golang-lru/branch/master/graph/badge.svg)](https://codecov.io/gh/hnlq715/golang-lru)

This provides the `lru` package which implements a fixed-size
thread safe LRU cache with expire feature. It is based on [golang-lru](https://github.com/hashicorp/golang-lru).

Documentation
=============

Full docs are available on [Godoc](http://godoc.org/github.com/hnlq715/golang-lru)

Example
=======

Using the LRU is very simple:

```go
l, _ := NewARCWithExpire(128, 30*time.Second)
for i := 0; i < 256; i++ {
    l.Add(i, nil)
}
if l.Len() != 128 {
    panic(fmt.Sprintf("bad len: %v", l.Len()))
}
```

Benchmarks
===

[without pool](https://github.com/hnlq715/golang-lru/tree/master)
---
```
Running tool: /home/liqi/workspace/go/bin/go test -benchmem -run=^$ -coverprofile=/tmp/vscode-govFyHq9/go-code-cover -bench . github.com/hnlq715/golang-lru

goos: linux
goarch: amd64
pkg: github.com/hnlq715/golang-lru
Benchmark2Q_Rand-4    	 1000000	      1415 ns/op	     158 B/op	       5 allocs/op
--- BENCH: Benchmark2Q_Rand-4
    2q_test.go:34: hit: 0 miss: 1 ratio: 0.000000
    2q_test.go:34: hit: 1 miss: 99 ratio: 0.010101
    2q_test.go:34: hit: 1381 miss: 8619 ratio: 0.160227
    2q_test.go:34: hit: 248530 miss: 751470 ratio: 0.330725
Benchmark2Q_Freq-4    	 1000000	      1017 ns/op	     143 B/op	       5 allocs/op
--- BENCH: Benchmark2Q_Freq-4
    2q_test.go:66: hit: 1 miss: 0 ratio: +Inf
    2q_test.go:66: hit: 100 miss: 0 ratio: +Inf
    2q_test.go:66: hit: 9840 miss: 160 ratio: 61.500000
    2q_test.go:66: hit: 332417 miss: 667583 ratio: 0.497941
BenchmarkARC_Rand-4   	 1000000	      1402 ns/op	     193 B/op	       6 allocs/op
--- BENCH: BenchmarkARC_Rand-4
    arc_test.go:39: hit: 0 miss: 1 ratio: 0.000000
    arc_test.go:39: hit: 1 miss: 99 ratio: 0.010101
    arc_test.go:39: hit: 1398 miss: 8602 ratio: 0.162520
    arc_test.go:39: hit: 249099 miss: 750901 ratio: 0.331733
BenchmarkARC_Freq-4   	  963909	      1190 ns/op	     166 B/op	       5 allocs/op
--- BENCH: BenchmarkARC_Freq-4
    arc_test.go:71: hit: 1 miss: 0 ratio: +Inf
    arc_test.go:71: hit: 100 miss: 0 ratio: +Inf
    arc_test.go:71: hit: 9860 miss: 140 ratio: 70.428571
    arc_test.go:71: hit: 310475 miss: 653434 ratio: 0.475144
BenchmarkLRU_Rand-4   	 2287102	       613 ns/op	      88 B/op	       3 allocs/op
--- BENCH: BenchmarkLRU_Rand-4
    lru_test.go:34: hit: 0 miss: 1 ratio: 0.000000
    lru_test.go:34: hit: 0 miss: 100 ratio: 0.000000
    lru_test.go:34: hit: 1379 miss: 8621 ratio: 0.159958
    lru_test.go:34: hit: 248489 miss: 751511 ratio: 0.330653
    lru_test.go:34: hit: 570640 miss: 1716462 ratio: 0.332451
BenchmarkLRU_Freq-4   	 2456690	       487 ns/op	      83 B/op	       3 allocs/op
--- BENCH: BenchmarkLRU_Freq-4
    lru_test.go:66: hit: 1 miss: 0 ratio: +Inf
    lru_test.go:66: hit: 100 miss: 0 ratio: +Inf
    lru_test.go:66: hit: 9846 miss: 154 ratio: 63.935065
    lru_test.go:66: hit: 312529 miss: 687471 ratio: 0.454607
    lru_test.go:66: hit: 752485 miss: 1704205 ratio: 0.441546
PASS
coverage: 54.9% of statements
ok  	github.com/hnlq715/golang-lru	9.138s

```


[with sync pool](https://github.com/hnlq715/golang-lru/tree/feature/syncpool)
---
```
Running tool: /home/liqi/workspace/go/bin/go test -benchmem -run=^$ -coverprofile=/tmp/vscode-govFyHq9/go-code-cover -bench . github.com/hnlq715/golang-lru

goos: linux
goarch: amd64
pkg: github.com/hnlq715/golang-lru
Benchmark2Q_Rand-4    	 1000000	      1090 ns/op	      92 B/op	       4 allocs/op
--- BENCH: Benchmark2Q_Rand-4
    2q_test.go:34: hit: 0 miss: 1 ratio: 0.000000
    2q_test.go:34: hit: 0 miss: 100 ratio: 0.000000
    2q_test.go:34: hit: 1375 miss: 8625 ratio: 0.159420
    2q_test.go:34: hit: 249496 miss: 750504 ratio: 0.332438
Benchmark2Q_Freq-4    	 1223035	       944 ns/op	      85 B/op	       4 allocs/op
--- BENCH: Benchmark2Q_Freq-4
    2q_test.go:66: hit: 1 miss: 0 ratio: +Inf
    2q_test.go:66: hit: 100 miss: 0 ratio: +Inf
    2q_test.go:66: hit: 9872 miss: 128 ratio: 77.125000
    2q_test.go:66: hit: 334464 miss: 665536 ratio: 0.502548
    2q_test.go:66: hit: 405282 miss: 817753 ratio: 0.495604
BenchmarkARC_Rand-4   	 1000000	      1330 ns/op	     111 B/op	       4 allocs/op
--- BENCH: BenchmarkARC_Rand-4
    arc_test.go:39: hit: 0 miss: 1 ratio: 0.000000
    arc_test.go:39: hit: 0 miss: 100 ratio: 0.000000
    arc_test.go:39: hit: 1368 miss: 8632 ratio: 0.158480
    arc_test.go:39: hit: 248419 miss: 751581 ratio: 0.330529
BenchmarkARC_Freq-4   	 1000000	      1090 ns/op	      93 B/op	       4 allocs/op
--- BENCH: BenchmarkARC_Freq-4
    arc_test.go:71: hit: 1 miss: 0 ratio: +Inf
    arc_test.go:71: hit: 100 miss: 0 ratio: +Inf
    arc_test.go:71: hit: 9876 miss: 124 ratio: 79.645161
    arc_test.go:71: hit: 337535 miss: 662465 ratio: 0.509514
BenchmarkLRU_Rand-4   	 2327682	       509 ns/op	      52 B/op	       2 allocs/op
--- BENCH: BenchmarkLRU_Rand-4
    lru_test.go:34: hit: 0 miss: 1 ratio: 0.000000
    lru_test.go:34: hit: 1 miss: 99 ratio: 0.010101
    lru_test.go:34: hit: 1478 miss: 8522 ratio: 0.173433
    lru_test.go:34: hit: 249019 miss: 750981 ratio: 0.331592
    lru_test.go:34: hit: 580746 miss: 1746936 ratio: 0.332437
BenchmarkLRU_Freq-4   	 2630702	       475 ns/op	      49 B/op	       2 allocs/op
--- BENCH: BenchmarkLRU_Freq-4
    lru_test.go:66: hit: 1 miss: 0 ratio: +Inf
    lru_test.go:66: hit: 100 miss: 0 ratio: +Inf
    lru_test.go:66: hit: 9784 miss: 216 ratio: 45.296296
    lru_test.go:66: hit: 311783 miss: 688217 ratio: 0.453030
    lru_test.go:66: hit: 810266 miss: 1820436 ratio: 0.445094
PASS
coverage: 55.3% of statements
ok  	github.com/hnlq715/golang-lru	9.714s
```



[with list pool](https://github.com/hnlq715/golang-lru/tree/feature/listpool)
---
```
Running tool: /home/liqi/workspace/go/bin/go test -benchmem -run=^$ -coverprofile=/tmp/vscode-govFyHq9/go-code-cover -bench . github.com/hnlq715/golang-lru

goos: linux
goarch: amd64
pkg: github.com/hnlq715/golang-lru
Benchmark2Q_Rand-4    	 1000000	      1311 ns/op	      26 B/op	       2 allocs/op
--- BENCH: Benchmark2Q_Rand-4
    2q_test.go:34: hit: 0 miss: 1 ratio: 0.000000
    2q_test.go:34: hit: 0 miss: 100 ratio: 0.000000
    2q_test.go:34: hit: 1320 miss: 8680 ratio: 0.152074
    2q_test.go:34: hit: 249497 miss: 750503 ratio: 0.332440
Benchmark2Q_Freq-4    	 1515253	       811 ns/op	      25 B/op	       2 allocs/op
--- BENCH: Benchmark2Q_Freq-4
    2q_test.go:66: hit: 1 miss: 0 ratio: +Inf
    2q_test.go:66: hit: 100 miss: 0 ratio: +Inf
    2q_test.go:66: hit: 9880 miss: 120 ratio: 82.333333
    2q_test.go:66: hit: 333341 miss: 666659 ratio: 0.500017
    2q_test.go:66: hit: 416877 miss: 839436 ratio: 0.496616
    2q_test.go:66: hit: 500425 miss: 1014828 ratio: 0.493113
BenchmarkARC_Rand-4   	 1000000	      1162 ns/op	      27 B/op	       2 allocs/op
--- BENCH: BenchmarkARC_Rand-4
    arc_test.go:39: hit: 0 miss: 1 ratio: 0.000000
    arc_test.go:39: hit: 0 miss: 100 ratio: 0.000000
    arc_test.go:39: hit: 1437 miss: 8563 ratio: 0.167815
    arc_test.go:39: hit: 248890 miss: 751110 ratio: 0.331363
BenchmarkARC_Freq-4   	 1418329	       860 ns/op	      26 B/op	       2 allocs/op
--- BENCH: BenchmarkARC_Freq-4
    arc_test.go:71: hit: 1 miss: 0 ratio: +Inf
    arc_test.go:71: hit: 100 miss: 0 ratio: +Inf
    arc_test.go:71: hit: 9847 miss: 153 ratio: 64.359477
    arc_test.go:71: hit: 337917 miss: 662083 ratio: 0.510385
    arc_test.go:71: hit: 472261 miss: 946068 ratio: 0.499183
BenchmarkLRU_Rand-4   	 2970676	       397 ns/op	      16 B/op	       1 allocs/op
--- BENCH: BenchmarkLRU_Rand-4
    lru_test.go:34: hit: 0 miss: 1 ratio: 0.000000
    lru_test.go:34: hit: 0 miss: 100 ratio: 0.000000
    lru_test.go:34: hit: 1375 miss: 8625 ratio: 0.159420
    lru_test.go:34: hit: 249937 miss: 750063 ratio: 0.333221
    lru_test.go:34: hit: 741873 miss: 2228803 ratio: 0.332857
BenchmarkLRU_Freq-4   	 3186918	       359 ns/op	      16 B/op	       1 allocs/op
--- BENCH: BenchmarkLRU_Freq-4
    lru_test.go:66: hit: 1 miss: 0 ratio: +Inf
    lru_test.go:66: hit: 100 miss: 0 ratio: +Inf
    lru_test.go:66: hit: 9874 miss: 126 ratio: 78.365079
    lru_test.go:66: hit: 312430 miss: 687570 ratio: 0.454397
    lru_test.go:66: hit: 977757 miss: 2209161 ratio: 0.442592
PASS
coverage: 55.3% of statements
ok  	github.com/hnlq715/golang-lru	11.639s
```