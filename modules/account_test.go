// Copyright 2009 The Das Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package modules_test

import (
	"../repository"
	"github.com/go-redis/redis"
	"github.com/wonderivan/logger"
	"testing"
)

func initMg() *repo.MongoClient {
	rp, err := repo.NewMongoClient("mongodb://localhost:27017", 15)
	if err != nil {
		logger.Debug("repo init fail", err)
	} else {
		repo.InitMongoSource(rp, "dgc")
	}
	return rp
}

func initRedis() *redis.Client {
	//url :=cfg.Config["url"].(string)
	opt := &redis.Options{
		Addr:     "192.168.28.239:6379",
		Password: "",
		DB:       0,
	}
	rp := redis.NewClient(opt)

	return rp
}

func TestAccount_Load(t *testing.T) {
	//init mgclient
	client := initMg()
	if client == nil {
		return
	}

	//account :=&Account{}
	//account.Query("noexist")
	//
	//account.CreateTime=time.Now().Unix()
	//account.Name="jack"
	//account.Save()

}
