package main

import (
	"log"
	"math/rand"

	"github.com/go-gl/gl/v4.1-core/gl"
)

var (
	squarePoints = []float32{
		// Bottom left right-angle triangle
		-1, 1, 0,
		1, -1, 0,
		-1, -1, 0,

		// Top right right-angle triangle
		-1, 1, 0,
		1, 1, 0,
		1, -1, 0,
	}

	squarePointCount = int32(len(squarePoints) / 3)
)

type Cell struct {
	drawable uint32

	x int
	y int

	alive     bool
	nextState bool
}

// checkState determines the state of the Cell for the next tick of the game.
func (c *Cell) checkState(cells [][]*Cell) {
	c.alive = c.nextState
	c.nextState = c.alive

	liveCount := c.liveNeighbors(cells)
	if c.alive {
		// 1. Any live cell with fewer than two live neighbours dies, as if caused by underpopulation.
		if liveCount < 2 {
			c.nextState = false
		}

		// 2. Any live cell with two or three live neighbours lives on to the next generation.
		if liveCount == 2 || liveCount == 3 {
			c.nextState = true
		}

		// 3. Any live cell with more than three live neighbours dies, as if by overpopulation.
		if liveCount > 3 {
			c.nextState = false
		}
	} else {
		// 4. Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.
		if liveCount == 3 {
			c.nextState = true
		}
	}
}

// liveNeighbors returns the number of live neighbors for a Cell.
func (c *Cell) liveNeighbors(cells [][]*Cell) int {
	var liveCount int
	add := func(x, y int) {
		// If we're at an edge, check the other side of the board.
		if x == len(cells) {
			x = 0
		} else if x == -1 {
			x = len(cells) - 1
		}
		if y == len(cells[x]) {
			y = 0
		} else if y == -1 {
			y = len(cells[x]) - 1
		}

		if cells[x][y].alive {
			liveCount++
		}
	}

	add(c.x-1, c.y)   // To the left
	add(c.x+1, c.y)   // To the right
	add(c.x, c.y+1)   // up
	add(c.x, c.y-1)   // down
	add(c.x-1, c.y+1) // top-left
	add(c.x+1, c.y+1) // top-right
	add(c.x-1, c.y-1) // bottom-left
	add(c.x+1, c.y-1) // bottom-right

	return liveCount
}

// draw draws the Cell if it is alive.
func (c *Cell) draw() {
	if !c.alive {
		return
	}

	gl.BindVertexArray(c.drawable)
	gl.DrawArrays(gl.TRIANGLES, 0, squarePointCount)
}

// newCell initializes and returns a Cell with the given x/y coordinates.
func newCell(x, y int) *Cell {
	points := make([]float32, len(squarePoints), len(squarePoints))
	copy(points, squarePoints)

	for i := 0; i < len(points); i++ {
		var factor float32
		var size float32
		switch i % 3 {
		case 0:
			size = 1.0 / float32(columns)
			factor = float32(x) * size
		case 1:
			size = 1.0 / float32(rows)
			factor = float32(y) * size
		default:
			continue
		}

		if points[i] < 0 {
			points[i] = (factor * 2) - 1
		} else {
			points[i] = ((factor + size) * 2) - 1
		}
	}

	return &Cell{
		drawable: makeVao(points),

		x: x,
		y: y,
	}
}

// makeCells creates the Cell matrix and sets the initial state of the game.
func makeCells(seed int64, threshold float64) [][]*Cell {
	log.Printf("Using seed=%v, threshold=%v", seed, threshold)
	rand.Seed(seed)

	cells := make([][]*Cell, rows, rows)
	for x := 0; x < rows; x++ {
		for y := 0; y < columns; y++ {
			c := newCell(x, y)

			c.alive = rand.Float64() < threshold
			c.nextState = c.alive

			cells[x] = append(cells[x], c)
		}
	}

	return cells
}
