package mock_service

import (
	reflect "reflect"

	domain "github.com/Aleksey-Andris/go-yandex-shortener/internal/app/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockLinkStorage is a mock of LinkStorage interface.
type MockLinkStorage struct {
	ctrl     *gomock.Controller
	recorder *MockLinkStorageMockRecorder
}

// MockLinkStorageMockRecorder is the mock recorder for MockLinkStorage.
type MockLinkStorageMockRecorder struct {
	mock *MockLinkStorage
}

// NewMockLinkStorage creates a new mock instance.
func NewMockLinkStorage(ctrl *gomock.Controller) *MockLinkStorage {
	mock := &MockLinkStorage{ctrl: ctrl}
	mock.recorder = &MockLinkStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLinkStorage) EXPECT() *MockLinkStorageMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockLinkStorage) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockLinkStorageMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockLinkStorage)(nil).Close))
}

// Create mocks base method.
func (m *MockLinkStorage) Create(idemt, fulLink string) (domain.Link, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", idemt, fulLink)
	ret0, _ := ret[0].(domain.Link)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockLinkStorageMockRecorder) Create(idemt, fulLink interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockLinkStorage)(nil).Create), idemt, fulLink)
}

// GetMaxID mocks base method.
func (m *MockLinkStorage) GetMaxID() (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMaxID")
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMaxID indicates an expected call of GetMaxID.
func (mr *MockLinkStorageMockRecorder) GetMaxID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMaxID", reflect.TypeOf((*MockLinkStorage)(nil).GetMaxID))
}

// GetOneByIdent mocks base method.
func (m *MockLinkStorage) GetOneByIdent(ident string) (domain.Link, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOneByIdent", ident)
	ret0, _ := ret[0].(domain.Link)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOneByIdent indicates an expected call of GetOneByIdent.
func (mr *MockLinkStorageMockRecorder) GetOneByIdent(ident interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOneByIdent", reflect.TypeOf((*MockLinkStorage)(nil).GetOneByIdent), ident)
}
