# Joyful Configuration

Configuration is divided into three sections: `devices`, `modes`, and `rules`. Each yaml file can have any number of these sections; joyful will combine the configuration from all files at runtime.

### Device configuration

Each entry in `devices` must have a couple of parameters:

* `name` - This is an identifier that your rules will use to refer to the device. It is recommended to avoid spaces or special characters.
* `type` - Should be `physical` for an input device, and `virtual` for an output device.

`physical` devices must additionally define these parameters:

* `device_name` - The name of the device as reported by the included `evlist` command. If your device name ends with a space, use quotation marks (`""`) around the name.

`virtual` devices must additionally define these parameters:

* `buttons` - a number between 0 and 74. Linux may not recognize buttons greater than 56.
* `axes` - a number between 0 and 8.

Virtual devices can also define a `relative_axes` parameter; this must be a list of `REL_` event keycodes, and can be useful for a simulated mouse device. Some environments will only register mouse events if the device *only* supports mouse-like events, so it can be useful to isolate your `relative_axes` to their own virtual device.

### Rules configuration

All `rules` must have a `type` parameter. Valid values for this parameter are:

* `button` - a single button mapping
* `button-combo` - multiple input buttons mapped to a single output. The output event will trigger when all the input conditions are met.
* `button-latched` - a single button mapped to a single output, but each time the input is pressed, the output will toggle.
* `axis` - a simple axis mapping
* `axis-to-button` - causes an axis input to produce a button output. This can be repeated with variable speed proportional to the axis' input value
* `axis-to-relaxis` - like axis-to-button, but produces a "relative axis" output value. This is useful for simulating mouse scrollwheel and movement events.

Configuration options for each rule type vary. See <examples/ruletypes.yml> for an example of each type with all options specified.

### Keycodes

Keycodes are the values that identify buttons and axes. There are several ways to configure keycodes. All of them are case-insensitive.

Ways to specify keycodes are:

* Using evdev's Keycodes. This is the best way to be absolutely certain about which axis you're referencing. You can specify these keycodes in two forms:
  * Using the code's identifier from <https://github.com/holoplot/go-evdev/blob/master/codes.go>. e.g., `ABS_X`, `REL_WHEEL`, `BTN_TRIGGER`.
  * Alternately, you can omit the `ABS_` type prefix, and Joyful will automatically add it from context. So for a button input, you can simply specify `button: trigger` instead of `BTN_TRIGGER`.
* You can use the hexadecimal value of the keycode directly, via `"0x<hex value>"`. This can be useful if you want to force a specific numeric value that isn't represented by a Linux keycode directly. Note however that not all keycodes will work. Only the first 8 axes are available, and see <internal/config/variables.go> for a list of valid button outputs. This is most useful with input configurations. **Note: You must use quotation marks around the hex value to prevent the yaml parser from automatically converting it to decimal.**
* For buttons, you can specify the button number, as in `button: 3`. There are 74 buttons available, and the first button is button number `0`. As a result, valid values are 0-73. Note that buttons 12-14 and buttons 55-73 may not work in all Linux-native games.

For input, you can figure out what keycodes your device is emitting by running the Linux utility `evtest`. `evtest` works well with `grep`, so if you just want to see button inputs, you can do:

```
evtest | grep BTN_
```

The authors of this tool recognize that this is currently a pain in the ass. Easier ways to represent keycodes (as well as outputting additional keycodes) is planned for the future.

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
