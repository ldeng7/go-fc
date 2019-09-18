package core

import "fmt"

type Mapper interface {
	reset()
	read(addr uint16, data byte)
	write(addr uint16, data byte)
	readLow(addr uint16) byte
	writeLow(addr uint16, data byte)
	readEx(addr uint16) byte
	writeEx(addr uint16, data byte)
	hSync(scanline uint16)
	vSync()
	clock(nCycle int64)
	ppuLatch(addr uint16)
	ppuChrLatch(addr uint16)
	ppuExtLatchX(x uint8)
	ppuExtLatch(addr uint16, chL *byte, chH *byte, attr *byte)
}

var mapperTable = [256]func(bm *baseMapper) Mapper{
	newMapper00, newMapper01, newMapper02, newMapper03, newMapper04, newMapper05, newMapper06, newMapper07,
	newMapper08, newMapper09, newMapper0a, newMapper0b, newMapper0c, newMapper0d, nil, newMapper0f,
	newMapper10, newMapper11, newMapper12, newMapper13, nil, newMapper15, newMapper16, newMapper17,
	newMapper18, newMapper19, newMapper1a, newMapper1b, nil, nil, nil, nil,
	newMapper20, newMapper21, newMapper22, nil, nil, nil, nil, nil,
	newMapper28, newMapper29, newMapper2a, newMapper2b, newMapper2c, newMapper2d, newMapper2e, newMapper2f,
	newMapper30, newMapper31, newMapper32, newMapper33, newMapper34, newMapper35, newMapper36, newMapper37,
	newMapper38, newMapper39, newMapper3a, newMapper3b, newMapper3c, newMapper3d, newMapper3e, newMapper3f,
	newMapper40, newMapper41, newMapper42, newMapper43, newMapper44, newMapper45, newMapper46, newMapper47,
	newMapper48, newMapper49, newMapper4a, newMapper4b, newMapper4c, newMapper4d, newMapper4e, newMapper4f,
	newMapper50, newMapper51, newMapper52, newMapper53, newMapper54, newMapper55, newMapper56, newMapper57,
	newMapper58, newMapper59, newMapper5a, newMapper5b, newMapper5c, newMapper5d, newMapper5e, newMapper5f,
	newMapper60, newMapper61, newMapper62, newMapper63, newMapper64, newMapper65, newMapper66, newMapper67,
	newMapper68, newMapper69, newMapper6a, newMapper6b, newMapper6c, newMapper6d, newMapper6e, newMapper6f,
	newMapper70, newMapper71, newMapper72, newMapper73, newMapper74, newMapper75, newMapper76, newMapper77,
	newMapper78, newMapper79, newMapper7a, newMapper7b, newMapper7c, newMapper7d, newMapper7e, newMapper7f,
	newMapper80, newMapper81, newMapper82, newMapper83, newMapper84, newMapper85, newMapper86, newMapper87,
	newMapper88, newMapper89, newMapper8a, newMapper8b, newMapper8c, newMapper8d, newMapper8e, newMapper8f,
	newMapper90, newMapper91, newMapper92, newMapper93, newMapper94, newMapper95, newMapper96, newMapper97,
	newMapper98, newMapper99, newMapper9a, newMapper9b, newMapper9c, newMapper9d, newMapper9e, newMapper9f,
	newMappera0, newMappera1, newMappera2, newMappera3, newMappera4, newMappera5, newMappera6, newMappera7,
	newMappera8, newMappera9, newMapperaa, newMapperab, newMapperac, newMapperad, newMapperae, newMapperaf,
	newMapperb0, newMapperb1, newMapperb2, newMapperb3, newMapperb4, newMapperb5, newMapperb6, newMapperb7,
	newMapperb8, newMapperb9, newMapperba, newMapperbb, newMapperbc, newMapperbd, newMapperbe, newMapperbf,
	newMapperc0, newMapperc1, newMapperc2, newMapperc3, newMapperc4, newMapperc5, newMapperc6, newMapperc7,
	newMapperc8, newMapperc9, newMapperca, nil, nil, nil, nil, nil,
	nil, nil, nil, nil, nil, nil, nil, nil,
	nil, nil, nil, nil, nil, nil, newMapperde, newMapperdf,
	newMappere0, newMappere1, newMappere2, newMappere3, newMappere4, newMappere5, newMappere6, newMappere7,
	newMappere8, newMappere9, newMapperea, newMappereb, newMapperec, newMappered, newMapperee, newMapperef,
	newMapperf0, newMapperf1, newMapperf2, newMapperf3, newMapperf4, newMapperf5, newMapperf6, newMapperf7,
	newMapperf8, newMapperf9, newMapperfa, newMapperfb, newMapperfc, newMapperfd, newMapperfe, newMapperff,
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

	f := mapperTable[sys.rom.mapperNo]
	if nil == f {
		return nil, fmt.Errorf("unsupported mapper #%d", sys.rom.mapperNo)
	}
	m := f(bm)
	m.reset()
	return m, nil
}

func (m *baseMapper) reset()                         {}
func (m *baseMapper) read(addr uint16, data byte)    {}
func (m *baseMapper) write(addr uint16, data byte)   {}
func (m *baseMapper) readEx(addr uint16) byte        { return 0 }
func (m *baseMapper) writeEx(addr uint16, data byte) {}
func (m *baseMapper) hSync(scanline uint16)          {}
func (m *baseMapper) vSync()                         {}
func (m *baseMapper) clock(nCycle int64)             {}
func (m *baseMapper) ppuLatch(addr uint16)           {}
func (m *baseMapper) ppuChrLatch(addr uint16)        {}
func (m *baseMapper) ppuExtLatchX(x uint8)           {}

func (m *baseMapper) readLow(addr uint16) byte {
	if addr >= 0x6000 && addr <= 0x7fff {
		return m.cpuBanks[addr>>13][addr&0x1fff]
	}
	return byte(addr)
}

func (m *baseMapper) writeLow(addr uint16, data byte) {
	if addr >= 0x6000 && addr <= 0x7fff {
		m.cpuBanks[addr>>13][addr&0x1fff] = data
	}
}

func (m *baseMapper) ppuExtLatch(addr uint16, chL *byte, chH *byte, attr *byte) {}
