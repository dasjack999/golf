// Copyright 2010 The Das Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//the base structs and interface
package base

//
import (
	"../errors"
	"github.com/wonderivan/logger"
	"reflect"
	"runtime"
)

//
type Str2Any map[string]interface{}

//the root object,all struct should extends from it
type Base struct {
	id int
}

//the dctor interface
func (b *Base) SetFinal(dtor func(*Base)) {
	runtime.SetFinalizer(b, dtor)
	//func(obj *Base) {
	//	logger.Debug("obj gc:%v", obj)
	//}
}

//the global class map for runtime to get class type
var clsMap = map[string]reflect.Type{}

//register global class info
//for dynmic create objects
func RegisterClass(cls interface{}) (name string, err error) {
	t := reflect.TypeOf(cls)
	if t == nil || t.Kind() != reflect.Ptr {
		logger.Debug("RegisterClass:should use point of init class instance", t)
		return "", errors.New(errors.ErrInvalidData, "class without elem")
	}
	//if te ==nil{
	//
	//	return "",errors.New(errors.Err_invalid_data,"class without elem")
	//}
	name = t.Elem().Name()
	//
	clsMap[name] = t

	return
}

//get class type in runtime
//
func GetClass(id string) (t reflect.Type, ok bool) {
	t, ok = clsMap[id]
	return
}

//get instance by class name
func NewInstance(id string) interface{} {
	t, ok := GetClass(id)
	if !ok {
		logger.Debug("UnPack:without regclass,so as interface{}")
		return nil
	}
	obj := reflect.New(t.Elem()).Interface()
	return obj
}
