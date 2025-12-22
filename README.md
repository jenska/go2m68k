# go2m68k

> **⚠️ PROJECT DISCONTINUED**
>
> This project is no longer maintained.
>
> **Please use the new emulator: [m68kemu](https://github.com/jenska/m68kemu)**
>
> `m68kemu` is the designated successor to this project. It provides a newer, improved implementation for Motorola 68000 emulation. We highly recommend using `m68kemu` for all future needs.

Fresh start for a Motorola 68000 emulator in Go, inspired by [Musashi](https://github.com/kstenerud/Musashi). The goal is a clean API that can execute individual instructions and integrate with tests generated through the `m68kasm` assembler API.

## Project goals
- Minimal, extensible CPU core with clear data structures
- Execute single instructions through a public API
- Ready for tests driven by `m68kasm`
- MIT licensed and open to contributions

## Requirements
- Go 1.22 or newer

## Quickstart
1. Install the Go toolchain.
2. Run the test suite with the Makefile helper:
   ```bash
   make test
   ```
3. Load program bytes, register the built-in instructions, and step through them:
   ```go
   cpu := cpu.New(1024)
   _ = instructions.RegisterDefaults(cpu)
   _ = cpu.LoadProgram(0x100, []byte{0x4e, 0x71}) // NOP
   cpu.Reset(0x100)
   _ = cpu.Step()
   ```

## Makefile tasks
The repository includes a small Makefile to streamline local development:

- `make fmt` — format the codebase with `gofmt`.
- `make vet` — run `go vet` for static analysis.
- `make test` — execute the full test suite.
- `make tidy` — ensure module dependencies are tidy.
- `make all` — run formatting, vetting, and tests in sequence.

## Architecture
- `internal/cpu`: CPU state, registers, memory, and instruction dispatch.
- `internal/instructions`: Built-in instruction implementations grouped by function (control, arithmetic, move).
- Register custom instructions with `RegisterInstruction` and run them directly via `ExecuteInstruction`.
- `Memory.ReadWord` and `Memory.WriteWord` simplify instruction implementations.

## Testing with m68kasm
Instruction handlers can be tested with machine code assembled through `m68kasm`. Assemble code, load it via `LoadProgram`, and step through instructions to validate behavior.

## License
Released under the MIT License. See `LICENSE` for details.
