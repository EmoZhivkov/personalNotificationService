// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	common "personalNotificationService/common"

	mock "github.com/stretchr/testify/mock"
)

// MocksenderConcreteImpl is an autogenerated mock type for the senderConcreteImpl type
type MocksenderConcreteImpl struct {
	mock.Mock
}

type MocksenderConcreteImpl_Expecter struct {
	mock *mock.Mock
}

func (_m *MocksenderConcreteImpl) EXPECT() *MocksenderConcreteImpl_Expecter {
	return &MocksenderConcreteImpl_Expecter{mock: &_m.Mock}
}

// Send provides a mock function with given fields: _a0, _a1
func (_m *MocksenderConcreteImpl) Send(_a0 interface{}, _a1 common.Content) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Send")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}, common.Content) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MocksenderConcreteImpl_Send_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Send'
type MocksenderConcreteImpl_Send_Call struct {
	*mock.Call
}

// Send is a helper method to define mock.On call
//   - _a0 interface{}
//   - _a1 common.Content
func (_e *MocksenderConcreteImpl_Expecter) Send(_a0 interface{}, _a1 interface{}) *MocksenderConcreteImpl_Send_Call {
	return &MocksenderConcreteImpl_Send_Call{Call: _e.mock.On("Send", _a0, _a1)}
}

func (_c *MocksenderConcreteImpl_Send_Call) Run(run func(_a0 interface{}, _a1 common.Content)) *MocksenderConcreteImpl_Send_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}), args[1].(common.Content))
	})
	return _c
}

func (_c *MocksenderConcreteImpl_Send_Call) Return(_a0 error) *MocksenderConcreteImpl_Send_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MocksenderConcreteImpl_Send_Call) RunAndReturn(run func(interface{}, common.Content) error) *MocksenderConcreteImpl_Send_Call {
	_c.Call.Return(run)
	return _c
}

// NewMocksenderConcreteImpl creates a new instance of MocksenderConcreteImpl. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMocksenderConcreteImpl(t interface {
	mock.TestingT
	Cleanup(func())
}) *MocksenderConcreteImpl {
	mock := &MocksenderConcreteImpl{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}