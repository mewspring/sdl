// Package render provides support for 2D accelerated rendering. At least one
// active rendering context is required for texture creation and drawing
// operations. A call to window.Open will implicitly add its renderer to the
// list of active renderers, and a subsequent call to win.Close will remove said
// renderer.
package render

import (
	"unsafe"
)

// active keeps track of all active renderers.
var active = make(map[unsafe.Pointer]bool)

// Add adds the renderer to the list of active renderers.
func Add(renderer unsafe.Pointer) {
	if !active[renderer] {
		active[renderer] = true
	}
}

// Del deletes the renderer from the list of active renderers.
func Del(renderer unsafe.Pointer) {
	delete(active, renderer)
}
