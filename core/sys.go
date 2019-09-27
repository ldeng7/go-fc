package core

import "io"

const (
	RenderModePre byte = iota
	RenderModePost
	RenderModePreAll
	RenderModePostAll
	RenderModeTile
)

type Conf struct {
	PatchTyp      uint64
	IsPal         bool
	AllSprite     bool
	RenderMode    byte
	AudioSampRate uint16
}

type tvFormat struct {
	cpuRate           float32
	nScanline         uint16
	nScanlineCycle    int64
	nHDrawCycle       int64
	nHBlankCycle      int64
	nScanlineEndCycle int64
	framePeriod       float32
}

var tvFormatPal = tvFormat{1662607.125, 312, 1278, 960, 318, 2, 1000.0 / 50.0}
var tvFormatNtsc = tvFormat{1789772.5, 262, 1364, 1024, 340, 4, 1000.0 / 60.0}

type Sys struct {
	rom    *Rom
	mem    *Mem
	mapper Mapper
	cpu    *Cpu
	ppu    *Ppu
	apu    *Apu
	pad    *Pad

	tvFormat   tvFormat
	renderMode byte
	conf       Conf

	scanline  uint16
	nCycle    int64
	nCycleReq int64
}

func NewSys(file io.Reader, conf *Conf) (*Sys, error) {
	var err error
	sys := &Sys{}

	sys.renderMode = conf.RenderMode
	if conf.IsPal {
		sys.tvFormat = tvFormatPal
	} else {
		sys.tvFormat = tvFormatNtsc
	}
	sys.conf = *conf

	if sys.rom, err = newRom(file); err != nil {
		return nil, err
	}
	sys.mem = newMem(sys)
	if sys.mapper, err = newMapper(sys); err != nil {
		return nil, err
	}
	sys.cpu = newCpu(sys)
	sys.ppu = newPpu(sys)
	sys.apu = newApu(sys)
	sys.pad = newPad()

	sys.reset(true)
	return sys, nil
}

func (sys *Sys) Reset() {
	sys.reset(false)
}

func (sys *Sys) GetFramePeriod() float32 {
	return sys.tvFormat.framePeriod
}

func (sys *Sys) SetFrameBuffer(fb *FrameBuffer) {
	sys.ppu.screen = fb
}

func (sys *Sys) SetPadKey(p byte, k byte, down bool) {
	sys.pad.setKey(p, k, down)
}

func (sys *Sys) GetAudioDataQueue() *ApuDataQueue {
	return &sys.apu.dq
}

func (sys *Sys) reset(init bool) {
	sys.mem.reset(init)
	sys.mapper.reset()
	sys.cpu.reset()
	sys.ppu.reset(init)
	sys.apu.reset(init)
	sys.pad.reset()
}

func (sys *Sys) read(addr uint16) byte {
	switch addr >> 13 {
	case 0x00:
		return sys.mem.ram[addr&0x07ff]
	case 0x01:
		return sys.ppu.read(addr & 0xe007)
	case 0x02:
		if addr >= 0x4100 {
			return sys.mapper.readLow(addr)
		}
		switch byte(addr) {
		case 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
			0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x15:
			return sys.apu.read(addr)
		case 0x14:
			return byte(addr)
		case 0x16:
			return sys.pad.read(addr) | 0x40
		case 0x17:
			return sys.pad.read(addr) | sys.apu.read(addr)
		default:
			return sys.mapper.readEx(addr)
		}
	case 0x03:
		return sys.mapper.readLow(addr)
	default:
		return sys.mem.cpuBanks[addr>>13][addr&0x1fff]
	}
}

func (sys *Sys) write(addr uint16, b byte) {
	switch addr >> 13 {
	case 0x00:
		sys.mem.ram[addr&0x07ff] = b
	case 0x01:
		sys.ppu.write(addr&0xe007, b)
	case 0x02:
		if addr >= 0x4100 {
			sys.mapper.writeLow(addr, b)
			break
		}
		switch byte(addr) {
		case 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
			0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x15:
			sys.apu.write(addr, b)
			sys.mem.cpuReg[byte(addr)] = b
		case 0x14:
			sys.ppu.dma(b)
			sys.cpu.nCycleDma += 514
			sys.mem.cpuReg[byte(addr)] = b
		case 0x16:
			sys.mapper.writeEx(addr, b)
			sys.pad.write(addr, b)
			sys.mem.cpuReg[byte(addr)] = b
		case 0x17:
			sys.pad.write(addr, b)
			sys.apu.write(addr, b)
			sys.mem.cpuReg[byte(addr)] = b
		case 0x18:
			sys.apu.write(addr, b)
		default:
			sys.mapper.writeEx(addr, b)
		}
	case 0x03:
		sys.mapper.writeLow(addr, b)
	default:
		sys.mapper.write(addr, b)
	}
}

func (sys *Sys) runCpu(nCycleReq int64) {
	sys.nCycleReq += nCycleReq
	if nCycleReqCpu := (sys.nCycleReq - sys.nCycle) / 12; nCycleReqCpu > 0 {
		exec := sys.cpu.run(nCycleReqCpu)
		sys.apu.sync(int32(exec))
		sys.nCycle += exec * 12
	}
}

func (sys *Sys) RunFrame() {
	ppu := sys.ppu
	bAllSprite := sys.conf.AllSprite
	nScanline := sys.tvFormat.nScanline - 1

	sys.scanline, ppu.iScanline = 0, 0
	switch sys.renderMode {
	case RenderModePostAll, RenderModePreAll:
		sys.runCpu(sys.tvFormat.nScanlineCycle)
		ppu.frameStart()
		ppu.scanlineNext()
		sys.mapper.hSync(sys.scanline)
		ppu.scanlineStart()
	case RenderModePost, RenderModePre:
		sys.runCpu(sys.tvFormat.nHDrawCycle)
		ppu.frameStart()
		ppu.scanlineNext()
		sys.mapper.hSync(sys.scanline)
		sys.runCpu(256)
		ppu.scanlineStart()
		sys.runCpu(80 + sys.tvFormat.nScanlineEndCycle)
	case RenderModeTile:
		sys.runCpu(1024)
		ppu.frameStart()
		ppu.scanlineNext()
		sys.runCpu(80)
		sys.mapper.hSync(sys.scanline)
		sys.runCpu(176)
		ppu.scanlineStart()
		sys.runCpu(80 + sys.tvFormat.nScanlineEndCycle)
	}

	for sys.scanline = 1; sys.scanline < 240; sys.scanline++ {
		ppu.iScanline = ScreenWidth * sys.scanline
		switch sys.renderMode {
		case RenderModePostAll:
			sys.runCpu(sys.tvFormat.nScanlineCycle)
			ppu.scanlineRender(byte(sys.scanline), bAllSprite)
			ppu.scanlineNext()
			sys.mapper.hSync(sys.scanline)
			ppu.scanlineStart()
		case RenderModePreAll:
			ppu.scanlineRender(byte(sys.scanline), bAllSprite)
			ppu.scanlineNext()
			sys.runCpu(sys.tvFormat.nScanlineCycle)
			sys.mapper.hSync(sys.scanline)
			ppu.scanlineStart()
		case RenderModePost:
			sys.runCpu(sys.tvFormat.nHDrawCycle)
			ppu.scanlineRender(byte(sys.scanline), bAllSprite)
			ppu.scanlineNext()
			sys.mapper.hSync(sys.scanline)
			sys.runCpu(256)
			ppu.scanlineStart()
			sys.runCpu(80 + sys.tvFormat.nScanlineEndCycle)
		case RenderModePre:
			ppu.scanlineRender(byte(sys.scanline), bAllSprite)
			sys.runCpu(sys.tvFormat.nHDrawCycle)
			ppu.scanlineNext()
			sys.mapper.hSync(sys.scanline)
			sys.runCpu(256)
			ppu.scanlineStart()
			sys.runCpu(80 + sys.tvFormat.nScanlineEndCycle)
		case RenderModeTile:
			ppu.scanlineRender(byte(sys.scanline), bAllSprite)
			ppu.scanlineNext()
			sys.runCpu(80)
			sys.mapper.hSync(sys.scanline)
			sys.runCpu(176)
			ppu.scanlineStart()
			sys.runCpu(80 + sys.tvFormat.nScanlineEndCycle)
		}
	}

	for ; sys.scanline <= nScanline; sys.scanline++ {
		switch sys.scanline {
		case 240:
			sys.mapper.vSync()
			sys.pad.vSync()
		case 241:
			ppu.reg2 |= ppuReg2VBlank
			if ppu.reg0&ppuReg0VBlank != 0 {
				sys.cpu.intr |= cpuIntrTypNmi
			}
		case nScanline:
			ppu.reg2 &^= ppuReg2VBlank | ppuReg2SpHit
		}
		if sys.renderMode == RenderModePreAll || sys.renderMode == RenderModePostAll {
			sys.runCpu(sys.tvFormat.nScanlineCycle)
			sys.mapper.hSync(sys.scanline)
		} else {
			sys.runCpu(sys.tvFormat.nHDrawCycle)
			sys.mapper.hSync(sys.scanline)
			sys.runCpu(sys.tvFormat.nHBlankCycle)
		}
	}

	sys.apu.render()
}
