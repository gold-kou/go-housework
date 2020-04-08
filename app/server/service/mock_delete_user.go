// Code generated by MockGen. DO NOT EDIT.
// Source: delete_user.go

// Package service is a generated GoMock package.
package service

import (
	gomock "github.com/golang/mock/gomock"
	model "github.com/gold-kou/go-housework/app/model"
	reflect "reflect"
)

// MockDeleteUserServiceInterface is a mock of DeleteUserServiceInterface interface
type MockDeleteUserServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockDeleteUserServiceInterfaceMockRecorder
}

// MockDeleteUserServiceInterfaceMockRecorder is the mock recorder for MockDeleteUserServiceInterface
type MockDeleteUserServiceInterfaceMockRecorder struct {
	mock *MockDeleteUserServiceInterface
}

// NewMockDeleteUserServiceInterface creates a new mock instance
func NewMockDeleteUserServiceInterface(ctrl *gomock.Controller) *MockDeleteUserServiceInterface {
	mock := &MockDeleteUserServiceInterface{ctrl: ctrl}
	mock.recorder = &MockDeleteUserServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDeleteUserServiceInterface) EXPECT() *MockDeleteUserServiceInterfaceMockRecorder {
	return m.recorder
}

// Execute mocks base method
func (m *MockDeleteUserServiceInterface) Execute(arg0 *model.Auth) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Execute indicates an expected call of Execute
func (mr *MockDeleteUserServiceInterfaceMockRecorder) Execute(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockDeleteUserServiceInterface)(nil).Execute), arg0)
}
