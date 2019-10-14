package core

// 000

type mapper000 struct {
	baseMapper
}

func newMapper000(bm *baseMapper) Mapper {
	return &mapper000{baseMapper: *bm}
}

func (m *mapper000) reset() {
	switch m.nProm8kPage >> 1 {
	case 1:
		m.mem.setProm16kBank(4, 0)
		m.mem.setProm16kBank(6, 0)
	case 2:
		m.mem.setProm32kBank(0)
	}
}

// 001

type mapper001 struct {
	baseMapper
	largeTyp bool
	wramTyp  byte
	reg      [4]byte
	regBuf   byte
	shift    byte
	wramCnt  byte
	wramBank byte
	prevAddr uint16
}

func newMapper001(bm *baseMapper) Mapper {
	return &mapper001{baseMapper: *bm}
}

func (m *mapper001) reset() {
	patch := m.sys.conf.PatchTyp
	if patch&0x01 != 0 {
		m.wramTyp = 1
		m.wramCnt, m.wramBank = 0, 0
	} else if patch&0x02 != 0 {
		m.wramTyp = 2
	}

	m.reg[0], m.reg[1], m.reg[2], m.reg[3] = 0x0c, 0, 0, 0
	m.regBuf, m.shift = 0, 0
	if m.nProm8kPage < 64 {
		m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
		m.largeTyp = false
	} else {
		m.mem.setProm16kBank(4, 0)
		m.mem.setProm16kBank(6, 15)
		m.largeTyp = true
	}
}

func (m *mapper001) setVramMirror() {
	if m.reg[0]&0x02 != 0 {
		if m.reg[0]&0x01 != 0 {
			m.mem.setVramMirror(memVramMirrorH)
		} else {
			m.mem.setVramMirror(memVramMirrorV)
		}
	} else {
		if m.reg[0]&0x01 != 0 {
			m.mem.setVramMirror(memVramMirror4H)
		} else {
			m.mem.setVramMirror(memVramMirror4L)
		}
	}
}

func (m *mapper001) write(addr uint16, data byte) {
	if m.wramTyp == 1 && addr == 0xbfff {
		m.wramCnt++
		m.wramBank += data & 0x01
		if m.wramCnt == 5 {
			if m.wramBank != 0 {
				m.mem.setCpuBank(3, m.mem.wram[0x2000:], memBankTypRam)
			} else {
				m.mem.setCpuBank(3, m.mem.wram[:], memBankTypRam)
			}
			m.wramCnt, m.wramBank = 0, 0
		}
	}

	if !m.largeTyp {
		if addr&0x6000 != m.prevAddr&0x6000 {
			m.shift, m.regBuf = 0, 0
		}
		m.prevAddr = addr
	}
	if data&0x80 != 0 {
		m.shift, m.regBuf = 0, 0
		m.reg[0] |= 0x0c
		return
	}
	if data&0x01 != 0 {
		m.regBuf |= 1 << m.shift
	}
	m.shift++
	if m.shift < 5 {
		return
	}

	addr = (addr & 0x7fff) >> 13
	m.reg[addr] = m.regBuf
	m.shift, m.regBuf = 0, 0
	if !m.largeTyp {
		switch addr {
		case 0:
			m.setVramMirror()
		case 1, 2:
			if m.nVrom1kPage != 0 {
				if m.reg[0]&0x10 != 0 {
					m.mem.setVrom4kBank(0, uint32(m.reg[1]))
					m.mem.setVrom4kBank(4, uint32(m.reg[2]))
				} else {
					m.mem.setVrom8kBank(uint32(m.reg[1] >> 1))
				}
			} else if m.reg[0]&0x10 != 0 {
				m.mem.setCram4kBank(0, uint32(m.reg[addr]))
			}
		case 3:
			if m.reg[0]&0x08 != 0 {
				if m.reg[0]&0x04 != 0 {
					m.mem.setProm16kBank(4, uint32(m.reg[3]))
					m.mem.setProm16kBank(6, (m.nProm8kPage>>1)-1)
				} else {
					m.mem.setProm16kBank(4, 0)
					m.mem.setProm16kBank(6, uint32(m.reg[3]))
				}
			} else {
				m.mem.setProm32kBank(uint32(m.reg[3] >> 1))
			}
		}
	} else {
		if m.wramTyp == 2 {
			if m.reg[1]&0x18 != 0 {
				m.mem.setCpuBank(3, m.mem.wram[0x2000:], memBankTypRam)
			} else {
				m.mem.setCpuBank(3, m.mem.wram[0x2000:], memBankTypRam)
			}
		}

		var promBase uint32 = 0
		if m.nProm8kPage >= 64 {
			promBase = uint32(m.reg[1] & 0x10)
		}
		if m.reg[0]&0x08 != 0 {
			if m.reg[0]&0x04 != 0 {
				m.mem.setProm16kBank(4, promBase+uint32(m.reg[3]&0x0f))
				if m.nProm8kPage >= 64 {
					m.mem.setProm16kBank(6, promBase+15)
				}
			} else {
				if m.nProm8kPage >= 64 {
					m.mem.setProm16kBank(4, promBase)
				}
				m.mem.setProm16kBank(6, promBase+uint32(m.reg[3]&0x0f))
			}
		} else {
			m.mem.setProm32kBank(uint32(m.reg[3] >> 1))
		}

		if m.nVrom1kPage != 0 {
			if m.reg[0]&0x10 != 0 {
				m.mem.setVrom4kBank(0, uint32(m.reg[1]))
				m.mem.setVrom4kBank(4, uint32(m.reg[2]))
			} else {
				m.mem.setVrom8kBank(uint32(m.reg[1] >> 1))
			}
		} else if m.reg[0]&0x10 != 0 {
			m.mem.setCram4kBank(0, uint32(m.reg[1]))
			m.mem.setCram4kBank(4, uint32(m.reg[2]))
		}
		if addr == 0 {
			m.setVramMirror()
		}
	}
}

// 002

type mapper002 struct {
	baseMapper
	patchTyp byte
}

func newMapper002(bm *baseMapper) Mapper {
	return &mapper002{baseMapper: *bm}
}

func (m *mapper002) reset() {
	patch := byte(m.sys.conf.PatchTyp)
	if patch&0x01 != 0 {
		m.patchTyp = 1
	} else if patch&0x02 != 0 {
		m.patchTyp = 2
	}
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
}

func (m *mapper002) writeLow(addr uint16, data byte) {
	if m.sys.rom.bSaveRam {
		m.baseMapper.writeLow(addr, data)
	} else if m.patchTyp == 1 && addr >= 0x5000 {
		m.mem.setProm16kBank(4, uint32(data))
	}
}

func (m *mapper002) write(addr uint16, data byte) {
	if m.patchTyp != 2 {
		m.mem.setProm16kBank(4, uint32(data))
	} else {
		m.mem.setProm16kBank(4, uint32(data>>4))
	}
}

// 003

type mapper003 struct {
	baseMapper
}

func newMapper003(bm *baseMapper) Mapper {
	return &mapper003{baseMapper: *bm}
}

func (m *mapper003) reset() {
	switch m.nProm8kPage >> 1 {
	case 1:
		m.mem.setProm16kBank(4, 0)
		m.mem.setProm16kBank(6, 0)
	case 2:
		m.mem.setProm32kBank(0)
	}
}

func (m *mapper003) write(addr uint16, data byte) {
	m.mem.setVrom8kBank(uint32(data))
}

// 004

type mapper004 struct {
	baseMapper
	irqTyp    byte
	irqEn     bool
	irqCnt    byte
	irqLatch  byte
	irqReq    bool
	irqPre    bool
	irqPreVbl bool
	p0, p1    byte
	r         byte
	c         [8]byte
}

func newMapper004(bm *baseMapper) Mapper {
	return &mapper004{baseMapper: *bm}
}

func (m *mapper004) setCpuBanks() {
	if m.r&0x40 != 0 {
		m.mem.setProm32kBank4(m.nProm8kPage-2, uint32(m.p1), uint32(m.p0), m.nProm8kPage-1)
	} else {
		m.mem.setProm32kBank4(uint32(m.p0), uint32(m.p1), m.nProm8kPage-2, m.nProm8kPage-1)
	}
}

func (m *mapper004) setPpuBanks() {
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

func (m *mapper004) reset() {
	patch := m.sys.conf.PatchTyp
	m.irqTyp = 0
	if patch&0x01 != 0 {
		m.irqTyp = 1
	} else if patch&0x02 != 0 {
		m.irqTyp = 2
	} else if patch&0x04 != 0 {
		m.irqTyp = 3
	} else if patch&0x08 != 0 {
		m.irqTyp = 4
	} else if patch&0x10 != 0 {
		m.irqTyp = 5
	}
	if patch&0x20 != 0 {
		m.sys.renderMode = RenderModeTile
	} else if patch&0x40 != 0 {
		m.sys.renderMode = RenderModePost
	}

	m.irqEn, m.irqCnt, m.irqLatch = false, 0, 0xff
	m.irqReq, m.irqPre, m.irqPreVbl = false, false, false
	m.r, m.p0, m.p1 = 0, 0, 1
	for i := byte(0); i < 8; i++ {
		m.c[i] = i
	}
	m.setCpuBanks()
	m.setPpuBanks()
}

func (m *mapper004) readLow(addr uint16) byte {
	if addr >= 0x5000 && addr < 0x6000 {
		return m.mem.xram[addr&0x1fff]
	}
	return m.baseMapper.readLow(addr)
}

func (m *mapper004) writeLow(addr uint16, data byte) {
	if addr >= 0x5000 && addr < 0x6000 {
		m.mem.xram[addr&0x1fff] = data
		return
	}
	m.baseMapper.writeLow(addr, data)
}

func (m *mapper004) write(addr uint16, data byte) {
	switch addr & 0xE001 {
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
			m.setPpuBanks()
		case 0x06:
			m.p0 = data
			m.setCpuBanks()
		case 0x07:
			m.p1 = data
			m.setCpuBanks()
		}
	case 0xa000:
		if m.sys.rom.b4Screen {
			if data&0x01 != 0 {
				m.mem.setVramMirror(memVramMirrorH)
			} else {
				m.mem.setVramMirror(memVramMirrorV)
			}
		}
	case 0xc000:
		switch m.irqTyp {
		case 1, 5:
			m.irqCnt = data
		case 4:
			m.irqLatch = 0x07
		default:
			m.irqLatch = data
		}
	case 0xc001:
		if m.irqTyp == 1 || m.irqTyp == 5 {
			m.irqLatch = data
		} else if m.sys.scanline < ScreenHeight || m.irqTyp == 2 {
			m.irqCnt, m.irqPre = m.irqCnt|0x80, true
		} else {
			m.irqCnt, m.irqPre, m.irqPreVbl = m.irqCnt|0x80, false, true
		}
	case 0xe000:
		m.irqEn, m.irqReq = false, false
		m.clearIntr()
	case 0xe001:
		m.irqEn, m.irqReq = true, false
	}
}

func (m *mapper004) hSync(scanline uint16) {
	switch m.irqTyp {
	case 1:
		if scanline < ScreenHeight && m.isPpuDisp() && m.irqEn {
			if m.irqCnt == 0 {
				m.irqCnt, m.irqReq = m.irqLatch, true
			}
			if m.irqCnt > 0 {
				m.irqCnt--
			}
		}
		if m.irqReq {
			m.setIntr()
		}
	case 5:
		if scanline < ScreenHeight && m.isPpuDisp() && m.irqEn {
			m.irqCnt--
			if m.irqCnt == 0 {
				m.irqCnt, m.irqReq = m.irqLatch, true
			}
		}
		if m.irqReq {
			m.setIntr()
		}
	default:
		if scanline < ScreenHeight && m.isPpuDisp() {
			if m.irqPreVbl {
				m.irqCnt, m.irqPreVbl = m.irqLatch, false
			}
			if m.irqPre {
				m.irqCnt, m.irqPre = m.irqLatch, false
				if m.irqTyp == 3 && scanline == 0 {
					m.irqCnt--
				}
			} else if m.irqCnt > 0 {
				m.irqCnt--
			}
			if m.irqCnt == 0 {
				if m.irqEn {
					m.irqReq = true
					m.setIntr()
				}
				m.irqPre = true
			}
		}
	}
}

// 005

type mapper005 struct {
	baseMapper
	sramSize byte
	irqPatch bool
	chrPatch bool

	prgSize, chrSize     byte
	sramWEA, sramWEB     bool
	graphMode            byte
	fillChr, fillPal     byte
	chrMode              bool
	splitCont            byte
	splitScrl, splitPage byte
	irqEn                bool
	irqStatus, irqClear  byte
	irqLine, irqScanline byte
	multA, multB         byte
	splitX, splitY       byte
	splitAddr            uint16

	ntTyps  [4]byte
	c0, c1  [8]byte
	bgBanks [8][]byte
}

func newMapper005(bm *baseMapper) Mapper {
	return &mapper005{baseMapper: *bm}
}

func (m *mapper005) setCpuBankAlt(iBank byte, data byte) {
	switch m.sramSize {
	case 0:
		if data > 3 {
			data = 8
		} else {
			data = 0
		}
	case 1:
		if data > 3 {
			data = 1
		} else {
			data = 0
		}
	case 2:
		if data > 3 {
			data = 8
		}
	}
	if data != 8 {
		m.mem.setCpuBank(iBank, m.mem.wram[uint32(data)<<13:], memBankTypRam)
	} else {
		m.mem.cpuBanksTyp[iBank] = memBankTypRom
	}
}

func (m *mapper005) reset() {
	patch := m.sys.conf.PatchTyp
	if patch&0x01 != 0 {
		m.sramSize = 1
	} else if patch&0x02 != 0 {
		m.sramSize = 2
	}
	if patch&0x04 != 0 {
		m.irqPatch = true
	}
	if patch&0x08 != 0 {
		m.chrPatch = true
	}

	m.prgSize, m.chrSize = 3, 3
	m.sramWEA, m.sramWEB = false, false
	m.graphMode = 0
	m.fillChr, m.fillPal = 0, 0
	m.chrMode = false
	m.splitCont, m.splitScrl, m.splitPage = 0, 0, 0
	m.irqEn = false
	m.irqStatus, m.irqClear, m.irqLine, m.irqScanline = 0, 0, 0, 0
	m.multA, m.multB = 0, 0

	m.ntTyps[0], m.ntTyps[1], m.ntTyps[2], m.ntTyps[3] = 0, 0, 0, 0
	for i := byte(0); i < 8; i++ {
		m.c0[i] = i
		m.c1[i] = (i & 0x03) | 0x04
		b := uint32(i) << 10
		m.bgBanks[i] = m.mem.vrom[b : b+0x0400 : b+0x0400]
	}

	b := m.nProm8kPage - 1
	m.mem.setProm32kBank4(b, b, b, b)
	m.mem.setVrom8kBank(0)
	m.setCpuBankAlt(3, 0)
	m.sys.ppu.bExtLatch = true
}

func (m *mapper005) readLow(addr uint16) byte {
	switch addr {
	case 0x5015:
		return 0
	case 0x5204:
		data := m.irqStatus
		m.irqStatus &^= 0x80
		m.clearIntr()
		return data
	case 0x5205:
		return m.multA * m.multB
	case 0x5206:
		w := uint16(m.multA) * uint16(m.multB)
		return byte(w >> 8)
	}
	if addr >= 0x5c00 && addr < 0x6000 && m.graphMode >= 2 {
		return m.mem.vram[(addr&0x03ff)+0x0800]
	}
	return m.baseMapper.readLow(addr)
}

func (m *mapper005) setPpuBanks() {
	if !m.chrMode {
		switch m.chrSize {
		case 0:
			m.mem.setVrom8kBank(uint32(m.c0[7]))
		case 1:
			m.mem.setVrom4kBank(0, uint32(m.c0[3]))
			m.mem.setVrom4kBank(4, uint32(m.c0[7]))
		case 2:
			m.mem.setVrom2kBank(0, uint32(m.c0[1]))
			m.mem.setVrom2kBank(2, uint32(m.c0[3]))
			m.mem.setVrom2kBank(4, uint32(m.c0[5]))
			m.mem.setVrom2kBank(6, uint32(m.c0[7]))
		case 3:
			for i := byte(0); i < 8; i++ {
				m.mem.setVrom1kBank(i, uint32(m.c0[i]))
			}
		}
	} else {
		switch m.chrSize {
		case 0:
			b := (uint32(m.c1[7]) % (m.nVrom1kPage << 3)) << 13
			for i := 0; i < 8; i++ {
				m.bgBanks[i] = m.mem.vrom[b : b+0x0400 : b+0x0400]
				b += 0x0400
			}
		case 1:
			mod := m.nVrom1kPage << 2
			b0, b1 := (uint32(m.c1[3])%mod)<<12, (uint32(m.c1[7])%mod)<<12
			for i := 0; i < 4; i++ {
				m.bgBanks[i] = m.mem.vrom[b0 : b0+0x0400 : b0+0x0400]
				m.bgBanks[i+4] = m.mem.vrom[b1 : b1+0x0400 : b1+0x0400]
				b0, b1 = b0+0x0400, b1+0x0400
			}
		case 2:
			mod := m.nVrom1kPage << 1
			b0, b1 := (uint32(m.c1[1])%mod)<<11, (uint32(m.c1[3])%mod)<<11
			b2, b3 := (uint32(m.c1[5])%mod)<<11, (uint32(m.c1[7])%mod)<<11
			for i := 0; i < 2; i++ {
				m.bgBanks[i] = m.mem.vrom[b0 : b0+0x0400 : b0+0x0400]
				m.bgBanks[i+2] = m.mem.vrom[b1 : b1+0x0400 : b1+0x0400]
				m.bgBanks[i+4] = m.mem.vrom[b2 : b2+0x0400 : b2+0x0400]
				m.bgBanks[i+6] = m.mem.vrom[b3 : b3+0x0400 : b3+0x0400]
				b0, b1, b2, b3 = b0+0x0400, b1+0x0400, b2+0x0400, b3+0x0400
			}
		case 3:
			mod := m.nVrom1kPage
			for i := 0; i < 8; i++ {
				b := (uint32(m.c1[i]) % mod) << 10
				m.bgBanks[i] = m.mem.vrom[b : b+0x0400 : b+0x0400]
			}
		}
	}
}

func (m *mapper005) writeLow(addr uint16, data byte) {
	switch addr {
	case 0x5100:
		m.prgSize = data & 0x03
	case 0x5101:
		m.chrSize = data & 0x03
	case 0x5102:
		m.sramWEA = data&0x02 != 0
	case 0x5103:
		m.sramWEB = data&0x01 != 0
	case 0x5104:
		m.graphMode = data & 0x03
	case 0x5105:
		for i := byte(0); i < 4; i++ {
			m.ntTyps[i] = data & 0x03
			m.mem.setVram1kBank(8+i, uint32(data&0x03))
			data >>= 2
		}
	case 0x5106:
		m.fillChr = data
	case 0x5107:
		m.fillPal = data & 0x03
	case 0x5113:
		m.setCpuBankAlt(3, data&0x07)
	case 0x5114, 0x5115, 0x5116, 0x5117:
		if data&0x80 != 0 {
			p := uint32(data & 0x7f)
			switch addr & 0x07 {
			case 0x04:
				if m.prgSize == 3 {
					m.mem.setProm8kBank(4, p)
				}
			case 0x05:
				switch m.prgSize {
				case 1, 2:
					m.mem.setProm16kBank(4, p>>1)
				case 3:
					m.mem.setProm8kBank(5, p)
				}
			case 0x06:
				switch m.prgSize {
				case 2, 3:
					m.mem.setProm8kBank(6, p)
				}
			case 0x07:
				switch m.prgSize {
				case 0:
					m.mem.setProm32kBank(p >> 2)
				case 1:
					m.mem.setProm16kBank(6, p>>1)
				case 2, 3:
					m.mem.setProm8kBank(7, p)
				}
			}
		} else {
			switch addr & 0x07 {
			case 0x04:
				if m.prgSize == 3 {
					m.setCpuBankAlt(4, data&0x07)
				}
			case 0x05:
				switch m.prgSize {
				case 1, 2:
					m.setCpuBankAlt(4, data&0x06)
					m.setCpuBankAlt(5, (data&0x06)+1)
				case 3:
					m.setCpuBankAlt(5, data&0x07)
				}
			case 0x06:
				switch m.prgSize {
				case 2, 3:
					m.setCpuBankAlt(6, data&0x07)
				}
			}
		}
	case 0x5120, 0x5121, 0x5122, 0x5123, 0x5124, 0x5125, 0x5126, 0x5127:
		m.chrMode, m.c0[addr&0x07] = false, data
		m.setPpuBanks()
	case 0x5128, 0x5129, 0x512a, 0x512b:
		m.chrMode, m.c1[addr&0x03], m.c1[(addr&0x03)+4] = true, data, data
		m.setPpuBanks()
	case 0x5200:
		m.splitCont = data
	case 0x5201:
		m.splitScrl = data
	case 0x5202:
		m.splitPage = data & 0x3f
	case 0x5203:
		m.irqLine = data
		m.clearIntr()
	case 0x5204:
		m.irqEn = data&0x80 != 0
		m.clearIntr()
	case 0x5205:
		m.multA = data
	case 0x5206:
		m.multB = data
	default:
		if addr >= 0x5C00 && addr < 0x6000 {
			switch m.graphMode {
			case 2:
				m.mem.vram[(addr&0x03ff)+0x0800] = data
			case 3:
			default:
				if m.irqStatus&0x40 != 0 {
					m.mem.vram[(addr&0x03ff)+0x0800] = data
				} else {
					m.mem.vram[(addr&0x03ff)+0x0800] = 0
				}
			}
		} else if addr >= 0x6000 {
			if m.sramWEA && m.sramWEB && m.mem.cpuBanksTyp[3] == memBankTypRam {
				m.cpuBanks[3][addr&0x1fff] = data
			}
		}
	}
}

func (m *mapper005) write(addr uint16, data byte) {
	if addr >= 0x8000 && addr < 0xe000 {
		i := addr >> 13
		if m.sramWEA && m.sramWEB && m.mem.cpuBanksTyp[i] == memBankTypRam {
			m.cpuBanks[i][addr&0x1fff] = data
		}
	}
}

func (m *mapper005) hSync(scanline uint16) {
	if m.irqPatch && m.irqScanline == m.irqLine {
		m.irqStatus |= 0x80
	}
	if m.isPpuDisp() && scanline < ScreenHeight {
		m.irqScanline++
		m.irqStatus |= 0x40
		m.irqClear = 0
	} else if m.irqPatch {
		m.irqScanline = 0
		m.irqStatus &^= 0xc0
	}
	if !m.irqPatch {
		if m.irqScanline == m.irqLine {
			m.irqStatus |= 0x80
		}
		m.irqClear++
		if m.irqClear > 2 {
			m.irqScanline = 0
			m.irqStatus &^= 0xc0
			m.clearIntr()
		}
	}
	if m.irqEn && m.irqStatus&0xc0 != 0 {
		m.setIntr()
	}

	if scanline == 0 {
		m.splitY = m.splitScrl & 0x07
		m.splitAddr = uint16(m.splitScrl&0xf8) << 2
	} else if m.isPpuDisp() {
		if m.splitY == 7 {
			m.splitY = 0
			switch m.splitAddr & 0x03e0 {
			case 0x03a0, 0x03e0:
				m.splitAddr &= 0x001f
			default:
				m.splitAddr += 0x0020
			}
		} else {
			m.splitY++
		}
	}
}

func (m *mapper005) ppuExtLatchX(x byte) {
	m.splitX = x
}

func (m *mapper005) ppuExtLatch(iNameTbl uint16, chL *byte, chH *byte, attr *byte) {
	ppu, ppuBanks, vram, vrom := m.sys.ppu, m.mem.ppuBanks, m.mem.vram[:], m.mem.vrom
	bSplit := false
	if m.splitCont&0x80 != 0 {
		bSplit = (m.splitCont&0x40 == 0) == (m.splitCont&0x1f > m.splitX)
	}

	var tile uint32
	if !bSplit {
		bFill := m.ntTyps[(iNameTbl&0x0c00)>>10] == 3
		if m.graphMode == 1 {
			iNameTbl = (iNameTbl & 0x0fff) | 0x2000
			if bFill {
				tile = (uint32(m.fillChr) << 4) + uint32(ppu.loopyY)
				*attr = (m.fillPal << 2) & 0x0c
			} else {
				tile = (uint32(ppuBanks[iNameTbl>>10][iNameTbl&0x03ff]) << 4) + uint32(ppu.loopyY)
				*attr = (vram[(iNameTbl&0x03ff)+0x0800] & 0xc0) >> 4
			}
			b := vram[(iNameTbl&0x03ff)+0x0800]
			tile += (uint32(b&0x3f) % (m.nVrom1kPage >> 2)) << 12
			*chL, *chH = vrom[tile], vrom[tile+8]
		} else {
			tile = (uint32(ppu.reg0&ppuReg0BgTbl) << 8)
			if bFill {
				tile += (uint32(m.fillChr) << 4) + uint32(ppu.loopyY)
				*attr = (m.fillPal << 2) & 0x0c
			} else {
				iAttr := (iNameTbl & 0x0c00) + ((iNameTbl & 0x0380) >> 4) + ((iNameTbl & 0x001c) >> 2) + 0x23c0
				iNameTbl = (iNameTbl & 0x0fff) | 0x2000
				tile += (uint32(ppuBanks[iNameTbl>>10][iNameTbl&0x03ff]) << 4) + uint32(ppu.loopyY)
				a := ppuBanks[iAttr>>10][iAttr&0x03ff]
				if iNameTbl&0x0002 != 0 {
					a >>= 2
				}
				if iNameTbl&0x0040 != 0 {
					a >>= 4
				}
				*attr = (a & 0x03) << 2
			}
			bank := m.bgBanks[tile>>10]
			if m.chrPatch {
				bank = ppuBanks[tile>>10]
			}
			tile &= 0x03ff
			*chL, *chH = bank[tile], bank[tile+8]
		}
	} else {
		iNameTbl = ((m.splitAddr & 0x03e0) | uint16(m.splitX&0x1f)) & 0x03ff
		tile = ((uint32(m.splitPage) % (m.nVrom1kPage >> 2)) << 12) +
			(uint32(vram[iNameTbl+0x0800]) << 4) + uint32(m.splitY)
		*chL, *chH = vrom[tile], vrom[tile+8]
		iAttr := ((iNameTbl & 0x0380) >> 4) + ((iNameTbl & 0x001c) >> 2) + 0x0bc0
		a := vram[iAttr]
		if iNameTbl&0x0002 != 0 {
			a >>= 2
		}
		if iNameTbl&0x0040 != 0 {
			a >>= 4
		}
		*attr = (a & 0x03) << 2
	}
}

// 006

type mapper006 struct {
	baseMapper
	irqEn  bool
	irqCnt uint64
}

func newMapper006(bm *baseMapper) Mapper {
	return &mapper006{baseMapper: *bm}
}

func (m *mapper006) reset() {
	m.mem.setProm32kBank4(0, 1, 14, 15)
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	} else {
		m.mem.setCram8kBank(0)
	}
	m.irqEn, m.irqCnt = false, 0
}

func (m *mapper006) writeLow(addr uint16, data byte) {
	switch addr {
	case 0x42fe:
		if data&0x10 != 0 {
			m.mem.setVramMirror(memVramMirror4H)
		} else {
			m.mem.setVramMirror(memVramMirror4L)
		}
	case 0x42ff:
		if data&0x10 != 0 {
			m.mem.setVramMirror(memVramMirrorH)
		} else {
			m.mem.setVramMirror(memVramMirrorV)
		}
	case 0x4501:
		m.irqEn = false
		m.clearIntr()
	case 0x4502:
		m.irqCnt = (m.irqCnt & 0xff00) | uint64(data)
	case 0x4503:
		m.irqCnt = (m.irqCnt & 0x00ff) | (uint64(data) << 8)
		m.irqEn = true
		m.clearIntr()
	default:
		m.baseMapper.writeLow(addr, data)
	}
}

func (m *mapper006) write(addr uint16, data byte) {
	m.mem.setProm16kBank(4, uint32(data&0x3c)>>2)
	m.mem.setCram8kBank(uint32(data & 0x03))
}

func (m *mapper006) hSync(scanline uint16) {
	if m.irqEn {
		m.irqCnt += 133
		if m.irqCnt >= 0xffff {
			m.irqCnt = 0
			m.setIntr()
		}
	}
}

// 007

type mapper007 struct {
	baseMapper
}

func newMapper007(bm *baseMapper) Mapper {
	return &mapper007{baseMapper: *bm}
}

func (m *mapper007) reset() {
	m.mem.setProm32kBank(0)
	m.mem.setVramMirror(memVramMirror4L)
}

func (m *mapper007) write(addr uint16, data byte) {
	m.mem.setProm32kBank(uint32(data & 0x07))
	if data&0x10 != 0 {
		m.mem.setVramMirror(memVramMirror4H)
	} else {
		m.mem.setVramMirror(memVramMirror4L)
	}
}

// 008

type mapper008 struct {
	baseMapper
}

func newMapper008(bm *baseMapper) Mapper {
	return &mapper008{baseMapper: *bm}
}

func (m *mapper008) reset() {
	m.mem.setProm32kBank(0)
	m.mem.setVrom8kBank(0)
}

func (m *mapper008) write(addr uint16, data byte) {
	m.mem.setProm16kBank(4, uint32((data&0xf8)>>3))
	m.mem.setVrom8kBank(uint32(data & 0x07))
}

// 009

type mapper009 struct {
	baseMapper
	r      [4]byte
	latchA bool
	latchB bool
}

func newMapper009(bm *baseMapper) Mapper {
	return &mapper009{baseMapper: *bm}
}

func (m *mapper009) reset() {
	m.mem.setProm32kBank4(0, m.nProm8kPage-3, m.nProm8kPage-2, m.nProm8kPage-1)
	m.r[0], m.r[1], m.r[2], m.r[3] = 0, 4, 0, 0
	m.latchA, m.latchB = false, false
	m.mem.setVrom4kBank(0, 4)
	m.mem.setVrom4kBank(4, 0)
	m.sys.ppu.bChrLatch = true
}

func (m *mapper009) write(addr uint16, data byte) {
	switch addr & 0xf000 {
	case 0xa000:
		m.mem.setProm8kBank(4, uint32(data))
	case 0xb000:
		m.r[0] = data
		if m.latchA {
			m.mem.setVrom4kBank(0, uint32(data))
		}
	case 0xc000:
		m.r[1] = data
		if !m.latchA {
			m.mem.setVrom4kBank(0, uint32(data))
		}
	case 0xd000:
		m.r[2] = data
		if m.latchB {
			m.mem.setVrom4kBank(4, uint32(data))
		}
	case 0xe000:
		m.r[3] = data
		if !m.latchB {
			m.mem.setVrom4kBank(4, uint32(data))
		}
	case 0xf000:
		if data&0x01 != 0 {
			m.mem.setVramMirror(memVramMirrorH)
		} else {
			m.mem.setVramMirror(memVramMirrorV)
		}
	}
}

func (m *mapper009) ppuChrLatch(addr uint16) {
	switch addr & 0x1ff0 {
	case 0x0fd0:
		if !m.latchA {
			m.latchA = true
			m.mem.setVrom4kBank(0, uint32(m.r[0]))
		}
	case 0x0fe0:
		if m.latchA {
			m.latchA = false
			m.mem.setVrom4kBank(0, uint32(m.r[1]))
		}
	case 0x1fd0:
		if !m.latchB {
			m.latchB = true
			m.mem.setVrom4kBank(0, uint32(m.r[2]))
		}
	case 0x1fe0:
		if m.latchB {
			m.latchB = false
			m.mem.setVrom4kBank(0, uint32(m.r[3]))
		}
	}
}

// 010

type mapper010 struct {
	baseMapper
	latchA byte
	latchB byte
	r      [4]byte
}

func newMapper010(bm *baseMapper) Mapper {
	return &mapper010{baseMapper: *bm}
}

func (m *mapper010) reset() {
	m.latchA, m.latchB = 0xfe, 0xfe
	m.r[0], m.r[1], m.r[2], m.r[3] = 0, 4, 0, 0
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
	m.mem.setVrom4kBank(0, 4)
	m.mem.setVrom4kBank(4, 0)
	m.sys.ppu.bChrLatch = true
}

func (m *mapper010) write(addr uint16, data byte) {
	switch addr & 0xf000 {
	case 0xa000:
		m.mem.setProm16kBank(4, uint32(data))
	case 0xb000:
		m.r[0] = data
		if m.latchA == 0xfd {
			m.mem.setVrom4kBank(0, uint32(m.r[0]))
		}
	case 0xc000:
		m.r[1] = data
		if m.latchA == 0xfe {
			m.mem.setVrom4kBank(0, uint32(m.r[1]))
		}
	case 0xd000:
		m.r[2] = data
		if m.latchB == 0xfd {
			m.mem.setVrom4kBank(4, uint32(m.r[2]))
		}
	case 0xe000:
		m.r[3] = data
		if m.latchB == 0xfe {
			m.mem.setVrom4kBank(4, uint32(m.r[3]))
		}
	case 0xf000:
		if data&0x01 != 0 {
			m.mem.setVramMirror(memVramMirrorH)
		} else {
			m.mem.setVramMirror(memVramMirrorV)
		}
	}
}

func (m *mapper010) ppuChrLatch(addr uint16) {
	if addr&0x1ff0 == 0x0fd0 && m.latchA != 0xfd {
		m.latchA = 0xfd
		m.mem.setVrom4kBank(0, uint32(m.r[0]))
	} else if addr&0x1ff0 == 0x0fe0 && m.latchA != 0xfe {
		m.latchA = 0xfe
		m.mem.setVrom4kBank(0, uint32(m.r[1]))
	} else if addr&0x1ff0 == 0x1fd0 && m.latchB != 0xfd {
		m.latchB = 0xfd
		m.mem.setVrom4kBank(4, uint32(m.r[2]))
	} else if addr&0x1ff0 == 0x1fe0 && m.latchB != 0xfe {
		m.latchB = 0xfe
		m.mem.setVrom4kBank(4, uint32(m.r[3]))
	}
}

// 011

type mapper011 struct {
	baseMapper
}

func newMapper011(bm *baseMapper) Mapper {
	return &mapper011{baseMapper: *bm}
}

func (m *mapper011) reset() {
	m.mem.setProm32kBank(0)
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
	m.mem.setVramMirror(memVramMirrorV)
}

func (m *mapper011) write(addr uint16, data byte) {
	m.mem.setProm32kBank(uint32(data))
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(uint32(data >> 4))
	}
}

// 012

type mapper012 struct {
	baseMapper
	irqEn        bool
	irqCnt       byte
	irqLatch     byte
	irqPreset    bool
	irqPresetVbl bool
	r            byte
	vb           [2]uint32
	p            [2]byte
	c            [8]byte
}

func newMapper012(bm *baseMapper) Mapper {
	return &mapper012{baseMapper: *bm}
}

func (m *mapper012) setCpuBanks() {
	if m.r&0x40 != 0 {
		m.mem.setProm32kBank4(m.nProm8kPage-2, uint32(m.p[1]), uint32(m.p[0]), m.nProm8kPage-1)
	} else {
		m.mem.setProm32kBank4(uint32(m.p[0]), uint32(m.p[1]), m.nProm8kPage-2, m.nProm8kPage-1)
	}
}

func (m *mapper012) setPpuBanks() {
	if m.nVrom1kPage != 0 {
		if m.r&0x80 != 0 {
			for i := byte(0); i < 4; i++ {
				m.mem.setVrom1kBank(i, m.vb[0]+uint32(m.c[i+4]))
			}
			for i := byte(0); i < 4; i++ {
				m.mem.setVrom1kBank(i+4, m.vb[1]+uint32(m.c[i]))
			}
		} else {
			for i := byte(0); i < 8; i++ {
				m.mem.setVrom1kBank(i, m.vb[i&0x04]+uint32(m.c[i]))
			}
		}
	} else {
		if m.r&0x80 != 0 {
			for i := byte(0); i < 4; i++ {
				m.mem.setCram1kBank(i, uint32(m.c[i+4]&0x07))
			}
			for i := byte(0); i < 4; i++ {
				m.mem.setCram1kBank(i+4, uint32(m.c[i]&0x07))
			}
		} else {
			for i := byte(0); i < 8; i++ {
				m.mem.setCram1kBank(i, uint32(m.c[i]&0x07))
			}
		}
	}
}

func (m *mapper012) reset() {
	m.irqEn, m.irqCnt, m.irqLatch = false, 0, 0xff
	m.irqPreset, m.irqPresetVbl = false, false
	m.vb[0], m.vb[1] = 0, 0
	m.p[0], m.p[1] = 0, 0
	for i := byte(0); i < 8; i++ {
		m.c[i] = i
	}
	m.setCpuBanks()
	m.setPpuBanks()
}

func (m *mapper012) writeLow(addr uint16, data byte) {
	if addr > 0x4100 && addr < 0x6000 {
		m.vb[0], m.vb[1] = uint32(data&0x01)<<8, uint32(data&0x10)<<4
		m.setPpuBanks()
	} else if addr >= 0x6000 {
		m.cpuBanks[addr>>13][addr&0x1fff] = data
	}
}

func (m *mapper012) readLow(addr uint16) byte {
	return 0x01
}

func (m *mapper012) write(addr uint16, data byte) {
	switch addr & 0xe001 {
	case 0x8000:
		m.r = data
		m.setCpuBanks()
		m.setPpuBanks()
	case 0x8001:
		r := m.r & 0x07
		switch r {
		case 0x00, 0x01:
			m.c[r<<1], m.c[(r<<1)+1] = data&0xfe, (data&0xfe)|0x01
			m.setPpuBanks()
		case 0x02, 0x03, 0x04, 0x05:
			m.c[r+2] = data
			m.setPpuBanks()
		case 0x06, 0x07:
			m.p[r&0x01] = data
			m.setCpuBanks()
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
		m.irqLatch = data
	case 0xc001:
		if m.sys.scanline < ScreenHeight {
			m.irqCnt |= 0x80
			m.irqPreset = true
		} else {
			m.irqCnt |= 0x80
			m.irqPreset, m.irqPresetVbl = false, true
		}
	case 0xe000:
		m.irqEn = false
		m.clearIntr()
	case 0xe001:
		m.irqEn = true
	}
}

func (m *mapper012) hSync(scanline uint16) {
	if scanline < ScreenHeight && m.isPpuDisp() {
		if m.irqPresetVbl {
			m.irqCnt, m.irqPresetVbl = m.irqLatch, false
		}
		if m.irqPreset {
			m.irqCnt, m.irqPreset = m.irqLatch, false
		} else if m.irqCnt != 0 {
			m.irqCnt--
		}
		if m.irqCnt == 0 {
			if m.irqEn && m.irqLatch != 0 {
				m.setIntr()
			}
			m.irqPreset = true
		}
	}
}

// 013

type mapper013 struct {
	baseMapper
}

func newMapper013(bm *baseMapper) Mapper {
	return &mapper013{baseMapper: *bm}
}

func (m *mapper013) reset() {
	m.mem.setProm32kBank(0)
	m.mem.setCram4kBank(0, 0)
	m.mem.setCram4kBank(4, 0)
}

func (m *mapper013) write(addr uint16, data byte) {
	m.mem.setProm32kBank(uint32((data & 0x30) >> 4))
	m.mem.setCram4kBank(4, uint32(data&0x03))
}

// 015

type mapper015 struct {
	baseMapper
}

func newMapper015(bm *baseMapper) Mapper {
	return &mapper015{baseMapper: *bm}
}

func (m *mapper015) reset() {
	m.mem.setProm32kBank(0)
}

func (m *mapper015) write(addr uint16, data byte) {
	b := uint32(data&0x3f) << 1
	switch addr {
	case 0x8000:
		if data&0x80 != 0 {
			m.mem.setProm32kBank4(b+1, b, b+3, b+2)
		} else {
			m.mem.setProm32kBank4(b, b+1, b+2, b+3)
		}
		if data&0x40 != 0 {
			m.mem.setVramMirror(memVramMirrorH)
		} else {
			m.mem.setVramMirror(memVramMirrorV)
		}
	case 0x8001:
		if data&0x80 != 0 {
			m.mem.setProm8kBank(6, b+1)
			m.mem.setProm8kBank(7, b)
		} else {
			m.mem.setProm8kBank(6, b)
			m.mem.setProm8kBank(7, b+1)
		}
	case 0x8002:
		if data&0x80 != 0 {
			b++
		}
		m.mem.setProm32kBank4(b, b, b, b)
	case 0x8003:
		if data&0x80 != 0 {
			m.mem.setProm8kBank(6, b+1)
			m.mem.setProm8kBank(7, b)
		} else {
			m.mem.setProm8kBank(6, b)
			m.mem.setProm8kBank(7, b+1)
		}
		if data&0x40 != 0 {
			m.mem.setVramMirror(memVramMirrorH)
		} else {
			m.mem.setVramMirror(memVramMirrorV)
		}
	}
}
