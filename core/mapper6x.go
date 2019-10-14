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
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
	if m.nVrom1kPage != 0 {
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
	if m.nProm8kPage > 2 {
		m.mem.setProm32kBank4(0, 1, 2, 3)
	} else if m.nProm8kPage > 1 {
		m.mem.setProm32kBank4(0, 1, 0, 1)
	} else {
		m.mem.setProm32kBank4(0, 0, 0, 0)
	}
	if m.nVrom1kPage != 0 {
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

// 100

type mapper100 struct {
	baseMapper
	irqEn    bool
	irqCnt   byte
	irqLatch byte
	r        byte
	p        [4]byte
	c        [8]byte
}

func newMapper100(bm *baseMapper) Mapper {
	return &mapper100{baseMapper: *bm}
}

func (m *mapper100) setPpuBanks() {
	if m.nVrom1kPage != 0 {
		for i := byte(0); i < 8; i++ {
			m.mem.setVrom1kBank(i, uint32(m.c[i]))
		}
	}
}

func (m *mapper100) reset() {
	m.irqEn, m.irqCnt, m.irqLatch = false, 0, 0
	m.r = 0
	m.p[0], m.p[1], m.p[2], m.p[3] = 0, 1, byte(m.nProm8kPage)-2, byte(m.nProm8kPage)-1
	if m.nVrom1kPage != 0 {
		for i := byte(0); i < 8; i++ {
			m.c[i] = i
		}
		m.setPpuBanks()
	} else {
		for i := byte(0); i < 8; i++ {
			m.c[i] = 0
		}
		m.c[1], m.c[3] = 1, 1
	}
	m.mem.setProm32kBank4(uint32(m.p[0]), uint32(m.p[1]), uint32(m.p[2]), uint32(m.p[3]))
}

func (m *mapper100) write(addr uint16, data byte) {
	switch addr & 0xe001 {
	case 0x8000:
		m.r = data
	case 0x8001:
		r := m.r & 0xc7
		switch r {
		case 0x00, 0x01:
			if m.nVrom1kPage != 0 {
				i := r << 1
				m.c[i] = data & 0xfe
				m.c[i+1] = m.c[i] + 1
				m.setPpuBanks()
			}
		case 0x02, 0x03, 0x04, 0x05:
			if m.nVrom1kPage != 0 {
				m.c[r+2] = data
				m.setPpuBanks()
			}
		case 0x06, 0x07, 0x46, 0x47:
			m.p[(r>>5)|(r&0x01)] = data
			m.setPpuBanks()
		case 0x80, 0x81:
			if m.nVrom1kPage != 0 {
				i := ((r & 0x01) << 1) + 4
				m.c[i] = data & 0xfe
				m.c[i+1] = m.c[i] + 1
				m.setPpuBanks()
			}
		case 0x82, 0x83, 0x84, 0x85:
			if m.nVrom1kPage != 0 {
				i := (r & 0x0f) - 2
				m.c[i] = data
				m.setPpuBanks()
			}
		}
	case 0xa000:
		if !m.sys.rom.b4Screen {
			if data&0x01 != 0 {
				m.mem.setVramMirror(memVramMirrorH)
			} else {
				m.mem.setVramMirror(memVramMirrorV)
			}
		}
	case 0xc000:
		m.irqCnt = data
	case 0xc001:
		m.irqLatch = data
	case 0xe000:
		m.irqEn = false
		m.clearIntr()
	case 0xe001:
		m.irqEn = true
	}
}

func (m *mapper100) hSync(scanline uint16) {
	if scanline < ScreenHeight && m.isPpuDisp() && m.irqEn {
		m.irqCnt--
		if m.irqCnt == 0xff {
			m.irqCnt = m.irqLatch
			m.setIntr()
		}
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

// 105

type mapper105 struct {
	baseMapper
	irqEn    bool
	irqCnt   int32
	state    byte
	writeCnt byte
	bits     byte
	r        [4]byte
}

func newMapper105(bm *baseMapper) Mapper {
	return &mapper105{baseMapper: *bm}
}

func (m *mapper105) reset() {
	m.irqEn, m.irqCnt = false, 0
	m.state, m.writeCnt, m.bits = 0, 0, 0
	m.r[0], m.r[1], m.r[2], m.r[3] = 0x0c, 0, 0, 0x10
	m.mem.setProm32kBank(0)
}

func (m *mapper105) write(addr uint16, data byte) {
	i := (addr & 0x7fff) >> 13
	if data&0x80 != 0 {
		m.bits, m.writeCnt = 0, 0
		if i == 0 {
			m.r[0] |= 0x0c
		}
	} else {
		m.bits |= (data & 0x01) << m.writeCnt
		m.writeCnt++
		if m.writeCnt == 5 {
			m.r[i], m.bits, m.writeCnt = m.bits&0x1f, 0, 0
		}
	}
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

	switch m.state {
	case 0, 1:
		m.state++
	case 2:
		if m.r[1]&0x08 != 0 {
			if m.r[0]&0x08 != 0 {
				r := uint32(m.r[3]&0x07) << 1
				if m.r[0]&0x04 != 0 {
					m.mem.setProm32kBank4(r+16, r+17, 30, 31)
				} else {
					m.mem.setProm32kBank4(16, 17, r+16, r+17)
				}
			} else {
				r := uint32(m.r[3]&0x06) << 1
				m.mem.setProm32kBank4(r+16, r+17, r+18, r+19)
			}
		} else {
			r := uint32(m.r[1]&0x06) << 1
			m.mem.setProm32kBank4(r, r+1, r+2, r+3)
		}
		if m.r[1]&0x10 != 0 {
			m.irqEn, m.irqCnt = false, 0
		} else {
			m.irqEn = true
		}
	}
}

func (m *mapper105) hSync(scanline uint16) {
	if scanline == 0 {
		if m.irqEn {
			m.irqCnt += 29781
		}
		if (m.irqCnt|0x21ffffff)&0x3e000000 == 0x3e000000 {
			m.sys.cpu.intr |= cpuIntrTypTrig2
		}
	}
}

// 107

type mapper107 struct {
	baseMapper
}

func newMapper107(bm *baseMapper) Mapper {
	return &mapper107{baseMapper: *bm}
}

func (m *mapper107) reset() {
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
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

// 109

type mapper109 struct {
	baseMapper
	r     byte
	c     [4]byte
	mode0 byte
	mode1 byte
}

func newMapper109(bm *baseMapper) Mapper {
	return &mapper109{baseMapper: *bm}
}

func (m *mapper109) setPpuBanks() {
	if m.nVrom1kPage != 0 {
		m.mem.setVrom1kBank(0, uint32(m.c[0]))
		m.mem.setVrom1kBank(1, uint32(m.c[1]|((m.mode1<<3)&0x08)))
		m.mem.setVrom1kBank(2, uint32(m.c[2]|((m.mode1<<2)&0x08)))
		m.mem.setVrom1kBank(3, uint32(m.c[3]|((m.mode1<<1)&0x08)|(m.mode0<<4)))
		m.mem.setVrom1kBank(4, m.nVrom1kPage-4)
		m.mem.setVrom1kBank(5, m.nVrom1kPage-3)
		m.mem.setVrom1kBank(6, m.nVrom1kPage-2)
		m.mem.setVrom1kBank(7, m.nVrom1kPage-1)
	}
}

func (m *mapper109) reset() {
	m.r = 0
	m.c[0], m.c[1], m.c[2], m.c[3] = 0, 0, 0, 0
	m.mode0, m.mode1 = 0, 0
	m.mem.setProm32kBank(0)
	m.setPpuBanks()
}

func (m *mapper109) writeLow(addr uint16, data byte) {
	switch addr {
	case 0x4100:
		m.r = data
	case 0x4101:
		switch m.r {
		case 0x00, 0x01, 0x02, 0x03:
			m.c[m.r] = data
			m.setPpuBanks()
		case 0x04:
			m.mode0 = data & 0x01
			m.setPpuBanks()
		case 0x05:
			m.mem.setProm32kBank(uint32(data & 0x07))
		case 0x06:
			m.mode1 = data & 0x07
			m.setPpuBanks()
		case 0x07:
			if data&0x01 != 0 {
				m.mem.setVramMirror(memVramMirrorH)
			} else {
				m.mem.setVramMirror(memVramMirrorV)
			}
		}
	}

}

// 110

type mapper110 struct {
	baseMapper
	r0, r1 byte
}

func newMapper110(bm *baseMapper) Mapper {
	return &mapper110{baseMapper: *bm}
}

func (m *mapper110) reset() {
	m.r0, m.r1 = 0, 0
	m.mem.setProm32kBank(0)
	m.mem.setVrom8kBank(0)
}

func (m *mapper110) writeLow(addr uint16, data byte) {
	switch addr {
	case 0x4100:
		m.r1 = data & 0x07
	case 0x4101:
		switch m.r1 {
		case 0x00:
			m.r0 = data & 0x01
			m.mem.setVrom8kBank(uint32(m.r0))
		case 0x02:
			m.r0 = data
			m.mem.setVrom8kBank(uint32(m.r0))
		case 0x04:
			m.r0 = m.r0 | (data << 1)
			m.mem.setVrom8kBank(uint32(m.r0))
		case 0x05:
			m.mem.setProm32kBank(uint32(data))
		case 0x06:
			m.r0 = m.r0 | (data << 2)
			m.mem.setVrom8kBank(uint32(m.r0))
		}
	}
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
	if m.nProm8kPage < 64 {
		m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
	} else {
		m.mem.setProm16kBank(4, 0)
		m.mem.setProm16kBank(6, 15)
		m.largeTyp = true
	}
	if m.nVrom1kPage != 0 {
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
			if m.nVrom1kPage != 0 {
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
				m.mem.setProm16kBank(6, (m.nProm8kPage>>1)-1)
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

		if m.nVrom1kPage != 0 {
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
