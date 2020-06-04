// Copyright 2009 The Das Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package the interface of protocl
package base

import (
	"context"
	"github.com/wonderivan/logger"
)

//the remote client agent,provides the inters to remote
//
type Client interface {
	//
	Id() uint64
	//
	Write(Cmd)
	//
	WriteTo(Cmd, ...uint64)
	//
	GetSession(key string) (interface{}, bool)
	//
	SetSession(key string, v interface{})
}

//the base interface for a pack entity
type Cmd interface {
	//
}

//
type HandlerFunc func(context.Context, Cmd, Client)

//
type BaseCmd struct {
	//
	Id string `json:"id"`
}

//
type Connected struct {
	BaseCmd
}

//
type DisConnected struct {
	BaseCmd
}

//the pack/unpack
type Protocoler interface {
	//
	Pack(cmd Cmd) (data []byte, err error)
	//
	UnPack(data []byte) (cmd Cmd, len int, err error)
}

//
func init() {
	logger.Debug("protocol init")
}
