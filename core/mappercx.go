package core

// 0xc0

type mapperc0 struct {
	baseMapper
}

func newMapperc0(bm *baseMapper) Mapper {
	return &mapperc0{baseMapper: *bm}
}

func (m *mapperc0) init() {
}

// 0xc1

type mapperc1 struct {
	baseMapper
}

func newMapperc1(bm *baseMapper) Mapper {
	return &mapperc1{baseMapper: *bm}
}

func (m *mapperc1) init() {
}

// 0xc2

type mapperc2 struct {
	baseMapper
}

func newMapperc2(bm *baseMapper) Mapper {
	return &mapperc2{baseMapper: *bm}
}

func (m *mapperc2) init() {
	m.mem.setProm32kBank(m.mem.prom32kSize - 1)
}

func (m *mapperc2) write(addr uint16, data byte) {
	m.mem.setProm8kBank(3, uint32(data))
}

// 0xc3

type mapperc3 struct {
	baseMapper
}

func newMapperc3(bm *baseMapper) Mapper {
	return &mapperc3{baseMapper: *bm}
}

func (m *mapperc3) init() {
}

// 0xc4

type mapperc4 struct {
	baseMapper
}

func newMapperc4(bm *baseMapper) Mapper {
	return &mapperc4{baseMapper: *bm}
}

func (m *mapperc4) init() {
}

// 0xc5

type mapperc5 struct {
	baseMapper
}

func newMapperc5(bm *baseMapper) Mapper {
	return &mapperc5{baseMapper: *bm}
}

func (m *mapperc5) init() {
}

// 0xc6

type mapperc6 struct {
	baseMapper
}

func newMapperc6(bm *baseMapper) Mapper {
	return &mapperc6{baseMapper: *bm}
}

func (m *mapperc6) init() {
}

// 0xc7

type mapperc7 struct {
	baseMapper
}

func newMapperc7(bm *baseMapper) Mapper {
	return &mapperc7{baseMapper: *bm}
}

func (m *mapperc7) init() {
}

// 0xc8

type mapperc8 struct {
	baseMapper
}

func newMapperc8(bm *baseMapper) Mapper {
	return &mapperc8{baseMapper: *bm}
}

func (m *mapperc8) init() {
}

// 0xc9

type mapperc9 struct {
	baseMapper
}

func newMapperc9(bm *baseMapper) Mapper {
	return &mapperc9{baseMapper: *bm}
}

func (m *mapperc9) init() {
}

// 0xca

type mapperca struct {
	baseMapper
}

func newMapperca(bm *baseMapper) Mapper {
	return &mapperca{baseMapper: *bm}
}

func (m *mapperca) init() {
}

// 0xcb

type mappercb struct {
	baseMapper
}

func newMappercb(bm *baseMapper) Mapper {
	return &mappercb{baseMapper: *bm}
}

func (m *mappercb) init() {
}

// 0xcc

type mappercc struct {
	baseMapper
}

func newMappercc(bm *baseMapper) Mapper {
	return &mappercc{baseMapper: *bm}
}

func (m *mappercc) init() {
}

// 0xcd

type mappercd struct {
	baseMapper
}

func newMappercd(bm *baseMapper) Mapper {
	return &mappercd{baseMapper: *bm}
}

func (m *mappercd) init() {
}

// 0xce

type mapperce struct {
	baseMapper
}

func newMapperce(bm *baseMapper) Mapper {
	return &mapperce{baseMapper: *bm}
}

func (m *mapperce) init() {
}

// 0xcf

type mappercf struct {
	baseMapper
}

func newMappercf(bm *baseMapper) Mapper {
	return &mappercf{baseMapper: *bm}
}

func (m *mappercf) init() {
}
