// Copyright 2009 The Das Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package route implements task service
package service

import (
	"../base"
	"../protocol"
	"context"
	"github.com/wonderivan/logger"
)

//
func init() {
	logger.Debug("s_chat init")

}

//
type sTask struct {
	//
	BaseService
}

//
var STask = &sTask{}

//
func (s *sTask) OnReqTaskList(ctx context.Context, req *protocol.Chat, client base.Client) {
	logger.Debug("handle req ", req.From, req.To, req.Word)

}
