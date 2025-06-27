// Functions for cleaning up stale virtual devices

package virtualdevice

import (
	"fmt"
	"strings"

	"github.com/holoplot/go-evdev"
)

func CleanupStaleVirtualDevices() {
	devices, err := evdev.ListDevicePaths()
	if err != nil {
		fmt.Printf("Couldn't list devices while running cleanup: %s\n", err.Error())
		return
	}

	for _, devicePath := range devices {
		if strings.HasPrefix(devicePath.Name, "joyful-joystick") {
			device, err := evdev.Open(devicePath.Path)
			if err != nil {
				fmt.Printf("Failed to open existing joyful device at '%s': %s\n", devicePath.Path, err.Error())
				continue
			}

			err = evdev.DestroyDevice(device)
			if err != nil {
				fmt.Printf("Failed to destroy existing joyful device '%s' at '%s': %s\n", devicePath.Name, devicePath.Path, err.Error())
			} else {
				fmt.Printf("Destroyed stale joyful device '%s'\n", devicePath.Path)
			}
		}
	}
}
