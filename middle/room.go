// Copyright 2009 The Das Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package route implements dispatch request.
//
package middle

import (
	"../base"
	"../errors"
	"../protocol"
	"context"
	"github.com/wonderivan/logger"
	"reflect"
)

//
func init() {
	logger.Debug("MdRoom init")
	sRoom := MdRoom{
		rooms: map[int]ClientList{},
	}
	sRoom.Name = reflect.TypeOf(sRoom).Name()
	base.RegisterMiddle(&sRoom)
}

//client id list
type ClientList []uint64

//
type MdRoom struct {
	base.BaseMiddle
	//
	rooms map[int]ClientList
}

//
func (s *MdRoom) HandleRequest(ctx context.Context, cmd base.Cmd, client base.Client) (handled bool, err error) {

	switch cmd.(type) {
	case *base.DisConnected:
		s.OnReqLeaveRoom(ctx, &protocol.ReqLeaveRoom{
			RoomIds: nil,
		}, client)
	case *protocol.ReqCreateRoom:
		s.OnReqCreateRoom(ctx, cmd.(*protocol.ReqCreateRoom), client)
	case *protocol.ReqLeaveRoom:
		s.OnReqLeaveRoom(ctx, cmd.(*protocol.ReqLeaveRoom), client)
	case *protocol.ReqJoinRoom:
		s.OnReqJoinRoom(ctx, cmd.(*protocol.ReqJoinRoom), client)
	}
	return
}

//
func (s *MdRoom) OnReqJoinRoom(ctx context.Context, req *protocol.ReqJoinRoom, client base.Client) {
	logger.Debug("handle join room ", req.RoomIds, client.Id())
	//
	joins, _ := s.Join(client.Id(), req.RoomIds...)
	res := protocol.NewRespJoinRoom(client.Id(), joins)

	others := make([]uint64, 0)
	for _, rid := range req.RoomIds {
		others = append(others, s.GetRoomMeets(rid)...)
	}

	client.WriteTo(res, others...)
}

//
func (s *MdRoom) OnReqCreateRoom(ctx context.Context, req *protocol.ReqCreateRoom, client base.Client) {
	logger.Debug("handle create room ", req.RoomId, client.Id())
	//
	err := s.Create(client.Id(), req.RoomId)
	eCode := 0
	if err != nil {
		eCode = err.Code
	}
	res := protocol.NewRespCreatRoom(eCode)
	client.Write(res)
}

//
func (s *MdRoom) OnReqLeaveRoom(ctx context.Context, req *protocol.ReqLeaveRoom, client base.Client) {
	logger.Debug("handle leave room ", req.RoomIds, client.Id())
	//others := make([]uint64, 0)
	//for _, rid := range req.RoomIds {
	//	others = append(others, s.GetRoomMeets(rid)...)
	//}
	//
	leaves, _ := s.Leave(client.Id(), req.RoomIds...)
	res := protocol.NewRespLeaveRoom(client.Id(), leaves)
	others := make([]uint64, 0)
	for _, rid := range leaves {
		others = append(others, s.GetRoomMeets(rid)...)
	}
	client.WriteTo(res, others...)
}

//
func (s *MdRoom) GetRoomMeets(roomId int) (joins []uint64) {
	joins = s.rooms[roomId]
	return
}

//
func (s *MdRoom) Join(clientId uint64, roomIds ...int) (joins []int, err error) {
	for _, rid := range roomIds {
		list, ok := s.rooms[rid]
		if !ok {
			logger.Debug("join room:room not exist", rid)
			//return errors.New("join room:room not exist")
		} else {
			joins = append(joins, rid)
			s.rooms[rid] = append(list, clientId)
		}
	}
	return
}

//
func (s *MdRoom) Leave(clientId uint64, roomIds ...int) (leaves []int, err error) {
	//all rooms
	if len(roomIds) == 0 {
		roomIds = s.getRoomIds()
	}
	//
	for _, rid := range roomIds {
		list, ok := s.rooms[rid]
		if !ok {
			logger.Debug("leave room:room not exist", rid)
		} else {
			for i := 0; i < len(list); i++ {
				if list[i] == clientId {
					leaves = append(leaves, rid)
					list = append(list[:i], list[i+1:]...)
					i-- // maintain the correct index
				}
			}
			//all leaves
			if len(list) == 0 {
				delete(s.rooms, rid)
			}
		}
	}
	return
}

//
func (s *MdRoom) Create(clientId uint64, roomId int) *errors.BaseError {
	_, ok := s.rooms[roomId]
	if ok {
		logger.Debug("create room:room exist already", roomId)
		return errors.New(errors.ErrInvalidData, "create room:room exist already")
	}
	//
	s.rooms[roomId] = make([]uint64, 1)
	s.rooms[roomId] = append(s.rooms[roomId], clientId)
	return nil
}

//
func (s *MdRoom) getRoomIds() (res []int) {
	for id, _ := range s.rooms {
		res = append(res, id)
	}
	return
}
