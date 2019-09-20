package core

// 0x50

type mapper50 struct {
	baseMapper
}

func newMapper50(bm *baseMapper) Mapper {
	return &mapper50{baseMapper: *bm}
}

func (m *mapper50) reset() {
}

// 0x51

type mapper51 struct {
	baseMapper
}

func newMapper51(bm *baseMapper) Mapper {
	return &mapper51{baseMapper: *bm}
}

func (m *mapper51) reset() {
}

// 0x52

type mapper52 struct {
	baseMapper
}

func newMapper52(bm *baseMapper) Mapper {
	return &mapper52{baseMapper: *bm}
}

func (m *mapper52) reset() {
}

// 0x53

type mapper53 struct {
	baseMapper
}

func newMapper53(bm *baseMapper) Mapper {
	return &mapper53{baseMapper: *bm}
}

func (m *mapper53) reset() {
}

// 0x54

type mapper54 struct {
	baseMapper
}

func newMapper54(bm *baseMapper) Mapper {
	return &mapper54{baseMapper: *bm}
}

func (m *mapper54) reset() {
}

// 0x55

type mapper55 struct {
	baseMapper
}

func newMapper55(bm *baseMapper) Mapper {
	return &mapper55{baseMapper: *bm}
}

func (m *mapper55) reset() {
}

// 0x56

type mapper56 struct {
	baseMapper
	r, c byte
}

func newMapper56(bm *baseMapper) Mapper {
	return &mapper56{baseMapper: *bm}
}

func (m *mapper56) reset() {
	m.mem.setProm32kBank(0)
	m.mem.setVrom8kBank(0)
	m.r, m.c = 0xff, 0
}

func (m *mapper56) write(addr uint16, data byte) {
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

func (m *mapper56) vSync() {
	if m.c != 0 {
		m.c--
	}
}

// 0x57

type mapper57 struct {
	baseMapper
}

func newMapper57(bm *baseMapper) Mapper {
	return &mapper57{baseMapper: *bm}
}

func (m *mapper57) reset() {
	m.mem.setProm32kBank(0)
	m.mem.setVrom8kBank(0)
}

func (m *mapper57) write(addr uint16, data byte) {
	if addr == 0x6000 {
		m.mem.setVrom8kBank(uint32(data&0x02) >> 1)
	}
}

// 0x58

type mapper58 struct {
	baseMapper
	r byte
}

func newMapper58(bm *baseMapper) Mapper {
	return &mapper58{baseMapper: *bm}
}

func (m *mapper58) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	if m.mem.nVrom1kPage >= 8 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper58) write(addr uint16, data byte) {
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

// 0x59

type mapper59 struct {
	baseMapper
}

func newMapper59(bm *baseMapper) Mapper {
	return &mapper59{baseMapper: *bm}
}

func (m *mapper59) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	m.mem.setVrom8kBank(0)
}

func (m *mapper59) write(addr uint16, data byte) {
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

// 0x5a

type mapper5a struct {
	baseMapper
}

func newMapper5a(bm *baseMapper) Mapper {
	return &mapper5a{baseMapper: *bm}
}

func (m *mapper5a) reset() {
}

// 0x5b

type mapper5b struct {
	baseMapper
}

func newMapper5b(bm *baseMapper) Mapper {
	return &mapper5b{baseMapper: *bm}
}

func (m *mapper5b) reset() {
}

// 0x5c

type mapper5c struct {
	baseMapper
}

func newMapper5c(bm *baseMapper) Mapper {
	return &mapper5c{baseMapper: *bm}
}

func (m *mapper5c) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper5c) write(addr uint16, data byte) {
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

// 0x5d

type mapper5d struct {
	baseMapper
}

func newMapper5d(bm *baseMapper) Mapper {
	return &mapper5d{baseMapper: *bm}
}

func (m *mapper5d) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper5d) write(addr uint16, data byte) {
	if addr == 0x6000 {
		m.mem.setProm16kBank(4, uint32(data))
	}
}

// 0x5e

type mapper5e struct {
	baseMapper
}

func newMapper5e(bm *baseMapper) Mapper {
	return &mapper5e{baseMapper: *bm}
}

func (m *mapper5e) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
}

func (m *mapper5e) write(addr uint16, data byte) {
	if (addr & 0xfff0) == 0xff00 {
		m.mem.setProm16kBank(4, uint32(data>>2)&0x07)
	}
}

// 0x5f

type mapper5f struct {
	baseMapper
}

func newMapper5f(bm *baseMapper) Mapper {
	return &mapper5f{baseMapper: *bm}
}

func (m *mapper5f) reset() {
}
