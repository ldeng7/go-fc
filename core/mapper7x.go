package core

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
	if m.mem.nProm8kPage == m.mem.nVrom1kPage>>3 {
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
		m.setCpuBanks0(4, byte(m.mem.nProm8kPage-2))
		m.setCpuBanks0(6, m.d[6])
	} else {
		m.setCpuBanks0(4, m.d[6])
		m.setCpuBanks0(6, byte(m.mem.nProm8kPage-2))
	}
	m.setCpuBanks0(5, m.d[7])
	m.setCpuBanks0(7, byte(m.mem.nProm8kPage-1))
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
