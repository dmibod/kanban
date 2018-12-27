// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "context"
import db "github.com/dmibod/kanban/shared/tools/db"
import mock "github.com/stretchr/testify/mock"

// RepositoryFactory is an autogenerated mock type for the RepositoryFactory type
type RepositoryFactory struct {
	mock.Mock
}

// CreateRepository provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *RepositoryFactory) CreateRepository(_a0 context.Context, _a1 string, _a2 db.InstanceFactory, _a3 db.InstanceIdentity) db.Repository {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	var r0 db.Repository
	if rf, ok := ret.Get(0).(func(context.Context, string, db.InstanceFactory, db.InstanceIdentity) db.Repository); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(db.Repository)
		}
	}

	return r0
}