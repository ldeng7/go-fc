package core

// 160

type mapper160 struct {
	baseMapper
	irqEn    bool
	irqCnt   byte
	irqLatch byte
	rTyp     byte
}

func newMapper160(bm *baseMapper) Mapper {
	return &mapper160{baseMapper: *bm}
}

func (m *mapper160) reset() {
	m.irqEn, m.irqCnt, m.irqLatch, m.rTyp = false, 0, 0, 0
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
}

func (m *mapper160) write(addr uint16, data byte) {
	switch addr {
	case 0x8000, 0x8001, 0x8002:
		m.mem.setProm8kBank(byte(addr)+4, uint32(data))
	case 0x9000:
		switch data {
		case 0x2b:
			m.rTyp = 1
		case 0xa8:
			m.rTyp = 2
		case 0x1f:
			m.rTyp = 3
		case 0x7c:
			m.rTyp = 4
		case 0x18:
			m.rTyp = 5
		case 0x60:
			m.rTyp = 6
		default:
			m.rTyp = 0
		}
		m.mem.setVrom1kBank(0, uint32(data))
	case 0x9001, 0x9003, 0x9004, 0x9005, 0x9006, 0x9007:
		m.mem.setVrom1kBank(byte(addr), uint32(data))
	case 0x9002:
		if m.rTyp == 2 && data != 0xe8 {
			m.rTyp = 0
		}
		m.mem.setVrom1kBank(2, uint32(data))
	case 0xc000:
		m.irqEn, m.irqCnt = m.irqLatch != 0, m.irqLatch
	case 0xc001:
		m.irqLatch = data
	case 0xc002:
		m.irqEn = false
		m.clearIntr()
	case 0xc003:
		m.irqCnt = data
	}
}

func (m *mapper160) hSync(scanline uint16) {
	switch scanline {
	case 0, 239:
		switch m.rTyp {
		case 1:
			m.mem.setVrom8kBank8(88, 89, 90, 91, 88, 89, 90, 91)
		case 2:
			m.mem.setVrom8kBank8(120, 121, 122, 123, 120, 121, 122, 123)
		case 3:
			m.mem.setVrom8kBank8(124, 125, 126, 127, 124, 125, 126, 127)
		case 5:
			m.mem.setVrom8kBank8(112, 113, 114, 115, 116, 117, 118, 119)
		case 6:
			m.mem.setVrom8kBank8(92, 93, 94, 95, 124, 125, 126, 127)
		}
	case 64:
		if m.rTyp == 4 {
			var b uint32 = 104
			if m.mem.ppuBanks[8][0x0150] != 0x0a {
				b += 4
			}
			m.mem.setVrom1kBank(0, b)
			m.mem.setVrom1kBank(1, b+1)
			m.mem.setVrom1kBank(2, b+2)
			m.mem.setVrom1kBank(3, b+3)
		}
	case 128:
		switch m.rTyp {
		case 4:
			m.mem.setVrom1kBank(0, 104)
			m.mem.setVrom1kBank(1, 105)
			m.mem.setVrom1kBank(2, 106)
			m.mem.setVrom1kBank(3, 107)
		case 5:
			m.mem.setVrom8kBank8(116, 117, 118, 119, 116, 117, 118, 119)
		}
	case 160:
		if m.rTyp == 6 {
			m.mem.setVrom8kBank8(96, 97, 98, 99, 124, 125, 126, 127)
		}
	}
	if m.irqEn {
		if m.irqCnt == 0xff {
			m.irqEn, m.irqCnt = false, 0
			m.setIntr()
		} else {
			m.irqCnt++
		}
	}
}

// 162

type mapper162 struct {
	baseMapper
	r [4]byte
}

func newMapper162(bm *baseMapper) Mapper {
	return &mapper162{baseMapper: *bm}
}

func (m *mapper162) setCpuBanks() {
	var b byte
	switch m.r[3] {
	case 0x04:
		b = ((m.r[0] & 0x0f) + ((m.r[1] & 0x02) >> 1)) | ((m.r[2] & 0x03) << 4)
	case 0x07:
		b = (m.r[0] & 0x0f) | ((m.r[1] & 0x01) << 4) | ((m.r[2] & 0x03) << 4)
	}
	m.mem.setProm32kBank(uint32(b))
}

func (m *mapper162) reset() {
	m.r[0], m.r[1], m.r[2], m.r[3] = 3, 0, 0, 7
	m.setCpuBanks()
	m.mem.setCram8kBank(0)
}

func (m *mapper162) writeLow(addr uint16, data byte) {
	switch addr {
	case 0x5000, 0x5100, 0x5200:
		m.r[(addr&0x0300)>>8] = data
		m.setCpuBanks()
		m.mem.setCram8kBank(0)
	case 0x5300:
		m.r[3] = data
	}
	if addr >= 0x6000 {
		m.cpuBanks[addr>>13][addr&0x1fff] = data
	}
}

func (m *mapper162) hSync(scanline uint16) {
	if m.r[0]&0x80 != 0 && m.isPpuDisp() {
		if scanline < 127 {
			m.mem.setCram4kBank(4, 0)
		} else if scanline < 240 {
			m.mem.setCram4kBank(4, 1)
		}
	}
}

// 163

type mapper163 struct {
	baseMapper
	typ        byte
	trig       bool
	strobe     bool
	secur      byte
	r0, r1, r2 byte
}

func newMapper163(bm *baseMapper) Mapper {
	return &mapper163{baseMapper: *bm}
}

func (m *mapper163) reset() {
	patch := m.sys.conf.PatchTyp
	if patch&0x01 != 0 {
		m.typ = 1
	} else if patch&0x02 != 0 {
		m.typ = 2
	}
	m.trig, m.strobe = false, true
	m.secur = 0
	m.r0, m.r1 = 0, 0xff
	m.mem.setProm32kBank(15)
	if patch&0x04 != 0 {
		m.mem.setProm32kBank(0)
	}
}

func (m *mapper163) readLow(addr uint16) byte {
	if addr >= 0x5000 && addr < 0x6000 {
		switch addr & 0x7700 {
		case 0x5100:
			return m.secur
		case 0x5500:
			if m.trig {
				return m.secur
			}
			return 0
		}
		return 0x04
	}
	return m.baseMapper.readLow(addr)
}

func (m *mapper163) writeLow(addr uint16, data byte) {
	if addr == 0x5101 {
		if m.strobe && data == 0 {
			m.trig = !m.trig
		}
		m.strobe = data != 0
	} else if addr == 0x5100 && data == 0x06 {
		m.mem.setProm32kBank(3)
	} else if addr >= 0x4020 && addr < 0x6000 {
		switch addr & 0x7300 {
		case 0x5000:
			m.r1 = data
			m.mem.setProm32kBank(uint32((m.r1 & 0x0f) | (m.r0 << 4)))
			if (m.r1&0x80 == 0 && m.sys.scanline < 128) || m.typ == 1 {
				m.mem.setCram8kBank(0)
			}
		case 0x5100:
			m.r2 = data
		case 0x5200:
			m.r0 = data
			m.mem.setProm32kBank(uint32((m.r1 & 0x0f) | (m.r0 << 4)))
		case 0x5300:
			m.secur = data
		}
	} else if addr >= 0x6000 {
		m.cpuBanks[addr>>13][addr&0x1fff] = data
		if addr >= 0x7900 && addr < 0x7a00 {
			m.sys.apu.write(0x4011, data)
		}
	}
}

func (m *mapper163) hSync(scanline uint16) {
	if m.r1&0x80 != 0 && m.isPpuDisp() {
		if scanline < 127 {
			if m.typ == 1 {
				m.mem.setCram4kBank(0, 0)
				m.mem.setCram4kBank(4, 0)
			}
		} else if scanline == 127 {
			m.mem.setCram4kBank(0, 1)
			m.mem.setCram4kBank(4, 1)
		} else if scanline == ScreenHeight-1 {
			if m.typ != 1 {
				m.mem.setCram4kBank(0, 0)
				m.mem.setCram4kBank(4, 0)
				if m.typ == 2 {
					m.mem.setCram4kBank(4, 1)
				}
			}
		}
	}
}

// 164

type mapper164 struct {
	baseMapper
	pMode  bool
	a      byte
	r0, r1 byte
}

func newMapper164(bm *baseMapper) Mapper {
	return &mapper164{baseMapper: *bm}
}

func (m *mapper164) setCpuBanks() {
	b := uint32(m.r0&0x01) << 5
	switch (m.r0 >> 4) & 0x07 {
	case 0x00, 0x02, 0x04, 0x06:
		bank := uint32(m.r0&0x0f) + uint32((m.r0&0x20)>>1)
		m.mem.setProm16kBank(4, b+bank)
		m.mem.setProm16kBank(6, b+0x1f)
	case 0x05:
		bank := uint32(m.r0 & 0x0f)
		m.mem.setProm32kBank(bank + (b >> 1))
	case 0x07:
		bank := uint32(m.r0 & 0x0f)
		bank += (bank & 0x08) << 1
		m.mem.setProm16kBank(4, b+bank)
		m.mem.setProm16kBank(6, b+0x1f)
	}
}

func (m *mapper164) reset() {
	m.r0, m.r1 = 0, 0
	m.sys.ppu.bExtLatch = true
	m.setCpuBanks()
	m.mem.setCram8kBank(0)
}

func (m *mapper164) writeLow(addr uint16, data byte) {
	switch addr {
	case 0x5000:
		m.pMode, m.r0 = (data>>7) != 0, data
		m.setCpuBanks()
		m.mem.setCram8kBank(0)
	case 0x5100:
		m.r1 = data
		m.setCpuBanks()
		m.mem.setCram8kBank(0)
	}
	if addr >= 0x6000 {
		m.cpuBanks[addr>>13][addr&0x1fff] = data
	}
}

func (m *mapper164) ppuExtLatchX(x byte) {
	m.a = (x & 0x01) << 3
}

func (m *mapper164) ppuExtLatch(iNameTbl uint16, chL *byte, chH *byte, attr *byte) {
	ppu := m.sys.ppu
	bank := m.mem.ppuBanks[iNameTbl>>10]

	tile := (uint16(ppu.reg0&ppuReg0BgTbl) << 8) + (uint16(bank[iNameTbl&0x03ff]) << 4) + ppu.loopyY
	if !m.pMode {
		*chL, *chH = bank[tile&0x03ff], bank[(tile&0x03ff)+8]
	} else {
		b := bank[(tile&0x03f7)|uint16(m.a)]
		*chL, *chH = b, b
	}

	iAttr := ((ppu.loopyV & 0x0380) >> 4) | 0x03c0
	x := byte(iNameTbl) & 0x1f
	sh := (byte(iNameTbl) & 0x40) >> 4
	*attr = ((bank[iAttr+uint16(x>>2)] >> ((x & 0x02) | sh)) & 0x03) << 2
}

// 165

type mapper165 struct {
	baseMapper
	sw bool
	r  byte
	p  [2]byte
	c  [4]byte
}

func newMapper165(bm *baseMapper) Mapper {
	return &mapper165{baseMapper: *bm}
}

func (m *mapper165) setPpuBanksSub(iBank byte, iPage uint32) {
	if iPage != 0 {
		m.mem.setVrom4kBank(iBank, iPage>>2)
	} else {
		m.mem.setCram4kBank(iBank, 0)
	}
}

func (m *mapper165) setPpuBanks() {
	i := 0
	if m.sw {
		i++
	}
	m.setPpuBanksSub(0, uint32(m.c[i]))
	m.setPpuBanksSub(4, uint32(m.c[i+2]))
}

func (m *mapper165) reset() {
	m.sw = false
	m.r, m.p[0], m.p[1] = 0, 0, 1
	m.c[0], m.c[1], m.c[2], m.c[3] = 0, 0, 4, 4
	m.sys.ppu.bChrLatch = true
	m.mem.setProm32kBank4(uint32(m.p[0]), uint32(m.p[1]), m.nProm8kPage-2, m.nProm8kPage-1)
	m.setPpuBanks()
}

func (m *mapper165) write(addr uint16, data byte) {
	switch addr & 0xe001 {
	case 0x8000:
		m.r = data
		m.mem.setProm32kBank4(uint32(m.p[0]), uint32(m.p[1]), m.nProm8kPage-2, m.nProm8kPage-1)
		m.setPpuBanks()
	case 0x8001:
		r := m.r & 0x07
		switch r {
		case 0x00, 0x02:
			m.c[r] = data & 0xfc
			if !m.sw {
				m.setPpuBanks()
			}
		case 0x01, 0x03:
			m.c[(r>>1)+1] = data & 0xfc
			if m.sw {
				m.setPpuBanks()
			}
		case 0x06, 0x07:
			m.p[r&0x01] = data
			m.mem.setProm32kBank4(uint32(m.p[0]), uint32(m.p[1]), m.nProm8kPage-2, m.nProm8kPage-1)
		}
	case 0xa000:
		if data&0x01 != 0 {
			m.mem.setVramMirror(memVramMirrorH)
		} else {
			m.mem.setVramMirror(memVramMirrorV)
		}
	}
}

func (m *mapper165) ppuChrLatch(addr uint16) {
	switch addr & 0x1ff0 {
	case 0x1fd0:
		m.sw = false
		m.setPpuBanks()
	case 0x1fe0:
		m.sw = true
		m.setPpuBanks()
	}
}

// 167

type mapper167 struct {
	baseMapper
	r [4]byte
}

func newMapper167(bm *baseMapper) Mapper {
	return &mapper167{baseMapper: *bm}
}

func (m *mapper167) setCpuBanks() {
	b := ((uint32(m.r[0]^m.r[1]) & 0x10) << 1) | (uint32(m.r[2]^m.r[3]) & 0x1f)
	if m.r[1]&0x08 != 0 {
		b &= 0xfe
		m.mem.setProm16kBank(4, b+1)
		m.mem.setProm16kBank(6, b)
	} else if m.r[1]&0x04 != 0 {
		m.mem.setProm16kBank(4, 31)
		m.mem.setProm16kBank(6, b)
	} else {
		m.mem.setProm16kBank(4, b)
		m.mem.setProm16kBank(6, 32)
	}
}

func (m *mapper167) reset() {
	m.r[0], m.r[1], m.r[2], m.r[3] = 0, 0, 0, 0
	m.setCpuBanks()
	m.mem.setCram8kBank(0)
}

func (m *mapper167) write(addr uint16, data byte) {
	m.r[(addr>>13)&0x03] = data
	m.setCpuBanks()
	m.mem.setCram8kBank(0)
}

// 168

type mapper168 struct {
	baseMapper
	typ    byte
	sw     bool
	r0, r1 byte
}

func newMapper168(bm *baseMapper) Mapper {
	m := &mapper168{baseMapper: *bm}
	patch := bm.sys.conf.PatchTyp
	if patch&0x01 != 0 {
		m.typ = 1
		bm.sys.tvFormat = tvFormats[2]
	} else if patch&0x02 != 0 {
		m.typ = 2
	}
	return m
}

func (m *mapper168) reset() {
	m.sw = false
	m.r0, m.r1 = 0, 0
	m.mem.setProm16kBank(4, 0)
	m.mem.setProm16kBank(6, 0)
	if m.typ == 1 {
		m.mem.setProm32kBank(0)
	}
	m.sys.ppu.bExtLatch = true
}

func (m *mapper168) readLow(addr uint16) byte {
	if addr == 0x5300 {
		return 0x8f
	}
	return m.baseMapper.readLow(addr)
}

func (m *mapper168) writeLow(addr uint16, data byte) {
	if addr == 0x5000 || addr == 0x5200 {
		if addr == 0x5000 {
			m.r0 = data
		} else {
			m.r1 = data & 0x07
		}
		if m.r1 < 0x04 {
			m.mem.setProm16kBank(4, uint32(m.r0))
		} else {
			m.mem.setProm32kBank(uint32(m.r0))
		}
		switch m.r1 {
		case 0x00:
			m.mem.setVramMirror(memVramMirrorV)
			m.sw = false
		case 0x01:
		case 0x03:
			m.mem.setVramMirror(memVramMirrorH)
			m.sw = false
		case 0x05:
			if m.typ == 2 && m.r0 == 0x04 {
				m.sys.ppu.bExtLatch = false
				m.mem.setVramMirror(memVramMirrorH)
			}
		}
	} else if addr >= 0x6000 {
		m.cpuBanks[addr>>13][addr&0x1fff] = data
	}
}

func (m *mapper168) write(addr uint16, data byte) {
	if m.typ == 1 {
		m.mem.setProm32kBank(uint32(data & 0x1f))
		if data&0x40 != 0 {
			m.mem.setVramMirror(memVramMirrorH)
		} else {
			m.mem.setVramMirror(memVramMirrorV)
		}
		m.sw = data&0xc0 != 0
	}
}

func (m *mapper168) ppuExtLatch(iNameTbl uint16, chL *byte, chH *byte, attr *byte) {
	ppu := m.sys.ppu
	bank := m.mem.ppuBanks[iNameTbl>>10]

	tile := (uint16(ppu.reg0&ppuReg0BgTbl) << 8)
	j := byte(iNameTbl>>10) & 0x03
	if j == 0x02 || (j != 0 && m.sw) {
		tile |= 0x1000
	}
	tile += (uint16(bank[iNameTbl&0x03ff]) << 4) + ppu.loopyY
	*chL, *chH = bank[tile&0x03ff], bank[(tile&0x03ff)+8]

	iAttr := ((ppu.loopyV & 0x0380) >> 4) | 0x03c0
	x := byte(iNameTbl) & 0x1f
	sh := (byte(iNameTbl) & 0x40) >> 4
	*attr = ((bank[iAttr+uint16(x>>2)] >> ((x & 0x02) | sh)) & 0x03) << 2
}

// 174

type mapper174 struct {
	baseMapper
	r [4]byte
	p [2]uint32
	c [8]byte
}

func newMapper174(bm *baseMapper) Mapper {
	return &mapper174{baseMapper: *bm}
}

func (m *mapper174) reset() {
	m.r[0], m.r[1], m.r[2], m.r[3] = 0, 0, 0, 0
	m.p[0], m.p[1] = 0, 0
	for i := 0; i < 8; i++ {
		m.c[i] = 0
	}
	m.mem.setProm32kBank4(m.p[0], m.p[1], 62, 63)
	m.setPpuBanks()
}

func (m *mapper174) writeLow(addr uint16, data byte) {
	switch addr {
	case 0x5010:
		m.r[0] = data
		break
	case 0x5011, 0x5012:
		m.r[addr&0x03] = data
		switch m.r[0] & 0x07 {
		case 0x00:
			m.mem.setProm16kBank(6, uint32(m.r[1]&0x70)+uint32(m.r[2]))
			m.mem.setProm16kBank(6, uint32(m.r[1]&0x70)+uint32(m.r[2])+31)
		case 0x01:
			m.mem.setProm16kBank(6, uint32(m.r[1]&0x70)+uint32(m.r[2]))
			m.mem.setProm16kBank(6, uint32(m.r[1]&0x70)+uint32(m.r[2])+15)
		case 0x02:
			m.mem.setProm16kBank(4, uint32(m.r[1]&0x7f)+uint32(m.r[2])+6)
			m.mem.setProm16kBank(6, uint32(m.r[1]&0x7f)+uint32(m.r[2])+7)
		case 0x03:
			m.mem.setProm16kBank(4, uint32(m.r[1]&0x7f)+uint32(m.r[2]))
			m.mem.setProm16kBank(6, uint32(m.r[1]&0x7f)+uint32(m.r[2]))
		case 0x04:
			m.mem.setProm32kBank((uint32(m.r[1]&0x7f) >> 1) + (uint32(m.r[2]) >> 1))
		case 0x05:
			m.mem.setProm16kBank(4, uint32(m.r[1]&0x7f)+uint32(m.r[2]))
			m.mem.setProm16kBank(6, uint32(m.r[1]&0x7f)+uint32(m.r[2])+7)
		}
	}
	if addr >= 0x6000 {
		m.cpuBanks[addr>>13][addr&0x1fff] = data
	}
}

func (m *mapper174) write(addr uint16, data byte) {
	if addr == 0xa000 {
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
	if r := m.r[0] & 0x07; r < 0x05 {
		switch addr {
		case 0x8000:
			m.r[3] = data
		case 0x8001:
			r := m.r[3] & 0x0f
			switch r {
			case 0x00, 0x01:
				m.c[r<<1] = data
				m.c[(r<<1)+1] = data + 1
				m.setPpuBanks()
			case 0x02, 0x03, 0x04, 0x05:
				m.c[r+2] = data
				m.setPpuBanks()
			case 0x06, 0x07:
				m.p[r&0x01] = uint32(data) + (uint32(m.r[1]&0x7f) << 1) + uint32(m.r[2]<<1)
				m.mem.setProm8kBank(4, m.p[0])
				m.mem.setProm8kBank(5, m.p[1])
			}
		}
	} else if r == 0x05 {
		m.mem.setProm16kBank(4, uint32(data)+uint32(m.r[1]&0x7f)+uint32(m.r[2]))
	}
}

func (m *mapper174) setPpuBanks() {
	for i := byte(0); i < 8; i++ {
		m.mem.setCram1kBank(i, uint32(m.c[i]&0x07))
	}
}

// 175

type mapper175 struct {
	baseMapper
	r byte
}

func newMapper175(bm *baseMapper) Mapper {
	return &mapper175{baseMapper: *bm}
}

func (m *mapper175) reset() {
	m.r = 0
	m.mem.setProm16kBank(4, 0)
	m.mem.setProm16kBank(6, 0)
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper175) read(addr uint16) byte {
	if addr == 0xfffc {
		m.mem.setProm16kBank(4, uint32(m.r&0x0f))
		m.mem.setProm16kBank(6, uint32(m.r&0x0f)<<1)
	}
	return m.cpuBanks[addr>>13][addr&0x1fff]
}

func (m *mapper175) write(addr uint16, data byte) {
	switch addr {
	case 0x8000:
		if data&0x04 != 0 {
			m.mem.setVramMirror(memVramMirrorH)
		} else {
			m.mem.setVramMirror(memVramMirrorV)
		}
	case 0xa000:
		m.r = data
		m.mem.setProm8kBank(7, (uint32(m.r&0x0f)<<1)+1)
		m.mem.setVrom8kBank(uint32(m.r & 0x0f))
	}
}
