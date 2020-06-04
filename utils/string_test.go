// Copyright 2010 The Das Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package utils_test

import (
	"./"
	"testing"
)

func TestJoinStr(t *testing.T) {
	var want = "abcd"
	var res = utils.JoinStr("a", "b", "c", "d")

	if want != res {
		t.Errorf("JoinStr failed,%s != %s", want, res)
	}

}

/*
 */
func TestAppend(t *testing.T) {
	//a :=[]int {1,2,3}
	//b :=[]int {4,5,6}
	//
	//c :=[]int{}
	//c =append(c,a...)
	//t.Errorf("c=%d",c)
	//c =append(c,b...)
	//t.Errorf("c=%d",c)

}
func TestMapKeyOrder(t *testing.T) {
	mp := map[string]int{}
	mp["z"] = 1
	mp["a"] = 2

	if true {
		for k, v := range mp {
			t.Error("order", k, v)
		}

	}
}

/*
bentest
*/
func BenchmarkReverse(b *testing.B) {

}
