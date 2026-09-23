# Converted Joystick Gremlin Profile

This directory is generated deterministically from `joystick_gremlin_profile.xml` by `go run ./cmd/gremlin-convert`.

Regenerate and run it from the repository root:

```sh
go run ./cmd/gremlin-convert --input joystick_gremlin_profile.xml --output joystick_gremlin_profile
joyful --config ./joystick_gremlin_profile
```

The profile uses two physical VKB Gladiator devices and separate left and right virtual joysticks. Mode inheritance has been flattened and equivalent resolved behavior is grouped into rules with explicit mode lists.

The converter translates VKB axis polarity from Gremlin/DirectInput to Linux evdev and accounts for the different `ABS_RX` / `ABS_RY` mini-stick ordering.

Modes, in startup order: SCM Mode, Auxiliary Mode, Modifier, Nav Mode.

On each controller, the mini-stick maps to virtual `ABS_Z` / `ABS_RX` (KDE axes 3/4), twist maps to `ABS_RY` (axis 5), and the remaining throttle/lever axis maps to `ABS_RZ` (axis 6). The original right-controller virtual mouse mapping and its acceleration curve are intentionally omitted.

Gremlin vJoy buttons 124 through 128 exceed Joyful's output range and are remapped to virtual buttons 60 through 64 (Joyful indices 59 through 63). Rebind those five controls in the target game.

The freelook recenter and Light Amplification macros are intentionally omitted. See `conversion-report.txt` for all omissions, remaps, and timing nuances.

Generated rules: 139.
