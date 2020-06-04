// Copyright 2009 The Das Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol

//
func NewRespLogin() *RespLogin {
	res := &RespLogin{}
	res.Id = "RespLogin"
	return res
}

//
func NewRespJoinRoom(clientId uint64, joins []int) *RespJoinRoom {
	res := &RespJoinRoom{
		ClientId: clientId,
		Joins:    joins,
	}

	res.Id = "RespJoinRoom"
	return res
}

//
func NewRespLeaveRoom(clientId uint64, leaves []int) *RespLeaveRoom {
	res := &RespLeaveRoom{
		ClientId: clientId,
		Leaves:   leaves,
	}

	res.Id = "RespLeaveRoom"
	return res
}

//
func NewRespCreatRoom(errId int) *RespCreateRoom {
	res := &RespCreateRoom{}

	res.Id = "RespCreateRoom"
	res.Err = 0
	return res
}

//
func NewRespResource(errId int, result interface{}) *RespResource {
	res := &RespResource{}

	res.Id = "RespResource"
	res.Err = 0
	res.Result = result
	return res
}
