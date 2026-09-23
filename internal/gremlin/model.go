package gremlin

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Profile struct {
	Devices map[string]Device
	Modes   map[string]string
	Inputs  map[inputKey]Input
	Actions map[string]*node
}

type Device struct {
	ID   string
	Name string
	Side string
}

type Input struct {
	DeviceID string
	Type     string
	Mode     string
	ID       int
	Root     string
}

type inputKey struct {
	DeviceID string
	Type     string
	Mode     string
	ID       int
}

var knownActionTypes = map[string]bool{
	"root": true, "description": true, "response-curve": true, "map-to-vjoy": true,
	"map-to-mouse": true, "change-mode": true, "tempo": true, "macro": true,
}

func Parse(r io.Reader) (*Profile, error) {
	root, err := parseXML(r)
	if err != nil {
		return nil, err
	}
	if root.Name != "profile" || root.Attr["version"] != "14" {
		return nil, fmt.Errorf("expected Joystick Gremlin profile version 14")
	}
	p := &Profile{Devices: map[string]Device{}, Modes: map[string]string{}, Inputs: map[inputKey]Input{}, Actions: map[string]*node{}}

	devices, err := root.child("devices")
	if err != nil {
		return nil, err
	}
	for _, n := range devices.children("device") {
		id, e1 := n.value("device-id")
		name, e2 := n.value("device-name")
		if e1 != nil || e2 != nil {
			return nil, fmt.Errorf("invalid device: %v %v", e1, e2)
		}
		if _, exists := p.Devices[id]; exists {
			return nil, fmt.Errorf("duplicate device UUID %q", id)
		}
		norm := strings.ToLower(strings.Join(strings.Fields(name), " "))
		side := ""
		if strings.Contains(norm, "ot l") || strings.HasSuffix(norm, " l") {
			side = "left"
		} else if strings.HasSuffix(norm, " r") || strings.Contains(norm, "ot r") {
			side = "right"
		}
		if side == "" {
			return nil, fmt.Errorf("cannot identify physical device side from %q", name)
		}
		p.Devices[id] = Device{ID: id, Name: strings.TrimSpace(name), Side: side}
	}

	modes, err := root.child("modes")
	if err != nil {
		return nil, err
	}
	for _, n := range modes.children("mode") {
		if n.Text == "" {
			return nil, fmt.Errorf("empty mode name")
		}
		if _, exists := p.Modes[n.Text]; exists {
			return nil, fmt.Errorf("duplicate mode %q", n.Text)
		}
		p.Modes[n.Text] = n.Attr["parent"]
	}
	if err := validateModes(p.Modes); err != nil {
		return nil, err
	}

	library, err := root.child("library")
	if err != nil {
		return nil, err
	}
	for _, action := range library.children("action") {
		id, typ := action.Attr["id"], action.Attr["type"]
		if id == "" || typ == "" {
			return nil, fmt.Errorf("library action missing id or type")
		}
		if !knownActionTypes[typ] {
			return nil, fmt.Errorf("action %s: unknown executable action type %q", id, typ)
		}
		if _, exists := p.Actions[id]; exists {
			return nil, fmt.Errorf("duplicate action UUID %q", id)
		}
		if err := validateActionShape(action); err != nil {
			return nil, fmt.Errorf("action %s: %w", id, err)
		}
		p.Actions[id] = action
	}

	inputs, err := root.child("inputs")
	if err != nil {
		return nil, err
	}
	for _, n := range inputs.children("input") {
		deviceID, e1 := n.value("device-id")
		typ, e2 := n.value("input-type")
		mode, e3 := n.value("mode")
		idText, e4 := n.value("input-id")
		cfg, e5 := n.child("action-configuration")
		if e1 != nil || e2 != nil || e3 != nil || e4 != nil || e5 != nil {
			return nil, fmt.Errorf("invalid input structure")
		}
		id, err := strconv.Atoi(idText)
		if err != nil || id < 1 {
			return nil, fmt.Errorf("invalid %s input id %q", typ, idText)
		}
		if typ != "button" && typ != "axis" && typ != "hat" {
			return nil, fmt.Errorf("unknown input type %q", typ)
		}
		if typ == "axis" && (id < 1 || id > 6) {
			return nil, fmt.Errorf("unsupported physical axis %d", id)
		}
		if typ == "hat" && id != 1 {
			return nil, fmt.Errorf("unsupported physical hat %d", id)
		}
		if typ == "button" && id > 128 {
			return nil, fmt.Errorf("unsupported physical button %d", id)
		}
		if _, ok := p.Devices[deviceID]; !ok {
			return nil, fmt.Errorf("input references unknown device %q", deviceID)
		}
		if _, ok := p.Modes[mode]; !ok {
			return nil, fmt.Errorf("input references unknown mode %q", mode)
		}
		rootID, err := cfg.value("root-action")
		if err != nil {
			return nil, err
		}
		behavior, err := cfg.value("behavior")
		if err != nil {
			return nil, err
		}
		if behavior != typ {
			return nil, fmt.Errorf("input %s/%s/%s/%d has behavior %q", deviceID, mode, typ, id, behavior)
		}
		action, ok := p.Actions[rootID]
		if !ok {
			return nil, fmt.Errorf("input references unknown root action %q", rootID)
		}
		if action.Attr["type"] != "root" {
			return nil, fmt.Errorf("input root reference %q has type %q", rootID, action.Attr["type"])
		}
		key := inputKey{deviceID, typ, mode, id}
		if _, exists := p.Inputs[key]; exists {
			return nil, fmt.Errorf("duplicate input key %s/%s/%s/%d", deviceID, mode, typ, id)
		}
		p.Inputs[key] = Input{deviceID, typ, mode, id, rootID}
	}
	if err := validateReferences(p); err != nil {
		return nil, err
	}
	return p, nil
}

func validateModes(modes map[string]string) error {
	for mode := range modes {
		seen := map[string]bool{}
		for current := mode; current != ""; {
			if seen[current] {
				return fmt.Errorf("mode inheritance cycle at %q", current)
			}
			seen[current] = true
			parent, exists := modes[current]
			if !exists {
				return fmt.Errorf("mode %q inherits unknown mode %q", mode, current)
			}
			current = parent
		}
	}
	return nil
}

func validateReferences(p *Profile) error {
	graph := make(map[string][]string, len(p.Actions))
	for id, action := range p.Actions {
		for _, group := range []string{"actions", "short-actions", "long-actions"} {
			for _, container := range action.children(group) {
				for _, ref := range container.children("action-id") {
					if _, ok := p.Actions[ref.Text]; !ok {
						return fmt.Errorf("action %s references unknown action %q", id, ref.Text)
					}
					graph[id] = append(graph[id], ref.Text)
				}
			}
		}
		if action.Attr["type"] == "change-mode" {
			mode, err := targetMode(action)
			if err != nil {
				return fmt.Errorf("action %s: %w", id, err)
			}
			if _, ok := p.Modes[mode]; !ok {
				return fmt.Errorf("action %s references unknown mode %q", id, mode)
			}
		}
	}
	state := make(map[string]uint8, len(p.Actions))
	var visit func(string) error
	visit = func(id string) error {
		switch state[id] {
		case 1:
			return fmt.Errorf("action reference cycle at %q", id)
		case 2:
			return nil
		}
		state[id] = 1
		for _, child := range graph[id] {
			if err := visit(child); err != nil {
				return err
			}
		}
		state[id] = 2
		return nil
	}
	for id := range p.Actions {
		if err := visit(id); err != nil {
			return err
		}
	}
	return nil
}
