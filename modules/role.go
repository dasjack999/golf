// Copyright 2009 The Das Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package modules

import (
	"../repository"
)

type Role struct {
	//
	repo.BaseEntity `json:",inline" bson:",inline"`
	//
	Name string `json:"name" bson:"name"`
	//
	CreateTime int64 `json:"ctime" bson:"ctime"`
	//
	TableId string `json:"tid" bson:"tid"`
	//
	OwnerId string `json:"oid" bson:"oid"`
	//
	//Props map[string]interface{} `json:"props" bson:"props"`
	//
	owner *Player
}

//
func NewRole(source string) *Role {
	a := &Role{}

	ent := repo.GetDataSource(source)
	if err := a.Init(ent, map[string]interface{}{
		"CollName": "role",
	}); err != nil {
		return nil
	}

	return a
}

//
func (r *Role) LoadById(id string) bool {
	return r.BaseEntity.LoadById(id, r)
}

//
func (r *Role) Save() bool {
	return r.BaseEntity.Save(r)
}

//
func (r *Role) SetOwner(a *Player) {
	r.owner = a
	r.OwnerId = a.Id
}

//
func (r *Role) GetOwner() *Player {
	if r.owner == nil {
		ac := NewPlayer("redis")
		if !ac.LoadById(r.OwnerId) {
			return nil
		}
		r.owner = ac
	}
	return r.owner
}
