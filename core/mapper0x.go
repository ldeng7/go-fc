package core

// 0x00

type mapper00 struct {
	baseMapper
}

func newMapper00(bm *baseMapper) Mapper {
	return &mapper00{baseMapper: *bm}
}

func (m *mapper00) reset() {
	switch m.mem.nProm8kPage >> 1 {
	case 1:
		m.mem.setProm16kBank(4, 0)
		m.mem.setProm16kBank(6, 0)
	case 2:
		m.mem.setProm32kBank(0)
	}
}

// 0x01

type mapper01 struct {
	baseMapper
}

func newMapper01(bm *baseMapper) Mapper {
	return &mapper01{baseMapper: *bm}
}

func (m *mapper01) reset() {
}

// 0x02

type mapper02 struct {
	baseMapper
}

func newMapper02(bm *baseMapper) Mapper {
	return &mapper02{baseMapper: *bm}
}

func (m *mapper02) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
}

func (m *mapper02) writeLow(addr uint16, data byte) {
	if m.sys.rom.bSaveRam {
		m.baseMapper.writeLow(addr, data)
	}
}

func (m *mapper02) write(addr uint16, data byte) {
	m.mem.setProm16kBank(4, uint32(data>>4))
}

// 0x03

type mapper03 struct {
	baseMapper
}

func newMapper03(bm *baseMapper) Mapper {
	return &mapper03{baseMapper: *bm}
}

func (m *mapper03) reset() {
	switch m.mem.nProm8kPage >> 1 {
	case 1:
		m.mem.setProm16kBank(4, 0)
		m.mem.setProm16kBank(6, 0)
	case 2:
		m.mem.setProm32kBank(0)
	}
}

func (m *mapper03) write(addr uint16, data byte) {
	m.mem.setVrom8kBank(uint32(data))
}

// 0x04

type mapper04 struct {
	baseMapper
}

func newMapper04(bm *baseMapper) Mapper {
	return &mapper04{baseMapper: *bm}
}

func (m *mapper04) reset() {
}

// 0x05

type mapper05 struct {
	baseMapper
}

func newMapper05(bm *baseMapper) Mapper {
	return &mapper05{baseMapper: *bm}
}

func (m *mapper05) reset() {
}

// 0x06

type mapper06 struct {
	baseMapper
}

func newMapper06(bm *baseMapper) Mapper {
	return &mapper06{baseMapper: *bm}
}

func (m *mapper06) reset() {
}

// 0x07

type mapper07 struct {
	baseMapper
}

func newMapper07(bm *baseMapper) Mapper {
	return &mapper07{baseMapper: *bm}
}

func (m *mapper07) reset() {
	m.mem.setProm32kBank(0)
	m.mem.setVramMirror(memVramMirror4L)
}

func (m *mapper07) write(addr uint16, data byte) {
	m.mem.setProm32kBank(uint32(data & 0x07))
	if data&0x10 != 0 {
		m.mem.setVramMirror(memVramMirror4H)
	} else {
		m.mem.setVramMirror(memVramMirror4L)
	}
}

// 0x08

type mapper08 struct {
	baseMapper
}

func newMapper08(bm *baseMapper) Mapper {
	return &mapper08{baseMapper: *bm}
}

func (m *mapper08) reset() {
	m.mem.setProm32kBank(0)
	m.mem.setVrom8kBank(0)
}

func (m *mapper08) write(addr uint16, data byte) {
	m.mem.setProm16kBank(4, uint32((data&0xf8)>>3))
	m.mem.setVrom8kBank(uint32(data & 0x07))
}

// 0x09

type mapper09 struct {
	baseMapper
}

func newMapper09(bm *baseMapper) Mapper {
	return &mapper09{baseMapper: *bm}
}

func (m *mapper09) reset() {
}

// 0x0a

type mapper0a struct {
	baseMapper
}

func newMapper0a(bm *baseMapper) Mapper {
	return &mapper0a{baseMapper: *bm}
}

func (m *mapper0a) reset() {
}

// 0x0b

type mapper0b struct {
	baseMapper
}

func newMapper0b(bm *baseMapper) Mapper {
	return &mapper0b{baseMapper: *bm}
}

func (m *mapper0b) reset() {
	m.mem.setProm32kBank(0)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
	m.mem.setVramMirror(memVramMirrorV)
}

func (m *mapper0b) write(addr uint16, data byte) {
	m.mem.setProm32kBank(uint32(data))
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(uint32(data >> 4))
	}
}

// 0x0c

type mapper0c struct {
	baseMapper
}

func newMapper0c(bm *baseMapper) Mapper {
	return &mapper0c{baseMapper: *bm}
}

func (m *mapper0c) reset() {
}

// 0x0d

type mapper0d struct {
	baseMapper
}

func newMapper0d(bm *baseMapper) Mapper {
	return &mapper0d{baseMapper: *bm}
}

func (m *mapper0d) reset() {
	m.mem.setProm32kBank(0)
	m.mem.setCram4kBank(0, 0)
	m.mem.setCram4kBank(4, 0)
}

func (m *mapper0d) write(addr uint16, data byte) {
	m.mem.setProm32kBank(uint32((data & 0x30) >> 4))
	m.mem.setCram4kBank(4, uint32(data&0x03))
}

// 0x0f

type mapper0f struct {
	baseMapper
}

func newMapper0f(bm *baseMapper) Mapper {
	return &mapper0f{baseMapper: *bm}
}

func (m *mapper0f) reset() {
}
