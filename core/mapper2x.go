package core

// 032

type mapper032 struct {
	baseMapper
	patchTyp bool
	r        byte
}

func newMapper032(bm *baseMapper) Mapper {
	return &mapper032{baseMapper: *bm}
}

func (m *mapper032) reset() {
	patch := m.sys.conf.PatchTyp
	if patch&0x01 != 0 {
		m.patchTyp = true
	}
	m.r = 0
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
	if patch&0x02 != 0 {
		m.mem.setProm32kBank4(30, 31, 30, 31)
	}
}

func (m *mapper032) write(addr uint16, data byte) {
	switch addr & 0xf000 {
	case 0x8000:
		if m.r&0x02 == 0 {
			m.mem.setProm8kBank(4, uint32(data))
		} else {
			m.mem.setProm8kBank(6, uint32(data))
		}
	case 0x9000:
		m.r = data
		if data&0x01 != 0 {
			m.mem.setVramMirror(memVramMirrorH)
		} else {
			m.mem.setVramMirror(memVramMirrorV)
		}
	case 0xa000:
		m.mem.setProm8kBank(5, uint32(data))
	}
	switch addr & 0xf007 {
	case 0xb000, 0xb001, 0xb002, 0xb003, 0xb004, 0xb005:
		m.mem.setVrom1kBank(byte(addr&0x07), uint32(data))
	case 0xb006:
		m.mem.setVrom1kBank(6, uint32(data))
		if m.patchTyp && data&0x40 != 0 {
			m.mem.setVramBank(0, 0, 0, 1)
		}
	case 0xb007:
		m.mem.setVrom1kBank(7, uint32(data))
		if m.patchTyp && data&0x40 != 0 {
			m.mem.setVramMirror(memVramMirror4L)
		}
	}
}

// 033

type mapper033 struct {
	baseMapper
	patch    bool
	irqEn    bool
	irqCnt   byte
	irqLatch byte
	r        [7]byte
}

func newMapper033(bm *baseMapper) Mapper {
	return &mapper033{baseMapper: *bm}
}

func (m *mapper033) setPpuBanks() {
	m.mem.setVrom2kBank(0, uint32(m.r[0]))
	m.mem.setVrom2kBank(2, uint32(m.r[1]))
	for i := byte(4); i < 8; i++ {
		m.mem.setVrom1kBank(i, uint32(m.r[i-2]))
	}
}

func (m *mapper033) reset() {
	if m.sys.conf.PatchTyp&0x01 != 0 {
		m.patch = true
	}
	m.irqEn, m.irqCnt, m.irqLatch = false, 0, 0
	m.r[0], m.r[1], m.r[2], m.r[3], m.r[4], m.r[5], m.r[6] = 0, 2, 4, 5, 6, 7, 1
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	if m.mem.nVrom1kPage != 0 {
		m.setPpuBanks()
	}
}

func (m *mapper033) write(addr uint16, data byte) {
	switch addr {
	case 0x8000:
		if m.patch {
			if data&0x40 != 0 {
				m.mem.setVramMirror(memVramMirrorH)
			} else {
				m.mem.setVramMirror(memVramMirrorV)
			}
			m.mem.setProm8kBank(4, uint32(data&0x1f))
		} else {
			m.mem.setProm8kBank(4, uint32(data))
		}
	case 0x8001:
		if m.patch {
			m.mem.setProm8kBank(5, uint32(data&0x1f))
		} else {
			m.mem.setProm8kBank(5, uint32(data))
		}
	case 0x8002, 0x8003:
		m.r[addr-0x8002] = data
		m.setPpuBanks()
	case 0xa000, 0xa001, 0xa002, 0xa003:
		m.r[addr-0x9ffe] = data
		m.setPpuBanks()
	case 0xc000:
		m.irqCnt, m.irqLatch = data, data
	case 0xc001:
		m.irqCnt = m.irqLatch
	case 0xc002:
		m.irqEn = true
	case 0xc003:
		m.irqEn = false
	case 0xe000:
		if data&0x40 != 0 {
			m.mem.setVramMirror(memVramMirrorH)
		} else {
			m.mem.setVramMirror(memVramMirrorV)
		}
	}
}

func (m *mapper033) hSync(scanline uint16) {
	if scanline < ScreenHeight && m.isPpuDisp() && m.irqEn {
		m.irqCnt++
		if m.irqCnt == 0 {
			m.irqEn = false
			m.setIntr()
		}
	}
}

// 034

type mapper034 struct {
	baseMapper
}

func newMapper034(bm *baseMapper) Mapper {
	return &mapper034{baseMapper: *bm}
}

func (m *mapper034) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper034) writeLow(addr uint16, data byte) {
	switch addr {
	case 0x7ffd:
		m.mem.setProm32kBank(uint32(data))
	case 0x7ffe:
		m.mem.setVrom4kBank(0, uint32(data))
	case 0x7fff:
		m.mem.setVrom4kBank(4, uint32(data))
	}
}

func (m *mapper034) write(addr uint16, data byte) {
	m.mem.setProm32kBank(uint32(data))
}

// 040

type mapper040 struct {
	baseMapper
	irqEn   bool
	irqLine int32
}

func newMapper040(bm *baseMapper) Mapper {
	return &mapper040{baseMapper: *bm}
}

func (m *mapper040) reset() {
	m.mem.setProm8kBank(3, 6)
	m.mem.setProm32kBank4(4, 5, 0, 7)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper040) write(addr uint16, data byte) {
	switch addr & 0xe000 {
	case 0x8000:
		m.irqEn = false
		m.clearIntr()
	case 0xa000:
		m.irqEn, m.irqLine = true, 37
		m.clearIntr()
	case 0xe000:
		m.mem.setProm8kBank(6, uint32(data)&0x07)
	}
}

func (m *mapper040) hSync(scanline uint16) {
	if m.irqEn {
		m.irqLine--
		if m.irqLine <= 0 {
			m.setIntr()
		}
	}
}

// 041

type mapper041 struct {
	baseMapper
	r0, r1 byte
}

func newMapper041(bm *baseMapper) Mapper {
	return &mapper041{baseMapper: *bm}
}

func (m *mapper041) reset() {
	m.mem.setProm32kBank(0)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper041) writeLow(addr uint16, data byte) {
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

func (m *mapper041) write(addr uint16, data byte) {
	if m.r0 != 0 {
		m.r1 = (m.r1 & 0x0c) | (byte(addr) & 0x03)
		m.mem.setVrom8kBank(uint32(m.r1))
	}
}

// 042

type mapper042 struct {
	baseMapper
	irqEn  bool
	irqCnt byte
}

func newMapper042(bm *baseMapper) Mapper {
	return &mapper042{baseMapper: *bm}
}

func (m *mapper042) reset() {
	m.irqEn, m.irqCnt = false, 0
	m.mem.setProm8kBank(3, 0)
	b := m.mem.nProm8kPage
	m.mem.setProm32kBank4(b-4, b-3, b-2, b-1)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper042) write(addr uint16, data byte) {
	switch addr & 0xe003 {
	case 0xe000:
		m.mem.setProm8kBank(3, uint32(data&0x0f))
	case 0xe001:
		if data&0x08 != 0 {
			m.mem.setVramMirror(memVramMirrorH)
		} else {
			m.mem.setVramMirror(memVramMirrorV)
		}
	case 0xe002:
		if data&0x02 != 0 {
			m.irqEn = true
		} else {
			m.irqEn, m.irqCnt = false, 0
		}
		m.clearIntr()
	}
}

func (m *mapper042) hSync(scanline uint16) {
	m.clearIntr()
	if m.irqEn {
		if m.irqCnt < 215 {
			m.irqCnt++
		}
		if m.irqCnt == 215 {
			m.irqEn = false
			m.setIntr()
		}
	}
}

// 043

type mapper043 struct {
	baseMapper
	irqEn  bool
	irqCnt uint32
}

func newMapper043(bm *baseMapper) Mapper {
	return &mapper043{baseMapper: *bm}
}

func (m *mapper043) reset() {
	m.irqEn, m.irqCnt = true, 0
	m.mem.setProm8kBank(3, 2)
	m.mem.setProm32kBank4(1, 0, 4, 9)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper043) readLow(addr uint16) byte {
	if addr >= 0x5000 && addr < 0x6000 {
		return m.mem.prom[0xc000+addr]
	}
	return byte(addr >> 8)
}

func (m *mapper043) writeEx(addr uint16, data byte) {
	if addr&0xf0ff == 0x4022 {
		switch data & 0x07 {
		case 0x00, 0x02, 0x03, 0x04:
			m.mem.setProm8kBank(6, 4)
		case 0x01:
			m.mem.setProm8kBank(6, 3)
		case 0x05:
			m.mem.setProm8kBank(6, 7)
		case 0x06:
			m.mem.setProm8kBank(6, 5)
		case 0x07:
			m.mem.setProm8kBank(6, 6)
		}
	}
}

func (m *mapper043) writeLow(addr uint16, data byte) {
	m.writeEx(addr, data)
}

func (m *mapper043) write(addr uint16, data byte) {
	if addr == 0x8122 {
		if data&0x03 != 0 {
			m.irqEn = true
		} else {
			m.irqEn, m.irqCnt = false, 0
		}
	}
}

func (m *mapper043) hSync(scanline uint16) {
	if m.irqEn {
		m.irqCnt += 114
		if m.irqCnt >= 4096 {
			m.irqCnt -= 4096
			m.setIntr()
		}
	}
}

// 044

type mapper044 struct {
	baseMapper
	bank     byte
	p0, p1   byte
	c, r     [8]byte
	irqEn    bool
	irqCnt   byte
	irqLatch byte
}

func newMapper044(bm *baseMapper) Mapper {
	return &mapper044{baseMapper: *bm}
}

func (m *mapper044) setCpuBanks() {
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

func (m *mapper044) setPpuBanks() {
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

func (m *mapper044) reset() {
	m.p1 = 1
	if m.mem.nVrom1kPage != 0 {
		for i := byte(0); i < 8; i++ {
			m.c[i] = i
		}
	} else {
		for i := byte(0); i < 8; i++ {
			m.c[i] = 0
		}
		m.c[1], m.c[3] = 1, 1
	}
	m.setCpuBanks()
	m.setPpuBanks()
}

func (m *mapper044) writeLow(addr uint16, data byte) {
	if addr == 0x6000 {
		m.bank = (data & 0x01) << 1
		m.setCpuBanks()
		m.setPpuBanks()
	}
}

func (m *mapper044) write(addr uint16, data byte) {
	switch addr & 0xe001 {
	case 0x8000:
		m.r[0] = data
		m.setCpuBanks()
		m.setPpuBanks()
	case 0x8001:
		m.r[1] = data
		r := m.r[0] & 0x07
		switch r {
		case 0x00, 0x01:
			r <<= 1
			m.c[r] = data & 0xfe
			m.c[r+1] = m.c[r] + 1
			m.setPpuBanks()
		case 0x02, 0x03, 0x04, 0x05:
			m.c[r+2] = data
			m.setPpuBanks()
		case 0x06:
			m.p0 = data
			m.setCpuBanks()
		case 0x07:
			m.p1 = data
			m.setCpuBanks()
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
		m.setCpuBanks()
		m.setPpuBanks()
	case 0xc000:
		m.r[4], m.irqCnt = data, data
	case 0xc001:
		m.r[5], m.irqLatch = data, data
	case 0xe000:
		m.r[6], m.irqEn = data, false
		m.clearIntr()
	case 0xe001:
		m.r[7], m.irqEn = data, true
	}
}

func (m *mapper044) hSync(scanline uint16) {
	if scanline < ScreenHeight && m.isPpuDisp() && m.irqEn {
		m.irqCnt--
		if m.irqCnt == 0 {
			m.irqCnt = m.irqLatch
			m.sys.cpu.intr &= cpuIntrTypMapper
		}
	}
}

// 045

type mapper045 struct {
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

func newMapper045(bm *baseMapper) Mapper {
	return &mapper045{baseMapper: *bm}
}

func (m *mapper045) setCpuBanks(i, data byte) {
	data = (data & ((m.r[3] & 0x3f) ^ 0xff) & 0x3f) | m.r[1]
	m.mem.setProm8kBank(i+4, uint32(data))
	m.p1[i] = data
}

func (m *mapper045) setPpuBanks() {
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

func (m *mapper045) reset() {
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

func (m *mapper045) writeLow(addr uint16, data byte) {
	if m.r[3]&0x40 == 0 {
		m.r[m.r[5]] = data
		m.r[5] = (m.r[5] + 1) & 0x03
		for i := byte(0); i < 4; i++ {
			m.setCpuBanks(i, m.p[i])
		}
		m.setPpuBanks()
	}
}

func (m *mapper045) write(addr uint16, data byte) {
	switch addr & 0xe001 {
	case 0x8000:
		if data&0x40 != m.r[6]&0x40 {
			m.p[0], m.p[2] = m.p[2], m.p[0]
			m.p1[0], m.p1[2] = m.p1[2], m.p1[0]
			m.setCpuBanks(0, m.p1[0])
			m.setCpuBanks(1, m.p1[1])
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
			m.setPpuBanks()
		case 0x02, 0x03, 0x04, 0x05:
			m.c[r+2] = data
			m.setPpuBanks()
		case 0x06:
			var i byte
			if m.r[6]&0x40 != 0 {
				i += 2
			}
			m.p[i] = data & 0x3f
			m.setCpuBanks(i, data)
		case 0x07:
			m.p[1] = data & 0x3f
			m.setCpuBanks(1, data)
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
	case 0xc001:
		m.irqCnt = m.irqLatch
	case 0xe000:
		m.irqEn, m.ireReset = false, true
		m.clearIntr()
	case 0xe001:
		m.irqEn = true
		if m.ireLatched {
			m.irqCnt = m.irqLatch
		}
	}
}

func (m *mapper045) hSync(scanline uint16) {
	m.ireReset = false
	if scanline < ScreenHeight && m.isPpuDisp() && m.irqCnt != 0 {
		m.irqCnt--
		if m.irqCnt == 0 && m.irqEn {
			m.sys.cpu.intr &= cpuIntrTypMapper
		}
	}
}

// 046

type mapper046 struct {
	baseMapper
	r0, r1, r2, r3 byte
}

func newMapper046(bm *baseMapper) Mapper {
	return &mapper046{baseMapper: *bm}
}

func (m *mapper046) reset() {
	m.setBank()
	m.mem.setVramMirror(memVramMirrorV)
}

func (m *mapper046) writeLow(addr uint16, data byte) {
	m.r0, m.r1 = data&0x0f, (data&0xf0)>>4
	m.setBank()
}

func (m *mapper046) write(addr uint16, data byte) {
	m.r2, m.r3 = data&0x01, (data&0x70)>>4
	m.setBank()
}

func (m *mapper046) setBank() {
	m.mem.setProm32kBank((uint32(m.r0) << 1) + uint32(m.r2))
	m.mem.setVrom8kBank((uint32(m.r1) << 3) + uint32(m.r3))
}

// 047

type mapper047 struct {
	baseMapper
	bank     byte
	p0, p1   byte
	c, r     [8]byte
	irqEn    bool
	irqCnt   byte
	irqLatch byte
}

func newMapper047(bm *baseMapper) Mapper {
	return &mapper047{baseMapper: *bm}
}

func (m *mapper047) setCpuBanks() {
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

func (m *mapper047) setPpuBanks() {
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

func (m *mapper047) reset() {
	m.p1 = 1
	if m.mem.nVrom1kPage != 0 {
		for i := byte(0); i < 8; i++ {
			m.c[i] = i
		}
	} else {
		for i := byte(0); i < 8; i++ {
			m.c[i] = 0
		}
		m.c[1], m.c[3] = 1, 1
	}
	m.setCpuBanks()
	m.setPpuBanks()
}

func (m *mapper047) writeLow(addr uint16, data byte) {
	if addr == 0x6000 {
		m.bank = (data & 0x01) << 1
		m.setCpuBanks()
		m.setPpuBanks()
	}
}

func (m *mapper047) write(addr uint16, data byte) {
	switch addr & 0xe001 {
	case 0x8000:
		m.r[0] = data
		m.setCpuBanks()
		m.setPpuBanks()
	case 0x8001:
		m.r[1] = data
		r := m.r[0] & 0x07
		switch r {
		case 0x00, 0x01:
			r <<= 1
			m.c[r] = data & 0xfe
			m.c[r+1] = m.c[r] + 1
			m.setPpuBanks()
		case 0x02, 0x03, 0x04, 0x05:
			m.c[r+2] = data
			m.setPpuBanks()
		case 0x06:
			m.p0 = data
			m.setCpuBanks()
		case 0x07:
			m.p1 = data
			m.setCpuBanks()
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
		m.clearIntr()
	case 0xe001:
		m.r[7], m.irqEn = data, true
	}
}

func (m *mapper047) hSync(scanline uint16) {
	if scanline < ScreenHeight && m.isPpuDisp() && m.irqEn {
		m.irqCnt--
		if m.irqCnt == 0 {
			m.irqCnt = m.irqLatch
			m.sys.cpu.intr &= cpuIntrTypMapper
		}
	}
}
