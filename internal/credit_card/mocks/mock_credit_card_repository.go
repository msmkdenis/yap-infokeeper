// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/msmkdenis/yap-infokeeper/internal/credit_card/service (interfaces: CreditCardRepository)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/msmkdenis/yap-infokeeper/internal/model"
)

// MockCreditCardRepository is a mock of CreditCardRepository interface.
type MockCreditCardRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCreditCardRepositoryMockRecorder
}

// MockCreditCardRepositoryMockRecorder is the mock recorder for MockCreditCardRepository.
type MockCreditCardRepositoryMockRecorder struct {
	mock *MockCreditCardRepository
}

// NewMockCreditCardRepository creates a new mock instance.
func NewMockCreditCardRepository(ctrl *gomock.Controller) *MockCreditCardRepository {
	mock := &MockCreditCardRepository{ctrl: ctrl}
	mock.recorder = &MockCreditCardRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCreditCardRepository) EXPECT() *MockCreditCardRepositoryMockRecorder {
	return m.recorder
}

// Insert mocks base method.
func (m *MockCreditCardRepository) Insert(arg0 context.Context, arg1 model.CreditCard) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockCreditCardRepositoryMockRecorder) Insert(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockCreditCardRepository)(nil).Insert), arg0, arg1)
}

// SelectAllByOwnerID mocks base method.
func (m *MockCreditCardRepository) SelectAllByOwnerID(arg0 context.Context, arg1 string) ([]model.CreditCard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectAllByOwnerID", arg0, arg1)
	ret0, _ := ret[0].([]model.CreditCard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectAllByOwnerID indicates an expected call of SelectAllByOwnerID.
func (mr *MockCreditCardRepositoryMockRecorder) SelectAllByOwnerID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectAllByOwnerID", reflect.TypeOf((*MockCreditCardRepository)(nil).SelectAllByOwnerID), arg0, arg1)
}

// SelectByOwnerIDCardNumber mocks base method.
func (m *MockCreditCardRepository) SelectByOwnerIDCardNumber(arg0 context.Context, arg1, arg2 string) (*model.CreditCard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectByOwnerIDCardNumber", arg0, arg1, arg2)
	ret0, _ := ret[0].(*model.CreditCard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectByOwnerIDCardNumber indicates an expected call of SelectByOwnerIDCardNumber.
func (mr *MockCreditCardRepositoryMockRecorder) SelectByOwnerIDCardNumber(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectByOwnerIDCardNumber", reflect.TypeOf((*MockCreditCardRepository)(nil).SelectByOwnerIDCardNumber), arg0, arg1, arg2)
}
