// Code generated by MockGen. DO NOT EDIT.
// Source: ./repository/payment_repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	reflect "reflect"

	model "github.com/andikabahari/eoplatform/model"
	request "github.com/andikabahari/eoplatform/request"
	gomock "github.com/golang/mock/gomock"
)

// MockPaymentRepository is a mock of PaymentRepository interface.
type MockPaymentRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPaymentRepositoryMockRecorder
}

// MockPaymentRepositoryMockRecorder is the mock recorder for MockPaymentRepository.
type MockPaymentRepositoryMockRecorder struct {
	mock *MockPaymentRepository
}

// NewMockPaymentRepository creates a new mock instance.
func NewMockPaymentRepository(ctrl *gomock.Controller) *MockPaymentRepository {
	mock := &MockPaymentRepository{ctrl: ctrl}
	mock.recorder = &MockPaymentRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPaymentRepository) EXPECT() *MockPaymentRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockPaymentRepository) Create(payment *model.Payment) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Create", payment)
}

// Create indicates an expected call of Create.
func (mr *MockPaymentRepositoryMockRecorder) Create(payment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockPaymentRepository)(nil).Create), payment)
}

// FindOnlyByOrderID mocks base method.
func (m *MockPaymentRepository) FindOnlyByOrderID(payment *model.Payment, orderID any) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "FindOnlyByOrderID", payment, orderID)
}

// FindOnlyByOrderID indicates an expected call of FindOnlyByOrderID.
func (mr *MockPaymentRepositoryMockRecorder) FindOnlyByOrderID(payment, orderID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOnlyByOrderID", reflect.TypeOf((*MockPaymentRepository)(nil).FindOnlyByOrderID), payment, orderID)
}

// GetOnlyByOrderID mocks base method.
func (m *MockPaymentRepository) GetOnlyByOrderID(payments *[]model.Payment, orderID any) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetOnlyByOrderID", payments, orderID)
}

// GetOnlyByOrderID indicates an expected call of GetOnlyByOrderID.
func (mr *MockPaymentRepositoryMockRecorder) GetOnlyByOrderID(payments, orderID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOnlyByOrderID", reflect.TypeOf((*MockPaymentRepository)(nil).GetOnlyByOrderID), payments, orderID)
}

// Update mocks base method.
func (m *MockPaymentRepository) Update(payment *model.Payment, req *request.MidtransTransactionNotificationRequest) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Update", payment, req)
}

// Update indicates an expected call of Update.
func (mr *MockPaymentRepositoryMockRecorder) Update(payment, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockPaymentRepository)(nil).Update), payment, req)
}
