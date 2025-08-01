// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/0xOnah/bank/internal/service (interfaces: SessionRepository)
//
// Generated by this command:
//
//	mockgen -package mockdb -destination internal/db/mock/session.go.go github.com/0xOnah/bank/internal/service SessionRepository
//

// Package mockdb is a generated GoMock package.
package mockdb

import (
	context "context"
	reflect "reflect"

	entity "github.com/0xOnah/bank/internal/entity"
	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockSessionRepository is a mock of SessionRepository interface.
type MockSessionRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSessionRepositoryMockRecorder
	isgomock struct{}
}

// MockSessionRepositoryMockRecorder is the mock recorder for MockSessionRepository.
type MockSessionRepositoryMockRecorder struct {
	mock *MockSessionRepository
}

// NewMockSessionRepository creates a new mock instance.
func NewMockSessionRepository(ctrl *gomock.Controller) *MockSessionRepository {
	mock := &MockSessionRepository{ctrl: ctrl}
	mock.recorder = &MockSessionRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionRepository) EXPECT() *MockSessionRepositoryMockRecorder {
	return m.recorder
}

// CreateSession mocks base method.
func (m *MockSessionRepository) CreateSession(ctx context.Context, arg entity.Session) (*entity.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", ctx, arg)
	ret0, _ := ret[0].(*entity.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockSessionRepositoryMockRecorder) CreateSession(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockSessionRepository)(nil).CreateSession), ctx, arg)
}

// GetSession mocks base method.
func (m *MockSessionRepository) GetSession(ctx context.Context, id uuid.UUID) (*entity.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSession", ctx, id)
	ret0, _ := ret[0].(*entity.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSession indicates an expected call of GetSession.
func (mr *MockSessionRepositoryMockRecorder) GetSession(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSession", reflect.TypeOf((*MockSessionRepository)(nil).GetSession), ctx, id)
}
