package core

// 132

type mapper132 struct {
	baseMapper
	r [4]byte
}

func newMapper132(bm *baseMapper) Mapper {
	return &mapper132{baseMapper: *bm}
}

func (m *mapper132) reset() {
	for i := 0; i < 4; i++ {
		m.r[i] = 0
	}
	m.mem.setProm32kBank(0)
	m.mem.setVrom8kBank(0)
}

func (m *mapper132) readLow(addr uint16) byte {
	if addr == 0x4100 && m.r[3] != 0 {
		return m.r[2]
	}
	return 0
}

func (m *mapper132) writeLow(addr uint16, data byte) {
	if addr >= 0x4100 && addr <= 0x4103 {
		m.r[addr&0x03] = data
	}
}

func (m *mapper132) write(addr uint16, data byte) {
	m.mem.setProm32kBank(uint32(m.r[2]>>2) & 0x01)
	m.mem.setVrom8kBank(uint32(m.r[2] & 0x03))
}

// 133

type mapper133 struct {
	baseMapper
}

func newMapper133(bm *baseMapper) Mapper {
	return &mapper133{baseMapper: *bm}
}

func (m *mapper133) reset() {
	m.mem.setProm32kBank(0)
	m.mem.setVrom8kBank(0)
}

func (m *mapper133) writeLow(addr uint16, data byte) {
	if addr == 0x4120 {
		m.mem.setProm32kBank(uint32(data&0x04) >> 2)
		m.mem.setVrom8kBank(uint32(data & 0x03))
	}
	m.cpuBanks[addr>>13][addr&0x1fff] = data
}

// 134

type mapper134 struct {
	baseMapper
	cmd  byte
	p, c byte
}

func newMapper134(bm *baseMapper) Mapper {
	return &mapper134{baseMapper: *bm}
}

func (m *mapper134) reset() {
	m.mem.setProm32kBank(0)
	m.mem.setVrom8kBank(0)
}

func (m *mapper134) writeLow(addr uint16, data byte) {
	switch addr & 0x4101 {
	case 0x4100:
		m.cmd = data & 0x07
	case 0x4101:
		switch m.cmd {
		case 0:
			m.p, m.c = 0, 3
		case 4:
			m.c = (m.c & 0x03) | ((data & 0x07) << 2)
		case 5:
			m.p = data & 0x07
		case 6:
			m.c = (m.c & 0x1c) | (data & 0x03)
		case 7:
			if data&0x01 != 0 {
				m.mem.setVramMirror(memVramMirrorH)
			} else {
				m.mem.setVramMirror(memVramMirrorV)
			}
		}
	}
	m.mem.setProm32kBank(uint32(m.p))
	m.mem.setVrom8kBank(uint32(m.c))
	m.cpuBanks[addr>>13][addr&0x1fff] = data
}

// 135

type mapper135 struct {
	baseMapper
	cmd byte
	c   [5]byte
}

func newMapper135(bm *baseMapper) Mapper {
	return &mapper135{baseMapper: *bm}
}

func (m *mapper135) setPpuBanks() {
	b := uint32(m.c[4]) << 4
	m.mem.setVrom2kBank(0, uint32(m.c[0]<<1)|b)
	m.mem.setVrom2kBank(2, uint32(m.c[1]<<1)|b|0x01)
	m.mem.setVrom2kBank(4, uint32(m.c[2]<<1)|b)
	m.mem.setVrom2kBank(6, uint32(m.c[3]<<1)|b|0x01)
}

func (m *mapper135) reset() {
	m.cmd = 0
	m.c[0], m.c[1], m.c[2], m.c[3], m.c[4] = 0, 0, 0, 0, 0
	m.mem.setProm32kBank(0)
	m.setPpuBanks()
}

func (m *mapper135) writeLow(addr uint16, data byte) {
	switch addr & 0x4101 {
	case 0x4100:
		m.cmd = data & 0x07
	case 0x4101:
		switch m.cmd {
		case 0, 1, 2, 3, 4:
			m.c[m.cmd] = data & 0x07
			m.setPpuBanks()
		case 5:
			m.mem.setProm32kBank(uint32(data & 0x07))
		case 7:
			switch (data >> 1) & 0x03 {
			case 0x00, 0x03:
				m.mem.setVramMirror(memVramMirror4L)
			case 0x01:
				m.mem.setVramMirror(memVramMirrorH)
			case 0x02:
				m.mem.setVramMirror(memVramMirrorV)
			}
		}
	}
	m.cpuBanks[addr>>13][addr&0x1fff] = data
}

// 140

type mapper140 struct {
	baseMapper
}

func newMapper140(bm *baseMapper) Mapper {
	return &mapper140{baseMapper: *bm}
}

func (m *mapper140) reset() {
	m.mem.setProm32kBank(0)
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper140) writeLow(addr uint16, data byte) {
	m.mem.setProm32kBank(uint32(data&0xf0) >> 4)
	m.mem.setVrom8kBank(uint32(data & 0x0f))
}

// 141

type mapper141 struct {
	baseMapper
	r   [8]byte
	cmd byte
}

func newMapper141(bm *baseMapper) Mapper {
	return &mapper141{baseMapper: *bm}
}

func (m *mapper141) reset() {
	for i := 0; i < 8; i++ {
		m.r[i] = 0
	}
	m.cmd = 0
	m.mem.setVrom8kBank(0)
	m.setBanks()
}

func (m *mapper141) writeLow(addr uint16, data byte) {
	addr &= 0x4101
	if addr&0x4101 == 0x4100 {
		m.cmd = data
	} else {
		m.r[m.cmd&0x07] = data
		m.setBanks()
	}
}

func (m *mapper141) setBanks() {
	m.mem.setProm32kBank(uint32(m.r[5] & 0x07))
	for i := byte(0); i < 4; i++ {
		j := i
		if m.r[7]&0x01 != 0 {
			j = 0
		}
		m.mem.setVrom2kBank(i<<1, uint32(((m.r[j]&0x07)<<1)|((m.r[4]&0x07)<<4)|(i&0x01)))
	}
	if m.r[7]&0x01 != 0 {
		m.mem.setVramMirror(memVramMirrorV)
	} else {
		m.mem.setVramMirror(memVramMirrorH)
	}
}

// 142

type mapper142 struct {
	baseMapper
	irqEn  bool
	irqCnt uint16
	p      byte
}

func newMapper142(bm *baseMapper) Mapper {
	return &mapper142{baseMapper: *bm}
}

func (m *mapper142) reset() {
	m.irqEn, m.irqCnt, m.p = false, 0, 0
	m.mem.setProm8kBank(3, 0)
	m.mem.setProm8kBank(7, 15)
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper142) write(addr uint16, data byte) {
	addr &= 0xf000
	switch addr {
	case 0x8000, 0x9000, 0xa000, 0xb000:
		i := (addr >> 10) & 0x0c
		m.irqCnt = (m.irqCnt &^ (0x000f << i)) | (uint16(data&0x0f) << i)
	case 0xc000:
		m.irqEn = data&0x0f != 0
		m.clearIntr()
	case 0xe000:
		m.p = data & 0x0f
	case 0xf000:
		switch m.p {
		case 0x01, 0x02, 0x03, 0x04:
			m.mem.setProm8kBank((m.p&0x03)+3, uint32(data&0x0f))
		}
	}
}

func (m *mapper142) hSync(scanline uint16) {
	if m.irqEn {
		if m.irqCnt > 65422 {
			m.irqCnt = 0
			m.setIntr()
		} else {
			m.irqCnt += 113
		}
	}
}
