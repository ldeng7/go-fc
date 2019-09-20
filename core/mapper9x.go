package core

// 0x97

type mapper97 struct {
	baseMapper
}

func newMapper97(bm *baseMapper) Mapper {
	return &mapper97{baseMapper: *bm}
}

func (m *mapper97) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
}

func (m *mapper97) write(addr uint16, data byte) {
	switch addr {
	case 0x8000, 0xa000, 0xc000:
		m.mem.setProm8kBank((byte(addr>>13)&0x03)+4, uint32(data))
	case 0xe000, 0xf000:
		m.mem.setVrom4kBank((byte(addr>>12)&0x01)<<2, uint32(data))
	}
}
