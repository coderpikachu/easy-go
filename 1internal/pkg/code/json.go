// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package code

//go:generate codegen -type=int

// json: json errors.
const (
	// ErrJsonUnmarshal - 401: User password incorrect.
	ErrJsonUnmarshal int = iota + 120001
)
