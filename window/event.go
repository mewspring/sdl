package window

import (
	"github.com/mewmew/we"
)

// PollEvent returns a pending event from the event queue or nil if the queue
// was empty. Note that more than one event may be present in the event queue.
//
// Note: PollEvent should only be called from the same thread that created the
// window.
func (win Window) PollEvent() (event we.Event) {
	panic("not yet implemented")
}
