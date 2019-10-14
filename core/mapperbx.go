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
	if m.nProm8kPage > 64 {
		m.mem.setProm32kBank4(0, 1, m.nProm8kPage/2-2, m.nProm8kPage/2-1)
	} else {
		m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
	}
	if m.nVrom1kPage != 0 {
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

// 180

type mapper180 struct {
	baseMapper
}

func newMapper180(bm *baseMapper) Mapper {
	return &mapper180{baseMapper: *bm}
}

func (m *mapper180) reset() {
	m.mem.setProm32kBank(0)
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper180) write(addr uint16, data byte) {
	m.mem.setProm16kBank(6, uint32(data&0x07))
}

// 181

type mapper181 struct {
	baseMapper
}

func newMapper181(bm *baseMapper) Mapper {
	return &mapper181{baseMapper: *bm}
}

func (m *mapper181) reset() {
	m.mem.setProm32kBank(0)
	m.mem.setVrom8kBank(0)
}

func (m *mapper181) writeLow(addr uint16, data byte) {
	if addr == 0x4120 {
		m.mem.setProm32kBank(uint32(data&0x08) >> 3)
		m.mem.setVrom8kBank(uint32(data & 0x07))
	}
}

// 182

type mapper182 struct {
	baseMapper
	irqEn  bool
	irqCnt byte
	r      byte
}

func newMapper182(bm *baseMapper) Mapper {
	return &mapper182{baseMapper: *bm}
}

func (m *mapper182) reset() {
	m.irqEn, m.irqCnt, m.r = false, 0, 0
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper182) write(addr uint16, data byte) {
	switch addr & 0xf003 {
	case 0x8001:
		if data&0x01 != 0 {
			m.mem.setVramMirror(memVramMirrorH)
		} else {
			m.mem.setVramMirror(memVramMirrorV)
		}
	case 0xa000:
		m.r = data & 0x07
	case 0xC000:
		switch m.r {
		case 0x00, 0x02:
			m.mem.setVrom1kBank(m.r, uint32(data&0xfe))
			m.mem.setVrom1kBank(m.r+1, uint32(data&0xfe)|0x01)
		case 0x01, 0x03:
			m.mem.setVrom1kBank(m.r+4, uint32(data))
		case 0x04, 0x05:
			m.mem.setProm8kBank(m.r, uint32(data))
		case 0x06, 0x07:
			m.mem.setVrom1kBank((m.r-4)<<1, uint32(data))
		}
	case 0xe003:
		m.irqEn, m.irqCnt = data != 0, data
		m.clearIntr()
	}
}

func (m *mapper182) hSync(scanline uint16) {
	if scanline < ScreenHeight && m.isPpuDisp() && m.irqEn {
		m.irqCnt--
		if m.irqCnt == 0 {
			m.irqEn = false
			m.setIntr()
		}
	}
}

// 183

type mapper183 struct {
	baseMapper
	irqEn  bool
	irqCnt uint16
	r      [8]byte
}

func newMapper183(bm *baseMapper) Mapper {
	return &mapper183{baseMapper: *bm}
}

func (m *mapper183) reset() {
	m.irqEn, m.irqCnt = false, 0
	for i := byte(0); i < 8; i++ {
		m.r[i] = i
	}
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
	m.mem.setVrom8kBank(0)
}

func (m *mapper183) write(addr uint16, data byte) {
	switch addr {
	case 0x8800:
		m.mem.setProm8kBank(4, uint32(data))
	case 0x9008:
		if data == 0x01 {
			for i := byte(0); i < 8; i++ {
				m.r[i] = i
			}
			m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
			m.mem.setVrom8kBank(0)
		}
	case 0x9800:
		switch data {
		case 0x00:
			m.mem.setVramMirror(memVramMirrorV)
		case 0x01:
			m.mem.setVramMirror(memVramMirrorH)
		case 0x02:
			m.mem.setVramMirror(memVramMirror4L)
		case 0x03:
			m.mem.setVramMirror(memVramMirror4H)
		}
	case 0xa000:
		m.mem.setProm8kBank(6, uint32(data))
	case 0xa800:
		m.mem.setProm8kBank(5, uint32(data))
	case 0xb000, 0xb008, 0xc000, 0xc008, 0xd000, 0xd008, 0xe000, 0xe008:
		i := (byte(addr>>11) - 0x16) | (byte(addr&0x08) >> 3)
		m.r[i] = (m.r[i] & 0xf0) | (data & 0x0f)
		m.mem.setVrom1kBank(i, uint32(m.r[i]))
	case 0xb004, 0xb00c, 0xc004, 0xc00c, 0xd004, 0xd00c, 0xe004, 0xe00c:
		i := (byte(addr>>11) - 0x16) | (byte(addr&0x08) >> 3)
		m.r[i] = (m.r[i] & 0x0f) | (data << 4)
		m.mem.setVrom1kBank(i, uint32(m.r[i]))
	case 0xf000:
		m.irqCnt = (m.irqCnt & 0xff00) | uint16(data)
	case 0xf004:
		m.irqCnt = (m.irqCnt & 0x00ff) | (uint16(data) << 8)
	case 0xf008:
		m.irqEn = data&0x02 != 0
		m.clearIntr()
	}
}

func (m *mapper183) hSync(scanline uint16) {
	if m.irqEn {
		if m.irqCnt <= 113 {
			m.irqCnt = 0
			m.setIntr()
		} else {
			m.irqCnt -= 113
		}
	}
}

// 184

type mapper184 struct {
	baseMapper
}

func newMapper184(bm *baseMapper) Mapper {
	return &mapper184{baseMapper: *bm}
}

func (m *mapper184) reset() {
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper184) writeLow(addr uint16, data byte) {
	if addr >= 0x6000 {
		m.mem.setVrom4kBank(0, uint32(((data&0x02)<<2)|(data&0x04)))
		m.mem.setVrom4kBank(1, uint32((data&0x20)>>2))
	}
}

// 185

type mapper185 struct {
	baseMapper
}

func newMapper185(bm *baseMapper) Mapper {
	return &mapper185{baseMapper: *bm}
}

func (m *mapper185) reset() {
	switch m.nProm8kPage >> 1 {
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

func (m *mapper185) write(addr uint16, data byte) {
	if data&0x03 != 0 {
		m.mem.setVrom8kBank(0)
	} else {
		for i := byte(0); i < 8; i++ {
			m.mem.setVram1kBank(i, 2)
		}
	}
}

// 187

type mapper187 struct {
	baseMapper
	extMode   bool // (ext_mode&0xA0)!=0xA0
	extEn     bool
	chrMode   byte
	irqEn     bool
	irqCnt    byte
	writePrev byte
	bank      [2]byte
	p         [4]byte
	c         [8]uint32
}

func newMapper187(bm *baseMapper) Mapper {
	return &mapper187{baseMapper: *bm}
}

func (m *mapper187) reset() {
	m.extMode, m.extEn, m.chrMode = true, false, 0
	m.irqEn, m.irqCnt = false, 0
	m.writePrev = 0
	m.bank[0], m.bank[1] = 0, 0
	b := byte(m.nProm8kPage)
	m.p[0], m.p[1], m.p[2], m.p[3] = b-4, b-3, b-2, b-1
	for i := 0; i < 8; i++ {
		m.c[i] = 0
	}
	m.sys.renderMode = RenderModePostAll
	m.mem.setProm32kBank4(uint32(m.p[0]), uint32(m.p[1]), uint32(m.p[2]), uint32(m.p[3]))
	m.mem.setVrom8kBank8(m.c[0], m.c[1], m.c[2], m.c[3], m.c[4], m.c[5], m.c[6], m.c[7])
}

func (m *mapper187) readLow(addr uint16) byte {
	switch m.writePrev {
	case 0x00:
		return 0x83
	case 0x01:
		return 0x83
	case 0x02:
		return 0x42
	case 0x03:
		return 0x00
	}
	return 0
}

func (m *mapper187) writeLow(addr uint16, data byte) {
	m.writePrev = data & 0x03
	if addr == 0x5000 {
		m.extMode = data&0xa0 != 0xa0
		if data&0x80 != 0 {
			if data&0x20 != 0 {
				m.p[0] = (data & 0x1e) << 1
				m.p[1] = ((data & 0x1e) << 1) | 0x01
				m.p[2] = ((data & 0x1e) << 1) | 0x02
				m.p[3] = ((data & 0x1e) << 1) | 0x03
			} else {
				m.p[2] = (data & 0x1f) << 1
				m.p[3] = ((data & 0x1f) << 1) | 0x01
			}
		} else {
			m.p[0] = m.bank[0]
			m.p[1] = m.bank[1]
			m.p[2] = byte(m.nProm8kPage) - 2
			m.p[3] = byte(m.nProm8kPage) - 1
		}
		m.mem.setProm32kBank4(uint32(m.p[0]), uint32(m.p[1]), uint32(m.p[2]), uint32(m.p[3]))
	}
}

func (m *mapper187) write(addr uint16, data byte) {
	m.writePrev = data & 0x03
	switch addr {
	case 0x8000:
		m.extEn, m.chrMode = false, data
	case 0x8001:
		c := m.chrMode & 0x07
		if !m.extEn {
			switch c {
			case 0x00, 0x01:
				i := 4 + (c << 1)
				m.c[i] = uint32(data&0xfe) | 0x0100
				m.c[i+1] = m.c[i] + 1
				m.mem.setVrom8kBank8(m.c[0], m.c[1], m.c[2], m.c[3], m.c[4], m.c[5], m.c[6], m.c[7])
			case 0x02, 0x03, 0x04, 0x05:
				m.c[c-2] = uint32(data)
				m.mem.setVrom8kBank8(m.c[0], m.c[1], m.c[2], m.c[3], m.c[4], m.c[5], m.c[6], m.c[7])
			case 0x06, 0x07:
				if m.extMode {
					m.p[c&0x01] = data
					m.mem.setProm32kBank4(uint32(m.p[0]), uint32(m.p[1]), uint32(m.p[2]), uint32(m.p[3]))
				}
			}
		} else {
			switch m.chrMode {
			case 0x2a:
				m.p[1] = 0x0f
			case 0x28:
				m.p[2] = 0x17
			}
			m.mem.setProm32kBank4(uint32(m.p[0]), uint32(m.p[1]), uint32(m.p[2]), uint32(m.p[3]))
		}
		if c >= 0x06 {
			m.bank[c-6] = data
		}
	case 0x8003:
		m.extEn, m.chrMode = true, data
		if data&0xf0 == 0 {
			m.p[2] = byte(m.nProm8kPage) - 2
			m.mem.setProm32kBank4(uint32(m.p[0]), uint32(m.p[1]), uint32(m.p[2]), uint32(m.p[3]))
		}
	case 0xa000:
		if data&0x01 != 0 {
			m.mem.setVramMirror(memVramMirrorH)
		} else {
			m.mem.setVramMirror(memVramMirrorV)
		}
	case 0xc000:
		m.irqCnt = data
		m.clearIntr()
	case 0xc001:
		m.clearIntr()
	case 0xe000, 0xe002:
		m.irqEn = false
		m.clearIntr()
	case 0xe001, 0xe003:
		m.irqEn = true
		m.clearIntr()
	}
}

func (m *mapper187) hSync(scanline uint16) {
	if scanline < ScreenHeight && m.isPpuDisp() && m.irqEn {
		if m.irqCnt == 0 {
			m.irqEn = true
			m.setIntr()
		}
		m.irqCnt--
	}
}

// 188

type mapper188 struct {
	baseMapper
}

func newMapper188(bm *baseMapper) Mapper {
	return &mapper188{baseMapper: *bm}
}

func (m *mapper188) reset() {
	if m.nProm8kPage > 16 {
		m.mem.setProm32kBank4(0, 1, 14, 15)
	} else {
		m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
	}
}

func (m *mapper188) write(addr uint16, data byte) {
	if data&0x10 != 0 {
		m.mem.setProm16kBank(4, uint32(data&0x07))
	} else if data != 0 {
		m.mem.setProm16kBank(4, uint32(data)+0x08)
	} else if m.nProm8kPage == 16 {
		m.mem.setProm16kBank(4, 7)
	} else {
		m.mem.setProm16kBank(4, 8)
	}
}

// 189

type mapper189 struct {
	baseMapper
	irqEn    bool
	irqCnt   byte
	irqLatch byte
	r        byte
	c        [8]byte
}

func newMapper189(bm *baseMapper) Mapper {
	return &mapper189{baseMapper: *bm}
}

func (m *mapper189) setPpuBanks() {
	if m.nVrom1kPage != 0 {
		if m.r&0x80 != 0 {
			for i := byte(0); i < 4; i++ {
				m.mem.setVrom1kBank(i, uint32(m.c[i+4]))
				m.mem.setVrom1kBank(i+4, uint32(m.c[i]))
			}
		} else {
			for i := byte(0); i < 8; i++ {
				m.mem.setVrom1kBank(i, uint32(m.c[i]))
			}
		}
	} else {
		if m.r&0x80 != 0 {
			for i := byte(0); i < 4; i++ {
				m.mem.setCram1kBank(i, uint32(m.c[i+4]&0x07))
				m.mem.setCram1kBank(i+4, uint32(m.c[i]&0x07))
			}
		} else {
			for i := byte(0); i < 8; i++ {
				m.mem.setCram1kBank(i, uint32(m.c[i]&0x07))
			}
		}
	}
}

func (m *mapper189) reset() {
	m.r = 0
	for i := byte(0); i < 8; i++ {
		m.c[i] = i
	}
	b := m.nProm8kPage
	m.mem.setProm32kBank4(b-4, b-3, b-2, b-1)
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
	m.setPpuBanks()
}

func (m *mapper189) writeLow(addr uint16, data byte) {
	switch addr & 0xff00 {
	case 0x4100:
		m.mem.setProm32kBank(uint32(data&0x30) >> 4)
	case 0x6100:
		m.mem.setProm32kBank(uint32(data & 0x03))
	}

}

func (m *mapper189) write(addr uint16, data byte) {
	switch addr & 0xe001 {
	case 0x8000:
		m.r = data
		m.setPpuBanks()
	case 0x8001:
		m.setPpuBanks()
		r := m.r & 0x07
		switch r {
		case 0x00, 0x01:
			i := r << 1
			m.c[i] = data & 0xfe
			m.c[i+1] = m.c[i] + 1
			m.setPpuBanks()
		case 0x02, 0x03, 0x04, 0x05:
			m.c[r+2] = data
			m.setPpuBanks()
		}
	case 0xa000:
		if data&0x01 != 0 {
			m.mem.setVramMirror(memVramMirrorH)
		} else {
			m.mem.setVramMirror(memVramMirrorV)
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

func (m *mapper189) hSync(scanline uint16) {
	if scanline < ScreenHeight && m.isPpuDisp() && m.irqEn {
		m.irqCnt--
		if m.irqCnt == 0 {
			m.irqCnt = m.irqLatch
			m.setIntr()
		}
	}
}

// 190

type mapper190 struct {
	baseMapper
	irqEn     bool
	irqCnt    byte
	irqLatch  byte
	cmdCh     bool
	cmd, cmdL byte
	lo        byte
}

func newMapper190(bm *baseMapper) Mapper {
	return &mapper190{baseMapper: *bm}
}

func (m *mapper190) reset() {
	m.irqEn, m.irqCnt = false, 0
	m.cmdCh, m.cmdL = false, 1
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
}

func (m *mapper190) readLow(addr uint16) byte {
	if addr == 0x5000 {
		return m.lo
	}
	return m.cpuBanks[addr>>13][addr&0x1fff]
}

func (m *mapper190) writeLow(addr uint16, data byte) {
	switch addr {
	case 0x5000:
		m.cmdL = data
		switch data {
		case 0xe0:
			m.mem.setProm32kBank(0)
		case 0xee:
			m.mem.setProm32kBank(3)
		}
	case 0x5001:
		if m.cmdL == 0 {
			m.mem.setProm32kBank(7)
		}
	case 0x5080:
		switch data {
		case 0x01:
			m.lo = 0x83
		case 0x02:
			m.lo = 0x42
		case 0x03:
			m.lo = 0x00
		}
	}
	m.cpuBanks[addr>>13][addr&0x1fff] = data
}

func (m *mapper190) write(addr uint16, data byte) {
	switch addr & 0xe003 {
	case 0x8000:
		m.cmd = data
	case 0x8003:
		m.cmdCh = data == 0x06
		switch data {
		case 0x06:
			m.mem.setProm32kBank4(30, 31, 31, 31)
		case 0x28:
			m.mem.setProm32kBank4(31, 31, 23, 31)
		case 0x2a:
			m.mem.setProm32kBank4(31, 15, 23, 31)
		}
	case 0x8001:
		if !m.cmdCh && m.cmdL != 0 {
			break
		}
		c := m.cmd & 0x07
		switch c {
		case 0x00, 0x01:
			i := c << 1
			if m.cmd&0x80 != 0 {
				i += 4
			}
			m.mem.setVrom1kBank(i, uint32(data)+256)
			m.mem.setVrom1kBank(i+1, uint32(data)+257)
		case 0x02, 0x03, 0x04, 0x05:
			i := c - 2
			if m.cmd&0x80 == 0 {
				i += 4
			}
			m.mem.setVrom1kBank(i, uint32(data))
		case 0x06:
			b := uint32(data) & ((m.nProm8kPage << 1) - 1)
			if m.cmdL&0x40 != 0 {
				m.mem.setProm8kBank(4, (m.nProm8kPage-1)<<1)
				m.mem.setProm8kBank(6, b)
			} else {
				m.mem.setProm8kBank(4, b)
				m.mem.setProm8kBank(6, (m.nProm8kPage-1)<<1)
			}
		case 0x07:
			b := uint32(data) & ((m.nProm8kPage << 1) - 1)
			if m.cmdL&0x40 != 0 {
				m.mem.setProm8kBank(4, (m.nProm8kPage-1)<<1)
				m.mem.setProm8kBank(5, b)
			} else {
				m.mem.setProm8kBank(5, b)
				m.mem.setProm8kBank(6, (m.nProm8kPage-1)<<1)
			}
		}
	case 0xa000:
		if data&0x01 != 0 {
			m.mem.setVramMirror(memVramMirrorH)
		} else {
			m.mem.setVramMirror(memVramMirrorV)
		}
	case 0xc000:
		m.irqCnt = data - 1
	case 0xc001:
		m.irqLatch = data - 1
	case 0xc002:
		m.irqCnt = data
	case 0xc003:
		m.irqLatch = data
	case 0xe000:
		m.irqEn, m.irqCnt = false, m.irqLatch
		m.clearIntr()
	case 0xe001:
		m.irqEn = true
	case 0xe002:
		m.irqEn, m.irqCnt = false, m.irqLatch
		m.clearIntr()
	case 0xe003:
		m.irqEn = true
	}
}

func (m *mapper190) hSync(scanline uint16) {
	if scanline < ScreenHeight && m.isPpuDisp() && m.irqEn {
		m.irqCnt--
		if m.irqCnt == 0xff {
			m.setIntr()
		}
	}
}

// 191

type mapper191 struct {
	baseMapper
	r, p, h byte
	c       [4]byte
}

func newMapper191(bm *baseMapper) Mapper {
	return &mapper191{baseMapper: *bm}
}

func (m *mapper191) setPpuBanks() {
	if m.nVrom1kPage != 0 {
		b := uint32(m.h) << 3
		for i := byte(0); i < 8; i++ {
			m.mem.setVrom1kBank(i, ((b+uint32(m.c[i>>1]))<<2)|uint32(i&0x03))
		}
	}
}

func (m *mapper191) reset() {
	m.r, m.p, m.h = 0, 0, 0
	m.c[0], m.c[1], m.c[2], m.c[3] = 0, 0, 0, 0
	m.mem.setProm32kBank(uint32(m.p))
	m.setPpuBanks()
}

func (m *mapper191) writeLow(addr uint16, data byte) {
	switch addr {
	case 0x4100:
		m.r = data
	case 0x4101:
		switch m.r {
		case 0x00, 0x01, 0x02, 0x03:
			m.c[m.r] = data & 0x07
			m.setPpuBanks()
		case 0x04:
			m.h = data & 0x07
			m.setPpuBanks()
		case 0x05:
			m.p = data & 0x07
			m.mem.setProm32kBank(uint32(m.p))
		case 0x07:
			if data&0x02 != 0 {
				m.mem.setVramMirror(memVramMirrorH)
			} else {
				m.mem.setVramMirror(memVramMirrorV)
			}
		}
	}
}
