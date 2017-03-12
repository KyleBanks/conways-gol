package main

var squarePoints = []float32{
	// Bottom left right-angle triangle
	-cellSize, 0, 0,
	0, -cellSize, 0,
	-cellSize, -cellSize, 0,

	// Top right right-angle triangle
	-cellSize, 0, 0,
	0, 0, 0,
	0, -cellSize, 0,
}

type Cell struct {
	pointCount int32
	vao        uint32

	alive     bool
	nextState bool
}

func (c *Cell) checkState(idx int, cells []*Cell) {
	c.alive = c.nextState
	c.nextState = c.alive

	liveCount := c.liveNeighbors(idx, cells)

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
		return
	}

	// 4. Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.
	if liveCount == 3 {
		c.nextState = true
	}
}

func (c *Cell) liveNeighbors(idx int, cells []*Cell) int {
	var liveCount int
	add := func(i int) {
		if i < len(cells) && i >= 0 && cells[i].alive {
			liveCount++
		}
	}

	add(idx - 1)          // To the left
	add(idx + 1)          // To the right
	add(idx + rowLen)     // up
	add(idx - rowLen)     // down
	add(idx + rowLen - 1) // top-left
	add(idx + rowLen + 1) // top-right
	add(idx - rowLen - 1) // bottom-left
	add(idx - rowLen + 1) // bottom-right

	return liveCount
}

func newCell(x, y float32) *Cell {
	points := make([]float32, len(squarePoints), len(squarePoints))
	copy(points, squarePoints)

	for i := 0; i < len(points); i++ {
		var factor float32
		switch i % 3 {
		case 0:
			factor = x
		case 1:
			factor = y
		default:
			continue
		}

		if points[i] < 0 {
			points[i] = (factor * 2) - 1
		} else {
			points[i] = ((factor + cellSize) * 2) - 1
		}
	}

	return &Cell{
		vao:        makeVao(points),
		pointCount: int32(len(points) / 3),
	}
}
