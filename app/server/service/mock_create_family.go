// Code generated by MockGen. DO NOT EDIT.
// Source: create_family.go

// Package service is a generated GoMock package.
package service

import (
	gomock "github.com/golang/mock/gomock"
	model "github.com/gold-kou/go-housework/app/model"
	db "github.com/gold-kou/go-housework/app/model/db"
	schemamodel "github.com/gold-kou/go-housework/app/model/schemamodel"
	reflect "reflect"
)

// MockCreateFamilyServiceInterface is a mock of CreateFamilyServiceInterface interface
type MockCreateFamilyServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockCreateFamilyServiceInterfaceMockRecorder
}

// MockCreateFamilyServiceInterfaceMockRecorder is the mock recorder for MockCreateFamilyServiceInterface
type MockCreateFamilyServiceInterfaceMockRecorder struct {
	mock *MockCreateFamilyServiceInterface
}

// NewMockCreateFamilyServiceInterface creates a new mock instance
func NewMockCreateFamilyServiceInterface(ctrl *gomock.Controller) *MockCreateFamilyServiceInterface {
	mock := &MockCreateFamilyServiceInterface{ctrl: ctrl}
	mock.recorder = &MockCreateFamilyServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCreateFamilyServiceInterface) EXPECT() *MockCreateFamilyServiceInterfaceMockRecorder {
	return m.recorder
}

// Execute mocks base method
func (m *MockCreateFamilyServiceInterface) Execute(arg0 *model.Auth, arg1 *schemamodel.RequestCreateFamily) (*db.Family, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", arg0, arg1)
	ret0, _ := ret[0].(*db.Family)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute
func (mr *MockCreateFamilyServiceInterfaceMockRecorder) Execute(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockCreateFamilyServiceInterface)(nil).Execute), arg0, arg1)
}
