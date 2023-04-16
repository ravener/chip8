package chip8

type Memory struct {
	RAM []byte
}

func NewMemory() *Memory {
	memory := new(Memory)
	memory.RAM = make([]byte, 4096)
	// Load fontset.
	copy(memory.RAM, Fontset[:])
	return memory
}

// Load a game ROM in the appropriate memory location.
func (m *Memory) LoadROM(rom []byte) {
	// Load ROM starting at address 0x200
	copy(m.RAM[0x200:], rom)
}

// Read a big-endian 16-bit value from the given address
func (m *Memory) ReadShort(address uint16) uint16 {
	return uint16(m.RAM[address])<<8 | uint16(m.RAM[address+1])
}
