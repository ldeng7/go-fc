package core

// 064

type mapper064 struct {
	baseMapper
	irqEn      bool
	irqMode    bool
	irqReset   bool
	irqLatch   byte
	irqCnt     int16
	irqCnt2    int16
	r0, r1, r2 byte
}

func newMapper064(bm *baseMapper) Mapper {
	return &mapper064{baseMapper: *bm}
}

func (m *mapper064) reset() {
	m.r0, m.r1, m.r2 = 0, 0, 0
	m.irqEn, m.irqMode, m.irqReset = false, false, false
	m.irqLatch, m.irqCnt, m.irqCnt2 = 0, 0, 0
	b := m.nProm8kPage - 1
	m.mem.setProm32kBank4(b, b, b, b)
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper064) write(addr uint16, data byte) {
	switch addr & 0xf003 {
	case 0x8000:
		m.r0, m.r1, m.r2 = data&0x0f, data&0x40, data&0x80
	case 0x8001:
		r := m.r0
		switch r {
		case 0x00, 0x01:
			i := r << 1
			if m.r2 != 0 {
				i += 4
			}
			m.mem.setVrom1kBank(i, uint32(data))
			m.mem.setVrom1kBank(i+1, uint32(data)+1)
		case 0x02, 0x03, 0x04, 0x05:
			i := r - 2
			if m.r2 == 0 {
				i += 4
			}
			m.mem.setVrom1kBank(i, uint32(data))
		case 0x06, 0x07:
			i := r - 2
			if m.r1 != 0 {
				i++
			}
			m.mem.setProm8kBank(i, uint32(data))
		case 0x08, 0x09:
			i := ((r & 0x01) << 1) + 1
			m.mem.setVrom1kBank(i, uint32(data))
		case 0x0f:
			if m.r1 != 0 {
				m.mem.setProm8kBank(4, uint32(data))
			} else {
				m.mem.setProm8kBank(6, uint32(data))
			}
		}
	case 0xa000:
		if data&0x01 != 0 {
			m.mem.setVramMirror(memVramMirrorH)
		} else {
			m.mem.setVramMirror(memVramMirrorV)
		}
	case 0xc000:
		m.irqLatch = data
		if m.irqReset {
			m.irqCnt = int16(m.irqLatch)
		}
	case 0xc001:
		m.irqMode, m.irqReset, m.irqCnt = data&0x01 != 0, true, int16(m.irqLatch)
	case 0xe000:
		m.irqEn = false
		if m.irqReset {
			m.irqCnt = int16(m.irqLatch)
		}
		m.clearIntr()
	case 0xe001:
		m.irqEn = true
		if m.irqReset {
			m.irqCnt = int16(m.irqLatch)
		}
	}
}

func (m *mapper064) hSync(scanline uint16) {
	if !m.irqMode {
		m.irqReset = false
		if scanline < ScreenHeight && m.isPpuDisp() && m.irqCnt >= 0 {
			m.irqCnt--
			if m.irqCnt < 0 && m.irqEn {
				m.irqReset = true
				m.setIntr()
			}
		}
	}
}

func (m *mapper064) clock(nCycle int64) {
	if m.irqMode {
		m.irqCnt2 += int16(nCycle)
		for m.irqCnt2 >= 4 {
			m.irqCnt2 -= 4
			if m.irqCnt >= 0 {
				m.irqCnt--
				if m.irqCnt < 0 && m.irqEn {
					m.setIntr()
				}
			}
		}
	}
}

// 065

type mapper065 struct {
	baseMapper
	patchTyp bool
	irqEn    bool
	irqCnt   int32
	irqLatch int32
}

func newMapper065(bm *baseMapper) Mapper {
	return &mapper065{baseMapper: *bm}
}

func (m *mapper065) reset() {
	patch := m.sys.conf.PatchTyp
	if patch&0x01 != 0 {
		m.patchTyp = true
	}
	m.irqEn, m.irqCnt = false, 0
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper065) write(addr uint16, data byte) {
	switch addr {
	case 0x8000, 0xa000, 0xc000:
		m.mem.setProm8kBank((byte(addr>>13)&0x03)+4, uint32(data))
	case 0x9000:
		if !m.patchTyp {
			if data&0x40 != 0 {
				m.mem.setVramMirror(memVramMirrorV)
			} else {
				m.mem.setVramMirror(memVramMirrorH)
			}
		}
	case 0x9001:
		if m.patchTyp {
			if data&0x80 != 0 {
				m.mem.setVramMirror(memVramMirrorH)
			} else {
				m.mem.setVramMirror(memVramMirrorV)
			}
		}
	case 0x9003:
		if !m.patchTyp {
			m.irqEn = data&0x80 != 0
			m.clearIntr()
		}
	case 0x9004:
		if !m.patchTyp {
			m.irqCnt = m.irqLatch
		}
	case 0x9005:
		if !m.patchTyp {
			m.irqLatch = (m.irqLatch & 0x00ff) | (int32(data) << 8)
		} else {
			m.irqEn, m.irqCnt = data != 0, int32(data<<1)
			m.clearIntr()
		}
	case 0x9006:
		if !m.patchTyp {
			m.irqLatch = (m.irqLatch & 0xff00) | int32(data)
		} else {
			m.irqEn = true
		}
	case 0xb000, 0xb001, 0xb002, 0xb003, 0xb004, 0xb005, 0xb006, 0xb007:
		m.mem.setVrom1kBank(byte(addr)&0x07, uint32(data))
	}
}

func (m *mapper065) hSync(scanline uint16) {
	if m.patchTyp && m.irqEn {
		if m.irqCnt == 0 {
			m.setIntr()
		} else {
			m.irqCnt--
		}
	}
}

func (m *mapper065) clock(nCycle int64) {
	if !m.patchTyp && m.irqEn {
		if m.irqCnt <= 0 {
			m.setIntr()
		} else {
			m.irqCnt -= int32(nCycle)
		}
	}
}

// 066

type mapper066 struct {
	baseMapper
}

func newMapper066(bm *baseMapper) Mapper {
	return &mapper066{baseMapper: *bm}
}

func (m *mapper066) reset() {
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
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

// 067

type mapper067 struct {
	baseMapper
	irqEn  bool
	irqTog bool
	irqCnt int32
}

func newMapper067(bm *baseMapper) Mapper {
	return &mapper067{baseMapper: *bm}
}

func (m *mapper067) reset() {
	m.irqEn, m.irqTog, m.irqCnt = false, false, 0
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
	m.mem.setVrom4kBank(0, 0)
	m.mem.setVrom4kBank(4, (m.nVrom1kPage>>2)-1)
}

func (m *mapper067) write(addr uint16, data byte) {
	addr &= 0xf800
	switch addr {
	case 0x8800, 0x9800, 0xa800, 0xb800:
		m.mem.setVrom2kBank(byte(addr>>11)&0x06, uint32(data))
	case 0xc800:
		if !m.irqTog {
			m.irqCnt = (m.irqCnt & 0x00ff) | (int32(data) << 8)
		} else {
			m.irqCnt = (m.irqCnt & 0xff00) | int32(data)
		}
		m.irqTog = !m.irqTog
		m.clearIntr()
	case 0xd800:
		m.irqEn, m.irqTog = data&0x10 != 0, false
		m.clearIntr()
	case 0xe800:
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
	case 0xf800:
		m.mem.setProm16kBank(4, uint32(data))
	}
}

func (m *mapper067) clock(nCycle int64) {
	if m.irqEn {
		m.irqCnt -= int32(nCycle)
		if m.irqCnt <= 0 {
			m.irqEn, m.irqCnt = false, 0xffff
			m.setIntr()
		}
	}
}

// 068

type mapper068 struct {
	baseMapper
	r [4]byte
}

func newMapper068(bm *baseMapper) Mapper {
	return &mapper068{baseMapper: *bm}
}

func (m *mapper068) reset() {
	m.r[0], m.r[1], m.r[2], m.r[3] = 0, 0, 0, 0
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
}

func (m *mapper068) setPpuBanks() {
	if m.r[0] != 0 {
		r2, r3 := uint32(m.r[2])+0x80, uint32(m.r[3])+0x80
		var b0, b1, b2, b3 uint32
		switch m.r[1] {
		case 0x00:
			b0, b1, b2, b3 = r2, r3, r2, r3
		case 0x01:
			b0, b1, b2, b3 = r2, r2, r3, r3
		case 0x02:
			b0, b1, b2, b3 = r2, r2, r2, r2
		case 0x03:
			b0, b1, b2, b3 = r3, r3, r3, r3
		}
		m.mem.setVrom1kBank(8, b0)
		m.mem.setVrom1kBank(9, b1)
		m.mem.setVrom1kBank(10, b2)
		m.mem.setVrom1kBank(11, b3)
	} else {
		switch m.r[1] {
		case 0x00:
			m.mem.setVramMirror(memVramMirrorV)
		case 0x01:
			m.mem.setVramMirror(memVramMirrorH)
		case 0x02:
			m.mem.setVramMirror(memVramMirror4L)
		case 0x03:
			m.mem.setVramMirror(memVramMirror4H)
		}
	}
}

func (m *mapper068) write(addr uint16, data byte) {
	addr &= 0xf000
	switch addr {
	case 0x8000, 0x9000, 0xa000, 0xb000:
		m.mem.setVrom2kBank(byte(addr>>11)&0x06, uint32(data))
	case 0xc000, 0xd000:
		m.r[(byte(addr>>12)&0x01)+2] = data
		m.setPpuBanks()
	case 0xe000:
		m.r[0], m.r[1] = (data&0x10)>>4, data&0x03
		m.setPpuBanks()
	case 0xf000:
		m.mem.setProm16kBank(4, uint32(data))
	}
}

// 069

type mapper069 struct {
	baseMapper
	patchTyp bool
	irqEn    bool
	irqCnt   int32
	r        byte
}

func newMapper069(bm *baseMapper) Mapper {
	return &mapper069{baseMapper: *bm}
}

func (m *mapper069) reset() {
	if m.sys.conf.PatchTyp&0x01 != 0 {
		m.patchTyp = true
	}
	m.irqEn, m.irqCnt = false, 0
	m.r = 0
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper069) write(addr uint16, data byte) {
	addr &= 0xe000
	switch addr {
	case 0x8000:
		m.r = data
	case 0xa000:
		r := m.r & 0x0f
		switch r {
		case 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07:
			m.mem.setVrom1kBank(r&0x07, uint32(data))
		case 0x08:
			if !m.patchTyp && data&0x40 == 0 {
				m.mem.setProm8kBank(3, uint32(data))
			}
		case 0x09, 0x0a, 0x0b:
			m.mem.setProm8kBank(r-0x05, uint32(data))
		case 0x0c:
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
		case 0x0d:
			m.irqEn = data != 0
			m.clearIntr()
		case 0x0e:
			m.irqCnt = (m.irqCnt & 0xff00) | int32(data)
			m.clearIntr()
		case 0x0F:
			m.irqCnt = (m.irqCnt & 0x00ff) | (int32(data) << 8)
			m.clearIntr()
		}
	}
}

func (m *mapper069) clock(nCycle int64) {
	if m.irqEn {
		m.irqCnt -= int32(nCycle)
		if m.irqCnt <= 0 {
			m.irqEn, m.irqCnt = false, 0xffff
			m.setIntr()
		}
	}
}

// 070

type mapper070 struct {
	baseMapper
}

func newMapper070(bm *baseMapper) Mapper {
	return &mapper070{baseMapper: *bm}
}

func (m *mapper070) reset() {
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
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
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
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
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
	if m.nVrom1kPage != 0 {
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
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
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

// 074

type mapper074 struct {
	baseMapper
	patchTyp bool
	irqEn    bool
	irqReq   bool
	irqCnt   byte
	irqLatch byte
	r        byte
	p        [2]byte
	c        [8]byte
}

func newMapper074(bm *baseMapper) Mapper {
	return &mapper074{baseMapper: *bm}
}

func (m *mapper074) setCpuBanks() {
	if m.r&0x40 != 0 {
		m.mem.setProm32kBank4(m.nProm8kPage-2, uint32(m.p[1]), uint32(m.p[0]), m.nProm8kPage-1)
	} else {
		m.mem.setProm32kBank4(uint32(m.p[0]), uint32(m.p[1]), m.nProm8kPage-2, m.nProm8kPage-1)
	}
}

func (m *mapper074) setPpuBankSub(iBank byte, iPage uint32) {
	if !m.patchTyp && (iPage == 8 || iPage == 9) {
		m.mem.setCram1kBank(iBank, iPage&0x07)
	} else if m.patchTyp && iPage >= 128 {
		m.mem.setCram1kBank(iBank, iPage&0x07)
	} else {
		m.mem.setVrom1kBank(iBank, iPage)
	}
}

func (m *mapper074) setPpuBanks() {
	if m.nVrom1kPage != 0 {
		if m.r&0x80 != 0 {
			for i := byte(0); i < 4; i++ {
				m.setPpuBankSub(i, uint32(m.c[i+4]))
				m.setPpuBankSub(i+4, uint32(m.c[i]))
			}
		} else {
			for i := byte(0); i < 8; i++ {
				m.setPpuBankSub(i, uint32(m.c[i]))
			}
		}
	} else {
		if m.r&0x80 != 0 {
			for i := byte(0); i < 4; i++ {
				m.mem.setCram1kBank(i, uint32(m.c[i+4]))
				m.mem.setCram1kBank(i+4, uint32(m.c[i]))
			}
		} else {
			for i := byte(0); i < 8; i++ {
				m.mem.setCram1kBank(i, uint32(m.c[i]))
			}
		}
	}
}

func (m *mapper074) reset() {
	if m.sys.conf.PatchTyp&0x01 != 0 {
		m.patchTyp = true
		m.sys.renderMode = RenderModeTile
	}
	m.irqEn, m.irqReq, m.irqCnt, m.irqLatch = false, false, 0, 0
	m.r, m.p[0], m.p[1] = 0, 0, 1
	for i := byte(0); i < 8; i++ {
		m.c[i] = i
	}
	m.setCpuBanks()
	m.setPpuBanks()
}

func (m *mapper074) write(addr uint16, data byte) {
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
			m.c[i] = data % 0xfe
			m.c[i+1] = m.c[i] + 1
			m.setPpuBanks()
		case 0x02, 0x03, 0x04, 0x05:
			m.c[r+2] = data
			m.setPpuBanks()
		case 0x06, 0x07:
			m.p[r-6] = data
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
		m.irqCnt, m.irqReq = data, false
	case 0xc001:
		m.irqLatch, m.irqReq = data, false
	case 0xe000:
		m.irqEn, m.irqReq = false, false
		m.clearIntr()
	case 0xe001:
		m.irqEn, m.irqReq = true, false
	}
}

func (m *mapper074) hSync(scanline uint16) {
	if scanline < ScreenHeight && m.isPpuDisp() && m.irqEn && !m.irqReq {
		if scanline == 0 && m.irqCnt != 0 {
			m.irqCnt--
		}
		m.irqCnt--
		if m.irqCnt == 0xff {
			m.irqReq, m.irqCnt = true, m.irqLatch
			m.setIntr()
		}
	}
}

// 075

type mapper075 struct {
	baseMapper
	r [2]byte
}

func newMapper075(bm *baseMapper) Mapper {
	return &mapper075{baseMapper: *bm}
}

func (m *mapper075) reset() {
	m.r[0], m.r[1] = 0, 1
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper075) write(addr uint16, data byte) {
	addr &= 0xf000
	switch addr {
	case 0x8000, 0xa000, 0xc000:
		m.mem.setProm8kBank((byte(addr>>13)&0x03)+4, uint32(data))
	case 0x9000:
		m.r[0] = (m.r[0] & 0x0f) | ((data & 0x02) << 3)
		m.r[1] = (m.r[1] & 0x0f) | ((data & 0x04) << 2)
		m.mem.setVrom4kBank(0, uint32(m.r[0]))
		m.mem.setVrom4kBank(4, uint32(m.r[1]))
		if data&0x01 != 0 {
			m.mem.setVramMirror(memVramMirrorH)
		} else {
			m.mem.setVramMirror(memVramMirrorV)
		}
	case 0xe000, 0xf000:
		i := byte(addr>>12) & 0x01
		m.r[i] = (m.r[i] & 0x10) | (data & 0x0f)
		m.mem.setVrom4kBank(i<<2, uint32(m.r[i]))
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
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
	if m.nVrom1kPage >= 8 {
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
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
	if m.nVrom1kPage != 0 {
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
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper079) write(addr uint16, data byte) {
	if addr&0x0100 != 0 {
		m.mem.setProm32kBank(uint32(data>>3) & 0x01)
		m.mem.setVrom8kBank(uint32(data & 0x07))
	}
}
