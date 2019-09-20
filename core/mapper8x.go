package core

// 0x85

type mapper85 struct {
	baseMapper
}

func newMapper85(bm *baseMapper) Mapper {
	return &mapper85{baseMapper: *bm}
}

func (m *mapper85) reset() {
	m.mem.setProm32kBank(0)
	m.mem.setVrom8kBank(0)
}

func (m *mapper85) writeLow(addr uint16, data byte) {
	if addr == 0x4120 {
		m.mem.setProm32kBank(uint32(data&0x04) >> 2)
		m.mem.setVrom8kBank(uint32(data & 0x03))
	}
	m.mem.cpuBanks[addr>>13][addr&0x1fff] = data
}

// 0x86

type mapper86 struct {
	baseMapper
}

func newMapper86(bm *baseMapper) Mapper {
	return &mapper86{baseMapper: *bm}
}

func (m *mapper86) reset() {
}

// 0x87

type mapper87 struct {
	baseMapper
}

func newMapper87(bm *baseMapper) Mapper {
	return &mapper87{baseMapper: *bm}
}

func (m *mapper87) reset() {
}

// 0x8c

type mapper8c struct {
	baseMapper
}

func newMapper8c(bm *baseMapper) Mapper {
	return &mapper8c{baseMapper: *bm}
}

func (m *mapper8c) reset() {
	m.mem.setProm32kBank(0)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper8c) writeLow(addr uint16, data byte) {
	m.mem.setProm32kBank(uint32(data&0xf0) >> 4)
	m.mem.setVrom8kBank(uint32(data & 0x0f))
}

// 0x8d

type mapper8d struct {
	baseMapper
}

func newMapper8d(bm *baseMapper) Mapper {
	return &mapper8d{baseMapper: *bm}
}

func (m *mapper8d) reset() {
}

// 0x8e

type mapper8e struct {
	baseMapper
}

func newMapper8e(bm *baseMapper) Mapper {
	return &mapper8e{baseMapper: *bm}
}

func (m *mapper8e) reset() {
}
