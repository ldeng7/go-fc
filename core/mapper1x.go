package core

// 0x10

type mapper10 struct {
	baseMapper
}

func newMapper10(bm *baseMapper) Mapper {
	return &mapper10{baseMapper: *bm}
}

func (m *mapper10) reset() {
}

// 0x11

type mapper11 struct {
	baseMapper
}

func newMapper11(bm *baseMapper) Mapper {
	return &mapper11{baseMapper: *bm}
}

func (m *mapper11) reset() {
}

// 0x12

type mapper12 struct {
	baseMapper
}

func newMapper12(bm *baseMapper) Mapper {
	return &mapper12{baseMapper: *bm}
}

func (m *mapper12) reset() {
}

// 0x13

type mapper13 struct {
	baseMapper
}

func newMapper13(bm *baseMapper) Mapper {
	return &mapper13{baseMapper: *bm}
}

func (m *mapper13) reset() {
}

// 0x14

type mapper14 struct {
	baseMapper
}

func newMapper14(bm *baseMapper) Mapper {
	return &mapper14{baseMapper: *bm}
}

func (m *mapper14) reset() {
}

// 0x15

type mapper15 struct {
	baseMapper
}

func newMapper15(bm *baseMapper) Mapper {
	return &mapper15{baseMapper: *bm}
}

func (m *mapper15) reset() {
}

// 0x16

type mapper16 struct {
	baseMapper
}

func newMapper16(bm *baseMapper) Mapper {
	return &mapper16{baseMapper: *bm}
}

func (m *mapper16) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
}

func (m *mapper16) write(addr uint16, data byte) {
	switch addr {
	case 0x8000:
		m.mem.setProm8kBank(4, uint32(data))
	case 0x9000:
		switch data & 0x03 {
		case 0x00:
			m.mem.setVramMirror(memVramMirrorV)
		case 0x01:
			m.mem.setVramMirror(memVramMirrorH)
		case 0x02:
			m.mem.setVramMirror(memVramMirror4H)
		case 0x03:
			m.mem.setVramMirror(memVramMirror4L)
		}
	case 0xa000:
		m.mem.setProm8kBank(5, uint32(data))
	case 0xb000, 0xb001, 0xc000, 0xc001, 0xd000, 0xd001, 0xe000, 0xe001:
		m.mem.setVrom1kBank((byte(addr>>11)-0x16)|(byte(addr)&0x01), uint32(data>>1))
	}
}

// 0x17

type mapper17 struct {
	baseMapper
}

func newMapper17(bm *baseMapper) Mapper {
	return &mapper17{baseMapper: *bm}
}

func (m *mapper17) reset() {
}

// 0x18

type mapper18 struct {
	baseMapper
}

func newMapper18(bm *baseMapper) Mapper {
	return &mapper18{baseMapper: *bm}
}

func (m *mapper18) reset() {
}

// 0x19

type mapper19 struct {
	baseMapper
}

func newMapper19(bm *baseMapper) Mapper {
	return &mapper19{baseMapper: *bm}
}

func (m *mapper19) reset() {
}

// 0x1a

type mapper1a struct {
	baseMapper
}

func newMapper1a(bm *baseMapper) Mapper {
	return &mapper1a{baseMapper: *bm}
}

func (m *mapper1a) reset() {
}

// 0x1b

type mapper1b struct {
	baseMapper
}

func newMapper1b(bm *baseMapper) Mapper {
	return &mapper1b{baseMapper: *bm}
}

func (m *mapper1b) reset() {
}
