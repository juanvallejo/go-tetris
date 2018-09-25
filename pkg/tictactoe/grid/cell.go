package grid

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"

	"github.com/juanvallejo/go-tictactoe/pkg/tictactoe/shape"
)

type Cell struct {
	color color.Color
	start pixel.Vec
	end   pixel.Vec
	width float64

	value *shape.Shape

	topLeft     *Cell
	topRight    *Cell
	top         *Cell
	left        *Cell
	right       *Cell
	bottom      *Cell
	bottomLeft  *Cell
	bottomRight *Cell
}

func (c *Cell) Start() pixel.Vec {
	return c.start
}

func (c *Cell) End() pixel.Vec {
	return c.end
}

func (c *Cell) Render(context *imdraw.IMDraw) {
	context.Color = c.color
	context.Push(c.start, c.end)
	context.Rectangle(c.width)

	if c.value != nil {
		c.value.Render(context)
	}
}

func (c *Cell) Set(shape *shape.Shape) bool {
	if c.value != nil {
		return false
	}

	c.value = shape
	return true
}
