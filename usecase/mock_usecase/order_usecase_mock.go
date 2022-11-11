// Code generated by MockGen. DO NOT EDIT.
// Source: ./usecase/order_usecase.go

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	reflect "reflect"

	helper "github.com/andikabahari/eoplatform/helper"
	model "github.com/andikabahari/eoplatform/model"
	request "github.com/andikabahari/eoplatform/request"
	gomock "github.com/golang/mock/gomock"
	echo "github.com/labstack/echo/v4"
)

// MockOrderUsecase is a mock of OrderUsecase interface.
type MockOrderUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockOrderUsecaseMockRecorder
}

// MockOrderUsecaseMockRecorder is the mock recorder for MockOrderUsecase.
type MockOrderUsecaseMockRecorder struct {
	mock *MockOrderUsecase
}

// NewMockOrderUsecase creates a new mock instance.
func NewMockOrderUsecase(ctrl *gomock.Controller) *MockOrderUsecase {
	mock := &MockOrderUsecase{ctrl: ctrl}
	mock.recorder = &MockOrderUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderUsecase) EXPECT() *MockOrderUsecaseMockRecorder {
	return m.recorder
}

// AcceptOrCompleteOrder mocks base method.
func (m *MockOrderUsecase) AcceptOrCompleteOrder(ctx echo.Context, order *model.Order) helper.APIError {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AcceptOrCompleteOrder", ctx, order)
	ret0, _ := ret[0].(helper.APIError)
	return ret0
}

// AcceptOrCompleteOrder indicates an expected call of AcceptOrCompleteOrder.
func (mr *MockOrderUsecaseMockRecorder) AcceptOrCompleteOrder(ctx, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AcceptOrCompleteOrder", reflect.TypeOf((*MockOrderUsecase)(nil).AcceptOrCompleteOrder), ctx, order)
}

// CancelOrder mocks base method.
func (m *MockOrderUsecase) CancelOrder(ctx echo.Context, order *model.Order) helper.APIError {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CancelOrder", ctx, order)
	ret0, _ := ret[0].(helper.APIError)
	return ret0
}

// CancelOrder indicates an expected call of CancelOrder.
func (mr *MockOrderUsecaseMockRecorder) CancelOrder(ctx, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CancelOrder", reflect.TypeOf((*MockOrderUsecase)(nil).CancelOrder), ctx, order)
}

// CreateOrder mocks base method.
func (m *MockOrderUsecase) CreateOrder(claims *helper.JWTCustomClaims, order *model.Order, req *request.CreateOrderRequest) helper.APIError {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrder", claims, order, req)
	ret0, _ := ret[0].(helper.APIError)
	return ret0
}

// CreateOrder indicates an expected call of CreateOrder.
func (mr *MockOrderUsecaseMockRecorder) CreateOrder(claims, order, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrder", reflect.TypeOf((*MockOrderUsecase)(nil).CreateOrder), claims, order, req)
}

// GetOrders mocks base method.
func (m *MockOrderUsecase) GetOrders(claims *helper.JWTCustomClaims, orders *[]model.Order, payments *[]model.Payment) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetOrders", claims, orders, payments)
}

// GetOrders indicates an expected call of GetOrders.
func (mr *MockOrderUsecaseMockRecorder) GetOrders(claims, orders, payments interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrders", reflect.TypeOf((*MockOrderUsecase)(nil).GetOrders), claims, orders, payments)
}

// PaymentStatus mocks base method.
func (m *MockOrderUsecase) PaymentStatus(req *request.MidtransTransactionNotificationRequest) helper.APIError {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PaymentStatus", req)
	ret0, _ := ret[0].(helper.APIError)
	return ret0
}

// PaymentStatus indicates an expected call of PaymentStatus.
func (mr *MockOrderUsecaseMockRecorder) PaymentStatus(req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PaymentStatus", reflect.TypeOf((*MockOrderUsecase)(nil).PaymentStatus), req)
}