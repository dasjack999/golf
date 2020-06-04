// Copyright 2009 The Das Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package route implements account service

package service

import (
	"../base"
	"../errors"
	"../modules"
	"../protocol"
	"context"
	"github.com/wonderivan/logger"
)

//
func init() {
	//register singleton
}

//
type sAccount struct {
	//
	BaseService
	//
	client2Account map[uint64]*modules.Account
}

//
var SAccount = &sAccount{
	client2Account: map[uint64]*modules.Account{},
}

//
func (s *sAccount) OnReqLogin(ctx context.Context, req *protocol.ReqLogin, client base.Client) {
	logger.Debug("handle login  ", req.UserId, client.Id())

	account := modules.NewAccount("redis")
	ok := account.Login(req.UserId, req.Pwd)
	//
	res := protocol.NewRespLogin()
	if !ok {
		res.Err = errors.ErrInvalidUser
	} else {
		s.client2Account[client.Id()] = account
		logger.Debug("account login", account.Name)
	}
	client.Write(res)
}

//
func (s *sAccount) OnReqLogout(ctx context.Context, req *protocol.ReqLogout, client base.Client) {
	logger.Debug("handle logout  ", client.Id())
	account := s.client2Account[client.Id()]
	if account != nil {
		account.Logout()
		delete(s.client2Account, client.Id())
	}
	client.Write(nil)
}

//
func (s *sAccount) OnReqRegister(ctx context.Context, req *protocol.ReqRegister, client base.Client) {
	logger.Debug("handle register  ", client.Id())
	account := modules.NewAccount("redis")
	ok := account.Register(req.UserId, req.Pwd)
	//
	res := protocol.NewRespLogin()
	if !ok {
		res.Err = errors.ErrInvalidUser
	} else {
		logger.Debug("account login", account.Name)
	}
	client.Write(res)
}
