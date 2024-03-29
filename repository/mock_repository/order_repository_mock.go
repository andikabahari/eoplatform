// Code generated by MockGen. DO NOT EDIT.
// Source: ./repository/order_repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	reflect "reflect"

	model "github.com/andikabahari/eoplatform/model"
	gomock "github.com/golang/mock/gomock"
)

// MockOrderRepository is a mock of OrderRepository interface.
type MockOrderRepository struct {
	ctrl     *gomock.Controller
	recorder *MockOrderRepositoryMockRecorder
}

// MockOrderRepositoryMockRecorder is the mock recorder for MockOrderRepository.
type MockOrderRepositoryMockRecorder struct {
	mock *MockOrderRepository
}

// NewMockOrderRepository creates a new mock instance.
func NewMockOrderRepository(ctrl *gomock.Controller) *MockOrderRepository {
	mock := &MockOrderRepository{ctrl: ctrl}
	mock.recorder = &MockOrderRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderRepository) EXPECT() *MockOrderRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockOrderRepository) Create(order *model.Order) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Create", order)
}

// Create indicates an expected call of Create.
func (mr *MockOrderRepositoryMockRecorder) Create(order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockOrderRepository)(nil).Create), order)
}

// Delete mocks base method.
func (m *MockOrderRepository) Delete(order *model.Order) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Delete", order)
}

// Delete indicates an expected call of Delete.
func (mr *MockOrderRepositoryMockRecorder) Delete(order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockOrderRepository)(nil).Delete), order)
}

// Find mocks base method.
func (m *MockOrderRepository) Find(order *model.Order, id string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Find", order, id)
}

// Find indicates an expected call of Find.
func (mr *MockOrderRepositoryMockRecorder) Find(order, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockOrderRepository)(nil).Find), order, id)
}

// FindOnly mocks base method.
func (m *MockOrderRepository) FindOnly(order *model.Order, id any) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "FindOnly", order, id)
}

// FindOnly indicates an expected call of FindOnly.
func (mr *MockOrderRepositoryMockRecorder) FindOnly(order, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOnly", reflect.TypeOf((*MockOrderRepository)(nil).FindOnly), order, id)
}

// GetOrdersForCustomer mocks base method.
func (m *MockOrderRepository) GetOrdersForCustomer(orders *[]model.Order, userID uint) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetOrdersForCustomer", orders, userID)
}

// GetOrdersForCustomer indicates an expected call of GetOrdersForCustomer.
func (mr *MockOrderRepositoryMockRecorder) GetOrdersForCustomer(orders, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrdersForCustomer", reflect.TypeOf((*MockOrderRepository)(nil).GetOrdersForCustomer), orders, userID)
}

// GetOrdersForOrganizer mocks base method.
func (m *MockOrderRepository) GetOrdersForOrganizer(orders *[]model.Order, userID uint) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetOrdersForOrganizer", orders, userID)
}

// GetOrdersForOrganizer indicates an expected call of GetOrdersForOrganizer.
func (mr *MockOrderRepositoryMockRecorder) GetOrdersForOrganizer(orders, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrdersForOrganizer", reflect.TypeOf((*MockOrderRepository)(nil).GetOrdersForOrganizer), orders, userID)
}

// Save mocks base method.
func (m *MockOrderRepository) Save(order *model.Order) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Save", order)
}

// Save indicates an expected call of Save.
func (mr *MockOrderRepositoryMockRecorder) Save(order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockOrderRepository)(nil).Save), order)
}
