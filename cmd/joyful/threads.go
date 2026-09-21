package main

import (
	"bufio"
	"context"
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
	ctx context.Context,
	wg *sync.WaitGroup) {

	defer wg.Done()

	for {
		event, err := device.ReadOne()
		if err != nil {
			logger.LogError(err, "Error while reading event. Disconnecting device.")
			return
		}

		select {
		case <-ctx.Done():
			return
		default:
			// Proceed
		}

		channel <- ChannelEvent{Device: device, Event: event, Type: ChannelEventInput}

		if event.Type == evdev.EV_SYN {
			time.Sleep(DeviceCheckIntervalMs * time.Millisecond)
		}
	}
}

func timerWatcher(
	rule mappingrules.MainLoopTimedEventEmitter,
	channel chan<- ChannelEvent,
	ctx context.Context,
	wg *sync.WaitGroup) {

	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Proceed
		}

		select {
		case channel <- ChannelEvent{Rule: rule, Type: ChannelEventTimer}:
		case <-ctx.Done():
			return
		}
		time.Sleep(TimerCheckIntervalMs * time.Millisecond)
	}
}

// consoleWatcher reads input from stdin, and on receiving anything,
// closes the current threading context
func consoleWatcher(channel chan<- ChannelEvent) {
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
