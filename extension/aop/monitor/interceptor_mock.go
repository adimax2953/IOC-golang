// Code generated by mockery v2.12.2. DO NOT EDIT.

package monitor

import (
	testing "testing"

	mock "github.com/stretchr/testify/mock"
)

// mockInterceptor is an autogenerated mock type for the interceptor type
type mockInterceptor struct {
	mock.Mock
}

// Monitor provides a mock function with given fields: monitorCtx
func (_m *mockInterceptor) Monitor(monitorCtx *context) {
	_m.Called(monitorCtx)
}

// StopMonitor provides a mock function with given fields:
func (_m *mockInterceptor) StopMonitor() {
	_m.Called()
}

// newMockInterceptor creates a new instance of mockInterceptor. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func newMockInterceptor(t testing.TB) *mockInterceptor {
	mock := &mockInterceptor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
