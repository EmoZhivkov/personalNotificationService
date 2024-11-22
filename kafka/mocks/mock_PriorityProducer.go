// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	kafka "personalNotificationService/kafka"

	mock "github.com/stretchr/testify/mock"
)

// MockPriorityProducer is an autogenerated mock type for the PriorityProducer type
type MockPriorityProducer struct {
	mock.Mock
}

type MockPriorityProducer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockPriorityProducer) EXPECT() *MockPriorityProducer_Expecter {
	return &MockPriorityProducer_Expecter{mock: &_m.Mock}
}

// Close provides a mock function with given fields:
func (_m *MockPriorityProducer) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockPriorityProducer_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type MockPriorityProducer_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *MockPriorityProducer_Expecter) Close() *MockPriorityProducer_Close_Call {
	return &MockPriorityProducer_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *MockPriorityProducer_Close_Call) Run(run func()) *MockPriorityProducer_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockPriorityProducer_Close_Call) Return(_a0 error) *MockPriorityProducer_Close_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockPriorityProducer_Close_Call) RunAndReturn(run func() error) *MockPriorityProducer_Close_Call {
	_c.Call.Return(run)
	return _c
}

// SendMessage provides a mock function with given fields: _a0
func (_m *MockPriorityProducer) SendMessage(_a0 kafka.MessageWithPriority) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for SendMessage")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(kafka.MessageWithPriority) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockPriorityProducer_SendMessage_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendMessage'
type MockPriorityProducer_SendMessage_Call struct {
	*mock.Call
}

// SendMessage is a helper method to define mock.On call
//   - _a0 kafka.MessageWithPriority
func (_e *MockPriorityProducer_Expecter) SendMessage(_a0 interface{}) *MockPriorityProducer_SendMessage_Call {
	return &MockPriorityProducer_SendMessage_Call{Call: _e.mock.On("SendMessage", _a0)}
}

func (_c *MockPriorityProducer_SendMessage_Call) Run(run func(_a0 kafka.MessageWithPriority)) *MockPriorityProducer_SendMessage_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(kafka.MessageWithPriority))
	})
	return _c
}

func (_c *MockPriorityProducer_SendMessage_Call) Return(_a0 error) *MockPriorityProducer_SendMessage_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockPriorityProducer_SendMessage_Call) RunAndReturn(run func(kafka.MessageWithPriority) error) *MockPriorityProducer_SendMessage_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockPriorityProducer creates a new instance of MockPriorityProducer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockPriorityProducer(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockPriorityProducer {
	mock := &MockPriorityProducer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
