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
	device *VirtualDeviceMock
	buffer *EventBuffer
}

// Mocks
type VirtualDeviceMock struct {
	mock.Mock
}

func (m *VirtualDeviceMock) WriteOne(event *evdev.InputEvent) error {
	args := m.Called()
	return args.Error(0)
}

// Setup
func TestRunnerEventBufferTests(t *testing.T) {
	suite.Run(t, new(EventBufferTests))
}

func (t *EventBufferTests) SetupSubTest() {
	t.device = new(VirtualDeviceMock)
	t.buffer = &EventBuffer{Device: t.device}
}

// Tests
func (t *EventBufferTests) TestNewEventBuffer() {
	t.Equal(t.device, t.buffer.Device)
	t.Len(t.buffer.events, 0)
}

func (t *EventBufferTests) TestEventBuffer() {

	t.Run("AddEvent", func() {
		t.buffer.AddEvent(&evdev.InputEvent{})
		t.buffer.AddEvent(&evdev.InputEvent{})
		t.buffer.AddEvent(&evdev.InputEvent{})
		t.Len(t.buffer.events, 3)
	})

	t.Run("SendEvents", func() {
		t.Run("3 Events", func() {
			writeOneCall := t.device.On("WriteOne").Return(nil)

			t.buffer.AddEvent(&evdev.InputEvent{})
			t.buffer.AddEvent(&evdev.InputEvent{})
			t.buffer.AddEvent(&evdev.InputEvent{})
			errs := t.buffer.SendEvents()

			t.Len(errs, 0)
			t.device.AssertNumberOfCalls(t.T(), "WriteOne", 4)

			writeOneCall.Unset()
		})

		t.Run("No Events", func() {
			writeOneCall := t.device.On("WriteOne").Return(nil)

			errs := t.buffer.SendEvents()

			t.Len(errs, 0)
			t.device.AssertNumberOfCalls(t.T(), "WriteOne", 0)

			writeOneCall.Unset()
		})

		t.Run("Bad Event", func() {
			writeOneCall := t.device.On("WriteOne").Return(errors.New("Fail"))

			t.buffer.AddEvent(&evdev.InputEvent{})
			errs := t.buffer.SendEvents()
			t.Len(errs, 2)

			writeOneCall.Unset()
		})
	})
}
