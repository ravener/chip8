package chip8

type Keys [16]bool

// Check if a key is currently pressed.
func (keys *Keys) IsPressed(key uint8) bool {
	return keys[key]
}

// Set a key to pressed state.
func (keys *Keys) Press(key uint8) {
	keys[key] = true
}

// Set a key to released state.
func (keys *Keys) Release(key uint8) {
	keys[key] = false
}
