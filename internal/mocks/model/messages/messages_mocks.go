// Code generated by MockGen. DO NOT EDIT.
// Source: internal/model/messages/incoming_msg.go

// Package mock_messages is a generated GoMock package.
package mock_messages

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

type Expense struct {
	Amount   float64
	Category string
	Date     int64
}

// MockMessageSender is a mock of MessageSender interface.
type MockMessageSender struct {
	ctrl     *gomock.Controller
	recorder *MockMessageSenderMockRecorder
}

// MockMessageSenderMockRecorder is the mock recorder for MockMessageSender.
type MockMessageSenderMockRecorder struct {
	mock *MockMessageSender
}

// NewMockMessageSender creates a new mock instance.
func NewMockMessageSender(ctrl *gomock.Controller) *MockMessageSender {
	mock := &MockMessageSender{ctrl: ctrl}
	mock.recorder = &MockMessageSenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMessageSender) EXPECT() *MockMessageSenderMockRecorder {
	return m.recorder
}

// SendMessage mocks base method.
func (m *MockMessageSender) SendMessage(text string, userID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", text, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMessage indicates an expected call of SendMessage.
func (mr *MockMessageSenderMockRecorder) SendMessage(text, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockMessageSender)(nil).SendMessage), text, userID)
}

// SendMessageWithKeyboard mocks base method.
func (m *MockMessageSender) SendMessageWithKeyboard(text, keyboardMarkup string, userID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessageWithKeyboard", text, keyboardMarkup, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMessageWithKeyboard indicates an expected call of SendMessageWithKeyboard.
func (mr *MockMessageSenderMockRecorder) SendMessageWithKeyboard(text, keyboardMarkup, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessageWithKeyboard", reflect.TypeOf((*MockMessageSender)(nil).SendMessageWithKeyboard), text, keyboardMarkup, userID)
}

// MockDataManipulator is a mock of DataManipulator interface.
type MockDataManipulator struct {
	ctrl     *gomock.Controller
	recorder *MockDataManipulatorMockRecorder
}

// MockDataManipulatorMockRecorder is the mock recorder for MockDataManipulator.
type MockDataManipulatorMockRecorder struct {
	mock *MockDataManipulator
}

// NewMockDataManipulator creates a new mock instance.
func NewMockDataManipulator(ctrl *gomock.Controller) *MockDataManipulator {
	mock := &MockDataManipulator{ctrl: ctrl}
	mock.recorder = &MockDataManipulatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDataManipulator) EXPECT() *MockDataManipulatorMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockDataManipulator) Add(ctx context.Context, userID int64, expense *Expense) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", ctx, userID, expense)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockDataManipulatorMockRecorder) Add(ctx, userID, expense interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockDataManipulator)(nil).Add), ctx, userID, expense)
}

// Get mocks base method.
func (m *MockDataManipulator) Get(ctx context.Context, userID int64) ([]*Expense, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, userID)
	ret0, _ := ret[0].([]*Expense)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockDataManipulatorMockRecorder) Get(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockDataManipulator)(nil).Get), ctx, userID)
}

// MockStateManipulator is a mock of StateManipulator interface.
type MockStateManipulator struct {
	ctrl     *gomock.Controller
	recorder *MockStateManipulatorMockRecorder
}

// MockStateManipulatorMockRecorder is the mock recorder for MockStateManipulator.
type MockStateManipulatorMockRecorder struct {
	mock *MockStateManipulator
}

// NewMockStateManipulator creates a new mock instance.
func NewMockStateManipulator(ctrl *gomock.Controller) *MockStateManipulator {
	mock := &MockStateManipulator{ctrl: ctrl}
	mock.recorder = &MockStateManipulatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStateManipulator) EXPECT() *MockStateManipulatorMockRecorder {
	return m.recorder
}

// GetState mocks base method.
func (m *MockStateManipulator) GetState(ctx context.Context, userID int64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetState", ctx, userID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetState indicates an expected call of GetState.
func (mr *MockStateManipulatorMockRecorder) GetState(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetState", reflect.TypeOf((*MockStateManipulator)(nil).GetState), ctx, userID)
}

// MockStorageManipulator is a mock of StorageManipulator interface.
type MockStorageManipulator struct {
	ctrl     *gomock.Controller
	recorder *MockStorageManipulatorMockRecorder
}

// MockStorageManipulatorMockRecorder is the mock recorder for MockStorageManipulator.
type MockStorageManipulatorMockRecorder struct {
	mock *MockStorageManipulator
}

// NewMockStorageManipulator creates a new mock instance.
func NewMockStorageManipulator(ctrl *gomock.Controller) *MockStorageManipulator {
	mock := &MockStorageManipulator{ctrl: ctrl}
	mock.recorder = &MockStorageManipulatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorageManipulator) EXPECT() *MockStorageManipulatorMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockStorageManipulator) Add(ctx context.Context, userID int64, expense *Expense) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", ctx, userID, expense)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockStorageManipulatorMockRecorder) Add(ctx, userID, expense interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockStorageManipulator)(nil).Add), ctx, userID, expense)
}

// Get mocks base method.
func (m *MockStorageManipulator) Get(ctx context.Context, userID int64) ([]*Expense, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, userID)
	ret0, _ := ret[0].([]*Expense)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockStorageManipulatorMockRecorder) Get(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockStorageManipulator)(nil).Get), ctx, userID)
}

// GetState mocks base method.
func (m *MockStorageManipulator) GetState(ctx context.Context, userID int64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetState", ctx, userID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetState indicates an expected call of GetState.
func (mr *MockStorageManipulatorMockRecorder) GetState(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetState", reflect.TypeOf((*MockStorageManipulator)(nil).GetState), ctx, userID)
}

// MockConverter is a mock of Converter interface.
type MockConverter struct {
	ctrl     *gomock.Controller
	recorder *MockConverterMockRecorder
}

// MockConverterMockRecorder is the mock recorder for MockConverter.
type MockConverterMockRecorder struct {
	mock *MockConverter
}

// NewMockConverter creates a new mock instance.
func NewMockConverter(ctrl *gomock.Controller) *MockConverter {
	mock := &MockConverter{ctrl: ctrl}
	mock.recorder = &MockConverterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConverter) EXPECT() *MockConverterMockRecorder {
	return m.recorder
}

// Exchange mocks base method.
func (m *MockConverter) Exchange(value float64, from, to string) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exchange", value, from, to)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exchange indicates an expected call of Exchange.
func (mr *MockConverterMockRecorder) Exchange(value, from, to interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exchange", reflect.TypeOf((*MockConverter)(nil).Exchange), value, from, to)
}