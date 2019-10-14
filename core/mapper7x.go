package core

// 112

type mapper112 struct {
	baseMapper
	r [4]byte
	p [2]byte
	c [8]byte
}

func (m *mapper112) setPpuBanks() {
	if m.r[2]&0x02 != 0 {
		for i := byte(0); i < 8; i++ {
			m.mem.setVrom1kBank(i, uint32(m.c[i]))
		}
	} else {
		r := uint32(m.r[3])
		for i := byte(0); i < 4; i++ {
			m.mem.setVrom1kBank(i, ((r<<(6-(i>>1)))&0x0100)|uint32(m.c[i]))
			m.mem.setVrom1kBank(i+4, ((r<<(4-i))&0x0100)|uint32(m.c[i+4]))
		}
	}
}

func newMapper112(bm *baseMapper) Mapper {
	return &mapper112{baseMapper: *bm}
}

func (m *mapper112) reset() {
	m.r[0], m.r[1], m.r[2], m.r[3] = 0, 0, 0, 0
	m.p[0], m.p[1] = 0, 1
	for i := byte(0); i < 8; i++ {
		m.c[i] = i
	}
	m.mem.setProm32kBank4(uint32(m.p[0]), uint32(m.p[1]), m.nProm8kPage-2, m.nProm8kPage-1)
	m.setPpuBanks()
}

func (m *mapper112) write(addr uint16, data byte) {
	switch addr {
	case 0x8000:
		m.r[0] = data
		m.mem.setProm32kBank4(uint32(m.p[0]), uint32(m.p[1]), m.nProm8kPage-2, m.nProm8kPage-1)
		m.setPpuBanks()
	case 0xa000:
		m.r[1] = data
		r := m.r[0] & 0x07
		switch r {
		case 0x00, 0x01:
			m.p[r] = data & (byte(m.nProm8kPage) - 1)
			m.mem.setProm32kBank4(uint32(m.p[0]), uint32(m.p[1]), m.nProm8kPage-2, m.nProm8kPage-1)
		case 0x02, 0x03:
			i := (r - 2) << 1
			m.c[i] = data & 0xfe
			m.c[i+1] = m.c[i]
			m.setPpuBanks()
		case 0x04, 0x05, 0x06, 0x07:
			m.c[r] = data
			m.setPpuBanks()
		}
	case 0xc000:
		m.r[3] = data
		m.setPpuBanks()
	case 0xe000:
		m.r[2] = data
		if !m.sys.rom.b4Screen {
			if data&0x01 != 0 {
				m.mem.setVramMirror(memVramMirrorH)
			} else {
				m.mem.setVramMirror(memVramMirrorV)
			}
		}
		m.setPpuBanks()
	}
}

// 113

type mapper113 struct {
	baseMapper
	patchTyp bool
}

func newMapper113(bm *baseMapper) Mapper {
	return &mapper113{baseMapper: *bm}
}

func (m *mapper113) reset() {
	if m.sys.conf.PatchTyp&0x01 != 0 {
		m.patchTyp = true
	}
	m.mem.setProm32kBank(0)
	m.mem.setVrom8kBank(0)
}

func (m *mapper113) writeLow(addr uint16, data byte) {
	switch addr {
	case 0x4100, 0x4111, 0x4120, 0x4194, 0x4195, 0x4900:
		if m.patchTyp {
			if data&0x80 != 0 {
				m.mem.setVramMirror(memVramMirrorV)
			} else {
				m.mem.setVramMirror(memVramMirrorH)
			}
		}
		m.mem.setProm32kBank(uint32(data >> 3))
		m.mem.setVrom8kBank((uint32(data>>3) & 0x08) + uint32(data&0x07))
	}
}

func (m *mapper113) write(addr uint16, data byte) {
	switch addr {
	case 0x8008, 0x8009:
		m.mem.setProm32kBank(uint32(data >> 3))
		m.mem.setVrom8kBank((uint32(data>>3) & 0x08) + uint32(data&0x07))
	case 0x8e66, 0x8e67:
		var b uint32
		if data&0x07 != 0 {
			b = 1
		}
		m.mem.setVrom8kBank(b)
	case 0xe00a:
		m.mem.setVramMirror(memVramMirror4L)
	}
}

// 114

type mapper114 struct {
	baseMapper
	irqOccur bool
	irqCnt   byte
	c        bool
	m, a     byte
	b        [8]byte
}

func newMapper114(bm *baseMapper) Mapper {
	return &mapper114{baseMapper: *bm}
}

func (m *mapper114) reset() {
	m.irqOccur, m.irqCnt = false, 0
	m.c, m.m, m.a = false, 0, 0
	for i := 0; i < 8; i++ {
		m.b[i] = 0
	}
	m.sys.renderMode = RenderModePost
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper114) setCpuBanks() {
	if m.m&0x80 != 0 {
		m.mem.setProm16kBank(4, uint32(m.m&0x1f))
	} else {
		m.mem.setProm8kBank(4, uint32(m.b[4]))
		m.mem.setProm8kBank(5, uint32(m.b[5]))
	}
}

func (m *mapper114) setPpuBanks() {
	m.mem.setVrom2kBank(0, uint32(m.b[0]>>1))
	m.mem.setVrom2kBank(2, uint32(m.b[2]>>1))
	m.mem.setVrom1kBank(4, uint32(m.b[6]))
	m.mem.setVrom1kBank(5, uint32(m.b[1]))
	m.mem.setVrom1kBank(6, uint32(m.b[7]))
	m.mem.setVrom1kBank(7, uint32(m.b[3]))
}

func (m *mapper114) writeLow(addr uint16, data byte) {
	m.m = data
	m.setPpuBanks()
}

func (m *mapper114) write(addr uint16, data byte) {
	switch addr {
	case 0xe002:
		m.irqOccur = false
		m.clearIntr()
	case 0xe003:
		m.irqCnt = data
	default:
		switch addr & 0xe000 {
		case 0x8000:
			if data&0x01 != 0 {
				m.mem.setVramMirror(memVramMirrorH)
			} else {
				m.mem.setVramMirror(memVramMirrorV)
			}
		case 0xa000:
			m.c, m.a = true, data
		case 0xc000:
			if m.c {
				i := m.a & 0x07
				m.c, m.b[i] = false, data
				switch i {
				case 0x00, 0x01, 0x02, 0x03, 0x06, 0x07:
					m.setPpuBanks()
				case 0x04, 0x05:
					m.setCpuBanks()
				}
			}
		}
	}
}

func (m *mapper114) hSync(scanline uint16) {
	if scanline < ScreenHeight && m.isPpuDisp() && m.irqCnt != 0 {
		m.irqCnt--
		if m.irqCnt == 0 {
			m.irqOccur = true
			m.setIntr()
		}
	}
}

// 115

type mapper115 struct {
	baseMapper
	cSwitch  bool
	pSwitch  byte
	irqEn    bool
	irqCnt   byte
	irqLatch byte
	r        byte
	p        [6]byte
	c        [8]byte
}

func newMapper115(bm *baseMapper) Mapper {
	return &mapper115{baseMapper: *bm}
}

func (m *mapper115) setCpuBanks() {
	if m.pSwitch&0x80 == 0 {
		m.p[0], m.p[1] = m.p[4], m.p[5]
		if m.r&0x40 != 0 {
			m.mem.setProm32kBank4(m.nProm8kPage-2, uint32(m.p[1]), uint32(m.p[0]), m.nProm8kPage-1)
		} else {
			m.mem.setProm32kBank4(uint32(m.p[0]), uint32(m.p[1]), m.nProm8kPage-2, m.nProm8kPage-1)
		}
	} else {
		m.p[0] = (m.pSwitch << 1) & 0x1e
		m.p[1] = m.p[0] + 1
		m.mem.setProm32kBank4(uint32(m.p[0]), uint32(m.p[1]), uint32(m.p[0])+2, uint32(m.p[1])+2)
	}
}

func (m *mapper115) setPpuBanks() {
	if m.nVrom1kPage != 0 {
		var b uint32
		if m.cSwitch {
			b = 0x0100
		}
		if m.r&0x80 != 0 {
			for i := byte(0); i < 4; i++ {
				m.mem.setVrom1kBank(i, uint32(m.c[i+4])|b)
				m.mem.setVrom1kBank(i+4, uint32(m.c[i])|b)
			}
		} else {
			for i := byte(0); i < 8; i++ {
				m.mem.setVrom1kBank(i, uint32(m.c[i])|b)
			}
		}
	}
}

func (m *mapper115) reset() {
	m.cSwitch, m.pSwitch, m.irqEn, m.irqCnt, m.irqLatch = false, 0, false, 0, 0
	m.r, m.p[0], m.p[1], m.p[2], m.p[3] = 0, 0, 1, byte(m.nProm8kPage)-2, byte(m.nProm8kPage)-1
	m.p[4], m.p[5] = 0, 1
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
	m.setCpuBanks()
}

func (m *mapper115) writeLow(addr uint16, data byte) {
	switch addr {
	case 0x6000:
		m.pSwitch = data
		m.setCpuBanks()
	case 0x6001:
		m.cSwitch = data&0x01 != 0
		m.setPpuBanks()
	}
	m.baseMapper.writeLow(addr, data)
}

func (m *mapper115) write(addr uint16, data byte) {
	switch addr & 0xe001 {
	case 0x8000:
		m.r = data
		m.setCpuBanks()
		m.setPpuBanks()
	case 0x8001:
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
		case 0x06, 0x07:
			m.p[r&0x01], m.p[r-2] = data, data
			m.setCpuBanks()
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

func (m *mapper115) hSync(scanline uint16) {
	if scanline < ScreenHeight && m.isPpuDisp() && m.irqEn {
		m.irqCnt--
		if m.irqCnt == 0xff {
			m.irqCnt = m.irqLatch
			m.setIntr()
		}
	}
}

// 116

type mapper116 struct {
	baseMapper
	cSwitch  bool
	irqEn    bool
	irqLatch byte
	irqCnt   int16
	r        byte
	p        [4]byte
	c        [8]byte
}

func newMapper116(bm *baseMapper) Mapper {
	return &mapper116{baseMapper: *bm}
}

func (m *mapper116) setCpuBanks() {
	if m.r&0x40 != 0 {
		m.mem.setProm32kBank4(m.nProm8kPage-2, uint32(m.p[1]), uint32(m.p[0]), m.nProm8kPage-1)
	} else {
		m.mem.setProm32kBank4(uint32(m.p[0]), uint32(m.p[1]), m.nProm8kPage-2, m.nProm8kPage-1)
	}
}

func (m *mapper116) setPpuBanks() {
	if m.nVrom1kPage != 0 {
		var b uint32
		if m.cSwitch {
			b = 0x0100
		}
		for i := byte(0); i < 4; i++ {
			m.mem.setVrom1kBank(i, uint32(m.c[i+4])|b)
			m.mem.setVrom1kBank(i+4, uint32(m.c[i]))
		}
	}
}

func (m *mapper116) reset() {
	m.cSwitch, m.irqEn, m.irqCnt, m.irqLatch = false, false, 0, 0
	m.r, m.p[0], m.p[1], m.p[2], m.p[3] = 0, 0, 1, byte(m.nProm8kPage)-2, byte(m.nProm8kPage)-1
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
	m.setCpuBanks()
}

func (m *mapper116) writeLow(addr uint16, data byte) {
	if addr&0x4100 == 0x4100 {
		m.cSwitch = data&0x04 != 0
		m.setPpuBanks()
	}
}

func (m *mapper116) write(addr uint16, data byte) {
	switch addr & 0xe001 {
	case 0x8000:
		m.r = data
		m.setCpuBanks()
		m.setPpuBanks()
	case 0x8001:
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
		case 0x06, 0x07:
			m.p[r&0x01] = data
			m.setCpuBanks()
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
		m.irqCnt = int16(data)
	case 0xc001:
		m.irqLatch = data
	case 0xe000:
		m.irqEn, m.irqCnt = false, int16(m.irqLatch)
		m.clearIntr()
	case 0xe001:
		m.irqEn = true
	}
}

func (m *mapper116) hSync(scanline uint16) {
	if scanline < ScreenHeight {
		if m.irqCnt <= 0 && m.irqEn {
			m.setIntr()
		} else if m.isPpuDisp() {
			m.irqCnt--
		}
	}
}

// 117

type mapper117 struct {
	baseMapper
	irqEn  bool
	irqCnt byte
}

func newMapper117(bm *baseMapper) Mapper {
	return &mapper117{baseMapper: *bm}
}

func (m *mapper117) reset() {
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper117) write(addr uint16, data byte) {
	switch addr {
	case 0x8000, 0x8001, 0x8002:
		m.mem.setProm8kBank(byte(addr)+4, uint32(data))
	case 0xa000, 0xa001, 0xa002, 0xa003, 0xa004, 0xa005, 0xa006, 0xa007:
		m.mem.setVrom1kBank(byte(addr), uint32(data))
	case 0xc001, 0xc002, 0xc003:
		m.irqCnt = data
	case 0xe000:
		m.irqEn = data&0x01 != 0
	}
}

func (m *mapper117) hSync(scanline uint16) {
	if scanline < ScreenHeight && m.isPpuDisp() && m.irqEn && m.irqCnt == byte(scanline) {
		m.irqEn = false
		m.sys.cpu.intr |= cpuIntrTypTrig
	}
}

// 118

type mapper118 struct {
	baseMapper
	irqEn    bool
	irqCnt   byte
	irqLatch byte
	r        byte
	p        [2]byte
	c        [8]byte
}

func newMapper118(bm *baseMapper) Mapper {
	return &mapper118{baseMapper: *bm}
}

func (m *mapper118) setCpuBanks() {
	if m.r&0x40 != 0 {
		m.mem.setProm32kBank4(m.nProm8kPage-2, uint32(m.p[1]), uint32(m.p[0]), m.nProm8kPage-1)
	} else {
		m.mem.setProm32kBank4(uint32(m.p[0]), uint32(m.p[1]), m.nProm8kPage-2, m.nProm8kPage-1)
	}
}

func (m *mapper118) setPpuBanks() {
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
	}
}

func (m *mapper118) reset() {
	m.irqEn, m.irqCnt, m.irqLatch = false, 0, 0
	m.r, m.p[0], m.p[1] = 0, 0, 1
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
	m.setCpuBanks()
}

func (m *mapper118) write(addr uint16, data byte) {
	switch addr & 0xe001 {
	case 0x8000:
		m.r = data
		m.setCpuBanks()
		m.setPpuBanks()
	case 0x8001:
		if (m.r&0x80 != 0 && m.r&0x07 == 0x02) || (m.r&0x80 == 0 && m.r&0x07 == 0) {
			if data&0x80 != 0 {
				m.mem.setVramMirror(memVramMirror4L)
			} else {
				m.mem.setVramMirror(memVramMirror4H)
			}
		}
		r := m.r & 0x07
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
		case 0x06, 0x07:
			m.p[r&0x01] = data
			m.setCpuBanks()
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

func (m *mapper118) hSync(scanline uint16) {
	if scanline < ScreenHeight && m.isPpuDisp() && m.irqEn {
		m.irqCnt--
		if m.irqCnt == 0xff {
			m.irqCnt = m.irqLatch
			m.setIntr()
		}
	}
}

// 119

type mapper119 struct {
	baseMapper
	irqEn    bool
	irqCnt   byte
	irqLatch byte
	r        byte
	p        [2]byte
	c        [8]byte
}

func newMapper119(bm *baseMapper) Mapper {
	return &mapper119{baseMapper: *bm}
}

func (m *mapper119) setCpuBanks() {
	if m.r&0x40 != 0 {
		m.mem.setProm32kBank4(m.nProm8kPage-2, uint32(m.p[1]), uint32(m.p[0]), m.nProm8kPage-1)
	} else {
		m.mem.setProm32kBank4(uint32(m.p[0]), uint32(m.p[1]), m.nProm8kPage-2, m.nProm8kPage-1)
	}
}

func (m *mapper119) setPpuBanks() {
	if m.r&0x80 != 0 {
		for i := byte(0); i < 4; i++ {
			if m.c[i+4]&0x40 != 0 {
				m.mem.setCram1kBank(i, uint32(m.c[i+4]&0x07))
			} else {
				m.mem.setVrom1kBank(i, uint32(m.c[i+4]))
			}
			if m.c[i]&0x40 != 0 {
				m.mem.setCram1kBank(i+4, uint32(m.c[i]&0x07))
			} else {
				m.mem.setVrom1kBank(i+4, uint32(m.c[i]))
			}
		}
	} else {
		for i := byte(0); i < 8; i++ {
			if m.c[i]&0x40 != 0 {
				m.mem.setCram1kBank(i, uint32(m.c[i]&0x07))
			} else {
				m.mem.setVrom1kBank(i, uint32(m.c[i]))
			}
		}
	}
}

func (m *mapper119) reset() {
	m.irqEn, m.irqCnt, m.irqLatch = false, 0, 0
	m.r, m.p[0], m.p[1] = 0, 0, 1
	for i := byte(0); i < 8; i++ {
		m.c[i] = i
	}
	m.setPpuBanks()
	m.setCpuBanks()
}

func (m *mapper119) write(addr uint16, data byte) {
	switch addr & 0xe001 {
	case 0x8000:
		m.r = data
		m.setCpuBanks()
		m.setPpuBanks()
	case 0x8001:
		r := m.r & 0x07
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
		case 0x06, 0x07:
			m.p[r&0x01] = data
			m.setCpuBanks()
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
		m.irqEn, m.irqCnt = false, m.irqLatch
		m.clearIntr()
	case 0xe001:
		m.irqEn = true
	}
}

func (m *mapper119) hSync(scanline uint16) {
	if scanline < ScreenHeight && m.isPpuDisp() && m.irqEn {
		m.irqCnt--
		if m.irqCnt == 0xff {
			m.irqCnt = m.irqLatch
			m.setIntr()
		}
	}
}

// 120

type mapper120 struct {
	baseMapper
}

func newMapper120(bm *baseMapper) Mapper {
	return &mapper120{baseMapper: *bm}
}

func (m *mapper120) reset() {
	m.mem.setProm8kBank(3, 8)
	m.mem.setProm32kBank(0)
}

func (m *mapper120) writeLow(addr uint16, data byte) {
	if addr == 0x41ff {
		m.mem.setProm8kBank(3, uint32(data)+8)
	} else {
		m.baseMapper.writeLow(addr, data)
	}
}

// 121

type mapper121 struct {
	baseMapper
	irqReq   bool
	irqRel   bool
	irqCnt   byte
	irqLatch byte
	cmd      byte
	a0, a1   byte
	r, d     [8]byte
}

func newMapper121(bm *baseMapper) Mapper {
	return &mapper121{baseMapper: *bm}
}

func (m *mapper121) setCpuBanks0(iBank byte, data byte) {
	m.mem.setProm8kBank(iBank, uint32(data&0x1f)|(uint32(m.r[3]&0x80)>>2))
	if m.r[5]&0x3f != 0 {
		r := uint32(m.r[3]&0x80) >> 2
		m.mem.setProm8kBank(5, uint32(m.r[2])|r)
		m.mem.setProm8kBank(6, uint32(m.r[1])|r)
		m.mem.setProm8kBank(7, uint32(m.r[0])|r)
	}
}

func (m *mapper121) setPpuBanks0(iBank byte, data byte) {
	if m.nProm8kPage == m.nVrom1kPage>>3 {
		m.mem.setVrom1kBank(iBank, uint32(data)|(uint32(m.r[3]&0x80)<<1))
	} else if iBank&0x04 == (m.cmd&0x80)>>5 {
		m.mem.setVrom1kBank(iBank, uint32(data)|0x0100)
	} else {
		m.mem.setVrom1kBank(iBank, uint32(data))
	}
}

func (m *mapper121) setPpuMirrors() {
	if m.a0&0x01 != 0 {
		m.mem.setVramMirror(memVramMirrorH)
	} else {
		m.mem.setVramMirror(memVramMirrorV)
	}
}

func (m *mapper121) setCpuBanks1(data byte) {
	if data&0x40 != 0 {
		m.setCpuBanks0(4, byte(m.nProm8kPage)-2)
		m.setCpuBanks0(6, m.d[6])
	} else {
		m.setCpuBanks0(4, m.d[6])
		m.setCpuBanks0(6, byte(m.nProm8kPage)-2)
	}
	m.setCpuBanks0(5, m.d[7])
	m.setCpuBanks0(7, byte(m.nProm8kPage)-1)
}

func (m *mapper121) setPpuBanks1(data byte) {
	b := (data & 0x80) >> 5
	m.setPpuBanks0(b^0x00, m.d[0]&^0x01)
	m.setPpuBanks0(b^0x01, m.d[0]|0x01)
	m.setPpuBanks0(b^0x02, m.d[1]&^0x01)
	m.setPpuBanks0(b^0x03, m.d[1]|0x01)
	m.setPpuBanks0(b^0x04, m.d[2])
	m.setPpuBanks0(b^0x05, m.d[3])
	m.setPpuBanks0(b^0x06, m.d[4])
	m.setPpuBanks0(b^0x07, m.d[5])
	m.setPpuMirrors()
}

func (m *mapper121) reset() {
	m.irqReq, m.irqRel, m.irqCnt, m.irqLatch = false, false, 0, 0
	m.cmd, m.a0, m.a1 = 0, 0, 0
	for i := 0; i < 8; i++ {
		m.r[i] = 0
	}
	m.r[3] = 0x80
	m.d[0], m.d[1], m.d[2], m.d[3], m.d[4], m.d[5], m.d[6], m.d[7] = 0, 2, 4, 5, 6, 7, 0, 1
	m.setCpuBanks1(m.cmd)
	m.setPpuBanks1(m.cmd)
}

func (m *mapper121) readLow(addr uint16) byte {
	if addr >= 0x5000 && addr < 0x6000 {
		return m.r[4]
	}
	return m.baseMapper.readLow(addr)
}

func (m *mapper121) writeLow(addr uint16, data byte) {
	if addr >= 0x5000 && addr < 0x6000 {
		switch data & 0x03 {
		case 0x00, 0x01:
			m.r[4] = 0x83
		case 0x02:
			m.r[4] = 0x42
		case 0x03:
			m.r[4] = 0
		}
		if addr&0x5180 == 0x5180 {
			m.r[3] = data
			m.setCpuBanks1(m.cmd)
			m.setPpuBanks1(m.cmd)
		}
	}
}

func (m *mapper121) setRegs() {
	switch m.r[5] & 0x3f {
	case 0x20, 0x29, 0x2b, 0x3f:
		m.r[0], m.r[7] = m.r[6], 1
	case 0x26, 0x28, 0x2a:
		m.r[((m.r[5]&0x0f)-0x06)>>1], m.r[7] = m.r[6], 0
	case 0x2c:
		m.r[7] = 1
		if m.r[6] != 0 {
			m.r[0] = m.r[6]
		}
	case 0x3c, 0x2f:
	default:
		m.r[5] = 0
	}
}

func (m *mapper121) write(addr uint16, data byte) {
	switch addr & 0xe003 {
	case 0x8001:
		m.r[6] = ((data & 0x01) << 5) | ((data & 0x02) << 3) | ((data & 0x04) << 1) | ((data & 0x08) >> 1) |
			((data & 0x10) >> 3) | ((data & 0x20) >> 5)
		if m.r[7] == 0 {
			m.setRegs()
		}
		i := m.cmd & 0x07
		m.d[i] = data
		b := (m.cmd & 0x80) >> 5
		switch i {
		case 0x00, 0x01:
			m.setPpuBanks0(b|(i<<1), data&^0x01)
			m.setPpuBanks0(b|(i<<1)|0x01, data|0x01)
		case 0x02, 0x03, 0x04, 0x05:
			m.setPpuBanks0(b^(0x04|(i-0x02)), data)
		}
		m.setCpuBanks1(m.cmd)
	case 0x8003:
		m.r[5] = data
		m.setRegs()
		fallthrough
	case 0x8000:
		if data&0x80 != m.cmd&0x80 {
			m.setPpuBanks1(data)
		}
		m.cmd = data
		m.setCpuBanks1(data)
	}

	switch addr & 0xe001 {
	case 0xa000:
		m.a0 = data
		m.setPpuMirrors()
	case 0xa001:
		m.a1 = data
	}
	switch addr & 0xe001 {
	case 0xc000:
		m.irqRel, m.irqCnt = false, data
	case 0xc001:
		m.irqRel, m.irqLatch = false, data
	case 0xe000:
		m.irqRel, m.irqReq = false, false
		m.clearIntr()
	case 0xe001:
		m.irqRel, m.irqReq = false, true
	}
}

func (m *mapper121) hSync(scanline uint16) {
	if scanline < ScreenHeight && m.isPpuDisp() && m.irqReq && !m.irqRel {
		if scanline == 0 && m.irqCnt != 0 {
			m.irqCnt--
		}
		if m.irqCnt == 0 {
			m.irqRel, m.irqCnt = true, m.irqLatch
			m.setIntr()
		}
		m.irqCnt--
	}
}

// 122

type mapper122 struct {
	baseMapper
}

func newMapper122(bm *baseMapper) Mapper {
	return &mapper122{baseMapper: *bm}
}

func (m *mapper122) reset() {
	m.mem.setProm32kBank(0)
}

func (m *mapper122) write(addr uint16, data byte) {
	if addr == 0x6000 {
		m.mem.setVrom4kBank(0, uint32(data&0x07))
		m.mem.setVrom4kBank(4, uint32(data&0x70)>>4)
	}
}
