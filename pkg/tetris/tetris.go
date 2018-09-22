package tetris

import (
	"fmt"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"

	"github.com/juanvallejo/tetris/pkg/tetris/grid"
	"github.com/juanvallejo/tetris/pkg/tetris/shape"
)

const (
	winWidth  = 800
	winHeight = 600

	cellMargin  = 45
	shapeMargin = 25
)

func NewGame() {
	config := pixelgl.WindowConfig{
		Title:  "Tetris",
		Bounds: pixel.R(0, 0, winWidth, winHeight),
		VSync:  true,
	}

	window, err := pixelgl.NewWindow(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	shapeDecider := shape.NewShapeDecider(shape.CrossShape)
	bounds := window.Bounds()
	context := imdraw.New(nil)
	g := grid.NewGrid(pixel.V(0, 0), bounds.Max.X, bounds.Max.Y, grid.MaxCells, cellMargin)

	for !window.Closed() {
		context.Clear()

		if window.JustPressed(pixelgl.MouseButtonLeft) {
			if cell := g.AtVector(window.MousePosition()); cell != nil {
				cell.Set(shape.NewShape(pixel.V(cell.Start().X, cell.Start().Y), shapeDecider.Next(), (bounds.Max.X-cellMargin*2)/grid.MaxCells, (bounds.Max.Y-cellMargin*2)/grid.MaxCells, shapeMargin))
				if g.CheckWin() {
					fmt.Printf("I WIN!!!\n")
				}
			}
		}

		g.Render(context)

		context.Draw(window)
		window.Update()
	}
}
