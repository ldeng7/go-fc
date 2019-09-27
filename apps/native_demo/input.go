package main

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/ldeng7/go-fc/core"
)

var keyMapP1 = map[glfw.Key]byte{
	glfw.KeyW:           core.PadKeyUp,
	glfw.KeyS:           core.PadKeyDown,
	glfw.KeyA:           core.PadKeyLeft,
	glfw.KeyD:           core.PadKeyRight,
	glfw.KeySpace:       core.PadKeyStart,
	glfw.KeyLeftControl: core.PadKeySelect,
	glfw.KeyK:           core.PadKeyA,
	glfw.KeyJ:           core.PadKeyB,
}
var keyMapP2 = map[glfw.Key]byte{
	glfw.KeyUp:        core.PadKeyUp,
	glfw.KeyDown:      core.PadKeyDown,
	glfw.KeyLeft:      core.PadKeyLeft,
	glfw.KeyRight:     core.PadKeyRight,
	glfw.KeyKPEnter:   core.PadKeyStart,
	glfw.KeyKPDecimal: core.PadKeySelect,
	glfw.KeyKP2:       core.PadKeyA,
	glfw.KeyKP1:       core.PadKeyB,
}

func (a *App) onKey(_ *glfw.Window, key glfw.Key, _ int, action glfw.Action, _ glfw.ModifierKey) {
	switch key {
	case glfw.KeyEscape:
		if action == glfw.Press {
			a.sys.Reset()
		}
	default:
		switch action {
		case glfw.Press:
			if pk, ok := keyMapP1[key]; ok {
				a.sys.SetPadKey(1, pk, true)
			} else if pk, ok := keyMapP2[key]; ok {
				a.sys.SetPadKey(2, pk, true)
			}
		case glfw.Release:
			if pk, ok := keyMapP1[key]; ok {
				a.sys.SetPadKey(1, pk, false)
			} else if pk, ok := keyMapP2[key]; ok {
				a.sys.SetPadKey(2, pk, false)
			}
		}
	}
}
