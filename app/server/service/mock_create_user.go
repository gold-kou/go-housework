// Code generated by MockGen. DO NOT EDIT.
// Source: create_user.go

// Package service is a generated GoMock package.
package service

import (
	gomock "github.com/golang/mock/gomock"
	db "github.com/gold-kou/go-housework/app/model/db"
	reflect "reflect"
)

// MockCreateUserServiceInterface is a mock of CreateUserServiceInterface interface
type MockCreateUserServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockCreateUserServiceInterfaceMockRecorder
}

// MockCreateUserServiceInterfaceMockRecorder is the mock recorder for MockCreateUserServiceInterface
type MockCreateUserServiceInterfaceMockRecorder struct {
	mock *MockCreateUserServiceInterface
}

// NewMockCreateUserServiceInterface creates a new mock instance
func NewMockCreateUserServiceInterface(ctrl *gomock.Controller) *MockCreateUserServiceInterface {
	mock := &MockCreateUserServiceInterface{ctrl: ctrl}
	mock.recorder = &MockCreateUserServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCreateUserServiceInterface) EXPECT() *MockCreateUserServiceInterfaceMockRecorder {
	return m.recorder
}

// Execute mocks base method
func (m *MockCreateUserServiceInterface) Execute() (*db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute")
	ret0, _ := ret[0].(*db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute
func (mr *MockCreateUserServiceInterfaceMockRecorder) Execute() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockCreateUserServiceInterface)(nil).Execute))
}