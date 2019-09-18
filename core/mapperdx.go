package core

// 0xde

type mapperde struct {
	baseMapper
}

func newMapperde(bm *baseMapper) Mapper {
	return &mapperde{baseMapper: *bm}
}

func (m *mapperde) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
	m.mem.setVramMirror(memVramMirrorV)
}

func (m *mapperde) write(addr uint16, data byte) {
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

// 0xdf

type mapperdf struct {
	baseMapper
}

func newMapperdf(bm *baseMapper) Mapper {
	return &mapperdf{baseMapper: *bm}
}

func (m *mapperdf) reset() {
}
