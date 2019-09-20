package core

const (
	ScreenWidth    = 256 + 16
	ScreenHeight   = 240
	FrameBufferLen = ScreenWidth * ScreenWidth
)

type FrameBuffer [FrameBufferLen]uint32

const (
	ppuReg0VBlank  byte = 0x80
	ppuReg0SpHit   byte = 0x40
	ppuReg0Sp16    byte = 0x20
	ppuReg0BgTbl   byte = 0x10
	ppuReg0SpTbl   byte = 0x08
	ppuReg0Inc32   byte = 0x04
	ppuReg0NameTbl byte = 0x03

	ppuReg1SpDisp    byte = 0x10
	ppuReg1BgDisp    byte = 0x08
	ppuReg1SpClip    byte = 0x04
	ppuReg1BgClip    byte = 0x02
	ppuReg1ColorMode byte = 0x01

	ppuReg2VBlank byte = 0x80
	ppuReg2SpHit  byte = 0x40
	ppuReg2SpMax  byte = 0x20

	ppuSpAttrVMirror  byte = 0x80
	ppuSpAttrHMirror  byte = 0x40
	ppuSpAttrPriority byte = 0x20
	ppuSpAttrColor    byte = 0x03
)

type Ppu struct {
	sys *Sys

	reg0    byte
	reg1    byte
	reg2    byte
	reg3    byte
	readBuf byte
	loopyT  uint16
	loopyV  uint16
	loopyX  uint16
	loopyY  uint16
	loopySh uint16
	bgPal   [16]byte
	spPal   [16]byte
	spram   [256]byte

	toggle        bool
	bExtLatch     bool
	bChrLatch     bool
	iScanline     uint16
	screen        *FrameBuffer
	palette       *[64]uint32
	spMirrorTable [256]byte
}

func newPpu(sys *Sys) *Ppu {
	ppu := &Ppu{}
	ppu.sys = sys

	for i := uint16(0); i < 256; i++ {
		var m, c byte = 0x80, 0
		for j := uint16(0); j < 8; j++ {
			if i&(1<<j) != 0 {
				c |= m
			}
			m >>= 1
		}
		ppu.spMirrorTable[i] = c
	}

	return ppu
}

func (ppu *Ppu) reset(clear bool) {
	if clear {
		l := len(ppu.bgPal)
		for i := 0; i < l; i++ {
			ppu.bgPal[i] = 0
		}
		l = len(ppu.spPal)
		for i := 0; i < l; i++ {
			ppu.spPal[i] = 0
		}
		l = len(ppu.spram)
		for i := 0; i < l; i++ {
			ppu.spram[i] = 0
		}
	}

	ppu.reg0, ppu.reg1, ppu.reg2, ppu.reg3, ppu.readBuf = 0, 0, 0, 0, 0xff
	ppu.loopyT, ppu.loopyV, ppu.loopyX, ppu.loopyY, ppu.loopySh = 0, 0, 0, 0, 0

	ppu.toggle, ppu.bExtLatch, ppu.bChrLatch = false, false, false
	ppu.iScanline = 0
	ppu.palette = &ppuPalette[0]
}

func (ppu *Ppu) read(addr uint16) byte {
	var data byte
	mem := ppu.sys.mem
	switch addr {
	case 0x2000, 0x2001, 0x2003, 0x2005, 0x2006:
		return ppu.readBuf
	case 0x2002:
		data = ppu.reg2
		ppu.toggle = false
		ppu.reg2 &^= ppuReg2VBlank
		return data
	case 0x2004:
		data = ppu.spram[ppu.reg3]
		ppu.reg3++
		return data
	case 0x2007:
		addr := ppu.loopyV & 0x3fff
		if ppu.reg0&ppuReg0Inc32 != 0 {
			ppu.loopyV += 32
		} else {
			ppu.loopyV++
		}
		if addr >= 0x3f00 {
			if addr&0x0010 == 0 {
				return ppu.bgPal[addr&0x000f]
			} else {
				return ppu.spPal[addr&0x000f]
			}
		}
		if addr >= 0x3000 {
			addr &= 0xefff
		}
		data = ppu.readBuf
		ppu.readBuf = mem.ppuBanks[addr>>10][addr&0x03ff]
		return data
	}
	return data
}

func (ppu *Ppu) write(addr uint16, data byte) {
	mem := ppu.sys.mem
	switch addr {
	case 0x2000:
		ppu.loopyT = ppu.loopyT&0xf3ff | ((uint16(data) & 0x0003) << 10)
		if (data&0x80 != 0) && (ppu.reg0&ppuReg0VBlank == 0) && (ppu.reg2&ppuReg2VBlank != 0) {
			ppu.sys.cpu.intr |= cpuIntrTypNmi
		}
		ppu.reg0 = data
	case 0x2001:
		ppu.reg1 = data
	case 0x2003:
		ppu.reg3 = data
	case 0x2004:
		ppu.spram[ppu.reg3] = data
		ppu.reg3++
	case 0x2005:
		if !ppu.toggle {
			ppu.loopyT = (ppu.loopyT & 0xffe0) | (uint16(data) >> 3)
			ppu.loopyX = uint16(data & 0x07)
		} else {
			ppu.loopyT = (ppu.loopyT & 0xfc1f) | (uint16(data&0xf8) << 2)
			ppu.loopyT = (ppu.loopyT & 0x8fff) | (uint16(data&0x07) << 12)
		}
		ppu.toggle = !ppu.toggle
	case 0x2006:
		if !ppu.toggle {
			ppu.loopyT = (ppu.loopyT & 0x00ff) | (uint16(data&0x3f) << 8)
		} else {
			ppu.loopyT = (ppu.loopyT & 0xff00) | uint16(data)
			ppu.loopyV = ppu.loopyT
			ppu.sys.mapper.ppuLatch(ppu.loopyV)
		}
		ppu.toggle = !ppu.toggle
	case 0x2007:
		vaddr := ppu.loopyV & 0x3fff
		if ppu.reg0&ppuReg0Inc32 != 0 {
			ppu.loopyV += 32
		} else {
			ppu.loopyV++
		}
		if vaddr >= 0x3f00 {
			data &= 0x3F
			if vaddr&0x000f == 0 {
				ppu.bgPal[0], ppu.spPal[0] = data, data
			} else if vaddr&0x0010 == 0 {
				ppu.bgPal[vaddr&0x000f] = data
			} else {
				ppu.spPal[vaddr&0x000f] = data
			}
			b := ppu.bgPal[0x00]
			ppu.bgPal[0x04], ppu.bgPal[0x08], ppu.bgPal[0x0c] = b, b, b
			ppu.spPal[0x00], ppu.spPal[0x04], ppu.spPal[0x08], ppu.spPal[0x0c] = b, b, b, b
			break
		}
		if vaddr >= 0x3000 {
			vaddr &= 0xefff
		}
		if mem.ppuBanksTyp[vaddr>>10] != memBankTypVrom {
			mem.ppuBanks[vaddr>>10][vaddr&0x03ff] = data
		}
	}
}

func (ppu *Ppu) dma(data byte) {
	addr := uint16(data) << 8
	sys, spram := ppu.sys, ppu.spram[:]
	for i := uint16(0); i < 256; i++ {
		spram[i] = sys.read(addr + i)
	}
}

func (ppu *Ppu) frameStart() {
	if ppu.reg1&(ppuReg1SpDisp|ppuReg1BgDisp) != 0 {
		ppu.loopyV, ppu.loopySh = ppu.loopyT, ppu.loopyX
		ppu.loopyY = (ppu.loopyV & 0x7000) >> 12
	}
	p := ppu.screen
	for i := 0; i < ScreenWidth; i++ {
		(*p)[i] = 0xff000000
	}
}

func (ppu *Ppu) scanlineStart() {
	if ppu.reg1&(ppuReg1BgDisp|ppuReg1SpDisp) != 0 {
		ppu.loopyV = (ppu.loopyV & 0xfbe0) | (ppu.loopyT & 0x041f)
		ppu.loopySh = ppu.loopyX
		ppu.loopyY = (ppu.loopyV & 0x7000) >> 12
		ppu.sys.mapper.ppuLatch((ppu.loopyV & 0x0fff) | 02000)
	}
}

func (ppu *Ppu) scanlineNext() {
	if ppu.reg1&(ppuReg1BgDisp|ppuReg1SpDisp) != 0 {
		if ppu.loopyV&0x7000 == 0x7000 {
			ppu.loopyV &= 0x8fff
			if ppu.loopyV&0x03e0 == 0x03a0 {
				ppu.loopyV ^= 0x0800
				ppu.loopyV &= 0xfc1f
			} else {
				if ppu.loopyV&0x03e0 == 0x03e0 {
					ppu.loopyV &= 0xfc1f
				} else {
					ppu.loopyV += 0x0020
				}
			}
		} else {
			ppu.loopyV += 0x1000
		}
		ppu.loopyY = (ppu.loopyV & 0x7000) >> 12
	}
}

func (ppu *Ppu) isSprite0(scanline uint16) bool {
	if ppu.reg1&(ppuReg1SpDisp|ppuReg1BgDisp) != ppuReg1SpDisp|ppuReg1BgDisp {
		return false
	}
	if ppu.reg2&ppuReg2SpHit != 0 {
		return false
	}
	sp := uint16(ppu.spram[0])
	if ppu.reg0&ppuReg0Sp16 == 0 {
		if scanline < sp+1 || scanline > sp+8 {
			return false
		}
	} else {
		if scanline < sp+1 || scanline > sp+16 {
			return false
		}
	}
	return true
}

func (ppu *Ppu) renderBgPal(attr byte, chL byte, chH byte, sl []uint32) {
	slPal := ppu.bgPal[attr:]
	c1 := ((chL >> 1) & 0x55) | (chH & 0xaa)
	c2 := (chL & 0x55) | ((chH << 1) & 0xaa)
	pal := ppu.palette
	sl[0] = (*pal)[slPal[c1>>6]]
	sl[1] = (*pal)[slPal[c2>>6]]
	sl[2] = (*pal)[slPal[(c1>>4)&0x03]]
	sl[3] = (*pal)[slPal[(c2>>4)&0x03]]
	sl[4] = (*pal)[slPal[(c1>>2)&0x03]]
	sl[5] = (*pal)[slPal[(c2>>2)&0x03]]
	sl[6] = (*pal)[slPal[c1&0x03]]
	sl[7] = (*pal)[slPal[c2&0x03]]
}

func (ppu *Ppu) scanlineRender(scanline uint8, bAllSp bool) {

	sys, mem := ppu.sys, ppu.sys.mem
	bgs := [34]byte{}

	if scanline == 1 {
		iPal := (ppu.reg1 >> 5) | ((ppu.reg1 & ppuReg1ColorMode) << 3)
		ppu.palette = &ppuPalette[iPal]
	}

	if ppu.reg1&ppuReg1BgDisp == 0 {
		p, c := ppu.screen, (*ppu.palette)[ppu.bgPal[0]]
		for i, ie := ppu.iScanline, ppu.iScanline+ScreenWidth; i < ie; i++ {
			(*p)[i] = c
		}
		if sys.renderMode == RenderModeTile {
			sys.runCpu(1024)
		}
	} else {
		sl := (*ppu.screen)[ppu.iScanline-ppu.loopySh:]
		iNameTbl := (ppu.loopyV & 0x0fff) | 0x2000
		bTileMode := sys.renderMode == RenderModeTile
		var prevTile uint16
		var prevAttr byte
		if !ppu.bExtLatch {
			iAttr := ((ppu.loopyV & 0x0380) >> 4) | 0x03c0
			x := byte(iNameTbl) & 0x1f
			attrSh := (byte(iNameTbl) & 0x40) >> 4
			bank := mem.ppuBanks[iNameTbl>>10]
			for i := byte(0); i < 33; i++ {
				tile := (uint16(ppu.reg0&ppuReg0BgTbl) << 8) +
					(uint16(bank[iNameTbl&0x03ff]) << 4) + ppu.loopyY
				if i != 0 && bTileMode {
					sys.runCpu(32)
				}
				attr := ((bank[iAttr+uint16(x>>2)] >> ((x & 0x02) | attrSh)) & 0x03) << 2
				if i != 0 && prevTile == tile && prevAttr == attr {
					copy(sl[8:16], sl[0:8])
					bgs[i] = bgs[i-1]
				} else {
					prevTile, prevAttr = tile, attr
					bank := mem.ppuBanks[tile>>10]
					chL, chH := bank[tile&0x03ff], bank[(tile&0x03ff)+8]
					bgs[i] = chH | chL
					ppu.renderBgPal(attr, chL, chH, sl[8:])
				}
				sl = sl[8:]
				if ppu.bChrLatch {
					sys.mapper.ppuChrLatch(tile)
				}
				x++
				if x == 32 {
					x = 0
					iNameTbl ^= 0x041f
					iAttr = ((iNameTbl & 0x0380) >> 4) | 0x03c0
					bank = mem.ppuBanks[iNameTbl>>10]
				} else {
					iNameTbl++
				}
			}
		} else {
			x := iNameTbl & 0x1f
			var chH, chL, exattr byte
			for i := byte(0); i < 33; i++ {
				if i != 0 && bTileMode {
					sys.runCpu(32)
				}
				sys.mapper.ppuExtLatchX(i)
				sys.mapper.ppuExtLatch(iNameTbl, &chL, &chH, &exattr)
				attr := exattr & 0x0c
				tile := (uint16(chH) << 8) | uint16(chL)
				if i != 0 && prevTile == tile && prevAttr == attr {
					copy(sl[8:16], sl[0:8])
					bgs[i] = bgs[i-1]
				} else {
					prevTile, prevAttr = tile, attr
					bgs[i] = chH | chL
					ppu.renderBgPal(attr, chL, chH, sl[8:])
				}
				sl = sl[8:]
				x++
				if x == 32 {
					x = 0
					iNameTbl ^= 0x041f
				} else {
					iNameTbl++
				}
			}
		}

		if ppu.reg1&ppuReg1BgClip == 0 {
			p, c := ppu.screen, (*ppu.palette)[ppu.bgPal[0]]
			for i, ie := ppu.iScanline+8, ppu.iScanline+16; i < ie; i++ {
				(*p)[i] = c
			}
		}
	}

	ppu.reg2 &^= ppuReg2SpMax
	if scanline > 239 || ppu.reg1&ppuReg1SpDisp == 0 {
		return
	}
	sps, spram := [34]byte{}, ppu.spram[:]
	nSp, spM := 0, byte(7)
	if ppu.reg0&ppuReg0Sp16 != 0 {
		spM = 15
	}
	if ppu.reg1&ppuReg1SpClip == 0 {
		sps[0] = 0xff
	}

	for i, j := 0, 0; i < 64; i, j = i+1, j+4 {
		spY, spTile, spAttr, spX := spram[j], spram[j+1], spram[j+2], spram[j+3]
		spDy := scanline - spY - 1
		if spY >= scanline || spDy > spM {
			continue
		}
		var spAddr uint16
		if ppu.reg0&ppuReg0Sp16 == 0 {
			spAddr = (uint16(ppu.reg0&ppuReg0SpTbl) << 9) | (uint16(spTile) << 4)
			if spAttr&ppuSpAttrVMirror == 0 {
				spAddr += uint16(spDy)
			} else {
				spAddr += uint16(7 - spDy)
			}
		} else {
			spAddr = (uint16(spTile&0x01) << 12) | (uint16(spTile&0xfe) << 4)
			if spAttr&ppuSpAttrVMirror == 0 {
				spAddr += uint16(((spDy & 0x08) << 1) | (spDy & 0x07))
			} else {
				spAddr += uint16(((^spDy & 0x08) << 1) | (7 - (spDy & 0x07)))
			}
		}

		bank := mem.ppuBanks[spAddr>>10]
		chL, chH := bank[spAddr&0x03ff], bank[(spAddr&0x03ff)+8]
		if ppu.bChrLatch {
			sys.mapper.ppuChrLatch(spAddr)
		}
		if spAttr&ppuSpAttrHMirror != 0 {
			chL, chH = ppu.spMirrorTable[chL], ppu.spMirrorTable[chH]
		}
		spPat := chL | chH

		var maskBg byte
		if i == 0 && ppu.reg2&ppuReg2SpHit == 0 {
			p16 := uint16(spX&0xf8) + ((ppu.loopySh + uint16(spX&0x7)) & 0x08)
			pos := byte(p16 >> 3)
			sh := 8 - byte((ppu.loopySh+uint16(spX))&0x07)
			m16 := (uint16(bgs[pos]) << 8) | uint16(bgs[pos+1])
			maskBg = byte(m16 >> sh)
			if spPat&maskBg != 0 {
				ppu.reg2 |= ppuReg2SpHit
			}
		}
		{
			pos, sh := spX>>3, 8-(spX&0x07)
			m16 := (uint16(sps[pos]) << 8) | uint16(sps[pos+1])
			maskSp := byte(m16 >> sh)
			spPat &^= maskSp
			wrt := uint16(spPat) << sh
			sps[pos] |= byte(wrt >> 8)
			sps[pos+1] |= byte(wrt)
		}
		if spAttr&ppuSpAttrPriority != 0 {
			spPat &^= maskBg
		}

		slPal := ppu.spPal[((spAttr & ppuSpAttrColor) << 2):]
		sl := (*ppu.screen)[ppu.iScanline+uint16(spX)+8:]
		c1 := ((chL >> 1) & 0x55) | (chH & 0xaa)
		c2 := (chL & 0x55) | ((chH << 1) & 0xaa)
		pal := ppu.palette
		for i := byte(0); i < 8; i += 2 {
			j := 6 - i
			if spPat&(0x02<<j) != 0 {
				sl[i] = (*pal)[slPal[(c1>>j)&0x03]]
			}
			if spPat&(0x01<<j) != 0 {
				sl[i+1] = (*pal)[slPal[(c2>>j)&0x03]]
			}
		}

		nSp++
		if nSp > 7 {
			ppu.reg2 |= ppuReg2SpMax
			if !bAllSp {
				break
			}
		}
	}
}
