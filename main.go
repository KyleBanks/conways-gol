package main

import (
	"fmt"
	"log"
	"math"
	"runtime"
	"strings"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	width  = 500.0
	height = 500.0
	fps    = 60

	vertexShaderSource = `
		#version 400

		in vec3 vp;
		void main() {
			gl_Position = vec4(vp, 1.0);
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 400

		out vec4 frag_colour;
		void main() {
  			frag_colour = vec4(1, 1, 1, 1.0);
		}
	` + "\x00"
)

var (
	// Determines the starting makeup of the game grid.
	//
	// Points that are 0 start dead, and points that are 1 start alive.
	grid = []int{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 1, 1, 1, 0, 0, 0, 0,
		0, 0, 0, 1, 1, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	rowLen   = int(math.Sqrt(float64(len(grid))))
	cellSize = float32(width) / float32(rowLen) / float32(len(grid)) / 5.0
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()

	// Initialize glfw
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
}

func main() {
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(width, height, "Conway's Game of Life", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	cells := makeCells()
	prog := makeProgram()
	t := time.Now()
	for !window.ShouldClose() {
		tick(cells)

		if err := draw(prog, window, cells); err != nil {
			panic(err)
		}

		// Limit the FPS for better viewing.
		time.Sleep(time.Second/fps - time.Since(t))
		t = time.Now()
	}
}

func tick(cells []*Cell) {
	for i, c := range cells {
		c.checkState(i, cells)
	}
}

func draw(prog uint32, window *glfw.Window, cells []*Cell) error {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(prog)

	for _, cell := range cells {
		if !cell.alive {
			continue
		}

		gl.BindVertexArray(cell.vao)
		gl.DrawArrays(gl.TRIANGLES, 0, cell.pointCount)
	}

	glfw.PollEvents()
	window.SwapBuffers()
	return nil
}

func makeCells() (cells []*Cell) {
	for y := 0; y < rowLen; y++ {
		for x := 0; x < rowLen; x++ {
			c := newCell(float32(x)*cellSize, float32(y)*cellSize)
			c.alive = grid[y*rowLen+x] == 1
			c.nextState = c.alive

			cells = append(cells, c)
		}
	}

	return
}

func makeProgram() uint32 {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)

	return prog
}

// https://github.com/go-gl/examples/blob/master/gl41core-cube/cube.go
func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func makeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}
