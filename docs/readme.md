# Joyful Configuration

Configuration is divided into three sections: `devices`, `modes`, and `rules`. Each yaml file can have any number of these sections; joyful will combine the configuration from all files at runtime.

### Device configuration

Each entry in `devices` must have a couple of parameters:

* `name` - This is an identifier that your rules will use to refer to the device. It is recommended to avoid spaces or special characters.
* `type` - Should be `physical` for an input device, and `virtual` for an output device.

`physical` devices must additionally define these parameters:

* `device_name` - The name of the device as reported by the included `evlist` command. If your device name ends with a space, use quotation marks (`""`) around the name.

`virtual` devices can additionally define these parameters:

* `buttons` or `num_buttons` - Either a list of explicit buttons or a number of buttons to create. (max 74 buttons) Linux-native games may not recognize all buttons created by Joyful.
* `axes` or `num_axes` - An explicit list of `ABS_` axes or a number to create.
* `relative_axes` or `num_relative_axes` - As above, but for `REL_` axes.

A couple of additional notes on virtual devices:

* For all 3 of the above options, an explicit list will override the `num_` parameters if both are present.
* Some environments will only register mouse events if the device *only* supports mouse-like events, so it can be useful to isolate your `relative_axes` to their own virtual device and explicitly define the axes.

### Rules configuration

All `rules` must have a `type` parameter. Valid values for this parameter are:

* `button` - a single button mapping
* `button-combo` - multiple input buttons mapped to a single output. The output event will trigger when all the input conditions are met.
* `button-latched` - a single button mapped to a single output, but each time the input is pressed, the output will toggle.
* `axis` - a simple axis mapping
* `axis-to-button` - causes an axis input to produce a button output. This can be repeated with variable speed proportional to the axis' input value
* `axis-to-relaxis` - like axis-to-button, but produces a "relative axis" output value. This is useful for simulating mouse scrollwheel and movement events.

Configuration options for each rule type vary. See <examples/ruletypes.yml> for an example of each type with all options specified.

### Event Codes

Event codes are the values that identify buttons and axes. There are several ways to configure these codes. All of them are case-insensitive, so `abs_x` and `ABS_X`.

Ways to specify event codes are:

* Using evdev's identifiers. This is the best way to be absolutely certain about which axis you're referencing. You can specify these in two forms:
  * Using the code's identifier from <https://github.com/holoplot/go-evdev/blob/master/codes.go>. e.g., `ABS_X`, `REL_WHEEL`, `BTN_TRIGGER`.
  * Alternately, you can omit the `ABS_` type prefix, and Joyful will automatically add it from context. So for a button input, you can simply specify `button: trigger` instead of `BTN_TRIGGER`.
* You can use the hexadecimal value of the code directly, via `"0x<hex value>"`. This can be useful if you want to force a specific numeric value that isn't represented by a Linux event code directly. Note however that not all output codes will work, especially in Windows games. Therefore, this option is most useful with input configurations. **Note: You must use quotation marks around the hex value to prevent the yaml parser from automatically converting it to decimal.**
* For buttons, you can specify them with the above methods, or use an integer index, as in `button: 3`. There are 74 buttons available, and the first button is button number `0`. As a result, valid values are 0-73. Note that buttons 12-14 and buttons 55-73 may not work in all Linux-native games.

For input, you can figure out what event codes your device is emitting by running the Linux utility `evtest`. `evtest` works well with `grep`, so if you just want to see button inputs, you can do:

```
evtest | grep BTN_
```

### Axis Deadzones

**NOTE: For most axis mappings, you probably don't want to specify a deadzone!** Use deadzone configurations in your target game instead. Joyful-configured deadzones are intended to be used in conjunction with the `axis-to-button` and `axis-to-relaxis` input types, or when splitting an axis into multiple outputs. Using them with standard `axis` mappings may result in a loss of fidelity and "stuck" inputs.

There are three ways to specify deadzones:

* Define `deadzone_start` and `deadzone_end` to explicitly set the deadzone bounds.
* Define `deadzone_center` and `deadzone_size`; this will create a deadzone of the indicated size centered at the given axis position.
* Define `deadzone_center` and `deadzone_size_percent` to use a percentage of the total axis size.

See <examples/ruletypes.yml> for usage examples.

## Modes

Modes are optional, and also have the simplest configuration. To define modes, add this to your configuration:

```
modes:
  - mode1
  - mode2
  - mode3
```

The first mode that Joyful reads will be the mode that Joyful starts up in. For that reason, it is recommended to define all your modes in the same file.

Once modes are defined, each rule may specify a `modes` parameter. That rule will only be processed if a matching mode is active. If a rule omits the `modes` parameter, it will be processed in all modes.

For example:

```
rules:
  - name: Test Rule 1 # This rule will be used when we are in mode1 or mode2
    modes:
      - mode1
      - mode2
    # define the rest of the rule here...
  - name: Test Rule 2 # This rule will be used when we are in mode3
    modes:
      - mode3
    # define the rest of the rule here...
```
