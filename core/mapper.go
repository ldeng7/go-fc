package core

import "fmt"

type Mapper interface {
	reset()
	readEx(addr uint16) byte
	writeEx(addr uint16, data byte)
	readLow(addr uint16) byte
	writeLow(addr uint16, data byte)
	read(addr uint16) byte
	write(addr uint16, data byte)

	hSync(scanline uint16)
	vSync()
	clock(nCycle int64)
	ppuLatch(addr uint16)
	ppuChrLatch(addr uint16)
	ppuExtLatchX(x byte)
	ppuExtLatchSpOfs() byte
	ppuExtLatch(iNameTbl uint16, chL *byte, chH *byte, attr *byte)
}

func newMapperNil(bm *baseMapper) Mapper { return nil }

var mapperTable = [256]func(bm *baseMapper) Mapper{
	// 0x
	newMapper000, newMapper001, newMapper002, newMapper003, newMapper004, newMapper005, newMapper006, newMapper007,
	newMapper008, newMapper009, newMapper010, newMapper011, newMapper012, newMapper013, newMapperNil, newMapper015,
	// 1x
	newMapper016, newMapper017, newMapper018, newMapper019, newMapperNil, newMapper021, newMapper022, newMapper023,
	newMapper024, newMapper025, newMapper026, newMapper027, newMapperNil, newMapperNil, newMapperNil, newMapperNil,
	// 2x
	newMapper032, newMapper033, newMapper034, newMapperNil, newMapperNil, newMapperNil, newMapperNil, newMapperNil,
	newMapper040, newMapper041, newMapper042, newMapper043, newMapper044, newMapper045, newMapper046, newMapper047,
	// 3x
	newMapperNil, newMapper049, newMapperNil, newMapperNil, newMapper052, newMapperNil, newMapperNil, newMapperNil,
	newMapperNil, newMapper057, newMapper058, newMapperNil, newMapper060, newMapper061, newMapper062, newMapperNil,
	// 4x
	newMapperNil, newMapperNil, newMapper066, newMapperNil, newMapperNil, newMapperNil, newMapper070, newMapper071,
	newMapper072, newMapper073, newMapperNil, newMapperNil, newMapper076, newMapper077, newMapper078, newMapper079,
	// 5x
	newMapperNil, newMapperNil, newMapperNil, newMapperNil, newMapperNil, newMapperNil, newMapper086, newMapper087,
	newMapper088, newMapper089, newMapperNil, newMapperNil, newMapper092, newMapper093, newMapper094, newMapperNil,
	// 6x
	newMapper096, newMapper097, newMapperNil, newMapper099, newMapperNil, newMapper101, newMapperNil, newMapperNil,
	newMapperNil, newMapperNil, newMapperNil, newMapper107, newMapper108, newMapperNil, newMapperNil, newMapper111,
	// 7x
	newMapperNil, newMapperNil, newMapperNil, newMapperNil, newMapperNil, newMapperNil, newMapperNil, newMapperNil,
	newMapper120, newMapper121, newMapper122, newMapperNil, newMapperNil, newMapperNil, newMapperNil, newMapperNil,
	// 8x
	newMapperNil, newMapperNil, newMapperNil, newMapperNil, newMapper132, newMapper133, newMapperNil, newMapperNil,
	newMapperNil, newMapperNil, newMapperNil, newMapperNil, newMapper140, newMapper141, newMapperNil, newMapperNil,
	// 9x
	newMapperNil, newMapperNil, newMapperNil, newMapperNil, newMapper148, newMapperNil, newMapper150, newMapper151,
	newMapperNil, newMapperNil, newMapperNil, newMapperNil, newMapperNil, newMapperNil, newMapperNil, newMapperNil,
	// ax
	newMapperNil, newMapperNil, newMapper162, newMapper163, newMapperNil, newMapperNil, newMapperNil, newMapperNil,
	newMapper168, newMapperNil, newMapperNil, newMapperNil, newMapperNil, newMapperNil, newMapper174, newMapper175,
	// bx
	newMapper176, nil, newMapper178, nil, newMapperb4, newMapperb5, newMapperb6, newMapperb7,
	nil, newMapperb9, newMapperba, newMapperbb, newMapperbc, newMapperbd, newMapperbe, newMapperbf,
	// cx
	newMapperNil, newMapper193, newMapper194, newMapperNil, newMapperNil, newMapperNil, newMapperNil, newMapper199,
	newMapper200, newMapper201, newMapper202, newMapperNil, newMapperNil, newMapperNil, newMapperNil, newMapperNil,
	// dx
	newMapperNil, newMapperNil, newMapperNil, newMapperNil, newMapper212, newMapperNil, newMapperNil, newMapperNil,
	newMapperNil, newMapperNil, newMapperNil, newMapperNil, newMapperNil, newMapperNil, newMapper222, newMapperNil,
	// ex
	newMapperNil, newMapper225, newMapper226, newMapper227, newMapper228, newMapper229, newMapper230, newMapper231,
	newMapper232, newMapper233, newMapper234, newMapper235, newMapper236, newMapper237, newMapperNil, newMapperNil,
	// fx
	newMapperf0, newMapperf1, newMapperf2, newMapperf3, newMapperf4, newMapperf5, newMapperf6, newMapperf7,
	newMapperf8, newMapperf9, nil, newMapperfb, newMapperfc, newMapper253, newMapperfe, newMapperff,
}

type baseMapper struct {
	sys      *Sys
	mem      *Mem
	cpuBanks [][]byte
}

func newMapper(sys *Sys) (Mapper, error) {
	bm := &baseMapper{}
	bm.sys = sys
	bm.mem = sys.mem
	bm.cpuBanks = sys.mem.cpuBanks[:]

	m := mapperTable[sys.rom.mapperNo](bm)
	if nil == m {
		return nil, fmt.Errorf("unsupported mapper #%d", sys.rom.mapperNo)
	}
	return m, nil
}

func (m *baseMapper) reset()                         {}
func (m *baseMapper) readEx(addr uint16) byte        { return 0 }
func (m *baseMapper) writeEx(addr uint16, data byte) {}
func (m *baseMapper) readLow(addr uint16) byte {
	if addr >= 0x6000 {
		return m.cpuBanks[addr>>13][addr&0x1fff]
	}
	return byte(addr >> 8)
}
func (m *baseMapper) writeLow(addr uint16, data byte) {
	if addr >= 0x6000 {
		m.cpuBanks[addr>>13][addr&0x1fff] = data
	}
}
func (m *baseMapper) read(addr uint16) byte {
	return m.cpuBanks[addr>>13][addr&0x1fff]
}
func (m *baseMapper) write(addr uint16, data byte) {}

func (m *baseMapper) hSync(scanline uint16)   {}
func (m *baseMapper) vSync()                  {}
func (m *baseMapper) clock(nCycle int64)      {}
func (m *baseMapper) ppuLatch(addr uint16)    {}
func (m *baseMapper) ppuChrLatch(addr uint16) {}
func (m *baseMapper) ppuExtLatchX(x byte)     {}
func (m *baseMapper) ppuExtLatchSpOfs() byte  { return 0 }
func (m *baseMapper) ppuExtLatch(
	iNameTbl uint16, chL *byte, chH *byte, attr *byte) {
}

func (m *baseMapper) setIntr() {
	m.sys.cpu.intr |= cpuIntrTypMapper
}

func (m *baseMapper) clearIntr() {
	m.sys.cpu.intr &^= cpuIntrTypMapper
}

func (m *baseMapper) isPpuDisp() bool {
	return m.sys.ppu.reg1&(ppuReg1BgDisp|ppuReg1SpDisp) != 0
}
