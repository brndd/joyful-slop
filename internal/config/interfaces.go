package config

import "github.com/holoplot/go-evdev"

type Device interface {
	AbsInfos() (map[evdev.EvCode]evdev.AbsInfo, error)
}
