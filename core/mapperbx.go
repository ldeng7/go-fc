package core

// 176

type mapper176 struct {
	baseMapper
	patch bool
	we    byte
	sb    bool
}

func newMapper176(bm *baseMapper) Mapper {
	return &mapper176{baseMapper: *bm}
}

func (m *mapper176) reset() {
	if m.sys.conf.PatchTyp&0x01 != 0 {
		m.patch = true
	}
	m.sb, m.we = false, 0
	if m.mem.nProm8kPage > 64 {
		m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage/2-2, m.mem.nProm8kPage/2-1)
	} else {
		m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	}
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper176) readLow(addr uint16) byte {
	if m.patch && addr >= 0x6000 {
		switch m.we {
		case 0xe4, 0xec, 0xe5, 0xed, 0xe6, 0xee, 0xe7, 0xef:
			return m.mem.wram[(addr&0x1fff)+(uint16(m.we)&0x03<<13)]
		default:
			return m.cpuBanks[addr>>13][addr&0x1fff]
		}
	}
	return m.baseMapper.readLow(addr)
}

func (m *mapper176) writeLow(addr uint16, data byte) {
	switch addr {
	case 0x5010:
		if data == 0x24 {
			m.sb = true
		}
	case 0x5011:
		if m.sb {
			m.mem.setProm32kBank(uint32(data >> 1))
		}
	case 0x5ff1:
		m.mem.setProm32kBank(uint32(data >> 1))
	case 0x5ff2:
		m.mem.setVrom8kBank(uint32(data))
	}
	if addr >= 0x6000 {
		if m.patch {
			switch m.we {
			case 0xe4, 0xec:
				m.cpuBanks[addr>>13][addr&0x1fff] = data
				m.mem.wram[addr&0x1fff] = data
			case 0xe5, 0xed, 0xe6, 0xee, 0xe7, 0xef:
				m.mem.wram[(addr&0x1fff)+(uint16(m.we)&0x03<<13)] = data
			default:
				m.cpuBanks[addr>>13][addr&0x1fff] = data
			}
		} else {
			m.cpuBanks[addr>>13][addr&0x1fff] = data
		}
	}
}

func (m *mapper176) write(addr uint16, data byte) {
	switch addr {
	case 0xa000:
		switch data & 0x03 {
		case 0x00:
			m.mem.setVramMirror(memVramMirrorV)
		case 0x01:
			m.mem.setVramMirror(memVramMirrorH)
		case 0x02:
			m.mem.setVramMirror(memVramMirror4L)
		case 0x03:
			m.mem.setVramMirror(memVramMirror4H)
		}
	case 0xa001:
		m.we = data
	}
}

// 178

type mapper178 struct {
	baseMapper
	patch      bool
	r0, r1, r2 byte
}

func newMapper178(bm *baseMapper) Mapper {
	return &mapper178{baseMapper: *bm}
}

func (m *mapper178) reset() {
	if m.sys.conf.PatchTyp&0x01 != 0 {
		m.patch = true
	}
	m.r0, m.r1, m.r2 = 0, 0, 0
	m.mem.setProm32kBank(uint32(m.r0 | m.r1))
}

func (m *mapper178) writeLow(addr uint16, data byte) {
	switch addr {
	case 0x4800:
		if data&0x01 != 0 {
			m.mem.setVramMirror(memVramMirrorH)
		} else {
			m.mem.setVramMirror(memVramMirrorV)
		}
	case 0x4801:
		m.r0 = (data >> 1) & 0x0f
		if m.patch {
			m.r0 = data << 2
		}
		m.mem.setProm32kBank(uint32(m.r0 | m.r1))
	case 0x4802:
		m.r1 = data << 2
		if m.patch {
			m.r1 = data
		}
		m.mem.setProm32kBank(uint32(m.r0 | m.r1))
	case 0x4803:
		m.r2 = data
	}
	if addr >= 0x6000 {
		m.baseMapper.writeLow(addr, data)
	}
}

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

// 184

type mapper184 struct {
	baseMapper
}

func newMapper184(bm *baseMapper) Mapper {
	return &mapper184{baseMapper: *bm}
}

func (m *mapper184) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper184) writeLow(addr uint16, data byte) {
	if addr >= 0x6000 {
		m.mem.setVrom4kBank(0, uint32(((data&0x02)<<2)|(data&0x04)))
		m.mem.setVrom4kBank(1, uint32((data&0x20)>>2))
	}
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
