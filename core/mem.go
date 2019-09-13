package core

const (
	memBankTypRom byte = iota
	memBankTypRam
	memBankTypDram
	memBankTypMapper
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

	cpuBanks    [8][]byte
	cpuBanksTyp [8]byte
	ppuBanks    [12][]byte
	ppuBanksTyp [12]byte

	ram    [8 * 1024]byte
	xram   [8 * 1024]byte
	eram   [32 * 1024]byte
	dram   [40 * 1024]byte
	wram   [128 * 1024]byte
	vram   [4 * 1024]byte
	cram   [32 * 1024]byte
	bgPal  [16]byte
	spPal  [16]byte
	spram  [256]byte
	cpuReg [24]byte
	prom   []byte
	vrom   []byte

	prom8kSize  uint32
	prom16kSize uint32
	prom32kSize uint32
	vrom1kSize  uint32
	vrom2kSize  uint32
	vrom4kSize  uint32
	vrom8kSize  uint32
}

func newMem(sys *Sys) *Mem {
	mem := &Mem{}
	mem.sys = sys

	mem.setPromBank(0, mem.ram[:0x2000], memBankTypRam)
	mem.setPromBank(1, mem.xram[:0x2000], memBankTypRom)
	mem.setPromBank(2, mem.xram[:0x2000], memBankTypRom)
	mem.setPromBank(3, mem.wram[:0x2000], memBankTypRom)

	rom := sys.rom
	mem.prom = rom.prom
	mem.vrom = rom.vrom
	if rom.bTrainer {
		copy(mem.wram[0x1000:0x1200], rom.trn)
	}

	mem.prom8kSize = uint32(rom.nPromBanks << 1)
	mem.prom16kSize = uint32(rom.nPromBanks)
	mem.prom32kSize = uint32(rom.nPromBanks >> 1)
	mem.vrom1kSize = uint32(rom.nVromBanks << 3)
	mem.vrom2kSize = uint32(rom.nVromBanks << 2)
	mem.vrom4kSize = uint32(rom.nVromBanks << 1)
	mem.vrom8kSize = uint32(rom.nVromBanks)

	if mem.vrom8kSize != 0 {
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

func (mem *Mem) setPromBank(iBank byte, slice []byte, typ byte) {
	mem.cpuBanks[iBank], mem.cpuBanksTyp[iBank] = slice, typ
}

func (mem *Mem) setProm8kBank(iBank byte, iPage uint32) {
	iPage %= mem.prom8kSize
	i := iPage << 13
	mem.cpuBanks[iBank], mem.cpuBanksTyp[iBank] = mem.prom[i:i+0x2000], memBankTypRom
}

func (mem *Mem) setProm16kBank(iBank byte, iPage uint32) {
	iPage *= 2
	mem.setProm8kBank(iBank, iPage)
	mem.setProm8kBank(iBank+1, iPage+1)
}

func (mem *Mem) setProm32kBank(iPage uint32) {
	iPage *= 4
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
	iPage %= mem.vrom1kSize
	i := iPage << 10
	mem.ppuBanks[iBank], mem.ppuBanksTyp[iBank] = mem.vrom[i:i+0x0400], memBankTypVrom
}

func (mem *Mem) setVrom2kBank(iBank byte, iPage uint32) {
	iPage *= 2
	mem.setVrom1kBank(iBank, iPage)
	mem.setVrom1kBank(iBank+1, iPage+1)
}

func (mem *Mem) setVrom4kBank(iBank byte, iPage uint32) {
	iPage *= 4
	mem.setVrom1kBank(iBank, iPage)
	mem.setVrom1kBank(iBank+1, iPage+1)
	mem.setVrom1kBank(iBank+2, iPage+2)
	mem.setVrom1kBank(iBank+3, iPage+3)
}

func (mem *Mem) setVrom8kBank(iPage uint32) {
	iPage *= 8
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
	mem.ppuBanks[iBank], mem.ppuBanksTyp[iBank] = mem.cram[i:i+0x0400], memBankTypCram
}

func (mem *Mem) setCram2kBank(iBank byte, iPage uint32) {
	iPage *= 2
	mem.setCram1kBank(iBank, iPage)
	mem.setCram1kBank(iBank+1, iPage+1)
}

func (mem *Mem) setCram4kBank(iBank byte, iPage uint32) {
	iPage *= 4
	mem.setCram1kBank(iBank, iPage)
	mem.setCram1kBank(iBank+1, iPage+1)
	mem.setCram1kBank(iBank+2, iPage+2)
	mem.setCram1kBank(iBank+3, iPage+3)
}

func (mem *Mem) setCram8kBank(iPage uint32) {
	iPage *= 8
	for i := byte(0); i < 8; i++ {
		mem.setCram1kBank(i, iPage+uint32(i))
	}
}

func (mem *Mem) setVram1kBank(iBank byte, iPage uint32) {
	iPage &= 3
	i := iPage << 10
	mem.ppuBanks[iBank], mem.ppuBanksTyp[iBank] = mem.vram[i:i+0x0400], memBankTypVram
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

func (mem *Mem) setVramMirror4(iPage0, iPage1, iPage2, iPage3 uint32) {
	mem.setVram1kBank(8, iPage0)
	mem.setVram1kBank(9, iPage1)
	mem.setVram1kBank(10, iPage2)
	mem.setVram1kBank(11, iPage3)
}
