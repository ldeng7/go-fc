package core

// 0x00

type mapper00 struct {
	baseMapper
}

func newMapper00(bm *baseMapper) Mapper {
	return &mapper00{baseMapper: *bm}
}

func (m *mapper00) init() {
	switch m.mem.prom16kSize {
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

func (m *mapper01) init() {
}

// 0x02

type mapper02 struct {
	baseMapper
}

func newMapper02(bm *baseMapper) Mapper {
	return &mapper02{baseMapper: *bm}
}

func (m *mapper02) init() {
}

// 0x03

type mapper03 struct {
	baseMapper
}

func newMapper03(bm *baseMapper) Mapper {
	return &mapper03{baseMapper: *bm}
}

func (m *mapper03) init() {
	switch m.mem.prom16kSize {
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

func (m *mapper04) init() {
}

// 0x05

type mapper05 struct {
	baseMapper
}

func newMapper05(bm *baseMapper) Mapper {
	return &mapper05{baseMapper: *bm}
}

func (m *mapper05) init() {
}

// 0x06

type mapper06 struct {
	baseMapper
}

func newMapper06(bm *baseMapper) Mapper {
	return &mapper06{baseMapper: *bm}
}

func (m *mapper06) init() {
}

// 0x07

type mapper07 struct {
	baseMapper
}

func newMapper07(bm *baseMapper) Mapper {
	return &mapper07{baseMapper: *bm}
}

func (m *mapper07) init() {
}

// 0x08

type mapper08 struct {
	baseMapper
}

func newMapper08(bm *baseMapper) Mapper {
	return &mapper08{baseMapper: *bm}
}

func (m *mapper08) init() {
}

// 0x09

type mapper09 struct {
	baseMapper
}

func newMapper09(bm *baseMapper) Mapper {
	return &mapper09{baseMapper: *bm}
}

func (m *mapper09) init() {
}

// 0x0a

type mapper0a struct {
	baseMapper
}

func newMapper0a(bm *baseMapper) Mapper {
	return &mapper0a{baseMapper: *bm}
}

func (m *mapper0a) init() {
}

// 0x0b

type mapper0b struct {
	baseMapper
}

func newMapper0b(bm *baseMapper) Mapper {
	return &mapper0b{baseMapper: *bm}
}

func (m *mapper0b) init() {
}

// 0x0c

type mapper0c struct {
	baseMapper
}

func newMapper0c(bm *baseMapper) Mapper {
	return &mapper0c{baseMapper: *bm}
}

func (m *mapper0c) init() {
}

// 0x0d

type mapper0d struct {
	baseMapper
}

func newMapper0d(bm *baseMapper) Mapper {
	return &mapper0d{baseMapper: *bm}
}

func (m *mapper0d) init() {
}

// 0x0e

type mapper0e struct {
	baseMapper
}

func newMapper0e(bm *baseMapper) Mapper {
	return &mapper0e{baseMapper: *bm}
}

func (m *mapper0e) init() {
}

// 0x0f

type mapper0f struct {
	baseMapper
}

func newMapper0f(bm *baseMapper) Mapper {
	return &mapper0f{baseMapper: *bm}
}

func (m *mapper0f) init() {
}
