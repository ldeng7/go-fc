package core

// 0x60

type mapper60 struct {
	baseMapper
	r0, r1 byte
}

func newMapper60(bm *baseMapper) Mapper {
	return &mapper60{baseMapper: *bm}
}

func (m *mapper60) reset() {
	m.r0, m.r1 = 0, 0
	m.mem.setProm32kBank(0)
	m.mem.setCram4kBank(0, uint32(m.r0)<<2+uint32(m.r1))
	m.mem.setCram4kBank(0, uint32(m.r0)<<2+0x03)
	m.mem.setVramMirror(memVramMirror4L)
}

func (m *mapper60) write(addr uint16, data byte) {
	m.mem.setProm32kBank(uint32(data & 0x03))
	m.r0 = (data & 0x04) >> 2
	m.mem.setCram4kBank(0, uint32(m.r0)<<2+uint32(m.r1))
	m.mem.setCram4kBank(0, uint32(m.r0)<<2+0x03)
}

func (m *mapper60) ppuLatch(addr uint16) {
	if (addr & 0xf000) == 0x2000 {
		m.r1 = byte(addr>>8) & 0x03
		m.mem.setCram4kBank(0, uint32(m.r0)<<2+uint32(m.r1))
		m.mem.setCram4kBank(0, uint32(m.r0)<<2+0x03)
	}
}

// 0x61

type mapper61 struct {
	baseMapper
}

func newMapper61(bm *baseMapper) Mapper {
	return &mapper61{baseMapper: *bm}
}

func (m *mapper61) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper61) write(addr uint16, data byte) {
	if addr < 0xc000 {
		m.mem.setProm16kBank(6, uint32(data&0x0f))
		if data&0x08 != 0 {
			m.mem.setVramMirror(memVramMirrorV)
		} else {
			m.mem.setVramMirror(memVramMirrorH)
		}
	}
}

// 0x62

type mapper62 struct {
	baseMapper
}

func newMapper62(bm *baseMapper) Mapper {
	return &mapper62{baseMapper: *bm}
}

func (m *mapper62) reset() {
}

// 0x63

type mapper63 struct {
	baseMapper
	c byte
}

func newMapper63(bm *baseMapper) Mapper {
	return &mapper63{baseMapper: *bm}
}

func (m *mapper63) reset() {
	m.c = 0
	if m.mem.nProm8kPage > 2 {
		m.mem.setProm32kBank4(0, 1, 2, 3)
	} else if m.mem.nProm8kPage > 1 {
		m.mem.setProm32kBank4(0, 1, 0, 1)
	} else {
		m.mem.setProm32kBank4(0, 0, 0, 0)
	}
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper63) readEx(addr uint16) byte {
	if addr == 0x4020 {
		return m.c
	}
	return byte(addr >> 8)
}

func (m *mapper63) writeEx(addr uint16, data byte) {
	switch addr {
	case 0x4016:
		if data&0x04 != 0 {
			m.mem.setVrom8kBank(1)
		} else {
			m.mem.setVrom8kBank(0)
		}
	case 0x4020:
		m.c = data
	}
}

// 0x64

type mapper64 struct {
	baseMapper
}

func newMapper64(bm *baseMapper) Mapper {
	return &mapper64{baseMapper: *bm}
}

func (m *mapper64) reset() {
}

// 0x65

type mapper65 struct {
	baseMapper
}

func newMapper65(bm *baseMapper) Mapper {
	return &mapper65{baseMapper: *bm}
}

func (m *mapper65) reset() {
	m.mem.setProm32kBank(0)
	m.mem.setVrom8kBank(0)
}

func (m *mapper65) writeLow(addr uint16, data byte) {
	if addr >= 0x6000 {
		m.mem.setVrom8kBank(uint32(data & 0x03))
	}
}

func (m *mapper65) write(addr uint16, data byte) {
	m.mem.setVrom8kBank(uint32(data & 0x03))
}

// 0x66

type mapper66 struct {
	baseMapper
}

func newMapper66(bm *baseMapper) Mapper {
	return &mapper66{baseMapper: *bm}
}

func (m *mapper66) reset() {
}

// 0x67

type mapper67 struct {
	baseMapper
}

func newMapper67(bm *baseMapper) Mapper {
	return &mapper67{baseMapper: *bm}
}

func (m *mapper67) reset() {
}

// 0x68

type mapper68 struct {
	baseMapper
}

func newMapper68(bm *baseMapper) Mapper {
	return &mapper68{baseMapper: *bm}
}

func (m *mapper68) reset() {
}

// 0x69

type mapper69 struct {
	baseMapper
}

func newMapper69(bm *baseMapper) Mapper {
	return &mapper69{baseMapper: *bm}
}

func (m *mapper69) reset() {
}

// 0x6a

type mapper6a struct {
	baseMapper
}

func newMapper6a(bm *baseMapper) Mapper {
	return &mapper6a{baseMapper: *bm}
}

func (m *mapper6a) reset() {
}

// 0x6b

type mapper6b struct {
	baseMapper
}

func newMapper6b(bm *baseMapper) Mapper {
	return &mapper6b{baseMapper: *bm}
}

func (m *mapper6b) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	m.mem.setVrom8kBank(0)
}

func (m *mapper6b) write(addr uint16, data byte) {
	m.mem.setProm32kBank(uint32(data>>1) & 0x03)
	m.mem.setVrom8kBank(uint32(data & 0x07))
}

// 0x6c

type mapper6c struct {
	baseMapper
}

func newMapper6c(bm *baseMapper) Mapper {
	return &mapper6c{baseMapper: *bm}
}

func (m *mapper6c) reset() {
	m.mem.setProm32kBank4(12, 13, 14, 15)
	m.mem.setProm8kBank(3, 0)
}

func (m *mapper6c) write(addr uint16, data byte) {
	m.mem.setProm8kBank(3, uint32(data))
}

// 0x6d

type mapper6d struct {
	baseMapper
}

func newMapper6d(bm *baseMapper) Mapper {
	return &mapper6d{baseMapper: *bm}
}

func (m *mapper6d) reset() {
}

// 0x6e

type mapper6e struct {
	baseMapper
}

func newMapper6e(bm *baseMapper) Mapper {
	return &mapper6e{baseMapper: *bm}
}

func (m *mapper6e) reset() {
}

// 0x6f

type mapper6f struct {
	baseMapper
}

func newMapper6f(bm *baseMapper) Mapper {
	return &mapper6f{baseMapper: *bm}
}

func (m *mapper6f) reset() {
}
