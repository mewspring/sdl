package font

import (
	"image/color"
)

// Text represent a graphical text entry with a specific font size, style and
// color. It implements the wandi.Image interface.
type Text struct {
}

// NewText returns a new graphical text entry based on the provided font and any
// optional customization arguments. The initial text, size, style and color of
// the graphical text entry can be customized through string, int, Style and
// color.Color arguments respectively, depending on the type of the argument.
//
// The default font size, style and color of the text is 12, regular (no style)
// and black respectively.
//
// Note: The Free method of the text entry must be called when finished using
// it.
func NewText(f Font, args ...interface{}) (text Text, err error) {
	panic("not yet implemented")
}

// Free frees the text entry.
func (text Text) Free() {
	panic("not yet implemented")
}

// SetText sets the text of the text entry.
func (text Text) SetText(s string) {
	panic("not yet implemented")
}

// SetSize sets the font size, in pixels, of the text.
func (text Text) SetSize(size int) {
	panic("not yet implemented")
}

// Style is a bitfield which represents the style of a text.
type Style uint32

// Text styles.
const (
	// Regular characters (no style).
	Regular Style = iota
	// Bold characters.
	Bold
	// Italic characters.
	Italic
	// Underlined characters.
	Underlined
)

// SetStyle sets the style of the text.
func (text Text) SetStyle(style Style) {
	panic("not yet implemented")
}

// SetColor sets the color of the text.
func (text Text) SetColor(c color.Color) {
	panic("not yet implemented")
}

// Width returns the width of the text entry.
func (text Text) Width() int {
	panic("not yet implemented")
}

// Height returns the height of the text entry.
func (text Text) Height() int {
	panic("not yet implemented")
}
