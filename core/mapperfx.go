package core

// 0xf0

type mapperf0 struct {
	baseMapper
}

func newMapperf0(bm *baseMapper) Mapper {
	return &mapperf0{baseMapper: *bm}
}

func (m *mapperf0) init() {
}

// 0xf1

type mapperf1 struct {
	baseMapper
}

func newMapperf1(bm *baseMapper) Mapper {
	return &mapperf1{baseMapper: *bm}
}

func (m *mapperf1) init() {
}

// 0xf2

type mapperf2 struct {
	baseMapper
}

func newMapperf2(bm *baseMapper) Mapper {
	return &mapperf2{baseMapper: *bm}
}

func (m *mapperf2) init() {
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
}

func newMapperf3(bm *baseMapper) Mapper {
	return &mapperf3{baseMapper: *bm}
}

func (m *mapperf3) init() {
}

// 0xf4

type mapperf4 struct {
	baseMapper
}

func newMapperf4(bm *baseMapper) Mapper {
	return &mapperf4{baseMapper: *bm}
}

func (m *mapperf4) init() {
}

// 0xf5

type mapperf5 struct {
	baseMapper
}

func newMapperf5(bm *baseMapper) Mapper {
	return &mapperf5{baseMapper: *bm}
}

func (m *mapperf5) init() {
}

// 0xf6

type mapperf6 struct {
	baseMapper
}

func newMapperf6(bm *baseMapper) Mapper {
	return &mapperf6{baseMapper: *bm}
}

func (m *mapperf6) init() {
}

// 0xf7

type mapperf7 struct {
	baseMapper
}

func newMapperf7(bm *baseMapper) Mapper {
	return &mapperf7{baseMapper: *bm}
}

func (m *mapperf7) init() {
}

// 0xf8

type mapperf8 struct {
	baseMapper
}

func newMapperf8(bm *baseMapper) Mapper {
	return &mapperf8{baseMapper: *bm}
}

func (m *mapperf8) init() {
}

// 0xf9

type mapperf9 struct {
	baseMapper
}

func newMapperf9(bm *baseMapper) Mapper {
	return &mapperf9{baseMapper: *bm}
}

func (m *mapperf9) init() {
}

// 0xfa

type mapperfa struct {
	baseMapper
}

func newMapperfa(bm *baseMapper) Mapper {
	return &mapperfa{baseMapper: *bm}
}

func (m *mapperfa) init() {
}

// 0xfb

type mapperfb struct {
	baseMapper
}

func newMapperfb(bm *baseMapper) Mapper {
	return &mapperfb{baseMapper: *bm}
}

func (m *mapperfb) init() {
}

// 0xfc

type mapperfc struct {
	baseMapper
}

func newMapperfc(bm *baseMapper) Mapper {
	return &mapperfc{baseMapper: *bm}
}

func (m *mapperfc) init() {
}

// 0xfd

type mapperfd struct {
	baseMapper
}

func newMapperfd(bm *baseMapper) Mapper {
	return &mapperfd{baseMapper: *bm}
}

func (m *mapperfd) init() {
}

// 0xfe

type mapperfe struct {
	baseMapper
}

func newMapperfe(bm *baseMapper) Mapper {
	return &mapperfe{baseMapper: *bm}
}

func (m *mapperfe) init() {
}

// 0xff

type mapperff struct {
	baseMapper
}

func newMapperff(bm *baseMapper) Mapper {
	return &mapperff{baseMapper: *bm}
}

func (m *mapperff) init() {
}
