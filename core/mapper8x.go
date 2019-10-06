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

// 140

type mapper140 struct {
	baseMapper
}

func newMapper140(bm *baseMapper) Mapper {
	return &mapper140{baseMapper: *bm}
}

func (m *mapper140) reset() {
	m.mem.setProm32kBank(0)
	if m.mem.nVrom1kPage != 0 {
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
