package core

// 0xf0

type mapperf0 struct {
	baseMapper
}

func newMapperf0(bm *baseMapper) Mapper {
	return &mapperf0{baseMapper: *bm}
}

func (m *mapperf0) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapperf2) writeLow(addr uint16, data byte) {
	if addr >= 0x4020 && addr < 0x6000 {
		m.mem.setProm32kBank(uint32(data&0xf0) >> 4)
		m.mem.setVrom8kBank(uint32(data) & 0x0f)
	}
}

// 0xf1

type mapperf1 struct {
	baseMapper
}

func newMapperf1(bm *baseMapper) Mapper {
	return &mapperf1{baseMapper: *bm}
}

func (m *mapperf1) reset() {
	m.mem.setProm32kBank(0)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapperf1) write(addr uint16, data byte) {
	if addr == 0x8000 {
		m.mem.setProm32kBank(uint32(data))
	}
}

// 0xf2

type mapperf2 struct {
	baseMapper
}

func newMapperf2(bm *baseMapper) Mapper {
	return &mapperf2{baseMapper: *bm}
}

func (m *mapperf2) reset() {
	m.mem.setProm32kBank(0)
}

func (m *mapperf2) write(addr uint16, data byte) {
	if addr&0x0001 != 0 {
		m.mem.setProm32kBank(uint32(addr&0xf8) >> 3)
	}
}

// 0xf3

type mapperf3 struct {
	baseMapper
	r [4]byte
}

func newMapperf3(bm *baseMapper) Mapper {
	return &mapperf3{baseMapper: *bm}
}

func (m *mapperf3) reset() {
	m.mem.setProm32kBank(0)
	if m.mem.nVrom1kPage>>3 > 4 {
		m.mem.setVrom8kBank(4)
	} else if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
	m.mem.setVramMirror(memVramMirrorH)
}

func (m *mapperf3) write(addr uint16, data byte) {
	switch addr & 0x4101 {
	case 0x4100:
		m.r[0] = data
	case 0x4101:
		switch m.r[0] & 0x07 {
		case 0x00:
			m.r[1], m.r[2] = 0, 3
		case 0x04:
			m.r[2] = (m.r[2] & 0x06) | (data & 0x01)
		case 0x05:
			m.r[1] = data & 0x01
		case 0x06:
			m.r[2] = (m.r[2] & 0x01) | ((data & 0x03) << 1)
		case 0x07:
			m.r[3] = data & 0x01
		}
		m.mem.setProm32kBank(uint32(m.r[1]))

		m.mem.setVrom8kBank(uint32(m.r[2]))
		if m.r[3] != 0 {
			m.mem.setVramMirror(memVramMirrorV)
		} else {
			m.mem.setVramMirror(memVramMirrorH)
		}
	}
}

// 0xf4

type mapperf4 struct {
	baseMapper
}

func newMapperf4(bm *baseMapper) Mapper {
	return &mapperf4{baseMapper: *bm}
}

func (m *mapperf4) reset() {
	m.mem.setProm32kBank(0)
}

func (m *mapperf4) write(addr uint16, data byte) {
	if addr >= 0x8065 && addr <= 0x80a4 {
		m.mem.setProm32kBank(uint32(addr-0x8065) & 0x03)
	}
	if addr >= 0x80a5 && addr <= 0x80e4 {
		m.mem.setVrom8kBank(uint32(addr-0x80a5) & 0x07)
	}
}

// 0xf5

type mapperf5 struct {
	baseMapper
}

func newMapperf5(bm *baseMapper) Mapper {
	return &mapperf5{baseMapper: *bm}
}

func (m *mapperf5) reset() {
}

// 0xf6

type mapperf6 struct {
	baseMapper
}

func newMapperf6(bm *baseMapper) Mapper {
	return &mapperf6{baseMapper: *bm}
}

func (m *mapperf6) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
}

func (m *mapperf6) write(addr uint16, data byte) {
	if addr < 0x6000 || addr >= 0x8000 {
		return
	}
	switch addr {
	case 0x6000, 0x6001, 0x6002, 0x6003:
		m.mem.setProm8kBank(byte(addr)+4, uint32(data))
	case 0x6004, 0x6005, 0x6006, 0x6007:
		m.mem.setVrom2kBank((byte(addr)-4)<<1, uint32(data))
	default:
		m.mem.cpuBanks[addr>>13][addr&0x1fff] = data
	}
}

// 0xf7

type mapperf7 struct {
	baseMapper
}

func newMapperf7(bm *baseMapper) Mapper {
	return &mapperf7{baseMapper: *bm}
}

func (m *mapperf7) reset() {
}

// 0xf8

type mapperf8 struct {
	baseMapper
}

func newMapperf8(bm *baseMapper) Mapper {
	return &mapperf8{baseMapper: *bm}
}

func (m *mapperf8) reset() {
}

// 0xf9

type mapperf9 struct {
	baseMapper
}

func newMapperf9(bm *baseMapper) Mapper {
	return &mapperf9{baseMapper: *bm}
}

func (m *mapperf9) reset() {
}

// 0xfa

type mapperfa struct {
	baseMapper
}

func newMapperfa(bm *baseMapper) Mapper {
	return &mapperfa{baseMapper: *bm}
}

func (m *mapperfa) reset() {
}

// 0xfb

type mapperfb struct {
	baseMapper
}

func newMapperfb(bm *baseMapper) Mapper {
	return &mapperfb{baseMapper: *bm}
}

func (m *mapperfb) reset() {
}

// 0xfc

type mapperfc struct {
	baseMapper
}

func newMapperfc(bm *baseMapper) Mapper {
	return &mapperfc{baseMapper: *bm}
}

func (m *mapperfc) reset() {
}

// 0xfd

type mapperfd struct {
	baseMapper
}

func newMapperfd(bm *baseMapper) Mapper {
	return &mapperfd{baseMapper: *bm}
}

func (m *mapperfd) reset() {
}

// 0xfe

type mapperfe struct {
	baseMapper
}

func newMapperfe(bm *baseMapper) Mapper {
	return &mapperfe{baseMapper: *bm}
}

func (m *mapperfe) reset() {
}

// 0xff

type mapperff struct {
	baseMapper
}

func newMapperff(bm *baseMapper) Mapper {
	return &mapperff{baseMapper: *bm}
}

func (m *mapperff) reset() {
}
