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
* Define keyboard, mouse, and gamepad outputs in addition to joysticks.
* Configure per-rule configurable deadzones for axes, with multiple ways to specify deadzones.
* Define multiple modes with per-mode behavior.
    * Text-to-speech engine that announces the current mode when it changes.

### Possible Future Features

* Macros - have a single input produce a sequence of button presses with configurable pauses.
* Sequence combos - Button1, Button2, Button3 -> VirtualButtonA
* Hat support
* HIDRAW support for more button options.
* Sensitivity Curves?
* Packaged builds for Arch and possibly other distributions.

## Configure

Configuration is handled via YAML files in `~/.config/joyful/`. Joyful will read every yaml file in this directory and combine them, so you can split your configuration up however you like. 

A configuration guide and examples can be found in the `docs/` directory.

Configuration can be fairly complicated and repetitive. If anyone wants to create a graphical interface to configure Joyful, we would love to link to it here.

## Install

If you are on Arch or an Arch-based distro, you can get the latest Joyful release from the AUR:

```
yay -S joyful
```

You may also need to add the user(s) who will be running joyful to the `input` group.

### Manual Install

To build joyful manually, first use your distribution's package manager to install the following dependencies:

* `go`
* `make`
* `alsa-lib` - this may be `libasound2-dev` or `libasound2-devel` depending on your distribution
* `espeak-ng` - if you want text-to-speech to announce mode changes

Then, clone this repository, e.g.:

```
git clone https://git.annabunches.net/anna/joyful.git
cd joyful
```

Then, to build and install, run:

```
go build -o build/ ./...
cp build/* ~/bin/
```

If you want to install Joyful system-wide, you can instead do:

```
go build -o build/ ./...
sudo cp build/* /usr/local/bin/
```


## Usage

After installing Joyful and writing your configuration (see above), run `joyful`. You can use `joyful --config <directory>` to specify different configuration profiles; just put all the YAML files for a given profile in a unique directory.

Pressing `<enter>` in the running terminal window will reload the `rules` section of your config files, so you can make changes to your rules without restarting the application. Applying any changes to `devices` or `modes` requires exiting and re-launching the program.

## Technical details

Joyful is written in golang, and uses `evdev`/`uinput` to manage devices and `espeak-ng` for TTS. See [cmd/joyful/main.go](cmd/joyful/main.go) for the program's entry point.

This was originally going to be a Rust project, but the author's Rust skills weren't quite up to the task yet. Please look forward to the inevitable Rust rewrite.

### Contributing

Issues and pull requests should be made on the [Codeberg mirror](https://codeberg.org/annabunches/joyful).