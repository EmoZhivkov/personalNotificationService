// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	slack "github.com/slack-go/slack"
	mock "github.com/stretchr/testify/mock"
)

// MockSlackClient is an autogenerated mock type for the SlackClient type
type MockSlackClient struct {
	mock.Mock
}

type MockSlackClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockSlackClient) EXPECT() *MockSlackClient_Expecter {
	return &MockSlackClient_Expecter{mock: &_m.Mock}
}

// PostMessage provides a mock function with given fields: channelID, options
func (_m *MockSlackClient) PostMessage(channelID string, options ...slack.MsgOption) (string, string, error) {
	_va := make([]interface{}, len(options))
	for _i := range options {
		_va[_i] = options[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, channelID)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for PostMessage")
	}

	var r0 string
	var r1 string
	var r2 error
	if rf, ok := ret.Get(0).(func(string, ...slack.MsgOption) (string, string, error)); ok {
		return rf(channelID, options...)
	}
	if rf, ok := ret.Get(0).(func(string, ...slack.MsgOption) string); ok {
		r0 = rf(channelID, options...)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string, ...slack.MsgOption) string); ok {
		r1 = rf(channelID, options...)
	} else {
		r1 = ret.Get(1).(string)
	}

	if rf, ok := ret.Get(2).(func(string, ...slack.MsgOption) error); ok {
		r2 = rf(channelID, options...)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockSlackClient_PostMessage_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PostMessage'
type MockSlackClient_PostMessage_Call struct {
	*mock.Call
}

// PostMessage is a helper method to define mock.On call
//   - channelID string
//   - options ...slack.MsgOption
func (_e *MockSlackClient_Expecter) PostMessage(channelID interface{}, options ...interface{}) *MockSlackClient_PostMessage_Call {
	return &MockSlackClient_PostMessage_Call{Call: _e.mock.On("PostMessage",
		append([]interface{}{channelID}, options...)...)}
}

func (_c *MockSlackClient_PostMessage_Call) Run(run func(channelID string, options ...slack.MsgOption)) *MockSlackClient_PostMessage_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]slack.MsgOption, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(slack.MsgOption)
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockSlackClient_PostMessage_Call) Return(_a0 string, _a1 string, _a2 error) *MockSlackClient_PostMessage_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *MockSlackClient_PostMessage_Call) RunAndReturn(run func(string, ...slack.MsgOption) (string, string, error)) *MockSlackClient_PostMessage_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockSlackClient creates a new instance of MockSlackClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockSlackClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockSlackClient {
	mock := &MockSlackClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}