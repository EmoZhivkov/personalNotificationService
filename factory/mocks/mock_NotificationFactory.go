// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	uuid "github.com/google/uuid"
	mock "github.com/stretchr/testify/mock"
)

// MockNotificationFactory is an autogenerated mock type for the NotificationFactory type
type MockNotificationFactory struct {
	mock.Mock
}

type MockNotificationFactory_Expecter struct {
	mock *mock.Mock
}

func (_m *MockNotificationFactory) EXPECT() *MockNotificationFactory_Expecter {
	return &MockNotificationFactory_Expecter{mock: &_m.Mock}
}

// CreateNotification provides a mock function with given fields: targetUsername, notificationID
func (_m *MockNotificationFactory) CreateNotification(targetUsername string, notificationID uuid.UUID) error {
	ret := _m.Called(targetUsername, notificationID)

	if len(ret) == 0 {
		panic("no return value specified for CreateNotification")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, uuid.UUID) error); ok {
		r0 = rf(targetUsername, notificationID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockNotificationFactory_CreateNotification_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateNotification'
type MockNotificationFactory_CreateNotification_Call struct {
	*mock.Call
}

// CreateNotification is a helper method to define mock.On call
//   - targetUsername string
//   - notificationID uuid.UUID
func (_e *MockNotificationFactory_Expecter) CreateNotification(targetUsername interface{}, notificationID interface{}) *MockNotificationFactory_CreateNotification_Call {
	return &MockNotificationFactory_CreateNotification_Call{Call: _e.mock.On("CreateNotification", targetUsername, notificationID)}
}

func (_c *MockNotificationFactory_CreateNotification_Call) Run(run func(targetUsername string, notificationID uuid.UUID)) *MockNotificationFactory_CreateNotification_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MockNotificationFactory_CreateNotification_Call) Return(_a0 error) *MockNotificationFactory_CreateNotification_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockNotificationFactory_CreateNotification_Call) RunAndReturn(run func(string, uuid.UUID) error) *MockNotificationFactory_CreateNotification_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockNotificationFactory creates a new instance of MockNotificationFactory. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockNotificationFactory(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockNotificationFactory {
	mock := &MockNotificationFactory{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
