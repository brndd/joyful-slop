package virtualdevice

import "github.com/holoplot/go-evdev"

type VirtualDevice interface {
	WriteOne(*evdev.InputEvent) error
}
