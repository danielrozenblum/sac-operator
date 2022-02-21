// Code generated by mockery 2.9.0. DO NOT EDIT.

package connector_deployer

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockConnectorDeployer is an autogenerated mock type for the ConnectorDeployer type
type MockConnectorDeployer struct {
	mock.Mock
}

// CreateConnector provides a mock function with given fields: ctx, inputs
func (_m *MockConnectorDeployer) CreateConnector(ctx context.Context, inputs *CreateConnectorInput) (string, error) {
	ret := _m.Called(ctx, inputs)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, *CreateConnectorInput) string); ok {
		r0 = rf(ctx, inputs)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *CreateConnectorInput) error); ok {
		r1 = rf(ctx, inputs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteConnector provides a mock function with given fields: ctx, Name
func (_m *MockConnectorDeployer) DeleteConnector(ctx context.Context, Name string) error {
	ret := _m.Called(ctx, Name)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, Name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetConnectorsForSite provides a mock function with given fields: ctx, siteName
func (_m *MockConnectorDeployer) GetConnectorsForSite(ctx context.Context, siteName string) ([]Connector, error) {
	ret := _m.Called(ctx, siteName)

	var r0 []Connector
	if rf, ok := ret.Get(0).(func(context.Context, string) []Connector); ok {
		r0 = rf(ctx, siteName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]Connector)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, siteName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}