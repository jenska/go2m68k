package m68k

/* Memory interface to implement read/write functions called by the CPU.
 * while values used are 32 bits, only the appropriate number
 * of bits are relevant (i.e. in write_memory_8, only the lower 8 bits
 * of value should be written to memory).
 */
type Memory interface {
	/* Read from anywhere */
	ReadMemoryByte(address Address) Byte
	ReadMemoryWord(address Address) Word
	ReadMemoryLong(address Address) Long

	/* Write to anywhere */
	WriteMemoryByte(address Address, value Byte)
	WriteMemoryWord(address Address, value Word)
	WriteMemoryLong(uaddress Address, value Long)
}

/* Memory access for the disassembler */
type DisassemblerMemroy interface {
	ReadMemoryByte(address Address) Byte
	ReadMemoryWord(address Address) Word
	ReadMemoryLong(address Address) Long
}
