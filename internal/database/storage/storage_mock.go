// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/database/storage/storage.go
//
// Generated by this command:
//
//	mockgen -source=./internal/database/storage/storage.go -destination=./internal/database/storage/storage_mock.go -package=storage
//

// Package storage is a generated GoMock package.
package storage

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
	isgomock struct{}
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// Del mocks base method.
func (m *MockStorage) Del(arg0 context.Context, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Del", arg0, arg1)
}

// Del indicates an expected call of Del.
func (mr *MockStorageMockRecorder) Del(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Del", reflect.TypeOf((*MockStorage)(nil).Del), arg0, arg1)
}

// Get mocks base method.
func (m *MockStorage) Get(arg0 context.Context, arg1 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(string)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockStorageMockRecorder) Get(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockStorage)(nil).Get), arg0, arg1)
}

// Set mocks base method.
func (m *MockStorage) Set(arg0 context.Context, arg1, arg2 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Set", arg0, arg1, arg2)
}

// Set indicates an expected call of Set.
func (mr *MockStorageMockRecorder) Set(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockStorage)(nil).Set), arg0, arg1, arg2)
}
