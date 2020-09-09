// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mock_usercollectionitem is a generated GoMock package.
package mock_usercollectionitem

import (
	usercollectionitem "20dojo-online/pkg/server/domain/model/usercollectionitem"
	sql "database/sql"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockUserCollectionItemRepo is a mock of UserCollectionItemRepo interface
type MockUserCollectionItemRepo struct {
	ctrl     *gomock.Controller
	recorder *MockUserCollectionItemRepoMockRecorder
}

// MockUserCollectionItemRepoMockRecorder is the mock recorder for MockUserCollectionItemRepo
type MockUserCollectionItemRepoMockRecorder struct {
	mock *MockUserCollectionItemRepo
}

// NewMockUserCollectionItemRepo creates a new mock instance
func NewMockUserCollectionItemRepo(ctrl *gomock.Controller) *MockUserCollectionItemRepo {
	mock := &MockUserCollectionItemRepo{ctrl: ctrl}
	mock.recorder = &MockUserCollectionItemRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUserCollectionItemRepo) EXPECT() *MockUserCollectionItemRepoMockRecorder {
	return m.recorder
}

// SelectSliceByUserID mocks base method
func (m *MockUserCollectionItemRepo) SelectSliceByUserID(userID string) ([]*usercollectionitem.UserCollectionItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectSliceByUserID", userID)
	ret0, _ := ret[0].([]*usercollectionitem.UserCollectionItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectSliceByUserID indicates an expected call of SelectSliceByUserID
func (mr *MockUserCollectionItemRepoMockRecorder) SelectSliceByUserID(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectSliceByUserID", reflect.TypeOf((*MockUserCollectionItemRepo)(nil).SelectSliceByUserID), userID)
}

// BulkInsert mocks base method
func (m *MockUserCollectionItemRepo) BulkInsert(arg0 []*usercollectionitem.UserCollectionItem, arg1 *sql.Tx) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BulkInsert", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// BulkInsert indicates an expected call of BulkInsert
func (mr *MockUserCollectionItemRepoMockRecorder) BulkInsert(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BulkInsert", reflect.TypeOf((*MockUserCollectionItemRepo)(nil).BulkInsert), arg0, arg1)
}
