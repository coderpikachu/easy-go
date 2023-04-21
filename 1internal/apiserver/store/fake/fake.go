// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package fake

import (
	"fmt"
	"sync"

	v1 "easy-go/3api/apiserver/v1"

	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"

	"easy-go/1internal/apiserver/store"
)

// ResourceCount defines the number of fake resources.
const ResourceCount = 1000

type datastore struct {
	sync.RWMutex
	users []*v1.User
}

func (ds *datastore) Users() store.UserStore {
	return newUsers(ds)
}

func (ds *datastore) Close() error {
	return nil
}

var (
	fakeFactory store.Factory
	once        sync.Once
)

// GetFakeFactoryOr create fake store.
func GetFakeFactoryOr() (store.Factory, error) {
	once.Do(func() {
		fakeFactory = &datastore{
			users: FakeUsers(ResourceCount),
		}
	})

	if fakeFactory == nil {
		return nil, fmt.Errorf("failed to get mysql store fatory, mysqlFactory: %+v", fakeFactory)
	}

	return fakeFactory, nil
}

// FakeUsers returns fake user data.
func FakeUsers(count int) []*v1.User {
	// init some user records
	users := make([]*v1.User, 0)
	for i := 1; i <= count; i++ {
		users = append(users, &v1.User{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("user%d", i),
				ID:   uint64(i),
			},
			Nickname: fmt.Sprintf("user%d", i),
			Password: fmt.Sprintf("User%d@2020", i),
			Email:    fmt.Sprintf("user%d@qq.com", i),
		})
	}

	return users
}
