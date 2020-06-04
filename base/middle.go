// Copyright 2009 The Das Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package the interface of middleware
// the aim of this package is to process the common work serially
// on each request or response and hold the global modules
package base

import (
	"context"
)

//
type RequestHandler interface {
	//
	HandleRequest(ctx context.Context, cmd Cmd, client Client) (handled bool, err error)
}

//
type ResponseHandler interface {
	//
	HandleResponse(cmd Cmd, client Client) Cmd
}

//
type MiddleMger interface {
	//
	GetMiddle(key string) Middler
	//
	SetMiddle(key string, m Middler)
}

//
type Middler interface {
	//do the service
	RequestHandler
	//
	ResponseHandler
	//
	GetName() string
	//
	Enable(enable ...bool) bool
	//
	SetBroadcaster(broadcaster Broadcaster)
}

//
var middlers = []Middler{}

//register midlle
func RegisterMiddle(middle Middler) {
	middlers = append(middlers, middle)
}

//find the global middle
func FindMiddle(name string) Middler {

	for _, middler := range middlers {
		//logger.Debug("check md",middler.GetName(),reflect.TypeOf(middler).Name())
		if middler.GetName() == name {
			return middler
		}
	}
	return nil
}

//setup the global midlle instance
func InitMiddle(name string, enable bool) Middler {
	m := FindMiddle(name)
	if m != nil {
		m.Enable(enable)
	}
	return m
}

//
type BaseMiddle struct {
	//
	Name string
	//
	enable bool
	//
	br Broadcaster
}

//
func (m *BaseMiddle) HandleResponse(cmd Cmd, client Client) Cmd {
	return cmd
}

//
func (m *BaseMiddle) GetName() string {
	return m.Name
}

//
func (m *BaseMiddle) Enable(enable ...bool) bool {
	if len(enable) > 0 {
		m.enable = enable[0]
	}
	return m.enable
}

//
func (m *BaseMiddle) SetBroadcaster(broadcaster Broadcaster) {
	m.br = broadcaster
}

//
func (m *BaseMiddle) Broadcast(cmd Cmd, ids []uint64) {
	m.br.Broadcast(cmd, ids)
}
