// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mock_collectionitem is a generated GoMock package.
package mock_collectionitem

import (
	collectionitem "20dojo-online/pkg/server/domain/model/collectionitem"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockCollectionItemRepo is a mock of CollectionItemRepo interface
type MockCollectionItemRepo struct {
	ctrl     *gomock.Controller
	recorder *MockCollectionItemRepoMockRecorder
}

// MockCollectionItemRepoMockRecorder is the mock recorder for MockCollectionItemRepo
type MockCollectionItemRepoMockRecorder struct {
	mock *MockCollectionItemRepo
}

// NewMockCollectionItemRepo creates a new mock instance
func NewMockCollectionItemRepo(ctrl *gomock.Controller) *MockCollectionItemRepo {
	mock := &MockCollectionItemRepo{ctrl: ctrl}
	mock.recorder = &MockCollectionItemRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCollectionItemRepo) EXPECT() *MockCollectionItemRepoMockRecorder {
	return m.recorder
}

// SelectAllCollectionItem mocks base method
func (m *MockCollectionItemRepo) SelectAllCollectionItem() ([]*collectionitem.CollectionItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectAllCollectionItem")
	ret0, _ := ret[0].([]*collectionitem.CollectionItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectAllCollectionItem indicates an expected call of SelectAllCollectionItem
func (mr *MockCollectionItemRepoMockRecorder) SelectAllCollectionItem() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectAllCollectionItem", reflect.TypeOf((*MockCollectionItemRepo)(nil).SelectAllCollectionItem))
}
