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

func (m *mapperf6) writeLow(addr uint16, data byte) {
	if addr >= 0x6000 {
		switch addr {
		case 0x6000, 0x6001, 0x6002, 0x6003:
			m.mem.setProm8kBank(byte(addr)+4, uint32(data))
		case 0x6004, 0x6005, 0x6006, 0x6007:
			m.mem.setVrom2kBank((byte(addr)-4)<<1, uint32(data))
		default:
			m.cpuBanks[addr>>13][addr&0x1fff] = data
		}
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

// 253

type mapper253 struct {
	baseMapper
	patch    bool
	vrsw     bool
	irqEn    byte
	irqCnt   byte
	irqLatch byte
	irqClk   uint16
	r        [8]byte
}

func newMapper253(bm *baseMapper) Mapper {
	return &mapper253{baseMapper: *bm}
}

func (m *mapper253) reset() {
	if m.sys.conf.PatchTyp&0x01 != 0 {
		m.patch = true
	}
	m.vrsw = false
	m.irqEn, m.irqCnt, m.irqLatch, m.irqClk = 0, 0, 0, 0
	for i := byte(0); i < 8; i++ {
		m.r[i] = i
	}
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	m.mem.setVrom8kBank(0)
}

func (m *mapper253) write(addr uint16, data byte) {
	switch addr {
	case 0x8010:
		m.mem.setProm8kBank(4, uint32(data))
	case 0xa010:
		m.mem.setProm8kBank(5, uint32(data))
	case 0x9400:
		switch data & 0x03 {
		case 0x00:
			m.mem.setVramMirror(memVramMirrorV)
		case 0x01:
			m.mem.setVramMirror(memVramMirrorH)
		case 0x02:
			m.mem.setVramMirror(memVramMirror4L)
		case 0x03:
			m.mem.setVramMirror(memVramMirror4H)
		}
	}
	a := addr & 0xf00c
	switch a {
	case 0xb000, 0xb008, 0xc000, 0xc008, 0xd000, 0xd008, 0xe000, 0xe008:
		i := (byte(addr>>11) - 0x16) | (byte(addr&0x08) >> 3)
		b := (m.r[i] & 0xf0) | (data & 0x0f)
		m.r[i] = b
		m.setPpuBanks(i, uint32(b))
	case 0xb004, 0xb00c, 0xc004, 0xc00c, 0xd004, 0xd00c, 0xe004, 0xe00c:
		i := (byte(addr>>11) - 0x16) | (byte(addr&0x08) >> 3)
		b := (m.r[i] & 0x0f) | ((data & 0x0f) << 4)
		m.r[i] = b
		m.setPpuBanks(i, uint32(b)|(uint32(data>>4)<<8))
	case 0xf000:
		m.irqLatch = (m.irqLatch & 0xf0) | (data & 0x0f)
	case 0xf004:
		m.irqLatch = (m.irqLatch & 0x0f) | ((data & 0x0f) << 4)
	case 0xf008:
		m.irqEn = data & 0x03
		if m.irqEn&0x02 != 0 {
			m.irqCnt, m.irqClk = m.irqLatch, 0
		}
		m.clearIntr()
	}
}

func (m *mapper253) setPpuBanks(iBank byte, iPage uint32) {
	if m.patch && (iPage == 0x88 || iPage == 0xc8) {
		m.vrsw = iPage&0x40 != 0
		return
	}
	if iPage == 4 || iPage == 5 {
		if m.patch && !m.vrsw {
			m.mem.setVrom1kBank(iBank, iPage)
		} else {
			m.mem.setCram1kBank(iBank, iPage)
		}
	} else {
		m.mem.setVrom1kBank(iBank, iPage)
	}
}

func (m *mapper253) clock(nCycle int64) {
	if m.irqEn&0x02 != 0 {
		m.irqClk += uint16(nCycle)
		if m.irqClk >= 114 {
			m.irqClk -= 114
			if m.irqCnt == 255 {
				m.irqCnt = m.irqLatch
				m.irqEn = (m.irqEn & 0x01) * 0x03
				m.setIntr()
			} else {
				m.irqCnt++
			}
		}
	}
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
