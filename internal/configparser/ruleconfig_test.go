package configparser

import (
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/require"
)

func TestButtonTempoAndModeShiftSchema(t *testing.T) {
	data := []byte(`rules:
  - type: button-tempo
    threshold_ms: 500
    input: {device: stick, button: BTN_TRIGGER}
    tap:
      outputs:
        - {device: gamepad, button: BTN_A}
        - {device: keyboard, button: KEY_F4}
      mode: nav
    hold:
      outputs:
        - {device: gamepad, button: BTN_B}
  - type: mode-shift
    input: {device: stick, button: BTN_PINKIE}
    mode: modifier
`)
	var config Config
	require.NoError(t, yaml.Unmarshal(data, &config))
	require.Len(t, config.Rules, 2)
	tempo := config.Rules[0].Config.(RuleConfigButtonTempo)
	require.Equal(t, 500, tempo.ThresholdMs)
	require.Len(t, tempo.Tap.Outputs, 2)
	require.Equal(t, "nav", tempo.Tap.Mode)
	shift := config.Rules[1].Config.(RuleConfigModeShift)
	require.Equal(t, "modifier", shift.Mode)
}
