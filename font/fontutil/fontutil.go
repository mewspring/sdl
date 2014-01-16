// Package fontutil provides font utility functions for word wrapping.
package fontutil

import (
	"image"
	"strings"

	"github.com/mewmew/sdl/font"
	"github.com/mewmew/sdl/win"
)

// A Box is a bounding box which confines texts to a fixed pixel width.
type Box struct {
	// The vertical space in between lines.
	spacing int
	// The font used for text rendering.
	f *font.Font
	// The width of the box.
	width int
}

// NewBox creates a new bounding box of the specified width. Its split and
// render logic wraps the text so that no single line overflows the pixel width
// of the box.
func NewBox(f *font.Font, width int) (box *Box) {
	box = &Box{
		f:     f,
		width: width,
	}
	return box
}

// Split splits the text into lines so that no single line overflows the pixel
// width of the box. It tries to wrap the text in between words if possible.
func (box *Box) Split(text string) (lines []string, err error) {
	// Wrap each individual line of text based on the width of the bounding box.
	textLines := strings.Split(text, "\n")
	for _, textLine := range textLines {
		ls, err := box.splitLine(textLine)
		if err != nil {
			return nil, err
		}
		lines = append(lines, ls...)
	}
	return lines, nil
}

// splitLine split the text into lines so that no single line ovreflows the
// pixel width. It tries to wrap the text in between words if possible.
func (box *Box) splitLine(text string) (lines []string, err error) {
	var start, end, spaceEnd int
	for i, r := range text {
		width, err := box.f.RenderWidth(text[start:i])
		if err != nil {
			return nil, err
		}
		if width <= box.width {
			if r == ' ' {
				spaceEnd = i
			}
			end = i
		} else {
			if spaceEnd != 0 {
				line := text[start:spaceEnd]
				lines = append(lines, line)
				start = spaceEnd + 1
				spaceEnd = 0
			} else {
				line := text[start:end]
				lines = append(lines, line)
				start = end
			}
			if r == ' ' {
				spaceEnd = i
			}
		}
	}
	line := text[start:]
	lines = append(lines, line)
	return lines, nil
}

// SetSpacing sets the vertical space in between lines to spacing. This
// information is used by the Render method while rendering lines. The default
// vertical space in between lines is 0.
func (box *Box) SetSpacing(spacing int) {
	box.spacing = spacing
}

// Render renders the text using the font of the box. It wraps the text so that
// no single line overflows the pixel width of the box.
func (box *Box) Render(text string) (img *win.Image, err error) {
	// Break the input text into lines based on the box width.
	lines, err := box.Split(text)
	if err != nil {
		return nil, err
	}

	// Create an image large enough to hold all lines.
	width := box.width
	height := box.f.Height*len(lines) + box.spacing*(len(lines)-1)
	img, err = win.NewImage(width, height)
	if err != nil {
		return nil, err
	}

	// Draw each line onto the final image.
	for i, line := range lines {
		if line == "" {
			continue
		}
		y := i * (box.f.Height + box.spacing)
		dp := image.Pt(0, y)
		partImg, err := box.f.Render(line)
		if err != nil {
			return nil, err
		}
		err = img.Draw(dp, partImg)
		if err != nil {
			return nil, err
		}
	}

	return img, nil
}
