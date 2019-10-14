package core

// 240

type mapper240 struct {
	baseMapper
}

func newMapper240(bm *baseMapper) Mapper {
	return &mapper240{baseMapper: *bm}
}

func (m *mapper240) reset() {
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper240) writeLow(addr uint16, data byte) {
	if addr >= 0x4020 && addr < 0x6000 {
		m.mem.setProm32kBank(uint32(data&0xf0) >> 4)
		m.mem.setVrom8kBank(uint32(data) & 0x0f)
	}
}

// 241

type mapper241 struct {
	baseMapper
}

func newMapper241(bm *baseMapper) Mapper {
	return &mapper241{baseMapper: *bm}
}

func (m *mapper241) reset() {
	m.mem.setProm32kBank(0)
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper241) write(addr uint16, data byte) {
	if addr == 0x8000 {
		m.mem.setProm32kBank(uint32(data))
	}
}

// 242

type mapper242 struct {
	baseMapper
}

func newMapper242(bm *baseMapper) Mapper {
	return &mapper242{baseMapper: *bm}
}

func (m *mapper242) reset() {
	m.mem.setProm32kBank(0)
}

func (m *mapper242) write(addr uint16, data byte) {
	if addr&0x0001 != 0 {
		m.mem.setProm32kBank(uint32(addr&0xf8) >> 3)
	}
}

// 243

type mapper243 struct {
	baseMapper
	r [4]byte
}

func newMapper243(bm *baseMapper) Mapper {
	return &mapper243{baseMapper: *bm}
}

func (m *mapper243) reset() {
	m.mem.setProm32kBank(0)
	if m.nVrom1kPage>>3 > 4 {
		m.mem.setVrom8kBank(4)
	} else if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
	m.mem.setVramMirror(memVramMirrorH)
}

func (m *mapper243) write(addr uint16, data byte) {
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

// 244

type mapper244 struct {
	baseMapper
}

func newMapper244(bm *baseMapper) Mapper {
	return &mapper244{baseMapper: *bm}
}

func (m *mapper244) reset() {
	m.mem.setProm32kBank(0)
}

func (m *mapper244) write(addr uint16, data byte) {
	if addr >= 0x8065 && addr <= 0x80a4 {
		m.mem.setProm32kBank(uint32(addr-0x8065) & 0x03)
	}
	if addr >= 0x80a5 && addr <= 0x80e4 {
		m.mem.setVrom8kBank(uint32(addr-0x80a5) & 0x07)
	}
}

// 245

type mapper245 struct {
	baseMapper
	irqEn    bool
	irqReq   bool
	irqCnt   byte
	irqLatch byte
	r0, r1   byte
	p        [2]byte
	c        [8]byte
}

func newMapper245(bm *baseMapper) Mapper {
	return &mapper245{baseMapper: *bm}
}

func (m *mapper245) reset() {
	m.irqEn, m.irqReq, m.irqCnt, m.irqLatch = false, false, 0, 0
	m.r0, m.r1, m.p[0], m.p[1] = 0, 0, 0, 1
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper245) write(addr uint16, data byte) {
	switch addr & 0xf7ff {
	case 0x8000:
		m.r0 = data
	case 0x8001:
		switch m.r0 {
		case 0x00:
			m.r1 = (data & 0x02) << 5
			m.mem.setProm8kBank(6, uint32(0x3e|m.r1))
			m.mem.setProm8kBank(7, uint32(0x3f|m.r1))
		case 0x06, 0x07:
			m.p[m.r0&0x01] = data
		}
		m.mem.setProm8kBank(4, uint32(m.p[0]|m.r1))
		m.mem.setProm8kBank(5, uint32(m.p[1]|m.r1))
	case 0xa000:
		if !m.sys.rom.b4Screen {
			if data&0x01 != 0 {
				m.mem.setVramMirror(memVramMirrorH)
			} else {
				m.mem.setVramMirror(memVramMirrorV)
			}
		}
	case 0xc000:
		m.irqReq, m.irqCnt = false, data
		m.clearIntr()
	case 0xc001:
		m.irqReq, m.irqLatch = false, data
		m.clearIntr()
	case 0xe000:
		m.irqEn, m.irqReq = false, false
		m.clearIntr()
	case 0xe001:
		m.irqEn, m.irqReq = true, false
		m.clearIntr()
	}
}

func (m *mapper245) hSync(scanline uint16) {
	if scanline < ScreenHeight && m.isPpuDisp() && m.irqEn && !m.irqReq {
		if scanline == 0 && m.irqCnt != 0 {
			m.irqCnt--
		}
		m.irqCnt--
		if m.irqCnt == 0xff {
			m.irqReq, m.irqCnt = true, m.irqLatch
			m.setIntr()
		}
	}
}

// 246

type mapper246 struct {
	baseMapper
}

func newMapper246(bm *baseMapper) Mapper {
	return &mapper246{baseMapper: *bm}
}

func (m *mapper246) reset() {
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
}

func (m *mapper246) writeLow(addr uint16, data byte) {
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

// 248

type mapper248 struct {
	baseMapper
	irqEn    bool
	irqCnt   byte
	irqLatch byte
	r        byte
	p        [2]byte
	c        [8]byte
}

func newMapper248(bm *baseMapper) Mapper {
	return &mapper248{baseMapper: *bm}
}

func (m *mapper248) setCpuBanks() {
	if m.r&0x40 != 0 {
		m.mem.setProm32kBank4(m.nProm8kPage-2, uint32(m.p[1]), uint32(m.p[0]), m.nProm8kPage-1)
	} else {
		m.mem.setProm32kBank4(uint32(m.p[0]), uint32(m.p[1]), m.nProm8kPage-2, m.nProm8kPage-1)
	}
}

func (m *mapper248) setPpuBanks() {
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

func (m *mapper248) reset() {
	m.irqEn, m.irqCnt, m.irqLatch = false, 0, 0
	m.r, m.p[0], m.p[1] = 0, 0, 1
	m.setCpuBanks()
	m.setPpuBanks()
}

func (m *mapper248) writeLow(addr uint16, data byte) {
	m.mem.setProm32kBank4(uint32(data)<<1, (uint32(data)<<1)+1, m.nProm8kPage-2, m.nProm8kPage-1)
}

func (m *mapper248) write(addr uint16, data byte) {
	addr &= 0xe001
	switch addr {
	case 0x8000:
		m.r = data
		m.setCpuBanks()
		m.setPpuBanks()
	case 0x8001:
		r := m.r & 0x07
		switch r {
		case 0x00, 0x01:
			i := r << 1
			m.c[i] = data & 0xfe
			m.c[i+1] = m.c[i] + 1
			m.setPpuBanks()
		case 0x02, 0x03, 0x04, 0x05:
			m.c[r+2] = data
		case 0x06, 0x07:
			m.p[r&0x01] = data
			m.setPpuBanks()
		}
	case 0xa000:
		if !m.sys.rom.b4Screen {
			if data&0x01 != 0 {
				m.mem.setVramMirror(memVramMirrorH)
			} else {
				m.mem.setVramMirror(memVramMirrorV)
			}
		}
	case 0xc000:
		m.irqEn, m.irqCnt, m.irqLatch = false, 0xbe, 0xbe
		m.clearIntr()
	case 0xc001:
		m.irqEn, m.irqCnt, m.irqLatch = true, 0xbe, 0xbe
		m.clearIntr()
	}
}

func (m *mapper248) hSync(scanline uint16) {
	if scanline < ScreenHeight && m.isPpuDisp() && m.irqEn {
		m.irqCnt--
		if m.irqCnt == 0xff {
			m.irqCnt = m.irqLatch
			m.setIntr()
		}
	}
}

// 249

type mapper249 struct {
	baseMapper
	sp       bool
	irqEn    bool
	irqReq   bool
	irqCnt   byte
	irqLatch byte
	r        byte
}

func newMapper249(bm *baseMapper) Mapper {
	return &mapper249{baseMapper: *bm}
}

func (m *mapper249) reset() {
	m.sp = false
	m.irqEn, m.irqReq, m.irqCnt, m.irqLatch = false, false, 0, 0
	m.r = 0
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
	m.mem.setVrom8kBank(0)
}

func (m *mapper249) writeLow(addr uint16, data byte) {
	if addr == 0x5000 {
		switch data {
		case 0x00:
			m.sp = false
		case 0x02:
			m.sp = true
		}
	} else if addr >= 0x6000 && addr < 0x8000 {
		m.cpuBanks[addr>>13][addr&0x1fff] = data
	}
}

func (m *mapper249) alterData(data byte) uint32 {
	if !m.sp {
		return uint32(data)
	}
	var ds [8]uint32
	for i := byte(0); i < 8; i++ {
		ds[i] = (uint32(data) & (0x01 << i)) >> i
	}
	return (ds[5] << 7) | (ds[4] << 6) | (ds[2] << 5) | (ds[6] << 4) | (ds[7] << 3) | (ds[3] << 2) | (ds[1] << 1) | ds[0]
}

func (m *mapper249) alterData2(data byte) uint32 {
	if !m.sp {
		return uint32(data)
	} else if data >= 0x20 {
		return m.alterData(data - 0x20)
	}
	var ds [5]uint32
	for i := byte(0); i < 5; i++ {
		ds[i] = (uint32(data) & (0x01 << i)) >> i
	}
	return (ds[2] << 4) | (ds[1] << 3) | (ds[3] << 2) | (ds[4] << 1) | ds[0]
}

func (m *mapper249) write(addr uint16, data byte) {
	switch addr & 0xff01 {
	case 0x8000, 0x8800:
		m.r = data
	case 0x8001, 0x8801:
		r := m.r & 0x07
		switch r {
		case 0x00, 0x01:
			i := r << 1
			d := m.alterData(data)
			m.mem.setVrom1kBank(i, d&0xfe)
			m.mem.setVrom1kBank(i+1, d|0x01)
		case 0x02, 0x03, 0x04, 0x05:
			m.mem.setVrom1kBank(r+2, m.alterData(data))
		case 0x06, 0x07:
			m.mem.setVrom1kBank(r-2, m.alterData2(data))
		}
	case 0xa000, 0xa800:
		if !m.sys.rom.b4Screen {
			if data&0x01 != 0 {
				m.mem.setVramMirror(memVramMirrorH)
			} else {
				m.mem.setVramMirror(memVramMirrorV)
			}
		}
	case 0xc000, 0xc800:
		m.irqReq, m.irqCnt = false, data
		m.clearIntr()
	case 0xc001, 0xc801:
		m.irqReq, m.irqLatch = false, data
		m.clearIntr()
	case 0xe000, 0xe800:
		m.irqEn, m.irqReq = false, false
		m.clearIntr()
	case 0xe001, 0xe801:
		m.irqEn, m.irqReq = true, false
		m.clearIntr()
	}
}

func (m *mapper249) hSync(scanline uint16) {
	if scanline < ScreenHeight && m.isPpuDisp() && m.irqEn && !m.irqReq {
		if scanline == 0 && m.irqCnt != 0 {
			m.irqCnt--
		}
		m.irqCnt--
		if m.irqCnt == 0xff {
			m.irqReq, m.irqCnt = true, m.irqLatch
			m.setIntr()
		}
	}
}

// 251

type mapper251 struct {
	baseMapper
	r [11]byte
	b [4]byte
}

func newMapper251(bm *baseMapper) Mapper {
	return &mapper251{baseMapper: *bm}
}

func (m *mapper251) reset() {
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
	m.mem.setVramMirror(memVramMirrorV)
	for i := 0; i < 11; i++ {
		m.r[i] = 0
	}
	for i := 0; i < 4; i++ {
		m.b[i] = 0
	}
}

func (m *mapper251) writeLow(addr uint16, data byte) {
	if addr&0xe001 == 0x6000 && m.r[9] != 0 {
		m.b[m.r[10]] = data
		m.r[10]++
		if m.r[10] == 4 {
			m.r[10] = 0
			m.setBanks()
		}
	}
}

func (m *mapper251) write(addr uint16, data byte) {
	switch addr & 0xe001 {
	case 0x8000:
		m.r[8] = data
		m.setBanks()
	case 0x8001:
		m.r[m.r[8]&0x07] = data
		m.setBanks()
	case 0xa001:
		if data&0x80 != 0 {
			m.r[9], m.r[10] = 1, 0
		} else {
			m.r[9] = 0
		}
	}

}
func (m *mapper251) setBanks() {
	var c [6]uint32
	for i := 0; i < 6; i++ {
		c[i] = (uint32(m.r[i]) | (uint32(m.b[1]) << 4)) & ((uint32(m.b[2]) << 4) | 0x0f)
	}
	if m.r[8]&0x80 != 0 {
		m.mem.setVrom8kBank8(c[2], c[3], c[4], c[5], c[0], c[0]+1, c[1], c[1]+1)
	} else {
		m.mem.setVrom8kBank8(c[0], c[0]+1, c[1], c[1]+1, c[2], c[3], c[4], c[5])
	}

	p0 := uint32(m.r[6]&((m.b[3]&0x3f)^0x3f)) | uint32(m.b[1])
	p1 := uint32(m.r[7]&((m.b[3]&0x3f)^0x3f)) | uint32(m.b[1])
	p3 := uint32((m.b[3]&0x3f)^0x3f) | uint32(m.b[1])
	p2 := p3 & (m.nProm8kPage - 1)
	if m.r[8]&0x40 != 0 {
		m.mem.setProm32kBank4(p2, p1, p0, p3)
	} else {
		m.mem.setProm32kBank4(p0, p1, p2, p3)
	}
}

// 252

type mapper252 struct {
	baseMapper
	irqOccur bool
	irqEn    byte
	irqCnt   byte
	irqLatch byte
	irqClk   uint16
	r        [9]byte
}

func newMapper252(bm *baseMapper) Mapper {
	return &mapper252{baseMapper: *bm}
}

func (m *mapper252) reset() {
	m.irqOccur, m.irqEn, m.irqCnt, m.irqLatch, m.irqClk = false, 0, 0, 0, 0
	for i := byte(0); i < 8; i++ {
		m.r[i] = i
	}
	m.r[8] = 0
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
	m.mem.setVrom8kBank(0)
	m.sys.renderMode = RenderModePost
}

func (m *mapper252) write(addr uint16, data byte) {
	switch addr & 0xf000 {
	case 0x8000, 0xa000:
		m.mem.setProm8kBank(byte(addr>>13), uint32(data))
		return
	}
	addr &= 0xf00c
	switch addr {
	case 0xb000, 0xb008, 0xc000, 0xc008, 0xd000, 0xd008, 0xe000, 0xe008:
		i := (byte(addr>>11) - 0x16) | (byte(addr&0x08) >> 3)
		m.r[i] = (m.r[i] & 0xf0) | (data & 0x0f)
		m.mem.setVrom1kBank(i, uint32(m.r[i]))
	case 0xb004, 0xb00c, 0xc004, 0xc00c, 0xd004, 0xd00c, 0xe004, 0xe00c:
		i := (byte(addr>>11) - 0x16) | (byte(addr&0x08) >> 3)
		m.r[i] = (m.r[i] & 0x0f) | (data << 4)
		m.mem.setVrom1kBank(i, uint32(m.r[i]))
	case 0xf000:
		m.irqOccur, m.irqLatch = false, (m.irqLatch&0xf0)|(data&0x0f)
	case 0xf004:
		m.irqOccur, m.irqLatch = false, (m.irqLatch&0x0f)|(data<<4)
	case 0xf008:
		m.irqOccur, m.irqEn = false, data&0x03
		if m.irqEn&0x02 != 0 {
			m.irqCnt, m.irqClk = m.irqLatch, 0
		}
		m.clearIntr()
	case 0xf00c:
		m.irqOccur, m.irqEn = false, (m.irqEn&0x01)*3
		m.clearIntr()
	}
}

func (m *mapper252) clock(nCycle int64) {
	if m.irqEn&0x02 != 0 {
		m.irqClk += uint16(nCycle)
		if m.irqClk >= 114 {
			m.irqClk -= 114
			if m.irqCnt == 0xFF {
				m.irqOccur = true
				m.irqEn, m.irqCnt = (m.irqEn&0x01)*3, m.irqLatch
				m.setIntr()
			} else {
				m.irqCnt++
			}
		}
	}
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
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
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

// 254

type mapper254 struct {
	baseMapper
	prot     bool
	irqEn    bool
	irqReq   bool
	irqCnt   byte
	irqLatch byte
	r        byte
	p        [2]byte
	c        [8]byte
}

func newMapper254(bm *baseMapper) Mapper {
	return &mapper254{baseMapper: *bm}
}

func (m *mapper254) setCpuBanks() {
	if m.r&0x40 != 0 {
		m.mem.setProm32kBank4(m.nProm8kPage-2, uint32(m.p[1]), uint32(m.p[0]), m.nProm8kPage-1)
	} else {
		m.mem.setProm32kBank4(uint32(m.p[0]), uint32(m.p[1]), m.nProm8kPage-2, m.nProm8kPage-1)
	}
}

func (m *mapper254) setPpuBanks() {
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
	} else {
		if m.r&0x80 != 0 {
			for i := byte(0); i < 4; i++ {
				m.mem.setCram1kBank(i, uint32(m.c[i+4]&0x07))
				m.mem.setCram1kBank(i+4, uint32(m.c[i]&0x07))
			}
		} else {
			for i := byte(0); i < 8; i++ {
				m.mem.setCram1kBank(i, uint32(m.c[i]&0x07))
			}
		}
	}
}

func (m *mapper254) reset() {
	m.prot = false
	m.irqEn, m.irqReq, m.irqCnt, m.irqLatch = false, false, 0, 0
	m.r, m.p[0], m.p[1] = 0, 0, 1
	m.setCpuBanks()
	m.setPpuBanks()
}

func (m *mapper254) readLow(addr uint16) byte {
	if addr >= 0x6000 {
		b := m.cpuBanks[addr>>13][addr&0x1fff]
		if !m.prot {
			b ^= 0x01
		}
		return b
	}
	return byte(addr >> 8)
}

func (m *mapper254) writeLow(addr uint16, data byte) {
	switch addr & 0xf000 {
	case 0x6000, 0x7000:
		m.cpuBanks[addr>>13][addr&0x1fff] = data
	}
}

func (m *mapper254) write(addr uint16, data byte) {
	addr &= 0xe001
	switch addr {
	case 0x8000:
		m.prot, m.r = true, data
		m.setCpuBanks()
		m.setPpuBanks()
	case 0x8001:
		r := m.r & 0x07
		switch r {
		case 0x00, 0x01:
			i := r << 1
			m.c[i] = data & 0xfe
			m.c[i+1] = m.c[i] + 1
			m.setPpuBanks()
		case 0x02, 0x03, 0x04, 0x05:
			m.c[r+2] = data
		case 0x06, 0x07:
			m.p[r&0x01] = data
			m.setPpuBanks()
		}
	case 0xa000:
		if !m.sys.rom.b4Screen {
			if data&0x01 != 0 {
				m.mem.setVramMirror(memVramMirrorH)
			} else {
				m.mem.setVramMirror(memVramMirrorV)
			}
		}
	case 0xc000:
		m.irqReq, m.irqCnt = false, data
		m.clearIntr()
	case 0xc001:
		m.irqReq, m.irqLatch = false, data
		m.clearIntr()
	case 0xe000:
		m.irqEn, m.irqReq = false, false
		m.clearIntr()
	case 0xe001:
		m.irqEn, m.irqReq = true, false
		m.clearIntr()
	}
}

func (m *mapper254) hSync(scanline uint16) {
	if scanline < ScreenHeight && m.isPpuDisp() && m.irqEn && !m.irqReq {
		if scanline == 0 && m.irqCnt != 0 {
			m.irqCnt--
		}
		m.irqCnt--
		if m.irqCnt == 0xff {
			m.irqReq, m.irqCnt = true, m.irqLatch
			m.setIntr()
		}
	}
}

// 255

type mapper255 struct {
	baseMapper
	r [4]byte
}

func newMapper255(bm *baseMapper) Mapper {
	return &mapper255{baseMapper: *bm}
}

func (m *mapper255) reset() {
	m.r[0], m.r[1], m.r[2], m.r[3] = 0, 0, 0, 0
	m.mem.setProm32kBank(0)
	m.mem.setVrom8kBank(0)
	m.mem.setVramMirror(memVramMirrorV)
}

func (m *mapper255) readLow(addr uint16) byte {
	if addr >= 0x5800 {
		return m.r[addr&0x03] & 0x0f
	}
	return byte(addr >> 8)
}

func (m *mapper255) writeLow(addr uint16, data byte) {
	if addr >= 0x5800 {
		m.r[addr&0x03] = data & 0x0f
	}
}

func (m *mapper255) write(addr uint16, data byte) {
	b := uint32(addr&0x4000) >> 7
	p := b + uint32(addr&0x0f80)>>5
	if addr&0x1000 != 0 {
		if addr&0x0040 != 0 {
			m.mem.setProm32kBank4(p+2, p+3, p+2, p+3)
		} else {
			m.mem.setProm32kBank4(p, p+1, p, p+1)
		}
	} else {
		m.mem.setProm32kBank4(p, p+1, p+2, p+3)
	}

	b = (b << 2) + (uint32(addr&0x003f) << 3)
	for i := byte(0); i < 8; i++ {
		m.mem.setVrom1kBank(i, b+uint32(i))
	}
	if addr&0x2000 != 0 {
		m.mem.setVramMirror(memVramMirrorH)
	} else {
		m.mem.setVramMirror(memVramMirrorV)
	}
}
