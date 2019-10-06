package core

const (
	memBankTypRom byte = iota
	memBankTypRam
	memBankTypVrom
	memBankTypCram
	memBankTypVram
)
const (
	memVramMirrorH byte = iota
	memVramMirrorV
	memVramMirror4
	memVramMirror4L
	memVramMirror4H
)

type Mem struct {
	sys *Sys

	nProm8kPage uint32
	nVrom1kPage uint32
	cpuBanks    [8][]byte
	cpuBanksTyp [8]byte
	ppuBanks    [12][]byte
	ppuBanksTyp [12]byte

	ram    [0x2000]byte
	xram   [0x2000]byte
	dram   [0xa000]byte
	wram   [0x20000]byte
	vram   [0x1000]byte
	cram   [0x8000]byte
	prom   []byte
	vrom   []byte
	cpuReg [24]byte
}

func newMem(sys *Sys) *Mem {
	mem := &Mem{}
	mem.sys = sys

	mem.setCpuBank(0, mem.ram[:], memBankTypRam)
	mem.setCpuBank(1, mem.xram[:], memBankTypRom)
	mem.setCpuBank(2, mem.xram[:], memBankTypRom)
	mem.setCpuBank(3, mem.wram[:], memBankTypRam)

	rom := sys.rom
	mem.prom = rom.prom
	mem.vrom = rom.vrom

	mem.nProm8kPage = uint32(rom.nPromPage) << 1
	mem.nVrom1kPage = uint32(rom.nVromPage) << 3
	if mem.nVrom1kPage != 0 {
		mem.setVrom8kBank(0)
	} else {
		mem.setCram8kBank(0)
	}
	if rom.b4Screen {
		mem.setVramMirror(memVramMirror4)
	} else if rom.bVMirror {
		mem.setVramMirror(memVramMirrorV)
	} else {
		mem.setVramMirror(memVramMirrorH)
	}

	return mem
}

func (mem *Mem) reset(init bool) {
	if !mem.sys.rom.bSaveRam {
		l := len(mem.wram)
		for i := 0; i < l; i++ {
			mem.wram[i] = 0xff
		}
	}
	if !init {
		l := len(mem.ram)
		for i := 0; i < l; i++ {
			mem.ram[i] = 0
		}
		l = len(mem.vram)
		for i := 0; i < l; i++ {
			mem.vram[i] = 0
		}
		l = len(mem.cram)
		for i := 0; i < l; i++ {
			mem.cram[i] = 0
		}
		l = len(mem.cpuReg)
		for i := 0; i < l; i++ {
			mem.cpuReg[i] = 0
		}
	}
	if mem.sys.rom.bTrainer {
		copy(mem.wram[0x1000:0x1200], mem.sys.rom.trn)
	}
}

func (mem *Mem) setCpuBank(iBank byte, slice []byte, typ byte) {
	mem.cpuBanks[iBank], mem.cpuBanksTyp[iBank] = slice[:0x2000:0x2000], typ
}

func (mem *Mem) setProm8kBank(iBank byte, iPage uint32) {
	iPage %= mem.nProm8kPage
	i := iPage << 13
	mem.cpuBanks[iBank], mem.cpuBanksTyp[iBank] = mem.prom[i:i+0x2000:i+0x2000], memBankTypRom
}

func (mem *Mem) setProm16kBank(iBank byte, iPage uint32) {
	iPage <<= 1
	mem.setProm8kBank(iBank, iPage)
	mem.setProm8kBank(iBank+1, iPage+1)
}

func (mem *Mem) setProm32kBank(iPage uint32) {
	iPage <<= 2
	mem.setProm8kBank(4, iPage)
	mem.setProm8kBank(5, iPage+1)
	mem.setProm8kBank(6, iPage+2)
	mem.setProm8kBank(7, iPage+3)
}

func (mem *Mem) setProm32kBank4(iPage0, iPage1, iPage2, iPage3 uint32) {
	mem.setProm8kBank(4, iPage0)
	mem.setProm8kBank(5, iPage1)
	mem.setProm8kBank(6, iPage2)
	mem.setProm8kBank(7, iPage3)
}

func (mem *Mem) setVrom1kBank(iBank byte, iPage uint32) {
	iPage %= mem.nVrom1kPage
	i := iPage << 10
	mem.ppuBanks[iBank], mem.ppuBanksTyp[iBank] = mem.vrom[i:i+0x0400:i+0x0400], memBankTypVrom
}

func (mem *Mem) setVrom2kBank(iBank byte, iPage uint32) {
	iPage <<= 1
	mem.setVrom1kBank(iBank, iPage)
	mem.setVrom1kBank(iBank+1, iPage+1)
}

func (mem *Mem) setVrom4kBank(iBank byte, iPage uint32) {
	iPage <<= 2
	mem.setVrom1kBank(iBank, iPage)
	mem.setVrom1kBank(iBank+1, iPage+1)
	mem.setVrom1kBank(iBank+2, iPage+2)
	mem.setVrom1kBank(iBank+3, iPage+3)
}

func (mem *Mem) setVrom8kBank(iPage uint32) {
	iPage <<= 3
	for i := byte(0); i < 8; i++ {
		mem.setVrom1kBank(i, iPage+uint32(i))
	}
}

func (mem *Mem) setVrom8kBank8(iPage0, iPage1, iPage2, iPage3, iPage4, iPage5, iPage6, iPage7 uint32) {
	mem.setVrom1kBank(0, iPage0)
	mem.setVrom1kBank(1, iPage1)
	mem.setVrom1kBank(2, iPage2)
	mem.setVrom1kBank(3, iPage3)
	mem.setVrom1kBank(4, iPage4)
	mem.setVrom1kBank(5, iPage5)
	mem.setVrom1kBank(6, iPage6)
	mem.setVrom1kBank(7, iPage7)
}

func (mem *Mem) setCram1kBank(iBank byte, iPage uint32) {
	iPage &= 0x1f
	i := iPage << 10
	mem.ppuBanks[iBank], mem.ppuBanksTyp[iBank] = mem.cram[i:i+0x0400:i+0x0400], memBankTypCram
}

func (mem *Mem) setCram2kBank(iBank byte, iPage uint32) {
	iPage <<= 1
	mem.setCram1kBank(iBank, iPage)
	mem.setCram1kBank(iBank+1, iPage+1)
}

func (mem *Mem) setCram4kBank(iBank byte, iPage uint32) {
	iPage <<= 2
	mem.setCram1kBank(iBank, iPage)
	mem.setCram1kBank(iBank+1, iPage+1)
	mem.setCram1kBank(iBank+2, iPage+2)
	mem.setCram1kBank(iBank+3, iPage+3)
}

func (mem *Mem) setCram8kBank(iPage uint32) {
	iPage <<= 3
	for i := byte(0); i < 8; i++ {
		mem.setCram1kBank(i, iPage+uint32(i))
	}
}

func (mem *Mem) setCram8kBank8(iPage0, iPage1, iPage2, iPage3, iPage4, iPage5, iPage6, iPage7 uint32) {
	mem.setCram1kBank(0, iPage0)
	mem.setCram1kBank(1, iPage1)
	mem.setCram1kBank(2, iPage2)
	mem.setCram1kBank(3, iPage3)
	mem.setCram1kBank(4, iPage4)
	mem.setCram1kBank(5, iPage5)
	mem.setCram1kBank(6, iPage6)
	mem.setCram1kBank(7, iPage7)
}

func (mem *Mem) setVram1kBank(iBank byte, iPage uint32) {
	iPage &= 3
	i := iPage << 10
	mem.ppuBanks[iBank], mem.ppuBanksTyp[iBank] = mem.vram[i:i+0x0400:i+0x0400], memBankTypVram
}

func (mem *Mem) setVramBank(iPage0, iPage1, iPage2, iPage3 uint32) {
	mem.setVram1kBank(8, iPage0)
	mem.setVram1kBank(9, iPage1)
	mem.setVram1kBank(10, iPage2)
	mem.setVram1kBank(11, iPage3)
}

func (mem *Mem) setVramMirror(typ byte) {
	switch typ {
	case memVramMirrorH:
		mem.setVramBank(0, 0, 1, 1)
	case memVramMirrorV:
		mem.setVramBank(0, 1, 0, 1)
	case memVramMirror4L:
		mem.setVramBank(0, 0, 0, 0)
	case memVramMirror4H:
		mem.setVramBank(1, 1, 1, 1)
	case memVramMirror4:
		mem.setVramBank(0, 1, 2, 3)
	}
}
