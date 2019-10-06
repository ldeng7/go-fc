package core

// 049

type mapper049 struct {
	baseMapper
	irqReq   bool
	irqRel   bool
	irqCnt   byte
	irqLatch byte
	cmd      byte
	a0, a1   byte
	tmp      byte
	d        [8]byte
}

func newMapper049(bm *baseMapper) Mapper {
	return &mapper049{baseMapper: *bm}
}

func (m *mapper049) setPpuBanks0a(iBank byte, data byte) { //cwrap
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom1kBank(iBank, uint32(data))
	} else {
		m.mem.setCram1kBank(iBank, uint32(data))
	}
}

func (m *mapper049) setPpuBanks0b(iBank byte, data byte) { //wrapc
	if m.tmp&0x02 == 0 {
		m.mem.setVrom1kBank(iBank, uint32(data))
	} else {
		m.mem.setCram8kBank(0)
	}
}

func (m *mapper049) setPpuMirrors() { //mwrap
	switch m.a0 {
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

func (m *mapper049) setCpuBanks1(data byte) { //FixMMC3PRG
	if data&0x40 != 0 {
		m.mem.setProm8kBank(4, 0x7e)
		m.mem.setProm8kBank(6, uint32(m.d[6]&0x7f))
	} else {
		m.mem.setProm8kBank(4, uint32(m.d[6]&0x7f))
		m.mem.setProm8kBank(6, 0x7e)
	}
	m.mem.setProm8kBank(5, uint32(m.d[7]&0x7f))
	m.mem.setProm8kBank(7, 0x7f)
}

func (m *mapper049) setPpuBanks1(data byte, alt bool) { //FixMMC3CHR / FixCHR
	f := m.setPpuBanks0a
	if alt {
		f = m.setPpuBanks0b
	}
	b := (data & 0x80) >> 5
	f(b^0x00, m.d[0]&^0x01)
	f(b^0x01, m.d[0]|0x01)
	f(b^0x02, m.d[1]&^0x01)
	f(b^0x03, m.d[1]|0x01)
	f(b^0x04, m.d[2])
	f(b^0x05, m.d[3])
	f(b^0x06, m.d[4])
	f(b^0x07, m.d[5])
	m.setPpuMirrors()
}

func (m *mapper049) reset() {
	m.irqReq, m.irqRel, m.irqCnt, m.irqLatch = false, false, 0, 0
	m.cmd, m.a0, m.a1, m.tmp = 0, 0, 0, 0
	m.d[0], m.d[1], m.d[2], m.d[3], m.d[4], m.d[5], m.d[6], m.d[7] = 0, 2, 4, 5, 6, 7, 0, 1
	m.setCpuBanks1(m.cmd)
	m.setPpuBanks1(m.cmd, false)
}

func (m *mapper049) writeLow(addr uint16, data byte) {
	if addr&0x4100 == 0x4100 {
		m.tmp = data
		m.setPpuBanks1(m.cmd, true)
	}
}

func (m *mapper049) write(addr uint16, data byte) {
	switch addr & 0xe001 {
	case 0x8000:
		if data&0x40 != m.cmd&0x40 {
			m.setCpuBanks1(data)
		}
		if data&0x80 != m.cmd&0x80 {
			m.setPpuBanks1(data, true)
		}
		m.cmd = data
	case 0x8001:
		i := m.cmd & 0x07
		m.d[i] = data
		b := (m.cmd & 0x80) >> 5
		switch i {
		case 0x00, 0x01:
			m.setPpuBanks0b(b|(i<<1), data&^0x01)
			m.setPpuBanks0b(b|(i<<1)|0x01, data|0x01)
		case 0x02, 0x03, 0x04, 0x05:
			m.setPpuBanks0b(b^(0x04|(i-0x02)), data)
		case 0x06:
			if m.cmd&0x40 == 0 {
				m.mem.setProm8kBank(4, uint32(data&0x7f))
			} else {
				m.mem.setProm8kBank(6, uint32(data&0x7f))
			}
		case 0x07:
			m.mem.setProm8kBank(5, uint32(data&0x7f))
		}
	case 0xa000:
		m.a0 = data
		m.setPpuMirrors()
	case 0xa001:
		m.a1 = data
	case 0xc000:
		m.irqRel = false
		m.irqCnt = data
	case 0xc001:
		m.irqRel = false
		m.irqLatch = data
	case 0xe000:
		m.irqRel = false
		m.irqReq = false
		m.clearIntr()
	case 0xe001:
		m.irqRel = false
		m.irqReq = true
	}
}

func (m *mapper049) hSync(scanline uint16) {
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

// 052

type mapper052 struct {
	baseMapper
	irqEn    bool
	irqCnt   byte
	irqLatch byte
	irqClk   uint16
	r        [9]byte
}

func newMapper052(bm *baseMapper) Mapper {
	return &mapper052{baseMapper: *bm}
}

func (m *mapper052) reset() {
	m.irqEn, m.irqCnt, m.irqLatch, m.irqClk = false, 0, 0, 0
	for i := byte(0); i < 8; i++ {
		m.r[i] = i
	}
	m.r[8] = 0
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	m.mem.setVrom8kBank(0)
	m.sys.renderMode = RenderModePreAll
}

func (m *mapper052) write(addr uint16, data byte) {
	switch addr {
	case 0x8000:
		if m.r[8] != 0 {
			m.mem.setProm8kBank(6, uint32(data))
		} else {
			m.mem.setProm8kBank(4, uint32(data))
		}
	case 0x9002:
		m.r[8] = data & 0x02
	case 0x9004:
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
	case 0xa000:
		m.mem.setProm8kBank(5, uint32(data))
	case 0xb000, 0xb002, 0xc000, 0xc002, 0xd000, 0xd002, 0xe000, 0xe002:
		i := ((byte((addr&0xf000)>>12) - 0x0b) << 1) | (byte(addr&0x02) >> 1)
		m.r[i] = (m.r[i] & 0xf0) | (data & 0x0f)
		m.mem.setVrom1kBank(i, uint32(m.r[i]))
	case 0xb001, 0xb003, 0xc001, 0xc003, 0xd001, 0xd003, 0xe001, 0xe003:
		i := ((byte((addr&0xf000)>>12) - 0x0b) << 1) | (byte(addr&0x02) >> 1)
		m.r[i] = (m.r[i] & 0x0f) | ((data & 0x0f) << 4)
		m.mem.setVrom1kBank(i, uint32(m.r[i]))
	case 0xf004, 0xff04:
		m.mem.ram[0x07f8] = 0x01
	case 0xf008, 0xff08:
		m.irqLatch = ^((m.mem.ram[0x1c0] << 1) + 0x11)
		m.irqEn, m.irqCnt, m.irqClk = true, m.irqLatch, 0
		m.clearIntr()
	case 0xF00C:
		m.irqEn = true
		m.clearIntr()
	}
}

func (m *mapper052) clock(nCycle int64) {
	if m.irqEn {
		m.irqClk += uint16(nCycle) * 3
		for m.irqClk >= 341 {
			m.irqClk -= 341
			m.irqCnt++
			if m.irqCnt == 0 {
				m.irqCnt = m.irqLatch
				m.setIntr()
			}
		}
	}
}

// 057

type mapper057 struct {
	baseMapper
	r byte
}

func newMapper057(bm *baseMapper) Mapper {
	return &mapper057{baseMapper: *bm}
}

func (m *mapper057) reset() {
	m.mem.setProm32kBank4(0, 1, 0, 1)
	m.mem.setVrom8kBank(0)
	m.r = 0
}

func (m *mapper057) write(addr uint16, data byte) {
	switch addr {
	case 0x8000, 0x8001, 0x8002, 0x8003:
		if data&0x40 != 0 {
			m.mem.setVrom8kBank(uint32(data&0x03) + (uint32(m.r&0x10) >> 1) + uint32(m.r&0x07))
		}
	case 0x8800:
		m.r = data
		if data&0x80 != 0 {
			m.mem.setProm32kBank((uint32(data&0x40) >> 6) | 0x02)
		} else {
			m.mem.setProm16kBank(4, uint32(data&0x60)>>5)
			m.mem.setProm16kBank(6, uint32(data&0x60)>>5)
		}
		m.mem.setVrom8kBank(uint32(data&0x07) | (uint32(data&0x10) >> 1))
		if data&0x08 != 0 {
			m.mem.setVramMirror(memVramMirrorH)
		} else {
			m.mem.setVramMirror(memVramMirrorV)
		}
	}
}

// 058

type mapper058 struct {
	baseMapper
}

func newMapper058(bm *baseMapper) Mapper {
	return &mapper058{baseMapper: *bm}
}

func (m *mapper058) reset() {
	m.mem.setProm32kBank4(0, 1, 0, 1)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper058) write(addr uint16, data byte) {
	if addr&0x40 != 0 {
		m.mem.setProm16kBank(4, uint32(addr)&0x07)
		m.mem.setProm16kBank(6, uint32(addr)&0x07)
	} else {
		m.mem.setProm32kBank(uint32(addr&0x06) >> 1)
	}
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(uint32(addr&0x38) >> 3)
	}
	if data&0x02 != 0 {
		m.mem.setVramMirror(memVramMirrorV)
	} else {
		m.mem.setVramMirror(memVramMirrorH)
	}
}

// 060

type mapper060 struct {
	baseMapper
	idx   byte
	patch bool
}

func newMapper060(bm *baseMapper) Mapper {
	return &mapper060{baseMapper: *bm}
}

func (m *mapper060) reset() {
	switch m.sys.conf.PatchTyp {
	case 1:
		m.patch = true
		i := uint32(m.idx)
		m.mem.setProm16kBank(4, i)
		m.mem.setProm16kBank(6, i)
		m.mem.setVrom8kBank(i)
		m.idx = (m.idx + 1) & 0x03
	default:
		m.mem.setProm32kBank(0)
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper060) write(addr uint16, data byte) {
	if m.patch {
		return
	}
	if addr&0x80 != 0 {
		m.mem.setProm16kBank(4, uint32(addr&0x70)>>4)
		m.mem.setProm16kBank(6, uint32(addr&0x70)>>4)
	} else {
		m.mem.setProm32kBank(uint32(addr&0x70) >> 5)
	}
	m.mem.setVrom8kBank(uint32(addr & 0x07))
	if data&0x08 != 0 {
		m.mem.setVramMirror(memVramMirrorV)
	} else {
		m.mem.setVramMirror(memVramMirrorH)
	}
}

// 061

type mapper061 struct {
	baseMapper
}

func newMapper061(bm *baseMapper) Mapper {
	return &mapper061{baseMapper: *bm}
}

func (m *mapper061) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
}

func (m *mapper061) write(addr uint16, data byte) {
	b := uint32(data)
	switch addr & 0x30 {
	case 0x00, 0x30:
		m.mem.setProm32kBank(b & 0x0f)
	case 0x10, 0x20:
		m.mem.setProm16kBank(4, ((b&0x0f)<<1)|((b&0x20)>>4))
		m.mem.setProm16kBank(6, ((b&0x0f)<<1)|((b&0x20)>>4))
	}

	if addr&0x80 != 0 {
		m.mem.setVramMirror(memVramMirrorH)
	} else {
		m.mem.setVramMirror(memVramMirrorV)
	}
}

// 062

type mapper062 struct {
	baseMapper
}

func newMapper062(bm *baseMapper) Mapper {
	return &mapper062{baseMapper: *bm}
}

func (m *mapper062) reset() {
	m.mem.setProm32kBank(0)
	m.mem.setVrom8kBank(0)
}

func (m *mapper062) write(addr uint16, data byte) {
	b := uint32(data)
	switch addr & 0xff00 {
	case 0x8100:
		m.mem.setProm8kBank(4, b)
		m.mem.setProm8kBank(5, b+1)
	case 0x8500:
		m.mem.setProm8kBank(4, b)
	case 0x8700:
		m.mem.setProm8kBank(5, b)
	default:
		for i := byte(0); i < 8; i++ {
			m.mem.setVrom1kBank(i, b+uint32(i))
		}
	}
}
