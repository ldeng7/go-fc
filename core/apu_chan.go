package core

var apuVblLen = [32]byte{
	5, 127, 10, 1, 19, 2, 40, 3, 80, 4, 30, 5, 7, 6, 13, 7,
	6, 8, 12, 9, 24, 10, 48, 11, 96, 12, 36, 13, 8, 14, 16, 15,
}
var apuFreqLimit = [8]int32{
	0x03ff, 0x0555, 0x0666, 0x071c, 0x0787, 0x07c1, 0x07e0, 0x07f0,
}
var apuDutyLut = [4]byte{
	2, 4, 8, 12,
}
var apuNoiseFreq = [16]int32{
	4, 8, 16, 32, 64, 96, 128, 160, 202, 254, 380, 508, 762, 1016, 2034, 4068,
}
var apuDpcmCyclesPal = [16]uint16{
	397, 353, 315, 297, 265, 235, 209, 198, 176, 148, 131, 118, 98, 78, 66, 50,
}
var apuDpcmCyclesNtsc = [16]uint16{
	428, 380, 340, 320, 286, 254, 226, 214, 190, 160, 142, 128, 106, 85, 72, 54,
}

type apuChanRect struct {
	apu    *Apu
	enMask byte

	en       bool
	holdnote bool
	volume   byte
	reg      [4]byte

	adder     byte
	duty      byte
	lenCount  byte
	freq      int32
	freqLimit int32
	curVolume int32
	phaseAcc  int32

	envFixed bool
	envDecay byte
	envCount byte
	envVol   byte

	swpOn    bool
	swpInc   bool
	swpShift byte
	swpDecay byte
	swpCount byte

	syncEn       bool
	syncHoldnote bool
	syncLenCount byte
	syncReg      [4]byte
}

func (ch *apuChanRect) writeAsync(addr uint16, data byte) {
	i := addr & 0x03
	ch.reg[i] = data
	switch i {
	case 0x00:
		ch.holdnote, ch.envFixed = data&0x20 != 0, data&0x10 != 0
		ch.volume, ch.envDecay, ch.duty = data&0x0f, (data&0x0f)+1, apuDutyLut[data>>6]
	case 0x01:
		ch.swpOn, ch.swpInc = data&0x80 != 0, data&0x08 != 0
		ch.swpShift, ch.swpDecay = data&0x07, ((data>>4)&0x07)+1
		ch.freqLimit = apuFreqLimit[data&0x07]
	case 0x02:
		ch.freq = (ch.freq & 0xff00) | int32(data)
	case 0x03:
		ch.freq = (int32(data&0x07) << 8) | (ch.freq & 0x00ff)
		ch.lenCount = apuVblLen[data>>3] << 1
		ch.envVol, ch.envCount, ch.adder = 0x0f, ch.envDecay+1, 0
		if ch.apu.reg&ch.enMask != 0 {
			ch.en = true
		}
	}
}

func (ch *apuChanRect) write(addr uint16, data byte) {
	i := addr & 0x03
	ch.syncReg[i] = data
	switch i {
	case 0x00:
		ch.syncHoldnote = data&0x20 != 0
	case 0x03:
		ch.syncLenCount = apuVblLen[data>>3] << 1
		if ch.apu.syncReg&ch.enMask != 0 {
			ch.syncEn = true
		}
	}
}

func (ch *apuChanRect) updateAsync(typ byte) {
	if !ch.en || ch.lenCount == 0 {
		return
	}
	if typ&0x01 == 0 {
		if !ch.holdnote {
			ch.lenCount--
		}
		if ch.swpOn && ch.swpShift != 0 {
			if ch.swpCount != 0 {
				ch.swpCount--
			}
			if ch.swpCount == 0 {
				ch.swpCount = ch.swpDecay
				if ch.swpInc {
					if ch.enMask == 0x01 {
						ch.freq += ^(ch.freq >> ch.swpShift)
					} else {
						ch.freq -= (ch.freq >> ch.swpShift)
					}
				} else {
					ch.freq += (ch.freq >> ch.swpShift)
				}
			}
		}
	}
	if ch.envCount != 0 {
		ch.envCount--
	}
	if ch.envCount == 0 {
		ch.envCount = ch.envDecay
		if ch.holdnote {
			ch.envVol = (ch.envVol - 1) & 0x0f
		} else if ch.envVol != 0 {
			ch.envVol--
		}
	}
	if !ch.envFixed {
		ch.curVolume = int32(ch.envVol) << 8
	}
}

func (ch *apuChanRect) update(typ byte) {
	if !ch.syncEn || ch.syncLenCount == 0 {
		return
	}
	if typ&0x01 == 0 && !ch.syncHoldnote {
		ch.syncLenCount--
	}
}

func (ch *apuChanRect) render() int32 {
	if !ch.en || ch.lenCount == 0 {
		return 0
	}
	if ch.freq < 8 || (!ch.swpInc && ch.freq > ch.freqLimit) {
		return 0
	}
	if ch.envFixed {
		ch.curVolume = int32(ch.volume) << 8
	}
	r, w, f := ch.apu.ratio, ch.phaseAcc, (ch.freq+1)<<16
	if w > r {
		w = r
	}
	s := int64(w)
	if ch.adder >= ch.duty {
		s = int64(-w)
	}
	ch.phaseAcc -= r
	for ch.phaseAcc < 0 {
		ch.phaseAcc += f
		ch.adder, w = (ch.adder+1)&0x0f, f
		if ch.phaseAcc > 0 {
			w -= ch.phaseAcc
		}
		if ch.adder < ch.duty {
			s += int64(w)
		} else {
			s -= int64(w)
		}
	}
	return int32(int64(ch.curVolume) * s / int64(r))
}

type apuChanTri struct {
	apu *Apu

	en       bool
	holdnote bool
	cntStart bool
	reg      [4]byte

	adder     byte
	lenCount  byte
	linCount  byte
	freq      int32
	curVolume int32
	phaseAcc  int32

	syncEn       bool
	syncHoldnote bool
	syncCntStart bool
	syncLenCount byte
	syncLinCount byte
	syncReg      [4]byte
}

func (ch *apuChanTri) writeAsync(addr uint16, data byte) {
	i := addr & 0x03
	ch.reg[i] = data
	switch i {
	case 0x00:
		ch.holdnote = data&0x80 != 0
	case 0x02:
		ch.freq = ((int32(ch.reg[3]&0x07) << 8) + int32(data) + 1) << 16
	case 0x03:
		ch.freq = ((int32(data&0x07) << 8) + int32(ch.reg[2]) + 1) << 16
		ch.lenCount = apuVblLen[data>>3] << 1
		ch.cntStart = true
		if ch.apu.reg&0x04 != 0 {
			ch.en = true
		}
	}
}

func (ch *apuChanTri) write(addr uint16, data byte) {
	i := addr & 0x03
	ch.syncReg[i] = data
	switch i {
	case 0x00:
		ch.syncHoldnote = data&0x80 != 0
	case 0x03:
		ch.syncLenCount = apuVblLen[data>>3] << 1
		ch.syncCntStart = true
		if ch.apu.syncReg&0x04 != 0 {
			ch.syncEn = true
		}
	}
}

func (ch *apuChanTri) updateAsync(typ byte) {
	if !ch.en {
		return
	}
	if typ&0x01 == 0 && !ch.holdnote && ch.lenCount != 0 {
		ch.lenCount--
	}
	if ch.cntStart {
		ch.linCount = ch.reg[0] & 0x7f
	} else if ch.linCount != 0 {
		ch.linCount--
	}
	if !ch.holdnote && ch.linCount != 0 {
		ch.cntStart = false
	}
}

func (ch *apuChanTri) update(typ byte) {
	if !ch.syncEn {
		return
	}
	if typ&0x01 == 0 && !ch.syncHoldnote && ch.syncLenCount != 0 {
		ch.syncLenCount--
	}
	if ch.syncCntStart {
		ch.syncLinCount = ch.syncReg[0] & 0x7f
	} else if ch.syncLinCount != 0 {
		ch.syncLinCount--
	}
	if !ch.syncHoldnote && ch.syncLinCount != 0 {
		ch.syncCntStart = false
	}
}

func (ch *apuChanTri) renderSub() {
	ch.phaseAcc += ch.freq
	ch.adder = (ch.adder + 1) & 0x1f
	if ch.adder < 0x10 {
		ch.curVolume = int32(ch.adder&0x0f) << 9
	} else {
		ch.curVolume = int32(0x0f-(ch.adder&0x0f)) << 9
	}
}

func (ch *apuChanTri) render() int32 {
	chd := ch.apu.ch4
	var vol int32 = 256 - (int32(chd.dpcmValue) << 1) - (int32(chd.reg[1]) & 0x01)
	if !ch.en || ch.lenCount == 0 || ch.linCount == 0 || ch.freq < 0x00080000 {
		return ch.curVolume * vol / 256
	}
	ch.phaseAcc -= ch.apu.ratio
	if ch.phaseAcc >= 0 {
		return ch.curVolume * vol / 256
	}
	if ch.freq > ch.apu.ratio {
		ch.renderSub()
		return ch.curVolume * vol / 256
	}
	var c, s int32
	for ch.phaseAcc < 0 {
		ch.renderSub()
		s += ch.curVolume
		c++
	}
	return (s / c) * vol / 256
}

type apuChanNoise struct {
	apu *Apu

	en       bool
	holdnote bool
	volume   byte
	xorTap   byte
	shiftReg uint16
	reg      [4]byte

	lenCount  byte
	freq      int32
	curVolume int32
	output    int32
	phaseAcc  int32

	envFixed bool
	envDecay byte
	envCount byte
	envVol   byte

	syncEn       bool
	syncHoldnote bool
	syncLenCount byte
	syncReg      [4]byte
}

func (ch *apuChanNoise) writeAsync(addr uint16, data byte) {
	i := addr & 0x03
	ch.reg[i] = data
	switch i {
	case 0x00:
		ch.holdnote, ch.envFixed = data&0x20 != 0, data&0x10 != 0
		ch.volume, ch.envDecay = data&0x0f, (data&0x0f)+1
	case 0x02:
		ch.freq = apuNoiseFreq[data&0x0f] << 16
		ch.xorTap = 0x40
		if data&0x80 == 0 {
			ch.xorTap = 0x02
		}
	case 0x03:
		ch.lenCount = apuVblLen[data>>3] << 1
		ch.envVol, ch.envCount = 0x0f, ch.envDecay+1
		if ch.apu.reg&0x08 != 0 {
			ch.en = true
		}
	}
}

func (ch *apuChanNoise) write(addr uint16, data byte) {
	i := addr & 0x03
	ch.reg[i] = data
	switch i {
	case 0x00:
		ch.syncHoldnote = data&0x20 != 0
	case 0x03:
		ch.syncLenCount = apuVblLen[data>>3] << 1
		if ch.apu.syncReg&0x08 != 0 {
			ch.syncEn = true
		}
	}
}

func (ch *apuChanNoise) updateAsync(typ byte) {
	if !ch.en || ch.lenCount == 0 {
		return
	}
	if typ&0x01 == 0 && !ch.holdnote {
		ch.lenCount--
	}
	if ch.envCount != 0 {
		ch.envCount--
	}
	if ch.envCount == 0 {
		ch.envCount = ch.envDecay
		if ch.holdnote {
			ch.envVol = (ch.envVol - 1) & 0x0f
		} else if ch.envVol != 0 {
			ch.envVol--
		}
	}
	if !ch.envFixed {
		ch.curVolume = int32(ch.envVol) << 8
	}
}

func (ch *apuChanNoise) update(typ byte) {
	if !ch.syncEn || ch.syncLenCount == 0 {
		return
	}
	if typ&0x01 == 0 && !ch.syncHoldnote {
		ch.syncLenCount--
	}
}

func (ch *apuChanNoise) renderSub() {
	ch.phaseAcc += ch.freq
	b := ch.shiftReg & 0x01
	b1 := b
	if ch.shiftReg&uint16(ch.xorTap) != 0 {
		b1 ^= 0x01
	}
	ch.shiftReg >>= 1
	ch.shiftReg |= b1 << 14
	if b == 0 {
		ch.output = ch.curVolume
	} else {
		ch.output -= ch.curVolume
	}
}

func (ch *apuChanNoise) render() int32 {
	if !ch.en || ch.lenCount == 0 {
		return 0
	}
	if ch.envFixed {
		ch.curVolume = int32(ch.volume) << 8
	}
	chd := ch.apu.ch4
	var vol int32 = 256 - (int32(chd.dpcmValue) << 1) - (int32(chd.reg[1]) & 0x01)
	ch.phaseAcc -= ch.apu.ratio
	if ch.phaseAcc >= 0 {
		return ch.output * vol / 256
	}
	if ch.freq > ch.apu.ratio {
		ch.renderSub()
		return ch.output * vol / 256
	}
	var c, s int32
	for ch.phaseAcc < 0 {
		ch.renderSub()
		s += ch.output
		c++
	}
	return (s / c) * vol / 256
}

type apuChanDpcm struct {
	apu      *Apu
	cycleTbl *[16]uint16

	en        bool
	looping   bool
	curByte   byte
	dpcmValue byte
	reg       [4]byte

	addr        uint16
	addrCache   uint16
	dmaLen      uint16
	dmaLenCache uint16
	freq        int32
	outputReal  int32
	outputFake  int32
	output      int32
	phaseAcc    int32

	syncEn          bool
	syncLooping     bool
	syncIrqGen      bool
	syncIrqEn       bool
	syncDmaLenCache uint16
	syncNCycleCache uint16
	syncDmaLen      uint16
	syncNCycle      int32
}

func (ch *apuChanDpcm) writeAsync(addr uint16, data byte) {
	i := addr & 0x03
	ch.reg[i] = data
	switch i {
	case 0x00:
		ch.freq = int32((*ch.cycleTbl)[data&0x0f]) << 16
		ch.looping = data&0x40 != 0
	case 0x01:
		ch.dpcmValue = (data & 0x7f) >> 1
	case 0x02:
		ch.addrCache = (uint16(data) << 6) | 0xc000
	case 0x03:
		ch.dmaLenCache = ((uint16(data) << 4) + 1) << 3
	}
}

func (ch *apuChanDpcm) write(addr uint16, data byte) {
	i := addr & 0x03
	ch.reg[i] = data
	switch i {
	case 0x00:
		ch.syncNCycleCache = (*ch.cycleTbl)[data&0x0f] << 3
		ch.syncLooping, ch.syncIrqGen = data&0x40 != 0, data&0x80 != 0
		if !ch.syncIrqGen {
			ch.syncIrqEn = false
			ch.apu.sys.cpu.intr &^= cpuIntrTypDpcm
		}
	case 0x03:
		ch.syncDmaLenCache = (uint16(data) << 4) + 1
	}
}

func (ch *apuChanDpcm) update(nCycle int32) {
	if !ch.syncEn {
		return
	}
	ch.syncNCycle -= nCycle
	for ch.syncNCycle < 0 {
		ch.syncNCycle += int32(ch.syncNCycleCache)
		if ch.syncDmaLen == 0 {
			continue
		}
		ch.syncDmaLen--
		if ch.syncDmaLen >= 2 {
			continue
		}
		if !ch.syncLooping {
			ch.syncDmaLen = 0
			if ch.syncIrqGen {
				ch.syncIrqEn = true
				ch.apu.sys.cpu.intr |= cpuIntrTypDpcm
			}
		} else {
			ch.syncDmaLen = ch.syncDmaLenCache
		}
	}
}

func (ch *apuChanDpcm) render() int32 {
	if ch.dmaLen != 0 {
		ch.phaseAcc -= ch.apu.ratio
		sys := ch.apu.sys
		for ch.phaseAcc < 0 {
			ch.phaseAcc += ch.freq
			if ch.dmaLen&0x07 == 0 {
				ch.curByte = sys.read(ch.addr)
				if ch.addr == 0xffff {
					ch.addr = 0x8000
				} else {
					ch.addr++
				}
			}
			ch.dmaLen--
			if ch.dmaLen == 0 {
				if ch.looping {
					ch.addr, ch.dmaLen = ch.addrCache, ch.dmaLenCache
				} else {
					ch.en = false
					break
				}
			}
			if ch.curByte&(1<<((ch.dmaLen&0x07)^0x07)) != 0 {
				if ch.dpcmValue < 0x3f {
					ch.dpcmValue++
				}
			} else {
				if ch.dpcmValue > 1 {
					ch.dpcmValue--
				}
			}
		}
	}
	ch.outputReal = (int32(ch.dpcmValue) << 1) + (int32(ch.reg[1]) & 0x01) - 64
	if d := ch.outputReal - ch.outputFake; d <= 8 && d >= -8 {
		ch.outputFake, ch.output = ch.outputReal, ch.outputReal<<8
	} else {
		if ch.outputReal > ch.outputFake {
			ch.outputFake += 8
		} else {
			ch.outputFake -= 8
		}
		ch.output = ch.outputFake << 8
	}
	return ch.output
}
