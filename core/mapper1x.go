package core

// 016

type mapper016Eeprom1 struct {
	sda       bool
	prevScl   bool
	prevSda   bool
	state     byte
	nextState byte
	n         byte
	addr      byte
	data      byte
	sl        []byte
}

func (r *mapper016Eeprom1) reset(sl []byte) {
	r.state, r.nextState = 0, 0
	r.addr, r.data = 0, 0
	r.sda, r.prevScl, r.prevSda = true, false, false
	r.sl = sl
}

func (r *mapper016Eeprom1) write(scl bool, sda bool) {
	sclRise, sclFall := !r.prevScl && scl, r.prevScl && !scl
	sdaRise, sdaFall := !r.prevSda && sda, r.prevSda && !sda
	prevScl := r.prevScl
	r.prevScl, r.prevSda = scl, sda

	if prevScl && sdaFall {
		r.state, r.n, r.addr = 1, 0, 0
		r.sda = true
		return
	} else if prevScl && sdaRise {
		r.state = 0
		r.sda = true
		return
	}

	if sclRise {
		switch r.state {
		case 1:
			if r.n < 7 {
				r.addr &^= 1 << r.n
				if sda {
					r.addr |= 1 << r.n
				}
			} else if sda {
				r.nextState = 2
				r.data = r.sl[r.addr&0x7f]
			} else {
				r.nextState = 3
			}
			r.n++
		case 2:
			if r.n < 8 {
				r.sda = r.data&(1<<r.n) != 0
			}
			r.n++
		case 3:
			if r.n < 8 {
				r.data &^= 1 << r.n
				if sda {
					r.data |= 1 << r.n
				}
			}
			r.n++
		case 4:
			r.sda = false
		case 5:
			if !sda {
				r.nextState = 0
			}
		}
	}
	if sclFall {
		switch r.state {
		case 1:
			if r.n >= 8 {
				r.state = 4
				r.sda = true
			}
		case 4:
			r.state, r.n = r.nextState, 0
			r.sda = true
		case 2:
			if r.n >= 8 {
				r.state, r.addr = 5, (r.addr+1)&0x7f
			}
		case 3:
			if r.n >= 8 {
				r.state, r.nextState = 4, 0
				r.sl[r.addr&0x7f] = r.data
				r.addr = (r.addr + 1) & 0x7f
			}
		}
	}
}

type mapper016Eeprom2 struct {
	sda       bool
	prevScl   bool
	prevSda   bool
	rw        bool
	state     byte
	nextState byte
	n         byte
	addr      byte
	data      byte
	sl        []byte
}

func (r *mapper016Eeprom2) reset(sl []byte) {
	r.state, r.nextState = 0, 0
	r.addr, r.data = 0, 0
	r.sda, r.prevScl, r.prevSda, r.rw = true, false, false, false
	r.sl = sl
}

func (r *mapper016Eeprom2) write(scl bool, sda bool) {
	sclRise, sclFall := !r.prevScl && scl, r.prevScl && !scl
	sdaRise, sdaFall := !r.prevSda && sda, r.prevSda && !sda
	prevScl := r.prevScl
	r.prevScl, r.prevSda = scl, sda

	if prevScl && sdaFall {
		r.state, r.n = 1, 0
		r.sda = true
		return
	} else if prevScl && sdaRise {
		r.state = 0
		r.sda = true
		return
	}

	if sclRise {
		switch r.state {
		case 1, 4:
			if r.n < 8 {
				r.data &^= 1 << (7 - r.n)
				if sda {
					r.data |= 1 << (7 - r.n)
				}
			}
			r.n++
		case 2:
			if r.n < 8 {
				r.addr &^= 1 << (7 - r.n)
				if sda {
					r.addr |= 1 << (7 - r.n)
				}
			}
			r.n++
		case 3:
			if r.n < 8 {
				r.sda = r.data&(1<<(7-r.n)) != 0
			}
			r.n++
		case 5:
			r.sda = false
		case 6:
			r.sda = true
		case 7:
			if !sda {
				r.nextState = 3
				r.data = r.sl[r.addr]
			}
		}
	}
	if sclFall {
		switch r.state {
		case 1:
			if r.n >= 8 {
				if (r.data & 0xa0) == 0xa0 {
					r.state, r.n = 5, 0
					r.rw, r.sda = r.data&0x01 == 0x01, true
					if r.rw {
						r.nextState = 3
						r.data = r.sl[r.addr]
					} else {
						r.nextState = 2
					}
				} else {
					r.state, r.nextState = 6, 0
					r.sda = true
				}
			}
		case 2:
			if r.n >= 8 {
				r.state, r.n = 5, 0
				r.sda = true
				if r.rw {
					r.nextState = 0
				} else {
					r.nextState = 4
				}
			}
		case 3:
			if r.n >= 8 {
				r.state = 7
				r.addr++
			}
		case 4:
			if r.n >= 8 {
				r.sl[r.addr] = r.data
				r.state, r.nextState, r.n = 5, 4, 0
				r.addr++
			}
		case 5, 7:
			r.state, r.n = r.nextState, 0
			r.sda = true
		case 6:
			r.state, r.n = 0, 0
			r.sda = true
		}
	}
}

type mapper016 struct {
	baseMapper
	patch1     bool
	eepTyp     byte
	irqEn      bool
	irqClkTyp  bool
	irqCnt     int32
	irqLatch   int32
	r0, r1, r2 byte
	eep1       mapper016Eeprom1
	eep2       mapper016Eeprom1
}

func newMapper016(bm *baseMapper) Mapper {
	return &mapper016{baseMapper: *bm}
}

func (m *mapper016) reset() {
	patch := m.sys.conf.PatchTyp
	m.patch1 = patch&0x01 != 0
	m.eepTyp = 0
	if patch&0x02 != 0 {
		m.eepTyp = 1
	} else if patch&0x04 != 0 {
		m.eepTyp = 2
	} else if patch&0x08 != 0 {
		m.eepTyp = 0xff
	}
	m.irqClkTyp = true
	if patch&0x10 != 0 {
		m.irqClkTyp = false
	}
	if patch&0x20 != 0 {
		m.sys.renderMode = RenderModePreAll
	}

	m.irqEn, m.irqClkTyp, m.irqCnt, m.irqLatch = false, true, 0, 0
	m.r0, m.r1, m.r2 = 0, 0, 0
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	switch m.eepTyp {
	case 0:
		m.eep1.reset(m.mem.wram[:])
	case 1:
		m.eep2.reset(m.mem.wram[:])
	case 2:
		m.eep1.reset(m.mem.wram[0x0100:])
		m.eep2.reset(m.mem.wram[:])
	}
}

func (m *mapper016) readLow(addr uint16) byte {
	if m.patch1 {
		return m.baseMapper.readLow(addr)
	}
	if (addr & 0x00ff) == 0x0000 {
		b := false
		switch m.eepTyp {
		case 0:
			b = m.eep1.sda
		case 1:
			b = m.eep2.sda
		case 2:
			b = m.eep1.sda && m.eep2.sda
		}
		if b {
			return 0x10
		}
	}
	return 0
}

func (m *mapper016) writeLow(addr uint16, data byte) {
	if !m.patch1 {
		m.write(addr, data)
	} else {
		m.baseMapper.writeLow(addr, data)
	}
}

func (m *mapper016) write(addr uint16, data byte) {
	if !m.patch1 {
		switch addr & 0x0f {
		case 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07:
			if m.mem.nVrom1kPage != 0 {
				m.mem.setVrom1kBank(byte(addr)&0x07, uint32(data))
			}
			if m.eepTyp == 2 {
				m.r0 = data
				m.eep1.write(data&0x08 != 0, m.r1&0x40 != 0)
			}
		case 0x08:
			m.mem.setProm16kBank(4, uint32(data))
		case 0x09:
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
		case 0x0a:
			m.irqEn = data&0x01 != 0
			m.irqCnt = m.irqLatch
			m.clearIntr()
		case 0x0b:
			m.irqLatch = (m.irqLatch & 0xff00) | int32(data)
			m.irqCnt = (m.irqCnt & 0xff00) | int32(data)
		case 0x0c:
			m.irqLatch = (m.irqLatch & 0x00ff) | (int32(data) << 8)
			m.irqCnt = (m.irqCnt & 0x00ff) | (int32(data) << 8)
		case 0x0d:
			switch m.eepTyp {
			case 0:
				m.eep1.write(data&0x20 != 0, data&0x40 != 0)
			case 1:
				m.eep2.write(data&0x20 != 0, data&0x40 != 0)
			case 2:
				m.r1 = data
				m.eep1.write(m.r0&0x08 != 0, data&0x40 != 0)
				m.eep2.write(data&0x20 != 0, data&0x40 != 0)
			}
		}
	} else {
		switch addr {
		case 0x8000, 0x8001, 0x8002, 0x8003:
			m.r0 = data & 0x01
			r0, r2 := uint32(m.r0)<<5, uint32(m.r2)<<1
			m.mem.setProm8kBank(4, r0+r2)
			m.mem.setProm8kBank(5, r0+r2+1)
		case 0x8004, 0x8005, 0x8006, 0x8007:
			m.r1 = data & 0x01
			r1 := uint32(m.r1) << 5
			m.mem.setProm8kBank(6, r1|0x1e)
			m.mem.setProm8kBank(7, r1|0x1f)
		case 0x8008:
			m.r2 = data
			r0, r1, r2 := uint32(m.r0)<<5, uint32(m.r1)<<5, uint32(m.r2)<<1
			m.mem.setProm8kBank(4, r0+r2)
			m.mem.setProm8kBank(5, r0+r2+1)
			m.mem.setProm8kBank(6, r1|0x1e)
			m.mem.setProm8kBank(7, r1|0x1f)
		case 0x8009:
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
		case 0x800a:
			m.irqEn = data&0x01 != 0
			m.irqCnt = m.irqLatch
			m.clearIntr()
		case 0x800b:
			m.irqLatch = (m.irqLatch & 0xff00) | int32(data)
		case 0x800c:
			m.irqLatch = (int32(data) << 8) | (m.irqLatch & 0x00ff)
		}
	}
}

func (m *mapper016) hSync(scanline uint16) {
	if m.irqEn && !m.irqClkTyp {
		if m.irqCnt <= 113 {
			m.setIntr()
			m.irqCnt &= 0xffff
		} else {
			m.irqCnt -= 113
		}
	}
}

func (m *mapper016) clock(nCycle int64) {
	if m.irqEn && m.irqClkTyp {
		m.irqCnt -= int32(nCycle)
		if m.irqCnt <= 0 {
			m.setIntr()
			m.irqCnt &= 0xffff
		}
	}
}

// 017

type mapper017 struct {
	baseMapper
	irqEn    bool
	irqCnt   uint32
	irqLatch uint32
}

func newMapper017(bm *baseMapper) Mapper {
	return &mapper017{baseMapper: *bm}
}

func (m *mapper017) reset() {
	m.irqEn, m.irqCnt, m.irqLatch = false, 0, 0
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper017) writeLow(addr uint16, data byte) {
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
		m.irqLatch = (m.irqLatch & 0xff00) | uint32(data)
	case 0x4503:
		m.irqLatch = (m.irqLatch & 0x00ff) | (uint32(data) << 8)
		m.irqEn, m.irqCnt = true, m.irqLatch
	case 0x4504, 0x4505, 0x4506, 0x4507:
		m.mem.setProm8kBank(byte(addr&0x07), uint32(data))
	case 0x4510, 0x4511, 0x4512, 0x4513, 0x4514, 0x4515, 0x4516, 0x4517:
		m.mem.setVrom1kBank(byte(addr&0x07), uint32(data))
	default:
		m.baseMapper.writeLow(addr, data)
	}
}

func (m *mapper017) hSync(scanline uint16) {
	if m.irqEn {
		if m.irqCnt >= 0xff8e {
			m.setIntr()
			m.irqCnt &= 0xffff
		} else {
			m.irqCnt += 113
		}
	}
}

// 018

type mapper018 struct {
	baseMapper
	irqEn    bool
	irqMode  byte
	irqCnt   int32
	irqLatch int32
	r        [11]byte
}

func newMapper018(bm *baseMapper) Mapper {
	return &mapper018{baseMapper: *bm}
}

func (m *mapper018) reset() {
	m.irqEn, m.irqMode = false, 0
	m.irqCnt, m.irqLatch = 0xffff, 0xffff
	for i := 0; i < 11; i++ {
		m.r[i] = 0
	}
	m.r[2], m.r[3] = byte(m.mem.nProm8kPage-2), byte(m.mem.nProm8kPage-1)
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
}

func (m *mapper018) write(addr uint16, data byte) {
	switch addr {
	case 0x8000, 0x8002, 0x9000:
		i := (byte(addr>>11) - 0x10) | (byte(addr&0x02) >> 1)
		m.r[i] = (m.r[i] & 0xf0) | (data & 0x0f)
		m.mem.setProm8kBank(i+4, uint32(m.r[i]))
	case 0x8001, 0x8003, 0x9001:
		i := (byte(addr>>11) - 0x10) | (byte(addr&0x02) >> 1)
		m.r[i] = (m.r[i] & 0x0f) | (data << 4)
		m.mem.setProm8kBank(i+4, uint32(m.r[i]))
	case 0xa000, 0xa002, 0xb000, 0xb002, 0xc000, 0xc002, 0xd000, 0xd002:
		i := (byte(addr>>11) - 0x14) | (byte(addr&0x02) >> 1)
		m.r[i+3] = (m.r[i+3] & 0xf0) | (data & 0x0f)
		m.mem.setVrom1kBank(i, uint32(m.r[i+3]))
	case 0xa001, 0xa003, 0xb001, 0xb003, 0xc001, 0xc003, 0xd001, 0xd003:
		i := (byte(addr>>11) - 0x14) | (byte(addr&0x02) >> 1)
		m.r[i+3] = (m.r[i+3] & 0x0f) | (data << 4)
		m.mem.setVrom1kBank(i, uint32(m.r[i+3]))
	case 0xe000, 0xe001, 0xe002, 0xe003:
		i := addr & 0x03
		m.irqLatch = (m.irqLatch &^ (0x0f << i)) | (int32(data&0x0f) << (i << 2))
	case 0xf000:
		m.irqCnt = m.irqLatch
	case 0xf001:
		m.irqMode = (data >> 1) & 0x07
		m.irqEn = data&0x01 != 0
		m.clearIntr()
	case 0xf002:
		switch data & 0x03 {
		case 0x00:
			m.mem.setVramMirror(memVramMirrorH)
		case 0x01:
			m.mem.setVramMirror(memVramMirrorV)
		default:
			m.mem.setVramMirror(memVramMirror4L)
		}
	}
}

func (m *mapper018) clock(nCycle int64) {
	if (!m.irqEn) && m.irqCnt == 0 {
		return
	}
	bIrq, prevCnt := false, m.irqCnt
	m.irqCnt -= int32(nCycle)
	switch m.irqMode {
	case 0:
		if m.irqCnt <= 0 {
			bIrq = true
		}
	case 1:
		if m.irqCnt&0xf000 != prevCnt&0xf000 {
			bIrq = true
		}
	case 2, 3:
		if m.irqCnt&0xff00 != prevCnt&0xff00 {
			bIrq = true
		}
	case 4, 5, 6, 7:
		if m.irqCnt&0xfff0 != prevCnt&0xfff0 {
			bIrq = true
		}
	}
	if bIrq {
		m.irqEn, m.irqCnt = false, 0
		m.setIntr()
	}
}

// 019

type mapper019 struct {
	baseMapper
	patchTyp byte
	irqEn    bool
	irqCnt   uint16
	r0, r1   byte
}

func newMapper019(bm *baseMapper) Mapper {
	return &mapper019{baseMapper: *bm}
}

func (m *mapper019) reset() {
	patch := m.sys.conf.PatchTyp
	if patch&0x01 != 0 {
		m.patchTyp = 1
	} else if patch&0x02 != 0 {
		m.patchTyp = 2
	} else if patch&0x04 != 0 {
		m.patchTyp = 3
	}
	m.irqEn, m.irqCnt = false, 0
	m.r0, m.r1 = 0, 0
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	if m.mem.nVrom1kPage >= 8 {
		m.mem.setVrom8kBank((m.mem.nVrom1kPage >> 3) - 1)
	}
}

func (m *mapper019) readLow(addr uint16) byte {
	switch addr & 0xf800 {
	case 0x4800:
		if addr == 0x4800 {
			data := m.mem.wram[m.r1&0x7f]
			if m.r1&0x80 != 0 {
				m.r1 = (m.r1 + 1) | 0x80
			}
			return data
		}
	case 0x5000:
		return byte(m.irqCnt)
	case 0x5800:
		return byte((m.irqCnt >> 8) & 0x7f)
	case 0x6000, 0x6800, 0x7000, 0x7800:
		return m.cpuBanks[addr>>13][addr&0x1fff]
	}
	return byte(addr >> 8)
}

func (m *mapper019) writeLow(addr uint16, data byte) {
	switch addr & 0xf800 {
	case 0x4800:
		if addr == 0x4800 {
			m.mem.wram[m.r1&0x7f] = data
			if m.r1&0x80 != 0 {
				m.r1 = (m.r1 + 1) | 0x80
			}
		}
	case 0x5000:
		m.irqCnt = (m.irqCnt & 0xff00) | uint16(data)
		m.clearIntr()
	case 0x5800:
		m.irqCnt = (m.irqCnt & 0x00ff) | (uint16(data&0x7f) << 8)
		m.irqEn = data&0x80 != 0
		m.clearIntr()
	case 0x6000, 0x6800, 0x7000, 0x7800:
		m.cpuBanks[addr>>13][addr&0x1fff] = data
	}
}

func (m *mapper019) write(addr uint16, data byte) {
	a := addr & 0xf800
	switch a {
	case 0x8000, 0x8800, 0x9000, 0x9800, 0xa000, 0xa800, 0xb000, 0xb800:
		i := byte(a>>11) & 0x07
		if data < 0xe0 || m.r0&(0x40<<(i>>2)) != 0 {
			m.mem.setVrom1kBank(i, uint32(data))
		} else {
			m.mem.setCram1kBank(i, uint32(data&0x1f))
		}
	case 0xc000, 0xc800, 0xd000, 0xd800:
		if m.patchTyp == 0 {
			i := (byte(a>>11) & 0x03) | 0x08
			if data <= 0xdf {
				m.mem.setVrom1kBank(i, uint32(data))
			} else {
				m.mem.setVram1kBank(i, uint32(data&0x01))
			}
		}
	case 0xe000:
		m.mem.setProm8kBank(4, uint32(data&0x3f))
		switch m.patchTyp {
		case 2:
			if data&0x40 != 0 {
				m.mem.setVramMirror(memVramMirrorV)
			} else {
				m.mem.setVramMirror(memVramMirror4L)
			}
		case 3:
			if data&0x80 != 0 {
				m.mem.setVramMirror(memVramMirrorH)
			} else {
				m.mem.setVramMirror(memVramMirrorV)
			}
		}
	case 0xe800:
		m.r0 = data & 0xc0
		m.mem.setProm8kBank(5, uint32(data&0x3f))
	case 0xf000:
		m.mem.setProm8kBank(5, uint32(data&0x3f))
	case 0xf800:
		if addr == 0xf800 {
			m.r1 = data
		}
	}
}

func (m *mapper019) clock(nCycle int64) {
	if m.irqEn {
		m.irqCnt += uint16(nCycle)
		if m.irqCnt >= 0x7fff {
			m.irqEn, m.irqCnt = false, 0x7fff
			m.setIntr()
		}
	}
}

// 021

type mapper021 struct {
	baseMapper
	irqEn    byte
	irqCnt   byte
	irqLatch byte
	irqClk   int16
	r        [9]byte
}

func newMapper021(bm *baseMapper) Mapper {
	return &mapper021{baseMapper: *bm}
}

func (m *mapper021) reset() {
	m.irqEn, m.irqCnt, m.irqLatch, m.irqClk = 0, 0, 0, 0
	for i := byte(0); i < 8; i++ {
		m.r[i] = i
	}
	m.r[8] = 0
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
}

func (m *mapper021) write(addr uint16, data byte) {
	addr &= 0xf0cf
	switch addr {
	case 0x8000:
		if m.r[8]&0x02 == 0 {
			m.mem.setProm8kBank(4, uint32(data))
		} else {
			m.mem.setProm8kBank(6, uint32(data))
		}
	case 0x9000:
		if data != 0xff {
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
	case 0x9002, 0x9080:
		m.r[8] = data
	case 0xa000:
		m.mem.setProm8kBank(5, uint32(data))
	case 0xb000, 0xb001, 0xb004, 0xb080, 0xc000, 0xc001, 0xc004, 0xc080,
		0xd000, 0xd001, 0xd004, 0xd080, 0xe000, 0xe001, 0xe004, 0xe080:
		i := (byte(addr>>11) - 0x16)
		switch addr & 0x00ff {
		case 0x01, 0x04, 0x80:
			i |= 0x01
		}
		m.r[i] = (m.r[i] & 0xf0) | (data & 0x0f)
		m.mem.setVrom1kBank(i, uint32(m.r[i]))
	case 0xb002, 0xb040, 0xb003, 0xb006, 0xb0c0, 0xc002, 0xc040, 0xc003, 0xc006, 0xc0c0,
		0xd002, 0xd040, 0xd003, 0xd006, 0xd0c0, 0xe002, 0xe040, 0xe003, 0xe006, 0xe0c0:
		i := (byte(addr>>11) - 0x16)
		switch addr & 0x00ff {
		case 0x03, 0x06, 0xc0:
			i |= 0x01
		}
		m.r[i] = (m.r[i] & 0x0f) | (data << 4)
		m.mem.setVrom1kBank(i, uint32(m.r[i]))
	case 0xf000:
		m.irqLatch = (m.irqLatch & 0xf0) | (data & 0x0f)
	case 0xf002, 0xf040:
		m.irqLatch = (m.irqLatch & 0x0f) | (data << 4)
	case 0xf003, 0xf006, 0xf0c0:
		m.irqEn, m.irqClk = (m.irqEn&0x01)*3, 0
		m.clearIntr()
	case 0xf004, 0xf080:
		m.irqEn = data & 0x03
		if m.irqEn&0x02 != 0 {
			m.irqCnt, m.irqClk = m.irqLatch, 0
		}
		m.clearIntr()
	}
}

func (m *mapper021) clock(nCycle int64) {
	if m.irqEn&0x02 != 0 {
		m.irqClk -= int16(nCycle)
		if m.irqClk < 0 {
			m.irqClk += 114
			if m.irqCnt == 0xff {
				m.irqCnt = m.irqLatch
				m.setIntr()
			} else {
				m.irqCnt++
			}
		}
	}
}

// 022

type mapper022 struct {
	baseMapper
}

func newMapper022(bm *baseMapper) Mapper {
	return &mapper022{baseMapper: *bm}
}

func (m *mapper022) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
}

func (m *mapper022) write(addr uint16, data byte) {
	switch addr {
	case 0x8000:
		m.mem.setProm8kBank(4, uint32(data))
	case 0x9000:
		switch data & 0x03 {
		case 0x00:
			m.mem.setVramMirror(memVramMirrorV)
		case 0x01:
			m.mem.setVramMirror(memVramMirrorH)
		case 0x02:
			m.mem.setVramMirror(memVramMirror4H)
		case 0x03:
			m.mem.setVramMirror(memVramMirror4L)
		}
	case 0xa000:
		m.mem.setProm8kBank(5, uint32(data))
	case 0xb000, 0xb001, 0xc000, 0xc001, 0xd000, 0xd001, 0xe000, 0xe001:
		m.mem.setVrom1kBank((byte(addr>>11)-0x16)|(byte(addr)&0x01), uint32(data>>1))
	}
}

// 023

type mapper023 struct {
	baseMapper
	addrMask uint16
	irqEn    byte
	irqCnt   byte
	irqLatch byte
	irqClk   uint16
	r        [9]byte
}

func newMapper023(bm *baseMapper) Mapper {
	return &mapper023{baseMapper: *bm}
}

func (m *mapper023) reset() {
	m.addrMask = 0xffff
	if m.sys.conf.PatchTyp&0x01 != 0 {
		m.addrMask = 0xf00c
	}
	m.irqEn, m.irqCnt, m.irqLatch, m.irqClk = 0, 0, 0, 0
	for i := byte(0); i < 8; i++ {
		m.r[i] = i
	}
	m.r[8] = 0
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	m.mem.setVrom8kBank(0)
}

func (m *mapper023) write(addr uint16, data byte) {
	addr &= m.addrMask
	switch addr {
	case 0x8000, 0x8004, 0x8008, 0x800c:
		if m.r[8] == 0 {
			m.mem.setProm8kBank(4, uint32(data))
		} else {
			m.mem.setProm8kBank(6, uint32(data))
		}
	case 0x9000:
		if data != 0xff {
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
	case 0x9008:
		m.r[8] = data & 0x02
	case 0xa000, 0xa004, 0xa008, 0xa00c:
		m.mem.setProm8kBank(5, uint32(data))
	case 0xb000, 0xb002, 0xb008, 0xc000, 0xc002, 0xc008,
		0xd000, 0xd002, 0xd008, 0xe000, 0xe002, 0xe008:
		i := (byte(addr>>11) - 0x16) | (byte(addr&0x02) >> 1) | (byte(addr&0x08) >> 3)
		m.r[i] = (m.r[i] & 0xf0) | (data & 0x0f)
		m.mem.setVrom1kBank(i, uint32(m.r[i]))
	case 0xb001, 0xb004, 0xb003, 0xb00c, 0xc001, 0xc004, 0xc003, 0xc00c,
		0xd001, 0xd004, 0xd003, 0xd00c, 0xe001, 0xe004, 0xe003, 0xe00c:
		i := (byte(addr>>11) - 0x16) | (byte(addr&0x02) >> 1) | (byte(addr&0x08) >> 3)
		m.r[i] = (m.r[i] & 0x0f) | (data << 4)
		m.mem.setVrom1kBank(i, uint32(m.r[i]))
	case 0xf000:
		m.irqLatch = (m.irqLatch & 0xf0) | (data & 0x0f)
		m.clearIntr()
	case 0xf004:
		m.irqLatch = (m.irqLatch & 0x0f) | (data << 4)
		m.clearIntr()
	case 0xf008:
		m.irqEn, m.irqCnt, m.irqClk = data&0x03, m.irqLatch, 0
		m.clearIntr()
	case 0xf00c:
		m.irqEn = (m.irqEn & 0x01) * 3
		m.clearIntr()
	}
}

func (m *mapper023) clock(nCycle int64) {
	if m.irqEn&0x02 != 0 {
		m.irqClk += uint16(nCycle) * 3
		for m.irqClk >= 341 {
			m.irqClk -= 341
			m.irqCnt++
			if m.irqCnt == 0 {
				m.irqCnt = m.irqLatch
				m.setIntr()
			}
		}
	}
}

// 024

type mapper024 struct {
	baseMapper
	irqEn    byte
	irqCnt   byte
	irqLatch byte
	irqClk   uint16
}

func newMapper024(bm *baseMapper) Mapper {
	return &mapper024{baseMapper: *bm}
}

func (m *mapper024) reset() {
	m.irqEn, m.irqCnt, m.irqLatch, m.irqClk = 0, 0, 0, 0
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
	m.sys.renderMode = RenderModePost
}

func (m *mapper024) write(addr uint16, data byte) {
	addr &= 0xf003
	switch addr {
	case 0x8000:
		m.mem.setProm16kBank(4, uint32(data))
	case 0xb003:
		switch data & 0x0c {
		case 0x00:
			m.mem.setVramMirror(memVramMirrorV)
		case 0x04:
			m.mem.setVramMirror(memVramMirrorH)
		case 0x08:
			m.mem.setVramMirror(memVramMirror4L)
		case 0x0c:
			m.mem.setVramMirror(memVramMirror4H)
		}
	case 0xc000:
		m.mem.setProm8kBank(6, uint32(data))
	case 0xd000, 0xd001, 0xd002, 0xd003:
		i := byte(addr)
		m.mem.setVrom1kBank(((i&0x01)<<1)|((i&0x02)>>1), uint32(data))
	case 0xe000, 0xe001, 0xe002, 0xe003:
		i := byte(addr)
		m.mem.setVrom1kBank(0x04|((i&0x01)<<1)|((i&0x02)>>1), uint32(data))
	case 0xf000:
		m.irqLatch = data
	case 0xf001:
		m.irqEn = data & 0x03
		if m.irqEn&0x02 != 0 {
			m.irqCnt, m.irqClk = m.irqLatch, 0
		}
		m.clearIntr()
	case 0xf002:
		m.irqEn = (m.irqEn & 0x01) * 3
		m.clearIntr()
	}
}

func (m *mapper024) clock(nCycle int64) {
	if m.irqEn&0x02 != 0 {
		m.irqClk += uint16(nCycle)
		if m.irqClk >= 114 {
			m.irqClk -= 114
			if m.irqCnt == 0xff {
				m.irqCnt = m.irqLatch
				m.setIntr()
			} else {
				m.irqCnt++
			}
		}
	}
}

// 025

type mapper025 struct {
	baseMapper
	irqEn      byte
	irqCnt     byte
	irqLatch   byte
	irqClk     uint16
	r1, r2, r3 byte
	r          [8]byte
}

func newMapper025(bm *baseMapper) Mapper {
	return &mapper025{baseMapper: *bm}
}

func (m *mapper025) reset() {
	m.irqEn, m.irqCnt, m.irqLatch, m.irqClk = 0, 0, 0, 0
	for i := 0; i < 8; i++ {
		m.r[i] = 0
	}
	m.r1, m.r2, m.r3 = 0, byte(m.mem.nProm8kPage-2), 0
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
}

func (m *mapper025) write(addr uint16, data byte) {
	switch addr & 0xf000 {
	case 0x8000:
		if m.r3&0x02 == 0 {
			m.r1 = data
			m.mem.setProm8kBank(4, uint32(data))
		} else {
			m.r2 = data
			m.mem.setProm8kBank(6, uint32(data))
		}
	case 0xa000:
		m.mem.setProm8kBank(5, uint32(data))
	}
	addr &= 0xf00f
	switch addr {
	case 0x9000:
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
	case 0x9001, 0x9004:
		if m.r3&0x02 != data&0x02 {
			m.r1, m.r2 = m.r2, m.r1
			m.mem.setProm8kBank(4, uint32(m.r1))
			m.mem.setProm8kBank(6, uint32(m.r2))
		}
		m.r3 = data
	case 0xb000, 0xb001, 0xb004, 0xc000, 0xc001, 0xc004,
		0xd000, 0xd001, 0xd004, 0xe000, 0xe001, 0xe004:
		i := (byte(addr>>11) - 0x16) | byte(addr&0x01) | (byte(addr&0x04) >> 2)
		m.r[i] = (m.r[i] & 0xf0) | (data & 0x0f)
		m.mem.setVrom1kBank(i, uint32(m.r[i]))
	case 0xb002, 0xb008, 0xb003, 0xb00c, 0xc002, 0xc008, 0xc003, 0xc00c,
		0xd002, 0xd008, 0xd003, 0xd00c, 0xe002, 0xe008, 0xe003, 0xe00c:
		i := (byte(addr>>11) - 0x16) | byte(addr&0x01) | (byte(addr&0x04) >> 2)
		m.r[i] = (m.r[i] & 0x0f) | (data << 4)
		m.mem.setVrom1kBank(i, uint32(m.r[i]))
	case 0xf000:
		m.irqLatch = (m.irqLatch & 0xf0) | (data & 0x0f)
		m.clearIntr()
	case 0xf002, 0xf008:
		m.irqLatch = (m.irqLatch & 0x0f) | (data << 4)
		m.clearIntr()
	case 0xf001, 0xf004:
		m.irqEn, m.irqCnt, m.irqClk = data&0x03, m.irqLatch, 0
		m.clearIntr()
	case 0xf003, 0xf00c:
		m.irqEn = (m.irqEn & 0x01) * 3
		m.clearIntr()
	}
}

func (m *mapper025) clock(nCycle int64) {
	if m.irqEn&0x02 != 0 {
		m.irqClk += uint16(nCycle) * 3
		for m.irqClk >= 341 {
			m.irqClk -= 341
			m.irqCnt++
			if m.irqCnt == 0 {
				m.irqCnt = m.irqLatch
				m.setIntr()
			}
		}
	}
}

// 026

type mapper026 struct {
	baseMapper
	irqEn    byte
	irqCnt   byte
	irqLatch byte
	irqClk   uint16
}

func newMapper026(bm *baseMapper) Mapper {
	return &mapper026{baseMapper: *bm}
}

func (m *mapper026) reset() {
	m.irqEn, m.irqCnt, m.irqLatch, m.irqClk = 0, 0, 0, 0
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper026) write(addr uint16, data byte) {
	addr &= 0xf003
	switch addr {
	case 0x8000:
		m.mem.setProm16kBank(4, uint32(data))
	case 0xb003:
		switch data & 0x7f {
		case 0x08, 0x2c:
			m.mem.setVramMirror(memVramMirror4H)
		case 0x020:
			m.mem.setVramMirror(memVramMirrorV)
		case 0x024:
			m.mem.setVramMirror(memVramMirrorH)
		case 0x028:
			m.mem.setVramMirror(memVramMirror4L)
		}
	case 0xc000:
		m.mem.setProm8kBank(6, uint32(data))
	case 0xd000, 0xd001, 0xd002, 0xd003:
		i := byte(addr)
		m.mem.setVrom1kBank(((i&0x01)<<1)|((i&0x02)>>1), uint32(data))
	case 0xe000, 0xe001, 0xe002, 0xe003:
		i := byte(addr)
		m.mem.setVrom1kBank(0x04|((i&0x01)<<1)|((i&0x02)>>1), uint32(data))
	case 0xf000:
		m.irqLatch = data
	case 0xf001:
		m.irqEn = (m.irqEn & 0x01) * 3
		m.clearIntr()
	case 0xf002:
		m.irqEn = data & 0x03
		if m.irqEn&0x02 != 0 {
			m.irqCnt, m.irqClk = m.irqLatch, 0
		}
		m.clearIntr()
	}
}

func (m *mapper026) clock(nCycle int64) {
	if m.irqEn&0x02 != 0 {
		m.irqClk += uint16(nCycle)
		if m.irqClk >= 114 {
			m.irqClk -= 114
			if m.irqCnt >= 0xff {
				m.irqCnt = m.irqLatch
				m.setIntr()
			} else {
				m.irqCnt++
			}
		}
	}
}

// 027

type mapper027 struct {
	baseMapper
	irqEn    byte
	irqCnt   byte
	irqLatch byte
	r        [9]byte
}

func newMapper027(bm *baseMapper) Mapper {
	return &mapper027{baseMapper: *bm}
}

func (m *mapper027) reset() {
	m.irqEn, m.irqCnt, m.irqLatch = 0, 0, 0
	for i := byte(0); i < 8; i++ {
		m.r[i] = i
	}
	m.r[8] = 0
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
}

func (m *mapper027) write(addr uint16, data byte) {
	addr &= 0xf0cf
	switch addr {
	case 0x8000:
		if m.r[8] == 0 {
			m.mem.setProm8kBank(4, uint32(data))
		} else {
			m.mem.setProm8kBank(6, uint32(data))
		}
	case 0x9000:
		if data != 0xff {
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
	case 0x9002, 0x9080:
		m.r[8] = data
	case 0xa000:
		m.mem.setProm8kBank(5, uint32(data))
	case 0xb000, 0xb002, 0xc000, 0xc002, 0xd000, 0xd002, 0xe000, 0xe002:
		i := (byte(addr>>11) - 0x16) | (byte(addr&0x02) >> 1)
		m.r[i] = (m.r[i] & 0xf0) | (data & 0x0f)
		m.mem.setVrom1kBank(i, uint32(m.r[i]))
	case 0xb001, 0xb003, 0xc001, 0xc003, 0xd001, 0xd003, 0xe001, 0xe003:
		i := (byte(addr>>11) - 0x16) | (byte(addr&0x02) >> 1)
		m.r[i] = (m.r[i] & 0x0f) | (data << 4)
		m.mem.setVrom1kBank(i, uint32(m.r[i]))
	case 0xf000:
		m.irqLatch = (m.irqLatch & 0xf0) | (data & 0x0f)
		m.clearIntr()
	case 0xf001:
		m.irqLatch = (m.irqLatch & 0x0f) | (data << 4)
		m.clearIntr()
	case 0xf002:
		m.irqEn = (m.irqEn & 0x01) * 3
		m.clearIntr()
	case 0xf003:
		m.irqEn = m.irqEn & 0x03
		if m.irqEn&0x02 != 0 {
			m.irqCnt = m.irqLatch
		}
		m.clearIntr()
	}
}

func (m *mapper027) hSync(scanline uint16) {
	if m.irqEn&0x02 != 0 {
		if m.irqCnt == 0xff {
			m.irqCnt = m.irqLatch
			m.setIntr()
		} else {
			m.irqCnt++
		}
	}
}
