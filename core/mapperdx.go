package core

// 212

type mapper212 struct {
	baseMapper
}

func newMapper212(bm *baseMapper) Mapper {
	return &mapper212{baseMapper: *bm}
}

func (m *mapper212) reset() {
	np := m.mem.nProm8kPage
	m.mem.setProm32kBank4(np-4, np-3, np-2, np-1)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper212) readLow(addr uint16) byte {
	return ^(byte(addr&0x10) << 3)
}

func (m *mapper212) write(addr uint16, data byte) {
	if addr&0x4000 != 0 {
		m.mem.setProm32kBank(uint32(addr>>1) & 0x03)
	} else {
		m.mem.setProm16kBank(4, uint32(addr)&0x07)
		m.mem.setProm16kBank(6, uint32(addr)&0x07)
	}
	m.mem.setVrom8kBank(uint32(addr) & 0x07)
	if addr&0x80 != 0 {
		m.mem.setVramMirror(memVramMirrorH)
	} else {
		m.mem.setVramMirror(memVramMirrorV)
	}
}

// 222

type mapper222 struct {
	baseMapper
}

func newMapper222(bm *baseMapper) Mapper {
	return &mapper222{baseMapper: *bm}
}

func (m *mapper222) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
	m.mem.setVramMirror(memVramMirrorV)
}

func (m *mapper222) write(addr uint16, data byte) {
	addr &= 0xf003
	switch addr {
	case 0x8000:
		m.mem.setProm8kBank(4, uint32(data))
	case 0xa000:
		m.mem.setProm8kBank(5, uint32(data))
	case 0xb000, 0xb002, 0xc000, 0xc002, 0xd000, 0xd002, 0xe000, 0xe002:
		m.mem.setVrom1kBank((byte(addr>>11)-0x16)|(byte(addr&0x02)>>1), uint32(data))
	}
}
