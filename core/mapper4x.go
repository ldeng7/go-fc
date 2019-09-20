package core

// 0x40

type mapper40 struct {
	baseMapper
}

func newMapper40(bm *baseMapper) Mapper {
	return &mapper40{baseMapper: *bm}
}

func (m *mapper40) reset() {
}

// 0x41

type mapper41 struct {
	baseMapper
}

func newMapper41(bm *baseMapper) Mapper {
	return &mapper41{baseMapper: *bm}
}

func (m *mapper41) reset() {
}

// 0x42

type mapper42 struct {
	baseMapper
}

func newMapper42(bm *baseMapper) Mapper {
	return &mapper42{baseMapper: *bm}
}

func (m *mapper42) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	m.mem.setVrom8kBank(0)
}

func (m *mapper42) writeLow(addr uint16, data byte) {
	if addr >= 0x6000 {
		m.mem.setProm32kBank(uint32(data&0xf0) >> 4)
		m.mem.setVrom8kBank(uint32(data & 0x0f))
	}
}

func (m *mapper42) write(addr uint16, data byte) {
	m.mem.setProm32kBank(uint32(data&0xf0) >> 4)
	m.mem.setVrom8kBank(uint32(data & 0x0f))
}

// 0x43

type mapper43 struct {
	baseMapper
}

func newMapper43(bm *baseMapper) Mapper {
	return &mapper43{baseMapper: *bm}
}

func (m *mapper43) reset() {
}

// 0x44

type mapper44 struct {
	baseMapper
}

func newMapper44(bm *baseMapper) Mapper {
	return &mapper44{baseMapper: *bm}
}

func (m *mapper44) reset() {
}

// 0x45

type mapper45 struct {
	baseMapper
}

func newMapper45(bm *baseMapper) Mapper {
	return &mapper45{baseMapper: *bm}
}

func (m *mapper45) reset() {
}

// 0x46

type mapper46 struct {
	baseMapper
}

func newMapper46(bm *baseMapper) Mapper {
	return &mapper46{baseMapper: *bm}
}

func (m *mapper46) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	m.mem.setVrom8kBank(0)
}

func (m *mapper46) write(addr uint16, data byte) {
	m.mem.setProm16kBank(4, uint32(data&0x70)>>4)
	m.mem.setVrom8kBank(uint32(data & 0x0f))
	if data&0x80 != 0 {
		m.mem.setVramMirror(memVramMirror4H)
	} else {
		m.mem.setVramMirror(memVramMirror4L)
	}
}

// 0x47

type mapper47 struct {
	baseMapper
}

func newMapper47(bm *baseMapper) Mapper {
	return &mapper47{baseMapper: *bm}
}

func (m *mapper47) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
}

func (m *mapper47) writeLow(addr uint16, data byte) {
	if (addr & 0xe000) == 0x6000 {
		m.mem.setProm16kBank(4, uint32(data))
	}
}

func (m *mapper47) write(addr uint16, data byte) {
	switch addr & 0xf000 {
	case 0x9000:
		if data&0x10 != 0 {
			m.mem.setVramMirror(memVramMirror4H)
		} else {
			m.mem.setVramMirror(memVramMirror4L)
		}
	case 0xc000, 0xd000, 0xe000, 0xf000:
		m.mem.setProm16kBank(4, uint32(data))
	}
}

// 0x48

type mapper48 struct {
	baseMapper
}

func newMapper48(bm *baseMapper) Mapper {
	return &mapper48{baseMapper: *bm}
}

func (m *mapper48) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper48) write(addr uint16, data byte) {
	if data&0x80 != 0 {
		m.mem.setProm16kBank(4, uint32(data&0x0f))
	} else if data&0x40 != 0 {
		m.mem.setVrom8kBank(uint32(data & 0x0f))
	}
}

// 0x49

type mapper49 struct {
	baseMapper
}

func newMapper49(bm *baseMapper) Mapper {
	return &mapper49{baseMapper: *bm}
}

func (m *mapper49) reset() {
}

// 0x4a

type mapper4a struct {
	baseMapper
}

func newMapper4a(bm *baseMapper) Mapper {
	return &mapper4a{baseMapper: *bm}
}

func (m *mapper4a) reset() {
}

// 0x4b

type mapper4b struct {
	baseMapper
}

func newMapper4b(bm *baseMapper) Mapper {
	return &mapper4b{baseMapper: *bm}
}

func (m *mapper4b) reset() {
}

// 0x4c

type mapper4c struct {
	baseMapper
	r byte
}

func newMapper4c(bm *baseMapper) Mapper {
	return &mapper4c{baseMapper: *bm}
}

func (m *mapper4c) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	if m.mem.nVrom1kPage >= 8 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper4c) write(addr uint16, data byte) {
	switch addr {
	case 0x8000:
		m.r = data
	case 0x8001:
		b, r := uint32(data), m.r&0x07
		switch r {
		case 0x02, 0x03, 0x04, 0x05:
			m.mem.setVrom2kBank((r-2)<<1, b)
		case 0x06, 0x07:
			m.mem.setProm8kBank(r+2, b)
		}
	}
}

// 0x4d

type mapper4d struct {
	baseMapper
}

func newMapper4d(bm *baseMapper) Mapper {
	return &mapper4d{baseMapper: *bm}
}

func (m *mapper4d) reset() {
	m.mem.setProm32kBank(0)
	m.mem.setVrom2kBank(0, 0)
	m.mem.setCram2kBank(2, 1)
	m.mem.setCram2kBank(4, 2)
	m.mem.setCram2kBank(6, 3)
}

func (m *mapper4d) write(addr uint16, data byte) {
	m.mem.setProm32kBank(uint32(data & 0x07))
	m.mem.setVrom2kBank(0, uint32(data&0xf0)>>4)
}

// 0x4e

type mapper4e struct {
	baseMapper
}

func newMapper4e(bm *baseMapper) Mapper {
	return &mapper4e{baseMapper: *bm}
}

func (m *mapper4e) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper4e) write(addr uint16, data byte) {
	m.mem.setProm16kBank(4, uint32(data&0x0f))
	m.mem.setVrom8kBank(uint32(data&0xf0) >> 4)
	if (addr & 0xfe00) != 0xfe00 {
		if data&0x08 != 0 {
			m.mem.setVramMirror(memVramMirror4H)
		} else {
			m.mem.setVramMirror(memVramMirror4L)
		}
	}
}

// 0x4f

type mapper4f struct {
	baseMapper
}

func newMapper4f(bm *baseMapper) Mapper {
	return &mapper4f{baseMapper: *bm}
}

func (m *mapper4f) reset() {
	m.mem.setProm32kBank(0)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper4f) write(addr uint16, data byte) {
	if addr&0x0100 != 0 {
		m.mem.setProm32kBank(uint32(data>>3) & 0x01)
		m.mem.setVrom8kBank(uint32(data & 0x07))
	}
}
