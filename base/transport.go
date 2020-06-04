// Copyright 2009 The Das Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package transport implements net layer.
package base

import (
	"context"
)

//
type Processer func(ctx context.Context, cmd Cmd, client Client)

//
type Transporter interface {
	Broadcaster
	//
	Name() string
	//
	Init(cfg map[string]interface{})
	//
	Start(handle RequestHandler, rhandle ResponseHandler, pctl Protocoler)
	//
	Stop()
}

//
type Broadcaster interface {
	//
	Broadcast(cmd Cmd, ids []uint64)
}
