package core

import (
	"encoding/binary"
	"errors"
	"io"
)

type NesFileHeader struct {
	Magic     uint32
	NPromPage byte
	NVromPage byte
	Control1  byte
	Control2  byte
	Reserved  [8]byte
}

type Rom struct {
	nPromPage byte
	nVromPage byte
	bVMirror  bool
	bSaveRam  bool
	bTrainer  bool
	b4Screen  bool
	mapperNo  byte
	prom      []byte
	vrom      []byte
	trn       []byte
}

func newRom(file io.Reader) (*Rom, error) {
	header := &NesFileHeader{}
	if err := binary.Read(file, binary.LittleEndian, header); err != nil {
		return nil, err
	}
	if header.Magic != 0x1a53454e { // "\x1aNES"
		return nil, errors.New("unsupported file type")
	}

	rom := &Rom{}
	rom.nPromPage = header.NPromPage
	rom.nVromPage = header.NVromPage
	rom.bVMirror = header.Control1&0x01 != 0
	rom.bSaveRam = header.Control1&0x02 != 0
	rom.bTrainer = header.Control1&0x04 != 0
	rom.b4Screen = header.Control1&0x08 != 0
	rom.mapperNo = (header.Control1 >> 4) | (header.Control2 & 0xf0)
	println("mapper: ", rom.mapperNo) //ldeng7

	if rom.bTrainer {
		rom.trn = make([]byte, 512)
		if _, err := io.ReadFull(file, rom.trn); err != nil {
			return nil, err
		}
	}
	rom.prom = make([]byte, uint32(rom.nPromPage)*0x4000)
	if _, err := io.ReadFull(file, rom.prom); err != nil {
		return nil, err
	}
	if rom.nVromPage != 0 {
		rom.vrom = make([]byte, uint32(rom.nVromPage)*0x2000)
		if _, err := io.ReadFull(file, rom.vrom); err != nil {
			return nil, err
		}
	}

	return rom, nil
}
