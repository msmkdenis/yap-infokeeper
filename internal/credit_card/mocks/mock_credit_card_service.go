// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/msmkdenis/yap-infokeeper/internal/credit_card/api/grpchandlers (interfaces: CreditCardService)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	specification "github.com/msmkdenis/yap-infokeeper/internal/credit_card/specification"
	model "github.com/msmkdenis/yap-infokeeper/internal/model"
)

// MockCreditCardService is a mock of CreditCardService interface.
type MockCreditCardService struct {
	ctrl     *gomock.Controller
	recorder *MockCreditCardServiceMockRecorder
}

// MockCreditCardServiceMockRecorder is the mock recorder for MockCreditCardService.
type MockCreditCardServiceMockRecorder struct {
	mock *MockCreditCardService
}

// NewMockCreditCardService creates a new mock instance.
func NewMockCreditCardService(ctrl *gomock.Controller) *MockCreditCardService {
	mock := &MockCreditCardService{ctrl: ctrl}
	mock.recorder = &MockCreditCardServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCreditCardService) EXPECT() *MockCreditCardServiceMockRecorder {
	return m.recorder
}

// Load mocks base method.
func (m *MockCreditCardService) Load(arg0 context.Context, arg1 *specification.CreditCardSpecification) ([]model.CreditCard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Load", arg0, arg1)
	ret0, _ := ret[0].([]model.CreditCard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Load indicates an expected call of Load.
func (mr *MockCreditCardServiceMockRecorder) Load(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Load", reflect.TypeOf((*MockCreditCardService)(nil).Load), arg0, arg1)
}

// Save mocks base method.
func (m *MockCreditCardService) Save(arg0 context.Context, arg1 model.CreditCard) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockCreditCardServiceMockRecorder) Save(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockCreditCardService)(nil).Save), arg0, arg1)
}
