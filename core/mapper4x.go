package core

// 066

type mapper066 struct {
	baseMapper
}

func newMapper066(bm *baseMapper) Mapper {
	return &mapper066{baseMapper: *bm}
}

func (m *mapper066) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	m.mem.setVrom8kBank(0)
}

func (m *mapper066) writeLow(addr uint16, data byte) {
	if addr >= 0x6000 {
		m.mem.setProm32kBank(uint32(data&0xf0) >> 4)
		m.mem.setVrom8kBank(uint32(data & 0x0f))
	}
}

func (m *mapper066) write(addr uint16, data byte) {
	m.mem.setProm32kBank(uint32(data&0xf0) >> 4)
	m.mem.setVrom8kBank(uint32(data & 0x0f))
}

// 070

type mapper070 struct {
	baseMapper
}

func newMapper070(bm *baseMapper) Mapper {
	return &mapper070{baseMapper: *bm}
}

func (m *mapper070) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	m.mem.setVrom8kBank(0)
}

func (m *mapper070) write(addr uint16, data byte) {
	m.mem.setProm16kBank(4, uint32(data&0x70)>>4)
	m.mem.setVrom8kBank(uint32(data & 0x0f))
	if data&0x80 != 0 {
		m.mem.setVramMirror(memVramMirror4H)
	} else {
		m.mem.setVramMirror(memVramMirror4L)
	}
}

// 071

type mapper071 struct {
	baseMapper
}

func newMapper071(bm *baseMapper) Mapper {
	return &mapper071{baseMapper: *bm}
}

func (m *mapper071) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
}

func (m *mapper071) writeLow(addr uint16, data byte) {
	if (addr & 0xe000) == 0x6000 {
		m.mem.setProm16kBank(4, uint32(data))
	}
}

func (m *mapper071) write(addr uint16, data byte) {
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

// 072

type mapper072 struct {
	baseMapper
}

func newMapper072(bm *baseMapper) Mapper {
	return &mapper072{baseMapper: *bm}
}

func (m *mapper072) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper072) write(addr uint16, data byte) {
	if data&0x80 != 0 {
		m.mem.setProm16kBank(4, uint32(data&0x0f))
	} else if data&0x40 != 0 {
		m.mem.setVrom8kBank(uint32(data & 0x0f))
	}
}

// 073

type mapper073 struct {
	baseMapper
	irqEn  bool
	irqCnt uint32
}

func newMapper073(bm *baseMapper) Mapper {
	return &mapper073{baseMapper: *bm}
}

func (m *mapper073) reset() {
	m.irqEn, m.irqCnt = false, 0
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
}

func (m *mapper073) write(addr uint16, data byte) {
	switch addr {
	case 0x8000:
		m.irqCnt = (m.irqCnt & 0xfff0) | uint32(data&0x0f)
	case 0x9000:
		m.irqCnt = (m.irqCnt & 0xff0f) | (uint32(data&0x0f) << 4)
	case 0xa000:
		m.irqCnt = (m.irqCnt & 0xf0ff) | (uint32(data&0x0f) << 8)
	case 0xb000:
		m.irqCnt = (m.irqCnt & 0x0fff) | (uint32(data&0x0f) << 12)
	case 0xc000:
		m.irqEn = data&0x02 != 0
		m.clearIntr()
	case 0xd000:
		m.clearIntr()
	case 0xf000:
		m.mem.setProm16kBank(4, uint32(data))
	}
}

func (m *mapper073) clock(nCycle int64) {
	if m.irqEn {
		m.irqCnt += uint32(nCycle)
		if m.irqCnt >= 0xffff {
			m.irqEn = false
			m.irqCnt &= 0xffff
			m.setIntr()
		}
	}
}

// 076

type mapper076 struct {
	baseMapper
	r byte
}

func newMapper076(bm *baseMapper) Mapper {
	return &mapper076{baseMapper: *bm}
}

func (m *mapper076) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	if m.mem.nVrom1kPage >= 8 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper076) write(addr uint16, data byte) {
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

// 077

type mapper077 struct {
	baseMapper
}

func newMapper077(bm *baseMapper) Mapper {
	return &mapper077{baseMapper: *bm}
}

func (m *mapper077) reset() {
	m.mem.setProm32kBank(0)
	m.mem.setVrom2kBank(0, 0)
	m.mem.setCram2kBank(2, 1)
	m.mem.setCram2kBank(4, 2)
	m.mem.setCram2kBank(6, 3)
}

func (m *mapper077) write(addr uint16, data byte) {
	m.mem.setProm32kBank(uint32(data & 0x07))
	m.mem.setVrom2kBank(0, uint32(data&0xf0)>>4)
}

// 078

type mapper078 struct {
	baseMapper
}

func newMapper078(bm *baseMapper) Mapper {
	return &mapper078{baseMapper: *bm}
}

func (m *mapper078) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper078) write(addr uint16, data byte) {
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

// 079

type mapper079 struct {
	baseMapper
}

func newMapper079(bm *baseMapper) Mapper {
	return &mapper079{baseMapper: *bm}
}

func (m *mapper079) reset() {
	m.mem.setProm32kBank(0)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper079) write(addr uint16, data byte) {
	if addr&0x0100 != 0 {
		m.mem.setProm32kBank(uint32(data>>3) & 0x01)
		m.mem.setVrom8kBank(uint32(data & 0x07))
	}
}
