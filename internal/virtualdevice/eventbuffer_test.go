package virtualdevice

import (
	"errors"
	"testing"

	"github.com/holoplot/go-evdev"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type EventBufferTests struct {
	suite.Suite
	device       *VirtualDeviceMock
	writeOneCall *mock.Call
}

type VirtualDeviceMock struct {
	mock.Mock
}

func (m *VirtualDeviceMock) WriteOne(event *evdev.InputEvent) error {
	args := m.Called()
	return args.Error(0)
}

func TestRunnerEventBufferTests(t *testing.T) {
	suite.Run(t, new(EventBufferTests))
}

func (t *EventBufferTests) SetupTest() {
	t.device = new(VirtualDeviceMock)
}

func (t *EventBufferTests) SetupSubTest() {
	t.device = new(VirtualDeviceMock)
	t.writeOneCall = t.device.On("WriteOne").Return(nil)
}

func (t *EventBufferTests) TearDownSubTest() {
	t.writeOneCall.Unset()
}

func (t *EventBufferTests) TestNewEventBuffer() {
	buffer := NewEventBuffer(t.device)
	t.Equal(t.device, buffer.Device)
	t.Len(buffer.events, 0)
}

func (t *EventBufferTests) TestEventBufferAddEvent() {
	buffer := NewEventBuffer(t.device)
	buffer.AddEvent(&evdev.InputEvent{})
	buffer.AddEvent(&evdev.InputEvent{})
	buffer.AddEvent(&evdev.InputEvent{})
	t.Len(buffer.events, 3)
}

func (t *EventBufferTests) TestEventBufferSendEvents() {
	t.Run("3 Events", func() {
		buffer := NewEventBuffer(t.device)
		buffer.AddEvent(&evdev.InputEvent{})
		buffer.AddEvent(&evdev.InputEvent{})
		buffer.AddEvent(&evdev.InputEvent{})
		errs := buffer.SendEvents()

		t.Len(errs, 0)
		t.device.AssertNumberOfCalls(t.T(), "WriteOne", 4)
	})

	t.Run("No Events", func() {
		buffer := NewEventBuffer(t.device)
		errs := buffer.SendEvents()

		t.Len(errs, 0)
		t.device.AssertNumberOfCalls(t.T(), "WriteOne", 0)
	})

	t.Run("Bad Event", func() {
		t.writeOneCall.Unset()
		t.writeOneCall = t.device.On("WriteOne").Return(errors.New("Fail"))

		buffer := NewEventBuffer(t.device)
		buffer.AddEvent(&evdev.InputEvent{})
		errs := buffer.SendEvents()
		t.Len(errs, 2)
	})

}
