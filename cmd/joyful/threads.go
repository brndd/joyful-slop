package main

import (
	"bufio"
	"os"
	"sync"
	"time"

	"git.annabunches.net/annabunches/joyful/internal/logger"
	"git.annabunches.net/annabunches/joyful/internal/mappingrules"
	"github.com/holoplot/go-evdev"
)

const (
	TimerCheckIntervalMs  = 1
	DeviceCheckIntervalMs = 1
)

func eventWatcher(
	device *evdev.InputDevice,
	channel chan<- ChannelEvent,
	done chan bool,
	wg *sync.WaitGroup) {

	defer wg.Done()

	for {
		select {
		case cancel := <-done:
			if cancel {
				done <- true
				return
			}
		default:
		}

		event, err := device.ReadOne()
		if err != nil {
			logger.LogError(err, "Error while reading event. Disconnecting device.")
			return
		}
		channel <- ChannelEvent{Device: device, Event: event, Type: ChannelEventInput}

		if event.Type == evdev.EV_SYN {
			time.Sleep(DeviceCheckIntervalMs * time.Millisecond)
		}
	}
}

func timerWatcher(
	rule mappingrules.TimedEventEmitter,
	channel chan<- ChannelEvent,
	done chan bool,
	wg *sync.WaitGroup) {

	defer wg.Done()

	for {
		select {
		case cancel := <-done:
			if cancel {
				done <- true
				return
			}
		default:
		}

		event := rule.TimerEvent()
		if event != nil {
			channel <- ChannelEvent{
				Device: rule.GetOutputDevice(),
				Event:  event,
				Type:   ChannelEventTimer,
			}
		}
		time.Sleep(TimerCheckIntervalMs * time.Millisecond)
	}
}

// consoleWatcher reads input from stdin, and on receiving anything
func consoleWatcher(channel chan<- ChannelEvent, wg *sync.WaitGroup) {
	defer wg.Done()
	stdin := bufio.NewReader(os.Stdin)
	for {
		_, err := stdin.ReadString('\n')
		if err != nil {
			logger.LogErrorf(err, "Error in console input thread")
			continue
		}

		channel <- ChannelEvent{
			Type: ChannelEventReload,
		}
		return
	}
}
