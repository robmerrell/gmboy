package ui

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/robmerrell/gmboy/system/debugger"
)

// InputState is the current state of all buttons. Since everything can be considered
// pressed (true) or not pressed (false) booleans work well for us here.
type InputState struct {
	// buttons!
	Up     bool
	Down   bool
	Left   bool
	Right  bool
	A      bool
	B      bool
	Select bool
	Start  bool

	// window where we want to watch for input events
	window *glfw.Window

	// debugger
	debugger *debugger.Debugger
}

// NewInput creates a new input state for the game controls. A display is expected so that we know which window to look for keypresses in.
func NewInput(display *Display) *InputState {
	return &InputState{window: display.window}
}

// AttachDebugger attaches a javascript debugger to the InputState and sets up hotkeys
// for the debugger.
func (i *InputState) AttachDebugger(dbg *debugger.Debugger) {
	i.debugger = dbg

	i.window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		// set the N key to "next" the debugger
		if key == glfw.KeyN && action == glfw.Release {
			i.debugger.Next()
		}

		// set the K key to "continue" the debugger
		if key == glfw.KeyK && action == glfw.Release {
			i.debugger.Continue()
		}

		// set the M key to advance one frame
		if key == glfw.KeyM && action == glfw.Release {
		}
	})
}

// updateState updates the input state of a controller.
func (i *InputState) updateState() {
	i.Up = i.window.GetKey(glfw.KeyUp) == glfw.Press
	i.Down = i.window.GetKey(glfw.KeyDown) == glfw.Press
	i.Left = i.window.GetKey(glfw.KeyLeft) == glfw.Press
	i.Right = i.window.GetKey(glfw.KeyRight) == glfw.Press
	i.A = i.window.GetKey(glfw.KeyX) == glfw.Press
	i.B = i.window.GetKey(glfw.KeyZ) == glfw.Press
	i.Select = i.window.GetKey(glfw.KeyRightShift) == glfw.Press
	i.Start = i.window.GetKey(glfw.KeyEnter) == glfw.Press
}
