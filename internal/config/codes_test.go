package config

import (
	"fmt"
	"testing"

	"github.com/holoplot/go-evdev"
	"github.com/stretchr/testify/suite"
)

type EventCodeParserTests struct {
	suite.Suite
}

func TestRunnerEventCodeParserTests(t *testing.T) {
	suite.Run(t, new(EventCodeParserTests))
}

func parseCodeTestCase(t *EventCodeParserTests, in string, out evdev.EvCode, prefix string) {
	t.Run(fmt.Sprintf("%s: %s", prefix, in), func() {
		code, err := parseCode(in, prefix)
		t.Nil(err)
		t.EqualValues(out, code)
	})
}

func (t *EventCodeParserTests) TestParseCodeButton() {
	testCases := []struct {
		in  string
		out evdev.EvCode
	}{
		{"BTN_A", evdev.BTN_A},
		{"A", evdev.BTN_A},
		{"BTN_TRIGGER_HAPPY", evdev.BTN_TRIGGER_HAPPY},
		{"KEY_A", evdev.KEY_A},
		{"KEY_ESC", evdev.KEY_ESC},
	}

	for _, testCase := range testCases {
		t.Run(testCase.in, func() {
			code, err := parseCodeButton(testCase.in)
			t.Nil(err)
			t.EqualValues(code, testCase.out)
		})
	}
}

func (t *EventCodeParserTests) TestParseCode() {

	t.Run("ABS", func() {
		testCases := []struct {
			in  string
			out evdev.EvCode
		}{
			{"ABS_X", evdev.ABS_X},
			{"ABS_Y", evdev.ABS_Y},
			{"ABS_Z", evdev.ABS_Z},
			{"ABS_RX", evdev.ABS_RX},
			{"ABS_RY", evdev.ABS_RY},
			{"ABS_RZ", evdev.ABS_RZ},
			{"ABS_THROTTLE", evdev.ABS_THROTTLE},
			{"ABS_RUDDER", evdev.ABS_RUDDER},
			{"x", evdev.ABS_X},
			{"y", evdev.ABS_Y},
			{"z", evdev.ABS_Z},
			{"throttle", evdev.ABS_THROTTLE},
			{"rudder", evdev.ABS_RUDDER},
			{"0x0", evdev.ABS_X},
			{"0x1", evdev.ABS_Y},
			{"0x2", evdev.ABS_Z},
		}

		for _, testCase := range testCases {
			parseCodeTestCase(t, testCase.in, testCase.out, "ABS")
		}
	})

	t.Run("REL", func() {
		testCases := []struct {
			in  string
			out evdev.EvCode
		}{
			{"REL_X", evdev.REL_X},
			{"REL_Y", evdev.REL_Y},
			{"REL_Z", evdev.REL_Z},
			{"REL_RX", evdev.REL_RX},
			{"REL_RY", evdev.REL_RY},
			{"REL_RZ", evdev.REL_RZ},
			{"REL_WHEEL", evdev.REL_WHEEL},
			{"REL_HWHEEL", evdev.REL_HWHEEL},
			{"REL_MISC", evdev.REL_MISC},
			{"x", evdev.REL_X},
			{"y", evdev.REL_Y},
			{"wheel", evdev.REL_WHEEL},
			{"0x0", evdev.REL_X},
			{"0x1", evdev.REL_Y},
			{"0x2", evdev.REL_Z},
		}

		for _, testCase := range testCases {
			parseCodeTestCase(t, testCase.in, testCase.out, "REL")
		}
	})

	t.Run("BTN", func() {
		testCases := []struct {
			in  string
			out evdev.EvCode
		}{
			{"BTN_TRIGGER", evdev.BTN_TRIGGER},
			{"trigger", evdev.BTN_TRIGGER},
			{"0", evdev.BTN_TRIGGER},
			{"0x120", evdev.BTN_TRIGGER},
		}

		for _, testCase := range testCases {
			parseCodeTestCase(t, testCase.in, testCase.out, "BTN")
		}
	})

	t.Run("Invalid", func() {
		testCases := []struct {
			in     string
			prefix string
		}{
			{"badbutton", "BTN"},
			{"ABS_X", "BTN"},
			{"!@#$%^&*(){}-_", "BTN"},
			{"REL_X", "ABS"},
			{"ABS_W", "ABS"},
			{"0", "ABS"},
			{"0xg", "ABS"},
		}

		for _, testCase := range testCases {
			t.Run(fmt.Sprintf("%s - '%s'", testCase.prefix, testCase.in), func() {
				_, err := parseCode(testCase.in, testCase.prefix)
				t.NotNil(err)
			})
		}
	})
}
