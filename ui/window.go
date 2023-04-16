package ui

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/ravener/go-gl"
)

const (
	width  = 64 // original width of CHIP-8 screen.
	height = 32 // original height of CHIP-8 screen.
	scale  = 10 // scale 10x for the window size.
)

// Create window and initialize OpenGL context.
func createWindow() *glfw.Window {
	// OpenGL 2.1
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	// Create window.
	window, err := glfw.CreateWindow(width*scale, height*scale, "CHIP-8", nil, nil)
	if err != nil {
		panic(err)
	}

	// Make this OpenGL context the current context.
	window.MakeContextCurrent()
	// Initialize OpenGL functions for this context.
	if err := gl.Init(); err != nil {
		panic(err)
	}

	window.SetSizeCallback(func(w *glfw.Window, width, height int) {
		// Update the OpenGL viewport when the window size changes.
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	// Enable V-Sync
	glfw.SwapInterval(1)

	return window
}
