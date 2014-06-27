package window

// #cgo pkg-config: sdl2
// #include <SDL2/SDL.h>
//
// int getEventType(SDL_Event *e) {
//    return e->type;
// }
//
// SDL_WindowEvent * getWindowEvent(SDL_Event *e) {
//    return &e->window;
// }
//
// SDL_KeyboardEvent * getKeyboardEvent(SDL_Event *e) {
//    return &e->key;
// }
//
// SDL_TextInputEvent * getTextInputEvent(SDL_Event *e) {
//    return &e->text;
// }
//
// SDL_MouseMotionEvent * getMouseMotionEvent(SDL_Event *e) {
//    return &e->motion;
// }
//
// SDL_MouseButtonEvent * getMouseButtonEvent(SDL_Event *e) {
//    return &e->button;
// }
//
// SDL_MouseWheelEvent * getMouseWheelEvent(SDL_Event *e) {
//    return &e->wheel;
// }
import "C"

import (
	"container/list"
	"image"
	"unicode/utf8"

	"github.com/mewmew/we"
)

// TODO(u): Protect queue from race conditions.

// queue maps window IDs to their corresponding event queue.
var queue = make(map[uint32]*list.List)

// PollEvent returns a pending event from the event queue or nil if the queue
// was empty. Note that more than one event may be present in the event queue.
//
// Note: PollEvent should only be called from the same thread that created the
// window.
func (win Window) PollEvent() (event we.Event) {
	// Poll events from the event queue until it's empty.
	sdlEvent := new(C.SDL_Event)
	for C.SDL_PollEvent(sdlEvent) != 0 {
		event, winID := weEvent(sdlEvent)
		// Add non-nil events to the event queue of their corresponding window.
		if event == nil {
			continue
		}
		l, ok := queue[winID]
		if !ok {
			l = list.New()
			queue[winID] = l
		}
		l.PushBack(event)
	}

	// Pop the first event from the event queue of the current window.
	curID := uint32(C.SDL_GetWindowID(win.win))
	l, ok := queue[curID]
	if !ok {
		return nil
	}
	e := l.Front()
	if e == nil {
		return nil
	}
	event = l.Remove(e)
	return event
}

// weEvent returns the corresponding we.Event for the provided SDL event or nil
// if no such event exists.
func weEvent(sdlEvent *C.SDL_Event) (event we.Event, winID uint32) {
	typ := C.getEventType(sdlEvent)
	switch typ {
	// Window events.
	case C.SDL_WINDOWEVENT:
		e := C.getWindowEvent(sdlEvent)
		winID = uint32(e.windowID)
		switch e.event {
		case C.SDL_WINDOWEVENT_CLOSE:
			return we.Close{}, winID
		case C.SDL_WINDOWEVENT_RESIZED:
			event = we.Resize{
				Width:  int(e.data1),
				Height: int(e.data2),
			}
			return event, winID
		case C.SDL_WINDOWEVENT_ENTER:
			return we.MouseEnter(true), winID
		case C.SDL_WINDOWEVENT_LEAVE:
			return we.MouseEnter(false), winID
		}

	// Keyboard events.
	case C.SDL_KEYDOWN:
		e := C.getKeyboardEvent(sdlEvent)
		winID = uint32(e.windowID)
		if e.repeat == 1 {
			event = we.KeyRepeat{
				Key: weKey(e.keysym.sym),
				Mod: weMod(C.SDL_Keymod(e.keysym.mod)),
			}
			return event, winID
		}
		event = we.KeyPress{
			Key: weKey(e.keysym.sym),
			Mod: weMod(C.SDL_Keymod(e.keysym.mod)),
		}
		return event, winID
	case C.SDL_KEYUP:
		e := C.getKeyboardEvent(sdlEvent)
		winID = uint32(e.windowID)
		event = we.KeyRelease{
			Key: weKey(e.keysym.sym),
			Mod: weMod(C.SDL_Keymod(e.keysym.mod)),
		}
		return event, winID
	case C.SDL_TEXTINPUT:
		e := C.getTextInputEvent(sdlEvent)
		winID = uint32(e.windowID)
		text := C.GoString(&e.text[0])
		r, _ := utf8.DecodeRuneInString(text)
		event = we.KeyRune(r)
		return event, winID

	// Mouse events.
	case C.SDL_MOUSEMOTION:
		e := C.getMouseMotionEvent(sdlEvent)
		winID = uint32(e.windowID)
		if e.state != 0 {
			event = we.MouseDrag{
				Point:  image.Pt(int(e.x), int(e.y)),
				From:   image.Pt(int(e.x-e.xrel), int(e.y-e.yrel)),
				Button: weButtonFromState(e.state),
				Mod:    getMod(),
			}
			return event, winID
		}
		event = we.MouseMove{
			Point: image.Pt(int(e.x), int(e.y)),
			From:  image.Pt(int(e.x-e.xrel), int(e.y-e.yrel)),
		}
		return event, winID
	case C.SDL_MOUSEBUTTONDOWN:
		e := C.getMouseButtonEvent(sdlEvent)
		winID = uint32(e.windowID)
		event = we.MousePress{
			Point:  image.Pt(int(e.x), int(e.y)),
			Button: weButton(e.button),
			Mod:    getMod(),
		}
		return event, winID
	case C.SDL_MOUSEBUTTONUP:
		e := C.getMouseButtonEvent(sdlEvent)
		winID = uint32(e.windowID)
		event = we.MouseRelease{
			Point:  image.Pt(int(e.x), int(e.y)),
			Button: weButton(e.button),
			Mod:    getMod(),
		}
		return event, winID
	case C.SDL_MOUSEWHEEL:
		e := C.getMouseWheelEvent(sdlEvent)
		winID = uint32(e.windowID)
		switch {
		case e.x != 0:
			event = we.ScrollX{
				Off: int(e.x),
				Mod: getMod(),
			}
			return event, winID
		case e.y != 0:
			event = we.ScrollY{
				Off: int(e.y),
				Mod: getMod(),
			}
			return event, winID
		}

	// Unknown event types.
	default:
		return nil, 0
	}

	// Ignore event.
	return nil, 0
}
