// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/msmkdenis/yap-infokeeper/internal/text_data/api/grpchandlers (interfaces: TextDataService)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/msmkdenis/yap-infokeeper/internal/model"
	specification "github.com/msmkdenis/yap-infokeeper/internal/text_data/specification"
)

// MockTextDataService is a mock of TextDataService interface.
type MockTextDataService struct {
	ctrl     *gomock.Controller
	recorder *MockTextDataServiceMockRecorder
}

// MockTextDataServiceMockRecorder is the mock recorder for MockTextDataService.
type MockTextDataServiceMockRecorder struct {
	mock *MockTextDataService
}

// NewMockTextDataService creates a new mock instance.
func NewMockTextDataService(ctrl *gomock.Controller) *MockTextDataService {
	mock := &MockTextDataService{ctrl: ctrl}
	mock.recorder = &MockTextDataServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTextDataService) EXPECT() *MockTextDataServiceMockRecorder {
	return m.recorder
}

// Load mocks base method.
func (m *MockTextDataService) Load(arg0 context.Context, arg1 *specification.TextDataSpecification) ([]model.TextData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Load", arg0, arg1)
	ret0, _ := ret[0].([]model.TextData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Load indicates an expected call of Load.
func (mr *MockTextDataServiceMockRecorder) Load(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Load", reflect.TypeOf((*MockTextDataService)(nil).Load), arg0, arg1)
}

// Save mocks base method.
func (m *MockTextDataService) Save(arg0 context.Context, arg1 model.TextData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockTextDataServiceMockRecorder) Save(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockTextDataService)(nil).Save), arg0, arg1)
}
