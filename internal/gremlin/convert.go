package gremlin

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

const (
	DefaultLeftPath  = "/dev/input/by-id/usb-VKB-Sim__C__Alex_Oz_2023_VKBsim_Gladiator_EVO_OT_L-event-joystick"
	DefaultRightPath = "/dev/input/by-id/usb-VKB-Sim__C__Alex_Oz_2023_VKBsim_Gladiator_EVO_R-event-joystick"
)

var (
	// DirectInput and Linux expose the VKB mini-stick rotation axes in the
	// opposite order. Keep physical input numbering separate from vJoy output
	// numbering so Gremlin axis 4/5 retain their intended controls.
	physicalAxisNames = map[int]string{1: "X", 2: "Y", 3: "Z", 4: "RY", 5: "RX", 6: "RZ"}
	vjoyAxisNames     = map[int]string{1: "X", 2: "Y", 3: "Z", 4: "RX", 5: "RY", 6: "RZ"}
)

type Options struct {
	LeftPath  string
	RightPath string
}

type Result struct {
	Files    map[string][]byte
	Warnings []string
	Counts   map[string]int
}

type rule struct {
	Type         string
	Name         string
	Modes        []string
	InputDevice  string
	InputKind    string
	InputValue   string
	Inverted     bool
	Deadzones    [][2]int
	OutputDevice string
	OutputKind   string
	OutputValue  string
	Threshold    int
	TapOutputs   []buttonTarget
	TapMode      string
	HoldOutputs  []buttonTarget
	HoldMode     string
	RepeatMin    int
	RepeatMax    int
	Increment    int
	ShiftMode    string
}

type buttonTarget struct {
	Device string
	Button int
}

type converter struct {
	p        *Profile
	warnings map[string]bool
}

func Convert(p *Profile, opts Options) (*Result, error) {
	if opts.LeftPath == "" {
		opts.LeftPath = DefaultLeftPath
	}
	if opts.RightPath == "" {
		opts.RightPath = DefaultRightPath
	}
	c := &converter{p: p, warnings: map[string]bool{}}
	modes := orderedModes(p.Modes)

	type baseKey struct {
		DeviceID string
		Type     string
		ID       int
	}
	keySet := map[baseKey]bool{}
	for key := range p.Inputs {
		keySet[baseKey{key.DeviceID, key.Type, key.ID}] = true
	}
	keys := make([]baseKey, 0, len(keySet))
	for key := range keySet {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		di, dj := p.Devices[keys[i].DeviceID].Side, p.Devices[keys[j].DeviceID].Side
		if di != dj {
			return di < dj
		}
		if keys[i].Type != keys[j].Type {
			return keys[i].Type < keys[j].Type
		}
		return keys[i].ID < keys[j].ID
	})

	grouped := map[string]*rule{}
	order := []string{}
	for _, key := range keys {
		if key.Type == "button" && key.ID >= 125 && key.ID <= 128 {
			c.warn(fmt.Sprintf("omitted unsupported physical %s button %d (the connected Linux device exposes buttons 1..79)", p.Devices[key.DeviceID].Side, key.ID))
			continue
		}
		for _, mode := range modes {
			input, ok := p.resolvedInput(key.DeviceID, key.Type, key.ID, mode)
			if !ok {
				continue
			}
			rules, err := c.convertInput(input)
			if err != nil {
				return nil, fmt.Errorf("%s %s %d in %s: %w", p.Devices[key.DeviceID].Side, key.Type, key.ID, mode, err)
			}
			for _, r := range rules {
				// Gremlin serializes the release side of this temporary shift as a
				// second action in Modifier. Joyful's mode-shift restores the mode it
				// remembered on press, so generating this action would immediately
				// undo the shift on the original press event.
				if r.Type == "mode-shift" && mode == "Modifier" && r.ShiftMode == "SCM Mode" {
					c.warn("collapsed the Modifier-to-SCM temporary action into the momentary Modifier mode-shift restore behavior")
					continue
				}
				sig := r.signature()
				if existing, ok := grouped[sig]; ok {
					existing.Modes = append(existing.Modes, mode)
				} else {
					r.Modes = []string{mode}
					grouped[sig] = &r
					order = append(order, sig)
				}
			}
		}
	}

	rules := make([]rule, 0, len(order))
	for _, sig := range order {
		rules = append(rules, *grouped[sig])
	}
	sort.SliceStable(rules, func(i, j int) bool {
		if rules[i].InputDevice != rules[j].InputDevice {
			return rules[i].InputDevice < rules[j].InputDevice
		}
		if rules[i].InputKind != rules[j].InputKind {
			return rules[i].InputKind < rules[j].InputKind
		}
		if rules[i].InputValue != rules[j].InputValue {
			return numericLess(rules[i].InputValue, rules[j].InputValue)
		}
		return rules[i].signature() < rules[j].signature()
	})

	warnings := make([]string, 0, len(c.warnings))
	for warning := range c.warnings {
		warnings = append(warnings, warning)
	}
	sort.Strings(warnings)
	counts := map[string]int{}
	for _, r := range rules {
		counts[r.Type]++
	}
	files := renderFiles(opts, modes, rules, warnings, counts)
	return &Result{Files: files, Warnings: warnings, Counts: counts}, nil
}

func (p *Profile) resolvedInput(deviceID, typ string, id int, mode string) (Input, bool) {
	for current := mode; current != ""; current = p.Modes[current] {
		if input, ok := p.Inputs[inputKey{deviceID, typ, current, id}]; ok {
			return input, true
		}
	}
	return Input{}, false
}

func orderedModes(modes map[string]string) []string {
	preferred := []string{"SCM Mode", "Auxiliary Mode", "Modifier", "Nav Mode"}
	out := []string{}
	seen := map[string]bool{}
	for _, mode := range preferred {
		if _, ok := modes[mode]; ok {
			out = append(out, mode)
			seen[mode] = true
		}
	}
	rest := []string{}
	for mode := range modes {
		if !seen[mode] {
			rest = append(rest, mode)
		}
	}
	sort.Strings(rest)
	return append(out, rest...)
}

func (c *converter) convertInput(input Input) ([]rule, error) {
	root := c.p.Actions[input.Root]
	if err := validateActionShape(root); err != nil {
		return nil, err
	}
	actions, err := actionRefs(root, "actions")
	if err != nil {
		return nil, err
	}
	var executable []*node
	var curve *node
	for _, id := range actions {
		a := c.p.Actions[id]
		if err := validateActionShape(a); err != nil {
			return nil, fmt.Errorf("action %s: %w", id, err)
		}
		switch a.Attr["type"] {
		case "description":
		case "response-curve":
			curve = a
		default:
			executable = append(executable, a)
		}
	}
	if len(executable) != 1 {
		return nil, fmt.Errorf("root must contain exactly one executable action, found %d", len(executable))
	}
	a := executable[0]
	name := actionLabel(a)
	if name == "" || name == "None" {
		name = fmt.Sprintf("%s %s %d", c.p.Devices[input.DeviceID].Side, input.Type, input.ID)
	}
	base := rule{Name: name, InputDevice: c.p.Devices[input.DeviceID].Side, InputKind: input.Type, InputValue: inputValue(input.Type, input.ID)}
	switch a.Attr["type"] {
	case "map-to-vjoy":
		return c.convertVJoy(base, a, curve)
	case "map-to-mouse":
		return c.convertMouse(base, a, curve)
	case "change-mode":
		props, _ := properties(a)
		changeType, err := oneProperty(props, "change-type")
		if err != nil {
			return nil, err
		}
		if changeType != "Temporary" {
			return nil, fmt.Errorf("root change-mode must be Temporary, got %q", changeType)
		}
		mode, err := targetMode(a)
		if err != nil {
			return nil, err
		}
		base.Type, base.ShiftMode = "mode-shift", mode
		return []rule{base}, nil
	case "tempo":
		return c.convertTempo(base, a)
	case "macro":
		c.warn(fmt.Sprintf("omitted macro action %s on %s %s %d", a.Attr["id"], base.InputDevice, input.Type, input.ID))
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupported root executable type %q", a.Attr["type"])
	}
}

func (c *converter) convertVJoy(base rule, a, curve *node) ([]rule, error) {
	props, _ := properties(a)
	deviceText, err := oneProperty(props, "vjoy-device-id")
	if err != nil {
		return nil, err
	}
	idText, err := oneProperty(props, "vjoy-input-id")
	if err != nil {
		return nil, err
	}
	typ, err := oneProperty(props, "vjoy-input-type")
	if err != nil {
		return nil, err
	}
	deviceID, e1 := strconv.Atoi(deviceText)
	id, e2 := strconv.Atoi(idText)
	if e1 != nil || e2 != nil || (deviceID != 1 && deviceID != 2) {
		return nil, fmt.Errorf("invalid vJoy target %q/%q", deviceText, idText)
	}
	base.OutputDevice = map[int]string{1: "left-vjoy", 2: "right-vjoy"}[deviceID]
	switch typ {
	case "button":
		button, err := c.vjoyButton(id)
		if err != nil {
			return nil, err
		}
		base.Type, base.OutputKind, base.OutputValue = "button", "button", strconv.Itoa(button)
		return []rule{base}, nil
	case "axis":
		axis, ok := vjoyAxisNames[id]
		if !ok {
			return nil, fmt.Errorf("unsupported vJoy axis %d", id)
		}
		inverted, err := curveInverted(curve, false)
		if err != nil {
			return nil, err
		}
		if base.InputValue == "RX" || base.InputValue == "RY" {
			// Use the requested six-axis layout: the mini-stick occupies KDE
			// axes 3/4 (ABS_Z/ABS_RX), leaving ABS_RY for twist.
			miniOutput := map[string]string{"RX": "Z", "RY": "RX"}[base.InputValue]
			base.Type, base.OutputKind, base.OutputValue = "axis", "axis", miniOutput
			base.Inverted = base.InputValue == "RY"
			base.Name = fmt.Sprintf("%s mini-stick %s", title(base.InputDevice), map[string]string{"RX": "X", "RY": "Y"}[base.InputValue])
			return []rule{base}, nil
		}
		// Linux evdev and Gremlin/DirectInput use opposite joystick axis
		// polarity for these devices. Toggle the Gremlin curve's inversion when
		// producing the Joyful input transform.
		base.Type, base.OutputKind, base.OutputValue, base.Inverted = "axis", "axis", axis, !inverted
		if (base.InputDevice == "left" && base.InputValue == "Y") || base.InputValue == "Z" {
			base.Inverted = false
		}
		if base.InputValue == "RZ" {
			base.OutputValue = "RY"
			base.Name = fmt.Sprintf("%s twist", title(base.InputDevice))
		}
		return []rule{base}, nil
	case "hat":
		if curve != nil {
			return nil, fmt.Errorf("hat cannot have response curve")
		}
		hat := id - 1
		if hat < 0 || hat > 3 {
			return nil, fmt.Errorf("unsupported vJoy hat %d", id)
		}
		out := make([]rule, 2)
		for i, suffix := range []string{"X", "Y"} {
			out[i] = base
			out[i].Type = "hat"
			out[i].InputValue = "HAT0" + suffix
			out[i].OutputKind = "hat"
			out[i].OutputValue = fmt.Sprintf("HAT%d%s", hat, suffix)
			out[i].Name += " " + suffix
		}
		return out, nil
	default:
		return nil, fmt.Errorf("unsupported vJoy input type %q", typ)
	}
}

func (c *converter) convertMouse(base rule, a, curve *node) ([]rule, error) {
	if base.InputKind != "axis" {
		return nil, fmt.Errorf("map-to-mouse Motion requires axis input")
	}
	if err := flatCenterCurve(curve); err != nil {
		return nil, err
	}
	props, _ := properties(a)
	mode, err := oneProperty(props, "mode")
	if err != nil || mode != "Motion" {
		return nil, fmt.Errorf("map-to-mouse mode must be Motion")
	}
	directionText, err := oneProperty(props, "direction")
	if err != nil {
		return nil, err
	}
	direction, err := strconv.Atoi(directionText)
	if err != nil || (direction != 0 && direction != 90) {
		return nil, fmt.Errorf("invalid mouse direction")
	}
	expectedInput := map[int]string{0: "RX", 90: "RY"}[direction]
	if base.InputValue != expectedInput {
		return nil, fmt.Errorf("mouse direction %d does not match translated physical axis %s", direction, base.InputValue)
	}
	base.Type = "axis"
	base.OutputDevice = base.InputDevice + "-vjoy"
	base.OutputKind = "axis"
	base.OutputValue = map[string]string{"RX": "Z", "RY": "RX"}[base.InputValue]
	base.Inverted = base.InputValue == "RY"
	base.Name = fmt.Sprintf("%s mini-stick %s", title(base.InputDevice), map[string]string{"RX": "X", "RY": "Y"}[base.InputValue])
	c.warn(fmt.Sprintf("converted Gremlin mouse Motion direction %d to virtual joystick axis %s and omitted its mouse acceleration curve", direction, base.OutputValue))
	return []rule{base}, nil
}

func (c *converter) convertTempo(base rule, a *node) ([]rule, error) {
	props, _ := properties(a)
	threshold, err := oneProperty(props, "threshold")
	if err != nil {
		return nil, err
	}
	seconds, err := strconv.ParseFloat(threshold, 64)
	if err != nil || seconds < 0 {
		return nil, fmt.Errorf("invalid tempo threshold %q", threshold)
	}
	activateOn, err := oneProperty(props, "activate-on")
	if err != nil {
		return nil, err
	}
	if activateOn != "press" && activateOn != "release" {
		return nil, fmt.Errorf("unsupported tempo activate-on %q", activateOn)
	}
	base.Type, base.Threshold = "button-tempo", int(math.Round(seconds*1000))
	short, err := actionRefs(a, "short-actions")
	if err != nil {
		return nil, err
	}
	long, err := actionRefs(a, "long-actions")
	if err != nil {
		return nil, err
	}
	base.TapOutputs, base.TapMode, _, err = c.tempoBranch(short, "tap")
	if err != nil {
		return nil, err
	}
	base.HoldOutputs, base.HoldMode, _, err = c.tempoBranch(long, "hold")
	if err != nil {
		return nil, err
	}
	if activateOn == "press" {
		c.warn("tempo activate-on=press is represented by button-tempo's release/threshold state machine; output press/release behavior is externally equivalent, but exact dispatch timing may differ")
	}
	return []rule{base}, nil
}

func (c *converter) tempoBranch(ids []string, branch string) ([]buttonTarget, string, bool, error) {
	var outputs []buttonTarget
	mode := ""
	macro := false
	for _, id := range ids {
		a := c.p.Actions[id]
		if err := validateActionShape(a); err != nil {
			return nil, "", false, err
		}
		switch a.Attr["type"] {
		case "map-to-vjoy":
			props, _ := properties(a)
			typ, err := oneProperty(props, "vjoy-input-type")
			if err != nil || typ != "button" {
				return nil, "", false, fmt.Errorf("tempo branch vJoy action must be a button")
			}
			devText, _ := oneProperty(props, "vjoy-device-id")
			idText, _ := oneProperty(props, "vjoy-input-id")
			dev, e1 := strconv.Atoi(devText)
			buttonID, e2 := strconv.Atoi(idText)
			if e1 != nil || e2 != nil || (dev != 1 && dev != 2) {
				return nil, "", false, fmt.Errorf("invalid tempo vJoy target")
			}
			button, err := c.vjoyButton(buttonID)
			if err != nil {
				return nil, "", false, err
			}
			outputs = append(outputs, buttonTarget{map[int]string{1: "left-vjoy", 2: "right-vjoy"}[dev], button})
		case "change-mode":
			props, _ := properties(a)
			changeType, err := oneProperty(props, "change-type")
			if err != nil {
				return nil, "", false, err
			}
			if changeType != "Cycle" {
				return nil, "", false, fmt.Errorf("tempo mode action must be Cycle, got %q", changeType)
			}
			if mode != "" {
				return nil, "", false, fmt.Errorf("multiple mode targets in tempo %s branch", branch)
			}
			mode, err = targetMode(a)
			if err != nil {
				return nil, "", false, err
			}
		case "macro":
			macro = true
			c.warn(fmt.Sprintf("omitted macro action %s from tempo %s branch; other branch actions were preserved", id, branch))
		default:
			return nil, "", false, fmt.Errorf("unsupported tempo %s action type %q", branch, a.Attr["type"])
		}
	}
	return outputs, mode, macro, nil
}

func (c *converter) vjoyButton(id int) (int, error) {
	if id >= 124 && id <= 128 {
		mapped := id - 65
		c.warn(fmt.Sprintf("remapped vJoy button %d to Joyful index %d (vJoy-style button %d)", id, mapped, mapped+1))
		return mapped, nil
	}
	if id < 1 || id > 74 {
		return 0, fmt.Errorf("unsupported vJoy button %d", id)
	}
	return id - 1, nil
}

func (c *converter) warn(message string) { c.warnings[message] = true }

func actionRefs(a *node, group string) ([]string, error) {
	containers := a.children(group)
	if len(containers) != 1 {
		return nil, fmt.Errorf("%s action requires one <%s>", a.Attr["type"], group)
	}
	refs := []string{}
	for _, child := range containers[0].Children {
		if child.Name != "action-id" || child.Text == "" {
			return nil, fmt.Errorf("invalid <%s> child <%s>", group, child.Name)
		}
		refs = append(refs, child.Text)
	}
	return refs, nil
}

func targetMode(a *node) (string, error) {
	target, err := a.child("target-mode")
	if err != nil {
		return "", err
	}
	props, err := properties(target)
	if err != nil {
		return "", err
	}
	mode, err := oneProperty(props, "name")
	if err != nil {
		return "", err
	}
	return mode, nil
}

func actionLabel(a *node) string {
	props, _ := properties(a)
	values := props["action-label"]
	if len(values) == 1 {
		return values[0]
	}
	return ""
}

func inputValue(typ string, id int) string {
	if typ == "axis" {
		return physicalAxisNames[id]
	}
	if typ == "hat" {
		return "HAT0"
	}
	return strconv.Itoa(id - 1)
}

func curveInverted(curve *node, required bool) (bool, error) {
	if curve == nil {
		if required {
			return false, fmt.Errorf("missing response curve")
		}
		return false, nil
	}
	props, _ := properties(curve)
	typ, err := oneProperty(props, "curve-type")
	if err != nil || typ != "PiecewiseLinear" {
		return false, fmt.Errorf("response curve must be PiecewiseLinear")
	}
	points, err := curvePoints(curve)
	if err != nil {
		return false, err
	}
	if equalPoints(points, [][2]float64{{-1, -1}, {1, 1}}) {
		return false, nil
	}
	if equalPoints(points, [][2]float64{{-1, 1}, {1, -1}}) {
		return true, nil
	}
	return false, fmt.Errorf("unsupported non-linear ordinary response curve")
}

func flatCenterCurve(curve *node) error {
	if curve == nil {
		return fmt.Errorf("mouse Motion requires flat-center response curve")
	}
	props, _ := properties(curve)
	typ, err := oneProperty(props, "curve-type")
	if err != nil || typ != "PiecewiseLinear" {
		return fmt.Errorf("mouse response curve must be PiecewiseLinear")
	}
	points, err := curvePoints(curve)
	if err != nil {
		return err
	}
	if !equalPoints(points, [][2]float64{{-1, -1}, {-.1, 0}, {.1, 0}, {1, 1}}) {
		return fmt.Errorf("unsupported mouse response curve")
	}
	return nil
}

func curvePoints(curve *node) ([][2]float64, error) {
	control, err := curve.child("control-points")
	if err != nil {
		return nil, err
	}
	props, err := properties(control)
	if err != nil {
		return nil, err
	}
	points := make([][2]float64, 0, len(props["point"]))
	for _, value := range props["point"] {
		parts := strings.Split(value, ",")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid curve point %q", value)
		}
		x, e1 := strconv.ParseFloat(parts[0], 64)
		y, e2 := strconv.ParseFloat(parts[1], 64)
		if e1 != nil || e2 != nil {
			return nil, fmt.Errorf("invalid curve point %q", value)
		}
		points = append(points, [2]float64{x, y})
	}
	return points, nil
}

func equalPoints(a, b [][2]float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if math.Abs(a[i][0]-b[i][0]) > 1e-9 || math.Abs(a[i][1]-b[i][1]) > 1e-9 {
			return false
		}
	}
	return true
}

func validateActionShape(a *node) error {
	allowed := map[string]map[string]bool{
		"root": {"actions": true, "property": true}, "description": {"property": true},
		"response-curve": {"deadzone": true, "control-points": true, "property": true},
		"map-to-vjoy":    {"property": true}, "map-to-mouse": {"property": true},
		"change-mode": {"property": true, "target-mode": true},
		"tempo":       {"short-actions": true, "long-actions": true, "property": true},
		"macro":       {"property": true, "macro-action": true},
	}
	typ := a.Attr["type"]
	for _, child := range a.Children {
		if !allowed[typ][child.Name] {
			return fmt.Errorf("unknown executable structure <%s> in %s", child.Name, typ)
		}
	}
	allowedProperties := map[string]map[string]bool{
		"root":           {"action-label": true, "activation-mode": true},
		"description":    {"description": true, "action-label": true, "activation-mode": true},
		"response-curve": {"curve-type": true, "action-label": true, "activation-mode": true},
		"map-to-vjoy":    {"vjoy-device-id": true, "vjoy-input-id": true, "vjoy-input-type": true, "button-inverted": true, "axis-mode": true, "axis-scaling": true, "action-label": true, "activation-mode": true},
		"map-to-mouse":   {"mode": true, "direction": true, "min-speed": true, "max-speed": true, "time-to-max-speed": true, "action-label": true, "activation-mode": true},
		"change-mode":    {"change-type": true, "action-label": true, "activation-mode": true},
		"tempo":          {"threshold": true, "activate-on": true, "action-label": true, "activation-mode": true},
	}
	if typ != "macro" {
		props, err := properties(a)
		if err != nil {
			return err
		}
		for name, values := range props {
			if !allowedProperties[typ][name] {
				return fmt.Errorf("unknown property %q in %s", name, typ)
			}
			if len(values) != 1 {
				return fmt.Errorf("duplicate property %q in %s", name, typ)
			}
		}
	}
	if typ == "response-curve" {
		deadzone, err := a.child("deadzone")
		if err != nil {
			return err
		}
		if err := validatePropertyContainer(deadzone, map[string]bool{"low": true, "center-low": true, "center-high": true, "high": true}, false); err != nil {
			return fmt.Errorf("deadzone: %w", err)
		}
		control, err := a.child("control-points")
		if err != nil {
			return err
		}
		if err := validatePropertyContainer(control, map[string]bool{"point": true}, true); err != nil {
			return fmt.Errorf("control-points: %w", err)
		}
	}
	return nil
}

func validatePropertyContainer(n *node, allowed map[string]bool, allowDuplicates bool) error {
	for _, child := range n.Children {
		if child.Name != "property" {
			return fmt.Errorf("unknown structure <%s>", child.Name)
		}
	}
	props, err := properties(n)
	if err != nil {
		return err
	}
	for name, values := range props {
		if !allowed[name] {
			return fmt.Errorf("unknown property %q", name)
		}
		if !allowDuplicates && len(values) != 1 {
			return fmt.Errorf("duplicate property %q", name)
		}
	}
	return nil
}

func (r rule) signature() string { copy := r; copy.Modes = nil; return fmt.Sprintf("%#v", copy) }
func numericLess(a, b string) bool {
	ai, e1 := strconv.Atoi(a)
	bi, e2 := strconv.Atoi(b)
	if e1 == nil && e2 == nil {
		return ai < bi
	}
	return a < b
}

func title(value string) string {
	if value == "" {
		return value
	}
	return strings.ToUpper(value[:1]) + value[1:]
}
