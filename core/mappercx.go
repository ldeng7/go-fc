package core

// 0xc0

type mapperc0 struct {
	baseMapper
}

func newMapperc0(bm *baseMapper) Mapper {
	return &mapperc0{baseMapper: *bm}
}

func (m *mapperc0) reset() {
}

// 0xc1

type mapperc1 struct {
	baseMapper
}

func newMapperc1(bm *baseMapper) Mapper {
	return &mapperc1{baseMapper: *bm}
}

func (m *mapperc1) reset() {
	m.mem.setProm32kBank((m.mem.nProm8kPage >> 2) - 1)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapperc1) write(addr uint16, data byte) {
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

// 0xc2

type mapperc2 struct {
	baseMapper
}

func newMapperc2(bm *baseMapper) Mapper {
	return &mapperc2{baseMapper: *bm}
}

func (m *mapperc2) reset() {
	m.mem.setProm32kBank(m.mem.nProm8kPage>>2 - 1)
}

func (m *mapperc2) write(addr uint16, data byte) {
	m.mem.setProm8kBank(3, uint32(data))
}

// 0xc3

type mapperc3 struct {
	baseMapper
}

func newMapperc3(bm *baseMapper) Mapper {
	return &mapperc3{baseMapper: *bm}
}

func (m *mapperc3) reset() {
}

// 0xc4

type mapperc4 struct {
	baseMapper
}

func newMapperc4(bm *baseMapper) Mapper {
	return &mapperc4{baseMapper: *bm}
}

func (m *mapperc4) reset() {
}

// 0xc5

type mapperc5 struct {
	baseMapper
}

func newMapperc5(bm *baseMapper) Mapper {
	return &mapperc5{baseMapper: *bm}
}

func (m *mapperc5) reset() {
}

// 0xc6

type mapperc6 struct {
	baseMapper
}

func newMapperc6(bm *baseMapper) Mapper {
	return &mapperc6{baseMapper: *bm}
}

func (m *mapperc6) reset() {
}

// 0xc7

type mapperc7 struct {
	baseMapper
}

func newMapperc7(bm *baseMapper) Mapper {
	return &mapperc7{baseMapper: *bm}
}

func (m *mapperc7) reset() {
}

// 0xc8

type mapperc8 struct {
	baseMapper
}

func newMapperc8(bm *baseMapper) Mapper {
	return &mapperc8{baseMapper: *bm}
}

func (m *mapperc8) reset() {
	m.mem.setProm16kBank(4, 0)
	m.mem.setProm16kBank(6, 0)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapperc8) write(addr uint16, data byte) {
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

// 0xc9

type mapperc9 struct {
	baseMapper
}

func newMapperc9(bm *baseMapper) Mapper {
	return &mapperc9{baseMapper: *bm}
}

func (m *mapperc9) reset() {
	m.mem.setProm16kBank(4, 0)
	m.mem.setProm16kBank(6, 0)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapperc9) write(addr uint16, data byte) {
	var b uint32
	if addr&0x08 != 0 {
		b = uint32(addr) & 0x03
	}
	m.mem.setProm32kBank(b)
	m.mem.setVrom8kBank(b)
}

// 0xca

type mapperca struct {
	baseMapper
}

func newMapperca(bm *baseMapper) Mapper {
	return &mapperca{baseMapper: *bm}
}

func (m *mapperca) reset() {
	m.mem.setProm16kBank(4, 6)
	m.mem.setProm16kBank(6, 7)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapperca) writeEx(addr uint16, data byte) {
	if addr >= 0x4020 {
		m.write(addr, data)
	}
}

func (m *mapperca) writeLow(addr uint16, data byte) {
	m.write(addr, data)
}

func (m *mapperca) write(addr uint16, data byte) {
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
