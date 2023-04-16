package ui

import "github.com/go-gl/glfw/v3.3/glfw"

// Map the given physical key to the CHIP-8 keypad.
// We use a mapping of:
// Keypad             Keyboard
// +-+-+-+-+          +-+-+-+-+
// |1|2|3|C|          |1|2|3|4|
// +-+-+-+-+          +-+-+-+-+
// |4|5|6|D|          |Q|W|E|R|
// +-+-+-+-+    =>    +-+-+-+-+
// |7|8|9|E|          |A|S|D|F|
// +-+-+-+-+          +-+-+-+-+
// |A|0|B|F|          |Z|X|C|V|
// +-+-+-+-+          +-+-+-+-+
func mapKey(key glfw.Key) int8 {
	switch key {
	case glfw.Key1:
		return 0x1
	case glfw.Key2:
		return 0x2
	case glfw.Key3:
		return 0x3
	case glfw.Key4:
		return 0xC
	case glfw.KeyQ:
		return 0x4
	case glfw.KeyW:
		return 0x5
	case glfw.KeyE:
		return 0x6
	case glfw.KeyR:
		return 0xD
	case glfw.KeyA:
		return 0x7
	case glfw.KeyS:
		return 0x8
	case glfw.KeyD:
		return 0x9
	case glfw.KeyF:
		return 0xE
	case glfw.KeyZ:
		return 0xA
	case glfw.KeyX:
		return 0x0
	case glfw.KeyC:
		return 0xB
	case glfw.KeyV:
		return 0xF
	default:
		// Ignore other keys.
		return -1
	}
}
