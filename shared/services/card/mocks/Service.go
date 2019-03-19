// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import card "github.com/dmibod/kanban/shared/services/card"
import context "context"
import kernel "github.com/dmibod/kanban/shared/kernel"
import mock "github.com/stretchr/testify/mock"

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// Create provides a mock function with given fields: _a0, _a1
func (_m *Service) Create(_a0 context.Context, _a1 *card.CreateModel) (kernel.ID, error) {
	ret := _m.Called(_a0, _a1)

	var r0 kernel.ID
	if rf, ok := ret.Get(0).(func(context.Context, *card.CreateModel) kernel.ID); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(kernel.ID)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *card.CreateModel) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Describe provides a mock function with given fields: _a0, _a1, _a2
func (_m *Service) Describe(_a0 context.Context, _a1 kernel.ID, _a2 string) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, kernel.ID, string) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAll provides a mock function with given fields: _a0
func (_m *Service) GetAll(_a0 context.Context) ([]*card.Model, error) {
	ret := _m.Called(_a0)

	var r0 []*card.Model
	if rf, ok := ret.Get(0).(func(context.Context) []*card.Model); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*card.Model)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: _a0, _a1
func (_m *Service) GetByID(_a0 context.Context, _a1 kernel.ID) (*card.Model, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *card.Model
	if rf, ok := ret.Get(0).(func(context.Context, kernel.ID) *card.Model); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*card.Model)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, kernel.ID) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByLaneID provides a mock function with given fields: _a0, _a1
func (_m *Service) GetByLaneID(_a0 context.Context, _a1 kernel.ID) ([]*card.Model, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []*card.Model
	if rf, ok := ret.Get(0).(func(context.Context, kernel.ID) []*card.Model); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*card.Model)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, kernel.ID) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Name provides a mock function with given fields: _a0, _a1, _a2
func (_m *Service) Name(_a0 context.Context, _a1 kernel.ID, _a2 string) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, kernel.ID, string) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Remove provides a mock function with given fields: _a0, _a1
func (_m *Service) Remove(_a0 context.Context, _a1 kernel.ID) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, kernel.ID) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
