package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	cpu000 = iota
	cpu070
	cpu010
	cpu020
	cpu030
	cpu040
	cpuFSCPU32
	cpuColdfire
	cpuCount
)

type (
	eaInfo struct {
		name  string
		mask  int
		value int
	}

	opcode struct {
		eaInfo
		body      string
		size      string
		eaAllowed string
		priv      []bool
		cycles    []int
	}

	opcodeHandler struct {
		eaInfo
		body   string
		bits   int
		cycles []int
	}
)

var (
	//                   000           010           020           030           040        FSCPU32      Coldfire
	eaCycleTable = map[string][8][3]int{
		"none": {{0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
		"ai":   {{0, 4, 8}, {0, 4, 8}, {0, 4, 8}, {0, 4, 4}, {0, 4, 4}, {0, 4, 4}, {0, 4, 4}, {0, 0, 0}},
		"pi":   {{0, 4, 8}, {0, 4, 8}, {0, 4, 8}, {0, 4, 4}, {0, 4, 4}, {0, 4, 4}, {0, 4, 4}, {0, 0, 0}},
		"pi7":  {{0, 4, 8}, {0, 4, 8}, {0, 4, 8}, {0, 4, 4}, {0, 4, 4}, {0, 4, 4}, {0, 4, 4}, {0, 0, 0}},
		"pd":   {{0, 6, 10}, {0, 6, 10}, {0, 6, 10}, {0, 5, 5}, {0, 5, 5}, {0, 5, 5}, {0, 5, 5}, {0, 0, 0}},
		"pd7":  {{0, 6, 10}, {0, 6, 10}, {0, 6, 10}, {0, 5, 5}, {0, 5, 5}, {0, 5, 5}, {0, 5, 5}, {0, 0, 0}},
		"di":   {{0, 8, 12}, {0, 8, 12}, {0, 8, 12}, {0, 5, 5}, {0, 5, 5}, {0, 5, 5}, {0, 5, 5}, {0, 0, 0}},
		"ix":   {{0, 10, 14}, {0, 10, 14}, {0, 10, 14}, {0, 7, 7}, {0, 7, 7}, {0, 7, 7}, {0, 7, 7}, {0, 0, 0}},
		"aw":   {{0, 8, 12}, {0, 8, 12}, {0, 8, 12}, {0, 4, 4}, {0, 4, 4}, {0, 4, 4}, {0, 4, 4}, {0, 0, 0}},
		"al":   {{0, 12, 16}, {0, 12, 16}, {0, 12, 16}, {0, 4, 4}, {0, 4, 4}, {0, 4, 4}, {0, 4, 4}, {0, 0, 0}},
		"pcdi": {{0, 8, 12}, {0, 8, 12}, {0, 8, 12}, {0, 5, 5}, {0, 5, 5}, {0, 5, 5}, {0, 5, 5}, {0, 0, 0}},
		"pcix": {{0, 10, 14}, {0, 10, 14}, {0, 10, 14}, {0, 7, 7}, {0, 7, 7}, {0, 7, 7}, {0, 7, 7}, {0, 0, 0}},
		"i":    {{0, 4, 8}, {0, 4, 8}, {0, 4, 8}, {0, 2, 4}, {0, 2, 4}, {0, 2, 4}, {0, 2, 4}, {0, 0, 0}},
	}

	jmpCycleMap   = map[string]int{"none": 0, "ai": 4, "pi": 0, "pi7": 0, "pd": 0, "pd7": 0, "di": 6, "ix": 10, "aw": 6, "al": 8, "pcdi": 6, "pcix": 10, "i": 0}
	leaCycleMap   = map[string]int{"none": 0, "ai": 4, "pi": 0, "pi7": 0, "pd": 0, "pd7": 0, "di": 8, "ix": 12, "aw": 8, "al": 12, "pcdi": 8, "pcix": 12, "i": 0}
	peaCycleMap   = map[string]int{"none": 0, "ai": 6, "pi": 0, "pi7": 0, "pd": 0, "pd7": 0, "di": 10, "ix": 14, "aw": 10, "al": 14, "pcdi": 10, "pcix": 14, "i": 0}
	movemCycleMap = map[string]int{"none": 0, "ai": 0, "pi": 0, "pi7": 0, "pd": 0, "pd7": 0, "di": 4, "ix": 6, "aw": 4, "al": 8, "pcdi": 0, "pcix": 0, "i": 0}
	movesCycleMap = map[string][3]int{
		"none": {0, 0, 0},
		"ai":   {0, 4, 6},
		"pi":   {0, 4, 6},
		"pi7":  {0, 4, 6},
		"pd":   {0, 6, 12},
		"pd7":  {0, 6, 12},
		"di":   {0, 12, 16},
		"ix":   {0, 16, 20},
		"aw":   {0, 12, 16},
		"al":   {0, 16, 20},
		"pcdi": {0, 0, 0},
		"pcix": {0, 0, 0},
		"i":    {0, 0, 0},
	}

	clrCycleTable = map[string][3]int{
		"none": {0, 0, 0},
		"ai":   {0, 4, 6},
		"pi":   {0, 4, 6},
		"pi7":  {0, 4, 6},
		"pd":   {0, 6, 8},
		"pd7":  {0, 6, 8},
		"di":   {0, 8, 10},
		"ix":   {0, 10, 14},
		"aw":   {0, 8, 10},
		"al":   {0, 10, 14},
		"pcdi": {0, 0, 0},
		"pcix": {0, 0, 0},
		"i":    {0, 0, 0},
	}

	eaInfoMap = map[string]eaInfo{
		"none": {"", 0x00, 0x00},
		"ai":   {"AyAi", 0x38, 0x10},
		"pi":   {"AyPi", 0x38, 0x18},
		"pi7":  {"A7Pi", 0x3f, 0x1f},
		"pd":   {"AyPd", 0x38, 0x20},
		"pd7":  {"A7Pd", 0x3f, 0x27},
		"di":   {"AyDi", 0x38, 0x28},
		"ix":   {"AyIx", 0x38, 0x30},
		"aw":   {"Aw", 0x3f, 0x38},
		"al":   {"Al", 0x3f, 0x39},
		"pcdi": {"PcDI", 0x3f, 0x3a},
		"pcix": {"PcIX", 0x3f, 0x3b},
		"i":    {"I", 0x3f, 0x3c},
	}

	ccTable  = []string{"t", "f", "hi", "ls", "cc", "cs", "ne", "eq", "vc", "vs", "pl", "mi", "ge", "lt", "gt", "le"}
	cpuNames = "071234fc"

	opcodes        = []*opcode{}
	opcodeHandlers = []opcodeHandler{}
)

func newOpcode(line string) *opcode {
	result := opcode{}
	entries := strings.Fields(line)
	result.value = toInt(entries[0], 16)
	result.mask = toInt(entries[1], 16)
	result.name = entries[2]
	result.size = entries[3]
	result.eaAllowed = entries[4]
	result.priv = make([]bool, cpuCount)
	result.cycles = make([]int, cpuCount)
	for i := 0; i < len(result.cycles); i++ {
		result.cycles[i] = -1
	}
	for _, entry := range entries[5:] {
		parts := strings.Split(entry, ":")
		ci := parts[1]
		priv := ci[len(ci)-1] == 'p'
		if priv {
			ci = ci[:len(ci)-1]
		}
		cycles := toInt(ci, 10)
		for _, c := range parts[0] {
			cpu := strings.IndexRune(cpuNames, c)
			result.cycles[cpu] = cycles
			result.priv[cpu] = priv
		}
	}
	return &result
}

func newOpcodeHandler(o opcode, eaMode string) opcodeHandler {
	result := opcodeHandler{}
	result.cycles = make([]int, cpuCount)
	for i := 0; i < len(result.cycles); i++ {
		result.cycles[i] = -1
	}
	sizeOrder := 2
	if o.size == "." {
		sizeOrder = 0
	} else if o.size == "b" || o.size == "w" {
		sizeOrder = 1
	}
	for i := 0; i < cpuCount; i++ {
		if o.cycles[i] < 0 {
			continue
		}
		if i == cpu010 && o.name == "moves" {
			result.cycles[i] = o.cycles[i] + movesCycleMap[eaMode][sizeOrder]
		} else if i == cpu010 && o.name == "clr" {
			result.cycles[i] = o.cycles[i] + clrCycleTable[eaMode][sizeOrder]
		} else if (i == cpu000 || i == cpu070) && (eaMode == "i" || eaMode == "none") && o.size == "l" &&
			((o.cycles[i] == 6 && (o.name == "add" || o.name == "and" || o.name == "or" || o.name == "sub")) || o.name == "adda" || o.name == "suba") {
			result.cycles[i] = o.cycles[i] + eaCycleTable[eaMode][i][sizeOrder] + 2
		} else if i < cpu020 && (o.name == "jmp" || o.name == "jsr") {
			result.cycles[i] = o.cycles[i] + jmpCycleMap[eaMode]
		} else if i < cpu020 && o.name == "lea" {
			result.cycles[i] = o.cycles[i] + leaCycleMap[eaMode]
		} else if i < cpu020 && o.name == "pea" {
			result.cycles[i] = o.cycles[i] + peaCycleMap[eaMode]
		} else if i < cpu020 && o.name == "movem" {
			result.cycles[i] = o.cycles[i] + movemCycleMap[eaMode]
		} else {
			result.cycles[i] = o.cycles[i] + eaCycleTable[eaMode][i][sizeOrder]
		}
	}
	result.value = o.value | eaInfoMap[eaMode].value
	result.mask = o.mask | eaInfoMap[eaMode].mask
	result.name = fmt.Sprintf("x%04x_%s", o.value, o.name)
	if o.size != "." {
		result.name += "_" + o.size
	}
	if eaMode != "none" {
		result.name += "_" + eaMode
	}
	result.name += "_"
	for i := 0; i < cpuCount; i++ {
		if result.cycles[i] >= 0 {
			result.name += string(cpuNames[i])
		}
	}

	result.bits = 0
	for i := 0; i < 16; i++ {
		if (result.mask & (1 << i)) != 0 {
			result.bits++
		}
	}
	if eaMode != "none" {
		n := eaInfoMap[eaMode].name
		body := o.body
		body = strings.ReplaceAll(body, "M68KMAKE_GET_EA_AY_8", "ea"+n+"(c, Byte)")
		body = strings.ReplaceAll(body, "M68KMAKE_GET_EA_AY_16", "ea"+n+"(c, Word)")
		body = strings.ReplaceAll(body, "M68KMAKE_GET_EA_AY_32", "ea"+n+"(c, Long)")
		body = strings.ReplaceAll(body, "M68KMAKE_GET_OPER_AY_8", "c.read(ea"+n+"(c, Byte))")
		body = strings.ReplaceAll(body, "M68KMAKE_GET_OPER_AY_16", "c.read(ea"+n+"(c, Word))")
		body = strings.ReplaceAll(body, "M68KMAKE_GET_OPER_AY_32", "c.read(ea"+n+"(c, Long))")
		result.body = body
	} else {
		result.body = o.body
	}
	return result
}

func (o *opcode) append(line string) {
	o.body += "\n" + line
}

func (o *opcode) generate() {
	if o.name == "bcc" || o.name == "scc" || o.name == "dbcc" || o.name == "trapcc" {
		o.ccVariants()
	} else {
		o.eaVariants()
	}
}

func (o opcode) ccVariants() {
	bname := o.name[:len(o.name)-2]
	for cc := 2; cc < len(ccTable); cc++ {
		o.name = bname + ccTable[cc]
		newString := "c.cc" + strings.ToUpper(ccTable[cc])
		o.body = strings.ReplaceAll(o.body, "M68KMAKE_CC", newString)
		newString = "c.ccNot" + strings.ToUpper(ccTable[cc])
		o.body = strings.ReplaceAll(o.body, "M68KMAKE_NOT_CC", newString)
		o.mask |= 0x0f00
		o.value = (o.value & 0x0f00) | (cc << 8)
		o.eaVariants()
	}
}

func (o opcode) eaVariants() {
	allowed := o.eaAllowed
	if allowed == "." {
		opcodeHandlers = append(opcodeHandlers, newOpcodeHandler(o, "none"))
		return
	}
	if strings.ContainsRune(allowed, 'A') {
		opcodeHandlers = append(opcodeHandlers, newOpcodeHandler(o, "ai"))
	}
	if strings.ContainsRune(allowed, '+') {
		opcodeHandlers = append(opcodeHandlers, newOpcodeHandler(o, "pi"))
		if o.size == "b" {
			opcodeHandlers = append(opcodeHandlers, newOpcodeHandler(o, "pi7"))
		}
	}
	if strings.ContainsRune(allowed, '-') {
		opcodeHandlers = append(opcodeHandlers, newOpcodeHandler(o, "pd"))
		if o.size == "b" {
			opcodeHandlers = append(opcodeHandlers, newOpcodeHandler(o, "pd7"))
		}
	}
	if strings.ContainsRune(allowed, 'D') {
		opcodeHandlers = append(opcodeHandlers, newOpcodeHandler(o, "di"))
	}
	if strings.ContainsRune(allowed, 'X') {
		opcodeHandlers = append(opcodeHandlers, newOpcodeHandler(o, "ix"))
	}
	if strings.ContainsRune(allowed, 'W') {
		opcodeHandlers = append(opcodeHandlers, newOpcodeHandler(o, "aw"))
	}
	if strings.ContainsRune(allowed, 'L') {
		opcodeHandlers = append(opcodeHandlers, newOpcodeHandler(o, "al"))
	}
	if strings.ContainsRune(allowed, 'd') {
		opcodeHandlers = append(opcodeHandlers, newOpcodeHandler(o, "pcdi"))
	}
	if strings.ContainsRune(allowed, 'x') {
		opcodeHandlers = append(opcodeHandlers, newOpcodeHandler(o, "pcix"))
	}
	if strings.ContainsRune(allowed, 'I') {
		opcodeHandlers = append(opcodeHandlers, newOpcodeHandler(o, "i"))
	}
}

func toInt(s string, base int) int {
	if v, x := strconv.ParseInt(s, base, 64); x == nil {
		return int(v)
	} else {
		panic(x)
	}
}

func main() {
	cmd := os.Args[0]
	input, err := os.Open("./gen/m68kin.txt")
	if err != nil {
		panic(fmt.Errorf("Cannot open file m68kin.txt {%s}", err))
	}
	defer input.Close()

	scanner := bufio.NewScanner(input)
	var opcode *opcode
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || line[0] == ' ' || line[0] == '\t' {
			if opcode != nil {
				opcode.append(line)
			}
		} else if line[0] != '#' {
			opcode = newOpcode(line)
			opcodes = append(opcodes, opcode)
		}
	}
	for _, opcode := range opcodes {
		opcode.generate()
	}

	var output strings.Builder
	fmt.Fprintf(&output, "package cpu\n")
	fmt.Fprintf(&output, "// Generated source, edits will be lost. Run %s instead\n\n", cmd)
	for _, oh := range opcodeHandlers {
		fmt.Fprintf(&output, "func %s(c *M68K) {%s}\n\n", oh.name, oh.body)
	}

	ioutil.WriteFile("instructions.go", []byte(output.String()), 0644)
}
