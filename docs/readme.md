# Joyful Configuration

Configuration is divided into three sections: `devices`, `modes`, and `rules`. Each yaml file can have any number of these sections; joyful will combine the configuration from all files at runtime.

## Device configuration

Each entry in `devices` must have these parameters:

* `name` - This is an identifier that your rules will use to refer to the device. It is recommended to avoid spaces or special characters.
* `type` - 'physical' for an input device, 'virtual' for an output device.

### Physical Devices

`physical` devices have these additional parameters:

* `device_name` - The name of the device as reported by the included `evinfo` command. If your device name ends with a space, use quotation marks (`""`) around the name.
* `device_path` - If you have multiple devices that report the same name, you can use `device_path` instead of `device_name`. Setting this will cause the device to be opened directly via the device file. 
    * It is recommended to use the `by-path` symlinks, e.g., `/dev/input/by-path/pci-0000:0d:00.0-usbv2-0:3:1.0-event-joystick`.
    * Note that this method may be slightly unreliable since these identifiers may change if they are plugged into different USB ports or in the rare case that the USB topology changes (e.g., you add a new USB hub).
    * On the other hand, this method causes the device to be opened considerably faster, lowering Joyful's startup time substantially. If this is important to you this method may be preferable.
* `lock` - If set to 'true', the device will be locked for exclusive access. This means that your game will not see any events from the device, so you'll need to make sure you map every button you want to use. Setting this to 'false' might be useful if you're just mapping a few joystick buttons to keyboard buttons. This value defaults to 'true'.

`device_path` is given higher priority than `device_name`; if both are specified, `device_path` will be used.

### Virtual Devices

`virtual` devices have these additional parameters:

* `preset` - Can be 'joystick', 'gamepad', 'mouse', or 'keyboard', and will configure the virtual device to look like and emit an appropriate set of outputs based on the name. For exactly which axes and buttons are defined for each type, see the `Capabilities` values in [internal/config/variables.go](internal/config/variables.go).
* `buttons` or `num_buttons` - Either a list of explicit buttons or a number of buttons to create. (max 74 buttons) Linux-native games may not recognize all buttons created by Joyful.
* `axes` or `num_axes` - An explicit list of `ABS_` axes or a number to create.
* `relative_axes` or `num_relative_axes` - As above, but for `REL_` axes.

A couple of additional notes on virtual devices:

* Users are encouraged to use the `preset` options whenever possible. They have the highest probability of working the way you expect. If you need to output to multiple types of device, the best approach is to create multiple virtual devices.
* For all 3 of the above options, there is a priority order. If you specify a `preset`, it will be used ignoring any other settings. An explicit list will override the corresponding `num_` parameter.
* Some environments/applications are prescriptive about what combinations make sense; for example, they will only register mouse events if the device *only* supports mouse-like events. The `presets` attempt to take this into account. If you are defining capabilities manually and attempt to mix and match button codes, you may also run into this problem.

## Rules configuration

All `rules` must have a `type` parameter. Valid values for this parameter are:

* `button` - a single button mapping
* `button-combo` - multiple input buttons mapped to a single output. The output event will trigger when all the input conditions are met.
* `button-latched` - a single button mapped to a single output, but each time the input is pressed, the output will toggle.
* `axis` - a simple axis mapping
* `axis-combined` - a mapping that combines 2 input axes into a single output axis.
* `axis-to-button` - causes an axis input to produce a button output. This can be repeated with variable speed proportional to the axis' input value
* `axis-to-relaxis` - like axis-to-button, but produces a "relative axis" output value. This is useful for simulating mouse scrollwheel and movement events.
* `hat` - a special type of axis with ternary output. Each joystick hat will typically be 2 hat axes named `ABS_HAT0X` / `ABS_HAT0Y`, where the `0` is an index between 0 - 3. So for a typical hat you would define 2 `hat` rules.

Configuration options for each rule type vary. See [examples/ruletypes.yml](examples/ruletypes.yml) for an example of each type with all options specified.

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

Axis inputs can define a list of deadzones. Each deadzone can be specified a few ways:

* Define `start` and `end` to explicitly set the deadzone bounds.
* Define `center` and `size`; this will create a deadzone of the indicated size centered at the given axis position.
* Define `center` and `size_percent` to use a percentage of the total axis size.

In addition, deadzones can set `emit` to `true` and `emit_value` to a value that should be emitted when inside the deadzone.

**Note**: The `emit_value` is the final output value and should be between -32,768 and 32,767.

See the <examples/> directory for usage examples.

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
