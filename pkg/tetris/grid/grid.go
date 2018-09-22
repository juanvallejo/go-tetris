package grid

import (
	"fmt"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"

	"github.com/juanvallejo/tetris/pkg/tetris/shape"
)

const (
	MaxCells = 3
)

type Grid []*Cell

func (g Grid) Render(context *imdraw.IMDraw) {
	for i := range g {
		g[i].Render(context)
	}
}

// AtVector receives a vector and returns the cell containing
// that point, or nil. If two overlapping cells contain the
// point, then the first cell found is returned.
func (g Grid) AtVector(v pixel.Vec) *Cell {
	for i := range g {
		if v.X > g[i].start.X && v.X < g[i].end.X && v.Y < g[i].start.Y && v.Y > g[i].end.Y {
			return g[i]
		}
	}
	return nil
}

func (g Grid) CheckWin() bool {
	if len(g) != MaxCells*MaxCells {
		panic(fmt.Sprintf("malformed grid: expected to check wins on a %d by %d grid, but total grid size was %d", MaxCells, MaxCells, len(g)))
	}

	// check vertical wins
	for i := 0; i < MaxCells; i++ {
		if checkVerticalCount(g[i], g[i].value) >= MaxCells {
			return true
		}
	}

	// check hortizontal wins
	for i := range g {
		if i%MaxCells != 0 {
			continue
		}
		if checkHorizontalCount(g[i], g[i].value) >= MaxCells {
			return true
		}
	}

	// check diagonal wins
	if checkDiagonalCount(g[0], g[0].value, true) >= MaxCells || checkDiagonalCount(g[MaxCells-1], g[MaxCells-1].value, false) >= MaxCells {
		return true
	}

	return false
}

func checkDiagonalCount(cell *Cell, target *shape.Shape, bottomRight bool) int {
	if cell == nil || target == nil || cell.value == nil || cell.value.Kind() != target.Kind() {
		return 0
	}
	if (cell.bottomRight == nil && bottomRight) || (cell.bottomLeft == nil && !bottomRight) {
		return 1
	}

	if bottomRight {
		return 1 + checkDiagonalCount(cell.bottomRight, target, bottomRight)
	}
	return 1 + checkDiagonalCount(cell.bottomLeft, target, bottomRight)
}

func checkVerticalCount(cell *Cell, target *shape.Shape) int {
	if cell == nil || target == nil || cell.value == nil || cell.value.Kind() != target.Kind() {
		return 0
	}
	if cell.bottom == nil {
		return 1
	}

	return 1 + checkVerticalCount(cell.bottom, target)
}

func checkHorizontalCount(cell *Cell, target *shape.Shape) int {
	if cell == nil || target == nil || cell.value == nil || cell.value.Kind() != target.Kind() {
		return 0
	}
	if cell.right == nil {
		return 1
	}

	return 1 + checkHorizontalCount(cell.right, target)
}

func NewGrid(origin pixel.Vec, maxX, maxY, ncells, mar float64) Grid {
	margin := pixel.V(mar, mar)
	cellWidth := math.Floor((maxX - (margin.X * 2)) / ncells)
	cellHeight := math.Floor((maxY - (margin.Y * 2)) / ncells)

	origin = origin.Add(pixel.V(margin.X, maxY-margin.Y))

	cells := []*Cell{}
	for y := 0; y < int(ncells); y++ {
		for x := 0; x < int(ncells); x++ {
			start := origin.Add(pixel.V(cellWidth*float64(x), -cellHeight*float64(y)))
			cells = append(cells, &Cell{
				color: colornames.Darkturquoise,
				start: start,
				end:   start.Add(pixel.V(cellWidth, -cellHeight)),
				width: 3,
			})
		}
	}

	// populate neighbors
	for i := range cells {
		x := i % int(ncells)
		y := int(i / int(ncells))

		// top left neighbor
		if x > 0 && y > 0 {
			cells[i].topLeft = cells[(int(ncells)*y-(int(ncells)-x))-1]
		}

		// top right neighbor
		if x < int(ncells)-1 && y > 0 {
			cells[i].topRight = cells[(int(ncells)*y-(int(ncells)-x))+1]
		}

		// top neighbor
		if y > 0 {
			cells[i].top = cells[int(ncells)*y-(int(ncells)-x)]
		}

		// left neighbor
		if x > 0 {
			cells[i].left = cells[(int(ncells)*(y+1)-(int(ncells)-x))-1]
		}

		// right neighbor
		if x < int(ncells)-1 {
			cells[i].right = cells[(int(ncells)*(y+1)-(int(ncells)-x))+1]
		}

		// bottom neighbor
		if y < int(ncells)-1 {
			cells[i].bottom = cells[int(ncells)*(y+2)-(int(ncells)-x)]
		}

		// bottom left neighbor
		if x > 0 && y < int(ncells)-1 {
			cells[i].bottomLeft = cells[(int(ncells)*(y+2)-(int(ncells)-x))-1]
		}

		// bottom right neighbor
		if x < int(ncells)-1 && y < int(ncells)-1 {
			cells[i].bottomRight = cells[(int(ncells)*(y+2)-(int(ncells)-x))+1]
		}
	}

	return Grid(cells)
}
