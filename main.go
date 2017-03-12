package main

import (
	"flag"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	width  = 500.0
	height = 500.0
)

var (
	columns   = 20
	rows      = 20
	seed      = time.Now().UnixNano()
	threshold = 0.2
	fps       = 20
)

func init() {
	flag.IntVar(&columns, "columns", columns, "Sets the number of columns.")
	flag.IntVar(&rows, "rows", rows, "Sets the number of columns.")
	flag.Int64Var(&seed, "seed", seed, "Sets the starting seed of the game, used to randomize the initial state.")
	flag.Float64Var(&threshold, "threshold", threshold, "A percentage between 0 and 1 used in conjunction with the -seed to determine if a cell starts alive. For example, 0.15 means each cell has a 15% chance of starting alive.")
	flag.IntVar(&fps, "fps", fps, "Sets the frames-per-second, used set the speed of the simulation.")
	flag.Parse()

}

func main() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
	window := initGlfw()
	defer glfw.Terminate()
	prog := initOpenGL()

	// Make the cells and start the loop
	cells := makeCells(seed, threshold)
	t := time.Now()
	for !window.ShouldClose() {
		tick(cells)

		if err := draw(prog, window, cells); err != nil {
			panic(err)
		}

		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
		t = time.Now()
	}
}

// tick updates the state of each cell in the game board.
func tick(cells [][]*cell) {
	for x := range cells {
		for _, c := range cells[x] {
			c.checkState(cells)
		}
	}
}

// draw redraws the game board and the cells within.
func draw(prog uint32, window *glfw.Window, cells [][]*cell) error {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(prog)

	for x := range cells {
		for _, c := range cells[x] {
			c.draw()
		}
	}

	glfw.PollEvents()
	window.SwapBuffers()
	return nil
}
