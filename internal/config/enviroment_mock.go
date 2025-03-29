// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/config/enviroment.go
//
// Generated by this command:
//
//	mockgen -source=./internal/config/enviroment.go -destination=./internal/config/enviroment_mock.go -package=config
//

// Package config is a generated GoMock package.
package config

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockOS is a mock of OS interface.
type MockOS struct {
	ctrl     *gomock.Controller
	recorder *MockOSMockRecorder
	isgomock struct{}
}

// MockOSMockRecorder is the mock recorder for MockOS.
type MockOSMockRecorder struct {
	mock *MockOS
}

// NewMockOS creates a new mock instance.
func NewMockOS(ctrl *gomock.Controller) *MockOS {
	mock := &MockOS{ctrl: ctrl}
	mock.recorder = &MockOSMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOS) EXPECT() *MockOSMockRecorder {
	return m.recorder
}

// GetEnv mocks base method.
func (m *MockOS) GetEnv(arg0 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEnv", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetEnv indicates an expected call of GetEnv.
func (mr *MockOSMockRecorder) GetEnv(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEnv", reflect.TypeOf((*MockOS)(nil).GetEnv), arg0)
}

// ReadFile mocks base method.
func (m *MockOS) ReadFile(arg0 string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadFile", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadFile indicates an expected call of ReadFile.
func (mr *MockOSMockRecorder) ReadFile(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadFile", reflect.TypeOf((*MockOS)(nil).ReadFile), arg0)
}
