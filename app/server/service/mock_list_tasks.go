// Code generated by MockGen. DO NOT EDIT.
// Source: list_tasks.go

// Package service is a generated GoMock package.
package service

import (
	gomock "github.com/golang/mock/gomock"
	model "github.com/gold-kou/go-housework/app/model"
	db "github.com/gold-kou/go-housework/app/model/db"
	reflect "reflect"
)

// MockListTasksServiceInterface is a mock of ListTasksServiceInterface interface
type MockListTasksServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockListTasksServiceInterfaceMockRecorder
}

// MockListTasksServiceInterfaceMockRecorder is the mock recorder for MockListTasksServiceInterface
type MockListTasksServiceInterfaceMockRecorder struct {
	mock *MockListTasksServiceInterface
}

// NewMockListTasksServiceInterface creates a new mock instance
func NewMockListTasksServiceInterface(ctrl *gomock.Controller) *MockListTasksServiceInterface {
	mock := &MockListTasksServiceInterface{ctrl: ctrl}
	mock.recorder = &MockListTasksServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockListTasksServiceInterface) EXPECT() *MockListTasksServiceInterfaceMockRecorder {
	return m.recorder
}

// Execute mocks base method
func (m *MockListTasksServiceInterface) Execute(arg0 *model.Auth, arg1 string) ([]*db.Task, *db.Family, []*db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", arg0, arg1)
	ret0, _ := ret[0].([]*db.Task)
	ret1, _ := ret[1].(*db.Family)
	ret2, _ := ret[2].([]*db.User)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// Execute indicates an expected call of Execute
func (mr *MockListTasksServiceInterfaceMockRecorder) Execute(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockListTasksServiceInterface)(nil).Execute), arg0, arg1)
}
