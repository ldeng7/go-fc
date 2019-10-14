package core

// 080

type mapper080 struct {
	baseMapper
	r byte
}

func newMapper080(bm *baseMapper) Mapper {
	return &mapper080{baseMapper: *bm}
}

func (m *mapper080) reset() {
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper080) writeLow(addr uint16, data byte) {
	switch addr {
	case 0x7ef0, 0x7ef1:
		i := (byte(addr) & 0x01) << 1
		m.mem.setVrom2kBank(i, uint32(data>>1)&0x3f)
		if m.nProm8kPage == 32 {
			b := uint32(data >> 7)
			m.mem.setVram1kBank(8+i, b)
			m.mem.setVram1kBank(9+i, b)
		}
	case 0x7EF2, 0x7EF3, 0x7EF4, 0x7EF5:
		i := (byte(addr) & 0x0f) + 2
		m.mem.setVrom1kBank(i, uint32(data))
	case 0x7ef6:
		if data&0x01 != 0 {
			m.mem.setVramMirror(memVramMirrorV)
		} else {
			m.mem.setVramMirror(memVramMirrorH)
		}
	case 0x7efa, 0x7efb, 0x7efc, 0x7efd, 0x7efe, 0x7eff:
		i := (byte(addr) & 0x06) + 3
		m.mem.setProm8kBank(i, uint32(data))
	default:
		m.baseMapper.writeLow(addr, data)
	}
}

// 082

type mapper082 struct {
	baseMapper
	r byte
}

func newMapper082(bm *baseMapper) Mapper {
	return &mapper082{baseMapper: *bm}
}

func (m *mapper082) reset() {
	m.r = 0
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
	m.mem.setVramMirror(memVramMirrorV)
}

func (m *mapper082) writeLow(addr uint16, data byte) {
	switch addr {
	case 0x7ef0, 0x7ef1:
		i := (byte(addr) & 0x01) << 1
		if m.r != 0 {
			i += 4
		}
		m.mem.setVrom2kBank(i, uint32(data>>1))
	case 0x7EF2, 0x7EF3, 0x7EF4, 0x7EF5:
		i := (byte(addr) & 0x0f) - 2
		if m.r == 0 {
			i += 4
		}
		m.mem.setVrom1kBank(i, uint32(data))
	case 0x7ef6:
		m.r = data & 0x02
		if data&0x01 != 0 {
			m.mem.setVramMirror(memVramMirrorV)
		} else {
			m.mem.setVramMirror(memVramMirrorH)
		}
	case 0x7efa, 0x7efb, 0x7efc:
		i := (byte(addr) & 0x07) + 2
		m.mem.setProm8kBank(i, uint32(data>>2))
	default:
		m.baseMapper.writeLow(addr, data)
	}
}

// 083

type mapper083 struct {
	baseMapper
	patchTyp   bool
	irqEn      bool
	irqCnt     uint16
	chrBank    uint32
	r0, r1, r2 byte
}

func newMapper083(bm *baseMapper) Mapper {
	return &mapper083{baseMapper: *bm}
}

func (m *mapper083) reset() {
	if m.sys.conf.PatchTyp&0x01 != 0 {
		m.patchTyp = true
	}
	m.irqEn, m.irqCnt, m.chrBank = false, 0, 0
	m.r0, m.r1, m.r2 = 0, 0, 0
	if m.nProm8kPage >= 32 {
		m.mem.setProm32kBank4(0, 1, 30, 31)
		m.r1 = 0x30
	} else {
		m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
	}
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper083) readLow(addr uint16) byte {
	if addr&0x5100 == 0x5100 {
		return m.r2
	}
	return m.baseMapper.readLow(addr)
}

func (m *mapper083) writeLow(addr uint16, data byte) {
	switch addr {
	case 0x5101, 0x5102, 0x5103:
		m.r2 = data
	}
	m.baseMapper.writeLow(addr, data)
}

func (m *mapper083) write(addr uint16, data byte) {
	switch addr {
	case 0x8000, 0xb000, 0xb0ff, 0xb1ff:
		m.chrBank, m.r0 = uint32(data&0x30)<<4, data
		m.mem.setProm16kBank(4, uint32(data))
		m.mem.setProm16kBank(6, uint32(data&0x30)|0x0f)
	case 0x8100:
		m.r1 = data & 0x80
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
	case 0x8200:
		m.irqCnt = (m.irqCnt & 0xff00) | uint16(data)
		break
	case 0x8201:
		m.irqCnt = (m.irqCnt & 0x00ff) | (uint16(data) << 8)
		m.irqEn = m.r1 != 0
	case 0x8300, 0x8301, 0x8302:
		m.mem.setProm8kBank(byte(addr&0x03)+4, uint32(data))
	case 0x8310, 0x8311:
		i, b := byte(addr&0x01), m.chrBank|uint32(data)
		if m.patchTyp {
			m.mem.setVrom2kBank(i<<1, b)
		} else {
			m.mem.setVrom1kBank(i, b)
		}
	case 0x8312, 0x8313, 0x8314, 0x8315:
		m.mem.setVrom1kBank(byte(addr&0x0f), m.chrBank|uint32(data))
	case 0x8316, 0x8317:
		i, b := byte(addr&0x0f), m.chrBank|uint32(data)
		if m.patchTyp {
			m.mem.setVrom2kBank((i-1)&0x06, b)
		} else {
			m.mem.setVrom1kBank(i, b)
		}
	case 0x8318:
		m.mem.setProm16kBank(4, uint32((m.r0&0x30)|data))
	}
}

func (m *mapper083) hSync(scanline uint16) {
	if m.irqEn {
		if m.irqCnt <= 113 {
			m.irqEn = false
			m.setIntr()
		} else {
			m.irqCnt -= 113
		}
	}
}

// 085

type mapper085 struct {
	baseMapper
	irqEn    byte
	irqCnt   byte
	irqLatch byte
	irqClk   uint16
}

func newMapper085(bm *baseMapper) Mapper {
	return &mapper085{baseMapper: *bm}
}

func (m *mapper085) reset() {
	m.irqEn, m.irqCnt, m.irqLatch, m.irqClk = 0, 0, 0, 0
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	} else {
		m.mem.setCram8kBank(0)
	}
}

func (m *mapper085) write(addr uint16, data byte) {
	addr &= 0xf038
	switch addr {
	case 0x8000:
		m.mem.setProm8kBank(4, uint32(data))
	case 0x8008, 0x8010:
		m.mem.setProm8kBank(5, uint32(data))
	case 0x9000:
		m.mem.setProm8kBank(6, uint32(data))
	case 0xa000, 0xa008, 0xa010, 0xb000, 0xb008, 0xb010,
		0xc000, 0xc008, 0xc010, 0xd000, 0xd008, 0xd010:
		i := (byte(addr>>11) - 0x14) | (byte(addr&0x08) >> 3) | (byte(addr&0x10) >> 4)
		if m.nVrom1kPage != 0 {
			m.mem.setVrom1kBank(i, uint32(data))
		} else {
			m.mem.setCram1kBank(i, uint32(data))
		}
	case 0xe000:
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
	case 0xe008, 0xe010:
		m.irqLatch = data
	case 0xf000:
		m.irqEn, m.irqCnt, m.irqClk = data&0x03, m.irqLatch, 0
		m.clearIntr()
	case 0xf008, 0xf010:
		m.irqEn = (m.irqEn & 0x01) * 3
		m.clearIntr()
	}
}

func (m *mapper085) clock(nCycle int64) {
	if m.irqEn&0x02 != 0 {
		m.irqClk += uint16(nCycle * 4)
		for m.irqClk >= 455 {
			m.irqClk -= 455
			m.irqCnt++
			if m.irqCnt == 0 {
				m.irqCnt = m.irqLatch
				m.setIntr()
			}
		}
	}
}

// 086

type mapper086 struct {
	baseMapper
	r, c byte
}

func newMapper086(bm *baseMapper) Mapper {
	return &mapper086{baseMapper: *bm}
}

func (m *mapper086) reset() {
	m.mem.setProm32kBank(0)
	m.mem.setVrom8kBank(0)
	m.r, m.c = 0xff, 0
}

func (m *mapper086) write(addr uint16, data byte) {
	switch addr {
	case 0x6000:
		m.mem.setProm32kBank(uint32(data&0x30) >> 4)
		m.mem.setVrom8kBank(uint32(data&0x03) | (uint32(data&0x40) >> 4))
	case 0x7000:
		if m.r&0x10 == 0 && data&0x10 != 0 && m.c == 0 && (data&0x0f == 0 || data&0x0f == 0x05) {
			m.c = 60
		}
		m.r = data
	}
}

func (m *mapper086) vSync() {
	if m.c != 0 {
		m.c--
	}
}

// 087

type mapper087 struct {
	baseMapper
}

func newMapper087(bm *baseMapper) Mapper {
	return &mapper087{baseMapper: *bm}
}

func (m *mapper087) reset() {
	m.mem.setProm32kBank(0)
	m.mem.setVrom8kBank(0)
}

func (m *mapper087) write(addr uint16, data byte) {
	if addr == 0x6000 {
		m.mem.setVrom8kBank(uint32(data&0x02) >> 1)
	}
}

// 088

type mapper088 struct {
	baseMapper
	r byte
}

func newMapper088(bm *baseMapper) Mapper {
	return &mapper088{baseMapper: *bm}
}

func (m *mapper088) reset() {
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
	if m.nVrom1kPage >= 8 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper088) write(addr uint16, data byte) {
	switch addr {
	case 0x8000:
		m.r = data
	case 0x8001:
		b, r := uint32(data), m.r&0x07
		switch r {
		case 0x00, 0x01:
			m.mem.setVrom2kBank(r<<1, b>>1)
		case 0x02, 0x03, 0x04, 0x05:
			m.mem.setVrom1kBank(r+2, b+0x40)
		case 0x06, 0x07:
			m.mem.setProm8kBank(r+2, b)
		}
	case 0xc000:
		if data != 0 {
			m.mem.setVramMirror(memVramMirror4H)
		} else {
			m.mem.setVramMirror(memVramMirror4L)
		}
	}
}

// 089

type mapper089 struct {
	baseMapper
}

func newMapper089(bm *baseMapper) Mapper {
	return &mapper089{baseMapper: *bm}
}

func (m *mapper089) reset() {
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
	m.mem.setVrom8kBank(0)
}

func (m *mapper089) write(addr uint16, data byte) {
	if (addr & 0xff00) == 0xc000 {
		m.mem.setProm16kBank(4, uint32(data&0x70)>>4)
		m.mem.setVrom8kBank((uint32(data&0x80) >> 4) | uint32(data&0x07))
		if data&0x08 != 0 {
			m.mem.setVramMirror(memVramMirror4H)
		} else {
			m.mem.setVramMirror(memVramMirror4L)
		}
	}
}

// 090

type mapper090 struct {
	baseMapper
	patchTyp bool

	irqEn    bool
	irqPre   bool
	irqCnt   byte
	irqLatch byte
	irqOfs   byte

	mirMode bool
	mirTyp  byte
	flags   byte
	k, sw   byte
	m0, m1  byte

	p, rl, rh [4]byte
	cl, ch    [8]byte
}

func newMapper090(bm *baseMapper) Mapper {
	return &mapper090{baseMapper: *bm}
}

func (m *mapper090) reset() {
	patch := m.sys.conf.PatchTyp
	if patch&0x01 != 0 {
		m.patchTyp = true
	}
	if patch&0x02 != 0 {
		m.sys.renderMode = RenderModeTile
	}

	m.irqEn, m.irqPre, m.irqCnt, m.irqLatch, m.irqOfs = false, false, 0, 0, 0
	m.mirMode, m.mirTyp, m.flags = false, 0, 0
	m.k, m.sw, m.m0, m.m1 = 0, ^m.sw, 0, 0

	for i := byte(0); i < 4; i++ {
		m.p[i] = byte(m.nProm8kPage) - 4 + i
		m.rl[i], m.rh[i] = 0, 0
		m.cl[i], m.ch[i] = i, 0
		m.cl[i+4], m.ch[i+4] = i+4, 0
	}

	b := m.nProm8kPage
	m.mem.setProm32kBank4(b-4, b-3, b-2, b-1)
	m.mem.setVrom8kBank(0)
}

func (m *mapper090) readLow(addr uint16) byte {
	switch addr {
	case 0x5000:
		return ^m.sw
	case 0x5800:
		return m.m0 * m.m1
	case 0x5801:
		return byte((uint16(m.m0) * uint16(m.m1)) >> 8)
	case 0x5803:
		return m.k
	}
	return m.baseMapper.readLow(addr)
}

func (m *mapper090) writeLow(addr uint16, data byte) {
	switch addr {
	case 0x5800:
		m.m0 = data
	case 0x5801:
		m.m1 = data
	case 0x5803:
		m.k = data
	}
}

func (m *mapper090) setCpuBanks() {
	switch m.flags & 0x03 {
	case 0x00:
		m.mem.setProm32kBank4(m.nProm8kPage-4, m.nProm8kPage-3, m.nProm8kPage-2, m.nProm8kPage-1)
	case 0x01:
		m.mem.setProm32kBank4(uint32(m.p[1])<<1, (uint32(m.p[1])<<1)+1, m.nProm8kPage-2, m.nProm8kPage-1)
	case 0x02:
		if m.flags&0x04 != 0 {
			m.mem.setProm32kBank4(uint32(m.p[0]), uint32(m.p[1]), uint32(m.p[2]), uint32(m.p[3]))
		} else {
			if m.flags&0x80 != 0 {
				m.mem.setProm8kBank(3, uint32(m.p[3]))
			}
			m.mem.setProm32kBank4(uint32(m.p[0]), uint32(m.p[1]), uint32(m.p[2]), m.nProm8kPage-1)
		}
	case 0x03:
		m.mem.setProm32kBank4(uint32(m.p[3]), uint32(m.p[2]), uint32(m.p[1]), uint32(m.p[0]))
	}
}

func (m *mapper090) setPpuBanks() {
	var b [8]uint32
	for i := 0; i < 8; i++ {
		b[i] = (uint32(m.ch[i]) << 8) | uint32(m.cl[i])
	}
	switch (m.flags & 0x18) >> 3 {
	case 0x00:
		m.mem.setVrom8kBank(uint32(b[0]))
	case 0x01:
		m.mem.setVrom4kBank(0, b[0])
		m.mem.setVrom4kBank(4, b[4])
	case 0x02:
		m.mem.setVrom2kBank(0, b[0])
		m.mem.setVrom2kBank(2, b[2])
		m.mem.setVrom2kBank(4, b[4])
		m.mem.setVrom2kBank(6, b[6])
	case 0x03:
		m.mem.setVrom8kBank8(b[0], b[1], b[2], b[3], b[4], b[5], b[6], b[7])
	}
}

func (m *mapper090) setPpuBanksAlt() {
	var b [4]uint32
	for i := 0; i < 4; i++ {
		b[i] = (uint32(m.rh[i]) << 8) | uint32(m.rl[i])
	}
	if !m.patchTyp && m.mirMode {
		for i := byte(0); i < 4; i++ {
			if m.rh[i] == 0 && m.rl[i] == i {
				m.mirMode = false
			}
		}
		if m.mirMode {
			m.mem.setVrom1kBank(8, uint32(b[0]))
			m.mem.setVrom1kBank(9, uint32(b[1]))
			m.mem.setVrom1kBank(10, uint32(b[2]))
			m.mem.setVrom1kBank(11, uint32(b[3]))
		}
	} else {
		switch m.mirTyp {
		case 0x00:
			m.mem.setVramMirror(memVramMirrorV)
		case 0x01:
			m.mem.setVramMirror(memVramMirrorH)
		default:
			m.mem.setVramMirror(memVramMirror4L)
		}
	}
}

func (m *mapper090) write(addr uint16, data byte) {
	switch addr & 0xf007 {
	case 0x8000, 0x8001, 0x8002, 0x8003:
		m.p[addr&0x03] = data
		m.setCpuBanks()
	case 0x9000, 0x9001, 0x9002, 0x9003, 0x9004, 0x9005, 0x9006, 0x9007:
		m.cl[addr&0x07] = data
		m.setPpuBanks()
	case 0xa000, 0xa001, 0xa002, 0xa003, 0xa004, 0xa005, 0xa006, 0xa007:
		m.ch[addr&0x07] = data
		m.setPpuBanks()
	case 0xb000, 0xb001, 0xb002, 0xb003:
		m.rl[addr&0x03] = data
		m.setPpuBanksAlt()
	case 0xb004, 0xb005, 0xb006, 0xb007:
		m.rh[addr&0x03] = data
		m.setPpuBanksAlt()
	case 0xc002:
		m.irqEn = false
		m.clearIntr()
	case 0xc003:
		m.irqEn, m.irqPre = true, true
	case 0xc005:
		if m.irqOfs&0x80 != 0 {
			m.irqLatch = data ^ (m.irqOfs | 0x01)
		} else {
			m.irqLatch = data | (m.irqOfs & 0x27)
		}
		m.irqPre = true
	case 0xC006:
		if m.patchTyp {
			m.irqOfs = data
		}
	case 0xd000:
		m.mirMode, m.flags = data&0x20 != 0, data
		m.setCpuBanks()
		m.setPpuBanks()
		m.setPpuBanksAlt()
	case 0xD001:
		m.mirTyp = data & 0x03
		m.setPpuBanksAlt()
	}
}

func (m *mapper090) hSync(scanline uint16) {
	if scanline < ScreenHeight && m.isPpuDisp() {
		if m.irqPre {
			m.irqPre, m.irqCnt = false, m.irqLatch
		}
		if m.irqCnt != 0 {
			m.irqCnt--
		}
		if m.irqCnt == 0 && m.irqEn {
			m.setIntr()
		}
	}
}

// 091

type mapper091 struct {
	baseMapper
	irqEn  bool
	irqCnt byte
}

func newMapper091(bm *baseMapper) Mapper {
	return &mapper091{baseMapper: *bm}
}

func (m *mapper091) reset() {
	m.irqEn, m.irqCnt = false, 0
	b := m.nProm8kPage
	m.mem.setProm32kBank4(b-2, b-1, b-2, b-1)
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank8(0, 0, 0, 0, 0, 0, 0, 0)
	}
	m.sys.renderMode = RenderModePostAll
}

func (m *mapper091) writeLow(addr uint16, data byte) {
	addr &= 0xf003
	switch addr {
	case 0x6000, 0x6001, 0x6002, 0x6003:
		m.mem.setVrom2kBank(byte(addr&0x03)<<1, uint32(data))
	case 0x7000, 0x7001:
		m.mem.setProm8kBank(byte(addr&0x01)+4, uint32(data))
	case 0x7002:
		m.irqEn, m.irqCnt = false, 0
		m.clearIntr()
	case 0x7003:
		m.irqEn = true
	}
}

func (m *mapper091) hSync(scanline uint16) {
	if scanline < ScreenHeight && m.isPpuDisp() {
		if m.irqEn {
			m.irqCnt++
		}
		if m.irqCnt >= 8 {
			m.setIntr()
		}
	}
}

// 092

type mapper092 struct {
	baseMapper
}

func newMapper092(bm *baseMapper) Mapper {
	return &mapper092{baseMapper: *bm}
}

func (m *mapper092) reset() {
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper092) write(addr uint16, data byte) {
	data = byte(addr)
	if addr >= 0x9000 {
		if (data & 0xf0) == 0xd0 {
			m.mem.setProm16kBank(6, uint32(data&0x0f))
		} else if (data & 0xf0) == 0xe0 {
			m.mem.setVrom8kBank(uint32(data & 0x0f))
		}
	} else {
		if (data & 0xf0) == 0xb0 {
			m.mem.setProm16kBank(6, uint32(data&0x0f))
		} else if (data & 0xf0) == 0x70 {
			m.mem.setVrom8kBank(uint32(data & 0x0f))
		}
	}
}

// 093

type mapper093 struct {
	baseMapper
}

func newMapper093(bm *baseMapper) Mapper {
	return &mapper093{baseMapper: *bm}
}

func (m *mapper093) reset() {
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper093) write(addr uint16, data byte) {
	if addr == 0x6000 {
		m.mem.setProm16kBank(4, uint32(data))
	}
}

// 094

type mapper094 struct {
	baseMapper
}

func newMapper094(bm *baseMapper) Mapper {
	return &mapper094{baseMapper: *bm}
}

func (m *mapper094) reset() {
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
}

func (m *mapper094) write(addr uint16, data byte) {
	if (addr & 0xfff0) == 0xff00 {
		m.mem.setProm16kBank(4, uint32(data>>2)&0x07)
	}
}

// 095

type mapper095 struct {
	baseMapper
	r byte
	p [2]byte
	c [8]byte
}

func newMapper095(bm *baseMapper) Mapper {
	return &mapper095{baseMapper: *bm}
}

func (m *mapper095) setCpuBanks() {
	if m.r&0x40 != 0 {
		m.mem.setProm32kBank4(m.nProm8kPage-2, uint32(m.p[1]), uint32(m.p[0]), m.nProm8kPage-1)
	} else {
		m.mem.setProm32kBank4(uint32(m.p[0]), uint32(m.p[1]), m.nProm8kPage-2, m.nProm8kPage-1)
	}
}

func (m *mapper095) setPpuBanks() {
	if m.nVrom1kPage != 0 {
		if m.r&0x80 != 0 {
			for i := byte(0); i < 4; i++ {
				m.mem.setVrom1kBank(i, uint32(m.c[i+4]))
				m.mem.setVrom1kBank(i+4, uint32(m.c[i]))
			}
		} else {
			for i := byte(0); i < 8; i++ {
				m.mem.setVrom1kBank(i, uint32(m.c[i]))
			}
		}
	}
}

func (m *mapper095) reset() {
	m.r, m.p[0], m.p[1] = 0, 0, 1
	if m.nVrom1kPage != 0 {
		for i := byte(0); i < 8; i++ {
			m.c[i] = i
		}
	} else {
		for i := byte(0); i < 8; i++ {
			m.c[i] = 0
		}
		m.c[1], m.c[3] = 1, 1
	}
	m.setCpuBanks()
	m.setPpuBanks()
	m.sys.renderMode = RenderModePost
}

func (m *mapper095) write(addr uint16, data byte) {
	switch addr & 0xe001 {
	case 0x8000:
		m.r = data
		m.setCpuBanks()
		m.setPpuBanks()
	case 0x8001:
		if m.r <= 0x05 {
			if data&0x20 != 0 {
				m.mem.setVramMirror(memVramMirror4H)
			} else {
				m.mem.setVramMirror(memVramMirror4L)
			}
			data &= 0x1f
		}
		r := m.r & 0x07
		switch r {
		case 0x00, 0x01:
			if m.nVrom1kPage != 0 {
				i := r << 1
				m.c[i] = data & 0xfe
				m.c[i+1] = m.c[i] + 1
				m.setPpuBanks()
			}
		case 0x02, 0x03, 0x04, 0x05:
			if m.nVrom1kPage != 0 {
				m.c[r+2] = data
				m.setPpuBanks()
			}
		case 0x06, 0x07:
			m.p[r&0x01] = data
			m.setPpuBanks()
		}
	}
}
