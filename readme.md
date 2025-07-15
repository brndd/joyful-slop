# Joyful - joystick remapper for Linux

Joyful is a Linux tool for mapping inputs from various joystick-like devices to "virtual" output devices. This is useful when playing games that don't support multiple joysticks, or for games that don't gracefully handle devices changing order (e.g., Star Citizen).

Joyful also allows you to map combinations of physical inputs to a single output, as well as creating other complex scenarios.

Joyful is ideal for Linux gamers who enjoy space and flight sims and miss the features of Joystick Gremlin.

## Features

### Current Features - try them today!

* Create virtual devices with up to 8 axes and 56 buttons.
* Make simple 1:1 mappings of buttons and axes: Button1 -> VirtualButtonA
* Make combination mappings: Button1 + Button2 -> VirtualButtonA
* Define multiple modes with per-mode behavior.
* "Split" axis mapping: map sections of an axis to different outputs.
* Configure per-mapping configurable deadzones for axes.
* Axis -> button mapping with optional "proportional" repeat speed (i.e. repeat faster as the axis is engaged further)
* Axis -> Relative Axis mapping, for converting a joystick axis to mouse movement and scrollwheel events.

### Future Features - try them at an unspecified point in the future!

* Macros - have a single input produce a sequence of button presses with configurable pauses.
* Sequence combos - Button1, Button2, Button3 -> VirtualButtonA
* More ways to specify keycodes
* Output keyboard button presses
* Input and output from gamepad-like devices.

## Configuration

Configuration is handled via YAML files in `~/.config/joyful/`. Joyful will read every yaml file in this directory and combine them, so you can split your configuration up however you like. 

A configuration guide and examples can be found in the `docs/` directory.

Configuration can be fairly complicated and repetitive. If anyone wants to create a graphical interface to configure Joyful, we would love to link to it here.

## Technical details

Joyful is written in golang, and uses evdev/uinput to manage devices. 

### Building

To build joyful, install `go` via your package manager, then run:

```
go build -o build/ ./...
```

Look for binaries in the `build/` directory.

### Contributing

Send patches and questions to [annabunches@gmail.com](mailto:annabunches@gmail.com). Make sure the subject of your email starts with `[Joyful]`.

If enough people show an interest in contributing, I'll consider mirroring the repository on Github.