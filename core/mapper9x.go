package core

// 148

type mapper148 struct {
	baseMapper
}

func newMapper148(bm *baseMapper) Mapper {
	return &mapper148{baseMapper: *bm}
}

func (m *mapper148) reset() {
	m.mem.setProm32kBank(0)
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper148) write(addr uint16, data byte) {
	m.mem.setProm32kBank(uint32(data>>3) & 0x01)
	m.mem.setVrom8kBank(uint32(data & 0x07))
}

// 150

type mapper150 struct {
	baseMapper
	r0, r1, r2, r3, r4 byte
	cmd                byte
}

func newMapper150(bm *baseMapper) Mapper {
	return &mapper150{baseMapper: *bm}
}

func (m *mapper150) reset() {
	m.r0, m.r1, m.r2, m.r3, m.r4 = 0, 0, 0, 0, 0
	m.cmd = 0
	m.mem.setProm32kBank(0)
	m.mem.setVrom8kBank(0)
}

func (m *mapper150) readLow(addr uint16) byte {
	if addr&0x4100 == 0x4100 {
		return (^m.cmd) & 0x3f
	}
	return 0
}

func (m *mapper150) write(addr uint16, data byte) {
	if addr&0x4101 == 0x4100 {
		m.cmd = data & 0x07
		return
	}
	switch m.cmd {
	case 0x02:
		m.r0 = data & 0x01
		m.r3 = (data & 0x01) << 3
	case 0x04:
		m.r4 = (data & 0x01) << 2
	case 0x05:
		m.r0 = data & 0x07
	case 0x06:
		m.r1 = data & 0x03
	case 0x07:
		m.r2 = data >> 1
	}
	m.mem.setProm32kBank(uint32(m.r0))
	m.mem.setVrom8kBank(uint32(m.r1 | m.r3 | m.r4))
	switch m.r2 & 0x03 {
	case 0x00:
		m.mem.setVramMirror(memVramMirrorV)
	case 0x01:
		m.mem.setVramMirror(memVramMirrorH)
	case 0x02:
		m.mem.setVramBank(0, 1, 1, 1)
	case 0x03:
		m.mem.setVramMirror(memVramMirror4L)
	}
}

// 151

type mapper151 struct {
	baseMapper
}

func newMapper151(bm *baseMapper) Mapper {
	return &mapper151{baseMapper: *bm}
}

func (m *mapper151) reset() {
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
}

func (m *mapper151) write(addr uint16, data byte) {
	switch addr {
	case 0x8000, 0xa000, 0xc000:
		m.mem.setProm8kBank((byte(addr>>13)&0x03)+4, uint32(data))
	case 0xe000, 0xf000:
		m.mem.setVrom4kBank((byte(addr>>12)&0x01)<<2, uint32(data))
	}
}
