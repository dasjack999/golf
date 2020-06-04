// Copyright 2010 The Das Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//file system utils
package utils

import (
	"github.com/wonderivan/logger"
	"io/ioutil"
	"os"
)

//
type WalkFunc func(file os.FileInfo)

//
func Walk(path string, w WalkFunc) {
	fileInfo, err := ioutil.ReadDir(path)
	if err != nil {
		logger.Debug(err)
		return
	}
	for _, info := range fileInfo {
		w(info)
		if info.IsDir() {
			Walk(path+"\\"+info.Name(), w)
		}
	}
}
