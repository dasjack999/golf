// Copyright 2009 The Das Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package dgc services container
package main

//
import (
	"./base"
	"./repository"
	"./route"
	"./server"
	"./server/pprocessor"
	"flag"
	"github.com/go-redis/redis"
	"github.com/json-iterator/go"
	"github.com/wonderivan/logger"

	"io/ioutil"
	"runtime"
	"time"
)

//
type Config struct {
	Name       string         `json:"name"`
	Cores      int            `json:"cores"`
	Middles    []CfgMiddle    `json:"middles"`
	Repos      []CfgRepo      `json:"repos"`
	Transports []CfgTransport `json:"transports"`
}

//
type CfgRepo struct {
	Name   string                 `json:"name"`
	Config map[string]interface{} `json:"config,omitempty"`
}

//
type CfgMiddle struct {
	Name   string                 `json:"name"`
	Config map[string]interface{} `json:"config,omitempty"`
}

//
type CfgTransport struct {
	Name   string                 `json:"name"`
	Config map[string]interface{} `json:"config"`
}

var (
	//
	cfgPath    string = ""
	logCfgPath string = ""
)

//
//the global json object
var json = jsoniter.ConfigCompatibleWithStandardLibrary

//
func init() {
	flag.StringVar(&cfgPath, "config", "config.json", "server config")
	flag.StringVar(&logCfgPath, "log", "log.json", "server config")
}

//
func main() {

	defer func() {
		logger.Debug("server shutdown")
	}()
	//config logger
	if len(logCfgPath) > 0 {
		_ = logger.SetLogger(logCfgPath)
	}
	//parse cmdline
	flag.Parse()
	//
	cfg, err := loadConfig(cfgPath)
	if err != nil {
		logger.Debug("parse config error:", err)
		return
	}
	logger.Debug("load config", cfg)
	//set cores
	runtime.GOMAXPROCS(cfg.Cores)
	//
	app := NewApp(cfg)
	//
	logger.Debug("server start")
	app.Start()
}

//
func loadConfig(path string) (cfg Config, err error) {

	buf, err := ioutil.ReadFile(path)
	if err == nil {
		err = json.Unmarshal(buf, &cfg)
	}

	return
}

//
type Application struct {
	server.Server
}

//
func NewApp(cfg Config) *Application {
	app := &Application{}
	app.Init(route.Process)
	app.initMidlleWares(cfg.Middles)
	app.initTransports(cfg.Transports)
	app.initRepos(cfg.Repos)
	return app
}

//
func (app *Application) initRepos(cfgs []CfgRepo) {

	for _, cfg := range cfgs {
		if cfg.Name == "mongo" {
			url := cfg.Config["url"].(string)
			timeout := cfg.Config["timeout"].(float64)

			rp, err := repo.NewMongoClient(url, time.Duration(timeout))
			if err != nil {
				logger.Debug("repo init fail", err)
			} else {
				dname := cfg.Config["db"].(string)

				repo.InitMongoSource(rp, dname)
			}
		} else if cfg.Name == "redis" {
			url := cfg.Config["url"].(string)
			opt := &redis.Options{
				Addr:     url,
				Password: "",
				DB:       0,
			}
			rp := redis.NewClient(opt)
			if rp == nil {
				logger.Debug("repo init fail")
			} else {
				repo.InitRedisSource(rp, "dgc")
			}
		}
	}
}

//
func (app *Application) initMidlleWares(cfgs []CfgMiddle) {

	for _, cfg := range cfgs {
		md := base.InitMiddle(cfg.Name, true)
		if md == nil {
			logger.Error("create middle fail", cfg)
		} else {
			logger.Debug("add middle", cfg)
			app.SetMiddle(cfg.Name, md.(base.Middler))
		}
	}
}

//
func (app *Application) initTransports(cfgs []CfgTransport) {
	//trans := transport.New("ws", websocket.TextMessage, ":8080", "/chat")
	//app.AddTransport(trans, cfg, protocol.GJsonPctl)
	for _, cfg := range cfgs {
		trans := base.NewInstance(cfg.Name)
		if trans == nil {
			logger.Error("create transport fail", cfg)
		} else {
			logger.Debug("add transport", cfg)
			pctl := pprocessor.GJsonPctl
			t := trans.(base.Transporter)
			t.Init(cfg.Config)
			app.AddTransport(t, pctl)
		}
	}
}
