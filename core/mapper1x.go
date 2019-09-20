package core

// 0x10

type mapper10Eeprom1 struct {
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

func (r *mapper10Eeprom1) reset(sl []byte) {
	r.state, r.nextState = 0, 0
	r.addr, r.data = 0, 0
	r.sda, r.prevScl, r.prevSda = true, false, false
	r.sl = sl
}

func (r *mapper10Eeprom1) write(scl bool, sda bool) {
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

type mapper10Eeprom2 struct {
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

func (r *mapper10Eeprom2) reset(sl []byte) {
	r.state, r.nextState = 0, 0
	r.addr, r.data = 0, 0
	r.sda, r.prevScl, r.prevSda, r.rw = true, false, false, false
	r.sl = sl
}

func (r *mapper10Eeprom2) write(scl bool, sda bool) {
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

type mapper10 struct {
	baseMapper
	patch1     bool
	eepTyp     byte
	r0, r1, r2 byte
	irqEn      bool
	irqClkTyp  bool
	irqCnt     int32
	irqLatch   int32
	eep1       mapper10Eeprom1
	eep2       mapper10Eeprom1
}

func newMapper10(bm *baseMapper) Mapper {
	return &mapper10{baseMapper: *bm}
}

func (m *mapper10) reset() {
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

	m.r0, m.r1, m.r2 = 0, 0, 0
	m.irqEn, m.irqClkTyp, m.irqCnt, m.irqLatch = false, true, 0, 0
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

func (m *mapper10) readLow(addr uint16) byte {
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

func (m *mapper10) writeLow(addr uint16, data byte) {
	if !m.patch1 {
		m.write(addr, data)
	} else {
		m.baseMapper.writeLow(addr, data)
	}
}

func (m *mapper10) write(addr uint16, data byte) {
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
			m.sys.cpu.intr &^= cpuIntrTypMapper
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
			m.sys.cpu.intr &^= cpuIntrTypMapper
		case 0x800b:
			m.irqLatch = (m.irqLatch & 0xff00) | int32(data)
		case 0x800c:
			m.irqLatch = (int32(data) << 8) | (m.irqLatch & 0x00ff)
		}
	}
}

func (m *mapper10) hSync(scanline uint16) {
	if m.irqEn && !m.irqClkTyp {
		if m.irqCnt <= 113 {
			m.sys.cpu.intr |= cpuIntrTypMapper
			m.irqCnt &= 0xffff
		} else {
			m.irqCnt -= 113
		}
	}
}

func (m *mapper10) clock(nCycle int64) {
	if m.irqEn && m.irqClkTyp {
		m.irqCnt -= int32(nCycle)
		if m.irqCnt <= 0 {
			m.sys.cpu.intr |= cpuIntrTypMapper
			m.irqCnt &= 0xffff
		}
	}
}

// 0x11

type mapper11 struct {
	baseMapper
}

func newMapper11(bm *baseMapper) Mapper {
	return &mapper11{baseMapper: *bm}
}

func (m *mapper11) reset() {
}

// 0x12

type mapper12 struct {
	baseMapper
}

func newMapper12(bm *baseMapper) Mapper {
	return &mapper12{baseMapper: *bm}
}

func (m *mapper12) reset() {
}

// 0x13

type mapper13 struct {
	baseMapper
}

func newMapper13(bm *baseMapper) Mapper {
	return &mapper13{baseMapper: *bm}
}

func (m *mapper13) reset() {
}

// 0x14

type mapper14 struct {
	baseMapper
}

func newMapper14(bm *baseMapper) Mapper {
	return &mapper14{baseMapper: *bm}
}

func (m *mapper14) reset() {
}

// 0x15

type mapper15 struct {
	baseMapper
}

func newMapper15(bm *baseMapper) Mapper {
	return &mapper15{baseMapper: *bm}
}

func (m *mapper15) reset() {
}

// 0x16

type mapper16 struct {
	baseMapper
}

func newMapper16(bm *baseMapper) Mapper {
	return &mapper16{baseMapper: *bm}
}

func (m *mapper16) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
}

func (m *mapper16) write(addr uint16, data byte) {
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

// 0x17

type mapper17 struct {
	baseMapper
	r        [9]byte
	irqEn    byte
	irqCnt   byte
	irqLatch byte
	irqClk   uint16
	addrMask uint16
}

func newMapper17(bm *baseMapper) Mapper {
	return &mapper17{baseMapper: *bm}
}

func (m *mapper17) reset() {
	for i := byte(0); i < 8; i++ {
		m.r[i] = i
	}
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
	m.mem.setVrom8kBank(0)
}

func (m *mapper17) write(addr uint16, data byte) {
	switch addr {
	case 0x8000, 0x8004, 0x8008, 0x800c:
		if m.r[8] != 0 {
			m.mem.setProm8kBank(6, uint32(data))
		} else {
			m.mem.setProm8kBank(4, uint32(data))
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
		m.r[i] = (m.r[i] & 0x0f) | ((data & 0x0f) << 4)
		m.mem.setVrom1kBank(i, uint32(m.r[i]))
	case 0xf000:
		m.irqLatch = (m.irqLatch & 0xf0) | (data & 0x0f)
		m.sys.cpu.intr &^= cpuIntrTypMapper
	case 0xf004:
		m.irqLatch = (m.irqLatch & 0x0f) | ((data & 0x0f) << 4)
		m.sys.cpu.intr &^= cpuIntrTypMapper
	case 0xf008:
		m.irqEn, m.irqCnt, m.irqClk = data&0x03, m.irqLatch, 0
		m.sys.cpu.intr &^= cpuIntrTypMapper
	case 0xf00c:
		m.irqEn = (m.irqEn & 0x01) * 3
		m.sys.cpu.intr &^= cpuIntrTypMapper
	}
}

func (m *mapper17) clock(nCycle int64) {
	if m.irqEn&0x02 != 0 {
		m.irqClk += uint16(nCycle * 3)
		for m.irqClk >= 341 {
			m.irqClk -= 341
			m.irqCnt++
			if m.irqCnt == 0 {
				m.irqCnt = m.irqLatch
				m.sys.cpu.intr |= cpuIntrTypMapper
			}
		}
	}
}

// 0x18

type mapper18 struct {
	baseMapper
}

func newMapper18(bm *baseMapper) Mapper {
	return &mapper18{baseMapper: *bm}
}

func (m *mapper18) reset() {
}

// 0x19

type mapper19 struct {
	baseMapper
}

func newMapper19(bm *baseMapper) Mapper {
	return &mapper19{baseMapper: *bm}
}

func (m *mapper19) reset() {
}

// 0x1a

type mapper1a struct {
	baseMapper
}

func newMapper1a(bm *baseMapper) Mapper {
	return &mapper1a{baseMapper: *bm}
}

func (m *mapper1a) reset() {
}

// 0x1b

type mapper1b struct {
	baseMapper
}

func newMapper1b(bm *baseMapper) Mapper {
	return &mapper1b{baseMapper: *bm}
}

func (m *mapper1b) reset() {
}
