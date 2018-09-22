package main

import (
	"github.com/faiface/pixel/pixelgl"

	"github.com/juanvallejo/tetris/pkg/tetris"
)

func main() {
	pixelgl.Run(tetris.NewGame)
}
