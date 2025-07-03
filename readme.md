# Joyful - virtual joystick remapper for Linux

Joyful is a Linux tool for creating virtual joysticks and mapping inputs from various
real devices to them. This is useful when playing games that don't support multiple joysticks,
or for games that don't handle devices changing order (like Star Citizen).

Perhaps more significantly, Joyful allows you to map combinations of physical inputs to a single output,
as well as creating other complex scenarios. Want a single button press to simultaneously produce multiple outputs?
Joyful can do that!

Are you a Linux gamer who misses the features of Joystick Gremlin? Wish you could map combo inputs,
or combine your input devices into one virtual device? Are you ok with writing a bunch of YAML?

Joyful might be the tool for you.

## Features

### Current Features - try them today!

* Create virtual devices with up to 8 axes and 80 buttons.
* Make simple 1:1 mappings of buttons and axes: Button1 -> VirtualButtonA
* Make combination mappings: Button1 + Button2 -> VirtualButtonA

### Future Features - try them at an unspecified point in the future!

* Multiple modes with per-mode behavior.
* Partial axis mapping: map sections of an axis to different outputs.
* Highly configurable deadzones
* Macros - have a single input produce a sequence of button presses with configurable pauses.
* Sequence combos - Button1, Button2, Button3 -> VirtualButtonA

## Configuration

Configuration is currently done via hand-written YAML files in `~/.config/joyful/`. Joyful will read every
yaml file in this directory and combine them, so you can split your configuration up however you like.

Configuration is divided into two sections: `devices` and `rules`. Each of these is a YAML list.
The options for each are described in some detail below. See the `examples/` directory for concrete examples.

### Device configuration

Each entry in `devices` must have a couple of fields:

* `name` - This is an identifier that your rules will use to refer to the device. It is recommended to avoid spaces or special characters.
* `type` - Should be `physical` for an input device, and `virtual` for an output device.

`physical` devices must additionally define these fields:

* `device_name` - The name of the device as reported by the included `evlist` command. If your device name ends with a space, use quotation marks (`""`) around the name.

`virtual` devices must additionally define these fields:

* `buttons` - a number between 0 and 80. Linux may not recognize buttons greater than 56.
* `axes` - a number between 0 and 8.

### Rules configuration

All `rules` must have a `type` field. Valid values for this field are:

* `simple` - a single input mapped to a single output
* `combo` - multiple inputs mapped to a single output. The output event will trigger when all the input conditions are met.
* `latched` - a single input mapped to a single output, but each time the input is pressed, the output will toggle.

Configuration options for each type vary. See <examples/ruletypes.yml> for an example of each type with all options specified.

### Modes

All rules can have a `modes` field that is a list of strings.


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