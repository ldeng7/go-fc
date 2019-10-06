package core

// 096

type mapper096 struct {
	baseMapper
	r0, r1 byte
}

func newMapper096(bm *baseMapper) Mapper {
	return &mapper096{baseMapper: *bm}
}

func (m *mapper096) reset() {
	m.r0, m.r1 = 0, 0
	m.mem.setProm32kBank(0)
	m.mem.setCram4kBank(0, uint32(m.r0)<<2+uint32(m.r1))
	m.mem.setCram4kBank(0, uint32(m.r0)<<2+0x03)
	m.mem.setVramMirror(memVramMirror4L)
}

func (m *mapper096) write(addr uint16, data byte) {
	m.mem.setProm32kBank(uint32(data & 0x03))
	m.r0 = (data & 0x04) >> 2
	m.mem.setCram4kBank(0, uint32(m.r0)<<2+uint32(m.r1))
	m.mem.setCram4kBank(0, uint32(m.r0)<<2+0x03)
}

func (m *mapper096) ppuLatch(addr uint16) {
	if (addr & 0xf000) == 0x2000 {
		m.r1 = byte(addr>>8) & 0x03
		m.mem.setCram4kBank(0, uint32(m.r0)<<2+uint32(m.r1))
		m.mem.setCram4kBank(0, uint32(m.r0)<<2+0x03)
	}
}

// 097

type mapper097 struct {
	baseMapper
}

func newMapper097(bm *baseMapper) Mapper {
	return &mapper097{baseMapper: *bm}
}

func (m *mapper097) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper097) write(addr uint16, data byte) {
	if addr < 0xc000 {
		m.mem.setProm16kBank(6, uint32(data&0x0f))
		if data&0x08 != 0 {
			m.mem.setVramMirror(memVramMirrorV)
		} else {
			m.mem.setVramMirror(memVramMirrorH)
		}
	}
}

// 099

type mapper099 struct {
	baseMapper
	c byte
}

func newMapper099(bm *baseMapper) Mapper {
	return &mapper099{baseMapper: *bm}
}

func (m *mapper099) reset() {
	m.c = 0
	if m.mem.nProm8kPage > 2 {
		m.mem.setProm32kBank4(0, 1, 2, 3)
	} else if m.mem.nProm8kPage > 1 {
		m.mem.setProm32kBank4(0, 1, 0, 1)
	} else {
		m.mem.setProm32kBank4(0, 0, 0, 0)
	}
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper099) readEx(addr uint16) byte {
	if addr == 0x4020 {
		return m.c
	}
	return byte(addr >> 8)
}

func (m *mapper099) writeEx(addr uint16, data byte) {
	switch addr {
	case 0x4016:
		if data&0x04 != 0 {
			m.mem.setVrom8kBank(1)
		} else {
			m.mem.setVrom8kBank(0)
		}
	case 0x4020:
		m.c = data
	}
}

// 101

type mapper101 struct {
	baseMapper
}

func newMapper101(bm *baseMapper) Mapper {
	return &mapper101{baseMapper: *bm}
}

func (m *mapper101) reset() {
	m.mem.setProm32kBank(0)
	m.mem.setVrom8kBank(0)
}

func (m *mapper101) writeLow(addr uint16, data byte) {
	if addr >= 0x6000 {
		m.mem.setVrom8kBank(uint32(data & 0x03))
	}
}

func (m *mapper101) write(addr uint16, data byte) {
	m.mem.setVrom8kBank(uint32(data & 0x03))
}

// 107

type mapper107 struct {
	baseMapper
}

func newMapper107(bm *baseMapper) Mapper {
	return &mapper107{baseMapper: *bm}
}

func (m *mapper107) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	m.mem.setVrom8kBank(0)
}

func (m *mapper107) write(addr uint16, data byte) {
	m.mem.setProm32kBank(uint32(data>>1) & 0x03)
	m.mem.setVrom8kBank(uint32(data & 0x07))
}

// 108

type mapper108 struct {
	baseMapper
}

func newMapper108(bm *baseMapper) Mapper {
	return &mapper108{baseMapper: *bm}
}

func (m *mapper108) reset() {
	m.mem.setProm32kBank4(12, 13, 14, 15)
	m.mem.setProm8kBank(3, 0)
}

func (m *mapper108) write(addr uint16, data byte) {
	m.mem.setProm8kBank(3, uint32(data))
}

// 111

type mapper111 struct {
	baseMapper
	largeTyp bool
	r        [4]byte
}

func newMapper111(bm *baseMapper) Mapper {
	return &mapper111{baseMapper: *bm}
}

func (m *mapper111) reset() {
	m.r[0], m.r[1], m.r[2], m.r[3] = 0x0c, 0, 0, 0
	if m.mem.nProm8kPage < 64 {
		m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	} else {
		m.mem.setProm16kBank(4, 0)
		m.mem.setProm16kBank(6, 15)
		m.largeTyp = true
	}
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper111) setPpuMirrors() {
	if m.r[0]&0x02 != 0 {
		if m.r[0]&0x01 != 0 {
			m.mem.setVramMirror(memVramMirrorH)
		} else {
			m.mem.setVramMirror(memVramMirrorV)
		}
	} else {
		if m.r[0]&0x01 != 0 {
			m.mem.setVramMirror(memVramMirror4H)
		} else {
			m.mem.setVramMirror(memVramMirror4L)
		}
	}
}

func (m *mapper111) write(addr uint16, data byte) {
	if data&0x80 != 0 {
		m.r[0] |= 0x0c
		return
	}

	addr = (addr >> 13) & 0x03
	m.r[addr] = data
	if !m.largeTyp {
		switch addr {
		case 0x00:
			m.setPpuMirrors()
		case 0x01, 0x02:
			if m.mem.nVrom1kPage != 0 {
				if m.r[0]&0x10 != 0 {
					m.mem.setVrom4kBank(0, uint32(m.r[1]))
					m.mem.setVrom4kBank(4, uint32(m.r[2]))
				} else {
					m.mem.setVrom8kBank(uint32(m.r[1] >> 1))
				}
			} else if m.r[0]&0x10 != 0 {
				m.mem.setCram4kBank((byte(addr)-0x01)<<2, uint32(m.r[addr]))
			}
		case 0x03:
			if m.r[0]&0x08 == 0 {
				m.mem.setProm32kBank(uint32(m.r[3] >> 1))
			} else if m.r[0]&0x04 != 0 {
				m.mem.setProm16kBank(4, uint32(m.r[3]))
				m.mem.setProm16kBank(6, (m.mem.nProm8kPage>>1)-1)
			} else {
				m.mem.setProm16kBank(4, 0)
				m.mem.setProm16kBank(6, uint32(m.r[3]))
			}
		}
	} else {
		b := uint32(m.r[1] & 0x10)
		if m.r[0]&0x08 == 0 {
			m.mem.setProm32kBank((uint32(m.r[3]) & (b | 0x0f)) >> 1)
		} else if m.r[0]&0x04 != 0 {
			m.mem.setProm16kBank(4, uint32(m.r[3]&0x0f)|b)
			m.mem.setProm16kBank(6, b+15)
		} else {
			m.mem.setProm16kBank(4, b)
			m.mem.setProm16kBank(6, uint32(m.r[3]&0x0f)|b)
		}

		if m.mem.nVrom1kPage != 0 {
			if m.r[0]&0x10 != 0 {
				m.mem.setVrom4kBank(0, uint32(m.r[1]))
				m.mem.setVrom4kBank(4, uint32(m.r[2]))
			} else {
				m.mem.setVrom8kBank(uint32(m.r[1] >> 1))
			}
		} else if m.r[0]&0x10 != 0 {
			m.mem.setCram4kBank(0, uint32(m.r[1]))
			m.mem.setCram4kBank(4, uint32(m.r[2]))
		}

		if addr == 0 {
			m.setPpuMirrors()
		}
	}
}
