// Package font handles graphical text entries with customizable font size,
// style and color. It uses a small subset of the features provided by the SDL
// library version 2.0 [1].
//
// [1]: http://www.libsdl.org/
package font

// A Font provides glyphs (visual characters) and metrics used for text
// rendering.
type Font struct {
}

// Load loads the provided TTF font.
//
// Note: The Free method of the font must be called when finished using it.
func Load(path string) (f Font, err error) {
	panic("not yet implemented")
}

// Free frees the font.
func (f Font) Free() {
	panic("not yet implemented")
}
