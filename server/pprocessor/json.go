// Copyright 2009 The Das Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package protocol implements by json
package pprocessor

import (
	"../../base"
	"github.com/json-iterator/go"
	"github.com/wonderivan/logger"
	"reflect"
)

//the global json object
var json = jsoniter.ConfigCompatibleWithStandardLibrary

//
type ptclJson struct {
	//
}

//
func (js *ptclJson) UnPack(data []byte) (cmd base.Cmd, len int, err error) {

	defer func() {
		if p := recover(); p != nil {
			logger.Error("internal error: %v", p)
		}
	}()

	n := jsoniter.Get(data, "id").ToString()
	t, ok := base.GetClass(n)
	if !ok {
		logger.Debug("UnPack:without regclass,so as interface{}")
		err = json.Unmarshal(data, &cmd)
		return
	}
	cmd = reflect.New(t.Elem()).Interface()
	err = json.Unmarshal(data, cmd)

	return
}

//
func (js *ptclJson) Pack(cmd base.Cmd) (data []byte, err error) {
	defer func() {
		if p := recover(); p != nil {
			logger.Error("internal error: %v", p)
		}
	}()
	data, err = json.Marshal(cmd)
	return
}

//export this as singleton
var GJsonPctl = &ptclJson{}

//
func init() {
	logger.Debug("json pctl init")
}
