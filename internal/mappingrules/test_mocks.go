package mappingrules

import (
	"github.com/holoplot/go-evdev"
	"github.com/stretchr/testify/mock"
)

type InputDeviceMock struct {
	mock.Mock
	calls []*mock.Call
}

func NewInputDeviceMock() *InputDeviceMock {
	m := new(InputDeviceMock)
	m.calls = make([]*mock.Call, 0, 10)
	return m
}

func (m *InputDeviceMock) AbsInfos() (map[evdev.EvCode]evdev.AbsInfo, error) {
	args := m.Called()
	return args.Get(0).(map[evdev.EvCode]evdev.AbsInfo), args.Error(1)
}

// TODO: this would make a great library class functioning as a slight improvement on testify's
// mocks for instances where we want to redefine behavior more often...
// (alternately, this is possibly an anti-pattern in Go, in which case find a cleaner way to do this and remove)
func (m *InputDeviceMock) Stub(method string) *mock.Call {
	call := m.On(method)
	m.calls = append(m.calls, call)
	return call
}

func (m *InputDeviceMock) Reset() {
	for _, call := range m.calls {
		call.Unset()
	}
	m.calls = make([]*mock.Call, 0, 10)
}
