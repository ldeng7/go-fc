package core

const (
	PadKeyA byte = 1 << iota
	PadKeyB
	PadKeySelect
	PadKeyStart
	PadKeyUp
	PadKeyDown
	PadKeyLeft
	PadKeyRight
)

type Pad struct {
	bStrobe  bool
	b1p, b2p byte
	b1, b2   byte
	b1u, b2u byte
}

func newPad() *Pad {
	return &Pad{}
}

func (pad *Pad) reset() {
	pad.bStrobe = false
	pad.b1p, pad.b2p = 0, 0
	pad.b1, pad.b2 = 0, 0
	pad.b1u, pad.b2u = 0, 0
}

func (pad *Pad) read(addr uint16) byte {
	var b byte
	switch addr {
	case 0x4016:
		b = pad.b1u & 0x01
		pad.b1u >>= 1
	case 0x4017:
		b = pad.b2u & 0x01
		pad.b2u >>= 1
	}
	return b
}

func (pad *Pad) write(addr uint16, data byte) {
	if addr == 0x4016 {
		if data&0x01 != 0 {
			pad.bStrobe = true
		} else if pad.bStrobe {
			pad.bStrobe = false
			pad.b1u, pad.b2u = pad.b1, pad.b2
		}
	}
}

func (pad *Pad) vSync() {
	pad.b1, pad.b2 = pad.b1p, pad.b2p
}

func (pad *Pad) setKey(p byte, k byte, down bool) {
	switch p {
	case 1:
		if down {
			pad.b1p |= k
		} else {
			pad.b1p &^= k
		}
	case 2:
		if down {
			pad.b2p |= k
		} else {
			pad.b2p &^= k
		}
	}
}
