package ui

import (
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

// Declare conformity with Layout interface
var _ fyne.Layout = (*barLayout)(nil)

const (
	iconZoomDistance = 2.5
	separatorWidth   = 2
)

//barLayout returns a layout used for zooming linear groups of icons
type barLayout struct {
	bar *bar

	mouseInside   bool          // Is the mouse inside of the layout?
	mousePosition fyne.Position // Current coordinates of the mouse cursor
}

//setPointerInside tells the barLayout that the mouse is inside of the Layout.
func (bl *barLayout) setPointerInside(inside bool) {
	bl.mouseInside = inside
}

//setPointerPosition tells the barLayout that the mouse position has been updated.
func (bl *barLayout) setPointerPosition(position fyne.Position) {
	bl.mousePosition = position
}

// Layout is called to pack all icons into a specified size.  It also handles the zooming effect of the icons.
func (bl *barLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	bg := objects[0]
	objects = objects[1:]

	offset := 0.0
	barWidth := len(objects)*(bl.bar.iconSize+theme.Padding()) - theme.Padding()
	if !bl.bar.disableTaskbar {
		barWidth = barWidth - bl.bar.iconSize + separatorWidth
	}
	barLeft := (size.Width - barWidth) / 2
	iconLeft := barLeft

	mouseX := bl.mousePosition.X
	zoom := !bl.bar.disableZoom && bl.mouseInside && mouseX >= barLeft && mouseX < barLeft+barWidth
	for _, child := range objects {
		if zoom {
			if _, ok := child.(*canvas.Rectangle); ok {
				child.Resize(fyne.NewSize(separatorWidth, bl.bar.iconSize))
				if iconLeft+separatorWidth+theme.Padding() < mouseX {
					offset += separatorWidth
				} else if iconLeft < mouseX {
					offset += float64(separatorWidth + theme.Padding())
				}
			} else {
				iconCenter := iconLeft + bl.bar.iconSize/2
				offsetX := float64(mouseX - iconCenter)

				scale := float64(bl.bar.iconScale) - (math.Abs(offsetX) / (float64(bl.bar.iconSize) * iconZoomDistance))
				newSize := float64(bl.bar.iconSize) * scale
				if int(newSize) < bl.bar.iconSize {
					newSize = float64(bl.bar.iconSize)
				}
				child.Resize(fyne.NewSize(int(newSize), int(newSize)))

				if iconLeft+bl.bar.iconSize+theme.Padding() < mouseX {
					offset += newSize - float64(bl.bar.iconSize)
				} else if iconLeft < mouseX {
					ratio := float64(mouseX-iconLeft) / float64(bl.bar.iconSize+theme.Padding())
					offset += (newSize-float64(bl.bar.iconSize))*ratio + float64(theme.Padding())
				}
			}
		} else {
			if _, ok := child.(*canvas.Rectangle); ok {
				child.Resize(fyne.NewSize(separatorWidth, bl.bar.iconSize))
			} else {
				child.Resize(fyne.NewSize(bl.bar.iconSize, bl.bar.iconSize))
			}
		}
		if _, ok := child.(*canvas.Rectangle); ok {
			iconLeft += separatorWidth
		} else {
			iconLeft += bl.bar.iconSize
		}
		iconLeft += theme.Padding()
	}

	x := barLeft - int(offset)
	zoomLeft := x
	tallHeight := int(float32(bl.bar.iconSize) * bl.bar.iconScale)
	for _, child := range objects {
		width := child.Size().Width
		height := child.Size().Height

		if zoom {
			if _, ok := child.(*canvas.Rectangle); ok {
				child.Move(fyne.NewPos(x, bl.bar.iconSize))
			} else {
				child.Move(fyne.NewPos(x, tallHeight-height))
			}
		} else {
			child.Move(fyne.NewPos(x, 0))
		}
		x += width + theme.Padding()
	}
	if zoom {
		bg.Move(fyne.NewPos(zoomLeft-theme.Padding(), bl.bar.iconSize))
	} else {
		bg.Move(fyne.NewPos(zoomLeft-theme.Padding(), 0))
	}
	bg.Resize(fyne.NewSize(x-zoomLeft+theme.Padding(), bl.bar.iconSize))
}

// MinSize finds the smallest size that satisfies all the child objects.
// For a barLayout this is the width of the widest item and the height is
// the sum of of all children combined with padding between each.
func (bl *barLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	barWidth := 0
	iconCount := len(objects)
	if !bl.bar.disableTaskbar {
		iconCount--
		barWidth = (iconCount * (bl.bar.iconSize + theme.Padding())) + separatorWidth
	} else {
		barWidth = iconCount * (bl.bar.iconSize + theme.Padding())
	}

	barLeft := (bl.bar.Size().Width - barWidth) / 2
	mouseX := bl.mousePosition.X
	if !bl.bar.disableZoom && bl.mouseInside && mouseX >= barLeft && mouseX < barLeft+barWidth {
		return fyne.NewSize(barWidth, int(float32(bl.bar.iconSize)*bl.bar.iconScale))
	}

	return fyne.NewSize(barWidth, bl.bar.iconSize)
}

// newBarLayout returns a horizontal icon bar
func newBarLayout(bar *bar) barLayout {
	return barLayout{bar, false, fyne.NewPos(0, 0)}
}
