package shape

import (
	"fmt"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

type ShapeKind string

var (
	ShapeColor            = colornames.Thistle
	CircleShape ShapeKind = "O"
	CrossShape  ShapeKind = "X"
)

type Shape struct {
	color color.Color
	start pixel.Vec
	end   pixel.Vec
	kind  ShapeKind
	width float64
}

func (s *Shape) Kind() ShapeKind {
	return s.kind
}

func (s *Shape) String() string {
	return string(s.kind)
}

func (s *Shape) Render(context *imdraw.IMDraw) {
	context.Color = s.color
	if s.kind == CrossShape {
		horLength := s.end.X - s.start.X
		verLength := s.end.Y - s.start.Y

		startX := s.start.X
		endX := s.end.X
		startY := s.start.Y
		endY := s.end.Y

		modifier := (horLength - verLength) / 2
		if modifier > 0 {
			startX += modifier
			endX -= modifier
		} else {
			startY += modifier
			endY -= modifier
		}

		context.Push(pixel.V(startX, startY))
		context.Push(pixel.V(endX, endY))
		context.Line(s.width)
		context.Push(pixel.V(endX, startY))
		context.Push(pixel.V(startX, endY))
		context.Line(s.width)
		return
	}

	if s.kind != CircleShape {
		panic(fmt.Sprintf("undefined shape: %s", s.kind))
	}

	context.Push(pixel.V((s.start.X+s.end.X)/2, (s.start.Y+s.end.Y)/2))
	context.Circle((s.end.Y-s.start.Y)/2, s.width)
}

func NewShape(origin pixel.Vec, shapeKind ShapeKind, width, height, mar float64) *Shape {
	margin := pixel.V(mar, mar)
	start := origin.Add(pixel.V(margin.X, -margin.Y))
	size := pixel.V(width-margin.X*2, -height+margin.Y*2)

	return &Shape{
		color: ShapeColor,
		start: start,
		end:   start.Add(size),
		kind:  shapeKind,
		width: 3,
	}
}

type ShapeDecider struct {
	next ShapeKind
}

func (n *ShapeDecider) Next() ShapeKind {
	next := n.next
	if n.next == CrossShape {
		n.next = CircleShape
	} else {
		n.next = CrossShape
	}

	return next
}

func NewShapeDecider(startingShape ShapeKind) *ShapeDecider {
	if startingShape != CrossShape && startingShape != CircleShape {
		panic(fmt.Sprintf("invalid shape-decider shape: %v", startingShape))
	}

	return &ShapeDecider{next: startingShape}
}
