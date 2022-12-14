// Code generated by MockGen. DO NOT EDIT.
// Source: internal/converter/converter.go

// Package mock_converter is a generated GoMock package.
package mock_converter

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	converter "gitlab.ozon.dev/alex1234562557/telegram-bot/internal/converter"
	domain "gitlab.ozon.dev/alex1234562557/telegram-bot/internal/domain"
)

// MockRateUpdater is a mock of RateUpdater interface.
type MockRateUpdater struct {
	ctrl     *gomock.Controller
	recorder *MockRateUpdaterMockRecorder
}

// MockRateUpdaterMockRecorder is the mock recorder for MockRateUpdater.
type MockRateUpdaterMockRecorder struct {
	mock *MockRateUpdater
}

// NewMockRateUpdater creates a new mock instance.
func NewMockRateUpdater(ctrl *gomock.Controller) *MockRateUpdater {
	mock := &MockRateUpdater{ctrl: ctrl}
	mock.recorder = &MockRateUpdaterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRateUpdater) EXPECT() *MockRateUpdaterMockRecorder {
	return m.recorder
}

// GetUpdate mocks base method.
func (m *MockRateUpdater) GetUpdate(ctx context.Context, date *int64) (*converter.Rates, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUpdate", ctx, date)
	ret0, _ := ret[0].(*converter.Rates)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUpdate indicates an expected call of GetUpdate.
func (mr *MockRateUpdaterMockRecorder) GetUpdate(ctx, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUpdate", reflect.TypeOf((*MockRateUpdater)(nil).GetUpdate), ctx, date)
}

// MockRateManipulator is a mock of RateManipulator interface.
type MockRateManipulator struct {
	ctrl     *gomock.Controller
	recorder *MockRateManipulatorMockRecorder
}

// MockRateManipulatorMockRecorder is the mock recorder for MockRateManipulator.
type MockRateManipulatorMockRecorder struct {
	mock *MockRateManipulator
}

// NewMockRateManipulator creates a new mock instance.
func NewMockRateManipulator(ctrl *gomock.Controller) *MockRateManipulator {
	mock := &MockRateManipulator{ctrl: ctrl}
	mock.recorder = &MockRateManipulatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRateManipulator) EXPECT() *MockRateManipulatorMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockRateManipulator) Add(ctx context.Context, date int64, code string, nominal float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", ctx, date, code, nominal)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockRateManipulatorMockRecorder) Add(ctx, date, code, nominal interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockRateManipulator)(nil).Add), ctx, date, code, nominal)
}

// Get mocks base method.
func (m *MockRateManipulator) Get(ctx context.Context, date int64, code string) (*domain.Rate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, date, code)
	ret0, _ := ret[0].(*domain.Rate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRateManipulatorMockRecorder) Get(ctx, date, code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRateManipulator)(nil).Get), ctx, date, code)
}

// MockUserManipulator is a mock of UserManipulator interface.
type MockUserManipulator struct {
	ctrl     *gomock.Controller
	recorder *MockUserManipulatorMockRecorder
}

// MockUserManipulatorMockRecorder is the mock recorder for MockUserManipulator.
type MockUserManipulatorMockRecorder struct {
	mock *MockUserManipulator
}

// NewMockUserManipulator creates a new mock instance.
func NewMockUserManipulator(ctrl *gomock.Controller) *MockUserManipulator {
	mock := &MockUserManipulator{ctrl: ctrl}
	mock.recorder = &MockUserManipulatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserManipulator) EXPECT() *MockUserManipulatorMockRecorder {
	return m.recorder
}

// GetAllUsers mocks base method.
func (m *MockUserManipulator) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllUsers", ctx)
	ret0, _ := ret[0].([]domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllUsers indicates an expected call of GetAllUsers.
func (mr *MockUserManipulatorMockRecorder) GetAllUsers(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllUsers", reflect.TypeOf((*MockUserManipulator)(nil).GetAllUsers), ctx)
}

// UpdateBudget mocks base method.
func (m *MockUserManipulator) UpdateBudget(ctx context.Context, userID int64, budget float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBudget", ctx, userID, budget)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBudget indicates an expected call of UpdateBudget.
func (mr *MockUserManipulatorMockRecorder) UpdateBudget(ctx, userID, budget interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBudget", reflect.TypeOf((*MockUserManipulator)(nil).UpdateBudget), ctx, userID, budget)
}

// MockReportCacher is a mock of ReportCacher interface.
type MockReportCacher struct {
	ctrl     *gomock.Controller
	recorder *MockReportCacherMockRecorder
}

// MockReportCacherMockRecorder is the mock recorder for MockReportCacher.
type MockReportCacherMockRecorder struct {
	mock *MockReportCacher
}

// NewMockReportCacher creates a new mock instance.
func NewMockReportCacher(ctrl *gomock.Controller) *MockReportCacher {
	mock := &MockReportCacher{ctrl: ctrl}
	mock.recorder = &MockReportCacherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReportCacher) EXPECT() *MockReportCacherMockRecorder {
	return m.recorder
}

// RemoveFromAll mocks base method.
func (m *MockReportCacher) RemoveFromAll(ctx context.Context, key []int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveFromAll", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveFromAll indicates an expected call of RemoveFromAll.
func (mr *MockReportCacherMockRecorder) RemoveFromAll(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveFromAll", reflect.TypeOf((*MockReportCacher)(nil).RemoveFromAll), ctx, key)
}
