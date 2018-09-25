package main

import (
	"github.com/faiface/pixel/pixelgl"

	"github.com/juanvallejo/go-tictactoe/pkg/tictactoe"
)

func main() {
	pixelgl.Run(tictactoe.NewGame)
}
