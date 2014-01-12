package win

// #cgo pkg-config: sdl2
// #include <SDL2/SDL.h>
import "C"

import (
	"github.com/mewmew/we"
)

// getMod returns the currently active keyboard modifiers.
func getMod() (mod we.Mod) {
	return goMod(C.SDL_GetModState())
}

// goMod returns the corresponding we.Mod for the provided SDL_Keymod.
func goMod(m C.SDL_Keymod) (mod we.Mod) {
	if m&C.KMOD_LSHIFT != 0 || m&C.KMOD_RSHIFT != 0 {
		mod |= we.ModShift
	}
	if m&C.KMOD_LCTRL != 0 || m&C.KMOD_RCTRL != 0 {
		mod |= we.ModControl
	}
	if m&C.KMOD_LALT != 0 || m&C.KMOD_RALT != 0 {
		mod |= we.ModAlt
	}
	if m&C.KMOD_LGUI != 0 || m&C.KMOD_RGUI != 0 {
		mod |= we.ModSuper
	}
	return mod
}

// goKey returns the corresponding we.Key for the provided SDL_Keycode.
func goKey(keycode C.SDL_Keycode) (key we.Key) {
	switch keycode {
	// Printable keys.
	case C.SDLK_SPACE:
		return we.KeySpace
	case C.SDLK_QUOTE:
		return we.KeyApostrophe
	case C.SDLK_COMMA:
		return we.KeyComma
	case C.SDLK_MINUS:
		return we.KeyMinus
	case C.SDLK_PERIOD:
		return we.KeyPeriod
	case C.SDLK_SLASH:
		return we.KeySlash
	case C.SDLK_0:
		return we.Key0
	case C.SDLK_1:
		return we.Key1
	case C.SDLK_2:
		return we.Key2
	case C.SDLK_3:
		return we.Key3
	case C.SDLK_4:
		return we.Key4
	case C.SDLK_5:
		return we.Key5
	case C.SDLK_6:
		return we.Key6
	case C.SDLK_7:
		return we.Key7
	case C.SDLK_8:
		return we.Key8
	case C.SDLK_9:
		return we.Key9
	case C.SDLK_SEMICOLON:
		return we.KeySemicolon
	case C.SDLK_EQUALS:
		return we.KeyEqual
	case C.SDLK_a:
		return we.KeyA
	case C.SDLK_b:
		return we.KeyB
	case C.SDLK_c:
		return we.KeyC
	case C.SDLK_d:
		return we.KeyD
	case C.SDLK_e:
		return we.KeyE
	case C.SDLK_f:
		return we.KeyF
	case C.SDLK_g:
		return we.KeyG
	case C.SDLK_h:
		return we.KeyH
	case C.SDLK_i:
		return we.KeyI
	case C.SDLK_j:
		return we.KeyJ
	case C.SDLK_k:
		return we.KeyK
	case C.SDLK_l:
		return we.KeyL
	case C.SDLK_m:
		return we.KeyM
	case C.SDLK_n:
		return we.KeyN
	case C.SDLK_o:
		return we.KeyO
	case C.SDLK_p:
		return we.KeyP
	case C.SDLK_q:
		return we.KeyQ
	case C.SDLK_r:
		return we.KeyR
	case C.SDLK_s:
		return we.KeyS
	case C.SDLK_t:
		return we.KeyT
	case C.SDLK_u:
		return we.KeyU
	case C.SDLK_v:
		return we.KeyV
	case C.SDLK_w:
		return we.KeyW
	case C.SDLK_x:
		return we.KeyX
	case C.SDLK_y:
		return we.KeyY
	case C.SDLK_z:
		return we.KeyZ
	case C.SDLK_LEFTBRACKET:
		return we.KeyLeftBracket
	case C.SDLK_BACKSLASH:
		return we.KeyBackslash
	case C.SDLK_RIGHTBRACKET:
		return we.KeyRightBracket
	case C.SDLK_BACKQUOTE:
		return we.KeyGraveAccent

	// Function keys.
	case C.SDLK_ESCAPE:
		return we.KeyEscape
	case C.SDLK_RETURN:
		return we.KeyEnter
	case C.SDLK_TAB:
		return we.KeyTab
	case C.SDLK_BACKSPACE:
		return we.KeyBackspace
	case C.SDLK_INSERT:
		return we.KeyInsert
	case C.SDLK_DELETE:
		return we.KeyDelete
	case C.SDLK_RIGHT:
		return we.KeyRight
	case C.SDLK_LEFT:
		return we.KeyLeft
	case C.SDLK_DOWN:
		return we.KeyDown
	case C.SDLK_UP:
		return we.KeyUp
	case C.SDLK_PAGEUP:
		return we.KeyPageUp
	case C.SDLK_PAGEDOWN:
		return we.KeyPageDown
	case C.SDLK_HOME:
		return we.KeyHome
	case C.SDLK_END:
		return we.KeyEnd
	case C.SDLK_CAPSLOCK:
		return we.KeyCapsLock
	case C.SDLK_SCROLLLOCK:
		return we.KeyScrollLock
	case C.SDLK_NUMLOCKCLEAR:
		return we.KeyNumLock
	case C.SDLK_PRINTSCREEN:
		return we.KeyPrintScreen
	case C.SDLK_PAUSE:
		return we.KeyPause
	case C.SDLK_F1:
		return we.KeyF1
	case C.SDLK_F2:
		return we.KeyF2
	case C.SDLK_F3:
		return we.KeyF3
	case C.SDLK_F4:
		return we.KeyF4
	case C.SDLK_F5:
		return we.KeyF5
	case C.SDLK_F6:
		return we.KeyF6
	case C.SDLK_F7:
		return we.KeyF7
	case C.SDLK_F8:
		return we.KeyF8
	case C.SDLK_F9:
		return we.KeyF9
	case C.SDLK_F10:
		return we.KeyF10
	case C.SDLK_F11:
		return we.KeyF11
	case C.SDLK_F12:
		return we.KeyF12
	case C.SDLK_F13:
		return we.KeyF13
	case C.SDLK_F14:
		return we.KeyF14
	case C.SDLK_F15:
		return we.KeyF15
	case C.SDLK_F16:
		return we.KeyF16
	case C.SDLK_F17:
		return we.KeyF17
	case C.SDLK_F18:
		return we.KeyF18
	case C.SDLK_F19:
		return we.KeyF19
	case C.SDLK_F20:
		return we.KeyF20
	case C.SDLK_F21:
		return we.KeyF21
	case C.SDLK_F22:
		return we.KeyF22
	case C.SDLK_F23:
		return we.KeyF23
	case C.SDLK_F24:
		return we.KeyF24
	case C.SDLK_KP_0:
		return we.KeyKp0
	case C.SDLK_KP_1:
		return we.KeyKp1
	case C.SDLK_KP_2:
		return we.KeyKp2
	case C.SDLK_KP_3:
		return we.KeyKp3
	case C.SDLK_KP_4:
		return we.KeyKp4
	case C.SDLK_KP_5:
		return we.KeyKp5
	case C.SDLK_KP_6:
		return we.KeyKp6
	case C.SDLK_KP_7:
		return we.KeyKp7
	case C.SDLK_KP_8:
		return we.KeyKp8
	case C.SDLK_KP_9:
		return we.KeyKp9
	case C.SDLK_KP_PERIOD:
		return we.KeyKpDecimal
	case C.SDLK_KP_DIVIDE:
		return we.KeyKpDivide
	case C.SDLK_KP_MULTIPLY:
		return we.KeyKpMultiply
	case C.SDLK_KP_MINUS:
		return we.KeyKpSubtract
	case C.SDLK_KP_PLUS:
		return we.KeyKpAdd
	case C.SDLK_KP_ENTER:
		return we.KeyKpEnter
	case C.SDLK_KP_EQUALS:
		return we.KeyKpEqual
	case C.SDLK_LSHIFT:
		return we.KeyLeftShift
	case C.SDLK_LCTRL:
		return we.KeyLeftControl
	case C.SDLK_LALT:
		return we.KeyLeftAlt
	case C.SDLK_LGUI:
		return we.KeyLeftSuper
	case C.SDLK_RSHIFT:
		return we.KeyRightShift
	case C.SDLK_RCTRL:
		return we.KeyRightControl
	case C.SDLK_RALT:
		return we.KeyRightAlt
	case C.SDLK_RGUI:
		return we.KeyRightSuper
	case C.SDLK_MENU:
		return we.KeyMenu
	}

	// Unknown key.
	return 0
}
