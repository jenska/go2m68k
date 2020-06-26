# M68K emulator

g02m68k will become an 680xx CPU emulator in Go language. 

## Project Status
go2m68k is a work-in-progress project; as such, expect bugs and a lot of missing features. 

## TODOS
- Perfomance improvements
- Documentation of the cpu builder
- Complete opcodes
- Include CPU cycle support
- Tracing feature (t0, t1)
- WaitGroups for bus access

## Build

Requires binutils-m68k-linux-gnu to run m68k-linux-gnu-as

### Performance

For performance measurement a bechmark is used that runs simple 68000 machine code.

```m68k
      moveq #100, d0
l0:   moveq #100, d1
l1:   dbra  d1, l0
      dbra  d0, l1
      stop  #2700
```

Using go*s profiler feature the perfomance leaks become obvious. 
Memory access consumes nearly half of the overall execution time.

```bash
go test -bench=. --cpuprofile cpu.prof
2020/05/02 07:40:34 added 2313 disassembler instructions
2020/05/02 07:40:34 added 2313 cpu instructions
2020/05/02 07:40:34 SR -----S---7 PC 00fc0030 USP 00000000 SSP 00001000
D0 00000000 D1 00000000 D2 00000000 D3 00000000 D4 00000000 D5 00000000 D6 00000000 D7 00000000 
A0 00000000 A1 00000000 A2 00000000 A3 00000000 A4 00000000 A5 00000000 A6 00000000 A7 00001000 

2020/05/02 07:40:34 20824 1
goos: linux
goarch: amd64
pkg: github.com/jenska/go2m68k
BenchmarkDbra-4   	2020/05/02 07:40:34 1061324 100
2020/05/02 07:40:35 22974254 2106
    2106	    506520 ns/op
PASS
2020/05/02 07:40:35 SR -----S---7 PC 00004010 USP 00000ffa SSP 00001000
D0 0000ffff D1 0000ffff D2 00000000 D3 00000000 D4 00000000 D5 00000000 D6 00000000 D7 00000000 
A0 00000000 A1 00000000 A2 00000000 A3 00000000 A4 00000000 A5 00000000 A6 00000000 A7 00001000 

ok  	github.com/jenska/go2m68k	1.343s

go tool pprof cpu.prof
File: go2m68k.test
Type: cpu
Time: May 2, 2020 at 7:40am (CEST)
Duration: 1.31s, Total samples = 1.14s (87.32%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 1010ms, 88.60% of 1140ms total
Showing top 10 nodes out of 35
      flat  flat%   sum%        cum   cum%
     220ms 19.30% 19.30%      730ms 64.04%  github.com/jenska/go2m68k.(*M68K).SetISA68000.func1
     180ms 15.79% 35.09%      510ms 44.74%  github.com/jenska/go2m68k.(*addressAreaQueue).read
     150ms 13.16% 48.25%      990ms 86.84%  github.com/jenska/go2m68k.(*M68K).step
     130ms 11.40% 59.65%      240ms 21.05%  github.com/jenska/go2m68k.NewBaseArea.func1
      90ms  7.89% 67.54%       90ms  7.89%  github.com/jenska/go2m68k.(*addressAreaQueue).findArea
      60ms  5.26% 72.81%       60ms  5.26%  encoding/binary.bigEndian.Uint16 (inline)
      50ms  4.39% 77.19%      110ms  9.65%  github.com/jenska/go2m68k.glob..func4
      50ms  4.39% 81.58%       50ms  4.39%  runtime.chanrecv
      40ms  3.51% 85.09%     1110ms 97.37%  github.com/jenska/go2m68k.(*M68K).Run
      40ms  3.51% 88.60%      770ms 67.54%  github.com/jenska/go2m68k.(*M68K).popPC (inline)
```

### Generate Disassmbly

Great for code review
```bash
go tool compile -S imisc.go cpu.go operands.go ea.go ssr.go builder.go bus.go  > cpu.asm
```
