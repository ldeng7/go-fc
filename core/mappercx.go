package core

// 193

type mapper193 struct {
	baseMapper
}

func newMapper193(bm *baseMapper) Mapper {
	return &mapper193{baseMapper: *bm}
}

func (m *mapper193) reset() {
	m.mem.setProm32kBank((m.nProm8kPage >> 2) - 1)
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper193) write(addr uint16, data byte) {
	switch addr {
	case 0x6000:
		m.mem.setVrom2kBank(0, uint32((data>>1)&0x7e))
		m.mem.setVrom2kBank(2, uint32((data>>1)&0x7e)+1)
	case 0x6001:
		m.mem.setVrom2kBank(4, uint32(data>>1))
	case 0x6002:
		m.mem.setVrom2kBank(6, uint32(data>>1))
	case 0x6003:
		m.mem.setProm32kBank(0)
	}
}

// 194

type mapper194 struct {
	baseMapper
}

func newMapper194(bm *baseMapper) Mapper {
	return &mapper194{baseMapper: *bm}
}

func (m *mapper194) reset() {
	m.mem.setProm32kBank(m.nProm8kPage>>2 - 1)
}

func (m *mapper194) write(addr uint16, data byte) {
	m.mem.setProm8kBank(3, uint32(data))
}

// 198

type mapper198 struct {
	baseMapper
	r   byte
	p   [2]byte
	c   [8]byte
	buf [4096]byte
}

func newMapper198(bm *baseMapper) Mapper {
	return &mapper198{baseMapper: *bm}
}

func (m *mapper198) setCpuBanks() {
	if m.r&0x40 != 0 {
		m.mem.setProm32kBank4(m.nProm8kPage-2, uint32(m.p[1]), uint32(m.p[0]), m.nProm8kPage-1)
	} else {
		m.mem.setProm32kBank4(uint32(m.p[0]), uint32(m.p[1]), m.nProm8kPage-2, m.nProm8kPage-1)
	}
}

func (m *mapper198) setPpuBanks() {
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

func (m *mapper198) reset() {
	m.r, m.p[0], m.p[1] = 0, 0, 1
	for i := byte(0); i < 8; i++ {
		m.c[i] = i
	}
	m.setCpuBanks()
	m.setPpuBanks()
}

func (m *mapper198) readLow(addr uint16) byte {
	if addr > 0x4018 && addr < 0x6000 {
		return m.cpuBanks[addr>>13][addr&0x1fff]
	}
	return m.buf[addr&0x0fff]
}

func (m *mapper198) writeLow(addr uint16, data byte) {
	if addr > 0x4018 && addr < 0x6000 {
		m.cpuBanks[addr>>13][addr&0x1fff] = data
	} else {
		m.buf[addr&0x0fff] = data
	}
}

func (m *mapper198) write(addr uint16, data byte) {
	addr &= 0xe001
	switch addr {
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
		case 0x06:
			if data >= 0x50 {
				data &= 0x4f
			}
			m.p[0] = data
			m.setPpuBanks()
		case 0x07:
			m.p[1] = data
			m.setPpuBanks()
		}
	case 0xa000:
		if !m.sys.rom.b4Screen {
			if data&0x01 != 0 {
				m.mem.setVramMirror(memVramMirrorH)
			} else {
				m.mem.setVramMirror(memVramMirrorV)
			}
		}
	}
}

// 199

type mapper199 struct {
	baseMapper
	irqEn    bool
	irqReq   bool
	irqCnt   byte
	irqLatch byte
	we       byte
	jm       bool
	r        byte
	jmData   [3]byte
	p        [4]byte
	c        [8]byte
}

func newMapper199(bm *baseMapper) Mapper {
	bm.sys.tvFormat = tvFormats[1]
	return &mapper199{baseMapper: *bm}
}

func (m *mapper199) setCpuBanks() {
	for i := byte(0); i < 4; i++ {
		j := i ^ ((m.r >> 5) & (^(i << 1)) & 0x02)
		m.mem.setProm8kBank(i+4, uint32(m.p[j]))
	}
}

func (m *mapper199) setPpuBanks() {
	b := (m.r & 0x80) >> 5
	for i, c := range m.c {
		if c <= 0x07 {
			m.mem.setCram1kBank(byte(i)^b, uint32(c))
		} else {
			m.mem.setVrom1kBank(byte(i)^b, uint32(c))
		}
	}
}

func (m *mapper199) reset() {
	m.irqEn, m.irqReq, m.irqCnt, m.irqLatch = false, false, 0, 0
	m.we = 0
	m.jm = false
	m.jmData[0], m.jmData[1], m.jmData[2] = 0, 0, 0
	m.r = 0
	m.p[0], m.p[1], m.p[2], m.p[3] = 0, 1, byte(m.nProm8kPage)-2, byte(m.nProm8kPage)-1
	for i := byte(0); i < 8; i++ {
		m.c[i] = i
	}
	m.setCpuBanks()
	m.setPpuBanks()
}

func (m *mapper199) readLow(addr uint16) byte {
	if addr >= 0x5000 && addr < 0x6000 {
		return m.mem.xram[addr&0x1fff]
	} else if addr >= 0x6000 {
		if m.jm {
			switch addr {
			case 0x6000:
				return m.jmData[0]
			case 0x6010:
				return m.jmData[1]
			case 0x6013:
				m.jm = false
				return m.jmData[2]
			}
		}
		switch m.we {
		case 0xe4, 0xe5, 0xe6, 0xe7, 0xec, 0xed, 0xee, 0xef:
			return m.mem.wram[(uint16(m.we&0x03)<<13)|(addr&0x1fff)]
		default:
			return m.cpuBanks[addr>>13][addr&0x1fff]
		}
	}
	return m.baseMapper.readLow(addr)
}

func (m *mapper199) writeLow(addr uint16, data byte) {
	if addr >= 0x5000 && addr < 0x6000 {
		m.mem.xram[addr&0x1fff] = data
		switch m.we {
		case 0xa1, 0xa5, 0xa9:
			m.jm = true
			switch addr {
			case 0x5000:
				m.jmData[0] = data
			case 0x5010:
				m.jmData[1] = data
			case 0x5013:
				m.jmData[2] = data
			}
		}
	} else if addr >= 0x6000 {
		switch m.we {
		case 0xe4, 0xec:
			m.mem.wram[addr&0x1fff] = data
			m.cpuBanks[addr>>13][addr&0x1fff] = data
		case 0xe5, 0xe6, 0xe7, 0xed, 0xee, 0xef:
			m.mem.wram[(uint16(m.we&0x03)<<13)|(addr&0x1fff)] = data
		default:
			m.cpuBanks[addr>>13][addr&0x1fff] = data
		}
	} else {
		m.baseMapper.writeLow(addr, data)
	}
}

func (m *mapper199) write(addr uint16, data byte) {
	switch addr & 0xe001 {
	case 0x8000:
		m.r = data
		m.setCpuBanks()
		m.setPpuBanks()
	case 0x8001:
		r := m.r & 0x0f
		switch r {
		case 0x00, 0x01:
			m.c[r<<1] = data
			m.setPpuBanks()
		case 0x02, 0x03, 0x04, 0x05:
			m.c[r+0x02] = data
			m.setPpuBanks()
		case 0x06, 0x07, 0x08, 0x09:
			m.p[r-0x06] = data
			m.setCpuBanks()
		case 0x0a, 0x0b:
			m.c[((r&0x01)<<1)+0x01] = data
			m.setPpuBanks()
		}
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

func (m *mapper199) hSync(scanline uint16) {
	if scanline < ScreenHeight && m.isPpuDisp() && m.irqEn && !m.irqReq {
		if scanline == 0 && m.irqCnt != 0 {
			m.irqCnt--
		}
		if m.irqCnt == 0 {
			m.irqReq, m.irqCnt = true, m.irqLatch
			m.setIntr()
		}
		m.irqCnt--
	}
}

// 200

type mapper200 struct {
	baseMapper
}

func newMapper200(bm *baseMapper) Mapper {
	return &mapper200{baseMapper: *bm}
}

func (m *mapper200) reset() {
	m.mem.setProm16kBank(4, 0)
	m.mem.setProm16kBank(6, 0)
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper200) write(addr uint16, data byte) {
	b := uint32(addr) & 0x07
	m.mem.setProm16kBank(4, b)
	m.mem.setProm16kBank(6, b)
	m.mem.setVrom8kBank(b)

	if addr&0x01 != 0 {
		m.mem.setVramMirror(memVramMirrorV)
	} else {
		m.mem.setVramMirror(memVramMirrorH)
	}
}

// 201

type mapper201 struct {
	baseMapper
}

func newMapper201(bm *baseMapper) Mapper {
	return &mapper201{baseMapper: *bm}
}

func (m *mapper201) reset() {
	m.mem.setProm16kBank(4, 0)
	m.mem.setProm16kBank(6, 0)
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper201) write(addr uint16, data byte) {
	var b uint32
	if addr&0x08 != 0 {
		b = uint32(addr) & 0x03
	}
	m.mem.setProm32kBank(b)
	m.mem.setVrom8kBank(b)
}

// 202

type mapper202 struct {
	baseMapper
}

func newMapper202(bm *baseMapper) Mapper {
	return &mapper202{baseMapper: *bm}
}

func (m *mapper202) reset() {
	m.mem.setProm16kBank(4, 6)
	m.mem.setProm16kBank(6, 7)
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper202) writeEx(addr uint16, data byte) {
	if addr >= 0x4020 {
		m.write(addr, data)
	}
}

func (m *mapper202) writeLow(addr uint16, data byte) {
	m.write(addr, data)
}

func (m *mapper202) write(addr uint16, data byte) {
	b := uint32(addr>>1) & 0x07
	m.mem.setProm16kBank(4, b)
	if addr&0x000c == 0x000c {
		m.mem.setProm16kBank(6, b+1)
	} else {
		m.mem.setProm16kBank(6, b)
	}
	m.mem.setVrom8kBank(b)

	if addr&0x01 != 0 {
		m.mem.setVramMirror(memVramMirrorH)
	} else {
		m.mem.setVramMirror(memVramMirrorV)
	}
}
