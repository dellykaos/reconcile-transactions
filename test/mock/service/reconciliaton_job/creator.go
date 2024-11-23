// Code generated by MockGen. DO NOT EDIT.
// Source: ./service/reconciliaton_job/creator.go
//
// Generated by this command:
//
//	mockgen -source=./service/reconciliaton_job/creator.go -destination=test/mock/service/./reconciliaton_job/creator.go
//

// Package mock_reconciliatonjob is a generated GoMock package.
package mock_reconciliatonjob

import (
	context "context"
	reflect "reflect"

	entity "github.com/delly/amartha/entity"
	filestorage "github.com/delly/amartha/repository/file_storage"
	dbgen "github.com/delly/amartha/repository/postgresql"
	reconciliatonjob "github.com/delly/amartha/service/reconciliaton_job"
	gomock "go.uber.org/mock/gomock"
)

// MockCreator is a mock of Creator interface.
type MockCreator struct {
	ctrl     *gomock.Controller
	recorder *MockCreatorMockRecorder
}

// MockCreatorMockRecorder is the mock recorder for MockCreator.
type MockCreatorMockRecorder struct {
	mock *MockCreator
}

// NewMockCreator creates a new mock instance.
func NewMockCreator(ctrl *gomock.Controller) *MockCreator {
	mock := &MockCreator{ctrl: ctrl}
	mock.recorder = &MockCreatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCreator) EXPECT() *MockCreatorMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockCreator) Create(ctx context.Context, params *reconciliatonjob.CreateParams) (*entity.ReconciliationJob, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, params)
	ret0, _ := ret[0].(*entity.ReconciliationJob)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockCreatorMockRecorder) Create(ctx, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCreator)(nil).Create), ctx, params)
}

// MockCreatorRepository is a mock of CreatorRepository interface.
type MockCreatorRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCreatorRepositoryMockRecorder
}

// MockCreatorRepositoryMockRecorder is the mock recorder for MockCreatorRepository.
type MockCreatorRepositoryMockRecorder struct {
	mock *MockCreatorRepository
}

// NewMockCreatorRepository creates a new mock instance.
func NewMockCreatorRepository(ctrl *gomock.Controller) *MockCreatorRepository {
	mock := &MockCreatorRepository{ctrl: ctrl}
	mock.recorder = &MockCreatorRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCreatorRepository) EXPECT() *MockCreatorRepositoryMockRecorder {
	return m.recorder
}

// CreateReconciliationJob mocks base method.
func (m *MockCreatorRepository) CreateReconciliationJob(ctx context.Context, job dbgen.CreateReconciliationJobParams) (dbgen.ReconciliationJob, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateReconciliationJob", ctx, job)
	ret0, _ := ret[0].(dbgen.ReconciliationJob)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateReconciliationJob indicates an expected call of CreateReconciliationJob.
func (mr *MockCreatorRepositoryMockRecorder) CreateReconciliationJob(ctx, job any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateReconciliationJob", reflect.TypeOf((*MockCreatorRepository)(nil).CreateReconciliationJob), ctx, job)
}

// MockFileRepository is a mock of FileRepository interface.
type MockFileRepository struct {
	ctrl     *gomock.Controller
	recorder *MockFileRepositoryMockRecorder
}

// MockFileRepositoryMockRecorder is the mock recorder for MockFileRepository.
type MockFileRepositoryMockRecorder struct {
	mock *MockFileRepository
}

// NewMockFileRepository creates a new mock instance.
func NewMockFileRepository(ctrl *gomock.Controller) *MockFileRepository {
	mock := &MockFileRepository{ctrl: ctrl}
	mock.recorder = &MockFileRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileRepository) EXPECT() *MockFileRepositoryMockRecorder {
	return m.recorder
}

// Store mocks base method.
func (m *MockFileRepository) Store(ctx context.Context, file *filestorage.File) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Store", ctx, file)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Store indicates an expected call of Store.
func (mr *MockFileRepositoryMockRecorder) Store(ctx, file any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockFileRepository)(nil).Store), ctx, file)
}
