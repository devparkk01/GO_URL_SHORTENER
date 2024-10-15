// Code generated by MockGen. DO NOT EDIT.
// Source: URL_SHORTENER/storage (interfaces: URLOperations)

// Package storage is a generated GoMock package.
package storage

import (
	models "URL_SHORTENER/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockURLOperations is a mock of URLOperations interface.
type MockURLOperations struct {
	ctrl     *gomock.Controller
	recorder *MockURLOperationsMockRecorder
}

// MockURLOperationsMockRecorder is the mock recorder for MockURLOperations.
type MockURLOperationsMockRecorder struct {
	mock *MockURLOperations
}

// NewMockURLOperations creates a new mock instance.
func NewMockURLOperations(ctrl *gomock.Controller) *MockURLOperations {
	mock := &MockURLOperations{ctrl: ctrl}
	mock.recorder = &MockURLOperationsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockURLOperations) EXPECT() *MockURLOperationsMockRecorder {
	return m.recorder
}

// CheckOriginalUrlExists mocks base method.
func (m *MockURLOperations) CheckOriginalUrlExists(arg0 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckOriginalUrlExists", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// CheckOriginalUrlExists indicates an expected call of CheckOriginalUrlExists.
func (mr *MockURLOperationsMockRecorder) CheckOriginalUrlExists(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckOriginalUrlExists", reflect.TypeOf((*MockURLOperations)(nil).CheckOriginalUrlExists), arg0)
}

// CheckShortUrlExists mocks base method.
func (m *MockURLOperations) CheckShortUrlExists(arg0 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckShortUrlExists", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// CheckShortUrlExists indicates an expected call of CheckShortUrlExists.
func (mr *MockURLOperationsMockRecorder) CheckShortUrlExists(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckShortUrlExists", reflect.TypeOf((*MockURLOperations)(nil).CheckShortUrlExists), arg0)
}

// DeleteShortUrl mocks base method.
func (m *MockURLOperations) DeleteShortUrl(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteShortUrl", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteShortUrl indicates an expected call of DeleteShortUrl.
func (mr *MockURLOperationsMockRecorder) DeleteShortUrl(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteShortUrl", reflect.TypeOf((*MockURLOperations)(nil).DeleteShortUrl), arg0)
}

// GetOriginalUrl mocks base method.
func (m *MockURLOperations) GetOriginalUrl(arg0 string) (*models.Url, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOriginalUrl", arg0)
	ret0, _ := ret[0].(*models.Url)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOriginalUrl indicates an expected call of GetOriginalUrl.
func (mr *MockURLOperationsMockRecorder) GetOriginalUrl(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOriginalUrl", reflect.TypeOf((*MockURLOperations)(nil).GetOriginalUrl), arg0)
}

// InsertUrl mocks base method.
func (m *MockURLOperations) InsertUrl(arg0 *models.Url) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertUrl", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertUrl indicates an expected call of InsertUrl.
func (mr *MockURLOperationsMockRecorder) InsertUrl(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertUrl", reflect.TypeOf((*MockURLOperations)(nil).InsertUrl), arg0)
}

// UpdateShortUrl mocks base method.
func (m *MockURLOperations) UpdateShortUrl(arg0, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateShortUrl", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateShortUrl indicates an expected call of UpdateShortUrl.
func (mr *MockURLOperationsMockRecorder) UpdateShortUrl(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateShortUrl", reflect.TypeOf((*MockURLOperations)(nil).UpdateShortUrl), arg0, arg1, arg2)
}