package gremlin

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

var outputNames = []string{"devices.yml", "modes.yml", "axes.yml", "buttons.yml", "README.md", "conversion-report.txt"}

func Generate(inputPath, outputDir string, opts Options) (*Result, error) {
	f, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("open input: %w", err)
	}
	defer f.Close()
	p, err := Parse(f)
	if err != nil {
		return nil, err
	}
	result, err := Convert(p, opts)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, fmt.Errorf("create output directory: %w", err)
	}
	for _, name := range outputNames {
		if err := os.WriteFile(filepath.Join(outputDir, name), result.Files[name], 0o644); err != nil {
			return nil, fmt.Errorf("write %s: %w", name, err)
		}
	}
	return result, nil
}

func renderFiles(opts Options, modes []string, rules []rule, warnings []string, counts map[string]int) map[string][]byte {
	files := map[string][]byte{}
	files["devices.yml"] = renderDevices(opts)
	files["modes.yml"] = renderModes(modes)
	var axes, buttons []rule
	for _, r := range rules {
		if r.Type == "axis" || r.Type == "axis-to-relaxis" {
			axes = append(axes, r)
		} else {
			buttons = append(buttons, r)
		}
	}
	files["axes.yml"] = renderRules(axes)
	files["buttons.yml"] = renderRules(buttons)
	files["README.md"] = renderReadme(modes, counts)
	files["conversion-report.txt"] = renderReport(counts, warnings)
	return files
}

func renderDevices(opts Options) []byte {
	return []byte(fmt.Sprintf(`devices:
  - name: left
    type: physical
    device_path: %s
  - name: right
    type: physical
    device_path: %s
  - name: left-vjoy
    type: virtual
    preset: joystick
    vendor_id: "0x4711"
    device_id: "0x0817"
  - name: right-vjoy
    type: virtual
    preset: joystick
    vendor_id: "0x4711"
    device_id: "0x0818"
`, yamlString(opts.LeftPath), yamlString(opts.RightPath)))
}

func renderModes(modes []string) []byte {
	var b bytes.Buffer
	b.WriteString("modes:\n")
	for _, mode := range modes {
		fmt.Fprintf(&b, "  - %s\n", yamlString(mode))
	}
	return b.Bytes()
}

func renderRules(rules []rule) []byte {
	var b bytes.Buffer
	b.WriteString("rules:\n")
	if len(rules) == 0 {
		return b.Bytes()
	}
	for _, r := range rules {
		fmt.Fprintf(&b, "  - type: %s\n    name: %s\n", r.Type, yamlString(r.Name))
		if len(r.Modes) > 0 {
			b.WriteString("    modes:\n")
			for _, mode := range r.Modes {
				fmt.Fprintf(&b, "      - %s\n", yamlString(mode))
			}
		}
		if r.Type == "button-tempo" {
			fmt.Fprintf(&b, "    threshold_ms: %d\n", r.Threshold)
		}
		b.WriteString("    input:\n")
		fmt.Fprintf(&b, "      device: %s\n", yamlString(r.InputDevice))
		switch r.InputKind {
		case "button":
			fmt.Fprintf(&b, "      button: %s\n", r.InputValue)
		case "axis":
			fmt.Fprintf(&b, "      axis: %s\n", r.InputValue)
		case "hat":
			fmt.Fprintf(&b, "      hat: %s\n", r.InputValue)
		}
		if r.Inverted {
			b.WriteString("      inverted: true\n")
		}
		if len(r.Deadzones) > 0 {
			b.WriteString("      deadzones:\n")
			for _, dz := range r.Deadzones {
				fmt.Fprintf(&b, "        - start: %d\n          end: %d\n", dz[0], dz[1])
			}
		}
		switch r.Type {
		case "button-tempo":
			renderTempoBranch(&b, "tap", r.TapOutputs, r.TapMode)
			renderTempoBranch(&b, "hold", r.HoldOutputs, r.HoldMode)
		case "mode-shift":
			fmt.Fprintf(&b, "    mode: %s\n", yamlString(r.ShiftMode))
		default:
			if r.Type == "axis-to-relaxis" {
				fmt.Fprintf(&b, "    repeat_rate_min: %d\n    repeat_rate_max: %d\n    increment: %d\n", r.RepeatMin, r.RepeatMax, r.Increment)
			}
			b.WriteString("    output:\n")
			fmt.Fprintf(&b, "      device: %s\n", yamlString(r.OutputDevice))
			switch r.OutputKind {
			case "button":
				fmt.Fprintf(&b, "      button: %s\n", r.OutputValue)
			case "axis":
				fmt.Fprintf(&b, "      axis: %s\n", r.OutputValue)
			case "hat":
				fmt.Fprintf(&b, "      hat: %s\n", r.OutputValue)
			}
		}
	}
	return b.Bytes()
}

func renderTempoBranch(b *bytes.Buffer, name string, outputs []buttonTarget, mode string) {
	fmt.Fprintf(b, "    %s:\n", name)
	if len(outputs) > 0 {
		b.WriteString("      outputs:\n")
		for _, output := range outputs {
			fmt.Fprintf(b, "        - device: %s\n          button: %d\n", yamlString(output.Device), output.Button)
		}
	} else {
		b.WriteString("      outputs: []\n")
	}
	if mode != "" {
		fmt.Fprintf(b, "      mode: %s\n", yamlString(mode))
	}
}

func renderReadme(modes []string, counts map[string]int) []byte {
	return []byte(fmt.Sprintf(`# Converted Joystick Gremlin Profile

This directory is generated deterministically from `+"`joystick_gremlin_profile.xml`"+` by `+"`go run ./cmd/gremlin-convert`"+`.

Regenerate and run it from the repository root:

`+"```sh"+`
go run ./cmd/gremlin-convert --input joystick_gremlin_profile.xml --output joystick_gremlin_profile
joyful --config ./joystick_gremlin_profile
`+"```"+`

The profile uses two physical VKB Gladiator devices and separate left and right virtual joysticks. Mode inheritance has been flattened and equivalent resolved behavior is grouped into rules with explicit mode lists.

The converter translates VKB axis polarity from Gremlin/DirectInput to Linux evdev and accounts for the different `+"`ABS_RX`"+` / `+"`ABS_RY`"+` mini-stick ordering.

Modes, in startup order: %s.

On each controller, the mini-stick maps to virtual `+"`ABS_Z`"+` / `+"`ABS_RX`"+` (KDE axes 3/4), twist maps to `+"`ABS_RY`"+` (axis 5), and the remaining throttle/lever axis maps to `+"`ABS_RZ`"+` (axis 6). The original right-controller virtual mouse mapping and its acceleration curve are intentionally omitted.

Gremlin vJoy buttons 124 through 128 exceed Joyful's output range and are remapped to virtual buttons 60 through 64 (Joyful indices 59 through 63). Rebind those five controls in the target game.

The freelook recenter and Light Amplification macros are intentionally omitted. See `+"`conversion-report.txt`"+` for all omissions, remaps, and timing nuances.

Generated rules: %d.
`, joinHuman(modes), totalCount(counts)))
}

func renderReport(counts map[string]int, warnings []string) []byte {
	var b bytes.Buffer
	b.WriteString("Joystick Gremlin v14 to Joyful conversion report\n\nRule counts:\n")
	types := make([]string, 0, len(counts))
	for typ := range counts {
		types = append(types, typ)
	}
	sort.Strings(types)
	for _, typ := range types {
		fmt.Fprintf(&b, "  %s: %d\n", typ, counts[typ])
	}
	fmt.Fprintf(&b, "  total: %d\n\nWarnings (%d):\n", totalCount(counts), len(warnings))
	if len(warnings) == 0 {
		b.WriteString("  none\n")
	} else {
		for _, warning := range warnings {
			fmt.Fprintf(&b, "  - %s\n", warning)
		}
	}
	return b.Bytes()
}

func yamlString(value string) string { return strconv.Quote(value) }
func totalCount(counts map[string]int) int {
	total := 0
	for _, count := range counts {
		total += count
	}
	return total
}
func joinHuman(values []string) string {
	out := ""
	for i, value := range values {
		if i > 0 {
			out += ", "
		}
		out += value
	}
	return out
}
