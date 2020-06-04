// Copyright 2010 The Das Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package repo

import (
	"crypto/md5"
	"encoding/hex"
)

//
type DataSource interface {
	//
	Query(collName string, query map[string]interface{}) (result interface{}, err error)
	//
	Remove(collName string, id string) bool
	//
	Update(collName string, id string, v interface{}) bool
	//
	Insert(collName string, id string, v interface{}) bool
	//
	Init(data map[string]interface{}) error
	//
	Load(collName string, id string, v interface{}) bool
}

//
var dsMap map[string]DataSource = map[string]DataSource{}

//
func RegisterDataSource(name string, entityer DataSource) {
	dsMap[name] = entityer
}

//
func GetDataSource(name string) DataSource {
	return dsMap[name]
}

//
func GetMd5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

//
type BaseEntity struct {
	//
	Loaded bool `json:"-" bson:"-"`
	//
	CollName string `json:"-" bson:"-"`
	//
	Id string `json:"id"`
	//
	Extension interface{} `json:"extension,omitempty" bson:"extension,omitempty"`
	//repo.MongoSource `json:",inline" bson:",inline"`
	ent DataSource `json:"-" bson:"-"`
}

//
func (e *BaseEntity) Init(ent DataSource, query map[string]interface{}) error {
	e.ent = ent
	e.CollName = query["CollName"].(string)
	return e.ent.Init(query)
}

//
func (e *BaseEntity) Load(id string, v interface{}) (ok bool) {

	ok = e.ent.Load(e.CollName, id, v)
	if ok {
		e.Loaded = true
	}
	return
}

//
//
func (e *BaseEntity) LoadById(id string, v interface{}) (ok bool) {
	return e.Load(id, v)
}

//
func (e *BaseEntity) Save(v interface{}) bool {
	if e.Loaded {
		return e.ent.Update(e.CollName, e.Id, v)
	} else {
		return e.ent.Insert(e.CollName, e.Id, v)
	}
}

//
func (e *BaseEntity) Remove() bool {
	ok := e.ent.Remove(e.CollName, e.Id)

	e.ent = nil
	return ok
}
