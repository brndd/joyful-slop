package mappingrules

import (
	"github.com/holoplot/go-evdev"
	"github.com/stretchr/testify/mock"
)

type InputDeviceMock struct {
	mock.Mock
}

func (m *InputDeviceMock) AbsInfos() (map[evdev.EvCode]evdev.AbsInfo, error) {
	args := m.Called()
	return args.Get(0).(map[evdev.EvCode]evdev.AbsInfo), args.Error(1)
}
