// Code generated by MockGen. DO NOT EDIT.
// Source: ./usecase/register_usecase.go

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	reflect "reflect"

	helper "github.com/andikabahari/eoplatform/helper"
	model "github.com/andikabahari/eoplatform/model"
	request "github.com/andikabahari/eoplatform/request"
	gomock "github.com/golang/mock/gomock"
)

// MockRegisterUsecase is a mock of RegisterUsecase interface.
type MockRegisterUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockRegisterUsecaseMockRecorder
}

// MockRegisterUsecaseMockRecorder is the mock recorder for MockRegisterUsecase.
type MockRegisterUsecaseMockRecorder struct {
	mock *MockRegisterUsecase
}

// NewMockRegisterUsecase creates a new mock instance.
func NewMockRegisterUsecase(ctrl *gomock.Controller) *MockRegisterUsecase {
	mock := &MockRegisterUsecase{ctrl: ctrl}
	mock.recorder = &MockRegisterUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRegisterUsecase) EXPECT() *MockRegisterUsecaseMockRecorder {
	return m.recorder
}

// Register mocks base method.
func (m *MockRegisterUsecase) Register(user *model.User, req *request.CreateUserRequest) helper.APIError {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", user, req)
	ret0, _ := ret[0].(helper.APIError)
	return ret0
}

// Register indicates an expected call of Register.
func (mr *MockRegisterUsecaseMockRecorder) Register(user, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockRegisterUsecase)(nil).Register), user, req)
}
