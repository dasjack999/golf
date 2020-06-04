package main

import (
	"../modules"
	"../repository"
	"crypto/md5"
	"encoding/hex"
	"github.com/go-redis/redis"
	"github.com/wonderivan/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
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
	if rp == nil {
		logger.Debug("repo init fail")
	} else {
		repo.InitRedisSource(rp, "")
	}
	return rp
}
func md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
func main() {
	//client :=initMg()
	client := initRedis()
	if client == nil {
		return
	}

	name := "jack"
	h := md5V(name)
	k, _ := primitive.ObjectIDFromHex(h)

	logger.Debug("id", k)
	keys, _ := client.Keys("*").Result()
	for _, k := range keys {
		logger.Debug("key", k)
	}
	account := modules.NewAccount("redis")
	account.LoadByName("noexist")

	account.Pwd = "jack"
	account.CreateTime = time.Now().Unix()
	account.Name = "jack2"

	logger.Debug("save account ", account.SaveByName())

	account2 := modules.NewAccount("redis")
	ok := account2.LoadByName("jack2")

	if !ok {
		logger.Error("load account fail")
	} else {
		account2.Pwd = "jack2"
		account2.LastLoginTime = time.Now().Unix()

		logger.Debug("update account ", account2.Save(account2))
	}

}
