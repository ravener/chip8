package ui

import (
	"fmt"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/ravener/chip8/chip8"
	"github.com/ravener/go-gl"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

// Read the ROM file and return the file contents in bytes.
func loadROM(file string) []byte {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return data
}

func Run(file string) {
	// Seed the random number generator.
	rand.Seed(time.Now().Unix())

	// Initialize GLFW.
	if err := glfw.Init(); err != nil {
		log.Fatalln(err)
	}
	defer glfw.Terminate()

	// Create Window.
	window := createWindow()
	defer window.Destroy()

	// Enable 2D texturing.
	gl.Enable(gl.TEXTURE_2D)

	// Initialize CPU and load the ROM.
	cpu := chip8.NewCPU(chip8.NewMemory())
	cpu.Memory.LoadROM(loadROM(file))
	paused := false // Whether we are paused.

	// Initialize the texture.
	texture := createTexture(cpu.Display)

	// Handle Input.
	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Press {
			switch key {
			// Take a screenshot with F11
			case glfw.KeyF11:
				file, err := os.Create("screenshot-" + fmt.Sprint(time.Now().Unix()) + ".png")
				if err != nil {
					panic(err)
				}
				png.Encode(file, cpu.Display)
			// Pause the emulator with space.
			case glfw.KeySpace:
				paused = !paused
				if paused {
					window.SetTitle("CHIP-8 (Paused)")
				} else {
					window.SetTitle("CHIP-8")
				}
			// Allow quitting with the escape button.
			case glfw.KeyEscape:
				window.SetShouldClose(true)
			}
		}

		// Map the physical key to the CHIP-8 keypad.
		k := mapKey(key)
		// Ignore other keys.
		if k == -1 {
			return
		}

		// Update the key's state in the CPU.
		switch action {
		case glfw.Press:
			cpu.Keys.Press(uint8(k))
		case glfw.Release:
			cpu.Keys.Release(uint8(k))
		}
	})

	// The main loop.
	// Note: We make use of V-Sync and assume this will run at 60Hz
	// If V-Sync is force-disabled or the monitor's refresh rate is over 60Hz
	// then this may not work as expected. That's something to fix.
	for !window.ShouldClose() {
		if !paused {
			// Execute a CPU cycle 9 times per-frame.
			// Assuming a perfect 60 FPS then 60*9 = 540
			// Since the CPU must run at ~500-540 Hz
			for i := 0; i < 9; i++ {
				cpu.Execute()
			}
		}

		// If the screen changed and needs to be redrawn.
		if cpu.Draw {
			// Clear the screen.
			gl.Clear(gl.COLOR_BUFFER_BIT)

			// Bind and update the GPU texture with the current display pixels.
			gl.BindTexture(gl.TEXTURE_2D, texture)
			updateTexture(cpu.Display)

			// Render a textured quad that covers the whole screen.
			// Warning: Deprecated/Legacy OpenGL functions used for simplicity.
			gl.Begin(gl.QUADS)
			gl.TexCoord2f(0, 1)
			gl.Vertex2f(-1, -1)
			gl.TexCoord2f(1, 1)
			gl.Vertex2f(1, -1)
			gl.TexCoord2f(1, 0)
			gl.Vertex2f(1, 1)
			gl.TexCoord2f(0, 0)
			gl.Vertex2f(-1, 1)
			gl.End()

			// Swap the back buffer.
			window.SwapBuffers()
			// We are done drawing, mark it off until the display needs to redrawn again.
			cpu.Draw = false
		}

		// Update timers at 60Hz
		cpu.UpdateTimers()
		// Poll OS Events
		glfw.PollEvents()
	}
}
