// Copyright 2009 The Das Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package route implements task service
package service

import (
	"../base"
	"../errors"
	"../protocol"
	"../repository"
	"../utils"
	"context"
	"github.com/wonderivan/logger"
)

//
func init() {
	logger.Debug("s_chat init")

}

//
type sResource struct {
	//
	BaseService
}

//
var SResource = &sResource{}

//
func (s *sResource) OnReqResource(ctx context.Context, req *protocol.ReqResource, client base.Client) {
	//logger.Debug("handle req ", req.From, req.To, req.Word)
	switch req.Method {
	case "query":
		s.QueryResource(ctx, req, client)
	case "add":
		s.AddResource(ctx, req, client)
	case "delete":
		s.DeleteResource(ctx, req, client)
	case "update":
		s.UpdateResource(ctx, req, client)
	default:
		logger.Error("no such method", req.Method)
	}
}

//
func (s *sResource) QueryResource(ctx context.Context, req *protocol.ReqResource, client base.Client) {
	//logger.Debug("handle req ", req.From, req.To, req.Word)
	ent := repo.GetDataSource("redis")
	if ent == nil {
		return
	}

	//if result, err := ent.Query(req.Url, req.Params); err == nil {
	//	client.Write(protocol.NewRespResource(0, result))
	//}

	result, err := ent.Query(req.Url, req.Params)
	errId := 0
	if err != nil {
		errId = errors.ErrResQuery
	}
	client.Write(protocol.NewRespResource(errId, result))
}

//
func (s *sResource) AddResource(ctx context.Context, req *protocol.ReqResource, client base.Client) {
	//logger.Debug("handle req ", req.From, req.To, req.Word)
	ent := repo.GetDataSource("redis")
	if ent == nil {
		return
	}

	ok := ent.Insert(req.Url, utils.GetGuid(), req.Params)
	errId := 0
	if !ok {
		errId = errors.ErrResAdd
	}
	client.Write(protocol.NewRespResource(errId, nil))
}

//
func (s *sResource) DeleteResource(ctx context.Context, req *protocol.ReqResource, client base.Client) {
	//logger.Debug("handle req ", req.From, req.To, req.Word)
	ent := repo.GetDataSource("redis")
	if ent == nil {
		return
	}

	ok := ent.Remove(req.Url, req.ResId)
	errId := 0
	if !ok {
		errId = errors.ErrResDelete
	}
	client.Write(protocol.NewRespResource(errId, nil))
}

//
func (s *sResource) UpdateResource(ctx context.Context, req *protocol.ReqResource, client base.Client) {
	//logger.Debug("handle req ", req.From, req.To, req.Word)
	ent := repo.GetDataSource("redis")
	if ent == nil {
		return
	}

	ok := ent.Update(req.Url, req.ResId, req.Params)
	errId := 0
	if !ok {
		errId = errors.ErrResUpdate
	}
	client.Write(protocol.NewRespResource(errId, nil))
}
