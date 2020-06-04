// Copyright 2009 The Das Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package route implements chat service
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
type sChat struct {
	//
	BaseService
}

//
var SChat = &sChat{}

//
func (s *sChat) OnChat(ctx context.Context, req *protocol.Chat, client base.Client) {
	logger.Debug("handle req ", req.From, req.To, req.Word)

	switch req.Type {
	//case 0://
	case 1: //single or all
		client.WriteTo(req, req.To)
	case 2: //group
		rm := s.GetRoomMiddler(ctx)
		if rm != nil {
			ments := rm.GetRoomMeets((int)(req.To))
			client.WriteTo(req, ments...)
		}

	}

}
