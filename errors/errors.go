// Copyright 2010 The Das Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package errors

// errorString is a trivial implementation of error.
type BaseError struct {
	Code int
	Text string
}

//
func (e *BaseError) Error() string {
	return e.Text
}

//
func New(code int, a ...interface{}) *BaseError {
	var text string = ""
	for _, e := range a {
		text += e.(string)
	}
	return &BaseError{
		Code: code,
		Text: text,
	}
}

//error ids
const (
	ErrInvalidData = 1
	ErrInvalidUser = 2
	ErrResDelete   = 3
	ErrResQuery    = 4
	ErrResUpdate   = 5
	ErrResAdd      = 6
)
