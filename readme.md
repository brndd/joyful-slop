# Joyful - joystick remapper for Linux

Joyful is a Linux tool for mapping inputs from various joystick-like devices to "virtual" output devices. This is useful when playing games that don't support multiple joysticks, or for games that don't gracefully handle devices changing order (e.g., Star Citizen).

Joyful also allows you to map combinations of physical inputs to a single output, as well as creating other complex scenarios.

Joyful is ideal for Linux gamers who enjoy space and flight sims and miss the features of Joystick Gremlin.

## Features

### Current Features

* Create virtual devices with up to 8 axes and 74 buttons.
* Flexible rule system that allows several different types of rules, including:
    * Simple 1:1 mappings of buttons and axes: Button1 -> VirtualButtonA
    * Combination mappings: Button1 + Button2 -> VirtualButtonA
    * "Split" axis mapping: map sections of an axis to different outputs using deadzones.
    * "Combined" axis mapping: map two physical axes to one virtual axis.
    * Axis -> button mapping with optional "proportional" repeat speed (i.e. repeat faster as the axis is engaged further)
    * Axis -> Relative Axis mapping, for converting a joystick axis to mouse movement and scrollwheel events.
* Configure per-rule configurable deadzones for axes, with multiple ways to specify deadzones.
* Define multiple modes with per-mode behavior.

### Possible Future Features

* Macros - have a single input produce a sequence of button presses with configurable pauses.
* Sequence combos - Button1, Button2, Button3 -> VirtualButtonA
* Output keyboard button presses
* Hat support
* HIDRAW support for more button options.
* Sensitivity Curves.

## Configuration

Configuration is handled via YAML files in `~/.config/joyful/`. Joyful will read every yaml file in this directory and combine them, so you can split your configuration up however you like. 

A configuration guide and examples can be found in the `docs/` directory.

Configuration can be fairly complicated and repetitive. If anyone wants to create a graphical interface to configure Joyful, we would love to link to it here.

## Usage

After building (see below) and writing your configuration (see above), just run `joyful`. You can use `joyful --config <directory>` to specify different configuration profiles; just put all the YAML files for a given profile in a unique directory.

Pressing `<enter>` in the running terminal window will reload the `rules` section of your config files, so you can make changes to your rules without restarting the application. Applying any changes to `devices` or `modes` requires exiting and re-launching the program.

## Technical details

Joyful is written in golang, and uses evdev/uinput to manage devices. See `cmd/joyful/main.go` for the program's entry point.

### Build & Install

To build joyful, install `go` via your package manager, then run:

```
CGO_ENABLED=0 go build -o build/ ./...
```

Copy the binaries in the `build/` directory to somewhere in your `$PATH`. (details depend on your setup, but typically somewhere like `/usr/local/bin` or `~/bin`)

### Contributing

Send patches and questions to [annabunches@gmail.com](mailto:annabunches@gmail.com). Make sure the subject of your email starts with `[Joyful]`.

If enough people show an interest in contributing, I'll consider mirroring the repository on Github.