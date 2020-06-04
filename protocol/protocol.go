// Copyright 2009 The Das Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package all the protocol cmd defines
package protocol

import (
	"../base"
	"github.com/wonderivan/logger"
)

//
type BaseResponse struct {
	base.BaseCmd
	//
	Err int `json:"err,omitempty"`
}

//
type ReqJoinRoom struct {
	//
	base.BaseCmd
	//
	RoomIds []int `json:"roomIds"`
}

//
type RespJoinRoom struct {
	BaseResponse
	//
	ClientId uint64 `json:"clientId"`
	//
	Joins []int `json:"joins"`
}

//
type ReqLeaveRoom struct {
	//
	base.BaseCmd
	//
	RoomIds []int `json:"roomIds"`
}

//
type RespLeaveRoom struct {
	BaseResponse
	//
	ClientId uint64 `json:"clientId"`
	//
	Leaves []int `json:"leaves"`
}

//
type ReqCreateRoom struct {
	//
	base.BaseCmd
	//
	RoomId int `json:"roomId"`
}

//
type RespCreateRoom struct {
	BaseResponse
}

//
type ReqLogin struct {
	base.BaseCmd
	ServiceToken string `json:"serviceToken"` //
	AppId        uint32 `json:"appId,omitempty"`
	UserId       string `json:"userId,omitempty"`
	Pwd          string `json:"pwd,omitempty"`
}

//
type RespLogin struct {
	BaseResponse
}

//
type ReqLogout struct {
	base.BaseCmd
}

//
type ReqRegister struct {
	base.BaseCmd
	ServiceToken string `json:"serviceToken"` //
	AppId        uint32 `json:"appId,omitempty"`
	UserId       string `json:"userId,omitempty"`
	Pwd          string `json:"pwd,omitempty"`
}

//
type RespRegister struct {
	BaseResponse
}

//
type HeartBeat struct {
	base.BaseCmd
	UserId string `json:"userId,omitempty"`
}

//chat
type Chat struct {
	base.BaseCmd
	From uint64 `json:"from"`
	To   uint64 `json:"to"`
	Type int    `json:"type"`
	Word string `json:"word"`
}

//
type ReqResource struct {
	base.BaseCmd
	UserId string `json:"userId,omitempty"`
	//[add|delete|update|query]
	Method string                 `json:"method,omitempty"`
	Url    string                 `json:"url,omitempty"`
	ResId  string                 `json:"resId,omitempty"`
	Params map[string]interface{} `json:"q,omitempty"`
}

//
type RespResource struct {
	BaseResponse
	UserId string `json:"userId,omitempty"`
	//[add|delete|update|query]
	Method string      `json:"method,omitempty"`
	Url    string      `json:"url,omitempty"`
	Result interface{} `json:"result,omitempty"`
}

//
func init() {
	//add json pctl
	_, _ = base.RegisterClass(&Chat{})
	_, _ = base.RegisterClass(&ReqLogin{})
	_, _ = base.RegisterClass(&ReqLeaveRoom{})
	_, _ = base.RegisterClass(&ReqCreateRoom{})
	_, _ = base.RegisterClass(&RespJoinRoom{})
	_, _ = base.RegisterClass(&ReqLeaveRoom{})
	_, _ = base.RegisterClass(&HeartBeat{})
	_, _ = base.RegisterClass(&ReqLogout{})
	_, _ = base.RegisterClass(&ReqRegister{})
	//add other pctl

	logger.Debug("cmd register init ")
}
