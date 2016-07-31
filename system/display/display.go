package display

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

// Display holds everything needed to manage windows and draw to the screen
type Display struct {
	window *glfw.Window
}

// NewDisplay creates a new window and initializes OpenGL
func NewDisplay(width, height, pixelScale int) (*Display, error) {
	err := glfw.Init()
	if err != nil {
		return nil, err
	}

	// set up the window
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	window, err := glfw.CreateWindow(width*pixelScale, height*pixelScale, "Chippy", nil, nil)
	if err != nil {
		return nil, err
	}

	if err := gl.Init(); err != nil {
		return nil, err
	}

	window.MakeContextCurrent()

	gl.Disable(gl.DEPTH_TEST)
	gl.ClearColor(0.0, 0.215, 0.361, 1.0)

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()

	gl.Ortho(0, float64(width), float64(height), 0, 0, -1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	// just remove this once I'm ready to really start drawing
	return &Display{window: window}, nil
}

// Stop terminates drawing to the window
func (d *Display) Stop() {
	glfw.Terminate()
}

// PollOSEvents polls for events on the system, so things don't hang and window events can be processed
func (d *Display) PollOSEvents() {
	glfw.PollEvents()
}

// Draw draws the screenstate to the screen
func (d *Display) Draw(screenState [][]byte) {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Color3f(1, 1, 1)

	for y := range screenState {
		for x := range screenState[y] {
			if screenState[y][x] == 1 {
				gl.Recti(int32(x), int32(y), int32(x+1), int32(y+1))
			}
		}
	}

	d.window.SwapBuffers()
}
