// Copyright 2010 The Das Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//string utils
package utils

//string utils
import "bytes"

/*
join string into one
*/
func JoinStr(strings ...string) (res string) {
	var buffer bytes.Buffer
	for _, s := range strings {
		buffer.WriteString(s)
	}
	res = buffer.String()
	return
}
