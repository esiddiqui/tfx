package progress

import (
	"fmt"
	"time"

	"github.com/esiddiqui/tfx/cursor"
)

const (
	ProgressWaiterTemplate1 string = "%v%% - %v %v"
)

// WaiterStatus encapsulates the status (% complete) &
type WaiterStatus struct {
	ProgressPercent int    // should be between 0-100, value below 0 is considered 0, and over 100 is considered 100
	Message         string // message to display for each progress update
	Data            any    // final data to return
	Error           error
}

// Status returns a WaiterStatus with supplied percent & message
func Status(pr int, msg string) WaiterStatus {
	return WaiterStatus{
		ProgressPercent: pr,
		Message:         msg,
		Data:            nil,
		Error:           nil,
	}
}

// Statusf returns a WaiterStatus with supplied percent & formatted msg
func Statusf(pr int, msg string, params ...any) WaiterStatus {
	return WaiterStatus{
		ProgressPercent: pr,
		Message:         fmt.Sprintf(msg, params...),
		Data:            nil,
		Error:           nil,
	}
}

// Return returns a WaiterStatus with 100% progress, a message & return data
func Done(data any) WaiterStatus {
	return WaiterStatus{
		ProgressPercent: 100,
		Data:            data,
		Error:           nil,
	}
}

// Err returns a WaiterStatus with 100% progress, a message & return data
func Err(err error) WaiterStatus {
	return WaiterStatus{
		ProgressPercent: 100,
		Data:            nil,
		Error:           err,
	}
}

type ComplexWaiterFnc func(chan WaiterStatus)

// ComplexWaiter displays a formatted progress message, i.e. %-complete & a message on the screen until the
// supplied worker aysnchronously completes the job. The job must supply status using a channel of
// type `WaiterStatus` which is used to update the message
//
// The supplied template is used to fill out the progress % value & the messages.
type ComplexWaiter struct {
	SimpleWaiter
	template string
}

// NewComplexWaiter creates & returns a new Progress Waiter
func NewComplexWaiter() ComplexWaiter {
	return ComplexWaiter{
		SimpleWaiter: NewSimpleWaiter(),
		template:     ProgressWaiterTemplate1,
	}
}

// WaitFor calls the supplied CompletxWaiterFnc &  waits to receive status on the
// shared channel. It displays a simple waiter animation or updated status percent
// complete & message
func (w ComplexWaiter) WaitFor(fn ComplexWaiterFnc) (any, error) {

	var idx int
	duration := time.Duration(1000/w.fps) * time.Millisecond

	// build a ticker that ticks every xms based on
	ticker := time.After(duration)

	rc := make(chan WaiterStatus) // make a quick channel

	// start the waiter func as a go routing
	cursor.Off()

	var lastPercent int
	var lastMsg string

	go fn(rc)

	for {

		select {
		case status := <-rc:

			percent := 100
			if status.ProgressPercent >= 100 {
				cursor.ClearToStartOfLine() // clear this ln
				cursor.Col(1)               // move cursor to beginning of ln
				cursor.On()
				return status.Data, nil
			}

			percent = status.ProgressPercent
			cursor.ClearToStartOfLine() // clear this ln
			cursor.Col(1)               // move cursor to beginning of ln
			fmt.Printf(w.template, percent, status.Message, "")
			lastPercent = percent
			lastMsg = status.Message

		case <-ticker:
			cursor.Col(1) // move cursor to beginning of ln
			// fmt.Print(w.Frames[idx])        // paint the next frame
			fmt.Printf(w.template, lastPercent, lastMsg, w.frames[idx])
			idx = (idx + 1) % len(w.frames) // move idx
			ticker = time.After(duration)   // set timer
		}
	}

}
