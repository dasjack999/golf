// Copyright 2010 The Das Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package repo

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/wonderivan/logger"
)

//
func init() {
	RegisterDataSource("redis", &RedisSource{})
}

//
type RedisSource struct {
}

//
func (e *RedisSource) Init(query map[string]interface{}) error {

	return nil
}

//
func (e *RedisSource) getName(collName string, id string) string {

	return dbName + ":" + collName + ":" + id
}

//
func (e *RedisSource) Query(collName string, query map[string]interface{}) (result interface{}, err error) {

	start := query["from"].(uint64)
	total := query["total"].(int64)
	cursor := start
	keys, cursor, err := dbclient.SScan(collName, cursor, "", total).Result()

	result = keys
	return
}

//
func (e *RedisSource) Load(collName string, id string, v interface{}) bool {

	key := e.getName(collName, id)

	str, err := dbclient.Get(key).Result()
	if err != nil {
		logger.Error("Query", err)
		return false
	}
	err = json.Unmarshal([]byte(str), v)
	if err != nil {
		logger.Error("Query", err)
	}
	return err == nil
}

//
func (e *RedisSource) Update(collName string, id string, v interface{}) bool {

	key := e.getName(collName, id)
	str, err := json.Marshal(v)
	if err != nil {
		logger.Error("Update", err)
		return false
	}
	s, err := dbclient.SetXX(key, str, 0).Result()
	return s
}

//
func (e *RedisSource) Insert(collName string, id string, v interface{}) bool {
	key := e.getName(collName, id)

	str, err := json.Marshal(v)
	if err != nil {
		logger.Error("Insert", err)
		return false
	}
	s, err := dbclient.SetNX(key, str, 0).Result()
	if err != nil {
		logger.Error("Insert", err)

	} else {
		dbclient.SAdd(collName, key)
	}

	return s
}

//
func (e *RedisSource) Remove(collName string, id string) bool {
	key := e.getName(collName, id)
	s, err := dbclient.Del(key).Result()
	if err != nil {
		logger.Error("Remove", err)
	} else {
		dbclient.SRem(collName, key)
	}
	return s == 1
}

//
var dbclient *redis.Client

//
var dbName string

//
func InitRedisSource(client *redis.Client, name string) {
	dbclient = client
	dbName = name
	logger.Debug("init redis client", name, client)
}
