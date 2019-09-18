package core

// 0xe0

type mappere0 struct {
	baseMapper
}

func newMappere0(bm *baseMapper) Mapper {
	return &mappere0{baseMapper: *bm}
}

func (m *mappere0) reset() {
}

// 0xe1

type mappere1 struct {
	baseMapper
}

func newMappere1(bm *baseMapper) Mapper {
	return &mappere1{baseMapper: *bm}
}

func (m *mappere1) reset() {
	m.mem.setProm32kBank(0)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mappere1) write(addr uint16, data byte) {
	i := uint32((addr & 0x0f80) >> 7)
	if addr&0x1000 != 0 {
		j := i << 2
		if addr&0x0040 != 0 {
			j += 2
		}
		m.mem.setProm32kBank4(j, j+1, j, j+1)
	} else {
		m.mem.setProm32kBank(i)
	}

	m.mem.setVrom8kBank(uint32(addr & 0x003f))
	if addr&0x2000 != 0 {
		m.mem.setVramMirror(memVramMirrorH)
	} else {
		m.mem.setVramMirror(memVramMirrorV)
	}
}

// 0xe2

type mappere2 struct {
	baseMapper
	r0, r1 byte
}

func newMappere2(bm *baseMapper) Mapper {
	return &mappere2{baseMapper: *bm}
}

func (m *mappere2) reset() {
	m.mem.setProm32kBank(0)
}

func (m *mappere2) write(addr uint16, data byte) {
	if addr&0x001 != 0 {
		m.r1 = data
	} else {
		m.r0 = data
	}

	i := uint32(((m.r0 & 0x1e) >> 1) | ((m.r0 & 0x80) >> 3) | ((m.r1 & 0x01) << 5))
	if m.r0&0x20 != 0 {
		j := i << 2
		if m.r0&0x01 != 0 {
			j += 2
		}
		m.mem.setProm32kBank4(j, j+1, j, j+1)
	} else {
		m.mem.setProm32kBank(i)
	}

	if m.r0&0x40 != 0 {
		m.mem.setVramMirror(memVramMirrorV)
	} else {
		m.mem.setVramMirror(memVramMirrorH)
	}
}

// 0xe3

type mappere3 struct {
	baseMapper
}

func newMappere3(bm *baseMapper) Mapper {
	return &mappere3{baseMapper: *bm}
}

func (m *mappere3) reset() {
	m.mem.setProm32kBank4(0, 1, 0, 1)
}

func (m *mappere3) write(addr uint16, data byte) {
	i := uint32(((addr & 0x0100) >> 4) | ((addr & 0x0078) >> 3))
	if addr&0x0001 != 0 {
		m.mem.setProm32kBank(i)
	} else {
		j := i << 2
		if addr&0x0004 != 0 {
			j += 2
		}
		m.mem.setProm32kBank4(j, j+1, j, j+1)
	}
	if addr&0x0080 == 0 {
		j := (i & 0x1c) << 2
		if addr&0x0200 != 0 {
			j += 14
		}
		m.mem.setProm8kBank(6, j)
		m.mem.setProm8kBank(7, j+1)
	}

	if addr&0x0002 != 0 {
		m.mem.setVramMirror(memVramMirrorH)
	} else {
		m.mem.setVramMirror(memVramMirrorV)
	}
}

// 0xe4

type mappere4 struct {
	baseMapper
}

func newMappere4(bm *baseMapper) Mapper {
	return &mappere4{baseMapper: *bm}
}

func (m *mappere4) reset() {
	m.mem.setProm32kBank(0)
	m.mem.setVrom8kBank(0)
}

func (m *mappere4) write(addr uint16, data byte) {
	b := (uint32(addr) & 0x0780) >> 7
	switch (addr & 0x1800) >> 11 {
	case 1:
		b |= 0x10
	case 3:
		b |= 0x20
	}
	if addr&0x0020 != 0 {
		b <<= 1
		if addr&0x0040 != 0 {
			b++
		}
		m.mem.setProm16kBank(4, b<<1)
		m.mem.setProm16kBank(6, b<<1)
	} else {
		m.mem.setProm32kBank(b)
	}

	m.mem.setVrom8kBank(((uint32(addr) & 0x0f) << 2) | (uint32(data) & 0x03))
	if addr&0x2000 != 0 {
		m.mem.setVramMirror(memVramMirrorH)
	} else {
		m.mem.setVramMirror(memVramMirrorV)
	}
}

// 0xe5

type mappere5 struct {
	baseMapper
}

func newMappere5(bm *baseMapper) Mapper {
	return &mappere5{baseMapper: *bm}
}

func (m *mappere5) reset() {
	m.mem.setProm32kBank(0)
	m.mem.setVrom8kBank(0)
}

func (m *mappere5) write(addr uint16, data byte) {
	if addr&0x001e != 0 {
		b := uint32(addr) & 0x1f
		m.mem.setProm16kBank(4, b)
		m.mem.setProm16kBank(6, b)
		m.mem.setVrom8kBank(uint32(addr) & 0x0fff)
	} else {
		m.mem.setProm32kBank(0)
		m.mem.setVrom8kBank(0)
	}
	if addr&0x0020 != 0 {
		m.mem.setVramMirror(memVramMirrorH)
	} else {
		m.mem.setVramMirror(memVramMirrorV)
	}
}

// 0xe6

type mappere6 struct {
	baseMapper
}

func newMappere6(bm *baseMapper) Mapper {
	return &mappere6{baseMapper: *bm}
}

func (m *mappere6) reset() {
}

// 0xe7

type mappere7 struct {
	baseMapper
}

func newMappere7(bm *baseMapper) Mapper {
	return &mappere7{baseMapper: *bm}
}

func (m *mappere7) reset() {
	m.mem.setProm32kBank(0)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mappere7) write(addr uint16, data byte) {
	if addr&0x0020 != 0 {
		m.mem.setProm32kBank((uint32(addr) >> 1) & 0xff)
	} else {
		b := uint32(addr) & 0x1e
		m.mem.setProm16kBank(4, b)
		m.mem.setProm16kBank(6, b)
	}
	if addr&0x0080 != 0 {
		m.mem.setVramMirror(memVramMirrorH)
	} else {
		m.mem.setVramMirror(memVramMirrorV)
	}
}

// 0xe8

type mappere8 struct {
	baseMapper
	r0, r1 byte
}

func newMappere8(bm *baseMapper) Mapper {
	return &mappere8{baseMapper: *bm}
}

func (m *mappere8) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	m.r0 = 0x0c
}

func (m *mappere8) writeLow(addr uint16, data byte) {
	if addr >= 0x6000 {
		m.write(addr, data)
	}
}

func (m *mappere8) write(addr uint16, data byte) {
	if addr <= 0x9fff {
		m.r0 = (data & 0x18) >> 1
	} else {
		m.r1 = data & 0x03
	}
	m.mem.setProm16kBank(4, uint32(m.r0)|uint32(m.r1))
	m.mem.setProm16kBank(6, uint32(m.r0)|0x03)
}

// 0xe9

type mappere9 struct {
	baseMapper
}

func newMappere9(bm *baseMapper) Mapper {
	return &mappere9{baseMapper: *bm}
}

func (m *mappere9) reset() {
	m.mem.setProm32kBank(0)
}

func (m *mappere9) write(addr uint16, data byte) {
	if data&0x20 != 0 {
		m.mem.setProm16kBank(4, uint32(data&0x1f))
		m.mem.setProm16kBank(6, uint32(data&0x1f))
	} else {
		m.mem.setProm32kBank(uint32(data&0x1e) >> 1)
	}
	switch data & 0xc0 {
	case 0x00:
		m.mem.setVramBank(0, 0, 0, 1)
	case 0x40:
		m.mem.setVramMirror(memVramMirrorV)
	case 0x80:
		m.mem.setVramMirror(memVramMirrorH)
	case 0xc0:
		m.mem.setVramMirror(memVramMirror4H)
	}
}

// 0xea

type mapperea struct {
	baseMapper
	r0, r1 byte
}

func newMapperea(bm *baseMapper) Mapper {
	return &mapperea{baseMapper: *bm}
}

func (m *mapperea) setBank() {
	r0, r1 := uint32(m.r0), uint32(m.r1)
	if r0&0x40 != 0 {
		m.mem.setProm32kBank((r0 & 0x0e) | (r1 & 0x01))
		m.mem.setVrom8kBank(((r0 & 0x0e) << 2) | ((r1 >> 4) & 0x07))
	} else {
		m.mem.setProm32kBank(r0 & 0x0f)
		m.mem.setVrom8kBank(((r0 & 0x0f) << 2) | ((r1 >> 4) & 0x03))
	}
	if r0&0x80 != 0 {
		m.mem.setVramMirror(memVramMirrorH)
	} else {
		m.mem.setVramMirror(memVramMirrorV)
	}
}

func (m *mapperea) reset() {
	m.mem.setProm32kBank(0)
}

func (m *mapperea) read(addr uint16, data byte) {
	if addr >= 0xff80 && addr <= 0xff9f {
		if m.r0 == 0 {
			m.r0 = data
			m.setBank()
		}
	}
	if addr >= 0xffe8 && addr <= 0xfff7 {
		m.r1 = data
		m.setBank()
	}
}

func (m *mapperea) write(addr uint16, data byte) {
	m.read(addr, data)
}

// 0xeb

type mappereb struct {
	baseMapper
}

func newMappereb(bm *baseMapper) Mapper {
	return &mappereb{baseMapper: *bm}
}

func (m *mappereb) reset() {
	for i := 0; i < 0x2000; i++ {
		m.mem.dram[i] = 0xff
	}
	m.mem.setProm32kBank(0)
}

func (m *mappereb) write(addr uint16, data byte) {
	i := uint32(((addr & 0x0300) >> 3) | (addr & 0x001f))
	b := false
	switch m.mem.nProm8kPage {
	case 128:
		if addr&0x0300 != 0 {
			b = true
		}
	case 256:
		switch addr & 0x0300 {
		case 0x0100, 0x0300:
			b = true
		case 0x0200:
			i = (i & 0x1f) | 0x20
		}
	case 384:
		switch addr & 0x0300 {
		case 0x0100:
			b = true
		case 0x0200:
			i = (i & 0x1f) | 0x20
		case 0x0300:
			i = (i & 0x1f) | 0x40
		}
	}

	if addr&0x0800 != 0 {
		j := i << 2
		if addr&0x1000 != 0 {
			j += 2
		}
		m.mem.setProm32kBank4(j, j+1, j, j+1)
	} else {
		m.mem.setProm32kBank(i)
	}
	if b {
		for i := byte(4); i < 8; i++ {
			m.mem.setCpuBank(i, m.mem.dram[:], memBankTypRom)
		}
	}

	if addr&0x0400 != 0 {
		m.mem.setVramMirror(memVramMirror4L)
	} else if addr&0x2000 != 0 {
		m.mem.setVramMirror(memVramMirrorH)
	} else {
		m.mem.setVramMirror(memVramMirrorV)
	}
}

// 0xec

type mapperec struct {
	baseMapper
	bank, mode byte
}

func newMapperec(bm *baseMapper) Mapper {
	return &mapperec{baseMapper: *bm}
}

func (m *mapperec) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
}

func (m *mapperec) write(addr uint16, data byte) {
	a8 := byte(addr)
	if addr >= 0x8000 && addr <= 0xbfff {
		m.bank = ((a8 & 0x03) << 4) | (m.bank & 0x07)
	} else {
		m.bank, m.mode = (a8&0x07)|(m.bank&0x30), a8&0x30
	}

	switch m.mode {
	case 0x00:
		m.bank |= 0x08
		b := uint32(m.bank)
		m.mem.setProm32kBank4(b<<1, (b<<1)+1, (b|0x07)<<1, ((b|0x07)<<1)+1)
	case 0x10:
		m.bank |= 0x37
		b := uint32(m.bank)
		m.mem.setProm32kBank4(b<<1, (b<<1)+1, (b|0x07)<<1, ((b|0x07)<<1)+1)
	case 0x20:
		m.bank |= 0x08
		b := uint32(m.bank) & 0xfe
		m.mem.setProm16kBank(4, b)
		m.mem.setProm16kBank(6, b+1)
	case 0x30:
		m.bank |= 0x08
		b := uint32(m.bank)
		m.mem.setProm16kBank(4, b)
		m.mem.setProm16kBank(6, b)
	}

	if a8&0x20 != 0 {
		m.mem.setVramMirror(memVramMirrorH)
	} else {
		m.mem.setVramMirror(memVramMirrorV)
	}
}

// 0xed

type mappered struct {
	baseMapper
}

func newMappered(bm *baseMapper) Mapper {
	return &mappered{baseMapper: *bm}
}

func (m *mappered) reset() {
}

// 0xee

type mapperee struct {
	baseMapper
}

func newMapperee(bm *baseMapper) Mapper {
	return &mapperee{baseMapper: *bm}
}

func (m *mapperee) reset() {
}

// 0xef

type mapperef struct {
	baseMapper
}

func newMapperef(bm *baseMapper) Mapper {
	return &mapperef{baseMapper: *bm}
}

func (m *mapperef) reset() {
}
