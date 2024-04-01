// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/msmkdenis/yap-infokeeper/internal/credential/api/grpchandlers (interfaces: CredentialService)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/msmkdenis/yap-infokeeper/internal/model"
)

// MockCredentialService is a mock of CredentialService interface.
type MockCredentialService struct {
	ctrl     *gomock.Controller
	recorder *MockCredentialServiceMockRecorder
}

// MockCredentialServiceMockRecorder is the mock recorder for MockCredentialService.
type MockCredentialServiceMockRecorder struct {
	mock *MockCredentialService
}

// NewMockCredentialService creates a new mock instance.
func NewMockCredentialService(ctrl *gomock.Controller) *MockCredentialService {
	mock := &MockCredentialService{ctrl: ctrl}
	mock.recorder = &MockCredentialServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCredentialService) EXPECT() *MockCredentialServiceMockRecorder {
	return m.recorder
}

// Save mocks base method.
func (m *MockCredentialService) Save(arg0 context.Context, arg1 model.Credential) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockCredentialServiceMockRecorder) Save(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockCredentialService)(nil).Save), arg0, arg1)
}
