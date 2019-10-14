package core

// 225

type mapper225 struct {
	baseMapper
}

func newMapper225(bm *baseMapper) Mapper {
	return &mapper225{baseMapper: *bm}
}

func (m *mapper225) reset() {
	m.mem.setProm32kBank(0)
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper225) write(addr uint16, data byte) {
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

// 226

type mapper226 struct {
	baseMapper
	r0, r1 byte
}

func newMapper226(bm *baseMapper) Mapper {
	return &mapper226{baseMapper: *bm}
}

func (m *mapper226) reset() {
	m.mem.setProm32kBank(0)
}

func (m *mapper226) write(addr uint16, data byte) {
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

// 227

type mapper227 struct {
	baseMapper
}

func newMapper227(bm *baseMapper) Mapper {
	return &mapper227{baseMapper: *bm}
}

func (m *mapper227) reset() {
	m.mem.setProm32kBank4(0, 1, 0, 1)
}

func (m *mapper227) write(addr uint16, data byte) {
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

// 228

type mapper228 struct {
	baseMapper
}

func newMapper228(bm *baseMapper) Mapper {
	return &mapper228{baseMapper: *bm}
}

func (m *mapper228) reset() {
	m.mem.setProm32kBank(0)
	m.mem.setVrom8kBank(0)
}

func (m *mapper228) write(addr uint16, data byte) {
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

// 229

type mapper229 struct {
	baseMapper
}

func newMapper229(bm *baseMapper) Mapper {
	return &mapper229{baseMapper: *bm}
}

func (m *mapper229) reset() {
	m.mem.setProm32kBank(0)
	m.mem.setVrom8kBank(0)
}

func (m *mapper229) write(addr uint16, data byte) {
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

// 230

type mapper230 struct {
	baseMapper
	idx bool
}

func newMapper230(bm *baseMapper) Mapper {
	return &mapper230{baseMapper: *bm, idx: true}
}

func (m *mapper230) reset() {
	m.idx = !m.idx
	if m.idx {
		m.mem.setProm32kBank4(0, 1, 14, 15)
	} else {
		m.mem.setProm32kBank4(16, 17, m.nProm8kPage-2, m.nProm8kPage-1)
	}
}

func (m *mapper230) write(addr uint16, data byte) {
	b := uint32(data)
	if m.idx {
		m.mem.setProm16kBank(4, b&0x07)
	} else {
		if data&0x20 != 0 {
			m.mem.setProm16kBank(4, (b&0x1f)+8)
			m.mem.setProm16kBank(6, (b&0x1f)+8)
		} else {
			m.mem.setProm16kBank(4, (b&0x1e)+8)
			m.mem.setProm16kBank(6, (b&0x1e)+9)
		}
		if data&0x40 != 0 {
			m.mem.setVramMirror(memVramMirrorV)
		} else {
			m.mem.setVramMirror(memVramMirrorH)
		}
	}
}

// 231

type mapper231 struct {
	baseMapper
}

func newMapper231(bm *baseMapper) Mapper {
	return &mapper231{baseMapper: *bm}
}

func (m *mapper231) reset() {
	m.mem.setProm32kBank(0)
	if m.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper231) write(addr uint16, data byte) {
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

// 232

type mapper232 struct {
	baseMapper
	r0, r1 byte
}

func newMapper232(bm *baseMapper) Mapper {
	return &mapper232{baseMapper: *bm}
}

func (m *mapper232) reset() {
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
	m.r0 = 0x0c
}

func (m *mapper232) writeLow(addr uint16, data byte) {
	if addr >= 0x6000 {
		m.write(addr, data)
	}
}

func (m *mapper232) write(addr uint16, data byte) {
	if addr < 0xa000 {
		m.r0 = (data & 0x18) >> 1
	} else {
		m.r1 = data & 0x03
	}
	m.mem.setProm16kBank(4, uint32(m.r0)|uint32(m.r1))
	m.mem.setProm16kBank(6, uint32(m.r0)|0x03)
}

// 233

type mapper233 struct {
	baseMapper
}

func newMapper233(bm *baseMapper) Mapper {
	return &mapper233{baseMapper: *bm}
}

func (m *mapper233) reset() {
	m.mem.setProm32kBank(0)
}

func (m *mapper233) write(addr uint16, data byte) {
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

// 234

type mapper234 struct {
	baseMapper
	r0, r1 byte
}

func newMapper234(bm *baseMapper) Mapper {
	return &mapper234{baseMapper: *bm}
}

func (m *mapper234) setBank() {
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

func (m *mapper234) reset() {
	m.mem.setProm32kBank(0)
}

func (m *mapper234) read(addr uint16) byte {
	data := m.cpuBanks[addr>>13][addr&0x1fff]
	m.write(addr, data)
	return data
}

func (m *mapper234) write(addr uint16, data byte) {
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

// 235

type mapper235 struct {
	baseMapper
}

func newMapper235(bm *baseMapper) Mapper {
	return &mapper235{baseMapper: *bm}
}

func (m *mapper235) reset() {
	for i := 0; i < 0x2000; i++ {
		m.mem.dram[i] = 0xff
	}
	m.mem.setProm32kBank(0)
}

func (m *mapper235) write(addr uint16, data byte) {
	i := uint32(((addr & 0x0300) >> 3) | (addr & 0x001f))
	b := false
	switch m.nProm8kPage {
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

// 236

type mapper236 struct {
	baseMapper
	bank, mode byte
}

func newMapper236(bm *baseMapper) Mapper {
	return &mapper236{baseMapper: *bm}
}

func (m *mapper236) reset() {
	m.mem.setProm32kBank4(0, 1, m.nProm8kPage-2, m.nProm8kPage-1)
}

func (m *mapper236) write(addr uint16, data byte) {
	a8 := byte(addr)
	if addr >= 0x8000 && addr < 0xc000 {
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

// 237

type mapper237 struct {
	baseMapper
	r    byte
	data byte
	addr uint16
}

func newMapper237(bm *baseMapper) Mapper {
	return &mapper237{baseMapper: *bm}
}

func (m *mapper237) reset() {
	m.mem.setProm16kBank(4, 0)
	m.mem.setProm16kBank(6, 7)
	m.r = (m.r + 1) & 0x03
	m.data, m.addr = 0, 0
}

func (m *mapper237) read(addr uint16) byte {
	if addr == 0xc000 {
		m.cpuBanks[6][0] = m.r + 1
		return m.r + 1
	}
	return m.cpuBanks[addr>>13][addr&0x1fff]
}

func (m *mapper237) write(addr uint16, data byte) {
	if m.addr&0x02 != 0 {
		d34, a2 := uint32(m.data&0x18), uint32((m.addr&0x04)<<3)
		m.mem.setProm16kBank(4, uint32(data)|d34|a2)
		m.mem.setProm16kBank(6, 0x07|d34|a2)
		return
	}
	m.data, m.addr = data, addr
	if data&0x20 != 0 {
		m.mem.setVramMirror(memVramMirrorH)
	} else {
		m.mem.setVramMirror(memVramMirrorV)
	}
	d012, d34, a2 := uint32(data&0x07), uint32(data&0x18), uint32((addr&0x04)<<3)
	switch data >> 6 {
	case 0x00:
		m.mem.setProm16kBank(4, d012|d34|a2)
		m.mem.setProm16kBank(6, 0x07|d34|a2)
	case 0x01:
		m.mem.setProm16kBank(4, (d012&0x06)|d34|a2)
		m.mem.setProm16kBank(6, 0x07|d34|a2)
	case 0x02:
		m.mem.setProm16kBank(4, d012|d34|a2)
		m.mem.setProm16kBank(6, d012|d34|a2)
	case 0x03:
		m.mem.setProm32kBank(uint32(data&0x06) | (d34 >> 1) | (a2 >> 1))
	}
}
