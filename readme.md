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

* Create virtual devices with up to 8 axes and 56 buttons.
* Make simple 1:1 mappings of buttons and axes: Button1 -> VirtualButtonA
* Make combination mappings: Button1 + Button2 -> VirtualButtonA
* Multiple modes with per-mode behavior.

### Future Features - try them at an unspecified point in the future!

* Partial axis mapping: map sections of an axis to different outputs.
* Highly configurable deadzones
* Macros - have a single input produce a sequence of button presses with configurable pauses.
* Sequence combos - Button1, Button2, Button3 -> VirtualButtonA
* Proportional axis to button mapping; repeatedly trigger a button with an axis, with frequency controlled by the axis value
* More ways to specify keycodes

## Configuration

Configuration is currently done via hand-written YAML files in `~/.config/joyful/`. Joyful will read every
yaml file in this directory and combine them, so you can split your configuration up however you like.

Configuration is divided into three sections: `devices`, `modes`, and `rules`. See the `examples/` directory for concrete examples.
Select options are explained in detail below.

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

### Keycodes

Currently, there is only one way to specify a button or axis: using evdev's Keycodes. These look like `ABS_X` for axes and `BTN_TRIGGER`
for buttons. See <https://github.com/holoplot/go-evdev/blob/master/codes.go> for a full list of these codes, but note that Joyful's virtual devices currently only uses a subset. Specifically, the axes from `ABS_X` to `ABS_RUDDER`, and the buttons from `BTN_JOYSTICK` to `BTN_DEAD`, as well as all of the `BTN_TRIGGER_HAPPY*` codes.

For input, you can figure out what keycodes your device is emitting by running the Linux utility `evtest`. `evtest` works well with `grep`, so if you just want to see button inputs, you can do:

```
evtest | grep 
```

The authors of this tool recognize that this is currently a pain in the ass. Easier ways to represent keycodes (as well as outputting additional keycodes) is planned for the future.

We don't have the cycles to develop tool-assisted configuration, but pull requests (or separate projects that produce compatible YAML) are very welcome!

### Modes

The top-level `modes` field is a simple list of strings, defining the different modes available to rules. The initial mode is always
the first one in the list. (TODO)

All rules can have a `modes` field that is a list of strings. If no `modes` field is present, the rule will be active in all modes.


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