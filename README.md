golang-lru
==========
[![Build Status](https://travis-ci.org/hnlq715/golang-lru.svg?branch=master)](https://travis-ci.org/hnlq715/golang-lru)
[![Coverage](https://codecov.io/gh/hnlq715/golang-lru/branch/master/graph/badge.svg)](https://codecov.io/gh/hnlq715/golang-lru)

This provides the `lru` package which implements a fixed-size
thread safe LRU cache with expire feature. It is based on [golang-lru](https://github.com/hashicorp/golang-lru).

Documentation
=============

Full docs are available on [Godoc](http://godoc.org/github.com/remones/golang-lru)

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
