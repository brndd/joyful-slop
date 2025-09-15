package mappingrules

import (
	"errors"
	"fmt"

	"git.annabunches.net/annabunches/joyful/internal/configparser"
	"github.com/holoplot/go-evdev"
)

// TODO: need tests for multiple deadzones
// TODO: need tests for emitting deadzones

type Deadzone struct {
	Start     int32
	End       int32
	Size      int32
	Emit      bool
	EmitValue int32
}

// DeadzoneState indicates whether a value is in a Deadzone and, if it is, whether the deadzone
// should emit an event
type DeadzoneState int

const (
	// DeadzoneClear indicates the value is *not* in the deadzone.
	DeadzoneClear DeadzoneState = iota
	DeadzoneEmit
	DeadzoneNoEmit
)

// calculateDeadzones produces the deadzone start and end values in absolute terms
func NewDeadzoneFromConfig(dzConfig configparser.DeadzoneConfig, device Device, axis evdev.EvCode) (Deadzone, error) {
	dz := Deadzone{}
	dz.Emit = dzConfig.Emit
	dz.EmitValue = dzConfig.Value

	var min, max int32
	absInfoMap, err := device.AbsInfos()

	if err != nil {
		return dz, err
	} else {
		absInfo := absInfoMap[axis]
		min = absInfo.Minimum
		max = absInfo.Maximum
	}

	if dzConfig.Start != 0 || dzConfig.End != 0 {
		dz.Start = Clamp(dzConfig.Start, min, max)
		dz.End = Clamp(dzConfig.End, min, max)
		if dz.Start > dz.End {
			return dz, errors.New("deadzone end must be greater than deadzone start")
		}
	} else {
		center := Clamp(dzConfig.Center, min, max)
		var deadzoneSize int32

		switch {
		case dzConfig.Size != 0:
			deadzoneSize = dzConfig.Size
		case dzConfig.SizePercent != 0:
			deadzoneSize = (max - min) / dzConfig.SizePercent
		default:
			return dz, fmt.Errorf("deadzone configured incorrectly; must define start and end or center and size")
		}

		dz.Start = center - deadzoneSize/2
		dz.End = center + deadzoneSize/2
		dz.Start, dz.End = clampAndShift(dz.Start, dz.End, min, max)
	}

	dz.Size = dz.End - dz.Start
	return dz, nil
}

func CalculateDeadzoneSize(dzs []Deadzone) int32 {
	var size int32

	for _, dz := range dzs {
		size += dz.Size
	}

	return size
}

// Match checks whether the target value is inside the deadzone.
// It returns a DeadzoneState enum and possibly an int32.
func (dz Deadzone) Match(value int32) (DeadzoneState, int32) {
	if value < dz.Start || value > dz.End {
		return DeadzoneClear, value
	}
	if dz.Emit {
		return DeadzoneEmit, dz.EmitValue
	}
	return DeadzoneNoEmit, value
}
