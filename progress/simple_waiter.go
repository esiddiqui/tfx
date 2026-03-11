package progress

import (
	"fmt"
	"time"

	"github.com/esiddiqui/tfx/cursor"
)

var (
	WaiterFrames1 []string = []string{"/", "-", "\\", "|", "-", "\\", "|"}
	WaiterFrames2 []string = []string{".", "..", "...", "....", " ...", "  ..", "   .", "", "   .", "  ..", " ...", "....", "... ", "..  ", ".  ", ""}
	WaiterFrames3 []string = []string{"¯", "¯\\", "¯\\_", "¯\\_(", "¯\\_(ツ", "¯\\_(ツ)", "¯\\_(ツ)_", "¯\\_(ツ)_/", "¯\\_(ツ)_/¯", "          "}
)

// a simple wrapper func type alias that takes in a `chan any` to indicate the worker finished processing & produced a result..
type SimpleWaiterFnc func(chan any)

// a simple post processor func type alias, which takes an argument of type any to process
type SimpleWaiterPostFnc func(any)

// SimpleWaiter displays a progress waiter text animation on the screen until the
// supplied worker aysnchronously completes the job & return a value over `channel any`
type SimpleWaiter struct {
	fps    int
	frames []string
}

type SimpleWaiterOptsFnc func(w *SimpleWaiter)

func WithFps(fps int) SimpleWaiterOptsFnc {
	return func(w *SimpleWaiter) {
		w.fps = fps
	}
}

func WithWaiterFrames(frames []string) SimpleWaiterOptsFnc {
	return func(w *SimpleWaiter) {
		w.frames = frames
	}
}

// NewSimpleWaiter creates & returns a SimpleWaiter
func NewSimpleWaiter(opts ...SimpleWaiterOptsFnc) SimpleWaiter {
	w := &SimpleWaiter{
		fps:    20,
		frames: WaiterFrames1,
	}

	for _, opt := range opts {
		opt(w)
	}

	return *w
}

// WaitFor calls the supplied SimpleWaiterFnc & starts a waiter animation using the supplied
// frames. When the `fn` completes it sends the return data over the channel
func (w SimpleWaiter) WaitFor(fn SimpleWaiterFnc) (any, error) {

	cursor.Off()

	var idx int
	duration := time.Duration(1000/w.fps) * time.Millisecond

	// build a ticker that ticks every xms based on
	ticker := time.After(duration)
	rc := make(chan any) // make a quick channel

	go fn(rc)

	for {
		select {
		case val := <-rc:
			cursor.ClearToStartOfLine() // clear this ln
			cursor.Col(1)               // move cursor to beginning of ln
			cursor.On()                 // set cursor visible
			return val, nil

		case <-ticker:
			cursor.Col(1)                   // move cursor to beginning of ln
			fmt.Print(w.frames[idx])        // paint the next frame
			idx = (idx + 1) % len(w.frames) // move idx
			ticker = time.After(duration)   // set timer
		}
	}
}
