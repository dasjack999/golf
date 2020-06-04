// Copyright 2010 The Das Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package utils

import (
	"github.com/wonderivan/logger"
	"os"
	"path/filepath"
	"strings"
)

//
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[1]))
	if err != nil {
		logger.Debug(err)
	}
	return strings.Replace(dir, "\\", "/", -1)

}
