// TODO(u): install a default event filter?
// TODO(u): implement a SetEventFilter function?

package win

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
	"image"
	"unicode/utf8"

	"github.com/mewmew/we"
)

// PollEvent returns a pending event from the event queue or nil if the queue
// was empty.
//
// Note: PollEvent must be called from the same thread that created the window.
func (_ *sdlWindow) PollEvent() (event we.Event) {
	e := new(C.SDL_Event)
	// Poll the event queue until we locate a non-nil event or the queue is
	// empty.
	for {
		if C.SDL_PollEvent(e) != 1 {
			// Return nil if the event queue is empty.
			return nil
		}
		event = goEvent(e)
		if event != nil {
			return event
		}
	}
}

// goEvent returns the corresponding Go event for the provided SDL_Event or nil
// if no such Go event exists.
func goEvent(cEvent *C.SDL_Event) (event we.Event) {
	typ := C.getEventType(cEvent)
	switch typ {
	// Close events.
	case C.SDL_QUIT:
		return we.Close{}

	// Window events.
	case C.SDL_WINDOWEVENT:
		e := C.getWindowEvent(cEvent)
		switch e.event {
		case C.SDL_WINDOWEVENT_RESIZED:
			event = we.Resize{
				Width:  int(e.data1),
				Height: int(e.data2),
			}
			return event
		case C.SDL_WINDOWEVENT_ENTER:
			return we.MouseEnter(true)
		case C.SDL_WINDOWEVENT_LEAVE:
			return we.MouseEnter(false)
		}

	// Keyboard events.
	case C.SDL_KEYDOWN:
		e := C.getKeyboardEvent(cEvent)
		if e.repeat == 1 {
			event = we.KeyRepeat{
				Key: goKey(e.keysym.sym),
				Mod: goMod(C.SDL_Keymod(e.keysym.mod)),
			}
			return event
		}
		event = we.KeyPress{
			Key: goKey(e.keysym.sym),
			Mod: goMod(C.SDL_Keymod(e.keysym.mod)),
		}
		return event
	case C.SDL_KEYUP:
		e := C.getKeyboardEvent(cEvent)
		event = we.KeyRelease{
			Key: goKey(e.keysym.sym),
			Mod: goMod(C.SDL_Keymod(e.keysym.mod)),
		}
		return event
	case C.SDL_TEXTINPUT:
		e := C.getTextInputEvent(cEvent)
		text := C.GoString(&e.text[0])
		r, _ := utf8.DecodeRuneInString(text)
		event = we.KeyRune(r)
		return event

	// Mouse events.
	case C.SDL_MOUSEMOTION:
		e := C.getMouseMotionEvent(cEvent)
		if e.state != 0 {
			event = we.MouseDrag{
				Point:  image.Pt(int(e.x), int(e.y)),
				From:   image.Pt(int(e.x-e.xrel), int(e.y-e.yrel)),
				Button: goButtonFromState(e.state),
				Mod:    getMod(),
			}
			return event
		}
		event = we.MouseMove{
			Point: image.Pt(int(e.x), int(e.y)),
			From:  image.Pt(int(e.x-e.xrel), int(e.y-e.yrel)),
		}
		return event
	case C.SDL_MOUSEBUTTONDOWN:
		e := C.getMouseButtonEvent(cEvent)
		event = we.MousePress{
			Point:  image.Pt(int(e.x), int(e.y)),
			Button: goButton(e.button),
			Mod:    getMod(),
		}
		return event
	case C.SDL_MOUSEBUTTONUP:
		e := C.getMouseButtonEvent(cEvent)
		event = we.MouseRelease{
			Point:  image.Pt(int(e.x), int(e.y)),
			Button: goButton(e.button),
			Mod:    getMod(),
		}
		return event
	case C.SDL_MOUSEWHEEL:
		e := C.getMouseWheelEvent(cEvent)
		switch {
		case e.x != 0:
			event = we.ScrollX{
				Off: int(e.x),
				Mod: getMod(),
			}
			return event
		case e.y != 0:
			event = we.ScrollY{
				Off: int(e.y),
				Mod: getMod(),
			}
			return event
		}
	}

	// Ignore event.
	return nil
}
