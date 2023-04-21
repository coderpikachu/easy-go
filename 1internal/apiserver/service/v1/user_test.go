// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1

import (
	"os"
	"testing"

	"easy-go/1internal/apiserver/store/fake"
)

func TestMain(m *testing.M) {
	_, _ = fake.GetFakeFactoryOr()
	os.Exit(m.Run())
}
