package core

// 0x20

type mapper20 struct {
	baseMapper
}

func newMapper20(bm *baseMapper) Mapper {
	return &mapper20{baseMapper: *bm}
}

func (m *mapper20) reset() {
}

// 0x21

type mapper21 struct {
	baseMapper
}

func newMapper21(bm *baseMapper) Mapper {
	return &mapper21{baseMapper: *bm}
}

func (m *mapper21) reset() {
}

// 0x22

type mapper22 struct {
	baseMapper
}

func newMapper22(bm *baseMapper) Mapper {
	return &mapper22{baseMapper: *bm}
}

func (m *mapper22) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper22) writeLow(addr uint16, data byte) {
	switch addr {
	case 0x7ffd:
		m.mem.setProm32kBank(uint32(data))
	case 0x7ffe:
		m.mem.setVrom4kBank(0, uint32(data))
	case 0x7fff:
		m.mem.setVrom4kBank(4, uint32(data))
	}
}

func (m *mapper22) write(addr uint16, data byte) {
	m.mem.setProm32kBank(uint32(data))
}

// 0x28

type mapper28 struct {
	baseMapper
	irqEn   bool
	irqLine int32
}

func newMapper28(bm *baseMapper) Mapper {
	return &mapper28{baseMapper: *bm}
}

func (m *mapper28) reset() {
	m.mem.setProm8kBank(3, 6)
	m.mem.setProm32kBank4(4, 5, 0, 7)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper28) write(addr uint16, data byte) {
	switch addr & 0xe000 {
	case 0x8000:
		m.irqEn = false
		m.sys.cpu.intr &^= cpuIntrTypMapper
	case 0xa000:
		m.irqEn, m.irqLine = true, 37
		m.sys.cpu.intr &^= cpuIntrTypMapper
	case 0xe000:
		m.mem.setProm8kBank(6, uint32(data)&0x07)
	}
}

func (m *mapper28) hSync(scanline uint16) {
	if m.irqEn {
		m.irqLine--
		if m.irqLine <= 0 {
			m.sys.cpu.intr |= cpuIntrTypMapper
		}
	}
}

// 0x29

type mapper29 struct {
	baseMapper
	r0, r1 byte
}

func newMapper29(bm *baseMapper) Mapper {
	return &mapper29{baseMapper: *bm}
}

func (m *mapper29) reset() {
	m.mem.setProm32kBank(0)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper29) writeLow(addr uint16, data byte) {
	if addr >= 0x6000 && addr < 0x6800 {
		m.mem.setProm32kBank(uint32(addr) & 0x07)
		m.r0 = byte(addr) & 0x04
		m.r1 = (m.r1 & 0x03) | (byte(addr>>1) & 0x0c)
		m.mem.setVrom8kBank(uint32(m.r1))
		if addr&0x20 != 0 {
			m.mem.setVramMirror(memVramMirrorH)
		} else {
			m.mem.setVramMirror(memVramMirrorV)
		}
	}
}

func (m *mapper29) write(addr uint16, data byte) {
	if m.r0 != 0 {
		m.r1 = (m.r1 & 0x0c) | (byte(addr) & 0x03)
		m.mem.setVrom8kBank(uint32(m.r1))
	}
}

// 0x2a

type mapper2a struct {
	baseMapper
}

func newMapper2a(bm *baseMapper) Mapper {
	return &mapper2a{baseMapper: *bm}
}

func (m *mapper2a) reset() {
}

// 0x2b

type mapper2b struct {
	baseMapper
}

func newMapper2b(bm *baseMapper) Mapper {
	return &mapper2b{baseMapper: *bm}
}

func (m *mapper2b) reset() {
}

// 0x2c

type mapper2c struct {
	baseMapper
	bank     byte
	p0, p1   byte
	c, r     [8]byte
	irqEn    bool
	irqCnt   byte
	irqLatch byte
}

func newMapper2c(bm *baseMapper) Mapper {
	return &mapper2c{baseMapper: *bm}
}

func (m *mapper2c) setPromBank() {
	ps := [4]byte{}
	if m.r[0]&0x40 != 0 {
		ps[0], ps[1], ps[2], ps[3] = 0x1e, 0x1f&m.p1, 0x1f&m.p0, 0x1f
	} else {
		ps[0], ps[1], ps[2], ps[3] = 0x1f&m.p0, 0x1f&m.p1, 0x1e, 0x1f
	}
	for i := byte(0); i < 4; i++ {
		p := ps[i]
		if m.bank != 6 {
			p &= 0x0f
		}
		m.mem.setProm8kBank(i+4, uint32(p)|(uint32(m.bank)<<4))
	}
}

func (m *mapper2c) setVromBank() {
	if m.mem.nVrom1kPage == 0 {
		return
	}
	br := m.r[0]&0x80 != 0
	for i := byte(0); i < 8; i++ {
		j := i
		if br {
			j ^= 0x04
		}
		c := m.c[j]
		if m.bank != 6 {
			c &= 0x7f
		}
		m.mem.setVrom1kBank(i, uint32(c)|(uint32(m.bank)<<7))
	}
}

func (m *mapper2c) reset() {
	m.p1 = 1
	if m.mem.nVrom1kPage != 0 {
		for i := byte(0); i < 8; i++ {
			m.c[i] = i
		}
	} else {
		m.c[1], m.c[3] = 1, 1
	}
	m.setPromBank()
	m.setVromBank()
}

func (m *mapper2c) writeLow(addr uint16, data byte) {
	if addr == 0x6000 {
		m.bank = (data & 0x01) << 1
		m.setPromBank()
		m.setVromBank()
	}
}

func (m *mapper2c) write(addr uint16, data byte) {
	switch addr & 0xe001 {
	case 0x8000:
		m.r[0] = data
		m.setPromBank()
		m.setVromBank()
	case 0x8001:
		m.r[1] = data
		r := m.r[0] & 0x07
		switch r {
		case 0x00, 0x01:
			r <<= 1
			m.c[r] = data & 0xfe
			m.c[r+1] = m.c[r] + 1
			m.setVromBank()
		case 0x02, 0x03, 0x04, 0x05:
			m.c[r+2] = data
			m.setVromBank()
		case 0x06:
			m.p0 = data
			m.setPromBank()
		case 0x07:
			m.p1 = data
			m.setPromBank()
		}
	case 0xa000:
		m.r[2] = data
		if m.sys.rom.b4Screen {
			if data&0x01 != 0 {
				m.mem.setVramMirror(memVramMirrorH)
			} else {
				m.mem.setVramMirror(memVramMirrorV)
			}
		}
	case 0xa001:
		m.r[3], m.bank = data, data&0x07
		if m.bank == 7 {
			m.bank = 6
		}
		m.setPromBank()
		m.setVromBank()
	case 0xc000:
		m.r[4], m.irqCnt = data, data
	case 0xc001:
		m.r[5], m.irqLatch = data, data
	case 0xe000:
		m.r[6], m.irqEn = data, false
		m.sys.cpu.intr &^= cpuIntrTypMapper
	case 0xe001:
		m.r[7], m.irqEn = data, true
	}
}

func (m *mapper2c) hSync(scanline uint16) {
	if scanline < ScreenHeight && (m.sys.ppu.reg1&(ppuReg1BgDisp|ppuReg1SpDisp) != 0) && m.irqEn {
		m.irqCnt--
		if m.irqCnt == 0 {
			m.irqCnt = m.irqLatch
			m.sys.cpu.intr &= cpuIntrTypMapper
		}
	}
}

// 0x2d

type mapper2d struct {
	baseMapper
	p, p1      [4]byte
	c, r       [8]byte
	c1         [8]uint32
	irqEn      bool
	ireReset   bool
	ireLatched bool
	irqCnt     byte
	irqLatch   byte
}

func newMapper2d(bm *baseMapper) Mapper {
	return &mapper2d{baseMapper: *bm}
}

func (m *mapper2d) setPromBank(i, data byte) {
	data = (data & ((m.r[3] & 0x3f) ^ 0xff) & 0x3f) | m.r[1]
	m.mem.setProm8kBank(i+4, uint32(data))
	m.p1[i] = data
}

func (m *mapper2d) setVromBank() {
	table := [16]uint32{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x01, 0x03, 0x07, 0x0f, 0x1f, 0x3f, 0x7f, 0xff,
	}
	r0, r2 := uint32(m.r[0]), m.r[2]
	for i := byte(0); i < 8; i++ {
		t := (uint32(m.c[i]) & table[r2&0x0f]) | r0
		m.c1[i] = t + (uint32(r2&0x10) << 4)
	}
	br := m.r[6]&0x80 != 0
	for i := byte(0); i < 8; i++ {
		j := i
		if br {
			j ^= 0x04
		}
		m.mem.setVrom1kBank(i, m.c1[j])
	}
}

func (m *mapper2d) reset() {
	m.p[1], m.p[2], m.p[3] = 1, byte(m.mem.nProm8kPage)-2, byte(m.mem.nProm8kPage)-1
	for i := byte(0); i < 4; i++ {
		m.mem.setProm8kBank(i, uint32(m.p[i]))
		m.p1[i] = m.p[i]
	}
	m.mem.setVrom8kBank(0)
	for i := byte(0); i < 8; i++ {
		m.c[i], m.c1[i] = i, uint32(i)
	}
}

func (m *mapper2d) writeLow(addr uint16, data byte) {
	if m.r[3]&0x40 == 0 {
		m.r[m.r[5]] = data
		m.r[5] = (m.r[5] + 1) & 0x03
		for i := byte(0); i < 4; i++ {
			m.setPromBank(i, m.p[i])
		}
		m.setVromBank()
	}
}

func (m *mapper2d) write(addr uint16, data byte) {
	switch addr & 0xe001 {
	case 0x8000:
		if data&0x40 != m.r[6]&0x40 {
			m.p[0], m.p[2] = m.p[2], m.p[0]
			m.p1[0], m.p1[2] = m.p1[2], m.p1[0]
			m.setPromBank(0, m.p1[0])
			m.setPromBank(1, m.p1[1])
		}
		if m.mem.nVrom1kPage != 0 && data&0x80 != m.r[6]&0x80 {
			for i := byte(0); i < 4; i++ {
				m.c[i], m.c[i+4] = m.c[i+4], m.c[i]
				m.c1[i], m.c1[i+4] = m.c1[i+4], m.c1[i]
				m.mem.setVrom1kBank(i, m.c1[i])
				m.mem.setVrom1kBank(i+4, m.c1[i+4])
			}
		}
		m.r[6] = data
	case 0x8001:
		r := m.r[6] & 0x07
		switch r {
		case 0x00, 0x01:
			r <<= 1
			m.c[r], m.c[r+1] = data&0xfe, (data&0xfe)+1
			m.setVromBank()
		case 0x02, 0x03, 0x04, 0x05:
			m.c[r+2] = data
			m.setVromBank()
		case 0x06:
			var i byte
			if m.r[6]&0x40 != 0 {
				i += 2
			}
			m.p[i] = data & 0x3f
			m.setPromBank(i, data)
		case 0x07:
			m.p[1] = data & 0x3f
			m.setPromBank(1, data)
		}
	case 0xa000:
		if data&0x01 != 0 {
			m.mem.setVramMirror(memVramMirrorH)
		} else {
			m.mem.setVramMirror(memVramMirrorV)
		}
	case 0xc000:
		m.irqLatch, m.ireLatched = data, true
		if m.ireReset {
			m.irqCnt, m.ireLatched = data, false
		}
	case 0xC001:
		m.irqCnt = m.irqLatch
	case 0xE000:
		m.irqEn, m.ireReset = false, true
		m.sys.cpu.intr &^= cpuIntrTypMapper
	case 0xE001:
		m.irqEn = true
		if m.ireLatched {
			m.irqCnt = m.irqLatch
		}
	}
}

func (m *mapper2d) hSync(scanline uint16) {
	m.ireReset = false
	if scanline < ScreenHeight && (m.sys.ppu.reg1&(ppuReg1BgDisp|ppuReg1SpDisp) != 0) && m.irqCnt != 0 {
		m.irqCnt--
		if m.irqCnt == 0 && m.irqEn {
			m.sys.cpu.intr &= cpuIntrTypMapper
		}
	}
}

// 0x2e

type mapper2e struct {
	baseMapper
	r0, r1, r2, r3 byte
}

func newMapper2e(bm *baseMapper) Mapper {
	return &mapper2e{baseMapper: *bm}
}

func (m *mapper2e) reset() {
	m.setBank()
	m.mem.setVramMirror(memVramMirrorV)
}

func (m *mapper2e) writeLow(addr uint16, data byte) {
	m.r0, m.r1 = data&0x0f, (data&0xf0)>>4
	m.setBank()
}

func (m *mapper2e) write(addr uint16, data byte) {
	m.r2, m.r3 = data&0x01, (data&0x70)>>4
	m.setBank()
}

func (m *mapper2e) setBank() {
	m.mem.setProm32kBank((uint32(m.r0) << 1) + uint32(m.r2))
	m.mem.setVrom8kBank((uint32(m.r1) << 3) + uint32(m.r3))
}

// 0x2f

type mapper2f struct {
	baseMapper
	bank     byte
	p0, p1   byte
	c, r     [8]byte
	irqEn    bool
	irqCnt   byte
	irqLatch byte
}

func newMapper2f(bm *baseMapper) Mapper {
	return &mapper2f{baseMapper: *bm}
}

func (m *mapper2f) setPromBank() {
	ps := [4]byte{}
	if m.r[0]&0x40 != 0 {
		ps[0], ps[1], ps[2], ps[3] = 0x0e, m.p1, m.p0, 0x0f
	} else {
		ps[0], ps[1], ps[2], ps[3] = m.p0, m.p1, 0x0e, 0x0f
	}
	b := uint32(m.bank) << 3
	for i := byte(0); i < 4; i++ {
		m.mem.setProm8kBank(i+4, b+uint32(ps[i]))
	}
}

func (m *mapper2f) setVromBank() {
	if m.mem.nVrom1kPage == 0 {
		return
	}
	br := m.r[0]&0x80 != 0
	b := uint32(m.bank&0x02) << 6
	for i := byte(0); i < 8; i++ {
		j := i
		if br {
			j ^= 0x04
		}
		m.mem.setVrom1kBank(i, b+uint32(m.c[j]))
	}
}

func (m *mapper2f) reset() {
	m.p1 = 1
	if m.mem.nVrom1kPage != 0 {
		for i := byte(0); i < 8; i++ {
			m.c[i] = i
		}
	} else {
		m.c[1], m.c[3] = 1, 1
	}
	m.setPromBank()
	m.setVromBank()
}

func (m *mapper2f) writeLow(addr uint16, data byte) {
	if addr == 0x6000 {
		m.bank = (data & 0x01) << 1
		m.setPromBank()
		m.setVromBank()
	}
}

func (m *mapper2f) write(addr uint16, data byte) {
	switch addr & 0xe001 {
	case 0x8000:
		m.r[0] = data
		m.setPromBank()
		m.setVromBank()
	case 0x8001:
		m.r[1] = data
		r := m.r[0] & 0x07
		switch r {
		case 0x00, 0x01:
			r <<= 1
			m.c[r] = data & 0xfe
			m.c[r+1] = m.c[r] + 1
			m.setVromBank()
		case 0x02, 0x03, 0x04, 0x05:
			m.c[r+2] = data
			m.setVromBank()
		case 0x06:
			m.p0 = data
			m.setPromBank()
		case 0x07:
			m.p1 = data
			m.setPromBank()
		}
	case 0xa000:
		m.r[2] = data
		if data&0x01 != 0 {
			m.mem.setVramMirror(memVramMirrorH)
		} else {
			m.mem.setVramMirror(memVramMirrorV)
		}
	case 0xa001:
		m.r[3] = data
	case 0xc000:
		m.r[4], m.irqCnt = data, data
	case 0xc001:
		m.r[5], m.irqLatch = data, data
	case 0xe000:
		m.r[6], m.irqEn = data, false
		m.sys.cpu.intr &^= cpuIntrTypMapper
	case 0xe001:
		m.r[7], m.irqEn = data, true
	}
}

func (m *mapper2f) hSync(scanline uint16) {
	if scanline < ScreenHeight && (m.sys.ppu.reg1&(ppuReg1BgDisp|ppuReg1SpDisp) != 0) && m.irqEn {
		m.irqCnt--
		if m.irqCnt == 0 {
			m.irqCnt = m.irqLatch
			m.sys.cpu.intr &= cpuIntrTypMapper
		}
	}
}
