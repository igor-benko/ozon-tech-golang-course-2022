// Code generated by MockGen. DO NOT EDIT.
// Source: ./repository.go

// Package storage is a generated GoMock package.
package storage

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "gitlab.ozon.dev/igor.benko.1991/homework/internal/entity"
)

// MockPersonRepo is a mock of PersonRepo interface.
type MockPersonRepo struct {
	ctrl     *gomock.Controller
	recorder *MockPersonRepoMockRecorder
}

// MockPersonRepoMockRecorder is the mock recorder for MockPersonRepo.
type MockPersonRepoMockRecorder struct {
	mock *MockPersonRepo
}

// NewMockPersonRepo creates a new mock instance.
func NewMockPersonRepo(ctrl *gomock.Controller) *MockPersonRepo {
	mock := &MockPersonRepo{ctrl: ctrl}
	mock.recorder = &MockPersonRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPersonRepo) EXPECT() *MockPersonRepoMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockPersonRepo) Create(ctx context.Context, item entity.Person) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, item)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockPersonRepoMockRecorder) Create(ctx, item interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockPersonRepo)(nil).Create), ctx, item)
}

// Delete mocks base method.
func (m *MockPersonRepo) Delete(ctx context.Context, id uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockPersonRepoMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockPersonRepo)(nil).Delete), ctx, id)
}

// Get mocks base method.
func (m *MockPersonRepo) Get(ctx context.Context, id uint64) (*entity.Person, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(*entity.Person)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockPersonRepoMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockPersonRepo)(nil).Get), ctx, id)
}

// List mocks base method.
func (m *MockPersonRepo) List(ctx context.Context, filter entity.PersonFilter) (*entity.PersonPage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, filter)
	ret0, _ := ret[0].(*entity.PersonPage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockPersonRepoMockRecorder) List(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockPersonRepo)(nil).List), ctx, filter)
}

// Update mocks base method.
func (m *MockPersonRepo) Update(ctx context.Context, item entity.Person) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, item)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockPersonRepoMockRecorder) Update(ctx, item interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockPersonRepo)(nil).Update), ctx, item)
}

// MockVehicleRepo is a mock of VehicleRepo interface.
type MockVehicleRepo struct {
	ctrl     *gomock.Controller
	recorder *MockVehicleRepoMockRecorder
}

// MockVehicleRepoMockRecorder is the mock recorder for MockVehicleRepo.
type MockVehicleRepoMockRecorder struct {
	mock *MockVehicleRepo
}

// NewMockVehicleRepo creates a new mock instance.
func NewMockVehicleRepo(ctrl *gomock.Controller) *MockVehicleRepo {
	mock := &MockVehicleRepo{ctrl: ctrl}
	mock.recorder = &MockVehicleRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVehicleRepo) EXPECT() *MockVehicleRepoMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockVehicleRepo) Create(ctx context.Context, item entity.Vehicle) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, item)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockVehicleRepoMockRecorder) Create(ctx, item interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockVehicleRepo)(nil).Create), ctx, item)
}

// Delete mocks base method.
func (m *MockVehicleRepo) Delete(ctx context.Context, id uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockVehicleRepoMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockVehicleRepo)(nil).Delete), ctx, id)
}

// Exists mocks base method.
func (m *MockVehicleRepo) Exists(ctx context.Context, regNum string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exists", ctx, regNum)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exists indicates an expected call of Exists.
func (mr *MockVehicleRepoMockRecorder) Exists(ctx, regNum interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*MockVehicleRepo)(nil).Exists), ctx, regNum)
}

// Get mocks base method.
func (m *MockVehicleRepo) Get(ctx context.Context, id uint64) (*entity.Vehicle, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(*entity.Vehicle)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockVehicleRepoMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockVehicleRepo)(nil).Get), ctx, id)
}

// GetByPersonID mocks base method.
func (m *MockVehicleRepo) GetByPersonID(ctx context.Context, personID uint64) ([]entity.Vehicle, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByPersonID", ctx, personID)
	ret0, _ := ret[0].([]entity.Vehicle)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByPersonID indicates an expected call of GetByPersonID.
func (mr *MockVehicleRepoMockRecorder) GetByPersonID(ctx, personID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByPersonID", reflect.TypeOf((*MockVehicleRepo)(nil).GetByPersonID), ctx, personID)
}

// List mocks base method.
func (m *MockVehicleRepo) List(ctx context.Context, filter entity.VehicleFilter) (*entity.VehiclePage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, filter)
	ret0, _ := ret[0].(*entity.VehiclePage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockVehicleRepoMockRecorder) List(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockVehicleRepo)(nil).List), ctx, filter)
}

// Update mocks base method.
func (m *MockVehicleRepo) Update(ctx context.Context, item entity.Vehicle) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, item)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockVehicleRepoMockRecorder) Update(ctx, item interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockVehicleRepo)(nil).Update), ctx, item)
}