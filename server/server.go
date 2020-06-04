// Copyright 2009 The Das Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package server implements Server container.
//
package server

import (
	"../base"
	_ "./transport"
	"context"
	"github.com/wonderivan/logger"
	"os"
	"os/signal"
)

//
func init() {
	logger.Debug("server init")
}

//the server object is a container of all services,
//the netlayer can be http,ws,rpc
//
type Server struct {
	//
	ctxMiddle context.Context
	//
	ctxLifecycle context.Context
	//
	ctxLifecycleClose context.CancelFunc
	//
	transports map[string]base.Transporter
	//
	middles map[string]base.Middler
	//
	processor base.Processer
}

//the config object
//type Config struct {
//	Name string `json:"Name,string"`
//	Log  string `json:"Log,omitempty"`
//}
//fact of server,this should be a singleton
//cfg is config map of server
func (server *Server) Init(processor base.Processer) {
	server.transports = map[string]base.Transporter{}
	server.middles = map[string]base.Middler{}
	server.ctxLifecycle, server.ctxLifecycleClose = context.WithCancel(context.Background())
	server.ctxMiddle = context.WithValue(server.ctxLifecycle, "mdmgr", server)
	server.processor = processor
}

//should override this method to instance transport
//
func (server *Server) Start() {
	defer func() {
		logger.Debug("server down")
		//remove all transports
		for _, trans := range server.transports {
			_ = server.DelTransport(trans)
		}
	}()
	//
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	//
	for {
		select {
		case <-server.ctxLifecycle.Done():
			logger.Debug("server closed ")
			break
		case <-c:
			_ = server.Stop()
			break
		}
	}
}

//
func (server *Server) Stop() (err error) {
	logger.Debug("server closed ")
	server.ctxLifecycleClose()
	return
}

//
func (server *Server) AddTransport(trans base.Transporter, pctl base.Protocoler) (err error) {
	server.transports[trans.Name()] = trans
	go trans.Start(server, server, pctl)
	return
}

//
func (server *Server) DelTransport(trans base.Transporter) (err error) {
	trans.Stop()
	delete(server.transports, trans.Name())
	return
}

//
func (server *Server) HandleRequest(ctx context.Context, cmd base.Cmd, client base.Client) (handled bool, err error) {
	//do the middle ware first
	for _, md := range server.middles {

		handled, err := md.HandleRequest(server.ctxMiddle, cmd, client)
		if err != nil {
			logger.Error("middleware error", err)
		}
		if handled {
			return handled, err
		}
	}
	//let route do the final job
	server.processor(server.ctxMiddle, cmd, client)
	return handled, err
}

//
func (server *Server) HandleResponse(cmd base.Cmd, client base.Client) base.Cmd {
	//do the middle ware first
	for _, md := range server.middles {

		cmd = md.HandleResponse(cmd, client)
		if cmd == nil {
			//break
			logger.Error("middleware error")
		}
	}
	return cmd
}

//
func (server *Server) Broadcast(cmd base.Cmd, ids []uint64) {
	for _, trans := range server.transports {
		trans.Broadcast(cmd, ids)
	}
}

//find middle in this server
func (server *Server) GetMiddle(key string) base.Middler {
	m, ok := server.middles[key]
	if !ok {
		return nil
	}
	return m
}

//add midlle in the server
func (server *Server) SetMiddle(key string, m base.Middler) {
	m.SetBroadcaster(server)
	//sort,make sure route is the last one
	server.middles[key] = m
}
