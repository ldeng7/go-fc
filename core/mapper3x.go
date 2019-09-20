package core

// 0x30

type mapper30 struct {
	baseMapper
}

func newMapper30(bm *baseMapper) Mapper {
	return &mapper30{baseMapper: *bm}
}

func (m *mapper30) reset() {
}

// 0x31

type mapper31 struct {
	baseMapper
}

func newMapper31(bm *baseMapper) Mapper {
	return &mapper31{baseMapper: *bm}
}

func (m *mapper31) reset() {
}

// 0x32

type mapper32 struct {
	baseMapper
}

func newMapper32(bm *baseMapper) Mapper {
	return &mapper32{baseMapper: *bm}
}

func (m *mapper32) reset() {
}

// 0x33

type mapper33 struct {
	baseMapper
}

func newMapper33(bm *baseMapper) Mapper {
	return &mapper33{baseMapper: *bm}
}

func (m *mapper33) reset() {
}

// 0x34

type mapper34 struct {
	baseMapper
}

func newMapper34(bm *baseMapper) Mapper {
	return &mapper34{baseMapper: *bm}
}

func (m *mapper34) reset() {
}

// 0x35

type mapper35 struct {
	baseMapper
}

func newMapper35(bm *baseMapper) Mapper {
	return &mapper35{baseMapper: *bm}
}

func (m *mapper35) reset() {
}

// 0x36

type mapper36 struct {
	baseMapper
}

func newMapper36(bm *baseMapper) Mapper {
	return &mapper36{baseMapper: *bm}
}

func (m *mapper36) reset() {
}

// 0x37

type mapper37 struct {
	baseMapper
}

func newMapper37(bm *baseMapper) Mapper {
	return &mapper37{baseMapper: *bm}
}

func (m *mapper37) reset() {
}

// 0x38

type mapper38 struct {
	baseMapper
}

func newMapper38(bm *baseMapper) Mapper {
	return &mapper38{baseMapper: *bm}
}

func (m *mapper38) reset() {
}

// 0x39

type mapper39 struct {
	baseMapper
}

func newMapper39(bm *baseMapper) Mapper {
	return &mapper39{baseMapper: *bm}
}

func (m *mapper39) reset() {
}

// 0x3a

type mapper3a struct {
	baseMapper
}

func newMapper3a(bm *baseMapper) Mapper {
	return &mapper3a{baseMapper: *bm}
}

func (m *mapper3a) reset() {
	m.mem.setProm32kBank4(0, 1, 0, 1)
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper3a) write(addr uint16, data byte) {
	if addr&0x40 != 0 {
		m.mem.setProm16kBank(4, uint32(addr)&0x07)
		m.mem.setProm16kBank(6, uint32(addr)&0x07)
	} else {
		m.mem.setProm32kBank(uint32(addr&0x06) >> 1)
	}
	if m.mem.nVrom1kPage != 0 {
		m.mem.setVrom8kBank(uint32(addr&0x38) >> 3)
	}
	if data&0x02 != 0 {
		m.mem.setVramMirror(memVramMirrorV)
	} else {
		m.mem.setVramMirror(memVramMirrorH)
	}
}

// 0x3b

type mapper3b struct {
	baseMapper
}

func newMapper3b(bm *baseMapper) Mapper {
	return &mapper3b{baseMapper: *bm}
}

func (m *mapper3b) reset() {
}

// 0x3c

type mapper3c struct {
	baseMapper
	idx   byte
	patch bool
}

func newMapper3c(bm *baseMapper) Mapper {
	return &mapper3c{baseMapper: *bm}
}

func (m *mapper3c) reset() {
	switch m.sys.conf.PatchTyp {
	case 1:
		m.patch = true
		i := uint32(m.idx)
		m.mem.setProm16kBank(4, i)
		m.mem.setProm16kBank(6, i)
		m.mem.setVrom8kBank(i)
		m.idx = (m.idx + 1) & 0x03
	default:
		m.mem.setProm32kBank(0)
		m.mem.setVrom8kBank(0)
	}
}

func (m *mapper3c) write(addr uint16, data byte) {
	if m.patch {
		return
	}
	if addr&0x80 != 0 {
		m.mem.setProm16kBank(4, uint32(addr&0x70)>>4)
		m.mem.setProm16kBank(6, uint32(addr&0x70)>>4)
	} else {
		m.mem.setProm32kBank(uint32(addr&0x70) >> 5)
	}
	m.mem.setVrom8kBank(uint32(addr & 0x07))
	if data&0x08 != 0 {
		m.mem.setVramMirror(memVramMirrorV)
	} else {
		m.mem.setVramMirror(memVramMirrorH)
	}
}

// 0x3d

type mapper3d struct {
	baseMapper
}

func newMapper3d(bm *baseMapper) Mapper {
	return &mapper3d{baseMapper: *bm}
}

func (m *mapper3d) reset() {
	m.mem.setProm32kBank4(0, 1, m.mem.nProm8kPage-2, m.mem.nProm8kPage-1)
}

func (m *mapper3d) write(addr uint16, data byte) {
	b := uint32(data)
	switch addr & 0x30 {
	case 0x00, 0x30:
		m.mem.setProm32kBank(b & 0x0f)
	case 0x10, 0x20:
		m.mem.setProm16kBank(4, ((b&0x0f)<<1)|((b&0x20)>>4))
		m.mem.setProm16kBank(6, ((b&0x0f)<<1)|((b&0x20)>>4))
	}

	if addr&0x80 != 0 {
		m.mem.setVramMirror(memVramMirrorH)
	} else {
		m.mem.setVramMirror(memVramMirrorV)
	}
}

// 0x3e

type mapper3e struct {
	baseMapper
}

func newMapper3e(bm *baseMapper) Mapper {
	return &mapper3e{baseMapper: *bm}
}

func (m *mapper3e) reset() {
	m.mem.setProm32kBank(0)
	m.mem.setVrom8kBank(0)
}

func (m *mapper3e) write(addr uint16, data byte) {
	b := uint32(data)
	switch addr & 0xff00 {
	case 0x8100:
		m.mem.setProm8kBank(4, b)
		m.mem.setProm8kBank(5, b+1)
	case 0x8500:
		m.mem.setProm8kBank(4, b)
	case 0x8700:
		m.mem.setProm8kBank(5, b)
	default:
		for i := byte(0); i < 8; i++ {
			m.mem.setVrom1kBank(i, b+uint32(i))
		}
	}
}

// 0x3f

type mapper3f struct {
	baseMapper
}

func newMapper3f(bm *baseMapper) Mapper {
	return &mapper3f{baseMapper: *bm}
}

func (m *mapper3f) reset() {
}
