package core

const (
	apuEventQueueLen     = 8192
	apuEventQueueLenMask = apuEventQueueLen - 1
)

type ApuDataQueue struct {
	rp, wp uint16
	data   [65536]float32
}

func (q *ApuDataQueue) reset() {
	q.rp, q.wp = 0, 0
}

func (q *ApuDataQueue) enqueue(d float32) {
	q.data[q.wp] = d
	q.wp++
}

func (q *ApuDataQueue) Dequeue(buf []float32) {
	b, e := q.rp, q.rp+uint16(len(buf))
	if e >= b {
		copy(buf, q.data[b:e])
	} else {
		copy(buf, q.data[b:])
		copy(buf[65536-uint32(b):], q.data[:e])
	}
	q.rp = e
}

type apuEventQueueNode struct {
	data byte
	addr uint16
	time int64
}

type apuEventQueue struct {
	rp, wp uint16
	data   [apuEventQueueLen]apuEventQueueNode
}

func (q *apuEventQueue) reset() {
	q.rp, q.wp = 0, 0
}

func (q *apuEventQueue) enqueue(n *apuEventQueueNode) {
	q.data[q.wp] = *n
	q.wp = (q.wp + 1) & apuEventQueueLenMask
}

func (q *apuEventQueue) dequeue(time int64, no *apuEventQueueNode) bool {
	n := &(q.data[q.rp])
	if q.wp == q.rp || n.time > time {
		return false
	}
	*no = *n
	q.rp = (q.rp + 1) & apuEventQueueLenMask
	return true
}

func (q *apuEventQueue) forceDequeue(no *apuEventQueueNode) bool {
	if q.wp == q.rp {
		return false
	}
	*no = q.data[q.rp]
	q.rp = (q.rp + 1) & apuEventQueueLenMask
	return true
}

type Apu struct {
	sys *Sys

	sampRate    uint16
	renderLen   uint16
	ratio       int32
	frameNCycle int64
	rate        float64
	cutoff      float64

	reg     byte
	syncReg byte
	time    float64
	outTmp  float64

	frameIrqOccur bool
	frameIrq      byte
	frameCnt      uint32
	frameCycle    int32

	ch0 *apuChanRect
	ch1 *apuChanRect
	ch2 *apuChanTri
	ch3 *apuChanNoise
	ch4 *apuChanDpcm
	dq  ApuDataQueue
	eq  apuEventQueue
}

func newApu(sys *Sys) *Apu {
	apu := &Apu{}
	apu.sys = sys

	tvFormat, conf := &sys.tvFormat, sys.conf
	apu.sampRate = conf.AudioSampRate
	apu.renderLen = apu.sampRate / 60
	apu.ratio = int32(tvFormat.cpuRate * 65536.0 / float32(apu.sampRate))
	apu.frameNCycle = tvFormat.nScanlineCycle * int64(tvFormat.nScanline)
	apu.rate = float64(apu.frameNCycle) * 5 / float64(apu.sampRate)
	apu.cutoff = 251.32741228718345 / float64(apu.sampRate)

	apu.ch0 = &apuChanRect{apu: apu, enMask: 0x01}
	apu.ch1 = &apuChanRect{apu: apu, enMask: 0x02}
	apu.ch2 = &apuChanTri{apu: apu}
	apu.ch3 = &apuChanNoise{apu: apu}
	apu.ch4 = &apuChanDpcm{apu: apu}

	return apu
}

func (apu *Apu) reset(init bool) {
	if !init {
		apu.ch0 = &apuChanRect{apu: apu, enMask: 0x01}
		apu.ch1 = &apuChanRect{apu: apu, enMask: 0x02}
		apu.ch2 = &apuChanTri{apu: apu}
		apu.ch3 = &apuChanNoise{apu: apu}
	}

	apu.reg, apu.syncReg = 0, 0
	apu.time = 0
	apu.dq.reset()
	apu.eq.reset()
	apu.ch3.shiftReg = 0x4000
	for i := uint16(0x4000); i <= 0x4010; i++ {
		apu.writeAsync(i, 0)
		apu.write(i, 0)
	}
	for i := uint16(0x4012); i <= 0x4015; i++ {
		apu.writeAsync(i, 0)
		apu.write(i, 0)
	}
	apu.frameIrqOccur, apu.frameIrq, apu.frameCnt, apu.frameCycle = false, 0xc0, 0, 0
}

func (apu *Apu) read(addr uint16) byte {
	var data byte
	switch addr {
	case 0x4015:
		if apu.ch0.syncEn && apu.ch0.syncLenCount != 0 {
			data |= 0x01
		}
		if apu.ch1.syncEn && apu.ch1.syncLenCount != 0 {
			data |= 0x02
		}
		if apu.ch2.syncEn && apu.ch2.syncLenCount != 0 {
			data |= 0x04
		}
		if apu.ch3.syncEn && apu.ch3.syncLenCount != 0 {
			data |= 0x08
		}
		if apu.ch4.syncEn && apu.ch4.syncDmaLen != 0 {
			data |= 0x10
		}
		if apu.frameIrqOccur {
			data |= 0x40
		}
		if apu.ch4.syncIrqEn {
			data |= 0x80
		}
		apu.frameIrqOccur = false
		apu.sys.cpu.intr &^= cpuIntrTypFrame
	case 0x4017:
		if !apu.frameIrqOccur {
			data = byte(addr>>8) | 0x40
		}
	default:
		data = byte(addr >> 8)
	}
	return data
}

func (apu *Apu) writeRegAsync(data byte) {
	apu.reg = data
	ch0, ch1, ch2, ch3, ch4 := apu.ch0, apu.ch1, apu.ch2, apu.ch3, apu.ch4
	if data&0x01 == 0 {
		ch0.en, ch0.lenCount = false, 0
	}
	if data&0x02 == 0 {
		ch1.en, ch1.lenCount = false, 0
	}
	if data&0x04 == 0 {
		ch2.en, ch2.lenCount, ch2.linCount, ch2.cntStart = false, 0, 0, false
	}
	if data&0x08 == 0 {
		ch3.en, ch3.lenCount = false, 0
	}
	if data&0x10 == 0 {
		ch4.en, ch4.dmaLen = false, 0
	} else {
		ch4.en = true
		if ch4.dmaLen == 0 {
			ch4.addr, ch4.dmaLen, ch4.phaseAcc = ch4.addrCache, ch4.dmaLenCache, 0
		}
	}
}

func (apu *Apu) writeAsync(addr uint16, data byte) {
	switch addr {
	case 0x4000, 0x4001, 0x4002, 0x4003:
		apu.ch0.writeAsync(addr, data)
	case 0x4004, 0x4005, 0x4006, 0x4007:
		apu.ch1.writeAsync(addr, data)
	case 0x4008, 0x4009, 0x400a, 0x400b:
		apu.ch2.writeAsync(addr, data)
	case 0x400c, 0x400d, 0x400e, 0x400f:
		apu.ch3.writeAsync(addr, data)
	case 0x4010, 0x4011, 0x4012, 0x4013:
		apu.ch4.writeAsync(addr, data)
	case 0x4015:
		apu.writeRegAsync(data)
	case 0x4018:
		apu.ch0.updateAsync(data)
		apu.ch1.updateAsync(data)
		apu.ch2.updateAsync(data)
		apu.ch3.updateAsync(data)
	}
}

func (apu *Apu) writeReg(data byte) {
	apu.syncReg = data
	ch0, ch1, ch2, ch3, ch4 := apu.ch0, apu.ch1, apu.ch2, apu.ch3, apu.ch4
	if data&0x01 == 0 {
		ch0.syncEn, ch0.syncLenCount = false, 0
	}
	if data&0x02 == 0 {
		ch1.syncEn, ch1.syncLenCount = false, 0
	}
	if data&0x04 == 0 {
		ch2.syncEn, ch2.syncLenCount, ch2.syncLinCount, ch2.syncCntStart = false, 0, 0, false
	}
	if data&0x08 == 0 {
		ch3.syncEn, ch3.syncLenCount = false, 0
	}
	if data&0x10 == 0 {
		ch4.syncEn, ch4.syncIrqEn, ch4.syncDmaLen = false, false, 0
		apu.sys.cpu.intr &^= cpuIntrTypDpcm
	} else {
		ch4.syncEn = true
		if ch4.syncDmaLen == 0 {
			ch4.syncDmaLen, ch4.syncNCycle = ch4.syncDmaLenCache, 0
		}
	}
}

func (apu *Apu) updateFrame() {
	switch apu.frameCnt {
	case 0:
		if apu.frameIrq&0xc0 == 0 {
			apu.frameIrqOccur = true
			apu.sys.cpu.intr |= cpuIntrTypFrame
		}
	case 3:
		if apu.frameIrq&0x80 != 0 {
			apu.frameCycle += 14915
		}
	}
	apu.write(0x4018, byte(apu.frameCnt))
	apu.frameCnt = (apu.frameCnt + 1) & 3
}

func (apu *Apu) write(addr uint16, data byte) {
	switch addr {
	case 0x4000, 0x4001, 0x4002, 0x4003:
		apu.ch0.write(addr, data)
	case 0x4004, 0x4005, 0x4006, 0x4007:
		apu.ch1.write(addr, data)
	case 0x4008, 0x4009, 0x400a, 0x400b:
		apu.ch2.write(addr, data)
	case 0x400c, 0x400d, 0x400e, 0x400f:
		apu.ch3.write(addr, data)
	case 0x4010, 0x4011, 0x4012, 0x4013:
		apu.ch4.write(addr, data)
	case 0x4015:
		apu.writeReg(data)
	case 0x4017:
		apu.frameCycle, apu.frameIrq, apu.frameIrqOccur = 0, data, false
		apu.sys.cpu.intr &^= cpuIntrTypFrame
		apu.frameCnt = 0
		if data&0x80 != 0 {
			apu.updateFrame()
		}
		apu.frameCnt, apu.frameCycle = 1, 14915
	case 0x4018:
		apu.ch0.update(data)
		apu.ch1.update(data)
		apu.ch2.update(data)
		apu.ch3.update(data)
	}
	n := &apuEventQueueNode{data, addr, apu.sys.cpu.nCycle}
	apu.eq.enqueue(n)
}

func (apu *Apu) sync(nCycle int32) {
	apu.frameCycle -= nCycle << 1
	if apu.frameCycle <= 0 {
		apu.frameCycle += 14915
		apu.updateFrame()
	}
	apu.ch4.update(nCycle)
}

func (apu *Apu) render() {
	cpuNCycle := apu.sys.cpu.nCycle
	if int64(apu.time) > cpuNCycle {
		n := apuEventQueueNode{}
		for apu.eq.forceDequeue(&n) {
			apu.writeAsync(n.addr, n.data)
		}
	}

	for i := uint16(0); i < apu.renderLen; i++ {
		t := int64(apu.time)
		n := apuEventQueueNode{}
		for apu.eq.dequeue(t, &n) {
			apu.writeAsync(n.addr, n.data)
		}

		o := (apu.ch0.render()*0x00f0 + apu.ch1.render()*0x00f0 + apu.ch2.render()*0x0130 +
			apu.ch3.render()*0x00c0 + apu.ch4.render()*0x00f0) >> 8
		o1 := float64(o) - apu.outTmp
		apu.outTmp += apu.cutoff * o1
		o1 /= 32768
		apu.dq.enqueue(float32(o1))
		apu.time += apu.rate
	}
	if d := int64(apu.time) - cpuNCycle; d > apu.frameNCycle/24 || d < -apu.frameNCycle/6 {
		apu.time = float64(cpuNCycle)
	}
}
