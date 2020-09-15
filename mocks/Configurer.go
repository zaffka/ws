// Code generated by mockery v2.2.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Configurer is an autogenerated mock type for the Configurer type
type Configurer struct {
	mock.Mock
}

// Header provides a mock function with given fields:
func (_m *Configurer) Header() map[string][]string {
	ret := _m.Called()

	var r0 map[string][]string
	if rf, ok := ret.Get(0).(func() map[string][]string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string][]string)
		}
	}

	return r0
}

// URL provides a mock function with given fields:
func (_m *Configurer) URL() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
