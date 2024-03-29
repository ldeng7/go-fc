package core

const (
	cpuRegC byte = 0x01 << iota
	cpuRegZ
	cpuRegI
	cpuRegD
	cpuRegB
	cpuRegR
	cpuRegV
	cpuRegN
)

const (
	cpuIntrTypNmi byte = 0x01 << iota
	cpuIntrTypIrq
	cpuIntrTypFrame
	cpuIntrTypDpcm
	cpuIntrTypMapper
	cpuIntrTypTrig
	cpuIntrTypTrig2
)

type Cpu struct {
	sys *Sys
	ram []byte

	regPC uint16
	regA  byte
	regX  byte
	regY  byte
	regP  byte
	regS  byte
	intr  byte

	nCycle    int64
	nCycleDma int64
	znTable   [256]byte
}

func newCpu(sys *Sys) *Cpu {
	cpu := &Cpu{}
	cpu.sys = sys
	cpu.ram = sys.mem.ram[:]

	cpu.znTable[0] = cpuRegZ
	for i := 1; i < 256; i++ {
		if i&0x80 == 0 {
			cpu.znTable[i] = 0
		} else {
			cpu.znTable[i] = cpuRegN
		}
	}

	return cpu
}

func (cpu *Cpu) reset() {
	cpu.regPC = cpu.readWord(0xfffc)
	cpu.regA, cpu.regX, cpu.regY, cpu.regS = 0, 0, 0, 0xff
	cpu.regP = cpuRegZ | cpuRegR
	cpu.intr = 0
	cpu.nCycle, cpu.nCycleDma = 0, 0
}

func (cpu *Cpu) readWord(addr uint16) uint16 {
	return uint16(cpu.sys.read(addr)) | (uint16(cpu.sys.read(addr+1)) << 8)
}

func (cpu *Cpu) run(nCycleReq int64) int64 {
	nCyclePrev := cpu.nCycle
	sys := cpu.sys
	for nCycleReq > 0 {
		var nCycleExec int64 = 0
		if cpu.nCycleDma != 0 {
			if cpu.nCycleDma >= nCycleReq {
				cpu.nCycleDma -= nCycleReq
				cpu.nCycle += nCycleReq
				cpu.sys.mapper.clock(nCycleReq)
				return cpu.nCycle - nCyclePrev
			} else {
				nCycleExec += cpu.nCycleDma
				cpu.nCycleDma = 0
			}
		}

		opcode := sys.read(cpu.regPC)
		cpu.regPC++
		intrNmi, intrIrq := false, false
		if cpu.intr&cpuIntrTypNmi != 0 {
			intrNmi = true
			cpu.intr &^= cpuIntrTypNmi
		} else if cpu.intr&0xfc != 0 {
			if (cpu.regP&cpuRegI == 0) && opcode != 0x40 {
				intrIrq = true
				cpu.intr &^= cpuIntrTypTrig
			}
			cpu.intr &^= cpuIntrTypTrig2
		}

		var data, data1 byte
		var addr, addr1, word uint16
		switch opcode {
		case 0x69:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			word = uint16(cpu.regA) + uint16(data) + uint16(cpu.regP&cpuRegC)
			cpu.regP &^= cpuRegC | cpuRegV | cpuRegZ | cpuRegN
			if word > 0x00ff {
				cpu.regP |= cpuRegC
			}
			if (^(cpu.regA ^ data))&(cpu.regA^byte(word))&0x80 != 0 {
				cpu.regP |= cpuRegV
			}
			cpu.regA = byte(word)
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 2
		case 0x65:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			data = cpu.ram[data]
			word = uint16(cpu.regA) + uint16(data) + uint16(cpu.regP&cpuRegC)
			cpu.regP &^= cpuRegC | cpuRegV | cpuRegZ | cpuRegN
			if word > 0x00ff {
				cpu.regP |= cpuRegC
			}
			if (^(cpu.regA ^ data))&(cpu.regA^byte(word))&0x80 != 0 {
				cpu.regP |= cpuRegV
			}
			cpu.regA = byte(word)
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 3
		case 0x75:
			data = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			data = cpu.ram[data]
			word = uint16(cpu.regA) + uint16(data) + uint16(cpu.regP&cpuRegC)
			cpu.regP &^= cpuRegC | cpuRegV | cpuRegZ | cpuRegN
			if word > 0x00ff {
				cpu.regP |= cpuRegC
			}
			if (^(cpu.regA ^ data))&(cpu.regA^byte(word))&0x80 != 0 {
				cpu.regP |= cpuRegV
			}
			cpu.regA = byte(word)
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 4
		case 0x6d:
			addr = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			data = sys.read(addr)
			word = uint16(cpu.regA) + uint16(data) + uint16(cpu.regP&cpuRegC)
			cpu.regP &^= cpuRegC | cpuRegV | cpuRegZ | cpuRegN
			if word > 0x00ff {
				cpu.regP |= cpuRegC
			}
			if (^(cpu.regA ^ data))&(cpu.regA^byte(word))&0x80 != 0 {
				cpu.regP |= cpuRegV
			}
			cpu.regA = byte(word)
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 4
		case 0x7d:
			addr1 = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			addr = addr1 + uint16(cpu.regX)
			data = sys.read(addr)
			word = uint16(cpu.regA) + uint16(data) + uint16(cpu.regP&cpuRegC)
			cpu.regP &^= cpuRegC | cpuRegV | cpuRegZ | cpuRegN
			if word > 0x00ff {
				cpu.regP |= cpuRegC
			}
			if (^(cpu.regA ^ data))&(cpu.regA^byte(word))&0x80 != 0 {
				cpu.regP |= cpuRegV
			}
			cpu.regA = byte(word)
			cpu.regP |= cpu.znTable[cpu.regA]
			if (addr1 & 0xff00) != (addr & 0xff00) {
				nCycleExec++
			}
			nCycleExec += 4
		case 0x79:
			addr1 = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			addr = addr1 + uint16(cpu.regY)
			data = sys.read(addr)
			word = uint16(cpu.regA) + uint16(data) + uint16(cpu.regP&cpuRegC)
			cpu.regP &^= cpuRegC | cpuRegV | cpuRegZ | cpuRegN
			if word > 0x00ff {
				cpu.regP |= cpuRegC
			}
			if (^(cpu.regA ^ data))&(cpu.regA^byte(word))&0x80 != 0 {
				cpu.regP |= cpuRegV
			}
			cpu.regA = byte(word)
			cpu.regP |= cpu.znTable[cpu.regA]
			if (addr1 & 0xff00) != (addr & 0xff00) {
				nCycleExec++
			}
			nCycleExec += 4
		case 0x61:
			data = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			addr = (uint16(cpu.ram[data+1]) << 8) | uint16(cpu.ram[data])
			data = sys.read(addr)
			word = uint16(cpu.regA) + uint16(data) + uint16(cpu.regP&cpuRegC)
			cpu.regP &^= cpuRegC | cpuRegV | cpuRegZ | cpuRegN
			if word > 0x00ff {
				cpu.regP |= cpuRegC
			}
			if (^(cpu.regA ^ data))&(cpu.regA^byte(word))&0x80 != 0 {
				cpu.regP |= cpuRegV
			}
			cpu.regA = byte(word)
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 6
		case 0x71:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			addr1 = (uint16(cpu.ram[data+1]) << 8) | uint16(cpu.ram[data])
			addr = addr1 + uint16(cpu.regY)
			data = sys.read(addr)
			word = uint16(cpu.regA) + uint16(data) + uint16(cpu.regP&cpuRegC)
			cpu.regP &^= cpuRegC | cpuRegV | cpuRegZ | cpuRegN
			if word > 0x00ff {
				cpu.regP |= cpuRegC
			}
			if (^(cpu.regA ^ data))&(cpu.regA^byte(word))&0x80 != 0 {
				cpu.regP |= cpuRegV
			}
			cpu.regA = byte(word)
			cpu.regP |= cpu.znTable[cpu.regA]
			if (addr1 & 0xff00) != (addr & 0xff00) {
				nCycleExec++
			}
			nCycleExec += 4
		case 0xe9:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			word = uint16(cpu.regA) - uint16(data) - uint16(^cpu.regP&cpuRegC)
			cpu.regP &^= cpuRegC | cpuRegV | cpuRegZ | cpuRegN
			if word < 0x0100 {
				cpu.regP |= cpuRegC
			}
			if (cpu.regA^data)&(cpu.regA^byte(word))&0x80 != 0 {
				cpu.regP |= cpuRegV
			}
			cpu.regA = byte(word)
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 2
		case 0xe5:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			data = cpu.ram[data]
			word = uint16(cpu.regA) - uint16(data) - uint16(^cpu.regP&cpuRegC)
			cpu.regP &^= cpuRegC | cpuRegV | cpuRegZ | cpuRegN
			if word < 0x0100 {
				cpu.regP |= cpuRegC
			}
			if (cpu.regA^data)&(cpu.regA^byte(word))&0x80 != 0 {
				cpu.regP |= cpuRegV
			}
			cpu.regA = byte(word)
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 3
		case 0xf5:
			data = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			data = cpu.ram[data]
			word = uint16(cpu.regA) - uint16(data) - uint16(^cpu.regP&cpuRegC)
			cpu.regP &^= cpuRegC | cpuRegV | cpuRegZ | cpuRegN
			if word < 0x0100 {
				cpu.regP |= cpuRegC
			}
			if (cpu.regA^data)&(cpu.regA^byte(word))&0x80 != 0 {
				cpu.regP |= cpuRegV
			}
			cpu.regA = byte(word)
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 4
		case 0xed:
			addr = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			data = sys.read(addr)
			word = uint16(cpu.regA) - uint16(data) - uint16(^cpu.regP&cpuRegC)
			cpu.regP &^= cpuRegC | cpuRegV | cpuRegZ | cpuRegN
			if word < 0x0100 {
				cpu.regP |= cpuRegC
			}
			if (cpu.regA^data)&(cpu.regA^byte(word))&0x80 != 0 {
				cpu.regP |= cpuRegV
			}
			cpu.regA = byte(word)
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 4
		case 0xfd:
			addr1 = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			addr = addr1 + uint16(cpu.regX)
			data = sys.read(addr)
			word = uint16(cpu.regA) - uint16(data) - uint16(^cpu.regP&cpuRegC)
			cpu.regP &^= cpuRegC | cpuRegV | cpuRegZ | cpuRegN
			if word < 0x0100 {
				cpu.regP |= cpuRegC
			}
			if (cpu.regA^data)&(cpu.regA^byte(word))&0x80 != 0 {
				cpu.regP |= cpuRegV
			}
			cpu.regA = byte(word)
			cpu.regP |= cpu.znTable[cpu.regA]
			if (addr1 & 0xff00) != (addr & 0xff00) {
				nCycleExec++
			}
			nCycleExec += 4
		case 0xf9:
			addr1 = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			addr = addr1 + uint16(cpu.regY)
			data = sys.read(addr)
			word = uint16(cpu.regA) - uint16(data) - uint16(^cpu.regP&cpuRegC)
			cpu.regP &^= cpuRegC | cpuRegV | cpuRegZ | cpuRegN
			if word < 0x0100 {
				cpu.regP |= cpuRegC
			}
			if (cpu.regA^data)&(cpu.regA^byte(word))&0x80 != 0 {
				cpu.regP |= cpuRegV
			}
			cpu.regA = byte(word)
			cpu.regP |= cpu.znTable[cpu.regA]
			if (addr1 & 0xff00) != (addr & 0xff00) {
				nCycleExec++
			}
			nCycleExec += 4
		case 0xe1:
			data = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			addr = (uint16(cpu.ram[data+1]) << 8) | uint16(cpu.ram[data])
			data = sys.read(addr)
			word = uint16(cpu.regA) - uint16(data) - uint16(^cpu.regP&cpuRegC)
			cpu.regP &^= cpuRegC | cpuRegV | cpuRegZ | cpuRegN
			if word < 0x0100 {
				cpu.regP |= cpuRegC
			}
			if (cpu.regA^data)&(cpu.regA^byte(word))&0x80 != 0 {
				cpu.regP |= cpuRegV
			}
			cpu.regA = byte(word)
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 6
		case 0xf1:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			addr1 = (uint16(cpu.ram[data+1]) << 8) | uint16(cpu.ram[data])
			addr = addr1 + uint16(cpu.regY)
			data = sys.read(addr)
			word = uint16(cpu.regA) - uint16(data) - uint16(^cpu.regP&cpuRegC)
			cpu.regP &^= cpuRegC | cpuRegV | cpuRegZ | cpuRegN
			if word < 0x0100 {
				cpu.regP |= cpuRegC
			}
			if (cpu.regA^data)&(cpu.regA^byte(word))&0x80 != 0 {
				cpu.regP |= cpuRegV
			}
			cpu.regA = byte(word)
			cpu.regP |= cpu.znTable[cpu.regA]
			if (addr1 & 0xff00) != (addr & 0xff00) {
				nCycleExec++
			}
			nCycleExec += 5
		case 0xc6:
			data1 = sys.read(cpu.regPC)
			cpu.regPC++
			data = cpu.ram[data1] - 1
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[data]
			cpu.ram[data1] = data
			nCycleExec += 5
		case 0xd6:
			data1 = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			data = cpu.ram[data1] - 1
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[data]
			cpu.ram[data1] = data
			nCycleExec += 6
		case 0xce:
			addr = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			data = sys.read(addr) - 1
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[data]
			sys.write(addr, data)
			nCycleExec += 6
		case 0xde:
			addr = cpu.readWord(cpu.regPC) + uint16(cpu.regX)
			cpu.regPC += 2
			data = sys.read(addr) - 1
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[data]
			sys.write(addr, data)
			nCycleExec += 7
		case 0xca:
			cpu.regX--
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regX]
			nCycleExec += 2
		case 0x88:
			cpu.regY--
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regY]
			nCycleExec += 2
		case 0xe6:
			data1 = sys.read(cpu.regPC)
			cpu.regPC++
			data = cpu.ram[data1] + 1
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[data]
			cpu.ram[data1] = data
			nCycleExec += 5
		case 0xf6:
			data1 = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			data = cpu.ram[data1] + 1
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[data]
			cpu.ram[data1] = data
			nCycleExec += 6
		case 0xee:
			addr = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			data = sys.read(addr) + 1
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[data]
			sys.write(addr, data)
			nCycleExec += 6
		case 0xfe:
			addr = cpu.readWord(cpu.regPC) + uint16(cpu.regX)
			cpu.regPC += 2
			data = sys.read(addr) + 1
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[data]
			sys.write(addr, data)
			nCycleExec += 7
		case 0xe8:
			cpu.regX++
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regX]
			nCycleExec += 2
		case 0xC8:
			cpu.regY++
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regY]
			nCycleExec += 2
		case 0x29:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			cpu.regA &= data
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 2
		case 0x25:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			cpu.regA &= cpu.ram[data]
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 3
		case 0x35:
			data = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			cpu.regA &= cpu.ram[data]
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 4
		case 0x2d:
			addr = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			cpu.regA &= sys.read(addr)
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 4
		case 0x3d:
			addr1 = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			addr = addr1 + uint16(cpu.regX)
			cpu.regA &= sys.read(addr)
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			if (addr1 & 0xff00) != (addr & 0xff00) {
				nCycleExec++
			}
			nCycleExec += 4
		case 0x39:
			addr1 = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			addr = addr1 + uint16(cpu.regY)
			cpu.regA &= sys.read(addr)
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			if (addr1 & 0xff00) != (addr & 0xff00) {
				nCycleExec++
			}
			nCycleExec += 4
		case 0x21:
			data = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			addr = (uint16(cpu.ram[data+1]) << 8) | uint16(cpu.ram[data])
			cpu.regA &= sys.read(addr)
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 6
		case 0x31:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			addr1 = (uint16(cpu.ram[data+1]) << 8) | uint16(cpu.ram[data])
			addr = addr1 + uint16(cpu.regY)
			cpu.regA &= sys.read(addr)
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			if (addr1 & 0xff00) != (addr & 0xff00) {
				nCycleExec++
			}
			nCycleExec += 5
		case 0x0a:
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if cpu.regA&0x80 != 0 {
				cpu.regP |= cpuRegC
			}
			cpu.regA <<= 1
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 2
		case 0x06:
			data1 = sys.read(cpu.regPC)
			cpu.regPC++
			data = cpu.ram[data1]
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x80 != 0 {
				cpu.regP |= cpuRegC
			}
			data <<= 1
			cpu.regP |= cpu.znTable[data]
			cpu.ram[data1] = data
			nCycleExec += 5
		case 0x16:
			data1 = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			data = cpu.ram[data1]
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x80 != 0 {
				cpu.regP |= cpuRegC
			}
			data <<= 1
			cpu.regP |= cpu.znTable[data]
			cpu.ram[data1] = data
			nCycleExec += 6
		case 0x0e:
			addr = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			data = sys.read(addr)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x80 != 0 {
				cpu.regP |= cpuRegC
			}
			data <<= 1
			cpu.regP |= cpu.znTable[data]
			sys.write(addr, data)
			nCycleExec += 6
		case 0x1e:
			addr = cpu.readWord(cpu.regPC) + uint16(cpu.regX)
			cpu.regPC += 2
			data = sys.read(addr)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x80 != 0 {
				cpu.regP |= cpuRegC
			}
			data <<= 1
			cpu.regP |= cpu.znTable[data]
			sys.write(addr, data)
			nCycleExec += 7
		case 0x24:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			data = cpu.ram[data]
			cpu.regP &^= cpuRegV | cpuRegZ | cpuRegN
			if data&0x40 != 0 {
				cpu.regP |= cpuRegV
			}
			if data&cpu.regA == 0 {
				cpu.regP |= cpuRegZ
			}
			if data&0x80 != 0 {
				cpu.regP |= cpuRegN
			}
			nCycleExec += 3
		case 0x2c:
			addr = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			data = sys.read(addr)
			cpu.regP &^= cpuRegV | cpuRegZ | cpuRegN
			if data&0x40 != 0 {
				cpu.regP |= cpuRegV
			}
			if data&cpu.regA == 0 {
				cpu.regP |= cpuRegZ
			}
			if data&0x80 != 0 {
				cpu.regP |= cpuRegN
			}
			nCycleExec += 4
		case 0x49:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			cpu.regA ^= data
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 2
		case 0x45:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			cpu.regA ^= cpu.ram[data]
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 3
		case 0x55:
			data = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			cpu.regA ^= cpu.ram[data]
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 4
		case 0x4d:
			addr = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			cpu.regA ^= sys.read(addr)
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 4
		case 0x5d:
			addr1 = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			addr = addr1 + uint16(cpu.regX)
			cpu.regA ^= sys.read(addr)
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			if (addr1 & 0xff00) != (addr & 0xff00) {
				nCycleExec++
			}
			nCycleExec += 4
		case 0x59:
			addr1 = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			addr = addr1 + uint16(cpu.regY)
			cpu.regA ^= sys.read(addr)
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			if (addr1 & 0xff00) != (addr & 0xff00) {
				nCycleExec++
			}
			nCycleExec += 4
		case 0x41:
			data = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			addr = (uint16(cpu.ram[data+1]) << 8) | uint16(cpu.ram[data])
			cpu.regA ^= sys.read(addr)
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 6
		case 0x51:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			addr1 = (uint16(cpu.ram[data+1]) << 8) | uint16(cpu.ram[data])
			addr = addr1 + uint16(cpu.regY)
			cpu.regA ^= sys.read(addr)
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			if (addr1 & 0xff00) != (addr & 0xff00) {
				nCycleExec++
			}
			nCycleExec += 5
		case 0x4a:
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if cpu.regA&0x01 != 0 {
				cpu.regP |= cpuRegC
			}
			cpu.regA >>= 1
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 2
		case 0x46:
			data1 = sys.read(cpu.regPC)
			cpu.regPC++
			data = cpu.ram[data1]
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x01 != 0 {
				cpu.regP |= cpuRegC
			}
			data >>= 1
			cpu.regP |= cpu.znTable[data]
			cpu.ram[data1] = data
			nCycleExec += 5
		case 0x56:
			data1 = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			data = cpu.ram[data1]
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x01 != 0 {
				cpu.regP |= cpuRegC
			}
			data >>= 1
			cpu.regP |= cpu.znTable[data]
			cpu.ram[data1] = data
			nCycleExec += (6)
		case 0x4e:
			addr = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			data = sys.read(addr)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x01 != 0 {
				cpu.regP |= cpuRegC
			}
			data >>= 1
			cpu.regP |= cpu.znTable[data]
			sys.write(addr, data)
			nCycleExec += 6
		case 0x5e:
			addr = cpu.readWord(cpu.regPC) + uint16(cpu.regX)
			cpu.regPC += 2
			data = sys.read(addr)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x01 != 0 {
				cpu.regP |= cpuRegC
			}
			data >>= 1
			cpu.regP |= cpu.znTable[data]
			sys.write(addr, data)
			nCycleExec += 7
		case 0x09:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			cpu.regA |= data
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 2
		case 0x05:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			cpu.regA |= cpu.ram[data]
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 3
		case 0x15:
			data = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			cpu.regA |= cpu.ram[data]
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 4
		case 0x0d:
			addr = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			cpu.regA |= sys.read(addr)
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 4
		case 0x1d:
			addr1 = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			addr = addr1 + uint16(cpu.regX)
			cpu.regA |= sys.read(addr)
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			if (addr1 & 0xff00) != (addr & 0xff00) {
				nCycleExec++
			}
			nCycleExec += 4
		case 0x19:
			addr1 = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			addr = addr1 + uint16(cpu.regY)
			cpu.regA |= sys.read(addr)
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			if (addr1 & 0xff00) != (addr & 0xff00) {
				nCycleExec++
			}
			nCycleExec += 4
		case 0x01:
			data = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			addr = (uint16(cpu.ram[data+1]) << 8) | uint16(cpu.ram[data])
			cpu.regA |= sys.read(addr)
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 6
		case 0x11:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			addr1 = (uint16(cpu.ram[data+1]) << 8) | uint16(cpu.ram[data])
			addr = addr1 + uint16(cpu.regY)
			cpu.regA |= sys.read(addr)
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			if (addr1 & 0xff00) != (addr & 0xff00) {
				nCycleExec++
			}
			nCycleExec += 5
		case 0x2a:
			b := cpu.regP&cpuRegC != 0
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if cpu.regA&0x80 != 0 {
				cpu.regP |= cpuRegC
			}
			cpu.regA <<= 1
			if b {
				cpu.regA |= 0x01
			}
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 2
		case 0x26:
			data1 = sys.read(cpu.regPC)
			cpu.regPC++
			data = cpu.ram[data1]
			b := cpu.regP&cpuRegC != 0
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x80 != 0 {
				cpu.regP |= cpuRegC
			}
			data <<= 1
			if b {
				data |= 0x01
			}
			cpu.regP |= cpu.znTable[data]
			cpu.ram[data1] = data
			nCycleExec += 5
		case 0x36:
			data1 = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			data = cpu.ram[data1]
			b := cpu.regP&cpuRegC != 0
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x80 != 0 {
				cpu.regP |= cpuRegC
			}
			data <<= 1
			if b {
				data |= 0x01
			}
			cpu.regP |= cpu.znTable[data]
			cpu.ram[data1] = data
			nCycleExec += 6
		case 0x2e:
			addr = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			data = sys.read(addr)
			b := cpu.regP&cpuRegC != 0
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x80 != 0 {
				cpu.regP |= cpuRegC
			}
			data <<= 1
			if b {
				data |= 0x01
			}
			cpu.regP |= cpu.znTable[data]
			sys.write(addr, data)
			nCycleExec += 6
		case 0x3e:
			addr = cpu.readWord(cpu.regPC) + uint16(cpu.regX)
			cpu.regPC += 2
			data = sys.read(addr)
			b := cpu.regP&cpuRegC != 0
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x80 != 0 {
				cpu.regP |= cpuRegC
			}
			data <<= 1
			if b {
				data |= 0x01
			}
			cpu.regP |= cpu.znTable[data]
			sys.write(addr, data)
			nCycleExec += 7
		case 0x6a:
			b := cpu.regP&cpuRegC != 0
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if cpu.regA&0x01 != 0 {
				cpu.regP |= cpuRegC
			}
			cpu.regA >>= 1
			if b {
				cpu.regA |= 0x80
			}
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 2
		case 0x66:
			data1 = sys.read(cpu.regPC)
			cpu.regPC++
			data = cpu.ram[data1]
			b := cpu.regP&cpuRegC != 0
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x01 != 0 {
				cpu.regP |= cpuRegC
			}
			data >>= 1
			if b {
				data |= 0x80
			}
			cpu.regP |= cpu.znTable[data]
			cpu.ram[data1] = data
			nCycleExec += 5
		case 0x76:
			data1 = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			data = cpu.ram[data1]
			b := cpu.regP&cpuRegC != 0
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x01 != 0 {
				cpu.regP |= cpuRegC
			}
			data >>= 1
			if b {
				data |= 0x80
			}
			cpu.regP |= cpu.znTable[data]
			cpu.ram[data1] = data
			nCycleExec += 6
		case 0x6e:
			addr = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			data = sys.read(addr)
			b := cpu.regP&cpuRegC != 0
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x01 != 0 {
				cpu.regP |= cpuRegC
			}
			data >>= 1
			if b {
				data |= 0x80
			}
			cpu.regP |= cpu.znTable[data]
			sys.write(addr, data)
			nCycleExec += 6
		case 0x7e:
			addr = cpu.readWord(cpu.regPC) + uint16(cpu.regX)
			cpu.regPC += 2
			data = sys.read(addr)
			b := cpu.regP&cpuRegC != 0
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x01 != 0 {
				cpu.regP |= cpuRegC
			}
			data >>= 1
			if b {
				data |= 0x80
			}
			cpu.regP |= cpu.znTable[data]
			sys.write(addr, data)
			nCycleExec += 7
		case 0xa9:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			cpu.regA = data
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 2
		case 0xa5:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			cpu.regA = cpu.ram[data]
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 3
		case 0xb5:
			data = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			cpu.regA = cpu.ram[data]
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 4
		case 0xad:
			addr = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			cpu.regA = sys.read(addr)
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 4
		case 0xbd:
			addr1 = cpu.readWord(cpu.regPC)
			addr = addr1 + uint16(cpu.regX)
			cpu.regPC += 2
			cpu.regA = sys.read(addr)
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			if (addr1 & 0xff00) != (addr & 0xff00) {
				nCycleExec++
			}
			nCycleExec += 4
		case 0xb9:
			addr1 = cpu.readWord(cpu.regPC)
			addr = addr1 + uint16(cpu.regY)
			cpu.regPC += 2
			cpu.regA = sys.read(addr)
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			if (addr1 & 0xff00) != (addr & 0xff00) {
				nCycleExec++
			}
			nCycleExec += 4
		case 0xa1:
			data = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			addr = (uint16(cpu.ram[data+1]) << 8) | uint16(cpu.ram[data])
			cpu.regA = sys.read(addr)
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 6
		case 0xb1:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			addr1 = (uint16(cpu.ram[data+1]) << 8) | uint16(cpu.ram[data])
			addr = addr1 + uint16(cpu.regY)
			cpu.regA = sys.read(addr)
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			if (addr1 & 0xff00) != (addr & 0xff00) {
				nCycleExec++
			}
			nCycleExec += 5
		case 0xa2:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			cpu.regX = data
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regX]
			nCycleExec += 2
		case 0xa6:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			cpu.regX = cpu.ram[data]
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regX]
			nCycleExec += 3
		case 0xb6:
			data = sys.read(cpu.regPC) + cpu.regY
			cpu.regPC++
			cpu.regX = cpu.ram[data]
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regX]
			nCycleExec += 4
		case 0xae:
			addr = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			cpu.regX = sys.read(addr)
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regX]
			nCycleExec += 4
		case 0xbe:
			addr1 = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			addr = addr1 + uint16(cpu.regY)
			cpu.regX = sys.read(addr)
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regX]
			if (addr1 & 0xff00) != (addr & 0xff00) {
				nCycleExec++
			}
			nCycleExec += 4
		case 0xa0:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			cpu.regY = data
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regY]
			nCycleExec += 2
		case 0xa4:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			cpu.regY = cpu.ram[data]
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regY]
			nCycleExec += 3
		case 0xb4:
			data = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			cpu.regY = cpu.ram[data]
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regY]
			nCycleExec += 4
		case 0xac:
			addr = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			cpu.regY = sys.read(addr)
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regY]
			nCycleExec += 4
		case 0xbc:
			addr1 = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			addr = addr1 + uint16(cpu.regX)
			cpu.regY = sys.read(addr)
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regY]
			if (addr1 & 0xff00) != (addr & 0xff00) {
				nCycleExec++
			}
			nCycleExec += 4
		case 0x85:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			cpu.ram[data] = cpu.regA
			nCycleExec += 3
		case 0x95:
			data = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			cpu.ram[data] = cpu.regA
			nCycleExec += 4
		case 0x8d:
			addr = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			sys.write(addr, cpu.regA)
			nCycleExec += 4
		case 0x9d:
			addr = cpu.readWord(cpu.regPC) + uint16(cpu.regX)
			cpu.regPC += 2
			sys.write(addr, cpu.regA)
			nCycleExec += 5
		case 0x99:
			addr = cpu.readWord(cpu.regPC) + uint16(cpu.regY)
			cpu.regPC += 2
			sys.write(addr, cpu.regA)
			nCycleExec += 5
		case 0x81:
			data = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			addr = (uint16(cpu.ram[data+1]) << 8) | uint16(cpu.ram[data])
			sys.write(addr, cpu.regA)
			nCycleExec += 6
		case 0x91:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			addr = (uint16(cpu.ram[data+1]) << 8) | uint16(cpu.ram[data]) + uint16(cpu.regY)
			sys.write(addr, cpu.regA)
			nCycleExec += 6
		case 0x86:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			cpu.ram[data] = cpu.regX
			nCycleExec += 3
		case 0x96:
			data = sys.read(cpu.regPC) + cpu.regY
			cpu.regPC++
			cpu.ram[data] = cpu.regX
			nCycleExec += 4
		case 0x8e:
			addr = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			sys.write(addr, cpu.regX)
			nCycleExec += 4
		case 0x84:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			cpu.ram[data] = cpu.regY
			nCycleExec += 3
		case 0x94:
			data = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			cpu.ram[data] = cpu.regY
			nCycleExec += 4
		case 0x8c:
			addr = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			sys.write(addr, cpu.regY)
			nCycleExec += 4
		case 0xaa:
			cpu.regX = cpu.regA
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regX]
			nCycleExec += 2
		case 0x8a:
			cpu.regA = cpu.regX
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 2
		case 0xa8:
			cpu.regY = cpu.regA
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regY]
			nCycleExec += 2
		case 0x98:
			cpu.regA = cpu.regY
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 2
		case 0xba:
			cpu.regX = cpu.regS
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regX]
			nCycleExec += 2
		case 0x9a:
			cpu.regS = cpu.regX
			nCycleExec += 2
		case 0xc9:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			word = uint16(cpu.regA) - uint16(data)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if word&0x8000 == 0 {
				cpu.regP |= cpuRegC
			}
			cpu.regP |= cpu.znTable[byte(word)]
			nCycleExec += 2
		case 0xc5:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			data = cpu.ram[data]
			word = uint16(cpu.regA) - uint16(data)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if word&0x8000 == 0 {
				cpu.regP |= cpuRegC
			}
			cpu.regP |= cpu.znTable[byte(word)]
			nCycleExec += 3
		case 0xd5:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			data = cpu.ram[data+cpu.regX]
			word = uint16(cpu.regA) - uint16(data)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if word&0x8000 == 0 {
				cpu.regP |= cpuRegC
			}
			cpu.regP |= cpu.znTable[byte(word)]
			nCycleExec += 4
		case 0xcd:
			addr = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			data = sys.read(addr)
			word = uint16(cpu.regA) - uint16(data)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if word&0x8000 == 0 {
				cpu.regP |= cpuRegC
			}
			cpu.regP |= cpu.znTable[byte(word)]
			nCycleExec += 4
		case 0xdd:
			addr1 = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			addr = addr1 + uint16(cpu.regX)
			data = sys.read(addr)
			word = uint16(cpu.regA) - uint16(data)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if word&0x8000 == 0 {
				cpu.regP |= cpuRegC
			}
			cpu.regP |= cpu.znTable[byte(word)]
			if (addr1 & 0xff00) != (addr & 0xff00) {
				nCycleExec++
			}
			nCycleExec += 4
		case 0xd9:
			addr1 = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			addr = addr1 + uint16(cpu.regY)
			data = sys.read(addr)
			word = uint16(cpu.regA) - uint16(data)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if word&0x8000 == 0 {
				cpu.regP |= cpuRegC
			}
			cpu.regP |= cpu.znTable[byte(word)]
			if (addr1 & 0xff00) != (addr & 0xff00) {
				nCycleExec++
			}
			nCycleExec += 4
		case 0xc1:
			data = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			addr = (uint16(cpu.ram[data+1]) << 8) | uint16(cpu.ram[data])
			data = sys.read(addr)
			word = uint16(cpu.regA) - uint16(data)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if word&0x8000 == 0 {
				cpu.regP |= cpuRegC
			}
			cpu.regP |= cpu.znTable[byte(word)]
			nCycleExec += 6
		case 0xd1:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			addr1 = (uint16(cpu.ram[data+1]) << 8) | uint16(cpu.ram[data])
			addr = addr1 + uint16(cpu.regY)
			data = sys.read(addr)
			word = uint16(cpu.regA) - uint16(data)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if word&0x8000 == 0 {
				cpu.regP |= cpuRegC
			}
			cpu.regP |= cpu.znTable[byte(word)]
			if (addr1 & 0xff00) != (addr & 0xff00) {
				nCycleExec++
			}
			nCycleExec += 5
		case 0xe0:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			word = uint16(cpu.regX) - uint16(data)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if word&0x8000 == 0 {
				cpu.regP |= cpuRegC
			}
			cpu.regP |= cpu.znTable[byte(word)]
			nCycleExec += 2
		case 0xe4:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			data = cpu.ram[data]
			word = uint16(cpu.regX) - uint16(data)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if word&0x8000 == 0 {
				cpu.regP |= cpuRegC
			}
			cpu.regP |= cpu.znTable[byte(word)]
			nCycleExec += 3
		case 0xec:
			addr = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			data = sys.read(addr)
			word = uint16(cpu.regX) - uint16(data)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if word&0x8000 == 0 {
				cpu.regP |= cpuRegC
			}
			cpu.regP |= cpu.znTable[byte(word)]
			nCycleExec += 4
		case 0xc0:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			word = uint16(cpu.regY) - uint16(data)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if word&0x8000 == 0 {
				cpu.regP |= cpuRegC
			}
			cpu.regP |= cpu.znTable[byte(word)]
			nCycleExec += 2
		case 0xc4:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			data = cpu.ram[data]
			word = uint16(cpu.regY) - uint16(data)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if word&0x8000 == 0 {
				cpu.regP |= cpuRegC
			}
			cpu.regP |= cpu.znTable[byte(word)]
			nCycleExec += 3
		case 0xcc:
			addr = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			data = sys.read(addr)
			word = uint16(cpu.regY) - uint16(data)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if word&0x8000 == 0 {
				cpu.regP |= cpuRegC
			}
			cpu.regP |= cpu.znTable[byte(word)]
			nCycleExec += 4
		case 0x90:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			if cpu.regP&cpuRegC == 0 {
				addr1, addr = cpu.regPC, cpu.regPC+uint16(int8(data))
				cpu.regPC = addr
				nCycleExec++
				if (addr1 & 0xff00) != (addr & 0xff00) {
					nCycleExec++
				}
			}
			nCycleExec += 2
		case 0xb0:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			if cpu.regP&cpuRegC != 0 {
				addr1, addr = cpu.regPC, cpu.regPC+uint16(int8(data))
				cpu.regPC = addr
				nCycleExec++
				if (addr1 & 0xff00) != (addr & 0xff00) {
					nCycleExec++
				}
			}
			nCycleExec += 2
		case 0xf0:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			if cpu.regP&cpuRegZ != 0 {
				addr1, addr = cpu.regPC, cpu.regPC+uint16(int8(data))
				cpu.regPC = addr
				nCycleExec++
				if (addr1 & 0xff00) != (addr & 0xff00) {
					nCycleExec++
				}
			}
			nCycleExec += 2
		case 0x30:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			if cpu.regP&cpuRegN != 0 {
				addr1, addr = cpu.regPC, cpu.regPC+uint16(int8(data))
				cpu.regPC = addr
				nCycleExec++
				if (addr1 & 0xff00) != (addr & 0xff00) {
					nCycleExec++
				}
			}
			nCycleExec += 2
		case 0xd0:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			if cpu.regP&cpuRegZ == 0 {
				addr1, addr = cpu.regPC, cpu.regPC+uint16(int8(data))
				cpu.regPC = addr
				nCycleExec++
				if (addr1 & 0xff00) != (addr & 0xff00) {
					nCycleExec++
				}
			}
			nCycleExec += 2
		case 0x10:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			if cpu.regP&cpuRegN == 0 {
				addr1, addr = cpu.regPC, cpu.regPC+uint16(int8(data))
				cpu.regPC = addr
				nCycleExec++
				if (addr1 & 0xff00) != (addr & 0xff00) {
					nCycleExec++
				}
			}
			nCycleExec += 2
		case 0x50:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			if cpu.regP&cpuRegV == 0 {
				addr1, addr = cpu.regPC, cpu.regPC+uint16(int8(data))
				cpu.regPC = addr
				nCycleExec++
				if (addr1 & 0xff00) != (addr & 0xff00) {
					nCycleExec++
				}
			}
			nCycleExec += 2
		case 0x70:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			if cpu.regP&cpuRegV != 0 {
				addr1, addr = cpu.regPC, cpu.regPC+uint16(int8(data))
				cpu.regPC = addr
				nCycleExec++
				if (addr1 & 0xff00) != (addr & 0xff00) {
					nCycleExec++
				}
			}
			nCycleExec += 2
		case 0x4c:
			cpu.regPC = cpu.readWord(cpu.regPC)
			nCycleExec += 3
		case 0x6c:
			word = cpu.readWord(cpu.regPC)
			addr = uint16(sys.read(word))
			word = (word & 0xff00) | ((word + 1) & 0x00ff)
			cpu.regPC = addr + (uint16(sys.read(word)) << 8)
			nCycleExec += 5
		case 0x20:
			addr = cpu.readWord(cpu.regPC)
			cpu.regPC++
			cpu.ram[0x0100+uint16(cpu.regS)] = byte(cpu.regPC >> 8)
			cpu.ram[0x0100+uint16(cpu.regS)-1] = byte(cpu.regPC & 0xff)
			cpu.regS -= 2
			cpu.regPC = addr
			nCycleExec += 6
		case 0x40:
			addr = 0x0100 + uint16(cpu.regS)
			cpu.regP = cpu.ram[addr+1] | cpuRegR
			cpu.regPC = (uint16(cpu.ram[addr+3]) << 8) | uint16(cpu.ram[addr+2])
			cpu.regS += 3
			nCycleExec += 6
		case 0x60:
			addr = 0x0100 + uint16(cpu.regS)
			cpu.regPC = (uint16(cpu.ram[addr+2]) << 8) | uint16(cpu.ram[addr+1]) + 1
			cpu.regS += 2
			nCycleExec += 6
		case 0x18:
			cpu.regP &^= cpuRegC
			nCycleExec += 2
		case 0xd8:
			cpu.regP &^= cpuRegD
			nCycleExec += 2
		case 0x58:
			cpu.regP &^= cpuRegI
			nCycleExec += 2
		case 0xb8:
			cpu.regP &^= cpuRegV
			nCycleExec += 2
		case 0x38:
			cpu.regP |= cpuRegC
			nCycleExec += 2
		case 0xf8:
			cpu.regP |= cpuRegD
			nCycleExec += 2
		case 0x78:
			cpu.regP |= cpuRegI
			nCycleExec += 2
		case 0x48:
			cpu.ram[0x0100|uint16(cpu.regS)] = cpu.regA
			cpu.regS--
			nCycleExec += 3
		case 0x08:
			cpu.ram[0x0100|uint16(cpu.regS)] = cpu.regP | cpuRegB
			cpu.regS--
			nCycleExec += 3
		case 0x68:
			cpu.regS++
			cpu.regA = cpu.ram[0x0100|uint16(cpu.regS)]
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 4
		case 0x28:
			cpu.regS++
			cpu.regP = cpu.ram[0x0100|uint16(cpu.regS)] | cpuRegR
			nCycleExec += 4
		case 0x00:
			cpu.regPC++
			cpu.ram[0x0100+uint16(cpu.regS)] = uint8(cpu.regPC >> 8)
			cpu.ram[0x0100+uint16(cpu.regS)-1] = uint8(cpu.regPC & 0xff)
			cpu.regP |= cpuRegB
			cpu.ram[0x0100+uint16(cpu.regS)-2] = cpu.regP
			cpu.regS -= 3
			cpu.regP |= cpuRegI
			cpu.regPC = cpu.readWord(0xfffe)
			nCycleExec += 7
		case 0x0b, 0x2b:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			cpu.regA &= data
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			if cpu.regP&cpuRegN != 0 {
				cpu.regP |= cpuRegC
			}
			nCycleExec += 2
		case 0x8b:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			cpu.regA = (cpu.regA | 0xee) & cpu.regX & data
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 2
		case 0x6b:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			data &= cpu.regA
			cpu.regA = (data >> 1) | ((cpu.regP & cpuRegC) << 7)
			cpu.regP &^= cpuRegC | cpuRegV | cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			if cpu.regA&0x40 != 0 {
				cpu.regP |= cpuRegC
			}
			if (cpu.regA>>6)^(cpu.regA>>5) != 0 {
				cpu.regP |= cpuRegV
			}
			nCycleExec += 2
		case 0x4b:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			data &= cpu.regA
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x01 != 0 {
				cpu.regP |= cpuRegC
			}
			cpu.regA = data >> 1
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 2
		case 0xc7:
			data1 = sys.read(cpu.regPC)
			cpu.regPC++
			data = cpu.ram[data1] - 1
			word = uint16(cpu.regA) - uint16(data)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if (word & 0x8000) == 0 {
				cpu.regP |= cpuRegC
			}
			cpu.regP |= cpu.znTable[byte(word)]
			cpu.ram[data1] = data
			nCycleExec += 5
		case 0xd7:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			data1 = data + cpu.regX
			data = cpu.ram[data1] - 1
			word = uint16(cpu.regA) - uint16(data)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if (word & 0x8000) == 0 {
				cpu.regP |= cpuRegC
			}
			cpu.regP |= cpu.znTable[byte(word)]
			cpu.ram[data1] = data
			nCycleExec += 6
		case 0xcf:
			addr = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			data = sys.read(addr) - 1
			word = uint16(cpu.regA) - uint16(data)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if (word & 0x8000) == 0 {
				cpu.regP |= cpuRegC
			}
			cpu.regP |= cpu.znTable[byte(word)]
			sys.write(addr, data)
			nCycleExec += 6
		case 0xdf:
			addr = cpu.readWord(cpu.regPC) + uint16(cpu.regX)
			cpu.regPC += 2
			data = sys.read(addr) - 1
			word = uint16(cpu.regA) - uint16(data)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if (word & 0x8000) == 0 {
				cpu.regP |= cpuRegC
			}
			cpu.regP |= cpu.znTable[byte(word)]
			sys.write(addr, data)
			nCycleExec += 7
		case 0xdb:
			addr = cpu.readWord(cpu.regPC) + uint16(cpu.regY)
			cpu.regPC += 2
			data = sys.read(addr) - 1
			word = uint16(cpu.regA) - uint16(data)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if (word & 0x8000) == 0 {
				cpu.regP |= cpuRegC
			}
			cpu.regP |= cpu.znTable[byte(word)]
			sys.write(addr, data)
			nCycleExec += 7
		case 0xc3:
			data = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			addr = (uint16(cpu.ram[data+1]) << 8) | uint16(cpu.ram[data])
			data = sys.read(addr) - 1
			word = uint16(cpu.regA) - uint16(data)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if (word & 0x8000) == 0 {
				cpu.regP |= cpuRegC
			}
			cpu.regP |= cpu.znTable[byte(word)]
			sys.write(addr, data)
			nCycleExec += 8
		case 0xd3:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			addr = (uint16(cpu.ram[data+1]) << 8) | uint16(cpu.ram[data]) + uint16(cpu.regY)
			data = sys.read(addr) - 1
			word = uint16(cpu.regA) - uint16(data)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if (word & 0x8000) == 0 {
				cpu.regP |= cpuRegC
			}
			cpu.regP |= cpu.znTable[byte(word)]
			sys.write(addr, data)
			nCycleExec += 8
		case 0xe7:
			data1 = sys.read(cpu.regPC)
			cpu.regPC++
			data = cpu.ram[data1] + 1
			word = uint16(cpu.regA) - uint16(data) - uint16(^cpu.regP&cpuRegC)
			cpu.regP &^= cpuRegV | cpuRegC | cpuRegZ | cpuRegN
			if (cpu.regA^data)&(cpu.regA^byte(word))&0x80 != 0 {
				cpu.regP |= cpuRegV
			}
			if word < 0x0100 {
				cpu.regP |= cpuRegC
			}
			cpu.regA = byte(word)
			cpu.regP |= cpu.znTable[cpu.regA]
			cpu.ram[data1] = data
			nCycleExec += 5
		case 0xf7:
			data1 = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			data = cpu.ram[data1] + 1
			word = uint16(cpu.regA) - uint16(data) - uint16(^cpu.regP&cpuRegC)
			cpu.regP &^= cpuRegV | cpuRegC | cpuRegZ | cpuRegN
			if (cpu.regA^data)&(cpu.regA^byte(word))&0x80 != 0 {
				cpu.regP |= cpuRegV
			}
			if word < 0x0100 {
				cpu.regP |= cpuRegC
			}
			cpu.regA = byte(word)
			cpu.regP |= cpu.znTable[cpu.regA]
			cpu.ram[data1] = data
			nCycleExec += 5
		case 0xef:
			addr = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			data = sys.read(addr) + 1
			word = uint16(cpu.regA) - uint16(data) - uint16(^cpu.regP&cpuRegC)
			cpu.regP &^= cpuRegV | cpuRegC | cpuRegZ | cpuRegN
			if (cpu.regA^data)&(cpu.regA^byte(word))&0x80 != 0 {
				cpu.regP |= cpuRegV
			}
			if word < 0x0100 {
				cpu.regP |= cpuRegC
			}
			cpu.regA = byte(word)
			cpu.regP |= cpu.znTable[cpu.regA]
			sys.write(addr, data)
			nCycleExec += 5
		case 0xff:
			addr = cpu.readWord(cpu.regPC) + uint16(cpu.regX)
			cpu.regPC += 2
			data = sys.read(addr) + 1
			word = uint16(cpu.regA) - uint16(data) - uint16(^cpu.regP&cpuRegC)
			cpu.regP &^= cpuRegV | cpuRegC | cpuRegZ | cpuRegN
			if (cpu.regA^data)&(cpu.regA^byte(word))&0x80 != 0 {
				cpu.regP |= cpuRegV
			}
			if word < 0x0100 {
				cpu.regP |= cpuRegC
			}
			cpu.regA = byte(word)
			cpu.regP |= cpu.znTable[cpu.regA]
			sys.write(addr, data)
			nCycleExec += 5
		case 0xfb:
			addr = cpu.readWord(cpu.regPC) + uint16(cpu.regY)
			cpu.regPC += 2
			data = sys.read(addr) + 1
			word = uint16(cpu.regA) - uint16(data) - uint16(^cpu.regP&cpuRegC)
			cpu.regP &^= cpuRegV | cpuRegC | cpuRegZ | cpuRegN
			if (cpu.regA^data)&(cpu.regA^byte(word))&0x80 != 0 {
				cpu.regP |= cpuRegV
			}
			if word < 0x0100 {
				cpu.regP |= cpuRegC
			}
			cpu.regA = byte(word)
			cpu.regP |= cpu.znTable[cpu.regA]
			sys.write(addr, data)
			nCycleExec += 5
		case 0xe3:
			data = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			addr = uint16(cpu.ram[data+1]<<8) | uint16(cpu.ram[data])
			data = sys.read(addr) + 1
			word = uint16(cpu.regA) - uint16(data) - uint16(^cpu.regP&cpuRegC)
			cpu.regP &^= cpuRegV | cpuRegC | cpuRegZ | cpuRegN
			if (cpu.regA^data)&(cpu.regA^byte(word))&0x80 != 0 {
				cpu.regP |= cpuRegV
			}
			if word < 0x0100 {
				cpu.regP |= cpuRegC
			}
			cpu.regA = byte(word)
			cpu.regP |= cpu.znTable[cpu.regA]
			sys.write(addr, data)
			nCycleExec += 5
		case 0xf3:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			addr = uint16(cpu.ram[data+1]<<8) | uint16(cpu.ram[data]) + uint16(cpu.regY)
			data = sys.read(addr) + 1
			word = uint16(cpu.regA) - uint16(data) - uint16(^cpu.regP&cpuRegC)
			cpu.regP &^= cpuRegV | cpuRegC | cpuRegZ | cpuRegN
			if (cpu.regA^data)&(cpu.regA^byte(word))&0x80 != 0 {
				cpu.regP |= cpuRegV
			}
			if word < 0x0100 {
				cpu.regP |= cpuRegC
			}
			cpu.regA = byte(word)
			cpu.regP |= cpu.znTable[cpu.regA]
			sys.write(addr, data)
			nCycleExec += 5
		case 0xbb:
			addr1 = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			addr = addr1 + uint16(cpu.regY)
			data = cpu.regS & sys.read(addr)
			cpu.regA, cpu.regX, cpu.regS = data, data, data
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			if (addr1 & 0xff00) != (addr & 0xff00) {
				nCycleExec++
			}
			nCycleExec += 4
		case 0xa7:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			cpu.regA = cpu.ram[data]
			cpu.regX = cpu.regA
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 3
		case 0xb7:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			cpu.regA = cpu.ram[data+cpu.regY]
			cpu.regX = cpu.regA
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 4
		case 0xaf:
			addr = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			cpu.regA = sys.read(addr)
			cpu.regX = cpu.regA
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 4
		case 0xbf:
			addr1 = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			addr = addr1 + uint16(cpu.regY)
			cpu.regA = sys.read(addr)
			cpu.regX = cpu.regA
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			if (addr1 & 0xff00) != (addr & 0xff00) {
				nCycleExec++
			}
			nCycleExec += 4
		case 0xa3:
			data = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			addr = (uint16(cpu.ram[data+1]) << 8) | uint16(cpu.ram[data])
			cpu.regA = sys.read(addr)
			cpu.regX = cpu.regA
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 6
		case 0xb3:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			addr1 = (uint16(cpu.ram[data+1]) << 8) | uint16(cpu.ram[data])
			addr = addr1 + uint16(cpu.regY)
			cpu.regA = sys.read(addr)
			cpu.regX = cpu.regA
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			if (addr1 & 0xff00) != (addr & 0xff00) {
				nCycleExec++
			}
			nCycleExec += 5
		case 0xab:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			data &= cpu.regA | 0xee
			cpu.regA, cpu.regX = data, data
			cpu.regP &^= cpuRegZ | cpuRegN
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 2
		case 0x27:
			data1 = sys.read(cpu.regPC)
			cpu.regPC++
			data = cpu.ram[data1]
			b := cpu.regP&cpuRegC != 0
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x80 != 0 {
				cpu.regP |= cpuRegC
			}
			data <<= 1
			if b {
				data |= 1
			}
			cpu.regA &= data
			cpu.regP |= cpu.znTable[cpu.regA]
			cpu.ram[data1] = data
			nCycleExec += 5
		case 0x37:
			data1 = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			data = cpu.ram[data1]
			b := cpu.regP&cpuRegC != 0
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x80 != 0 {
				cpu.regP |= cpuRegC
			}
			data <<= 1
			if b {
				data |= 1
			}
			cpu.regA &= data
			cpu.regP |= cpu.znTable[cpu.regA]
			cpu.ram[data1] = data
			nCycleExec += 6
		case 0x2f:
			addr = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			data = sys.read(addr)
			b := cpu.regP&cpuRegC != 0
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x80 != 0 {
				cpu.regP |= cpuRegC
			}
			data <<= 1
			if b {
				data |= 1
			}
			cpu.regA &= data
			cpu.regP |= cpu.znTable[cpu.regA]
			sys.write(addr, data)
			nCycleExec += 6
		case 0x3f:
			addr = cpu.readWord(cpu.regPC) + uint16(cpu.regX)
			cpu.regPC += 2
			data = sys.read(addr)
			b := cpu.regP&cpuRegC != 0
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x80 != 0 {
				cpu.regP |= cpuRegC
			}
			data <<= 1
			if b {
				data |= 1
			}
			cpu.regA &= data
			cpu.regP |= cpu.znTable[cpu.regA]
			sys.write(addr, data)
			nCycleExec += 7
		case 0x3b:
			addr = cpu.readWord(cpu.regPC) + uint16(cpu.regY)
			cpu.regPC += 2
			data = sys.read(addr)
			b := cpu.regP&cpuRegC != 0
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x80 != 0 {
				cpu.regP |= cpuRegC
			}
			data <<= 1
			if b {
				data |= 1
			}
			cpu.regA &= data
			cpu.regP |= cpu.znTable[cpu.regA]
			sys.write(addr, data)
			nCycleExec += 7
		case 0x23:
			data = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			addr = (uint16(cpu.ram[data+1]) << 8) | uint16(cpu.ram[data])
			data = sys.read(addr)
			b := cpu.regP&cpuRegC != 0
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x80 != 0 {
				cpu.regP |= cpuRegC
			}
			data <<= 1
			if b {
				data |= 1
			}
			cpu.regA &= data
			cpu.regP |= cpu.znTable[cpu.regA]
			sys.write(addr, data)
			nCycleExec += 8
		case 0x33:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			addr = (uint16(cpu.ram[data+1]) << 8) | uint16(cpu.ram[data]) + uint16(cpu.regY)
			data = sys.read(addr)
			b := cpu.regP&cpuRegC != 0
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x80 != 0 {
				cpu.regP |= cpuRegC
			}
			data <<= 1
			if b {
				data |= 1
			}
			cpu.regA &= data
			cpu.regP |= cpu.znTable[cpu.regA]
			sys.write(addr, data)
			nCycleExec += 8
		case 0x67:
			data1 = sys.read(cpu.regPC)
			cpu.regPC++
			data = cpu.ram[data1]
			b := cpu.regP&cpuRegC != 0
			cpu.regP &^= cpuRegC | cpuRegV | cpuRegZ | cpuRegN
			if data&0x01 != 0 {
				cpu.regP |= cpuRegC
			}
			data >>= 1
			if b {
				data |= 0x80
			}
			word = uint16(cpu.regA) + uint16(data) + uint16(cpu.regP&cpuRegC)
			cpu.regA &= byte(word)
			if word > 0x00ff {
				cpu.regP |= cpuRegC
			}
			if (^(cpu.regA ^ data))&(cpu.regA^byte(word))&0x80 != 0 {
				cpu.regP |= cpuRegV
			}
			cpu.regP |= cpu.znTable[cpu.regA]
			cpu.ram[data1] = data
			nCycleExec += 5
		case 0x77:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			data1 = data + cpu.regX
			data = cpu.ram[data1]
			b := cpu.regP&cpuRegC != 0
			cpu.regP &^= cpuRegC | cpuRegV | cpuRegZ | cpuRegN
			if data&0x01 != 0 {
				cpu.regP |= cpuRegC
			}
			data >>= 1
			if b {
				data |= 0x80
			}
			word = uint16(cpu.regA) + uint16(data) + uint16(cpu.regP&cpuRegC)
			cpu.regA &= byte(word)
			if word > 0x00ff {
				cpu.regP |= cpuRegC
			}
			if (^(cpu.regA ^ data))&(cpu.regA^byte(word))&0x80 != 0 {
				cpu.regP |= cpuRegV
			}
			cpu.regP |= cpu.znTable[cpu.regA]
			cpu.ram[data1] = data
			nCycleExec += 6
		case 0x6f:
			addr = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			data = sys.read(addr)
			b := cpu.regP&cpuRegC != 0
			cpu.regP &^= cpuRegC | cpuRegV | cpuRegZ | cpuRegN
			if data&0x01 != 0 {
				cpu.regP |= cpuRegC
			}
			data >>= 1
			if b {
				data |= 0x80
			}
			word = uint16(cpu.regA) + uint16(data) + uint16(cpu.regP&cpuRegC)
			cpu.regA &= byte(word)
			if word > 0x00ff {
				cpu.regP |= cpuRegC
			}
			if (^(cpu.regA ^ data))&(cpu.regA^byte(word))&0x80 != 0 {
				cpu.regP |= cpuRegV
			}
			cpu.regP |= cpu.znTable[cpu.regA]
			sys.write(addr, data)
			nCycleExec += 6
		case 0x7f:
			addr = cpu.readWord(cpu.regPC) + uint16(cpu.regX)
			cpu.regPC += 2
			data = sys.read(addr)
			b := cpu.regP&cpuRegC != 0
			cpu.regP &^= cpuRegC | cpuRegV | cpuRegZ | cpuRegN
			if data&0x01 != 0 {
				cpu.regP |= cpuRegC
			}
			data >>= 1
			if b {
				data |= 0x80
			}
			word = uint16(cpu.regA) + uint16(data) + uint16(cpu.regP&cpuRegC)
			cpu.regA &= byte(word)
			if word > 0x00ff {
				cpu.regP |= cpuRegC
			}
			if (^(cpu.regA ^ data))&(cpu.regA^byte(word))&0x80 != 0 {
				cpu.regP |= cpuRegV
			}
			cpu.regP |= cpu.znTable[cpu.regA]
			sys.write(addr, data)
			nCycleExec += 7
		case 0x7b:
			addr = cpu.readWord(cpu.regPC) + uint16(cpu.regY)
			cpu.regPC += 2
			data = sys.read(addr)
			b := cpu.regP&cpuRegC != 0
			cpu.regP &^= cpuRegC | cpuRegV | cpuRegZ | cpuRegN
			if data&0x01 != 0 {
				cpu.regP |= cpuRegC
			}
			data >>= 1
			if b {
				data |= 0x80
			}
			word = uint16(cpu.regA) + uint16(data) + uint16(cpu.regP&cpuRegC)
			cpu.regA &= byte(word)
			if word > 0x00ff {
				cpu.regP |= cpuRegC
			}
			if (^(cpu.regA ^ data))&(cpu.regA^byte(word))&0x80 != 0 {
				cpu.regP |= cpuRegV
			}
			cpu.regP |= cpu.znTable[cpu.regA]
			sys.write(addr, data)
			nCycleExec += 7
		case 0x63:
			data = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			addr = (uint16(cpu.ram[data+1]) << 8) | uint16(cpu.ram[data])
			data = sys.read(addr)
			b := cpu.regP&cpuRegC != 0
			cpu.regP &^= cpuRegC | cpuRegV | cpuRegZ | cpuRegN
			if data&0x01 != 0 {
				cpu.regP |= cpuRegC
			}
			data >>= 1
			if b {
				data |= 0x80
			}
			word = uint16(cpu.regA) + uint16(data) + uint16(cpu.regP&cpuRegC)
			cpu.regA &= byte(word)
			if word > 0x00ff {
				cpu.regP |= cpuRegC
			}
			if (^(cpu.regA ^ data))&(cpu.regA^byte(word))&0x80 != 0 {
				cpu.regP |= cpuRegV
			}
			cpu.regP |= cpu.znTable[cpu.regA]
			sys.write(addr, data)
			nCycleExec += 8
		case 0x73:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			addr = (uint16(cpu.ram[data+1]) << 8) | uint16(cpu.ram[data]) + uint16(cpu.regY)
			data = sys.read(addr)
			b := cpu.regP&cpuRegC != 0
			cpu.regP &^= cpuRegC | cpuRegV | cpuRegZ | cpuRegN
			if data&0x01 != 0 {
				cpu.regP |= cpuRegC
			}
			data >>= 1
			if b {
				data |= 0x80
			}
			word = uint16(cpu.regA) + uint16(data) + uint16(cpu.regP&cpuRegC)
			cpu.regA &= byte(word)
			if word > 0x00ff {
				cpu.regP |= cpuRegC
			}
			if (^(cpu.regA ^ data))&(cpu.regA^byte(word))&0x80 != 0 {
				cpu.regP |= cpuRegV
			}
			cpu.regP |= cpu.znTable[cpu.regA]
			sys.write(addr, data)
			nCycleExec += 8
		case 0x87:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			cpu.ram[data] = cpu.regA & cpu.regX
			nCycleExec += 3
		case 0x97:
			data = sys.read(cpu.regPC) + cpu.regY
			cpu.regPC++
			cpu.ram[data] = cpu.regA & cpu.regX
			nCycleExec += 4
		case 0x8f:
			addr = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			sys.write(addr, cpu.regA&cpu.regX)
			nCycleExec += 4
		case 0x83:
			data = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			addr = (uint16(cpu.ram[data+1]) << 8) | uint16(cpu.ram[data])
			sys.write(addr, cpu.regA&cpu.regX)
			nCycleExec += 6
		case 0xcb:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			word = uint16(cpu.regA&cpu.regX) - uint16(data)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if word < 0x0100 {
				cpu.regP |= cpuRegC
			}
			cpu.regX = byte(word)
			cpu.regP |= cpu.znTable[cpu.regX]
			nCycleExec += 2
		case 0x9f:
			addr = cpu.readWord(cpu.regPC) + uint16(cpu.regY)
			cpu.regPC += 2
			data = cpu.regA & cpu.regX & byte((addr>>8)+1)
			sys.write(addr, data)
			nCycleExec += 5
		case 0x93:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			addr = (uint16(cpu.ram[data+1]) << 8) | uint16(cpu.ram[data]) + uint16(cpu.regY)
			data = cpu.regA & cpu.regX & byte((addr>>8)+1)
			sys.write(addr, data)
			nCycleExec += 6
		case 0x9B:
			addr = cpu.readWord(cpu.regPC) + uint16(cpu.regY)
			cpu.regPC += 2
			cpu.regS = cpu.regA & cpu.regX
			data = cpu.regS & byte((addr>>8)+1)
			sys.write(addr, data)
			nCycleExec += 5
		case 0x9e:
			addr = cpu.readWord(cpu.regPC) + uint16(cpu.regY)
			cpu.regPC += 2
			data = cpu.regX & byte((addr>>8)+1)
			sys.write(addr, data)
			nCycleExec += 5
		case 0x9c:
			addr = cpu.readWord(cpu.regPC) + uint16(cpu.regX)
			cpu.regPC += 2
			data = cpu.regY & byte((addr>>8)+1)
			sys.write(addr, data)
			nCycleExec += 5
		case 0x07:
			data1 = sys.read(cpu.regPC)
			cpu.regPC++
			data = cpu.ram[data1]
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x80 != 0 {
				cpu.regP |= cpuRegC
			}
			data <<= 1
			cpu.regA |= data
			cpu.regP |= cpu.znTable[cpu.regA]
			cpu.ram[data1] = data
			nCycleExec += 5
		case 0x17:
			data1 = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			data = cpu.ram[data1]
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x80 != 0 {
				cpu.regP |= cpuRegC
			}
			data <<= 1
			cpu.regA |= data
			cpu.regP |= cpu.znTable[cpu.regA]
			cpu.ram[data1] = data
			nCycleExec += 6
		case 0x0f:
			addr = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			data = sys.read(addr)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x80 != 0 {
				cpu.regP |= cpuRegC
			}
			data <<= 1
			cpu.regA |= data
			cpu.regP |= cpu.znTable[cpu.regA]
			sys.write(addr, data)
			nCycleExec += 6
		case 0x1f:
			addr = cpu.readWord(cpu.regPC) + uint16(cpu.regX)
			cpu.regPC += 2
			data = sys.read(addr)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x80 != 0 {
				cpu.regP |= cpuRegC
			}
			data <<= 1
			cpu.regA |= data
			cpu.regP |= cpu.znTable[cpu.regA]
			sys.write(addr, data)
			nCycleExec += 7
		case 0x1b:
			addr = cpu.readWord(cpu.regPC) + uint16(cpu.regY)
			cpu.regPC += 2
			data = sys.read(addr)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x80 != 0 {
				cpu.regP |= cpuRegC
			}
			data <<= 1
			cpu.regA |= data
			cpu.regP |= cpu.znTable[cpu.regA]
			sys.write(addr, data)
			nCycleExec += 7
		case 0x03:
			data = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			addr = (uint16(cpu.ram[data+1]) << 8) | uint16(cpu.ram[data])
			data = sys.read(addr)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x80 != 0 {
				cpu.regP |= cpuRegC
			}
			data <<= 1
			cpu.regA |= data
			cpu.regP |= cpu.znTable[cpu.regA]
			sys.write(addr, data)
			nCycleExec += 8
		case 0x13:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			addr = (uint16(cpu.ram[data+1]) << 8) | uint16(cpu.ram[data]) + uint16(cpu.regY)
			data = sys.read(addr)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x80 != 0 {
				cpu.regP |= cpuRegC
			}
			data <<= 1
			cpu.regA |= data
			cpu.regP |= cpu.znTable[cpu.regA]
			sys.write(addr, data)
			nCycleExec += 8
		case 0x47:
			data1 = sys.read(cpu.regPC)
			cpu.regPC++
			data = cpu.ram[data1]
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x01 != 0 {
				cpu.regP |= cpuRegC
			}
			data >>= 1
			cpu.regA ^= data
			cpu.regP |= cpu.znTable[cpu.regA]
			cpu.ram[data1] = data
			nCycleExec += 5
		case 0x57:
			data1 = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			data = cpu.ram[data1]
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x01 != 0 {
				cpu.regP |= cpuRegC
			}
			data >>= 1
			cpu.regA ^= data
			cpu.regP |= cpu.znTable[cpu.regA]
			cpu.ram[data1] = data
			nCycleExec += 6
		case 0x4f:
			addr = cpu.readWord(cpu.regPC)
			cpu.regPC += 2
			data = sys.read(addr)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x01 != 0 {
				cpu.regP |= cpuRegC
			}
			data >>= 1
			cpu.regA ^= data
			cpu.regP |= cpu.znTable[cpu.regA]
			sys.write(addr, data)
			nCycleExec += 6
		case 0x5f:
			addr = cpu.readWord(cpu.regPC) + uint16(cpu.regX)
			cpu.regPC += 2
			data = sys.read(addr)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x01 != 0 {
				cpu.regP |= cpuRegC
			}
			data >>= 1
			cpu.regA ^= data
			cpu.regP |= cpu.znTable[cpu.regA]
			sys.write(addr, data)
			nCycleExec += 7
		case 0x5b:
			addr = cpu.readWord(cpu.regPC) + uint16(cpu.regY)
			cpu.regPC += 2
			data = sys.read(addr)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x01 != 0 {
				cpu.regP |= cpuRegC
			}
			data >>= 1
			cpu.regA ^= data
			cpu.regP |= cpu.znTable[cpu.regA]
			sys.write(addr, data)
			nCycleExec += 7
		case 0x43:
			data = sys.read(cpu.regPC) + cpu.regX
			cpu.regPC++
			addr = (uint16(cpu.ram[data+1]) << 8) | uint16(cpu.ram[data])
			data = sys.read(addr)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x01 != 0 {
				cpu.regP |= cpuRegC
			}
			data >>= 1
			cpu.regA ^= data
			cpu.regP |= cpu.znTable[cpu.regA]
			sys.write(addr, data)
			nCycleExec += 8
		case 0x53:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			addr = (uint16(cpu.ram[data+1]) << 8) | uint16(cpu.ram[data]) + uint16(cpu.regY)
			data = sys.read(addr)
			cpu.regP &^= cpuRegC | cpuRegZ | cpuRegN
			if data&0x01 != 0 {
				cpu.regP |= cpuRegC
			}
			data >>= 1
			cpu.regA ^= data
			cpu.regP |= cpu.znTable[cpu.regA]
			sys.write(addr, data)
			nCycleExec += 8
		case 0xeb:
			data = sys.read(cpu.regPC)
			cpu.regPC++
			word = uint16(cpu.regA) - uint16(data) - uint16(^cpu.regP&cpuRegC)
			cpu.regP &^= cpuRegV | cpuRegC | cpuRegZ | cpuRegN
			if (cpu.regA^data)&(cpu.regA^byte(word))&0x80 != 0 {
				cpu.regP |= cpuRegV
			}
			if word < 0x0100 {
				cpu.regP |= cpuRegC
			}
			cpu.regA = byte(word)
			cpu.regP |= cpu.znTable[cpu.regA]
			nCycleExec += 2
		case 0x1a, 0x3a, 0x5a, 0x7a, 0xda, 0xea, 0xfa:
			nCycleExec += 2
		case 0x80, 0x82, 0x89, 0xc2, 0xe2:
			cpu.regPC++
			nCycleExec += 2
		case 0x04, 0x44, 0x64:
			cpu.regPC++
			nCycleExec += 3
		case 0x14, 0x34, 0x54, 0x74, 0xd4, 0xf4:
			cpu.regPC++
			nCycleExec += 4
		case 0x0c, 0x1c, 0x3c, 0x5c, 0x7c, 0xdc, 0xfc:
			cpu.regPC += 2
			nCycleExec += 4
		case 0x02, 0x12, 0x22, 0x32, 0x42, 0x52, 0x62, 0x72, 0x92, 0xB2, 0xD2, 0xF2:
			panic("cpu jammed") //ldeng7
		}

		if intrNmi || intrIrq {
			addr := 0x0100 | uint16(cpu.regS)
			cpu.ram[addr] = byte(cpu.regPC >> 8)
			cpu.ram[addr-1] = byte(cpu.regPC & 0xff)
			cpu.regP &^= cpuRegB
			cpu.ram[addr-2] = cpu.regP
			cpu.regS -= 3
			cpu.regP |= cpuRegI
			addr = 0xfffa
			if intrIrq {
				addr = 0xfffe
			}
			cpu.regPC = cpu.readWord(addr)
			nCycleExec += 7
		}
		nCycleReq -= nCycleExec
		cpu.nCycle += nCycleExec
		/*
			if cpu.nCycle < 700001 {
				ppu := cpu.sys.ppu
				cpu.sys.logger.WriteString(fmt.Sprintf("%d %d %d %d %d %d %d loo:%d %d %d %d %d %d ppureg:%d %d %d %d %d\r\n",
					opcode, cpu.regX, cpu.regY, cpu.regA, cpu.regP,
					cpu.regS, cpu.regPC, ppu.loopyT, ppu.loopyV, ppu.loopyX,
					ppu.loopyY, ppu.loopySh, ppu.readBuf, ppu.reg0, ppu.reg1, ppu.reg2, ppu.reg3, cpu.nCycle))
			}
		*/
		cpu.sys.mapper.clock(nCycleExec)
	}
	return cpu.nCycle - nCyclePrev
}
