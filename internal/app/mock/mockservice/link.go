package mockservice

import (
	context "context"
	reflect "reflect"

	domain "github.com/Aleksey-Andris/go-yandex-shortener/internal/app/domain"
	dto "github.com/Aleksey-Andris/go-yandex-shortener/internal/app/dto"
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
func (m *MockLinkStorage) Create(ctx context.Context, idemt, fulLink string, userID int32) (domain.Link, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, idemt, fulLink, userID)
	ret0, _ := ret[0].(domain.Link)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockLinkStorageMockRecorder) Create(ctx, idemt, fulLink, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockLinkStorage)(nil).Create), ctx, idemt, fulLink, userID)
}

// CreateLinks mocks base method.
func (m *MockLinkStorage) CreateLinks(ctx context.Context, links []domain.Link, userID int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateLinks", ctx, links, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateLinks indicates an expected call of CreateLinks.
func (mr *MockLinkStorageMockRecorder) CreateLinks(ctx, links, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateLinks", reflect.TypeOf((*MockLinkStorage)(nil).CreateLinks), ctx, links, userID)
}

// DeleteByIdents mocks base method.
func (m *MockLinkStorage) DeleteByIdents(ctx context.Context, idents ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range idents {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteByIdents", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByIdents indicates an expected call of DeleteByIdents.
func (mr *MockLinkStorageMockRecorder) DeleteByIdents(ctx interface{}, idents ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, idents...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByIdents", reflect.TypeOf((*MockLinkStorage)(nil).DeleteByIdents), varargs...)
}

// GetByIdents mocks base method.
func (m *MockLinkStorage) GetByIdents(ctx context.Context, idents ...string) ([]domain.Link, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range idents {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetByIdents", varargs...)
	ret0, _ := ret[0].([]domain.Link)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByIdents indicates an expected call of GetByIdents.
func (mr *MockLinkStorageMockRecorder) GetByIdents(ctx interface{}, idents ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, idents...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByIdents", reflect.TypeOf((*MockLinkStorage)(nil).GetByIdents), varargs...)
}

// GetLinksByUserID mocks base method.
func (m *MockLinkStorage) GetLinksByUserID(ctx context.Context, userID int32) ([]dto.LinkListByUserIDRes, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLinksByUserID", ctx, userID)
	ret0, _ := ret[0].([]dto.LinkListByUserIDRes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLinksByUserID indicates an expected call of GetLinksByUserID.
func (mr *MockLinkStorageMockRecorder) GetLinksByUserID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLinksByUserID", reflect.TypeOf((*MockLinkStorage)(nil).GetLinksByUserID), ctx, userID)
}

// GetOneByIdent mocks base method.
func (m *MockLinkStorage) GetOneByIdent(ctx context.Context, ident string) (domain.Link, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOneByIdent", ctx, ident)
	ret0, _ := ret[0].(domain.Link)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOneByIdent indicates an expected call of GetOneByIdent.
func (mr *MockLinkStorageMockRecorder) GetOneByIdent(ctx, ident interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOneByIdent", reflect.TypeOf((*MockLinkStorage)(nil).GetOneByIdent), ctx, ident)
}
