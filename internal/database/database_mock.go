// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/database/database.go
//
// Generated by this command:
//
//	mockgen -source=./internal/database/database.go -destination=./internal/database/database_mock.go -package=database
//

// Package database is a generated GoMock package.
package database

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockDatabase is a mock of Database interface.
type MockDatabase struct {
	ctrl     *gomock.Controller
	recorder *MockDatabaseMockRecorder
	isgomock struct{}
}

// MockDatabaseMockRecorder is the mock recorder for MockDatabase.
type MockDatabaseMockRecorder struct {
	mock *MockDatabase
}

// NewMockDatabase creates a new mock instance.
func NewMockDatabase(ctrl *gomock.Controller) *MockDatabase {
	mock := &MockDatabase{ctrl: ctrl}
	mock.recorder = &MockDatabaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDatabase) EXPECT() *MockDatabaseMockRecorder {
	return m.recorder
}

// HandleQuery mocks base method.
func (m *MockDatabase) HandleQuery(ctx context.Context, queryStr string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleQuery", ctx, queryStr)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HandleQuery indicates an expected call of HandleQuery.
func (mr *MockDatabaseMockRecorder) HandleQuery(ctx, queryStr any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleQuery", reflect.TypeOf((*MockDatabase)(nil).HandleQuery), ctx, queryStr)
}
