package core

// 0xb4

type mapperb4 struct {
	baseMapper
}

func newMapperb4(bm *baseMapper) Mapper {
	return &mapperb4{baseMapper: *bm}
}

func (m *mapperb4) reset() {
	m.mem.setProm32kBank(0)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapperb4) write(addr uint16, data byte) {
	m.mem.setProm16kBank(6, uint32(data&0x07))
}

// 0xb5

type mapperb5 struct {
	baseMapper
}

func newMapperb5(bm *baseMapper) Mapper {
	return &mapperb5{baseMapper: *bm}
}

func (m *mapperb5) reset() {
	m.mem.setProm32kBank(0)
	m.mem.setVrom8kBank(0)
}

func (m *mapperb5) writeLow(addr uint16, data byte) {
	if addr == 0x4120 {
		m.mem.setProm32kBank(uint32(data&0x08) >> 3)
		m.mem.setVrom8kBank(uint32(data & 0x07))
	}
}

// 0xb6

type mapperb6 struct {
	baseMapper
}

func newMapperb6(bm *baseMapper) Mapper {
	return &mapperb6{baseMapper: *bm}
}

func (m *mapperb6) reset() {
}

// 0xb7

type mapperb7 struct {
	baseMapper
}

func newMapperb7(bm *baseMapper) Mapper {
	return &mapperb7{baseMapper: *bm}
}

func (m *mapperb7) reset() {
}

// 0xb8

type mapperb8 struct {
	baseMapper
}

func newMapperb8(bm *baseMapper) Mapper {
	return &mapperb8{baseMapper: *bm}
}

func (m *mapperb8) reset() {
}

// 0xb9

type mapperb9 struct {
	baseMapper
}

func newMapperb9(bm *baseMapper) Mapper {
	return &mapperb9{baseMapper: *bm}
}

func (m *mapperb9) reset() {
	switch m.mem.nProm8kPage >> 1 {
	case 1:
		m.mem.setProm16kBank(4, 0)
		m.mem.setProm16kBank(6, 0)
	case 2:
		m.mem.setProm32kBank(0)
	}
	sl := m.mem.vram[0x0800:]
	for i := 0; i < 0x0400; i++ {
		sl[i] = 0xff
	}
}

func (m *mapperb9) write(addr uint16, data byte) {
	if data&0x03 != 0 {
		m.mem.setVrom8kBank(0)
	} else {
		for i := byte(0); i < 8; i++ {
			m.mem.setVram1kBank(i, 2)
		}
	}
}

// 0xba

type mapperba struct {
	baseMapper
}

func newMapperba(bm *baseMapper) Mapper {
	return &mapperba{baseMapper: *bm}
}

func (m *mapperba) reset() {
}

// 0xbb

type mapperbb struct {
	baseMapper
}

func newMapperbb(bm *baseMapper) Mapper {
	return &mapperbb{baseMapper: *bm}
}

func (m *mapperbb) reset() {
}

// 0xbc

type mapperbc struct {
	baseMapper
}

func newMapperbc(bm *baseMapper) Mapper {
	return &mapperbc{baseMapper: *bm}
}

func (m *mapperbc) reset() {
	if m.mem.nProm8kPage > 16 {
		m.mem.setProm32kBank4(0, 1, 14, 15)
	} else {
		m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	}
}

func (m *mapperbc) write(addr uint16, data byte) {
	if data&0x10 != 0 {
		m.mem.setProm16kBank(4, uint32(data&0x07))
	} else if data != 0 {
		m.mem.setProm16kBank(4, uint32(data)+0x08)
	} else if m.mem.nProm8kPage == 16 {
		m.mem.setProm16kBank(4, 7)
	} else {
		m.mem.setProm16kBank(4, 8)
	}
}

// 0xbd

type mapperbd struct {
	baseMapper
}

func newMapperbd(bm *baseMapper) Mapper {
	return &mapperbd{baseMapper: *bm}
}

func (m *mapperbd) reset() {
}

// 0xbe

type mapperbe struct {
	baseMapper
}

func newMapperbe(bm *baseMapper) Mapper {
	return &mapperbe{baseMapper: *bm}
}

func (m *mapperbe) reset() {
}

// 0xbf

type mapperbf struct {
	baseMapper
}

func newMapperbf(bm *baseMapper) Mapper {
	return &mapperbf{baseMapper: *bm}
}

func (m *mapperbf) reset() {
}
