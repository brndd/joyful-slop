package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/holoplot/go-evdev"
)

func parseCode(code, prefix string) (evdev.EvCode, error) {
	code = strings.ToUpper(code)

	var codeLookup map[string]evdev.EvCode

	switch prefix {
	case CodePrefixButton:
		codeLookup = evdev.KEYFromString
	case CodePrefixAxis:
		codeLookup = evdev.ABSFromString
	case CodePrefixRelaxis:
		codeLookup = evdev.RELFromString
	default:
		return 0, fmt.Errorf("invalid EvCode prefix '%s'", prefix)
	}

	switch {
	case strings.HasPrefix(code, prefix+"_"):
		eventCode, ok := codeLookup[code]
		if !ok {
			return 0, fmt.Errorf("invalid keycode specification '%s'", code)
		}

		return eventCode, nil

	case strings.HasPrefix(code, "0X"):
		codeInt, err := strconv.ParseUint(code[2:], 16, 0)
		if err != nil {
			return 0, err
		}
		return evdev.EvCode(codeInt), nil

	case prefix == CodePrefixButton && !hasError(strconv.Atoi(code)):
		index, err := strconv.Atoi(code)
		if err != nil {
			return 0, err
		}

		if index >= len(ButtonFromIndex) {
			return 0, fmt.Errorf("button index '%d' out of bounds", index)
		}

		return ButtonFromIndex[index], nil

	default:
		eventCode, ok := codeLookup[prefix+"_"+code]
		if !ok {
			return 0, fmt.Errorf("invalid keycode specification '%s'", code)
		}
		return eventCode, nil
	}
}
