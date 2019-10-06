package core

// 086

type mapper086 struct {
	baseMapper
	r, c byte
}

func newMapper086(bm *baseMapper) Mapper {
	return &mapper086{baseMapper: *bm}
}

func (m *mapper086) reset() {
	m.mem.setProm32kBank(0)
	m.mem.setVrom8kBank(0)
	m.r, m.c = 0xff, 0
}

func (m *mapper086) write(addr uint16, data byte) {
	switch addr {
	case 0x6000:
		m.mem.setProm32kBank(uint32(data&0x30) >> 4)
		m.mem.setVrom8kBank(uint32(data&0x03) | (uint32(data&0x40) >> 4))
	case 0x7000:
		if m.r&0x10 == 0 && data&0x10 != 0 && m.c == 0 && (data&0x0f == 0 || data&0x0f == 0x05) {
			m.c = 60
		}
		m.r = data
	}
}

func (m *mapper086) vSync() {
	if m.c != 0 {
		m.c--
	}
}

// 087

type mapper087 struct {
	baseMapper
}

func newMapper087(bm *baseMapper) Mapper {
	return &mapper087{baseMapper: *bm}
}

func (m *mapper087) reset() {
	m.mem.setProm32kBank(0)
	m.mem.setVrom8kBank(0)
}

func (m *mapper087) write(addr uint16, data byte) {
	if addr == 0x6000 {
		m.mem.setVrom8kBank(uint32(data&0x02) >> 1)
	}
}

// 088

type mapper088 struct {
	baseMapper
	r byte
}

func newMapper088(bm *baseMapper) Mapper {
	return &mapper088{baseMapper: *bm}
}

func (m *mapper088) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	if m.mem.nVrom1kPage >= 8 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper088) write(addr uint16, data byte) {
	switch addr {
	case 0x8000:
		m.r = data
	case 0x8001:
		b, r := uint32(data), m.r&0x07
		switch r {
		case 0x00, 0x01:
			m.mem.setVrom2kBank(r<<1, b>>1)
		case 0x02, 0x03, 0x04, 0x05:
			m.mem.setVrom1kBank(r+2, b+0x40)
		case 0x06, 0x07:
			m.mem.setProm8kBank(r+2, b)
		}
	case 0xc000:
		if data != 0 {
			m.mem.setVramMirror(memVramMirror4H)
		} else {
			m.mem.setVramMirror(memVramMirror4L)
		}
	}
}

// 089

type mapper089 struct {
	baseMapper
}

func newMapper089(bm *baseMapper) Mapper {
	return &mapper089{baseMapper: *bm}
}

func (m *mapper089) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	m.mem.setVrom8kBank(0)
}

func (m *mapper089) write(addr uint16, data byte) {
	if (addr & 0xff00) == 0xc000 {
		m.mem.setProm16kBank(4, uint32(data&0x70)>>4)
		m.mem.setVrom8kBank((uint32(data&0x80) >> 4) | uint32(data&0x07))
		if data&0x08 != 0 {
			m.mem.setVramMirror(memVramMirror4H)
		} else {
			m.mem.setVramMirror(memVramMirror4L)
		}
	}
}

// 092

type mapper092 struct {
	baseMapper
}

func newMapper092(bm *baseMapper) Mapper {
	return &mapper092{baseMapper: *bm}
}

func (m *mapper092) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper092) write(addr uint16, data byte) {
	data = byte(addr)
	if addr >= 0x9000 {
		if (data & 0xf0) == 0xd0 {
			m.mem.setProm16kBank(6, uint32(data&0x0f))
		} else if (data & 0xf0) == 0xe0 {
			m.mem.setVrom8kBank(uint32(data & 0x0f))
		}
	} else {
		if (data & 0xf0) == 0xb0 {
			m.mem.setProm16kBank(6, uint32(data&0x0f))
		} else if (data & 0xf0) == 0x70 {
			m.mem.setVrom8kBank(uint32(data & 0x0f))
		}
	}
}

// 093

type mapper093 struct {
	baseMapper
}

func newMapper093(bm *baseMapper) Mapper {
	return &mapper093{baseMapper: *bm}
}

func (m *mapper093) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper093) write(addr uint16, data byte) {
	if addr == 0x6000 {
		m.mem.setProm16kBank(4, uint32(data))
	}
}

// 094

type mapper094 struct {
	baseMapper
}

func newMapper094(bm *baseMapper) Mapper {
	return &mapper094{baseMapper: *bm}
}

func (m *mapper094) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
}

func (m *mapper094) write(addr uint16, data byte) {
	if (addr & 0xfff0) == 0xff00 {
		m.mem.setProm16kBank(4, uint32(data>>2)&0x07)
	}
}
