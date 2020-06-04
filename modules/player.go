// Copyright 2009 The Das Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package modules

import (
	"../repository"
	"../utils"
)

type Player struct {
	//
	repo.BaseEntity `json:",inline" bson:",inline"`
	//
	Name string `json:"name" bson:"name"`
	//
	CreateTime int64 `json:"ctime" bson:"ctime"`
	//
	OwnerId string `json:"oid" bson:"oid"`
	//
	RoleIds []string `json:"roleids" bson:"roleids"`
	//
	ServerId string `json:"serverId" bson:"serverId"`
	//
	owner *Account
	//
	roles map[string]*Role
}

//
func NewPlayer(source string) *Player {
	a := &Player{}

	ent := repo.GetDataSource(source)
	if err := a.Init(ent, map[string]interface{}{
		"CollName": "player",
	}); err != nil {
		return nil
	}

	return a
}

//
func (r *Player) LoadById(id string) bool {
	return r.BaseEntity.LoadById(id, r)
}

//
func (r *Player) Save() bool {
	return r.BaseEntity.Save(r)
}

//
func (r *Player) SetOwner(a *Account) {
	r.owner = a
	r.OwnerId = a.Id
}

//
func (r *Player) GetOwner() *Account {
	if r.owner == nil {
		ac := NewAccount("redis")
		if !ac.LoadById(r.OwnerId, ac) {
			return nil
		}
		r.owner = ac
	}
	return r.owner
}

//
func (a *Player) AddRole(r *Role) {
	r.SetOwner(a)
	a.RoleIds = append(a.RoleIds, r.Id)
	a.roles[r.Id] = r
}

//
func (a *Player) DelRole(rid string) {
	role := a.roles[rid]
	if role != nil {
		role.SetOwner(nil)
	}

	utils.DelStrSlice(a.RoleIds, rid)
	delete(a.roles, rid)
}

//
func (a *Player) LoadRoles() {
	for _, rid := range a.RoleIds {
		r := NewRole("redis")
		if r.LoadById(rid) {
			a.AddRole(r)
		}
	}
}
