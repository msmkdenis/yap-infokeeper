// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/msmkdenis/yap-infokeeper/internal/credential/service (interfaces: CredentialRepository)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	specification "github.com/msmkdenis/yap-infokeeper/internal/credential/specification"
	model "github.com/msmkdenis/yap-infokeeper/internal/model"
)

// MockCredentialRepository is a mock of CredentialRepository interface.
type MockCredentialRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCredentialRepositoryMockRecorder
}

// MockCredentialRepositoryMockRecorder is the mock recorder for MockCredentialRepository.
type MockCredentialRepositoryMockRecorder struct {
	mock *MockCredentialRepository
}

// NewMockCredentialRepository creates a new mock instance.
func NewMockCredentialRepository(ctrl *gomock.Controller) *MockCredentialRepository {
	mock := &MockCredentialRepository{ctrl: ctrl}
	mock.recorder = &MockCredentialRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCredentialRepository) EXPECT() *MockCredentialRepositoryMockRecorder {
	return m.recorder
}

// Insert mocks base method.
func (m *MockCredentialRepository) Insert(arg0 context.Context, arg1 model.Credential) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockCredentialRepositoryMockRecorder) Insert(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockCredentialRepository)(nil).Insert), arg0, arg1)
}

// SelectAll mocks base method.
func (m *MockCredentialRepository) SelectAll(arg0 context.Context, arg1 *specification.CredentialSpecification) ([]model.Credential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectAll", arg0, arg1)
	ret0, _ := ret[0].([]model.Credential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectAll indicates an expected call of SelectAll.
func (mr *MockCredentialRepositoryMockRecorder) SelectAll(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectAll", reflect.TypeOf((*MockCredentialRepository)(nil).SelectAll), arg0, arg1)
}
