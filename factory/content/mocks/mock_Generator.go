// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	common "personalNotificationService/common"

	mock "github.com/stretchr/testify/mock"

	repositories "personalNotificationService/repositories"
)

// MockGenerator is an autogenerated mock type for the Generator type
type MockGenerator struct {
	mock.Mock
}

type MockGenerator_Expecter struct {
	mock *mock.Mock
}

func (_m *MockGenerator) EXPECT() *MockGenerator_Expecter {
	return &MockGenerator_Expecter{mock: &_m.Mock}
}

// GenerateContent provides a mock function with given fields: _a0, _a1
func (_m *MockGenerator) GenerateContent(_a0 common.Metadata, _a1 repositories.Template) (common.Content, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GenerateContent")
	}

	var r0 common.Content
	var r1 error
	if rf, ok := ret.Get(0).(func(common.Metadata, repositories.Template) (common.Content, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(common.Metadata, repositories.Template) common.Content); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(common.Content)
	}

	if rf, ok := ret.Get(1).(func(common.Metadata, repositories.Template) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGenerator_GenerateContent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GenerateContent'
type MockGenerator_GenerateContent_Call struct {
	*mock.Call
}

// GenerateContent is a helper method to define mock.On call
//   - _a0 common.Metadata
//   - _a1 repositories.Template
func (_e *MockGenerator_Expecter) GenerateContent(_a0 interface{}, _a1 interface{}) *MockGenerator_GenerateContent_Call {
	return &MockGenerator_GenerateContent_Call{Call: _e.mock.On("GenerateContent", _a0, _a1)}
}

func (_c *MockGenerator_GenerateContent_Call) Run(run func(_a0 common.Metadata, _a1 repositories.Template)) *MockGenerator_GenerateContent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(common.Metadata), args[1].(repositories.Template))
	})
	return _c
}

func (_c *MockGenerator_GenerateContent_Call) Return(_a0 common.Content, _a1 error) *MockGenerator_GenerateContent_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGenerator_GenerateContent_Call) RunAndReturn(run func(common.Metadata, repositories.Template) (common.Content, error)) *MockGenerator_GenerateContent_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockGenerator creates a new instance of MockGenerator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockGenerator(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockGenerator {
	mock := &MockGenerator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}