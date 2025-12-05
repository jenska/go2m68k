# go2m68k

Fresh start for a Motorola 68000 emulator in Go, inspired by [Musashi](https://github.com/kstenerud/Musashi). The goal is a clean API that can execute individual instructions and integrate with tests generated through the `m68kasm` assembler API.

## Goals
- Minimal, extensible CPU core with clear data structures
- Execute single instructions through a public API
- Ready for tests driven by `m68kasm`
- MIT licensed and open to contributions

## Quickstart
1. Install Go 1.22.
2. Run the test suite:
   ```bash
   go test ./...
   ```
3. Load program bytes, register the built-in instructions, and step through them:
   ```go
   cpu := cpu.New(1024)
   _ = instructions.RegisterDefaults(cpu)
   _ = cpu.LoadProgram(0x100, []byte{0x4e, 0x71}) // NOP
   cpu.Reset(0x100)
   _ = cpu.Step()
   ```

## Architecture
- `internal/cpu`: CPU state, registers, memory, and instruction dispatch.
- `internal/instructions`: Built-in instruction implementations grouped by function (control, arithmetic, move).
- Register custom instructions with `RegisterInstruction` and run them directly via `ExecuteInstruction`.
- `Memory.ReadWord` and `Memory.WriteWord` simplify instruction implementations.

## Testing with m68kasm
Instruction handlers can be tested with machine code assembled through `m68kasm`. Assemble code, load it via `LoadProgram`, and step through instructions to validate behavior.

## License
Released under the MIT License. See `LICENSE` for details.
