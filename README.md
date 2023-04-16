# CHIP-8

A CHIP-8 emulator written in Go.

![screenshot](https://cdn.discordapp.com/attachments/976862121784655933/1097170511051628644/invaders.gif)

## Install

This package depends on:
- [github.com/go-gl/glfw](https://github.com/go-gl/glfw)
- [github.com/ravener/go-gl](https://github.com/ravener/go-gl)

You can use the `go install` command to easily fetch the code and automatically build it alongside the dependencies and add it to your `$GOPATH/bin`
```sh
go install github.com/ravener/chip8
```

## Usage

```sh
chip8 <file.rom>
```

Controls:

```asciidoc
Keypad             Keyboard
+-+-+-+-+          +-+-+-+-+
|1|2|3|C|          |1|2|3|4|
+-+-+-+-+          +-+-+-+-+
|4|5|6|D|          |Q|W|E|R|
+-+-+-+-+    =>    +-+-+-+-+
|7|8|9|E|          |A|S|D|F|
+-+-+-+-+          +-+-+-+-+
|A|0|B|F|          |Z|X|C|V|
+-+-+-+-+          +-+-+-+-+
```

- ESC to close the window.
- Space to pause the emulator.
- F11 to take a screenshot in the current directory.

## Resources

- [Writing a CHIP-8 interpreter - Ravener](https://ravener.vercel.app/posts/writing-a-chip8-interpreter) (my own blog post on this)
- [CHIP-8 Instruction Set](https://github.com/mattmikolay/chip-8/wiki/CHIP%E2%80%908-Instruction-Set)
- [CHIP-8 Technical Reference](https://github.com/mattmikolay/chip-8/wiki/CHIP%E2%80%908-Technical-Reference)

## TODO

- Fix timing. (Currently I cheated by relying on V-Sync)
- Audio output.

## License

[MIT License](LICENSE)
