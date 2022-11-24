// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/KirillMironov/tau/runtimes (interfaces: ContainerRuntime)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	runtimes "github.com/KirillMironov/tau/runtimes"
	gomock "github.com/golang/mock/gomock"
)

// MockContainerRuntime is a mock of ContainerRuntime interface.
type MockContainerRuntime struct {
	ctrl     *gomock.Controller
	recorder *MockContainerRuntimeMockRecorder
}

// MockContainerRuntimeMockRecorder is the mock recorder for MockContainerRuntime.
type MockContainerRuntimeMockRecorder struct {
	mock *MockContainerRuntime
}

// NewMockContainerRuntime creates a new mock instance.
func NewMockContainerRuntime(ctrl *gomock.Controller) *MockContainerRuntime {
	mock := &MockContainerRuntime{ctrl: ctrl}
	mock.recorder = &MockContainerRuntimeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockContainerRuntime) EXPECT() *MockContainerRuntimeMockRecorder {
	return m.recorder
}

// Remove mocks base method.
func (m *MockContainerRuntime) Remove(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove.
func (mr *MockContainerRuntimeMockRecorder) Remove(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockContainerRuntime)(nil).Remove), arg0)
}

// Start mocks base method.
func (m *MockContainerRuntime) Start(arg0 runtimes.Container) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockContainerRuntimeMockRecorder) Start(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockContainerRuntime)(nil).Start), arg0)
}
