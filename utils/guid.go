// Copyright 2010 The Das Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//make the global guid
package utils

import (
	"github.com/sony/sonyflake"
	"github.com/wonderivan/logger"
	"strconv"
	. "sync/atomic"
)

//use automic ops
var guid uint64 = 0

//
var cache Value

//
var sf *sonyflake.Sonyflake

//
func init() {
	cache.Store(make([]uint64, 0))

	var st sonyflake.Settings
	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		panic("sonyflake not created")
	}
}

//
func GetGlobalId() uint64 {
	c := cache.Load().([]uint64)
	if len(c) > 0 {
		//logger.Debug("get from cache", c)
		g := c[0]
		cache.Store(c[1:])
		return g
	}
	//
	AddUint64(&guid, 1)
	return guid
}

//
func ReStoreGlobalId(gid uint64) {
	c := cache.Load().([]uint64)
	logger.Debug("restore cache", c)
	cache.Store(append(c, gid))
}

//
func GetGuid() string {
	id, _ := sf.NextID()

	return strconv.FormatUint(id, 10)
}
