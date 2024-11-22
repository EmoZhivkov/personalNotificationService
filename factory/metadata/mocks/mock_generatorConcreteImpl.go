// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	common "personalNotificationService/common"

	mock "github.com/stretchr/testify/mock"
)

// MockgeneratorConcreteImpl is an autogenerated mock type for the generatorConcreteImpl type
type MockgeneratorConcreteImpl struct {
	mock.Mock
}

type MockgeneratorConcreteImpl_Expecter struct {
	mock *mock.Mock
}

func (_m *MockgeneratorConcreteImpl) EXPECT() *MockgeneratorConcreteImpl_Expecter {
	return &MockgeneratorConcreteImpl_Expecter{mock: &_m.Mock}
}

// GenerateMetadata provides a mock function with given fields:
func (_m *MockgeneratorConcreteImpl) GenerateMetadata() (common.Metadata, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GenerateMetadata")
	}

	var r0 common.Metadata
	var r1 error
	if rf, ok := ret.Get(0).(func() (common.Metadata, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() common.Metadata); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Metadata)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockgeneratorConcreteImpl_GenerateMetadata_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GenerateMetadata'
type MockgeneratorConcreteImpl_GenerateMetadata_Call struct {
	*mock.Call
}

// GenerateMetadata is a helper method to define mock.On call
func (_e *MockgeneratorConcreteImpl_Expecter) GenerateMetadata() *MockgeneratorConcreteImpl_GenerateMetadata_Call {
	return &MockgeneratorConcreteImpl_GenerateMetadata_Call{Call: _e.mock.On("GenerateMetadata")}
}

func (_c *MockgeneratorConcreteImpl_GenerateMetadata_Call) Run(run func()) *MockgeneratorConcreteImpl_GenerateMetadata_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockgeneratorConcreteImpl_GenerateMetadata_Call) Return(_a0 common.Metadata, _a1 error) *MockgeneratorConcreteImpl_GenerateMetadata_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockgeneratorConcreteImpl_GenerateMetadata_Call) RunAndReturn(run func() (common.Metadata, error)) *MockgeneratorConcreteImpl_GenerateMetadata_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockgeneratorConcreteImpl creates a new instance of MockgeneratorConcreteImpl. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockgeneratorConcreteImpl(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockgeneratorConcreteImpl {
	mock := &MockgeneratorConcreteImpl{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
