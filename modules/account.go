// Copyright 2009 The Das Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package modules

import (
	"../repository"
	"../utils"
)

//
type Account struct {
	//
	repo.BaseEntity `json:",inline" bson:",inline"`
	//
	Name string `json:"name" bson:"name"`
	//
	Pwd string `json:"pwd" bson:"pwd"`
	//
	CreateTime int64 `json:"ctime" bson:"ctime"`
	//
	LastLoginTime int64 `json:"ltime" bson:"ltime"`
	//
	PlayerIds []string `json:"roleids" bson:"roleids"`
	//
	players map[string]*Player `json:"-" bson:"-"`
}

//
func NewAccount(source string) *Account {
	a := &Account{}
	a.players = make(map[string]*Player, 0)
	ent := repo.GetDataSource(source)
	if err := a.Init(ent, map[string]interface{}{
		"CollName": "account",
	}); err != nil {
		return nil
	}

	return a
}

//
func (a *Account) LoadByName(name string) bool {
	err := a.Load(map[string]interface{}{
		"id": repo.GetMd5(name),
	}, a)
	if err {

	}
	return err
}

//
func (a *Account) SaveByName() bool {
	a.Id = repo.GetMd5(a.Name)
	return a.Save(a)
}

//
func (a *Account) Register(name string, pwd string) bool {
	if a.CanRegister(name) {
		a.Name = name
		a.Pwd = pwd
		a.SaveByName()
	}
	return false
}

//
func (a *Account) Login(name string, pwd string) bool {
	if a.Loaded {
		return a.Name == name && a.Pwd == pwd
	} else {
		ok := a.LoadByName(name)
		if !ok {
			return false
		}
		//return a.Login(name,pwd)
		return a.Name == name && a.Pwd == pwd
	}
}

//
func (a *Account) Logout() {
	if !a.Loaded {

	} else {

	}

}

//
func (a *Account) CanRegister(name string) bool {
	return !a.LoadByName(name)
}

//
func (a *Account) AddPlayer(r *Player) {
	r.SetOwner(a)
	a.PlayerIds = append(a.PlayerIds, r.Id)
	a.players[r.Id] = r
}

//
func (a *Account) DelPlayer(rid string) {
	role := a.players[rid]
	if role != nil {
		role.SetOwner(nil)
	}
	utils.DelStrSlice(a.PlayerIds, rid)
	delete(a.players, rid)
}

//
func (a *Account) LoadPlayers() {
	for _, rid := range a.PlayerIds {
		r := NewPlayer("redis")
		if r.LoadById(rid) {
			r.LoadRoles()
			a.AddPlayer(r)
		}
	}
}

//
func (a *Account) LoadPlayerById(pid string) {
	//
	if _, ok := a.players[pid]; !ok {
		return
	}
	//
	if a.players[pid] == nil {
		r := NewPlayer("redis")
		if r.LoadById(pid) {
			r.LoadRoles()
			a.AddPlayer(r)
		}
	}
}
