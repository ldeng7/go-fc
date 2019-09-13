package core

type Apu struct {
	sys *Sys
}

func newApu(sys *Sys) *Apu {
	apu := &Apu{}
	apu.sys = sys

	return apu
}

func (apu *Apu) read(addr uint16) byte {
	return 0
}

func (apu *Apu) write(addr uint16, data byte) {
}

func (aou *Apu) syncDpcm(nCycle int64) {

}
