// Copyright 2009 The Das Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package route implements dispatch request.
//
package route

import (
	"../base"
	"../protocol"
	"../service"
	"context"
	"github.com/wonderivan/logger"
)

//
func Process(ctx context.Context, cmd base.Cmd, client base.Client) {

	logger.Debug("process", cmd, client)

	switch cmd.(type) {
	case *base.DisConnected:
		service.SAccount.OnReqLogout(ctx, &protocol.ReqLogout{}, client)
	case *protocol.ReqLogin:
		service.SAccount.OnReqLogin(ctx, cmd.(*protocol.ReqLogin), client)
	case *protocol.Chat:
		service.SChat.OnChat(ctx, cmd.(*protocol.Chat), client)
	case *protocol.ReqResource:
		service.SResource.OnReqResource(ctx, cmd.(*protocol.ReqResource), client)
	}

	return
}
