// Code generated by MockGen. DO NOT EDIT.
// Source: ./cache.go

// Package cache is a generated GoMock package.
package cache

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCacheClient is a mock of CacheClient interface.
type MockCacheClient struct {
	ctrl     *gomock.Controller
	recorder *MockCacheClientMockRecorder
}

// MockCacheClientMockRecorder is the mock recorder for MockCacheClient.
type MockCacheClientMockRecorder struct {
	mock *MockCacheClient
}

// NewMockCacheClient creates a new mock instance.
func NewMockCacheClient(ctrl *gomock.Controller) *MockCacheClient {
	mock := &MockCacheClient{ctrl: ctrl}
	mock.recorder = &MockCacheClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCacheClient) EXPECT() *MockCacheClientMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockCacheClient) Get(ctx context.Context, key string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, key)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockCacheClientMockRecorder) Get(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCacheClient)(nil).Get), ctx, key)
}

// Publish mocks base method.
func (m *MockCacheClient) Publish(ctx context.Context, pubData *Publication) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", ctx, pubData)
	ret0, _ := ret[0].(error)
	return ret0
}

// Publish indicates an expected call of Publish.
func (mr *MockCacheClientMockRecorder) Publish(ctx, pubData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockCacheClient)(nil).Publish), ctx, pubData)
}

// Set mocks base method.
func (m *MockCacheClient) Set(ctx context.Context, key string, value []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, key, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockCacheClientMockRecorder) Set(ctx, key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockCacheClient)(nil).Set), ctx, key, value)
}

// Subscribe mocks base method.
func (m *MockCacheClient) Subscribe(ctx context.Context) <-chan *Publication {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", ctx)
	ret0, _ := ret[0].(<-chan *Publication)
	return ret0
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockCacheClientMockRecorder) Subscribe(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockCacheClient)(nil).Subscribe), ctx)
}