package grid

import (
	"fmt"
	"image/color"
	"math"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/text"

	"github.com/juanvallejo/tictactoe/pkg/tictactoe/shape"
)

const (
	MaxCells      = 3
	gridLineWidth = 3
)

var gridLineColor = colornames.Antiquewhite

type Grid []*Cell

func (g Grid) Render(context *imdraw.IMDraw) {
	for i := range g {
		g[i].Render(context)
	}

	// render vertical lines
	for i := 0; i < MaxCells; i++ {
		if (i+1)%MaxCells == 0 {
			continue
		}

		context.Color = gridLineColor
		context.Push(pixel.V(g[i].end.X, g[i].start.Y))
		context.Push(pixel.V(g[i].end.X, g[len(g)-1].end.Y))
		context.Line(gridLineWidth)
	}

	// render horizontal grid lines
	for i := range g {
		if i%MaxCells != 0 || i/MaxCells+1 == MaxCells {
			continue
		}

		context.Color = gridLineColor
		context.Push(pixel.V(g[i].start.X, g[i].end.Y))
		context.Push(pixel.V(g[len(g)-1].end.X, g[i].end.Y))
		context.Line(gridLineWidth)
	}
}

func (g Grid) Reset() {
	for i := range g {
		g[i].value = nil
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

func (g Grid) CheckWin(context *imdraw.IMDraw, textContext *text.Text) bool {
	if len(g) != MaxCells*MaxCells {
		panic(fmt.Sprintf("malformed grid: expected to check wins on a %d by %d grid, but total grid size was %d", MaxCells, MaxCells, len(g)))
	}

	// check vertical wins
	for i := 0; i < MaxCells; i++ {
		if checkVerticalCount(g[i], g[i].value) >= MaxCells {
			return drawVerticalWin(context, g[i]) && drawText(textContext, getWinText(g[i]))
		}
	}

	// check hortizontal wins
	for i := range g {
		if i%MaxCells != 0 {
			continue
		}
		if checkHorizontalCount(g[i], g[i].value) >= MaxCells {
			return drawHorizontalWin(context, g[i]) && drawText(textContext, getWinText(g[i]))
		}
	}

	// check diagonal wins
	if checkDiagonalCount(g[0], g[0].value, true) >= MaxCells {
		return drawDiagonalWin(context, g[0], true) && drawText(textContext, getWinText(g[0]))
	}
	if checkDiagonalCount(g[MaxCells-1], g[MaxCells-1].value, false) >= MaxCells {
		return drawDiagonalWin(context, g[MaxCells-1], false) && drawText(textContext, getWinText(g[MaxCells-1]))
	}

	// check draw
	for i := range g {
		if g[i].value == nil {
			return false
		}
	}

	return drawText(textContext, "TIE!")
}

func drawText(context *text.Text, contents string) bool {
	context.Dot.X -= context.BoundsOf(contents).W() / 2
	context.Dot.Y -= context.BoundsOf(contents).H() / 2
	fmt.Fprintf(context, "%s\n", contents)
	return true
}

// getWinText returns the string of text presented on a win
// depending on the value of cell.
func getWinText(cell *Cell) string {
	playerWin := 1
	if cell.value != nil && cell.value.Kind() == shape.CircleShape {
		playerWin = 2
	}

	return fmt.Sprintf("PLAYER %d WINS!", playerWin)
}

func drawHorizontalWin(context *imdraw.IMDraw, from *Cell) bool {
	context.Color = shape.ShapeColor
	context.Push(pixel.V(from.start.X-gridLineWidth, from.start.Y-((from.start.Y-from.end.Y)/2)))
	context.Push(pixel.V(from.end.X+gridLineWidth, from.start.Y-((from.start.Y-from.end.Y)/2)))
	context.Line(gridLineWidth)

	if from.right == nil {
		return true
	}

	return drawHorizontalWin(context, from.right)
}

func drawVerticalWin(context *imdraw.IMDraw, from *Cell) bool {
	context.Color = shape.ShapeColor
	context.Push(pixel.V(from.start.X+((from.end.X-from.start.X)/2), from.start.Y+gridLineWidth))
	context.Push(pixel.V(from.start.X+((from.end.X-from.start.X)/2), math.Max(from.end.Y-gridLineWidth, 0)))
	context.Line(gridLineWidth)

	if from.bottom == nil {
		return true
	}

	return drawVerticalWin(context, from.bottom)
}

func drawDiagonalWin(context *imdraw.IMDraw, from *Cell, bottomRight bool) bool {
	context.Color = shape.ShapeColor
	if bottomRight {
		context.Push(pixel.V(from.start.X, from.start.Y+gridLineWidth))
		context.Push(pixel.V(from.end.X, math.Max(from.end.Y-gridLineWidth, 0)))
	} else {
		context.Push(pixel.V(from.end.X, from.start.Y+gridLineWidth))
		context.Push(pixel.V(from.start.X, math.Max(from.end.Y-gridLineWidth, 0)))
	}
	context.Line(gridLineWidth)

	if (from.bottomRight == nil && bottomRight) || (from.bottomLeft == nil && !bottomRight) {
		return true
	}
	if bottomRight {
		return drawDiagonalWin(context, from.bottomRight, bottomRight)
	}

	return drawDiagonalWin(context, from.bottomLeft, bottomRight)
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
				color: color.Transparent,
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
