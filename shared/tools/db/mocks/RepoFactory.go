// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import db "github.com/dmibod/kanban/shared/tools/db"
import mock "github.com/stretchr/testify/mock"

// RepoFactory is an autogenerated mock type for the RepoFactory type
type RepoFactory struct {
	mock.Mock
}

// Create provides a mock function with given fields: _a0, _a1
func (_m *RepoFactory) Create(_a0 string, _a1 db.InstanceFactory) db.Repository {
	ret := _m.Called(_a0, _a1)

	var r0 db.Repository
	if rf, ok := ret.Get(0).(func(string, db.InstanceFactory) db.Repository); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(db.Repository)
		}
	}

	return r0
}
