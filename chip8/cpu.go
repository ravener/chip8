package chip8

import (
	"image"
	"image/color"
	"log"
	"math/rand"
)

type CPU struct {
	PC         uint16      // program counter
	V          [16]uint8   // registers
	I          uint16      // index register
	SP         uint8       // stack pointer
	stack      [16]uint16  // call stack
	Memory     *Memory     // memory
	DelayTimer uint8       // delay timer
	SoundTimer uint8       // sound timer
	Display    *image.RGBA // display
	Draw       bool        // draw flag
	Keys       Keys        // input state
}

func NewCPU(memory *Memory) *CPU {
	cpu := &CPU{
		Memory:  memory,
		PC:      0x200,
		Display: image.NewRGBA(image.Rect(0, 0, 64, 32)),
	}

	// Clear the display to a black screen.
	cpu.ClearDisplay()

	return cpu
}

func (cpu *CPU) fetch() uint16 {
	instruction := cpu.Memory.ReadShort(cpu.PC)
	// Increment program counter to point to the next instruction.
	cpu.PC += 2
	return instruction
}

func (cpu *CPU) ClearDisplay() {
	for x := 0; x < 64; x++ {
		for y := 0; y < 32; y++ {
			// Set to a black pixel.
			cpu.Display.SetRGBA(x, y, color.RGBA{A: 255})
		}
	}
}

func (cpu *CPU) Execute() {
	op := cpu.fetch()

	X := op >> 8 & 0xF
	Y := op >> 4 & 0xF
	NN := uint8(op & 0xFF)
	NNN := op & 0x0FFF

	switch {
	case op == 0x00E0: // 00E0: Clear the screen
		cpu.ClearDisplay()
		cpu.Draw = true
	case op == 0x00EE: // 00EE: Return from subroutine
		cpu.SP--
		cpu.PC = cpu.stack[cpu.SP]
	case op&0xF000 == 0x1000: // 1NNN: Jump to address NNN
		cpu.PC = NNN
	case op&0xF000 == 0x2000: // 2NNN: Execute subroutine at NNN
		cpu.stack[cpu.SP] = cpu.PC
		cpu.SP++
		cpu.PC = op & NNN
	case op&0xF000 == 0x3000: // 3XNN: Skip if VX == NN
		if cpu.V[X] == NN {
			cpu.PC += 2
		}
	case op&0xF000 == 0x4000: // 4XNN: Skip if VX != NN
		if cpu.V[X] != NN {
			cpu.PC += 2
		}
	case op&0xF00F == 0x5000: // 5XY0: Skip if VX == VY
		if cpu.V[X] == cpu.V[Y] {
			cpu.PC += 2
		}
	case op&0xF000 == 0x6000: // 6XNN: Store NN in VX
		cpu.V[X] = NN
	case op&0xF000 == 0x7000: // 7XNN: Add NN to VX
		cpu.V[X] += NN
	case op&0xF00F == 0x8000: // 8XY0: Store VY in VX
		cpu.V[X] = cpu.V[Y]
	case op&0xF00F == 0x8001: // 8XY1: Set VX to VX OR VY
		cpu.V[X] |= cpu.V[Y]
	case op&0xF00F == 0x8002: // 8XY2: Set VX to VX AND VY
		cpu.V[X] &= cpu.V[Y]
	case op&0xF00F == 0x8003: // 8XY3: Set VX to VX XOR VY
		cpu.V[X] ^= cpu.V[Y]
	case op&0xF00F == 0x8004: // 8XY4: Add VY to VX with carry in VF
		if cpu.V[Y] > 0xFF-cpu.V[X] {
			cpu.V[0xF] = 1
		} else {
			cpu.V[0xF] = 0
		}

		cpu.V[X] += cpu.V[Y]
	case op&0xF00F == 0x8005:
		if cpu.V[X] > cpu.V[Y] {
			cpu.V[0xF] = 1
		} else {
			cpu.V[0xF] = 0
		}

		cpu.V[X] -= cpu.V[Y]
	case op&0xF00F == 0x8006: // 8XY6: Store VX's LSB in VF and Shift right VX by 1
		cpu.V[0xF] = cpu.V[X] & 0x1
		cpu.V[X] >>= 1
	case op&0xF00F == 0x800E: // 8XYE: Store VX's MSB in VF and shift left VX by 1
		cpu.V[0xF] = cpu.V[X] & 0x80
		cpu.V[X] <<= 1
	case op&0xF00F == 0x9000: // 9XY0: Skip if VX != VY
		if cpu.V[X] != cpu.V[Y] {
			cpu.PC += 2
		}
	case op&0xF000 == 0xA000: // ANNN: Set I to NNN
		cpu.I = NNN
	case op&0xF000 == 0xC000: // CXNN: Set VX to random number with mask of NN
		cpu.V[X] = uint8(rand.Intn(255)) & NN
	case op&0xF000 == 0xD000: // DXYN: Draw a sprite
		// Draw a sprite at position VX, VY with N bytes of sprite data starting at
		// the address stored in I
		// Set VF to 01 if any set pixels are changed to unset, and 00 otherwise.
		x := int(cpu.V[X])
		y := int(cpu.V[Y])

		// The coordinates are supposed to wrap around the screen.
		// We do this by taking the modulo of the point by the axis length.
		if x > 63 {
			x %= 64
		}

		// Do the same for the y-axis.
		if y > 31 {
			y %= 32
		}

		height := int(op & 0x000F)
		// Reset the VF flag.
		cpu.V[0xF] = 0
		// A sprite consists of 8 pixels per-row
		// of upto 15 rows tall.
		// each row is just a byte with pixel data in each bit.
		for i := 0; i < height; i++ {
			row := cpu.Memory.RAM[cpu.I+uint16(i)]
			bit := 0

			for bit < 8 {
				// If this bit is on.
				if row&0x80 != 0 {
					// Grab the current pixel that is on the display.
					r, g, b, a := cpu.Display.At(x+bit, y+i).RGBA()

					// If the current pixel is on.
					if r == 0xFFFF && g == 0xFFFF && b == 0xFFFF && a == 0xFFFF {
						// Set the VF flag to indicate this pixel collision.
						cpu.V[0xF] = 1
						// Turn off the pixel.
						cpu.Display.SetRGBA(x+bit, y+i, color.RGBA{A: 255})
					} else {
						// Turn on the pixel.
						cpu.Display.SetRGBA(x+bit, y+i, color.RGBA{R: 255, G: 255, B: 255, A: 255})
					}
				}

				// Check the next bit.
				bit++
				row <<= 1
			}
		}
		// Set the draw flag to notify that the display has changed and must be redrawn.
		cpu.Draw = true
	case op&0xF0FF == 0xE09E: // EX9E: Skip if key in VX is pressed.
		if cpu.Keys.IsPressed(cpu.V[X]) {
			cpu.PC += 2
		}
	case op&0xF0FF == 0xE0A1: // EXA1: Skip if key in VX is not pressed.
		if !cpu.Keys.IsPressed(cpu.V[X]) {
			cpu.PC += 2
		}
	case op&0xF0FF == 0xF007: // FX07: Store value of delay timer in VX
		cpu.V[X] = cpu.DelayTimer
	case op&0xF0FF == 0xF00A: // FX01: Wait for keypress and store key in VX
		recieved := false
		for key, pressed := range cpu.Keys {
			if pressed {
				cpu.V[X] = uint8(key)
				recieved = true
			}
		}

		if !recieved {
			// If we didn't recieve a keypress, go back one cycle to retry.
			cpu.PC -= 2
		}
	case op&0xF0FF == 0xF015: // FX15: Set delay timer to VX
		cpu.DelayTimer = cpu.V[X]
	case op&0xF0FF == 0xF018: // FX18: Set sound timer to VX
		cpu.SoundTimer = cpu.V[X]
	case op&0xF0FF == 0xF01E: // FX1E: Add VX to I
		cpu.I += uint16(cpu.V[X])
	case op&0xF0FF == 0xF029: // FX29: Set I to sprite data for VX
		cpu.I = uint16(cpu.V[X] * 5)
	case op&0xF0FF == 0xF033: // FX33: Store BCD of VX at I, I + 1, I + 2
		cpu.Memory.RAM[cpu.I] = cpu.V[X] / 100
		cpu.Memory.RAM[cpu.I+1] = cpu.V[X] / 10 % 10
		cpu.Memory.RAM[cpu.I+2] = cpu.V[X] % 100 % 10
	case op&0xF0FF == 0xF055: // FX55: Store values of V0-VX in I
		for i := uint16(0); i <= X; i++ {
			cpu.Memory.RAM[cpu.I+i] = cpu.V[i]
		}
		// I is set to I + X + 1 after operation.
		cpu.I = cpu.I + X + 1
	case op&0xF0FF == 0xF065: // FX65: Fill V0-VX with values stored in I
		for i := uint16(0); i <= X; i++ {
			cpu.V[i] = cpu.Memory.RAM[cpu.I+i]
		}
		// I is set to I + X + 1 after operation.
		cpu.I = cpu.I + X + 1
	default:
		log.Fatalf("Unhandled instruction: 0x%X\n", op)
	}
}

func (cpu *CPU) UpdateSoundTimer() {
	if cpu.SoundTimer > 0 {
		cpu.SoundTimer--
		// TODO: Audio
		if cpu.SoundTimer == 0 {
			log.Println("BEEP")
		}
	}
}

func (cpu *CPU) UpdateDelayTimer() {
	if cpu.DelayTimer > 0 {
		cpu.DelayTimer--
	}
}

// Update the timers. This must be called at 60 Hz
func (cpu *CPU) UpdateTimers() {
	cpu.UpdateDelayTimer()
	cpu.UpdateSoundTimer()
}
