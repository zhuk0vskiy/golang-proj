// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	dto "backend/src/internal/model/dto"
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "backend/src/internal/model"
)

// IInstrumentalistRepository is an autogenerated mock type for the IInstrumentalistRepository type
type IInstrumentalistRepository struct {
	mock.Mock
}

// Add provides a mock function with given fields: ctx, request
func (_m *IInstrumentalistRepository) Add(ctx context.Context, request *dto.AddInstrumentalistRequest) error {
	ret := _m.Called(ctx, request)

	if len(ret) == 0 {
		panic("no return value specified for Add")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *dto.AddInstrumentalistRequest) error); ok {
		r0 = rf(ctx, request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, request
func (_m *IInstrumentalistRepository) Delete(ctx context.Context, request *dto.DeleteInstrumentalistRequest) error {
	ret := _m.Called(ctx, request)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *dto.DeleteInstrumentalistRequest) error); ok {
		r0 = rf(ctx, request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, request
func (_m *IInstrumentalistRepository) Get(ctx context.Context, request *dto.GetInstrumentalistRequest) (*model.Instrumentalist, error) {
	ret := _m.Called(ctx, request)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *model.Instrumentalist
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *dto.GetInstrumentalistRequest) (*model.Instrumentalist, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *dto.GetInstrumentalistRequest) *model.Instrumentalist); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Instrumentalist)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *dto.GetInstrumentalistRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByStudio provides a mock function with given fields: ctx, request
func (_m *IInstrumentalistRepository) GetByStudio(ctx context.Context, request *dto.GetInstrumentalistByStudioRequest) ([]*model.Instrumentalist, error) {
	ret := _m.Called(ctx, request)

	if len(ret) == 0 {
		panic("no return value specified for GetByStudio")
	}

	var r0 []*model.Instrumentalist
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *dto.GetInstrumentalistByStudioRequest) ([]*model.Instrumentalist, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *dto.GetInstrumentalistByStudioRequest) []*model.Instrumentalist); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Instrumentalist)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *dto.GetInstrumentalistByStudioRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, request
func (_m *IInstrumentalistRepository) Update(ctx context.Context, request *dto.UpdateInstrumentalistRequest) error {
	ret := _m.Called(ctx, request)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *dto.UpdateInstrumentalistRequest) error); ok {
		r0 = rf(ctx, request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewIInstrumentalistRepository creates a new instance of IInstrumentalistRepository. It also registers a testing interface on the mock and a cleanup function to assert the pool expectations.
// The first argument is typically a *testing.T value.
func NewIInstrumentalistRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *IInstrumentalistRepository {
	mock := &IInstrumentalistRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
