// Package window handles window creation, drawing and events. It uses a small
// subset of the features provided by the SDL library version 2.0 [1].
//
// [1]: http://www.libsdl.org/
package window

import (
	"image"
	"image/color"

	"github.com/mewmew/wandi"
)

// A Window represents a graphical window capable of handling draw operations
// and window events. It implements the wandi.Window interface.
type Window struct {
}

// Open opens a new window of the specified dimensions.
//
// Note: The Close method of the window must be called when finished using it.
func Open(width, height int) (win Window, err error) {
	panic("not yet implemented")
}

// Close closes the window.
func (win Window) Close() {
	panic("not yet implemented")
}

// SetTitle sets the title of the window.
func (win Window) SetTitle(title string) {
	panic("not yet implemented")
}

// ShowCursor displays or hides the mouse cursor depending on the value of
// visible. It is visible by default.
func (win Window) ShowCursor(visible bool) {
	panic("not yet implemented")
}

// Width returns the width of the window.
func (win Window) Width() int {
	panic("not yet implemented")
}

// Height returns the height of the window.
func (win Window) Height() int {
	panic("not yet implemented")
}

// Draw draws the entire src image onto the window starting at the destination
// point dp.
func (win Window) Draw(dp image.Point, src wandi.Image) (err error) {
	panic("not yet implemented")
}

// DrawRect draws a subset of the src image, as defined by the source rectangle
// sr, onto the window starting at the destination point dp.
func (win Window) DrawRect(dp image.Point, src wandi.Image, sr image.Rectangle) (err error) {
	panic("not yet implemented")
}

// Fill fills the entire window with the provided color.
func (win Window) Fill(c color.Color) {
	panic("not yet implemented")
}

// Display displays what has been rendered so far to the window.
func (win Window) Display() {
	panic("not yet implemented")
}
